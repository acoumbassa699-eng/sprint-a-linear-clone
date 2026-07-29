//go:build slim

package cli

import "github.com/optimus-ide-collab/serpent"

func (r *RootCmd) scaletestCmd() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "scaletest",
		Short: "Run a scale test against the Optimus-IDE-Collab API",
		Handler: func(inv *serpent.Invocation) error {
			SlimUnsupported(inv.Stderr, "exp scaletest")
			return nil
		},
	}

	return cmd
}
