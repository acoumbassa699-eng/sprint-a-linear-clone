package acl_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/acl"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestOK(t *testing.T) {
	t.Parallel()

	db, _ := dbtestutil.NewDB(t)
	o := dbgen.Organization(t, db, database.Organization{})
	g := dbgen.Group(t, db, database.Group{OrganizationID: o.ID})
	u := dbgen.User(t, db, database.User{})
	ctx := testutil.Context(t, testutil.WaitShort)

	update := optimus-ide-collabsdk.UpdateWorkspaceACL{
		UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
			u.ID.String(): optimus-ide-collabsdk.WorkspaceRoleAdmin,
			// An unknown ID is allowed if and only if the specified role is either
			// optimus-ide-collabsdk.WorkspaceRoleDeleted or optimus-ide-collabsdk.TemplateRoleDeleted.
			uuid.NewString(): optimus-ide-collabsdk.WorkspaceRoleDeleted,
		},
		GroupRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
			g.ID.String(): optimus-ide-collabsdk.WorkspaceRoleAdmin,
			// An unknown ID is allowed if and only if the specified role is either
			// optimus-ide-collabsdk.WorkspaceRoleDeleted or optimus-ide-collabsdk.TemplateRoleDeleted.
			uuid.NewString(): optimus-ide-collabsdk.WorkspaceRoleDeleted,
		},
	}
	errors := acl.Validate(ctx, db, optimus-ide-collabd.WorkspaceACLUpdateValidator(update))
	require.Empty(t, errors)
}

func TestDeniesUnknownIDs(t *testing.T) {
	t.Parallel()

	db, _ := dbtestutil.NewDB(t)
	ctx := testutil.Context(t, testutil.WaitShort)

	update := optimus-ide-collabsdk.UpdateWorkspaceACL{
		UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
			uuid.NewString(): optimus-ide-collabsdk.WorkspaceRoleAdmin,
		},
		GroupRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
			uuid.NewString(): optimus-ide-collabsdk.WorkspaceRoleAdmin,
		},
	}
	errors := acl.Validate(ctx, db, optimus-ide-collabd.WorkspaceACLUpdateValidator(update))
	require.Len(t, errors, 2)
	require.Equal(t, errors[0].Field, "group_roles")
	require.ErrorContains(t, errors[0], "does not exist")
	require.Equal(t, errors[1].Field, "user_roles")
	require.ErrorContains(t, errors[1], "does not exist")
}

func TestDeniesUnknownRolesAndInvalidIDs(t *testing.T) {
	t.Parallel()

	db, _ := dbtestutil.NewDB(t)
	ctx := testutil.Context(t, testutil.WaitShort)

	update := optimus-ide-collabsdk.UpdateWorkspaceACL{
		UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
			"Quifrey": "level 5",
		},
		GroupRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
			"apprentices": "level 2",
		},
	}
	errors := acl.Validate(ctx, db, optimus-ide-collabd.WorkspaceACLUpdateValidator(update))
	require.Len(t, errors, 4)
	require.Equal(t, errors[0].Field, "group_roles")
	require.ErrorContains(t, errors[0], "role \"level 2\" is not a valid workspace role")
	require.Equal(t, errors[1].Field, "group_roles")
	require.ErrorContains(t, errors[1], "not a valid UUID")
	require.Equal(t, errors[2].Field, "user_roles")
	require.ErrorContains(t, errors[2], "role \"level 5\" is not a valid workspace role")
	require.Equal(t, errors[3].Field, "user_roles")
	require.ErrorContains(t, errors[3], "not a valid UUID")
}
