package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestRemoveOrganizationMembers(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		ownerClient, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})

		secondOrganization := optimus-ide-collabdenttest.CreateOrganization(t, ownerClient, optimus-ide-collabdenttest.CreateOrganizationOptions{})
		orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, secondOrganization.ID, rbac.ScopedRoleOrgAdmin(secondOrganization.ID))
		_, user := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, secondOrganization.ID)

		ctx := testutil.Context(t, testutil.WaitMedium)

		inv, root := clitest.New(t, "organization", "members", "remove", "-O", secondOrganization.ID.String(), user.Username)
		clitest.SetupConfig(t, orgAdminClient, root)

		buf := new(bytes.Buffer)
		inv.Stdout = buf
		err := inv.WithContext(ctx).Run()
		require.NoError(t, err)

		members, err := orgAdminClient.OrganizationMembers(ctx, secondOrganization.ID)
		require.NoError(t, err)

		require.Len(t, members, 2)
	})

	t.Run("UserNotExists", func(t *testing.T) {
		t.Parallel()

		ownerClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)
		orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.ScopedRoleOrgAdmin(owner.OrganizationID))

		ctx := testutil.Context(t, testutil.WaitMedium)

		inv, root := clitest.New(t, "organization", "members", "remove", "-O", owner.OrganizationID.String(), "random_name")
		clitest.SetupConfig(t, orgAdminClient, root)

		buf := new(bytes.Buffer)
		inv.Stdout = buf
		err := inv.WithContext(ctx).Run()
		require.ErrorContains(t, err, "Resource not found or you do not have access to this resource")
	})
}

func TestEnterpriseListOrganizationMembers(t *testing.T) {
	t.Parallel()

	t.Run("CustomRole", func(t *testing.T) {
		t.Parallel()

		ownerClient, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureCustomRoles: 1,
				},
			},
		})

		ctx := testutil.Context(t, testutil.WaitMedium)
		//nolint:gocritic // only owners can patch roles
		customRole, err := ownerClient.CreateOrganizationRole(ctx, optimus-ide-collabsdk.Role{
			Name:            "custom",
			OrganizationID:  owner.OrganizationID.String(),
			DisplayName:     "Custom Role",
			SitePermissions: nil,
			OrganizationPermissions: optimus-ide-collabsdk.CreatePermissions(map[optimus-ide-collabsdk.RBACResource][]optimus-ide-collabsdk.RBACAction{
				optimus-ide-collabsdk.ResourceWorkspace: {optimus-ide-collabsdk.ActionRead},
			}),
			UserPermissions: nil,
		})
		require.NoError(t, err)

		client, user := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.RoleUserAdmin(), rbac.RoleIdentifier{
			Name:           customRole.Name,
			OrganizationID: owner.OrganizationID,
		}, rbac.ScopedRoleOrgAdmin(owner.OrganizationID))

		inv, root := clitest.New(t, "organization", "members", "list", "-c", "user id,username,organization roles")
		clitest.SetupConfig(t, client, root)

		buf := new(bytes.Buffer)
		inv.Stdout = buf
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)
		require.Contains(t, buf.String(), user.Username)
		require.Contains(t, buf.String(), owner.UserID.String())
		// Check the display name is the value in the cli list
		require.Contains(t, buf.String(), customRole.DisplayName)
	})
}

func TestAssignOrganizationMemberRole(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		ownerClient, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureCustomRoles: 1,
				},
			},
		})
		_, user := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.RoleUserAdmin())

		ctx := testutil.Context(t, testutil.WaitMedium)
		// nolint:gocritic // requires owner role to create
		customRole, err := ownerClient.CreateOrganizationRole(ctx, optimus-ide-collabsdk.Role{
			Name:            "custom-role",
			OrganizationID:  owner.OrganizationID.String(),
			DisplayName:     "Custom Role",
			SitePermissions: nil,
			OrganizationPermissions: optimus-ide-collabsdk.CreatePermissions(map[optimus-ide-collabsdk.RBACResource][]optimus-ide-collabsdk.RBACAction{
				optimus-ide-collabsdk.ResourceWorkspace: {optimus-ide-collabsdk.ActionRead},
			}),
			UserPermissions: nil,
		})
		require.NoError(t, err)

		inv, root := clitest.New(t, "organization", "members", "edit-roles", user.Username, optimus-ide-collabsdk.RoleOrganizationAdmin, customRole.Name)
		// nolint:gocritic // you cannot change your own roles
		clitest.SetupConfig(t, ownerClient, root)

		buf := new(bytes.Buffer)
		inv.Stdout = buf
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)
		require.Contains(t, buf.String(), must(rbac.RoleByName(rbac.ScopedRoleOrgAdmin(owner.OrganizationID))).DisplayName)
		require.Contains(t, buf.String(), customRole.DisplayName)
	})
}

func must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}
	return v
}
