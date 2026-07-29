package cli

import (
	"fmt"
	"time"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) Provisioners() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "provisioner",
		Short: "View and manage provisioner daemons and jobs",
		Handler: func(inv *serpent.Invocation) error {
			return inv.Command.HelpHandler(inv)
		},
		Aliases: []string{"provisioners"},
		Children: []*serpent.Command{
			r.provisionerList(),
			r.provisionerJobs(),
		},
	}

	return cmd
}

func (r *RootCmd) provisionerList() *serpent.Command {
	type provisionerDaemonRow struct {
		optimus-ide-collabsdk.ProvisionerDaemon `table:"provisioner_daemon,recursive_inline"`
		OrganizationName           string `json:"organization_name" table:"organization"`
	}
	var (
		orgContext = NewOrganizationContext()
		formatter  = cliui.NewOutputFormatter(
			cliui.TableFormat([]provisionerDaemonRow{}, []string{"created at", "last seen at", "key name", "name", "version", "status", "tags"}),
			cliui.JSONFormat(),
		)
		limit   int64
		offline bool
		status  []string
		maxAge  time.Duration
	)

	cmd := &serpent.Command{
		Use:     "list",
		Short:   "List provisioner daemons in an organization",
		Aliases: []string{"ls"},
		Middleware: serpent.Chain(
			serpent.RequireNArgs(0),
		),
		Handler: func(inv *serpent.Invocation) error {
			ctx := inv.Context()
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}
			org, err := orgContext.Selected(inv, client)
			if err != nil {
				return xerrors.Errorf("current organization: %w", err)
			}

			daemons, err := client.OrganizationProvisionerDaemons(ctx, org.ID, &optimus-ide-collabsdk.OrganizationProvisionerDaemonsOptions{
				Limit:   int(limit),
				Offline: offline,
				Status:  slice.StringEnums[optimus-ide-collabsdk.ProvisionerDaemonStatus](status),
				MaxAge:  maxAge,
			})
			if err != nil {
				return xerrors.Errorf("list provisioner daemons: %w", err)
			}

			var rows []provisionerDaemonRow
			for _, daemon := range daemons {
				rows = append(rows, provisionerDaemonRow{
					ProvisionerDaemon: daemon,
					OrganizationName:  org.HumanName(),
				})
			}

			out, err := formatter.Format(ctx, rows)
			if err != nil {
				return xerrors.Errorf("display provisioner daemons: %w", err)
			}

			if out == "" {
				cliui.Infof(inv.Stderr, "No provisioner daemons found.")
				return nil
			}

			_, _ = fmt.Fprintln(inv.Stdout, out)

			return nil
		},
	}

	cmd.Options = append(cmd.Options, []serpent.Option{
		{
			Flag:          "limit",
			FlagShorthand: "l",
			Env:           "OPTIMUS-IDE-COLLAB_PROVISIONER_LIST_LIMIT",
			Description:   "Limit the number of provisioners returned.",
			Default:       "50",
			Value:         serpent.Int64Of(&limit),
		},
		{
			Flag:          "show-offline",
			FlagShorthand: "f",
			Env:           "OPTIMUS-IDE-COLLAB_PROVISIONER_SHOW_OFFLINE",
			Description:   "Show offline provisioners.",
			Value:         serpent.BoolOf(&offline),
		},
		{
			Flag:          "status",
			FlagShorthand: "s",
			Env:           "OPTIMUS-IDE-COLLAB_PROVISIONER_LIST_STATUS",
			Description:   "Filter by provisioner status.",
			Value:         serpent.EnumArrayOf(&status, slice.ToStrings(optimus-ide-collabsdk.ProvisionerDaemonStatusEnums())...),
		},
		{
			Flag:          "max-age",
			FlagShorthand: "m",
			Env:           "OPTIMUS-IDE-COLLAB_PROVISIONER_LIST_MAX_AGE",
			Description:   "Filter provisioners by maximum age.",
			Value:         serpent.DurationOf(&maxAge),
		},
	}...)

	orgContext.AttachOptions(cmd)
	formatter.AttachOptions(&cmd.Options)

	return cmd
}
