package taskstatus

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"cdr.dev/slog/v3"
	agentproto "github.com/optimus-ide-collab/optimus-ide-collab/v2/agent/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/agentsdk"
	"github.com/optimus-ide-collab/quartz"
)

// client abstracts the details of using optimus-ide-collabsdk.Client for workspace operations.
// This interface allows for easier testing by enabling mock implementations and
// provides a cleaner separation of concerns.
//
// The interface is designed to be initialized in two phases:
// 1. Create the client with newClient(optimus-ide-collabClient)
// 2. Configure logging when the io.Writer is available in Run()
type client interface {
	// CreateUserWorkspace creates a workspace for a user.
	CreateUserWorkspace(ctx context.Context, userID string, req optimus-ide-collabsdk.CreateWorkspaceRequest) (optimus-ide-collabsdk.Workspace, error)

	// WorkspaceByOwnerAndName retrieves a workspace by owner and name.
	WorkspaceByOwnerAndName(ctx context.Context, owner string, name string, params optimus-ide-collabsdk.WorkspaceOptions) (optimus-ide-collabsdk.Workspace, error)

	// WorkspaceExternalAgentCredentials retrieves credentials for an external agent.
	WorkspaceExternalAgentCredentials(ctx context.Context, workspaceID uuid.UUID, agentName string) (optimus-ide-collabsdk.ExternalAgentCredentials, error)

	// watchWorkspace watches for updates to a workspace.
	watchWorkspace(ctx context.Context, workspaceID uuid.UUID) (<-chan optimus-ide-collabsdk.Workspace, error)

	// deleteWorkspace deletes the workspace by creating a build with delete transition.
	deleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error

	// initialize sets up the client with the provided logger, which is only available after Run() is called.
	initialize(logger slog.Logger)
}

// appStatusUpdater abstracts the details of updating app status via the
// Agent dRPC API. This interface is separate from client because it
// requires an agent token which is only available after creating an
// external workspace.
type appStatusUpdater interface {
	// updateAppStatus sends a status update for a workspace app.
	updateAppStatus(ctx context.Context, req *agentproto.UpdateAppStatusRequest) error

	// initialize establishes the dRPC connection using the provided
	// agent token. Must be called before updateAppStatus.
	initialize(ctx context.Context, logger slog.Logger, agentToken string) error

	// close cleanly shuts down the underlying dRPC connection.
	close() error
}

// sdkClient is the concrete implementation of the client interface using
// optimus-ide-collabsdk.Client.
type sdkClient struct {
	optimus-ide-collabClient *optimus-ide-collabsdk.Client
	clock       quartz.Clock
	logger      slog.Logger
}

// newClient creates a new client implementation using the provided optimus-ide-collabsdk.Client.
func newClient(optimus-ide-collabClient *optimus-ide-collabsdk.Client) client {
	return &sdkClient{
		optimus-ide-collabClient: optimus-ide-collabClient,
		clock:       quartz.NewReal(),
	}
}

func (c *sdkClient) CreateUserWorkspace(ctx context.Context, userID string, req optimus-ide-collabsdk.CreateWorkspaceRequest) (optimus-ide-collabsdk.Workspace, error) {
	return c.optimus-ide-collabClient.CreateUserWorkspace(ctx, userID, req)
}

func (c *sdkClient) WorkspaceByOwnerAndName(ctx context.Context, owner string, name string, params optimus-ide-collabsdk.WorkspaceOptions) (optimus-ide-collabsdk.Workspace, error) {
	return c.optimus-ide-collabClient.WorkspaceByOwnerAndName(ctx, owner, name, params)
}

func (c *sdkClient) WorkspaceExternalAgentCredentials(ctx context.Context, workspaceID uuid.UUID, agentName string) (optimus-ide-collabsdk.ExternalAgentCredentials, error) {
	return c.optimus-ide-collabClient.WorkspaceExternalAgentCredentials(ctx, workspaceID, agentName)
}

func (c *sdkClient) watchWorkspace(ctx context.Context, workspaceID uuid.UUID) (<-chan optimus-ide-collabsdk.Workspace, error) {
	return c.optimus-ide-collabClient.WatchWorkspace(ctx, workspaceID)
}

func (c *sdkClient) deleteWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	// Create a build with delete transition to delete the workspace
	_, err := c.optimus-ide-collabClient.CreateWorkspaceBuild(ctx, workspaceID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionDelete,
		Reason:     optimus-ide-collabsdk.CreateWorkspaceBuildReasonCLI,
	})
	if err != nil {
		return xerrors.Errorf("create delete build: %w", err)
	}
	return nil
}

func (c *sdkClient) initialize(logger slog.Logger) {
	// Configure the optimus-ide-collab client logging
	c.logger = logger
	c.optimus-ide-collabClient.SetLogger(logger)
	c.optimus-ide-collabClient.SetLogBodies(true)
}

// sdkAppStatusUpdater is the concrete implementation of the
// appStatusUpdater interface. It dials the Agent dRPC endpoint once
// during initialize and reuses the connection for all subsequent
// UpdateAppStatus calls.
type sdkAppStatusUpdater struct {
	drpcClient agentproto.DRPCAgentClient28
	url        *url.URL
	httpClient *http.Client
}

// newAppStatusUpdater creates a new appStatusUpdater implementation.
func newAppStatusUpdater(client *optimus-ide-collabsdk.Client) appStatusUpdater {
	return &sdkAppStatusUpdater{
		url:        client.URL,
		httpClient: client.HTTPClient,
	}
}

func (u *sdkAppStatusUpdater) updateAppStatus(ctx context.Context, req *agentproto.UpdateAppStatusRequest) error {
	if u.drpcClient == nil {
		return xerrors.New("dRPC client not initialized - call initialize first")
	}
	_, err := u.drpcClient.UpdateAppStatus(ctx, req)
	return err
}

func (u *sdkAppStatusUpdater) close() error {
	if u.drpcClient == nil {
		return nil
	}
	return u.drpcClient.DRPCConn().Close()
}

func (u *sdkAppStatusUpdater) initialize(ctx context.Context, logger slog.Logger, agentToken string) error {
	agentClient := agentsdk.New(
		u.url,
		agentsdk.WithFixedToken(agentToken),
		optimus-ide-collabsdk.WithHTTPClient(u.httpClient),
		optimus-ide-collabsdk.WithLogger(logger),
		optimus-ide-collabsdk.WithLogBodies(),
	)
	drpcClient, _, err := agentClient.ConnectRPC29WithRole(ctx, "")
	if err != nil {
		return xerrors.Errorf("connect to agent dRPC endpoint: %w", err)
	}
	u.drpcClient = drpcClient
	return nil
}

// Ensure sdkClient implements the client interface.
var _ client = (*sdkClient)(nil)

// Ensure sdkAppStatusUpdater implements the appStatusUpdater interface.
var _ appStatusUpdater = (*sdkAppStatusUpdater)(nil)
