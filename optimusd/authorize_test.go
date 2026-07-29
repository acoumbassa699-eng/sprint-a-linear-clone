package optimus-ide-collabd_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestCheckPermissions(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	t.Cleanup(cancel)

	adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		IncludeProvisionerDaemon: true,
	})
	// Create adminClient, member, and org adminClient
	adminUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)
	memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)
	memberUser, err := memberClient.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)
	orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID, rbac.ScopedRoleOrgAdmin(adminUser.OrganizationID))
	orgAdminUser, err := orgAdminClient.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)

	version := optimus-ide-collabdtest.CreateTemplateVersion(t, adminClient, adminUser.OrganizationID, nil)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, adminClient, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, adminClient, adminUser.OrganizationID, version.ID)

	// With admin, member, and org admin
	const (
		readAllUsers           = "read-all-users"
		readOrgWorkspaces      = "read-org-workspaces"
		readMyself             = "read-myself"
		readOwnWorkspaces      = "read-own-workspaces"
		updateSpecificTemplate = "update-specific-template"
	)
	params := map[string]optimus-ide-collabsdk.AuthorizationCheck{
		readAllUsers: {
			Object: optimus-ide-collabsdk.AuthorizationObject{
				ResourceType: optimus-ide-collabsdk.ResourceUser,
			},
			Action: "read",
		},
		readOrgWorkspaces: {
			Object: optimus-ide-collabsdk.AuthorizationObject{
				ResourceType:   optimus-ide-collabsdk.ResourceWorkspace,
				OrganizationID: adminUser.OrganizationID.String(),
			},
			Action: "read",
		},
		readMyself: {
			Object: optimus-ide-collabsdk.AuthorizationObject{
				ResourceType: optimus-ide-collabsdk.ResourceUser,
				OwnerID:      "me",
			},
			Action: "read",
		},
		readOwnWorkspaces: {
			Object: optimus-ide-collabsdk.AuthorizationObject{
				ResourceType:   optimus-ide-collabsdk.ResourceWorkspace,
				OrganizationID: adminUser.OrganizationID.String(),
				OwnerID:        "me",
			},
			Action: "read",
		},
		updateSpecificTemplate: {
			Object: optimus-ide-collabsdk.AuthorizationObject{
				ResourceType: optimus-ide-collabsdk.ResourceTemplate,
				ResourceID:   template.ID.String(),
			},
			Action: "update",
		},
	}

	testCases := []struct {
		Name   string
		Client *optimus-ide-collabsdk.Client
		UserID uuid.UUID
		Check  optimus-ide-collabsdk.AuthorizationResponse
	}{
		{
			Name:   "Admin",
			Client: adminClient,
			UserID: adminUser.UserID,
			Check: map[string]bool{
				readAllUsers:           true,
				readOrgWorkspaces:      true,
				readMyself:             true,
				readOwnWorkspaces:      true,
				updateSpecificTemplate: true,
			},
		},
		{
			Name:   "OrgAdmin",
			Client: orgAdminClient,
			UserID: orgAdminUser.ID,
			Check: map[string]bool{
				readAllUsers:           true,
				readOrgWorkspaces:      true,
				readMyself:             true,
				readOwnWorkspaces:      true,
				updateSpecificTemplate: true,
			},
		},
		{
			Name:   "Member",
			Client: memberClient,
			UserID: memberUser.ID,
			Check: map[string]bool{
				readAllUsers:           false,
				readOrgWorkspaces:      false,
				readMyself:             true,
				readOwnWorkspaces:      true,
				updateSpecificTemplate: false,
			},
		},
	}

	for _, c := range testCases {
		t.Run("CheckAuthorization/"+c.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			t.Cleanup(cancel)

			resp, err := c.Client.AuthCheck(ctx, optimus-ide-collabsdk.AuthorizationRequest{Checks: params})
			require.NoError(t, err, "check perms")
			require.Equal(t, c.Check, resp)
		})
	}
}
