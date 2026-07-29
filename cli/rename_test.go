package cli_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
)

func TestRename(t *testing.T) {
	t.Parallel()
	logger := testutil.Logger(t)

	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true, AllowWorkspaceRenames: true})
	owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
	member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, member, template.ID)
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	want := optimus-ide-collabdtest.RandomUsername(t)
	inv, root := clitest.New(t, "rename", workspace.Name, want, "--yes")
	clitest.SetupConfig(t, member, root)
	stdout := expecter.NewAttachedToInvocation(t, inv)
	stdin := testutil.NewWriterAttachedToInvocation(t, logger.Named("stdin"), inv)
	clitest.Start(t, inv)

	stdout.ExpectMatch(ctx, "confirm rename:")
	stdin.WriteLine(workspace.Name)
	stdout.ExpectMatch(ctx, "renamed to")

	ws, err := client.Workspace(ctx, workspace.ID)
	assert.NoError(t, err)

	got := ws.Name
	assert.Equal(t, want, got, "workspace name did not change")
}
