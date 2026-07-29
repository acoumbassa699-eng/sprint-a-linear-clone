package cli

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

const defaultGroupDisplay = "-"

func (r *RootCmd) sharing() *serpent.Command {
	cmd := &serpent.Command{
		Use:     "sharing [subcommand]",
		Short:   "Commands for managing shared workspaces",
		Aliases: []string{"share"},
		Handler: func(inv *serpent.Invocation) error {
			return inv.Command.HelpHandler(inv)
		},
		Children: []*serpent.Command{
			r.shareWorkspace(),
			r.unshareWorkspace(),
			r.statusWorkspaceSharing(),
		},
		Hidden: true,
	}

	return cmd
}

func (r *RootCmd) statusWorkspaceSharing() *serpent.Command {
	cmd := &serpent.Command{
		Use:     "status <workspace>",
		Short:   "List all users and groups the given Workspace is shared with.",
		Aliases: []string{"list"},
		Middleware: serpent.Chain(
			serpent.RequireNArgs(1),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			workspace, err := client.ResolveWorkspace(inv.Context(), inv.Args[0])
			if err != nil {
				return xerrors.Errorf("unable to fetch Workspace %s: %w", inv.Args[0], err)
			}

			acl, err := client.WorkspaceACL(inv.Context(), workspace.ID)
			if err != nil {
				return xerrors.Errorf("unable to fetch ACL for Workspace: %w", err)
			}

			out, err := workspaceACLToTable(inv.Context(), &acl)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(inv.Stdout, out)
			return err
		},
	}

	return cmd
}

func (r *RootCmd) shareWorkspace() *serpent.Command {
	var (
		users  []string
		groups []string

		// Username regex taken from optimus-ide-collabsdk/name.go
		nameRoleRegex = regexp.MustCompile(`(^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*)+(?::([A-Za-z0-9-]+))?`)
	)

	cmd := &serpent.Command{
		Use:     "add <workspace> --user <user>:<role> --group <group>:<role>",
		Aliases: []string{"share"},
		Short:   "Share a workspace with a user or group.",
		Options: serpent.OptionSet{
			{
				Name:        "user",
				Description: "A comma separated list of users to share the workspace with.",
				Flag:        "user",
				Value:       serpent.StringArrayOf(&users),
			}, {
				Name:        "group",
				Description: "A comma separated list of groups to share the workspace with.",
				Flag:        "group",
				Value:       serpent.StringArrayOf(&groups),
			},
		},
		Middleware: serpent.Chain(
			serpent.RequireNArgs(1),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			if len(users) == 0 && len(groups) == 0 {
				return xerrors.New("at least one user or group must be provided")
			}

			workspace, err := client.ResolveWorkspace(inv.Context(), inv.Args[0])
			if err != nil {
				return xerrors.Errorf("could not fetch the workspace %s: %w", inv.Args[0], err)
			}

			userRoleStrings := make([][2]string, len(users))
			for index, user := range users {
				userAndRole := nameRoleRegex.FindStringSubmatch(user)
				if userAndRole == nil {
					return xerrors.Errorf("invalid user format %q: must match pattern 'username:role'", user)
				}

				userRoleStrings[index] = [2]string{userAndRole[1], userAndRole[2]}
			}

			groupRoleStrings := make([][2]string, len(groups))
			for index, group := range groups {
				groupAndRole := nameRoleRegex.FindStringSubmatch(group)
				if groupAndRole == nil {
					return xerrors.Errorf("invalid group format %q: must match pattern 'group:role'", group)
				}

				groupRoleStrings[index] = [2]string{groupAndRole[1], groupAndRole[2]}
			}

			userRoles, groupRoles, err := fetchUsersAndGroups(inv.Context(), fetchUsersAndGroupsParams{
				Client:      client,
				OrgID:       workspace.OrganizationID,
				OrgName:     workspace.OrganizationName,
				Users:       userRoleStrings,
				Groups:      groupRoleStrings,
				DefaultRole: optimus-ide-collabsdk.WorkspaceRoleUse,
			})
			if err != nil {
				return err
			}

			err = client.UpdateWorkspaceACL(inv.Context(), workspace.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
				UserRoles:  userRoles,
				GroupRoles: groupRoles,
			})
			if err != nil {
				return err
			}

			acl, err := client.WorkspaceACL(inv.Context(), workspace.ID)
			if err != nil {
				return xerrors.Errorf("could not fetch current workspace ACL after sharing %w", err)
			}

			out, err := workspaceACLToTable(inv.Context(), &acl)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(inv.Stdout, out)
			return err
		},
	}

	return cmd
}

func (r *RootCmd) unshareWorkspace() *serpent.Command {
	var (
		users  []string
		groups []string
	)

	cmd := &serpent.Command{
		Use:     "remove <workspace> --user <user> --group <group>",
		Aliases: []string{"unshare"},
		Short:   "Remove shared access for users or groups from a workspace.",
		Options: serpent.OptionSet{
			{
				Name:        "user",
				Description: "A comma separated list of users to share the workspace with.",
				Flag:        "user",
				Value:       serpent.StringArrayOf(&users),
			}, {
				Name:        "group",
				Description: "A comma separated list of groups to share the workspace with.",
				Flag:        "group",
				Value:       serpent.StringArrayOf(&groups),
			},
		},
		Middleware: serpent.Chain(
			serpent.RequireNArgs(1),
		),
		Handler: func(inv *serpent.Invocation) error {
			if len(users) == 0 && len(groups) == 0 {
				return xerrors.New("at least one user or group must be provided")
			}
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			workspace, err := client.ResolveWorkspace(inv.Context(), inv.Args[0])
			if err != nil {
				return xerrors.Errorf("could not fetch the workspace %s: %w", inv.Args[0], err)
			}

			userRoleStrings := make([][2]string, len(users))
			for index, user := range users {
				if !optimus-ide-collabsdk.UsernameValidRegex.MatchString(user) {
					return xerrors.Errorf("invalid username")
				}

				userRoleStrings[index] = [2]string{user, ""}
			}

			groupRoleStrings := make([][2]string, len(groups))
			for index, group := range groups {
				if !optimus-ide-collabsdk.UsernameValidRegex.MatchString(group) {
					return xerrors.Errorf("invalid group name")
				}

				groupRoleStrings[index] = [2]string{group, ""}
			}

			userRoles, groupRoles, err := fetchUsersAndGroups(inv.Context(), fetchUsersAndGroupsParams{
				Client:      client,
				OrgID:       workspace.OrganizationID,
				OrgName:     workspace.OrganizationName,
				Users:       userRoleStrings,
				Groups:      groupRoleStrings,
				DefaultRole: optimus-ide-collabsdk.WorkspaceRoleDeleted,
			})
			if err != nil {
				return err
			}

			err = client.UpdateWorkspaceACL(inv.Context(), workspace.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
				UserRoles:  userRoles,
				GroupRoles: groupRoles,
			})
			if err != nil {
				return err
			}

			acl, err := client.WorkspaceACL(inv.Context(), workspace.ID)
			if err != nil {
				return xerrors.Errorf("could not fetch current workspace ACL after sharing %w", err)
			}

			out, err := workspaceACLToTable(inv.Context(), &acl)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(inv.Stdout, out)
			return err
		},
	}

	return cmd
}

func stringToWorkspaceRole(role string) (optimus-ide-collabsdk.WorkspaceRole, error) {
	switch role {
	case string(optimus-ide-collabsdk.WorkspaceRoleUse):
		return optimus-ide-collabsdk.WorkspaceRoleUse, nil
	case string(optimus-ide-collabsdk.WorkspaceRoleAdmin):
		return optimus-ide-collabsdk.WorkspaceRoleAdmin, nil
	case string(optimus-ide-collabsdk.WorkspaceRoleDeleted):
		return optimus-ide-collabsdk.WorkspaceRoleDeleted, nil
	default:
		return "", xerrors.Errorf("invalid role %q: expected %q, %q, or \"%q\"",
			role, optimus-ide-collabsdk.WorkspaceRoleAdmin, optimus-ide-collabsdk.WorkspaceRoleUse, optimus-ide-collabsdk.WorkspaceRoleDeleted)
	}
}

func workspaceACLToTable(ctx context.Context, acl *optimus-ide-collabsdk.WorkspaceACL) (string, error) {
	type workspaceShareRow struct {
		User  string                 `table:"user"`
		Group string                 `table:"group,default_sort"`
		Role  optimus-ide-collabsdk.WorkspaceRole `table:"role"`
	}

	formatter := cliui.NewOutputFormatter(
		cliui.TableFormat(
			[]workspaceShareRow{}, []string{"User", "Group", "Role"}),
		cliui.JSONFormat())

	outputRows := make([]workspaceShareRow, 0)
	for _, user := range acl.Users {
		if user.Role == optimus-ide-collabsdk.WorkspaceRoleDeleted {
			continue
		}

		outputRows = append(outputRows, workspaceShareRow{
			User:  user.Username,
			Group: defaultGroupDisplay,
			Role:  user.Role,
		})
	}
	for _, group := range acl.Groups {
		if group.Role == optimus-ide-collabsdk.WorkspaceRoleDeleted {
			continue
		}

		// The ACL endpoint intentionally omits the group's member roster to
		// avoid leaking member PII, so we display one row per group rather
		// than one row per member.
		outputRows = append(outputRows, workspaceShareRow{
			User:  defaultGroupDisplay,
			Group: group.Name,
			Role:  group.Role,
		})
	}
	out, err := formatter.Format(ctx, outputRows)
	if err != nil {
		return "", err
	}

	return out, nil
}

type fetchUsersAndGroupsParams struct {
	Client      *optimus-ide-collabsdk.Client
	OrgID       uuid.UUID
	OrgName     string
	Users       [][2]string
	Groups      [][2]string
	DefaultRole optimus-ide-collabsdk.WorkspaceRole
}

func fetchUsersAndGroups(ctx context.Context, params fetchUsersAndGroupsParams) (userRoles map[string]optimus-ide-collabsdk.WorkspaceRole, groupRoles map[string]optimus-ide-collabsdk.WorkspaceRole, err error) {
	var (
		client      = params.Client
		orgID       = params.OrgID
		orgName     = params.OrgName
		users       = params.Users
		groups      = params.Groups
		defaultRole = params.DefaultRole
	)

	userRoles = make(map[string]optimus-ide-collabsdk.WorkspaceRole, len(users))
	if len(users) > 0 {
		orgMembers, err := client.OrganizationMembers(ctx, orgID)
		if err != nil {
			return nil, nil, err
		}

		for _, user := range users {
			username := user[0]
			role := user[1]
			if role == "" {
				role = string(defaultRole)
			}

			userID := ""
			for _, member := range orgMembers {
				if member.Username == username {
					userID = member.UserID.String()
					break
				}
			}
			if userID == "" {
				return nil, nil, xerrors.Errorf("could not find user %s in the organization %s", username, orgName)
			}

			workspaceRole, err := stringToWorkspaceRole(role)
			if err != nil {
				return nil, nil, err
			}

			userRoles[userID] = workspaceRole
		}
	}

	groupRoles = make(map[string]optimus-ide-collabsdk.WorkspaceRole)
	if len(groups) > 0 {
		orgGroups, err := client.Groups(ctx, optimus-ide-collabsdk.GroupArguments{
			Organization: orgID.String(),
		})
		if err != nil {
			return nil, nil, err
		}

		for _, group := range groups {
			groupName := group[0]
			role := group[1]
			if role == "" {
				role = string(defaultRole)
			}

			var orgGroup *optimus-ide-collabsdk.Group
			for _, og := range orgGroups {
				if og.Name == groupName {
					orgGroup = &og
					break
				}
			}

			if orgGroup == nil {
				return nil, nil, xerrors.Errorf("could not find group named %s belonging to the organization %s", groupName, orgName)
			}

			workspaceRole, err := stringToWorkspaceRole(role)
			if err != nil {
				return nil, nil, err
			}

			groupRoles[orgGroup.ID.String()] = workspaceRole
		}
	}

	return userRoles, groupRoles, nil
}
