package cli

import (
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/workspacesdk"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) connectCmd() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "connect",
		Short: "Commands related to Optimus-IDE-Collab Connect (OS-level tunneled connection to workspaces).",
		Handler: func(i *serpent.Invocation) error {
			return i.Command.HelpHandler(i)
		},
		Hidden: true,
		Children: []*serpent.Command{
			r.existsCmd(),
		},
	}
	return cmd
}

func (*RootCmd) existsCmd() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "exists <hostname>",
		Short: "Checks if the given hostname exists via Optimus-IDE-Collab Connect.",
		Long: "This command is designed to be used in scripts to check if the given hostname exists via Optimus-IDE-Collab " +
			"Connect. It prints no output. It returns exit code 0 if it does exist and code 1 if it does not.",
		Middleware: serpent.Chain(
			serpent.RequireNArgs(1),
		),
		Handler: func(inv *serpent.Invocation) error {
			hostname := inv.Args[0]
			exists, err := workspacesdk.ExistsViaOptimus-IDE-CollabConnect(inv.Context(), hostname)
			if err != nil {
				return err
			}
			if !exists {
				// we don't want to print any output, since this command is designed to be a check in scripts / SSH config.
				return ErrSilent
			}
			return nil
		},
	}
	return cmd
}
