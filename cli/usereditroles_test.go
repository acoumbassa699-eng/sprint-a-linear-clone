package cli_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

var roles = []string{"auditor", "user-admin"}

func TestUserEditRoles(t *testing.T) {
	t.Parallel()

	t.Run("UpdateUserRoles", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleOwner())
		_, member := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleMember())

		inv, root := clitest.New(t, "users", "edit-roles", member.Username, fmt.Sprintf("--roles=%s", strings.Join(roles, ",")))
		clitest.SetupConfig(t, userAdmin, root)

		// Create context with timeout
		ctx := testutil.Context(t, testutil.WaitShort)

		err := inv.WithContext(ctx).Run()
		require.NoError(t, err)

		memberRoles, err := client.UserRoles(ctx, member.Username)
		require.NoError(t, err)

		require.ElementsMatch(t, memberRoles.Roles, roles)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())

		// Setup command with non-existent user
		inv, root := clitest.New(t, "users", "edit-roles", "nonexistentuser")
		clitest.SetupConfig(t, userAdmin, root)

		// Create context with timeout
		ctx := testutil.Context(t, testutil.WaitShort)

		err := inv.WithContext(ctx).Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), "fetch user")
	})
}
