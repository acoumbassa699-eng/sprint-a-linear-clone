package optimus-ide-collabd_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// TestCompositeWorkspaceScopes verifies that the composite
// optimus-ide-collab:workspaces.* scopes grant the permissions needed for
// workspace lifecycle operations when used on scoped API tokens.
func TestCompositeWorkspaceScopes(t *testing.T) {
	t.Parallel()

	// setupWorkspace creates a server with a provisioner daemon, an
	// admin user, a template, and a workspace. It returns the admin
	// client and the workspace so sub-tests can create scoped tokens
	// and act on them.
	type setupResult struct {
		adminClient *optimus-ide-collabsdk.Client
		workspace   optimus-ide-collabsdk.Workspace
	}
	setup := func(t *testing.T) setupResult {
		t.Helper()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			IncludeProvisionerDaemon: true,
		})
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, firstUser.OrganizationID, &echo.Responses{
			Parse:          echo.ParseComplete,
			ProvisionPlan:  echo.PlanComplete,
			ProvisionApply: echo.ApplyComplete,
			ProvisionGraph: echo.GraphComplete,
		})
		template := optimus-ide-collabdtest.CreateTemplate(t, client, firstUser.OrganizationID, version.ID)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

		return setupResult{
			adminClient: client,
			workspace:   workspace,
		}
	}

	// scopedClient creates an API token restricted to the given scopes
	// and returns a new client authenticated with that token.
	scopedClient := func(t *testing.T, adminClient *optimus-ide-collabsdk.Client, scopes []optimus-ide-collabsdk.APIKeyScope) *optimus-ide-collabsdk.Client {
		t.Helper()
		ctx, cancel := context.WithTimeout(t.Context(), testutil.WaitShort)
		defer cancel()

		resp, err := adminClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
			Scopes: scopes,
		})
		require.NoError(t, err, "creating scoped token")

		scoped := optimus-ide-collabsdk.New(
			adminClient.URL,
			optimus-ide-collabsdk.WithSessionToken(resp.Key),
			optimus-ide-collabsdk.WithHTTPClient(optimus-ide-collabdtest.NewIsolatedHTTPClient(adminClient.URL)),
		)
		t.Cleanup(func() { scoped.HTTPClient.CloseIdleConnections() })
		return scoped
	}

	// optimus-ide-collab:workspaces.create — token should be able to create a
	// workspace via POST /users/{user}/workspaces.
	t.Run("WorkspacesCreate", func(t *testing.T) {
		t.Parallel()
		s := setup(t)

		scoped := scopedClient(t, s.adminClient, []optimus-ide-collabsdk.APIKeyScope{
			optimus-ide-collabsdk.APIKeyScopeOptimus-IDE-CollabWorkspacesCreate,
		})

		ctx, cancel := context.WithTimeout(t.Context(), testutil.WaitLong)
		defer cancel()

		// List workspaces (requires workspace:read, included in the
		// composite scope).
		workspaces, err := scoped.Workspaces(ctx, optimus-ide-collabsdk.WorkspaceFilter{})
		require.NoError(t, err, "listing workspaces with optimus-ide-collab:workspaces.create scope")
		require.NotEmpty(t, workspaces.Workspaces, "should see at least the existing workspace")

		_, err = scoped.CreateUserWorkspace(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateWorkspaceRequest{
			TemplateID: s.workspace.TemplateID,
			Name:       optimus-ide-collabdtest.RandomUsername(t),
		})
		require.NoError(t, err, "creating workspace with optimus-ide-collab:workspaces.create scope")
	})

	// optimus-ide-collab:workspaces.operate — token should be able to read and
	// update workspace metadata.
	t.Run("WorkspacesOperate", func(t *testing.T) {
		t.Parallel()
		s := setup(t)

		scoped := scopedClient(t, s.adminClient, []optimus-ide-collabsdk.APIKeyScope{
			optimus-ide-collabsdk.APIKeyScopeOptimus-IDE-CollabWorkspacesOperate,
		})

		ctx, cancel := context.WithTimeout(t.Context(), testutil.WaitLong)
		defer cancel()

		// Read the workspace by ID (requires workspace:read).
		ws, err := scoped.Workspace(ctx, s.workspace.ID)
		require.NoError(t, err, "reading workspace with optimus-ide-collab:workspaces.operate scope")
		require.Equal(t, s.workspace.ID, ws.ID)

		// Update the workspace metadata (requires workspace:update). This goes
		// through the PATCH /workspaces/{workspace} endpoint.
		err = scoped.UpdateWorkspaceTTL(ctx, s.workspace.ID, optimus-ide-collabsdk.UpdateWorkspaceTTLRequest{
			TTLMillis: ptr.Ref[int64]((time.Hour).Milliseconds()),
		})
		require.NoError(t, err, "updating workspace with optimus-ide-collab:workspaces.operate scope")

		// Trigger a start build (requires workspace:update). This goes
		// through POST /workspaces/{workspace}/builds.
		started, err := scoped.CreateWorkspaceBuild(ctx, s.workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
			TemplateVersionID: ws.LatestBuild.TemplateVersionID,
			Transition:        optimus-ide-collabsdk.WorkspaceTransitionStart,
		})
		require.NoError(t, err, "starting workspace with optimus-ide-collab:workspaces.operate scope")
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, scoped, started.ID)

		_, err = scoped.CreateWorkspaceBuild(ctx, s.workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
			TemplateVersionID: ws.LatestBuild.TemplateVersionID,
			Transition:        optimus-ide-collabsdk.WorkspaceTransitionStop,
		})
		require.NoError(t, err, "starting workspace with optimus-ide-collab:workspaces.operate scope")

		// Verify we cannot create a new workspace — the operate scope
		// should not include workspace:create or template:read/use.
		_, err = scoped.CreateUserWorkspace(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateWorkspaceRequest{
			TemplateID: s.workspace.TemplateID,
			Name:       optimus-ide-collabdtest.RandomUsername(t),
		})
		require.Error(t, err, "creating workspace should fail with optimus-ide-collab:workspaces.operate scope")
	})

	// optimus-ide-collab:workspaces.delete — token should be able to read
	// workspaces and trigger a delete build.
	t.Run("WorkspacesDelete", func(t *testing.T) {
		t.Parallel()
		s := setup(t)

		scoped := scopedClient(t, s.adminClient, []optimus-ide-collabsdk.APIKeyScope{
			optimus-ide-collabsdk.APIKeyScopeOptimus-IDE-CollabWorkspacesDelete,
		})

		ctx, cancel := context.WithTimeout(t.Context(), testutil.WaitLong)
		defer cancel()

		// Read the workspace by ID (requires workspace:read).
		ws, err := scoped.Workspace(ctx, s.workspace.ID)
		require.NoError(t, err, "reading workspace with optimus-ide-collab:workspaces.delete scope")
		require.Equal(t, s.workspace.ID, ws.ID)

		// Delete the workspace via a delete transition build.
		_, err = scoped.CreateWorkspaceBuild(ctx, s.workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
			TemplateVersionID: ws.LatestBuild.TemplateVersionID,
			Transition:        optimus-ide-collabsdk.WorkspaceTransitionDelete,
		})
		require.NoError(t, err, "deleting workspace with optimus-ide-collab:workspaces.delete scope")
	})
}
