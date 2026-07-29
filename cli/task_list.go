package cli

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

type taskListRow struct {
	Task optimus-ide-collabsdk.Task `table:"t,recursive_inline"`

	StateChangedAgo string `table:"state changed"`
}

func taskListRowFromTask(now time.Time, t optimus-ide-collabsdk.Task) taskListRow {
	var stateAgo string
	if t.CurrentState != nil {
		stateAgo = now.UTC().Sub(t.CurrentState.Timestamp).Truncate(time.Second).String() + " ago"
	}

	return taskListRow{
		Task: t,

		StateChangedAgo: stateAgo,
	}
}

func (r *RootCmd) taskList() *serpent.Command {
	var (
		statusFilter string
		all          bool
		user         string
		quiet        bool

		formatter = cliui.NewOutputFormatter(
			cliui.TableFormat(
				[]taskListRow{},
				[]string{
					"name",
					"status",
					"state",
					"state changed",
					"message",
				},
			),
			cliui.ChangeFormatterData(
				cliui.JSONFormat(),
				func(data any) (any, error) {
					rows, ok := data.([]taskListRow)
					if !ok {
						return nil, xerrors.Errorf("expected []taskListRow, got %T", data)
					}
					out := make([]optimus-ide-collabsdk.Task, len(rows))
					for i := range rows {
						out[i] = rows[i].Task
					}
					return out, nil
				},
			),
		)
	)

	cmd := &serpent.Command{
		Use:   "list",
		Short: "List tasks",
		Long: FormatExamples(
			Example{
				Description: "List tasks for the current user.",
				Command:     "optimus-ide-collab task list",
			},
			Example{
				Description: "List tasks for a specific user.",
				Command:     "optimus-ide-collab task list --user someone-else",
			},
			Example{
				Description: "List all tasks you can view.",
				Command:     "optimus-ide-collab task list --all",
			},
			Example{
				Description: "List all your running tasks.",
				Command:     "optimus-ide-collab task list --status running",
			},
			Example{
				Description: "As above, but only show IDs.",
				Command:     "optimus-ide-collab task list --status running --quiet",
			},
		),
		Aliases: []string{"ls"},
		Middleware: serpent.Chain(
			serpent.RequireNArgs(0),
		),
		Options: serpent.OptionSet{
			{
				Name:        "status",
				Description: "Filter by task status.",
				Flag:        "status",
				Default:     "",
				Value:       serpent.EnumOf(&statusFilter, slice.ToStrings(optimus-ide-collabsdk.AllTaskStatuses())...),
			},
			{
				Name:          "all",
				Description:   "List tasks for all users you can view.",
				Flag:          "all",
				FlagShorthand: "a",
				Default:       "false",
				Value:         serpent.BoolOf(&all),
			},
			{
				Name:        "user",
				Description: "List tasks for the specified user (username, \"me\").",
				Flag:        "user",
				Default:     "",
				Value:       serpent.StringOf(&user),
			},
			{
				Name:          "quiet",
				Description:   "Only display task IDs.",
				Flag:          "quiet",
				FlagShorthand: "q",
				Default:       "false",
				Value:         serpent.BoolOf(&quiet),
			},
		},
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			ctx := inv.Context()

			targetUser := strings.TrimSpace(user)
			if targetUser == "" && !all {
				targetUser = optimus-ide-collabsdk.Me
			}

			tasks, err := client.Tasks(ctx, &optimus-ide-collabsdk.TasksFilter{
				Owner:  targetUser,
				Status: optimus-ide-collabsdk.TaskStatus(statusFilter),
			})
			if err != nil {
				return xerrors.Errorf("list tasks: %w", err)
			}

			if quiet {
				for _, task := range tasks {
					_, _ = fmt.Fprintln(inv.Stdout, task.ID.String())
				}

				return nil
			}

			rows := make([]taskListRow, len(tasks))
			now := time.Now()
			for i := range tasks {
				rows[i] = taskListRowFromTask(now, tasks[i])
			}

			out, err := formatter.Format(ctx, rows)
			if err != nil {
				return xerrors.Errorf("format tasks: %w", err)
			}
			if out == "" {
				cliui.Infof(inv.Stderr, "No tasks found.")
				return nil
			}
			_, _ = fmt.Fprintln(inv.Stdout, out)
			return nil
		},
	}

	formatter.AttachOptions(&cmd.Options)
	return cmd
}
