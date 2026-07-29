package aibridged

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/yamux"
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/buildinfo"
	aibridgedproto "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/aibridged/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/drpcsdk"
	"github.com/optimus-ide-collab/websocket"
)

// NewWebsocketDialer returns a [Dialer] that connects a standalone AI
// Gateway to optimus-ide-collabd's /api/v2/ai-gateway/serve endpoint over a WebSocket,
// multiplexes it with yamux, and exposes the aibridged DRPC services
// (Recorder, MCPConfigurator, Authorizer, ProviderConfigurator) over it.
// This is the standalone counterpart to API.CreateInMemoryAIBridgeServer,
// which wires the same services over an in-memory pipe for the embedded
// daemon.
//
// The gateway authenticates with an AI Gateway key
// (optimus-ide-collabsdk.AIGatewayKeyHeader), advertises its API version via the
// "version" query parameter, and reports its build version via
// optimus-ide-collabsdk.BuildVersionHeader (used by optimus-ide-collabd for observability only).
// TLS for this connection is governed by the scheme of serverURL and any
// TLS configuration baked into transport.
//
// On a failed upgrade the optimus-ide-collabd HTTP error is returned as a
// *optimus-ide-collabsdk.Error so [Server.connect] can distinguish fatal
// auth/entitlement failures from transient ones.
func readAIGatewayServeError(res *http.Response) error {
	err := optimus-ide-collabsdk.ReadBodyAsError(res)

	var sdkErr *optimus-ide-collabsdk.Error
	if errors.As(err, &sdkErr) && res.StatusCode == http.StatusUnauthorized {
		// /ai-gateway/serve authenticates with an AI Gateway key, not a user
		// session. Generic user-login helpers are misleading here.
		sdkErr.Helper = ""
	}
	return err
}

func NewWebsocketDialer(serverURL *url.URL, transport http.RoundTripper, key string) Dialer {
	return func(ctx context.Context) (DRPCClient, error) {
		serveURL, err := serverURL.Parse("/api/v2/ai-gateway/serve")
		if err != nil {
			return nil, xerrors.Errorf("parse url: %w", err)
		}
		query := serveURL.Query()
		query.Add(aibridgedproto.VersionQueryParam, aibridgedproto.CurrentVersion.String())
		serveURL.RawQuery = query.Encode()

		headers := http.Header{}
		headers.Set(optimus-ide-collabsdk.BuildVersionHeader, buildinfo.Version())
		headers.Set(optimus-ide-collabsdk.AIGatewayKeyHeader, key)

		httpClient := &http.Client{
			Transport: transport,
		}
		// nolint:bodyclose // ReadBodyAsError closes the body; success path hands off to the websocket conn.
		conn, res, err := websocket.Dial(ctx, serveURL.String(), &websocket.DialOptions{
			HTTPClient:      httpClient,
			CompressionMode: websocket.CompressionDisabled,
			HTTPHeader:      headers,
		})
		if err != nil {
			if res == nil {
				return nil, err
			}
			return nil, readAIGatewayServeError(res)
		}
		config := yamux.DefaultConfig()
		config.LogOutput = io.Discard
		// Use a background context because the caller closes the client
		// (and thus the multiplexed session) explicitly.
		_, wsNetConn := optimus-ide-collabsdk.WebsocketNetConn(context.Background(), conn, websocket.MessageBinary)
		conn.SetReadLimit(drpcsdk.YamuxDefaultStreamWindowSize)
		session, err := yamux.Client(wsNetConn, config)
		if err != nil {
			_ = conn.Close(websocket.StatusGoingAway, "")
			_ = wsNetConn.Close()
			return nil, xerrors.Errorf("multiplex client: %w", err)
		}

		dconn := drpcsdk.MultiplexedConn(session)
		return &Client{
			Conn:                           dconn,
			DRPCRecorderClient:             aibridgedproto.NewDRPCRecorderClient(dconn),
			DRPCMCPConfiguratorClient:      aibridgedproto.NewDRPCMCPConfiguratorClient(dconn),
			DRPCAuthorizerClient:           aibridgedproto.NewDRPCAuthorizerClient(dconn),
			DRPCProviderConfiguratorClient: aibridgedproto.NewDRPCProviderConfiguratorClient(dconn),
		}, nil
	}
}
