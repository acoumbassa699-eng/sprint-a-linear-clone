package optimus-ide-collabd_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/agentsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// TestAgentGitSSHKeyCustomRoles tests that the agent can fetch its git ssh key when
// the user has a custom role in a second workspace.
func TestAgentGitSSHKeyCustomRoles(t *testing.T) {
	t.Parallel()

	owner, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		Options: &optimus-ide-collabdtest.Options{
			IncludeProvisionerDaemon: true,
		},
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureCustomRoles:                1,
				optimus-ide-collabsdk.FeatureMultipleOrganizations:      1,
				optimus-ide-collabsdk.FeatureExternalProvisionerDaemons: 1,
			},
		},
	})

	// When custom roles exist in a second organization
	org := optimus-ide-collabdenttest.CreateOrganization(t, owner, optimus-ide-collabdenttest.CreateOrganizationOptions{
		IncludeProvisionerDaemon: true,
	})

	ctx := testutil.Context(t, testutil.WaitShort)
	//nolint:gocritic // required to make orgs
	newRole, err := owner.CreateOrganizationRole(ctx, optimus-ide-collabsdk.Role{
		Name:            "custom",
		OrganizationID:  org.ID.String(),
		DisplayName:     "",
		SitePermissions: nil,
		OrganizationPermissions: optimus-ide-collabsdk.CreatePermissions(map[optimus-ide-collabsdk.RBACResource][]optimus-ide-collabsdk.RBACAction{
			optimus-ide-collabsdk.ResourceTemplate: {optimus-ide-collabsdk.ActionRead, optimus-ide-collabsdk.ActionCreate, optimus-ide-collabsdk.ActionUpdate},
		}),
		UserPermissions: nil,
	})
	require.NoError(t, err)

	// Create the new user
	client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, owner, org.ID, rbac.RoleIdentifier{Name: newRole.Name, OrganizationID: org.ID})

	// Create the workspace + agent
	authToken := uuid.NewString()
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, org.ID, &echo.Responses{
		Parse:          echo.ParseComplete,
		ProvisionPlan:  echo.PlanComplete,
		ProvisionGraph: echo.ProvisionGraphWithAgent(authToken),
	})
	project := optimus-ide-collabdtest.CreateTemplate(t, client, org.ID, version.ID)
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
