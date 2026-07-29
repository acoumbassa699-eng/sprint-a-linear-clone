//go:build slim

package cli

import (
	agplcli "github.com/optimus-ide-collab/optimus-ide-collab/v2/cli"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) proxyServer() *serpent.Command {
	root := &serpent.Command{
		Use:     "server",
		Short:   "Start a workspace proxy server",
		Aliases: []string{},
		// We accept RawArgs so all commands and flags are accepted.
		RawArgs: true,
		Hidden:  true,
		Handler: func(inv *serpent.Invocation) error {
			agplcli.SlimUnsupported(inv.Stderr, "workspace-proxy server")
			return nil
		},
	}

	return root
}
