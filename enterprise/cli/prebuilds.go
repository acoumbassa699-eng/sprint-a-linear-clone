package cli

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) prebuilds() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "prebuilds",
		Short: "Manage Optimus-IDE-Collab prebuilds",
		Long: "Administrators can use these commands to manage prebuilt workspace settings.\n" + cli.FormatExamples(
			cli.Example{
				Description: "Pause Optimus-IDE-Collab prebuilt workspace reconciliation.",
				Command:     "optimus-ide-collab prebuilds pause",
			},
			cli.Example{
				Description: "Resume Optimus-IDE-Collab prebuilt workspace reconciliation if it has been paused.",
				Command:     "optimus-ide-collab prebuilds resume",
			},
		),
		Aliases: []string{"prebuild"},
		Handler: func(inv *serpent.Invocation) error {
			return inv.Command.HelpHandler(inv)
		},
		Children: []*serpent.Command{
			r.pausePrebuilds(),
			r.resumePrebuilds(),
		},
	}
	return cmd
}

func (r *RootCmd) pausePrebuilds() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "pause",
		Short: "Pause prebuilds",
		Middleware: serpent.Chain(
			serpent.RequireNArgs(0),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			err = client.PutPrebuildsSettings(inv.Context(), optimus-ide-collabsdk.PrebuildsSettings{
				ReconciliationPaused: true,
			})
			if err != nil {
				return xerrors.Errorf("unable to pause prebuilds: %w", err)
			}

			_, _ = fmt.Fprintln(inv.Stderr, "Prebuilds are now paused.")
			return nil
		},
	}
	return cmd
}

func (r *RootCmd) resumePrebuilds() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "resume",
		Short: "Resume prebuilds",
		Middleware: serpent.Chain(
			serpent.RequireNArgs(0),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			err = client.PutPrebuildsSettings(inv.Context(), optimus-ide-collabsdk.PrebuildsSettings{
				ReconciliationPaused: false,
			})
			if err != nil {
				return xerrors.Errorf("unable to resume prebuilds: %w", err)
			}

			_, _ = fmt.Fprintln(inv.Stderr, "Prebuilds are now resumed.")
			return nil
		},
	}
	return cmd
}
