package cli

import (
	"fmt"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) listOrganizations() *serpent.Command {
	formatter := cliui.NewOutputFormatter(
		cliui.TableFormat([]optimus-ide-collabsdk.Organization{}, []string{"name", "display name", "id", "default"}),
		cliui.JSONFormat(),
	)

	cmd := &serpent.Command{
		Use:     "list",
		Short:   "List all organizations",
		Long:    "List all organizations. Requires a role which grants ResourceOrganization: read.",
		Aliases: []string{"ls"},
		Middleware: serpent.Chain(
			serpent.RequireNArgs(0),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			organizations, err := client.Organizations(inv.Context())
			if err != nil {
				return err
			}

			out, err := formatter.Format(inv.Context(), organizations)
			if err != nil {
				return err
			}

			if out == "" {
				cliui.Infof(inv.Stderr, "No organizations found.")
				return nil
			}

			_, err = fmt.Fprintln(inv.Stdout, out)
			return err
		},
	}

	formatter.AttachOptions(&cmd.Options)
	return cmd
}
