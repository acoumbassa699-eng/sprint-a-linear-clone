package optimus-ide-collabd_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestCheckACLPermissions(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	t.Cleanup(cancel)

	adminClient, adminUser := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		Options: &optimus-ide-collabdtest.Options{
			IncludeProvisionerDaemon: true,
		},
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
			},
		},
	})
	// Create member and org adminClient
	memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)
	memberUser, err := memberClient.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)
	orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID, rbac.ScopedRoleOrgAdmin(adminUser.OrganizationID))
	orgAdminUser, err := orgAdminClient.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)

	version := optimus-ide-collabdtest.CreateTemplateVersion(t, adminClient, adminUser.OrganizationID, nil)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, adminClient, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, adminClient, adminUser.OrganizationID, version.ID)

	err = adminClient.UpdateTemplateACL(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateACL{
		UserPerms: map[string]optimus-ide-collabsdk.TemplateRole{
			memberUser.ID.String(): optimus-ide-collabsdk.TemplateRoleAdmin,
		},
	})
	require.NoError(t, err)

	const (
		updateSpecificTemplate = "read-specific-template"
	)
	params := map[string]optimus-ide-collabsdk.AuthorizationCheck{
		updateSpecificTemplate: {
			Object: optimus-ide-collabsdk.AuthorizationObject{
				ResourceType: optimus-ide-collabsdk.ResourceTemplate,
				ResourceID:   template.ID.String(),
			},
			Action: optimus-ide-collabsdk.ActionUpdate,
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
				updateSpecificTemplate: true,
			},
		},
		{
			Name:   "OrgAdmin",
			Client: orgAdminClient,
			UserID: orgAdminUser.ID,
			Check: map[string]bool{
				updateSpecificTemplate: true,
			},
		},
		{
			Name:   "Member",
			Client: memberClient,
			UserID: memberUser.ID,
			Check: map[string]bool{
				updateSpecificTemplate: true,
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
