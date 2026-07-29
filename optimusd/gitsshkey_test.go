package optimus-ide-collabd_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/gitsshkey"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/agentsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestGitSSHKey(t *testing.T) {
	t.Parallel()
	t.Run("None", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		res := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		key, err := client.GitSSHKey(ctx, res.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, key.PublicKey)
	})
	t.Run("Ed25519", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			SSHKeygenAlgorithm: gitsshkey.AlgorithmEd25519,
		})
		res := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		key, err := client.GitSSHKey(ctx, res.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, key.PublicKey)
	})
	t.Run("ECDSA", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			SSHKeygenAlgorithm: gitsshkey.AlgorithmECDSA,
		})
		res := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		key, err := client.GitSSHKey(ctx, res.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, key.PublicKey)
	})
	t.Run("RSA4096", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			SSHKeygenAlgorithm: gitsshkey.AlgorithmRSA4096,
		})
		res := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		key, err := client.GitSSHKey(ctx, res.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, key.PublicKey)
	})
	t.Run("Regenerate", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			SSHKeygenAlgorithm: gitsshkey.AlgorithmEd25519,
			Auditor:            auditor,
		})
		res := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		key1, err := client.GitSSHKey(ctx, res.UserID.String())
		require.NoError(t, err)
		require.NotEmpty(t, key1.PublicKey)
		key2, err := client.RegenerateGitSSHKey(ctx, res.UserID.String())
		require.NoError(t, err)
		require.GreaterOrEqual(t, key2.UpdatedAt, key1.UpdatedAt)
		require.NotEmpty(t, key2.PublicKey)

		require.Len(t, auditor.AuditLogs(), 2)
		assert.Equal(t, database.AuditActionWrite, auditor.AuditLogs()[1].Action)
	})
}

func TestAgentGitSSHKey(t *testing.T) {
	t.Parallel()

	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		IncludeProvisionerDaemon: true,
	})
	user := optimus-ide-collabdtest.CreateFirstUser(t, client)
	authToken := uuid.NewString()
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, &echo.Responses{
		Parse:          echo.ParseComplete,
		ProvisionPlan:  echo.PlanComplete,
		ProvisionGraph: echo.ProvisionGraphWithAgent(authToken),
	})
	project := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, project.ID)
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

	agentClient := agentsdk.New(client.URL, agentsdk.WithFixedToken(authToken))

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	agentKey, err := agentClient.GitSSHKey(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, agentKey.PrivateKey)
}

func TestAgentGitSSHKey_APIKeyScopes(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		apiKeyScope string
		expectError bool
	}{
		{apiKeyScope: "all", expectError: false},
		{apiKeyScope: "no_user_data", expectError: true},
	} {
		t.Run(tt.apiKeyScope, func(t *testing.T) {
			t.Parallel()

			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			})
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			authToken := uuid.NewString()
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, &echo.Responses{
				Parse:          echo.ParseComplete,
				ProvisionPlan:  echo.PlanComplete,
				ProvisionGraph: echo.ProvisionGraphWithAgentAndAPIKeyScope(authToken, tt.apiKeyScope),
			})
			project := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
			optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
			workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, project.ID)
			optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

			agentClient := agentsdk.New(client.URL, agentsdk.WithFixedToken(authToken))

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			_, err := agentClient.GitSSHKey(ctx)

			if tt.expectError {
				require.Error(t, err)
				var sdkErr *optimus-ide-collabsdk.Error
				require.ErrorAs(t, err, &sdkErr)
				require.Equal(t, http.StatusForbidden, sdkErr.StatusCode())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
