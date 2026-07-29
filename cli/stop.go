package cli

import (
	"fmt"
	"time"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) stop() *serpent.Command {
	var bflags buildFlags
	cmd := &serpent.Command{
		Annotations: workspaceCommand,
		Use:         "stop <workspace>",
		Short:       "Stop a workspace",
		Middleware: serpent.Chain(
			serpent.RequireNArgs(1),
		),
		Options: serpent.OptionSet{
			cliui.SkipPromptOption(),
		},
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			_, err = cliui.Prompt(inv, cliui.PromptOptions{
				Text:      "Confirm stop workspace?",
				IsConfirm: true,
			})
			if err != nil {
				return err
			}

			workspace, err := client.ResolveWorkspace(inv.Context(), inv.Args[0])
			if err != nil {
				return err
			}

			build, err := stopWorkspace(inv, client, workspace, bflags)
			if err != nil {
				return err
			}

			err = cliui.WorkspaceBuild(inv.Context(), inv.Stdout, client, build.ID)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(
				inv.Stdout,
				"\nThe %s workspace has been stopped at %s!\n",
				cliui.Keyword(workspace.Name),
				cliui.Timestamp(time.Now()),
			)
			return nil
		},
	}
	cmd.Options = append(cmd.Options, bflags.cliOptions()...)

	return cmd
}

func stopWorkspace(inv *serpent.Invocation, client *optimus-ide-collabsdk.Client, workspace optimus-ide-collabsdk.Workspace, bflags buildFlags) (optimus-ide-collabsdk.WorkspaceBuild, error) {
	if workspace.LatestBuild.Job.Status == optimus-ide-collabsdk.ProvisionerJobPending {
		// cliutil.WarnMatchedProvisioners also checks if the job is pending
		// but we still want to avoid users spamming multiple builds that will
		// not be picked up.
		cliui.Warn(inv.Stderr, "The workspace is already stopping!")
		cliutil.WarnMatchedProvisioners(inv.Stderr, workspace.LatestBuild.MatchedProvisioners, workspace.LatestBuild.Job)
		if _, err := cliui.Prompt(inv, cliui.PromptOptions{
			Text:      "Enqueue another stop?",
			IsConfirm: true,
			Default:   cliui.ConfirmNo,
		}); err != nil {
			return optimus-ide-collabsdk.WorkspaceBuild{}, err
		}
	}
	wbr := optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
	}
	if bflags.provisionerLogDebug {
		wbr.LogLevel = optimus-ide-collabsdk.ProvisionerLogLevelDebug
	}
	return client.CreateWorkspaceBuild(inv.Context(), workspace.ID, wbr)
}
