//go:build slim

package cli

import (
	agplcli "github.com/optimus-ide-collab/optimus-ide-collab/v2/cli"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) aiGatewayStart() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "start",
		Short: "Run a standalone AI Gateway server",
		// We accept RawArgs so all commands and flags are accepted.
		RawArgs: true,
		Hidden:  true,
		Handler: func(inv *serpent.Invocation) error {
			agplcli.SlimUnsupported(inv.Stderr, "ai-gateway start")
			return nil
		},
	}

	return cmd
}
