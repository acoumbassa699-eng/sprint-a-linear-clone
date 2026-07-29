package cli

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/pretty"
	"github.com/optimus-ide-collab/serpent"
)

// createUserStatusCommand sets a user status.
func (r *RootCmd) createUserStatusCommand(sdkStatus optimus-ide-collabsdk.UserStatus) *serpent.Command {
	var verb string
	var pastVerb string
	var aliases []string
	var short string
	switch sdkStatus {
	case optimus-ide-collabsdk.UserStatusActive:
		verb = "activate"
		pastVerb = "activated"
		aliases = []string{"active"}
		short = "Update a user's status to 'active'. Active users can fully interact with the platform"
	case optimus-ide-collabsdk.UserStatusSuspended:
		verb = "suspend"
		pastVerb = "suspended"
		short = "Update a user's status to 'suspended'. A suspended user cannot log into the platform"
	default:
		panic(fmt.Sprintf("%s is not supported", sdkStatus))
	}

	var columns []string
	allColumns := []string{"username", "email", "created at", "status"}
	cmd := &serpent.Command{
		Use:     fmt.Sprintf("%s <username|user_id>", verb),
		Short:   short,
		Aliases: aliases,
		Long: FormatExamples(
			Example{
				Command: fmt.Sprintf("optimus-ide-collab users %s example_user", verb),
			},
		),
		Middleware: serpent.Chain(
			serpent.RequireNArgs(1),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}
			identifier := inv.Args[0]
			if identifier == "" {
				return xerrors.Errorf("user identifier cannot be an empty string")
			}

			user, err := client.User(inv.Context(), identifier)
			if err != nil {
				return xerrors.Errorf("fetch user: %w", err)
			}

			// Display the user. This uses cliui.DisplayTable directly instead
			// of cliui.NewOutputFormatter because we prompt immediately
			// afterwards.
			table, err := cliui.DisplayTable([]optimus-ide-collabsdk.User{user}, "", columns)
			if err != nil {
				return xerrors.Errorf("render user table: %w", err)
			}
			_, _ = fmt.Fprintln(inv.Stdout, table)

			// User status is already set to this
			if user.Status == sdkStatus {
				_, _ = fmt.Fprintf(inv.Stdout, "User status is already %q\n", sdkStatus)
				return nil
			}

			// Prompt to confirm the action
			_, err = cliui.Prompt(inv, cliui.PromptOptions{
				Text:      fmt.Sprintf("Are you sure you want to %s this user?", verb),
				IsConfirm: true,
				Default:   cliui.ConfirmYes,
			})
			if err != nil {
				return err
			}

			_, err = client.UpdateUserStatus(inv.Context(), user.ID.String(), sdkStatus)
			if err != nil {
				return xerrors.Errorf("%s user: %w", verb, err)
			}

			_, _ = fmt.Fprintf(inv.Stdout, "\nUser %s has been %s!\n", pretty.Sprint(cliui.DefaultStyles.Keyword, user.Username), pastVerb)
			return nil
		},
	}
	cmd.Options = serpent.OptionSet{
		{
			Flag:          "column",
			FlagShorthand: "c",
			Description:   "Specify a column to filter in the table.",
			Default:       strings.Join(allColumns, ","),
			Value:         serpent.EnumArrayOf(&columns, allColumns...),
		},
	}
	return cmd
}
