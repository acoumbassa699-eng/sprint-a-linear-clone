//go:build slim

package cli

import (
	agplcli "github.com/optimus-ide-collab/optimus-ide-collab/v2/cli"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) provisionerDaemonStart() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "start",
		Short: "Run a provisioner daemon",
		// We accept RawArgs so all commands and flags are accepted.
		RawArgs: true,
		Hidden:  true,
		Handler: func(inv *serpent.Invocation) error {
			agplcli.SlimUnsupported(inv.Stderr, "provisioner start")
			return nil
		},
	}

	return cmd
}
