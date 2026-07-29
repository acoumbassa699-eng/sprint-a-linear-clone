package cli

import (
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) users() *serpent.Command {
	cmd := &serpent.Command{
		Short:   "Manage users",
		Use:     "users [subcommand]",
		Aliases: []string{"user"},
		Handler: func(inv *serpent.Invocation) error {
			return inv.Command.HelpHandler(inv)
		},
		Children: []*serpent.Command{
			r.userCreate(),
			r.userList(),
			r.userSingle(),
			r.userDelete(),
			r.userEditRoles(),
			r.userOIDCClaims(),
			r.createUserStatusCommand(optimus-ide-collabsdk.UserStatusActive),
			r.createUserStatusCommand(optimus-ide-collabsdk.UserStatusSuspended),
		},
	}
	return cmd
}
