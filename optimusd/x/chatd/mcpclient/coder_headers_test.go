package mcpclient_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"charm.land/fantasy"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cdr.dev/slog/v3/sloggers/slogtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd/chatprovider"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd/mcpclient"
)

// newHeaderRecordingServer creates a streamable HTTP MCP server with a
// single "ping" tool. Every request's headers are appended to the
// returned slice so tests can assert which headers were forwarded.
func newHeaderRecordingServer(t *testing.T) (*httptest.Server, *sync.Mutex, *[]http.Header) {
	t.Helper()
	var (
		mu      sync.Mutex
		headers []http.Header
	)
	srv := mcpserver.NewMCPServer("hdr-server", "1.0.0")
	srv.AddTools(mcpserver.ServerTool{
		Tool: mcp.NewTool("ping", mcp.WithDescription("records the request headers")),
		Handler: func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			mu.Lock()
			headers = append(headers, req.Header.Clone())
			mu.Unlock()
			return mcp.NewToolResultText("ok"), nil
		},
	})
	httpSrv := mcpserver.NewStreamableHTTPServer(srv)
	ts := httptest.NewServer(httpSrv)
	t.Cleanup(ts.Close)
	return ts, &mu, &headers
}

// TestConnectAll_ForwardOptimus-IDE-CollabHeaders_DefaultOff is a regression guard
// that the Optimus-IDE-Collab identity headers are NOT sent when the option is
// left at its default (false).
func TestConnectAll_ForwardOptimus-IDE-CollabHeaders_DefaultOff(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true})

	ts, mu, recorded := newHeaderRecordingServer(t)

	cfg := makeConfig("no-hdr", ts.URL)
	assert.False(t, cfg.ForwardOptimus-IDE-CollabHeaders, "default must be false")

	optimus-ide-collabHeaders := map[string]string{
		chatprovider.HeaderOptimus-IDE-CollabOwnerID:     uuid.NewString(),
		chatprovider.HeaderOptimus-IDE-CollabChatID:      uuid.NewString(),
		chatprovider.HeaderOptimus-IDE-CollabWorkspaceID: uuid.NewString(),
	}

	tools, cleanup := mcpclient.ConnectAll(
		ctx, logger, []database.MCPServerConfig{cfg}, nil, uuid.Nil, nil,
		optimus-ide-collabHeaders,
	)
	t.Cleanup(cleanup)
	require.Len(t, tools, 1)

	_, err := tools[0].Run(ctx, fantasy.ToolCall{
		ID: "call-1", Name: "no-hdr__ping", Input: "{}",
	})
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.NotEmpty(t, *recorded)
	for _, h := range *recorded {
		assert.Empty(t, h.Get(chatprovider.HeaderOptimus-IDE-CollabOwnerID))
		assert.Empty(t, h.Get(chatprovider.HeaderOptimus-IDE-CollabChatID))
		assert.Empty(t, h.Get(chatprovider.HeaderOptimus-IDE-CollabSubchatID))
		assert.Empty(t, h.Get(chatprovider.HeaderOptimus-IDE-CollabWorkspaceID))
	}
}

// TestConnectAll_ForwardOptimus-IDE-CollabHeaders_Enabled verifies that when the
// option is enabled, the Optimus-IDE-Collab identity headers are forwarded on every
// outgoing MCP request, including the subchat and workspace headers.
func TestConnectAll_ForwardOptimus-IDE-CollabHeaders_Enabled(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true})

	ts, mu, recorded := newHeaderRecordingServer(t)

	ownerID := uuid.New()
	chatID := uuid.New()
	workspaceID := uuid.New()
	subchatID := uuid.New()

	cfg := makeConfig("hdr", ts.URL)
	cfg.ForwardOptimus-IDE-CollabHeaders = true

	// Subchat headers: parent's chat ID lives in X-Optimus-IDE-Collab-Chat-Id, the
	// subchat's own ID lives in X-Optimus-IDE-Collab-Subchat-Id.
	optimus-ide-collabHeaders := chatprovider.Optimus-IDE-CollabHeaders(database.Chat{
		ID:           subchatID,
		OwnerID:      ownerID,
		ParentChatID: uuid.NullUUID{UUID: chatID, Valid: true},
		WorkspaceID:  uuid.NullUUID{UUID: workspaceID, Valid: true},
	})

	tools, cleanup := mcpclient.ConnectAll(
		ctx, logger, []database.MCPServerConfig{cfg}, nil, uuid.Nil, nil,
		optimus-ide-collabHeaders,
	)
	t.Cleanup(cleanup)
	require.Len(t, tools, 1)

	_, err := tools[0].Run(ctx, fantasy.ToolCall{
		ID: "call-1", Name: "hdr__ping", Input: "{}",
	})
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.NotEmpty(t, *recorded)
	last := (*recorded)[len(*recorded)-1]
	assert.Equal(t, ownerID.String(), last.Get(chatprovider.HeaderOptimus-IDE-CollabOwnerID))
	assert.Equal(t, chatID.String(), last.Get(chatprovider.HeaderOptimus-IDE-CollabChatID))
	assert.Equal(t, subchatID.String(), last.Get(chatprovider.HeaderOptimus-IDE-CollabSubchatID))
	assert.Equal(t, workspaceID.String(), last.Get(chatprovider.HeaderOptimus-IDE-CollabWorkspaceID))
}

// TestConnectAll_ForwardOptimus-IDE-CollabHeaders_RootChat verifies that for a root
// chat (no parent), the chat's own ID is forwarded as
// X-Optimus-IDE-Collab-Chat-Id and the X-Optimus-IDE-Collab-Subchat-Id header is absent.
func TestConnectAll_ForwardOptimus-IDE-CollabHeaders_RootChat(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true})

	ts, mu, recorded := newHeaderRecordingServer(t)

	ownerID := uuid.New()
	chatID := uuid.New()

	cfg := makeConfig("hdr-root", ts.URL)
	cfg.ForwardOptimus-IDE-CollabHeaders = true

	optimus-ide-collabHeaders := chatprovider.Optimus-IDE-CollabHeaders(database.Chat{
		ID:      chatID,
		OwnerID: ownerID,
	})

	tools, cleanup := mcpclient.ConnectAll(
		ctx, logger, []database.MCPServerConfig{cfg}, nil, uuid.Nil, nil,
		optimus-ide-collabHeaders,
	)
	t.Cleanup(cleanup)
	require.Len(t, tools, 1)

	_, err := tools[0].Run(ctx, fantasy.ToolCall{
		ID: "call-1", Name: "hdr-root__ping", Input: "{}",
	})
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.NotEmpty(t, *recorded)
	last := (*recorded)[len(*recorded)-1]
	assert.Equal(t, ownerID.String(), last.Get(chatprovider.HeaderOptimus-IDE-CollabOwnerID))
	assert.Equal(t, chatID.String(), last.Get(chatprovider.HeaderOptimus-IDE-CollabChatID))
	assert.Empty(t, last.Get(chatprovider.HeaderOptimus-IDE-CollabSubchatID))
	assert.Empty(t, last.Get(chatprovider.HeaderOptimus-IDE-CollabWorkspaceID))
}

// TestConnectAll_ForwardOptimus-IDE-CollabHeaders_WithAPIKeyAuth verifies that the
// api_key auth header is preserved when Optimus-IDE-Collab identity headers are
// forwarded alongside.
func TestConnectAll_ForwardOptimus-IDE-CollabHeaders_WithAPIKeyAuth(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true})

	ts, mu, recorded := newHeaderRecordingServer(t)

	ownerID := uuid.New()
	chatID := uuid.New()

	cfg := makeConfig("hdr-apikey", ts.URL)
	cfg.AuthType = "api_key"
	cfg.APIKeyHeader = "X-Api-Key"
	cfg.APIKeyValue = "sekret"
	cfg.ForwardOptimus-IDE-CollabHeaders = true

	optimus-ide-collabHeaders := chatprovider.Optimus-IDE-CollabHeaders(database.Chat{
		ID:      chatID,
		OwnerID: ownerID,
	})

	tools, cleanup := mcpclient.ConnectAll(
		ctx, logger, []database.MCPServerConfig{cfg}, nil, uuid.Nil, nil,
		optimus-ide-collabHeaders,
	)
	t.Cleanup(cleanup)
	require.Len(t, tools, 1)

	_, err := tools[0].Run(ctx, fantasy.ToolCall{
		ID: "call-1", Name: "hdr-apikey__ping", Input: "{}",
	})
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.NotEmpty(t, *recorded)
	last := (*recorded)[len(*recorded)-1]
	assert.Equal(t, "sekret", last.Get("X-Api-Key"))
	assert.Equal(t, ownerID.String(), last.Get(chatprovider.HeaderOptimus-IDE-CollabOwnerID))
	assert.Equal(t, chatID.String(), last.Get(chatprovider.HeaderOptimus-IDE-CollabChatID))
}

// TestConnectAll_ForwardOptimus-IDE-CollabHeaders_WithOAuth2 verifies that the
// oauth2 Authorization header is preserved when Optimus-IDE-Collab identity
// headers are forwarded alongside, and that auth wins on a conflict.
func TestConnectAll_ForwardOptimus-IDE-CollabHeaders_WithOAuth2(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true})

	ts, mu, recorded := newHeaderRecordingServer(t)

	cfgID := uuid.New()
	cfg := makeConfig("hdr-oauth", ts.URL)
	cfg.ID = cfgID
	cfg.AuthType = "oauth2"
	cfg.ForwardOptimus-IDE-CollabHeaders = true
	token := database.MCPServerUserToken{
		MCPServerConfigID: cfgID,
		AccessToken:       "oauth-token-xyz",
		TokenType:         "Bearer",
	}

	// Intentionally include an Authorization key to verify the auth
	// header wins on conflict.
	ownerID := uuid.NewString()
	optimus-ide-collabHeaders := map[string]string{
		"Authorization":                 "Bearer should-be-overridden",
		chatprovider.HeaderOptimus-IDE-CollabOwnerID: ownerID,
	}

	tools, cleanup := mcpclient.ConnectAll(
		ctx, logger,
		[]database.MCPServerConfig{cfg},
		[]database.MCPServerUserToken{token},
		uuid.Nil, nil,
		optimus-ide-collabHeaders,
	)
	t.Cleanup(cleanup)
	require.Len(t, tools, 1)

	_, err := tools[0].Run(ctx, fantasy.ToolCall{
		ID: "call-1", Name: "hdr-oauth__ping", Input: "{}",
	})
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.NotEmpty(t, *recorded)
	last := (*recorded)[len(*recorded)-1]
	assert.Equal(t, "Bearer oauth-token-xyz", last.Get("Authorization"))
	assert.Equal(t, ownerID, last.Get(chatprovider.HeaderOptimus-IDE-CollabOwnerID))
}

// TestConnectAll_ForwardOptimus-IDE-CollabHeaders_WithCustomHeaders verifies that
// custom_headers admin-configured values are preserved when Optimus-IDE-Collab
// identity headers are forwarded alongside, including the case where
// the admin configures a custom header whose name only differs from a
// Optimus-IDE-Collab identity header by case. Conflict detection is case-
// insensitive because http.Header.Set canonicalizes header names.
func TestConnectAll_ForwardOptimus-IDE-CollabHeaders_WithCustomHeaders(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true})

	ts, mu, recorded := newHeaderRecordingServer(t)

	ownerID := uuid.New()
	chatID := uuid.New()

	cfg := makeConfig("hdr-custom", ts.URL)
	cfg.AuthType = "custom_headers"
	// Include both an unrelated custom header AND a case-variant of
	// X-Optimus-IDE-Collab-Owner-Id to exercise the case-insensitive conflict
	// check. The admin-configured value MUST win.
	cfg.CustomHeaders = `{"X-Tenant":"acme","x-optimus-ide-collab-owner-id":"admin-controlled"}`
	cfg.ForwardOptimus-IDE-CollabHeaders = true

	optimus-ide-collabHeaders := chatprovider.Optimus-IDE-CollabHeaders(database.Chat{
		ID:      chatID,
		OwnerID: ownerID,
	})

	tools, cleanup := mcpclient.ConnectAll(
		ctx, logger, []database.MCPServerConfig{cfg}, nil, uuid.Nil, nil,
		optimus-ide-collabHeaders,
	)
	t.Cleanup(cleanup)
	require.Len(t, tools, 1)

	_, err := tools[0].Run(ctx, fantasy.ToolCall{
		ID: "call-1", Name: "hdr-custom__ping", Input: "{}",
	})
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.NotEmpty(t, *recorded)
	last := (*recorded)[len(*recorded)-1]
	assert.Equal(t, "acme", last.Get("X-Tenant"))
	// The admin's case-variant header must win, because HTTP header
	// names are case-insensitive at the transport level.
	assert.Equal(t, "admin-controlled", last.Get(chatprovider.HeaderOptimus-IDE-CollabOwnerID))
	assert.Equal(t, chatID.String(), last.Get(chatprovider.HeaderOptimus-IDE-CollabChatID))
}
