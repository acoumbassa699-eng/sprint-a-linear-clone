package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
)

func TestFavoriteUnfavorite(t *testing.T) {
	t.Parallel()

	var (
		client, db           = optimus-ide-collabdtest.NewWithDatabase(t, nil)
		owner                = optimus-ide-collabdtest.CreateFirstUser(t, client)
		memberClient, member = optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		ws                   = dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{OwnerID: member.ID, OrganizationID: owner.OrganizationID}).Do()
	)

	inv, root := clitest.New(t, "favorite", ws.Workspace.Name)
	clitest.SetupConfig(t, memberClient, root)

	var buf bytes.Buffer
	inv.Stdout = &buf
	err := inv.Run()
	require.NoError(t, err)

	updated := optimus-ide-collabdtest.MustWorkspace(t, memberClient, ws.Workspace.ID)
	require.True(t, updated.Favorite)

	buf.Reset()

	inv, root = clitest.New(t, "unfavorite", ws.Workspace.Name)
	clitest.SetupConfig(t, memberClient, root)
	inv.Stdout = &buf
	err = inv.Run()
	require.NoError(t, err)
	updated = optimus-ide-collabdtest.MustWorkspace(t, memberClient, ws.Workspace.ID)
	require.False(t, updated.Favorite)
}
