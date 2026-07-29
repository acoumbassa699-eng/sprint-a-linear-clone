package cli_test

import (
	"bytes"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestUpdateGroupSync(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		owner, first := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})

		ctx := testutil.Context(t, testutil.WaitLong)
		inv, root := clitest.New(t, "organization", "settings", "set", "groupsync")
		//nolint:gocritic // Using the owner, testing the cli not perms
		clitest.SetupConfig(t, owner, root)

		expectedSettings := optimus-ide-collabsdk.GroupSyncSettings{
			Field: "groups",
			Mapping: map[string][]uuid.UUID{
				"test": {first.OrganizationID},
			},
			RegexFilter:       regexp.MustCompile("^foo"),
			AutoCreateMissing: true,
			LegacyNameMapping: nil,
		}
		expectedData, err := json.Marshal(expectedSettings)
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		inv.Stdout = buf
		inv.Stdin = bytes.NewBuffer(expectedData)
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)
		require.JSONEq(t, string(expectedData), buf.String())

		// Now read it back
		inv, root = clitest.New(t, "organization", "settings", "show", "groupsync")
		//nolint:gocritic // Using the owner, testing the cli not perms
		clitest.SetupConfig(t, owner, root)

		buf = new(bytes.Buffer)
		inv.Stdout = buf
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)
		require.JSONEq(t, string(expectedData), buf.String())
	})
}

func TestUpdateRoleSync(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		owner, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})

		ctx := testutil.Context(t, testutil.WaitLong)
		inv, root := clitest.New(t, "organization", "settings", "set", "rolesync")
		//nolint:gocritic // Using the owner, testing the cli not perms
		clitest.SetupConfig(t, owner, root)

		expectedSettings := optimus-ide-collabsdk.RoleSyncSettings{
			Field: "roles",
			Mapping: map[string][]string{
				"test": {rbac.RoleOrgAdmin()},
			},
		}
		expectedData, err := json.Marshal(expectedSettings)
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		inv.Stdout = buf
		inv.Stdin = bytes.NewBuffer(expectedData)
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)
		require.JSONEq(t, string(expectedData), buf.String())

		// Now read it back
		inv, root = clitest.New(t, "organization", "settings", "show", "rolesync")
		//nolint:gocritic // Using the owner, testing the cli not perms
		clitest.SetupConfig(t, owner, root)

		buf = new(bytes.Buffer)
		inv.Stdout = buf
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)
		require.JSONEq(t, string(expectedData), buf.String())
	})
}

func TestUpdateOrganizationSync(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		owner, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})

		ctx := testutil.Context(t, testutil.WaitLong)
		inv, root := clitest.New(t, "organization", "settings", "set", "organization-sync")
		//nolint:gocritic // Using the owner, testing the cli not perms
		clitest.SetupConfig(t, owner, root)

		expectedSettings := optimus-ide-collabsdk.OrganizationSyncSettings{
			Field: "organizations",
			Mapping: map[string][]uuid.UUID{
				"test": {uuid.New()},
			},
		}
		expectedData, err := json.Marshal(expectedSettings)
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		inv.Stdout = buf
		inv.Stdin = bytes.NewBuffer(expectedData)
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)
		require.JSONEq(t, string(expectedData), buf.String())

		// Now read it back
		inv, root = clitest.New(t, "organization", "settings", "show", "organization-sync")
		//nolint:gocritic // Using the owner, testing the cli not perms
		clitest.SetupConfig(t, owner, root)

		buf = new(bytes.Buffer)
		inv.Stdout = buf
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)
		require.JSONEq(t, string(expectedData), buf.String())
	})
}
