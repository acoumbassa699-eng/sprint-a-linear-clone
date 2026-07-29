package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
)

func TestList(t *testing.T) {
	t.Parallel()
	t.Run("Single", func(t *testing.T) {
		t.Parallel()
		client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, memberUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		// setup template
		r := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OrganizationID: owner.OrganizationID,
			OwnerID:        memberUser.ID,
		}).WithAgent().Do()

		inv, root := clitest.New(t, "ls")
		clitest.SetupConfig(t, member, root)
		stdout := expecter.NewAttachedToInvocation(t, inv)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()
		done := make(chan any)
		go func() {
			errC := inv.WithContext(ctx).Run()
			assert.NoError(t, errC)
			close(done)
		}()
		stdout.ExpectMatch(ctx, r.Workspace.Name)
		stdout.ExpectMatch(ctx, "Started")
		cancelFunc()
		<-done
	})

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()
		client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, memberUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		_ = dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OrganizationID: owner.OrganizationID,
			OwnerID:        memberUser.ID,
		}).WithAgent().Do()

		inv, root := clitest.New(t, "list", "--output=json")
		clitest.SetupConfig(t, member, root)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()

		out := bytes.NewBuffer(nil)
		inv.Stdout = out
		err := inv.WithContext(ctx).Run()
		require.NoError(t, err)

		var workspaces []optimus-ide-collabsdk.Workspace
		require.NoError(t, json.Unmarshal(out.Bytes(), &workspaces))
		require.Len(t, workspaces, 1)
	})

	t.Run("NoWorkspacesJSON", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		inv, root := clitest.New(t, "list", "--output=json")
		clitest.SetupConfig(t, member, root)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()

		stdout := bytes.NewBuffer(nil)
		stderr := bytes.NewBuffer(nil)
		inv.Stdout = stdout
		inv.Stderr = stderr
		err := inv.WithContext(ctx).Run()
		require.NoError(t, err)

		var workspaces []optimus-ide-collabsdk.Workspace
		require.NoError(t, json.Unmarshal(stdout.Bytes(), &workspaces))
		require.Len(t, workspaces, 0)

		require.Len(t, stderr.Bytes(), 0)
	})

	t.Run("SharedWorkspaces", func(t *testing.T) {
		t.Parallel()

		var (
			client, db           = optimus-ide-collabdtest.NewWithDatabase(t, nil)
			orgOwner             = optimus-ide-collabdtest.CreateFirstUser(t, client)
			memberClient, member = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID, rbac.ScopedRoleOrgAuditor(orgOwner.OrganizationID))
			sharedWorkspace      = dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
				Name:           "wibble",
				OwnerID:        orgOwner.UserID,
				OrganizationID: orgOwner.OrganizationID,
			}).Do().Workspace
			_ = dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
				Name:           "wobble",
				OwnerID:        orgOwner.UserID,
				OrganizationID: orgOwner.OrganizationID,
			}).Do().Workspace
		)

		ctx := testutil.Context(t, testutil.WaitMedium)

		client.UpdateWorkspaceACL(ctx, sharedWorkspace.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				member.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})

		inv, root := clitest.New(t, "list", "--shared-with-me", "--output=json")
		clitest.SetupConfig(t, memberClient, root)

		stdout := new(bytes.Buffer)
		inv.Stdout = stdout
		err := inv.WithContext(ctx).Run()
		require.NoError(t, err)

		var workspaces []optimus-ide-collabsdk.Workspace
		require.NoError(t, json.Unmarshal(stdout.Bytes(), &workspaces))
		require.Len(t, workspaces, 1)
		require.Equal(t, sharedWorkspace.ID, workspaces[0].ID)
	})
}
