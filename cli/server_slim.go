//go:build slim

package cli

import "github.com/optimus-ide-collab/serpent"

func (r *RootCmd) Server(_ func()) *serpent.Command {
	root := &serpent.Command{
		Use:   "server",
		Short: "Start a Optimus-IDE-Collab server",
		// We accept RawArgs so all commands and flags are accepted.
		RawArgs: true,
		Hidden:  true,
		Handler: func(inv *serpent.Invocation) error {
			SlimUnsupported(inv.Stderr, "server")
			return nil
		},
	}

	return root
}
