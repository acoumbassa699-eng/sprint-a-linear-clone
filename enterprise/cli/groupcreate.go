package cli

import (
	"fmt"

	"golang.org/x/xerrors"

	agpl "github.com/optimus-ide-collab/optimus-ide-collab/v2/cli"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/pretty"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) groupCreate() *serpent.Command {
	var (
		avatarURL   string
		displayName string
		orgContext  = agpl.NewOrganizationContext()
	)

	cmd := &serpent.Command{
		Use:   "create <name>",
		Short: "Create a user group",
		Middleware: serpent.Chain(
			serpent.RequireNArgs(1),
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

			err = optimus-ide-collabsdk.GroupNameValid(inv.Args[0])
			if err != nil {
				return xerrors.Errorf("group name %q is invalid: %w", inv.Args[0], err)
			}

			group, err := client.CreateGroup(ctx, org.ID, optimus-ide-collabsdk.CreateGroupRequest{
				Name:        inv.Args[0],
				DisplayName: displayName,
				AvatarURL:   avatarURL,
			})
			if err != nil {
				return xerrors.Errorf("create group: %w", err)
			}

			_, _ = fmt.Fprintf(inv.Stdout, "Successfully created group %s!\n", pretty.Sprint(cliui.DefaultStyles.Keyword, group.Name))
			return nil
		},
	}

	cmd.Options = serpent.OptionSet{
		{
			Flag:          "avatar-url",
			Description:   `Set an avatar for a group.`,
			FlagShorthand: "u",
			Env:           "OPTIMUS-IDE-COLLAB_AVATAR_URL",
			Value:         serpent.StringOf(&avatarURL),
		},
		{
			Flag:        "display-name",
			Description: `Optional human friendly name for the group.`,
			Env:         "OPTIMUS-IDE-COLLAB_DISPLAY_NAME",
			Value: serpent.Validate(serpent.StringOf(&displayName), func(_displayName *serpent.String) error {
				displayName := _displayName.String()
				if displayName != "" {
					err := optimus-ide-collabsdk.DisplayNameValid(displayName)
					if err != nil {
						return xerrors.Errorf("group display name %q is invalid: %w", displayName, err)
					}
				}
				return nil
			}),
		},
	}
	orgContext.AttachOptions(cmd)

	return cmd
}
