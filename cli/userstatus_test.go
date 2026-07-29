package cli_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestUserStatus(t *testing.T) {
	t.Parallel()

	t.Run("StatusSelf", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		optimus-ide-collabdtest.CreateFirstUser(t, client)

		inv, root := clitest.New(t, "users", "suspend", "me")
		clitest.SetupConfig(t, client, root)
		// Yes to the prompt
		inv.Stdin = bytes.NewReader([]byte("yes\n"))
		err := inv.Run()
		// Expect an error, as you cannot suspend yourself
		require.Error(t, err)
		require.ErrorContains(t, err, "cannot suspend yourself")
	})

	t.Run("StatusOther", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())
		other, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		otherUser, err := other.User(context.Background(), optimus-ide-collabsdk.Me)
		require.NoError(t, err, "fetch user")

		inv, root := clitest.New(t, "users", "suspend", otherUser.Username)
		clitest.SetupConfig(t, userAdmin, root)
		// Yes to the prompt
		inv.Stdin = bytes.NewReader([]byte("yes\n"))
		err = inv.Run()
		require.NoError(t, err, "suspend user")

		// Check the user status
		otherUser, err = client.User(context.Background(), otherUser.Username)
		require.NoError(t, err, "fetch suspended user")
		require.Equal(t, optimus-ide-collabsdk.UserStatusSuspended, otherUser.Status, "suspended user")

		// Set back to active. Try using a uuid as well
		inv, root = clitest.New(t, "users", "activate", otherUser.ID.String())
		clitest.SetupConfig(t, userAdmin, root)
		// Yes to the prompt
		inv.Stdin = bytes.NewReader([]byte("yes\n"))
		err = inv.Run()
		require.NoError(t, err, "suspend user")

		// Check the user status
		otherUser, err = client.User(context.Background(), otherUser.ID.String())
		require.NoError(t, err, "fetch active user")
		require.Equal(t, optimus-ide-collabsdk.UserStatusActive, otherUser.Status, "active user")
	})
}
