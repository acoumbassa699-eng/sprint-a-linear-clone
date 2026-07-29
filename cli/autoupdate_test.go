package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestAutoUpdate(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)
		workspace := optimus-ide-collabdtest.CreateWorkspace(t, member, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)
		require.Equal(t, optimus-ide-collabsdk.AutomaticUpdatesNever, workspace.AutomaticUpdates)

		expectedPolicy := optimus-ide-collabsdk.AutomaticUpdatesAlways
		inv, root := clitest.New(t, "autoupdate", workspace.Name, string(expectedPolicy))
		clitest.SetupConfig(t, member, root)
		var buf bytes.Buffer
		inv.Stdout = &buf
		err := inv.Run()
		require.NoError(t, err)
		require.Contains(t, buf.String(), fmt.Sprintf("Updated workspace %q auto-update policy to %q", workspace.Name, expectedPolicy))

		workspace = optimus-ide-collabdtest.MustWorkspace(t, client, workspace.ID)
		require.Equal(t, expectedPolicy, workspace.AutomaticUpdates)
	})

	t.Run("InvalidArgs", func(t *testing.T) {
		type testcase struct {
			Name          string
			Args          []string
			ErrorContains string
		}

		cases := []testcase{
			{
				Name:          "NoPolicy",
				Args:          []string{"autoupdate", "ws"},
				ErrorContains: "wanted 2 args but got 1",
			},
			{
				Name:          "InvalidPolicy",
				Args:          []string{"autoupdate", "ws", "sometimes"},
				ErrorContains: `invalid option "sometimes" must be either of`,
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				t.Parallel()
				client := optimus-ide-collabdtest.New(t, nil)
				_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

				inv, root := clitest.New(t, c.Args...)
				clitest.SetupConfig(t, client, root)
				err := inv.Run()
				require.Error(t, err)
				require.Contains(t, err.Error(), c.ErrorContains)
			})
		}
	})
}
