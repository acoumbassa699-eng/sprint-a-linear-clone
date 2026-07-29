//go:build !slim

package cli

import (
	"context"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/aibridge/intercept/apidump"
	agplaibridge "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/aibridge"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/aibridged"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/prometheusmetrics"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/aibridgeproxyd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd"
)

// aiBridgeProxyDaemon bundles the proxy server and its pubsub
// subscription so both are torn down by a single Close call.
type aiBridgeProxyDaemon struct {
	server      *aibridgeproxyd.Server
	unsubscribe func()
}

func (d *aiBridgeProxyDaemon) Close() error {
	if d.unsubscribe != nil {
		d.unsubscribe()
	}
	return d.server.Close()
}

// newAIBridgeProxyDaemon starts the enterprise aibridge proxy daemon,
// subscribes to ai_providers changes so the proxy's routing snapshot
// tracks the database, and registers the HTTP handler on the API.
// The returned io.Closer tears down both the subscription and server.
func newAIBridgeProxyDaemon(optimus-ide-collabAPI *optimus-ide-collabd.API) (io.Closer, error) {
	ctx := context.Background()
	optimus-ide-collabAPI.Logger.Debug(ctx, "starting in-memory aibridgeproxy daemon")

	logger := optimus-ide-collabAPI.Logger.Named("aibridgeproxyd")

	// TODO(deprecation): Remove "optimus-ide-collab_aibridgeproxyd_" in v2.37.
	// See AIGOV-447:
	// https://linear.app/optimus-ide-collabcom/issue/AIGOV-447/remove-legacy-ai-gateway-metric-aliases
	reg := prometheusmetrics.NewMetricAliasRegisterer(optimus-ide-collabAPI.PrometheusRegistry, "optimus-ide-collab_ai_gateway_proxy_", "optimus-ide-collab_aibridgeproxyd_")
	metrics := aibridgeproxyd.NewMetrics(reg)

	var newDumper func(provider, requestID string) aibridgeproxyd.RoundTripDumper
	if dumpDir := optimus-ide-collabAPI.DeploymentValues.AI.BridgeProxyConfig.APIDumpDir.String(); dumpDir != "" {
		newDumper = func(provider, requestID string) aibridgeproxyd.RoundTripDumper {
			return apidump.NewDumper(filepath.Join(dumpDir, provider, requestID), logger)
		}
	}

	target, err := resolveAIGatewayProxyTarget(
		optimus-ide-collabAPI.AccessURL,
		optimus-ide-collabAPI.DeploymentValues.AI.BridgeProxyConfig.Target.String(),
	)
	if err != nil {
		return nil, err
	}

	srv, err := aibridgeproxyd.New(ctx, logger, aibridgeproxyd.Options{
		ListenAddr:          optimus-ide-collabAPI.DeploymentValues.AI.BridgeProxyConfig.ListenAddr.String(),
		TLSCertFile:         optimus-ide-collabAPI.DeploymentValues.AI.BridgeProxyConfig.TLSCertFile.String(),
		TLSKeyFile:          optimus-ide-collabAPI.DeploymentValues.AI.BridgeProxyConfig.TLSKeyFile.String(),
		GatewayURL:          target,
		MITMCertFile:        optimus-ide-collabAPI.DeploymentValues.AI.BridgeProxyConfig.MITMCertFile.String(),
		MITMKeyFile:         optimus-ide-collabAPI.DeploymentValues.AI.BridgeProxyConfig.MITMKeyFile.String(),
		UpstreamProxy:       optimus-ide-collabAPI.DeploymentValues.AI.BridgeProxyConfig.UpstreamProxy.String(),
		UpstreamProxyCA:     optimus-ide-collabAPI.DeploymentValues.AI.BridgeProxyConfig.UpstreamProxyCA.String(),
		AllowedPrivateCIDRs: optimus-ide-collabAPI.DeploymentValues.AI.BridgeProxyConfig.AllowedPrivateCIDRs.Value(),
		NewDumper:           newDumper,
		Metrics:             metrics,
		RefreshProviders:    refreshProxyProviders(optimus-ide-collabAPI.Database),
	})
	if err != nil {
		return nil, xerrors.Errorf("failed to start in-memory aibridgeproxy daemon: %w", err)
	}

	unsubscribe, err := aibridged.SubscribeProviderReload(ctx, optimus-ide-collabAPI.Pubsub, srv, logger.Named("provider-reload"))
	if err != nil {
		// Without the subscription the proxy can never track provider changes,
		// so fail startup rather than serve a permanently stale snapshot.
		_ = srv.Close()
		return nil, xerrors.Errorf("subscribe aibridgeproxyd to ai providers change channel: %w", err)
	}

	// Register the handler so optimus-ide-collabd can serve the proxy endpoints.
	optimus-ide-collabAPI.RegisterInMemoryAIBridgeProxydHTTPHandler(srv.Handler())

	return &aiBridgeProxyDaemon{
		server:      srv,
		unsubscribe: unsubscribe,
	}, nil
}

// resolveAIGatewayProxyTarget returns the URL to which the aibridgeproxyd should forward requests.
func resolveAIGatewayProxyTarget(accessURL *url.URL, target string) (string, error) {
	if target != "" {
		return target, nil
	}

	target, err := url.JoinPath(accessURL.String(), agplaibridge.AIGatewayRootPath)
	if err != nil {
		return "", xerrors.Errorf("build embedded AI Gateway proxy target: %w", err)
	}
	return target, nil
}

// refreshProxyProviders classifies every ai_providers row as enabled,
// disabled, or error so the proxy router and any observers see the full
// configured set. Disabled rows are excluded from routing; errored rows
// are excluded from routing and surface their failure reason for
// metrics and logs.
func refreshProxyProviders(db database.Store) aibridgeproxyd.RefreshProvidersFunc {
	return func(ctx context.Context) (aibridgeproxyd.ProviderReload, error) {
		//nolint:gocritic // AsAIProviderMetadataReader is the correct subject for routing-only access.
		rows, err := db.GetAIProviders(dbauthz.AsAIProviderMetadataReader(ctx), database.GetAIProvidersParams{
			IncludeDisabled: true,
		})
		if err != nil {
			return aibridgeproxyd.ProviderReload{}, xerrors.Errorf("load ai providers: %w", err)
		}
		reload := aibridgeproxyd.ProviderReload{
			Providers: make([]aibridgeproxyd.ReloadedProvider, 0, len(rows)),
		}
		seenHost := make(map[string]string, len(rows))
		for _, row := range rows {
			reload.Providers = append(reload.Providers, classifyProviderRow(row, seenHost))
		}
		return reload, nil
	}
}

// classifyProviderRow evaluates a single ai_providers row for routing.
// seenHost is mutated to track the first provider that claimed each
// hostname so later duplicates can be flagged as errors.
func classifyProviderRow(row database.AIProvider, seenHost map[string]string) aibridgeproxyd.ReloadedProvider {
	out := aibridgeproxyd.ReloadedProvider{
		ProviderOutcome: aibridged.ProviderOutcome{
			Name: row.Name,
			Type: string(row.Type),
		},
	}
	if !row.Enabled {
		out.Status = aibridged.ProviderStatusDisabled
		return out
	}
	if strings.TrimSpace(row.BaseUrl) == "" {
		out.Status = aibridged.ProviderStatusError
		out.Err = xerrors.New("base url is empty")
		return out
	}
	u, err := url.Parse(row.BaseUrl)
	if err != nil {
		out.Status = aibridged.ProviderStatusError
		out.Err = xerrors.Errorf("invalid base url %q: %w", row.BaseUrl, err)
		return out
	}
	host := strings.ToLower(u.Hostname())
	if host == "" {
		out.Status = aibridged.ProviderStatusError
		out.Err = xerrors.Errorf("base url %q has no hostname", row.BaseUrl)
		return out
	}
	if claimedBy, taken := seenHost[host]; taken {
		out.Status = aibridged.ProviderStatusError
		out.Err = xerrors.Errorf("hostname %q already claimed by provider %q", host, claimedBy)
		return out
	}
	seenHost[host] = row.Name
	out.Host = host
	out.Status = aibridged.ProviderStatusEnabled
	return out
}
