package cli_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cryptorand"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
)

func TestUserDelete(t *testing.T) {
	t.Parallel()
	t.Run("Username", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())

		pw, err := cryptorand.String(16)
		require.NoError(t, err)

		_, err = client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           "colin5@optimus-ide-collab.com",
			Username:        "coolin",
			Password:        pw,
			UserLoginType:   optimus-ide-collabsdk.LoginTypePassword,
			OrganizationIDs: []uuid.UUID{owner.OrganizationID},
		})
		require.NoError(t, err)

		inv, root := clitest.New(t, "users", "delete", "coolin")
		clitest.SetupConfig(t, userAdmin, root)
		stdout := expecter.NewAttachedToInvocation(t, inv)
		errC := make(chan error)
		go func() {
			errC <- inv.Run()
		}()
		require.NoError(t, <-errC)
		stdout.ExpectMatch(ctx, "coolin")
	})

	t.Run("UserID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())

		pw, err := cryptorand.String(16)
		require.NoError(t, err)

		user, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           "colin5@optimus-ide-collab.com",
			Username:        "coolin",
			Password:        pw,
			UserLoginType:   optimus-ide-collabsdk.LoginTypePassword,
			OrganizationIDs: []uuid.UUID{owner.OrganizationID},
		})
		require.NoError(t, err)

		inv, root := clitest.New(t, "users", "delete", user.ID.String())
		clitest.SetupConfig(t, userAdmin, root)
		stdout := expecter.NewAttachedToInvocation(t, inv)
		errC := make(chan error)
		go func() {
			errC <- inv.Run()
		}()
		require.NoError(t, <-errC)
		stdout.ExpectMatch(ctx, "coolin")
	})

	t.Run("UserID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())

		pw, err := cryptorand.String(16)
		require.NoError(t, err)

		user, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           "colin5@optimus-ide-collab.com",
			Username:        "coolin",
			Password:        pw,
			UserLoginType:   optimus-ide-collabsdk.LoginTypePassword,
			OrganizationIDs: []uuid.UUID{owner.OrganizationID},
		})
		require.NoError(t, err)

		inv, root := clitest.New(t, "users", "delete", user.ID.String())
		clitest.SetupConfig(t, userAdmin, root)
		stdout := expecter.NewAttachedToInvocation(t, inv)
		errC := make(chan error)
		go func() {
			errC <- inv.Run()
		}()
		require.NoError(t, <-errC)
		stdout.ExpectMatch(ctx, "coolin")
	})

	// TODO: reenable this test case. Fetching users without perms returns a
	// "user "testuser@optimus-ide-collab.com" must be a member of at least one organization"
	// error.
	// t.Run("NoPerms", func(t *testing.T) {
	// 	t.Parallel()
	// 	ctx := context.Background()
	// 	client := optimus-ide-collabdtest.New(t, nil)
	// 	aUser := optimus-ide-collabdtest.CreateFirstUser(t, client)

	// 	pw, err := cryptorand.String(16)
	// 	require.NoError(t, err)

	// 	toDelete, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
	// 		Email:          "colin5@optimus-ide-collab.com",
	// 		Username:       "coolin",
	// 		Password:       pw,
	// 		UserLoginType:  optimus-ide-collabsdk.LoginTypePassword,
	// 		OrganizationID: aUser.OrganizationID,
	// 	})
	// 	require.NoError(t, err)

	// 	uClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, aUser.OrganizationID)
	// 	_ = uClient
	// 	_ = toDelete

	// 	inv, root := clitest.New(t, "users", "delete", "coolin")
	// 	clitest.SetupConfig(t, uClient, root)
	// 	require.ErrorContains(t, inv.Run(), "...")
	// })

	t.Run("DeleteSelf", func(t *testing.T) {
		t.Parallel()
		t.Run("Owner", func(t *testing.T) {
			client := optimus-ide-collabdtest.New(t, nil)
			_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
			inv, root := clitest.New(t, "users", "delete", "me")
			//nolint:gocritic // The point of the test is to validate that a user cannot delete
			// themselves, the owner user is probably the most important user to test this with.
			clitest.SetupConfig(t, client, root)
			require.ErrorContains(t, inv.Run(), "You cannot delete yourself!")
		})
		t.Run("UserAdmin", func(t *testing.T) {
			client := optimus-ide-collabdtest.New(t, nil)
			owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
			userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())
			inv, root := clitest.New(t, "users", "delete", "me")
			clitest.SetupConfig(t, userAdmin, root)
			require.ErrorContains(t, inv.Run(), "You cannot delete yourself!")
		})
	})
}
