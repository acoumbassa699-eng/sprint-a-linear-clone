package cli_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
	"github.com/optimus-ide-collab/pretty"
)

func TestGroupEdit(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		client, admin := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
			},
		}})
		anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID, rbac.RoleUserAdmin())

		_, user1 := optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID)
		_, user2 := optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID)
		_, user3 := optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID)

		group := optimus-ide-collabdtest.CreateGroup(t, client, admin.OrganizationID, "alpha", user3)

		expectedName := "beta"

		inv, conf := newCLI(
			t,
			"groups", "edit", group.Name,
			"--name", expectedName,
			"--avatar-url", "https://example.com",
			"-a", user1.ID.String(),
			"-a", user2.Email,
			"-r", user3.ID.String(),
		)

		stdout := expecter.NewAttachedToInvocation(t, inv)
		clitest.SetupConfig(t, anotherClient, conf)
		ctx := testutil.Context(t, testutil.WaitMedium)

		err := inv.Run()
		require.NoError(t, err)

		stdout.ExpectMatch(ctx, fmt.Sprintf("Successfully patched group %s", pretty.Sprint(cliui.DefaultStyles.Keyword, expectedName)))
	})

	t.Run("InvalidUserInput", func(t *testing.T) {
		t.Parallel()

		client, admin := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
			},
		}})

		// Create a group with no members.
		group := optimus-ide-collabdtest.CreateGroup(t, client, admin.OrganizationID, "alpha")

		inv, conf := newCLI(
			t,
			"groups", "edit", group.Name,
			"-a", "foo",
		)

		clitest.SetupConfig(t, client, conf) //nolint:gocritic // intentional usage of owner

		err := inv.Run()
		require.ErrorContains(t, err, "must be a valid UUID or email address")
	})

	t.Run("NoArg", func(t *testing.T) {
		t.Parallel()

		client, user := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
			},
		}})
		anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, user.OrganizationID, rbac.RoleUserAdmin())

		inv, conf := newCLI(t, "groups", "edit")
		clitest.SetupConfig(t, anotherClient, conf)

		err := inv.Run()
		require.ErrorContains(t, err, "wanted 1 args but got 0")
	})
}
