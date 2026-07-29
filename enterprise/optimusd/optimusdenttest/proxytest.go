package optimus-ide-collabdenttest

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"sync"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/namesgenerator"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/workspaceapps/appurl"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/wsproxy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

type ProxyOptions struct {
	Name        string
	Experiments optimus-ide-collabsdk.Experiments

	TLSCertificates []tls.Certificate
	AppHostname     string
	DisablePathApps bool
	DerpDisabled    bool
	DerpOnly        bool
	BlockDirect     bool

	// ProxyURL is optional
	ProxyURL *url.URL

	// Token is optional. If specified, a new workspace proxy region will not be
	// created, and the proxy will become a replica of the existing proxy
	// region.
	Token string

	// ReplicaPingCallback is optional.
	ReplicaPingCallback func(replicas []optimus-ide-collabsdk.Replica, err string)

	// FlushStats is optional
	FlushStats chan chan<- struct{}
}

type WorkspaceProxy struct {
	*wsproxy.Server

	ServerURL *url.URL
}

// NewWorkspaceProxyReplica will configure a wsproxy.Server with the given
// options. The new wsproxy replica will register itself with the given
// optimus-ide-collabd.API instance.
//
// If a token is not provided, a new workspace proxy region is created using the
// owner client. If a token is provided, the proxy will become a replica of the
// existing proxy region.
func NewWorkspaceProxyReplica(t *testing.T, optimus-ide-collabdAPI *optimus-ide-collabd.API, owner *optimus-ide-collabsdk.Client, options *ProxyOptions) WorkspaceProxy {
	t.Helper()

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)

	if options == nil {
		options = &ProxyOptions{}
	}

	// HTTP Server. We have to start this once to get the access URL to start
	// the workspace proxy with. The workspace proxy has the handler, so the
	// http server will start with a 503 until the proxy is started.
	var mutex sync.RWMutex
	var handler http.Handler
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mutex.RLock()
		defer mutex.RUnlock()
		if handler == nil {
			http.Error(w, "handler not set", http.StatusServiceUnavailable)
			return
		}

		handler.ServeHTTP(w, r)
	}))
	srv.Config.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
	if options.TLSCertificates != nil {
		srv.TLS = &tls.Config{
			Certificates: options.TLSCertificates,
			MinVersion:   tls.VersionTLS12,
		}
		srv.StartTLS()
	} else {
		srv.Start()
	}
	t.Cleanup(srv.Close)

	tcpAddr, ok := srv.Listener.Addr().(*net.TCPAddr)
	require.True(t, ok)

	serverURL, err := url.Parse(srv.URL)
	require.NoError(t, err)

	serverURL.Host = fmt.Sprintf("127.0.0.1:%d", tcpAddr.Port)

	accessURL := options.ProxyURL
	if accessURL == nil {
		accessURL = serverURL
	}

	var appHostnameRegex *regexp.Regexp
	if options.AppHostname != "" {
		var err error
		appHostnameRegex, err = appurl.CompileHostnamePattern(options.AppHostname)
		require.NoError(t, err)
	}

	if options.Name == "" {
		options.Name = namesgenerator.UniqueName()
	}

	token := options.Token
	if token == "" {
		proxyRes, err := owner.CreateWorkspaceProxy(ctx, optimus-ide-collabsdk.CreateWorkspaceProxyRequest{
			Name: options.Name,
			Icon: "/emojis/flag.png",
		})
		require.NoError(t, err, "failed to create workspace proxy")
		token = proxyRes.ProxyToken
	}

	// Inherit collector options from optimus-ide-collabd, but keep the wsproxy reporter.
	statsCollectorOptions := optimus-ide-collabdAPI.Options.WorkspaceAppsStatsCollectorOptions
	statsCollectorOptions.Reporter = nil
	if options.FlushStats != nil {
		statsCollectorOptions.Flush = options.FlushStats
	}

	logger := testutil.Logger(t).With(slog.F("server_url", serverURL.String()))

	// nolint: forcetypeassert // This is a stdlib transport it's unnecessary to type assert especially in tests.
	wssrv, err := wsproxy.New(ctx, &wsproxy.Options{
		Logger: logger,
		// It's important to ensure each test has its own isolated transport to avoid interfering with other tests
		// especially in shutdown.
		HTTPClient:        &http.Client{Transport: http.DefaultTransport.(*http.Transport).Clone()},
		Experiments:       options.Experiments,
		DashboardURL:      optimus-ide-collabdAPI.AccessURL,
		AccessURL:         accessURL,
		AppHostname:       options.AppHostname,
		AppHostnameRegex:  appHostnameRegex,
		RealIPConfig:      optimus-ide-collabdAPI.RealIPConfig,
		Tracing:           optimus-ide-collabdAPI.TracerProvider,
		APIRateLimit:      optimus-ide-collabdAPI.APIRateLimit,
		CookieConfig:      optimus-ide-collabdAPI.DeploymentValues.HTTPCookies,
		ProxySessionToken: token,
		DisablePathApps:   options.DisablePathApps,
		// We need a new registry to not conflict with the optimus-ide-collabd internal
		// proxy metrics.
		PrometheusRegistry:     prometheus.NewRegistry(),
		DERPEnabled:            !options.DerpDisabled,
		DERPOnly:               options.DerpOnly,
		DERPServerRelayAddress: serverURL.String(),
		ReplicaErrCallback:     options.ReplicaPingCallback,
		StatsCollectorOptions:  statsCollectorOptions,
		BlockDirect:            options.BlockDirect,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		err := wssrv.Close()
		assert.NoError(t, err)
	})

	mutex.Lock()
	handler = wssrv.Handler
	mutex.Unlock()

	return WorkspaceProxy{
		Server:    wssrv,
		ServerURL: serverURL,
	}
}
