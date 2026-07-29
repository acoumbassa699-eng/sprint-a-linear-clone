package exptest_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"cdr.dev/slog/v3/sloggers/slogtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// This test validates that the scaletest CLI filters out workspaces not owned
// when disable owner workspace access is set.
// This test is in its own package because it mutates a global variable that
// can influence other tests in the same package.
// nolint:paralleltest
func TestScaleTestWorkspaceTraffic_UseHostLogin(t *testing.T) {
	log := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true})
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		Logger:                   &log,
		IncludeProvisionerDaemon: true,
		DeploymentValues: optimus-ide-collabdtest.DeploymentValues(t, func(dv *optimus-ide-collabsdk.DeploymentValues) {
			dv.DisableOwnerWorkspaceExec = true
		}),
	})
	owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
	tv := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
	_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, tv.ID)
	tpl := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, tv.ID)
	// Create a workspace owned by a different user
	memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
	_ = optimus-ide-collabdtest.CreateWorkspace(t, memberClient, tpl.ID, func(cwr *optimus-ide-collabsdk.CreateWorkspaceRequest) {
		cwr.Name = "scaletest-workspace"
	})

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	// Test without --use-host-login first.g
	inv, root := clitest.New(t, "exp", "scaletest", "workspace-traffic",
		"--template", tpl.Name,
	)
	// nolint:gocritic // We are intentionally testing this as the owner.
	clitest.SetupConfig(t, client, root)
	var stdoutBuf bytes.Buffer
	inv.Stdout = &stdoutBuf

	err := inv.WithContext(ctx).Run()
	require.ErrorContains(t, err, "no scaletest workspaces exist")
	require.Contains(t, stdoutBuf.String(), `1 workspace(s) were skipped`)

	// Test once again with --use-host-login.
	inv, root = clitest.New(t, "exp", "scaletest", "workspace-traffic",
		"--template", tpl.Name,
		"--use-host-login",
	)
	// nolint:gocritic // We are intentionally testing this as the owner.
	clitest.SetupConfig(t, client, root)
	stdoutBuf.Reset()
	inv.Stdout = &stdoutBuf

	err = inv.WithContext(ctx).Run()
	require.ErrorContains(t, err, "no scaletest workspaces exist")
	require.NotContains(t, stdoutBuf.String(), `1 workspace(s) were skipped`)
}
