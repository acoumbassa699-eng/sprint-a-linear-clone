package enterprise_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestEnterpriseMembers(t *testing.T) {
	t.Parallel()

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		owner, first := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
					optimus-ide-collabsdk.FeatureTemplateRBAC:          1,
				},
			},
		})

		secondOrg := optimus-ide-collabdenttest.CreateOrganization(t, owner, optimus-ide-collabdenttest.CreateOrganizationOptions{})

		orgAdminClient, orgAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, owner, secondOrg.ID, rbac.ScopedRoleOrgAdmin(secondOrg.ID))
		_, user := optimus-ide-collabdtest.CreateAnotherUser(t, owner, secondOrg.ID)

		ctx := testutil.Context(t, testutil.WaitMedium)

		// Groups exist to ensure a user removed from the org loses their
		// group access.
		g1, err := orgAdminClient.CreateGroup(ctx, secondOrg.ID, optimus-ide-collabsdk.CreateGroupRequest{
			Name:        "foo",
			DisplayName: "Foo",
		})
		require.NoError(t, err)

		g2, err := orgAdminClient.CreateGroup(ctx, secondOrg.ID, optimus-ide-collabsdk.CreateGroupRequest{
			Name:        "bar",
			DisplayName: "Bar",
		})
		require.NoError(t, err)

		// Verify the org of 3 members
		members, err := orgAdminClient.OrganizationMembers(ctx, secondOrg.ID)
		require.NoError(t, err)
		require.Len(t, members, 3)
		require.ElementsMatch(t,
			[]uuid.UUID{first.UserID, user.ID, orgAdmin.ID},
			slice.List(members, onlyIDs))

		// Add the member to some groups
		_, err = orgAdminClient.PatchGroup(ctx, g1.ID, optimus-ide-collabsdk.PatchGroupRequest{
			AddUsers: []string{user.ID.String()},
		})
		require.NoError(t, err)

		_, err = orgAdminClient.PatchGroup(ctx, g2.ID, optimus-ide-collabsdk.PatchGroupRequest{
			AddUsers: []string{user.ID.String()},
		})
		require.NoError(t, err)

		// Verify group membership
		userGroups, err := orgAdminClient.Groups(ctx, optimus-ide-collabsdk.GroupArguments{
			HasMember: user.ID.String(),
		})
		require.NoError(t, err)
		// Everyone group + 2 groups
		require.Len(t, userGroups, 3)

		// Delete a member
		err = orgAdminClient.DeleteOrganizationMember(ctx, secondOrg.ID, user.Username)
		require.NoError(t, err)

		members, err = orgAdminClient.OrganizationMembers(ctx, secondOrg.ID)
		require.NoError(t, err)
		require.Len(t, members, 2)
		require.ElementsMatch(t,
			[]uuid.UUID{first.UserID, orgAdmin.ID},
			slice.List(members, onlyIDs))

		// User should now belong to 0 groups
		userGroups, err = orgAdminClient.Groups(ctx, optimus-ide-collabsdk.GroupArguments{
			HasMember: user.ID.String(),
		})
		require.NoError(t, err)
		require.Len(t, userGroups, 0)
	})

	t.Run("PostUser", func(t *testing.T) {
		t.Parallel()

		owner, first := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})

		ctx := testutil.Context(t, testutil.WaitMedium)
		org := optimus-ide-collabdenttest.CreateOrganization(t, owner, optimus-ide-collabdenttest.CreateOrganizationOptions{})

		// Make a user not in the second organization
		_, user := optimus-ide-collabdtest.CreateAnotherUser(t, owner, first.OrganizationID)

		// Use scoped user admin in org to add the user
		client, userAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, owner, org.ID, rbac.ScopedRoleOrgUserAdmin(org.ID))

		members, err := client.OrganizationMembers(ctx, org.ID)
		require.NoError(t, err)
		require.Len(t, members, 2) // Verify the 2 members at the start

		// Add user to org
		_, err = client.PostOrganizationMember(ctx, org.ID, user.Username)
		require.NoError(t, err)

		members, err = client.OrganizationMembers(ctx, org.ID)
		require.NoError(t, err)
		// Owner + user admin + new member
		require.Len(t, members, 3)
		require.ElementsMatch(t,
			[]uuid.UUID{first.UserID, user.ID, userAdmin.ID},
			slice.List(members, onlyIDs))
	})

	t.Run("PostUserNotExists", func(t *testing.T) {
		t.Parallel()
		owner, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})

		org := optimus-ide-collabdenttest.CreateOrganization(t, owner, optimus-ide-collabdenttest.CreateOrganizationOptions{})

		ctx := testutil.Context(t, testutil.WaitMedium)
		// Add user to org
		//nolint:gocritic // Using owner to ensure it's not a 404 error
		_, err := owner.PostOrganizationMember(ctx, org.ID, uuid.NewString())
		require.Error(t, err)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Contains(t, apiErr.Message, "Resource not found or you do not have access to this resource")
	})

	// Calling it from a user without the org access.
	t.Run("ListNotInOrg", func(t *testing.T) {
		t.Parallel()

		owner, first := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, owner, first.OrganizationID, rbac.ScopedRoleOrgAdmin(first.OrganizationID))
		org := optimus-ide-collabdenttest.CreateOrganization(t, owner, optimus-ide-collabdenttest.CreateOrganizationOptions{})

		ctx := testutil.Context(t, testutil.WaitShort)

		// 404 error is expected instead of a 403/401 to not leak existence of
		// an organization.
		_, err := client.OrganizationMembers(ctx, org.ID)
		require.ErrorContains(t, err, "404")
	})
}

func onlyIDs(u optimus-ide-collabsdk.OrganizationMemberWithUserData) uuid.UUID {
	return u.UserID
}
