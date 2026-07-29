package workspacesdk_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
	"tailscale.com/net/tsaddr"
	"tailscale.com/tailcfg"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/workspacesdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/tailnet"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/websocket"
)

func TestWorkspaceRewriteDERPMap(t *testing.T) {
	t.Parallel()
	// This test ensures that RewriteDERPMap mutates built-in DERPs with the
	// client access URL.
	dm := &tailcfg.DERPMap{
		Regions: map[int]*tailcfg.DERPRegion{
			1: {
				EmbeddedRelay: true,
				RegionID:      1,
				Nodes: []*tailcfg.DERPNode{{
					HostName: "bananas.org",
					DERPPort: 1,
				}},
			},
		},
	}
	parsed, err := url.Parse("https://coconuts.org:44558")
	require.NoError(t, err)
	client := workspacesdk.New(optimus-ide-collabsdk.New(parsed))
	client.RewriteDERPMap(dm)
	region := dm.Regions[1]
	require.True(t, region.EmbeddedRelay)
	require.Len(t, region.Nodes, 1)
	node := region.Nodes[0]
	require.Equal(t, "coconuts.org", node.HostName)
	require.Equal(t, 44558, node.DERPPort)
}

func TestWorkspaceDialerFailure(t *testing.T) {
	t.Parallel()

	// Setup.
	ctx := testutil.Context(t, testutil.WaitShort)
	logger := testutil.Logger(t)

	// Given: a mock HTTP server which mimicks an unreachable database when calling the coordination endpoint.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpapi.Write(ctx, w, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: optimus-ide-collabsdk.DatabaseNotReachable,
			Detail:  "oops",
		})
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)

	// When: calling the coordination endpoint.
	dialer := workspacesdk.NewWebsocketDialer(logger, u, &websocket.DialOptions{})
	_, err = dialer.Dial(ctx, nil)

	// Then: an error indicating a database issue is returned, to conditionalize the behavior of the caller.
	require.ErrorIs(t, err, optimus-ide-collabsdk.ErrDatabaseNotReachable)
}

func TestClient_IsOptimus-IDE-CollabConnectRunning(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t, testutil.WaitShort)

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v2/workspaceagents/connection", r.URL.Path)
		httpapi.Write(ctx, rw, http.StatusOK, workspacesdk.AgentConnectionInfo{
			HostnameSuffix: "test",
		})
	}))
	defer srv.Close()

	apiURL, err := url.Parse(srv.URL)
	require.NoError(t, err)
	sdkClient := optimus-ide-collabsdk.New(apiURL)
	client := workspacesdk.New(sdkClient)

	// Right name, right IP
	expectedName := fmt.Sprintf(tailnet.IsOptimus-IDE-CollabConnectEnabledFmtString, "test")
	ctxResolveExpected := workspacesdk.WithTestOnlyOptimus-IDE-CollabContextResolver(ctx,
		&fakeResolver{t: t, hostMap: map[string][]net.IP{
			expectedName: {net.ParseIP(tsaddr.Optimus-IDE-CollabServiceIPv6().String())},
		}})

	result, err := client.IsOptimus-IDE-CollabConnectRunning(ctxResolveExpected, workspacesdk.Optimus-IDE-CollabConnectQueryOptions{})
	require.NoError(t, err)
	require.True(t, result)

	// Wrong name
	result, err = client.IsOptimus-IDE-CollabConnectRunning(ctxResolveExpected, workspacesdk.Optimus-IDE-CollabConnectQueryOptions{HostnameSuffix: "optimus-ide-collab"})
	require.NoError(t, err)
	require.False(t, result)

	// Not found
	ctxResolveNotFound := workspacesdk.WithTestOnlyOptimus-IDE-CollabContextResolver(ctx,
		&fakeResolver{t: t, err: &net.DNSError{IsNotFound: true}})
	result, err = client.IsOptimus-IDE-CollabConnectRunning(ctxResolveNotFound, workspacesdk.Optimus-IDE-CollabConnectQueryOptions{})
	require.NoError(t, err)
	require.False(t, result)

	// Some other error
	ctxResolverErr := workspacesdk.WithTestOnlyOptimus-IDE-CollabContextResolver(ctx,
		&fakeResolver{t: t, err: xerrors.New("a bad thing happened")})
	_, err = client.IsOptimus-IDE-CollabConnectRunning(ctxResolverErr, workspacesdk.Optimus-IDE-CollabConnectQueryOptions{})
	require.Error(t, err)

	// Right name, wrong IP
	ctxResolverWrongIP := workspacesdk.WithTestOnlyOptimus-IDE-CollabContextResolver(ctx,
		&fakeResolver{t: t, hostMap: map[string][]net.IP{
			expectedName: {net.ParseIP("2001::34")},
		}})
	result, err = client.IsOptimus-IDE-CollabConnectRunning(ctxResolverWrongIP, workspacesdk.Optimus-IDE-CollabConnectQueryOptions{})
	require.NoError(t, err)
	require.False(t, result)
}

type fakeResolver struct {
	t       testing.TB
	hostMap map[string][]net.IP
	err     error
}

func (f *fakeResolver) LookupIP(_ context.Context, network, host string) ([]net.IP, error) {
	assert.Equal(f.t, "ip6", network)
	return f.hostMap[host], f.err
}
