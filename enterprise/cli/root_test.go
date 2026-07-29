package cli_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/config"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/cli"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/serpent"
)

func newCLI(t *testing.T, args ...string) (*serpent.Invocation, config.Root) {
	var root cli.RootCmd
	cmd, err := root.Command(root.EnterpriseSubcommands())
	require.NoError(t, err)
	return clitest.NewWithCommand(t, cmd, args...)
}

func TestEnterpriseHandlersOK(t *testing.T) {
	t.Parallel()

	var root cli.RootCmd
	cmd, err := root.Command(root.EnterpriseSubcommands())
	require.NoError(t, err)

	clitest.HandlersOK(t, cmd)
}

func TestCheckWarnings(t *testing.T) {
	t.Parallel()

	t.Run("LicenseWarningForPrivilegedRoles", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				ExpiresAt: time.Now().Add(time.Hour * 24),
			},
		})

		inv, conf := newCLI(t, "list")

		var buf bytes.Buffer
		inv.Stderr = &buf
		clitest.SetupConfig(t, client, conf) //nolint:gocritic // owners should see this

		err := inv.Run()
		require.NoError(t, err)

		require.Contains(t, buf.String(), "Your license expires in 1 day.")
	})

	t.Run("NoLicenseWarningForRegularUser", func(t *testing.T) {
		t.Parallel()
		adminClient, admin := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				ExpiresAt: time.Now().Add(time.Hour * 24),
			},
		})

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, admin.OrganizationID)

		inv, conf := newCLI(t, "list")

		var buf bytes.Buffer
		inv.Stderr = &buf
		clitest.SetupConfig(t, client, conf)

		err := inv.Run()
		require.NoError(t, err)

		require.NotContains(t, buf.String(), "Your license expires")
	})
}
