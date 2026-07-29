package optimus-ide-collabd_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestWorkspaceSharingSettings(t *testing.T) {
	t.Parallel()

	t.Run("DisabledDefaultsFalse", func(t *testing.T) {
		t.Parallel()

		dv := optimus-ide-collabdtest.DeploymentValues(t)

		client, first := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues: dv,
			},
		})

		ctx := testutil.Context(t, testutil.WaitMedium)

		// Use a regular user to make sure the setting is exposed to them.
		memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID)
		settings, err := memberClient.WorkspaceSharingSettings(ctx, first.OrganizationID.String())
		require.NoError(t, err)
		// Check the deprecated boolean field.
		require.False(t, settings.SharingDisabled)
		require.Equal(t, optimus-ide-collabsdk.ShareableWorkspaceOwnersEveryone, settings.ShareableWorkspaceOwners)
	})

	t.Run("DisabledTogglePersists", func(t *testing.T) {
		t.Parallel()

		dv := optimus-ide-collabdtest.DeploymentValues(t)

		client, first := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues: dv,
			},
		})

		ctx := testutil.Context(t, testutil.WaitMedium)

		orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID, rbac.ScopedRoleOrgAdmin(first.OrganizationID))

		// Disable sharing via the deprecated boolean field.
		settings, err := orgAdminClient.PatchWorkspaceSharingSettings(ctx, first.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			SharingDisabled: true,
		})
		require.NoError(t, err)
		require.True(t, settings.SharingDisabled)
		require.Equal(t, optimus-ide-collabsdk.ShareableWorkspaceOwnersNone, settings.ShareableWorkspaceOwners)

		settings, err = orgAdminClient.WorkspaceSharingSettings(ctx, first.OrganizationID.String())
		require.NoError(t, err)
		require.True(t, settings.SharingDisabled)
		require.Equal(t, optimus-ide-collabsdk.ShareableWorkspaceOwnersNone, settings.ShareableWorkspaceOwners)

		// Switch to service_accounts mode via the new field.
		settings, err = orgAdminClient.PatchWorkspaceSharingSettings(ctx, first.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			ShareableWorkspaceOwners: optimus-ide-collabsdk.ShareableWorkspaceOwnersServiceAccounts,
		})
		require.NoError(t, err)
		require.False(t, settings.SharingDisabled)
		require.Equal(t, optimus-ide-collabsdk.ShareableWorkspaceOwnersServiceAccounts, settings.ShareableWorkspaceOwners)

		settings, err = orgAdminClient.WorkspaceSharingSettings(ctx, first.OrganizationID.String())
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ShareableWorkspaceOwnersServiceAccounts, settings.ShareableWorkspaceOwners)

		// Re-enable full sharing.
		settings, err = orgAdminClient.PatchWorkspaceSharingSettings(ctx, first.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			ShareableWorkspaceOwners: optimus-ide-collabsdk.ShareableWorkspaceOwnersEveryone,
		})
		require.NoError(t, err)
		require.False(t, settings.SharingDisabled)
		require.Equal(t, optimus-ide-collabsdk.ShareableWorkspaceOwnersEveryone, settings.ShareableWorkspaceOwners)

		settings, err = orgAdminClient.WorkspaceSharingSettings(ctx, first.OrganizationID.String())
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ShareableWorkspaceOwnersEveryone, settings.ShareableWorkspaceOwners)
	})

	t.Run("InvalidValueRejected", func(t *testing.T) {
		t.Parallel()

		client, first := optimus-ide-collabdenttest.New(t, nil)

		ctx := testutil.Context(t, testutil.WaitMedium)

		orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID, rbac.ScopedRoleOrgAdmin(first.OrganizationID))
		_, err := orgAdminClient.PatchWorkspaceSharingSettings(ctx, first.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			ShareableWorkspaceOwners: "invalid",
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})

	t.Run("UpdateAuthz", func(t *testing.T) {
		t.Parallel()

		dv := optimus-ide-collabdtest.DeploymentValues(t)

		client, first := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues: dv,
			},
		})

		ctx := testutil.Context(t, testutil.WaitMedium)

		memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID)
		_, err := memberClient.PatchWorkspaceSharingSettings(ctx, first.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			SharingDisabled: true,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusForbidden, apiErr.StatusCode())
	})

	t.Run("AuditLog", func(t *testing.T) {
		t.Parallel()

		auditor := audit.NewMock()
		dv := optimus-ide-collabdtest.DeploymentValues(t)

		client, first := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			AuditLogging: true,
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues: dv,
				Auditor:          auditor,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAuditLog: 1,
				},
			},
		})

		ctx := testutil.Context(t, testutil.WaitMedium)

		orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID, rbac.ScopedRoleOrgAdmin(first.OrganizationID))
		auditor.ResetLogs()
		_, err := orgAdminClient.PatchWorkspaceSharingSettings(ctx, first.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			SharingDisabled: true,
		})
		require.NoError(t, err)

		require.Len(t, auditor.AuditLogs(), 1)
		alog := auditor.AuditLogs()[0]
		require.Equal(t, database.AuditActionWrite, alog.Action)
		require.Equal(t, database.ResourceTypeOrganization, alog.ResourceType)
		require.Equal(t, first.OrganizationID, alog.ResourceID)
	})
}

func TestWorkspaceSharingDisabled(t *testing.T) {
	t.Parallel()

	t.Run("ACLEndpointsForbidden", func(t *testing.T) {
		t.Parallel()

		dv := optimus-ide-collabdtest.DeploymentValues(t)

		client, db, owner := optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues: dv,
			},
		})

		workspaceOwnerClient, workspaceOwner := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		ws := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        workspaceOwner.ID,
			OrganizationID: owner.OrganizationID,
		}).Do().Workspace

		ctx := testutil.Context(t, testutil.WaitMedium)

		orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.ScopedRoleOrgAdmin(owner.OrganizationID))
		_, err := orgAdminClient.PatchWorkspaceSharingSettings(ctx, owner.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			ShareableWorkspaceOwners: optimus-ide-collabsdk.ShareableWorkspaceOwnersNone,
		})
		require.NoError(t, err)

		// Reading the ACL list remains allowed even when workspace sharing is
		// disabled, but mutating it is forbidden.
		_, err = workspaceOwnerClient.WorkspaceACL(ctx, ws.ID)
		require.NoError(t, err)

		// We don't allow mutating the ACL.
		assertSharingDisabled := func(t *testing.T, err error) {
			t.Helper()

			var apiErr *optimus-ide-collabsdk.Error
			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, http.StatusForbidden, apiErr.StatusCode())
			require.Equal(t, "Workspace sharing is disabled for this organization.", apiErr.Message)
		}

		// Despite the site-wide workspace.share permission for the owner,
		// the endpoint should return an authz error.
		err = client.UpdateWorkspaceACL(ctx, ws.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				uuid.NewString(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		assertSharingDisabled(t, err)

		err = workspaceOwnerClient.DeleteWorkspaceACL(ctx, ws.ID)
		assertSharingDisabled(t, err)
	})

	t.Run("ACLEndpointsForbiddenServiceAccountsMode", func(t *testing.T) {
		t.Parallel()

		client, db, owner := optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureServiceAccounts: 1,
				},
			},
		})

		regularClient, regularUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		regularWS := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        regularUser.ID,
			OrganizationID: owner.OrganizationID,
		}).Do().Workspace

		// Create an SA with a workspace.
		saClient, saUser := optimus-ide-collabdtest.CreateAnotherUserMutators(t, client, owner.OrganizationID, nil, func(r *optimus-ide-collabsdk.CreateUserRequestWithOrgs) {
			r.ServiceAccount = true
		})
		saWS := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        saUser.ID,
			OrganizationID: owner.OrganizationID,
		}).Do().Workspace

		ctx := testutil.Context(t, testutil.WaitMedium)

		orgAdminClient, orgAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.ScopedRoleOrgAdmin(owner.OrganizationID))
		_, err := orgAdminClient.PatchWorkspaceSharingSettings(ctx, owner.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			ShareableWorkspaceOwners: optimus-ide-collabsdk.ShareableWorkspaceOwnersServiceAccounts,
		})
		require.NoError(t, err)

		// Regular member cannot share their own workspace.
		err = regularClient.UpdateWorkspaceACL(ctx, regularWS.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				orgAdmin.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusForbidden, apiErr.StatusCode())

		// SA can share their own workspace.
		err = saClient.UpdateWorkspaceACL(ctx, saWS.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				regularUser.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		require.NoError(t, err)
	})

	// Future-proofing: if custom roles with member-scoped
	// workspace:share are ever allowed, the member-level negation
	// from the organization-member system role must block sharing in
	// service_accounts mode even with such custom role.
	t.Run("MemberCannotBypassWithCustomRole", func(t *testing.T) {
		t.Parallel()

		rawDB, pubsub, sqlDB := dbtestutil.NewDBWithSQLDB(t)
		client, _, _, owner := optimus-ide-collabdenttest.NewWithAPI(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				Database: rawDB,
				Pubsub:   pubsub,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureCustomRoles:  1,
					optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
				},
			},
		})

		ctx := testutil.Context(t, testutil.WaitMedium)

		// Create an empty custom role via the API, then add
		// member-scoped workspace:share via raw SQL (the API and
		// dbauthz both reject member permissions on custom roles).
		//nolint:gocritic // owner context required for role creation
		customRole, err := client.CreateOrganizationRole(ctx, optimus-ide-collabsdk.Role{
			Name:           "workspace-share-granter",
			OrganizationID: owner.OrganizationID.String(),
		})
		require.NoError(t, err)

		_, err = sqlDB.ExecContext(ctx,
			`UPDATE custom_roles SET member_permissions = $1 WHERE name = $2 AND organization_id = $3`,
			database.CustomRolePermissions{{
				ResourceType: rbac.ResourceWorkspace.Type,
				Action:       policy.ActionShare,
			}},
			customRole.Name,
			owner.OrganizationID,
		)
		require.NoError(t, err)

		// Create a member and assign the custom role.
		memberClient, memberUser := optimus-ide-collabdtest.CreateAnotherUserMutators(
			t, client, owner.OrganizationID,
			[]rbac.RoleIdentifier{{
				Name:           customRole.Name,
				OrganizationID: owner.OrganizationID,
			}},
		)
		memberWS := dbfake.WorkspaceBuild(t, rawDB, database.WorkspaceTable{
			OwnerID:        memberUser.ID,
			OrganizationID: owner.OrganizationID,
		}).Do().Workspace

		_, sharedUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		// Switch to service_accounts mode.
		orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.ScopedRoleOrgAdmin(owner.OrganizationID))
		_, err = orgAdminClient.PatchWorkspaceSharingSettings(ctx, owner.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			ShareableWorkspaceOwners: optimus-ide-collabsdk.ShareableWorkspaceOwnersServiceAccounts,
		})
		require.NoError(t, err)

		// Despite the custom role granting workspace:share at the
		// member level, the negation from organization-member should
		// block it.
		err = memberClient.UpdateWorkspaceACL(ctx, memberWS.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				sharedUser.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusForbidden, apiErr.StatusCode())
	})

	t.Run("ACLsPurged", func(t *testing.T) {
		t.Parallel()

		dv := optimus-ide-collabdtest.DeploymentValues(t)

		client, db, owner := optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues: dv,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
				},
			},
		})

		workspaceOwnerClient, workspaceOwner := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		_, sharedUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		// Create a group to test group ACL purging.
		group := optimus-ide-collabdtest.CreateGroup(t, client, owner.OrganizationID, "test-group")

		ws := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        workspaceOwner.ID,
			OrganizationID: owner.OrganizationID,
		}).Do().Workspace

		ctx := testutil.Context(t, testutil.WaitMedium)

		// Set both user and group ACLs.
		err := workspaceOwnerClient.UpdateWorkspaceACL(ctx, ws.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				sharedUser.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
			GroupRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				group.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		require.NoError(t, err)

		acl, err := workspaceOwnerClient.WorkspaceACL(ctx, ws.ID)
		require.NoError(t, err)
		require.Len(t, acl.Users, 1)
		require.Equal(t, sharedUser.ID, acl.Users[0].ID)
		require.Equal(t, optimus-ide-collabsdk.WorkspaceRoleUse, acl.Users[0].Role)
		require.Len(t, acl.Groups, 1)
		require.Equal(t, group.ID, acl.Groups[0].ID)
		require.Equal(t, optimus-ide-collabsdk.WorkspaceRoleUse, acl.Groups[0].Role)

		orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.ScopedRoleOrgAdmin(owner.OrganizationID))
		_, err = orgAdminClient.PatchWorkspaceSharingSettings(ctx, owner.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			ShareableWorkspaceOwners: optimus-ide-collabsdk.ShareableWorkspaceOwnersNone,
		})
		require.NoError(t, err)

		_, err = orgAdminClient.PatchWorkspaceSharingSettings(ctx, owner.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			ShareableWorkspaceOwners: optimus-ide-collabsdk.ShareableWorkspaceOwnersEveryone,
		})
		require.NoError(t, err)

		// Verify both user and group ACLs are purged.
		acl, err = workspaceOwnerClient.WorkspaceACL(ctx, ws.ID)
		require.NoError(t, err)
		require.Empty(t, acl.Users)
		require.Empty(t, acl.Groups)

		// Verify ACLs can be set again after re-enabling sharing.
		err = workspaceOwnerClient.UpdateWorkspaceACL(ctx, ws.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				sharedUser.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		require.NoError(t, err)
		acl, err = workspaceOwnerClient.WorkspaceACL(ctx, ws.ID)
		require.NoError(t, err)
		require.Len(t, acl.Users, 1)
		require.Equal(t, sharedUser.ID, acl.Users[0].ID)
	})

	t.Run("ACLsPurgedExceptServiceAccounts", func(t *testing.T) {
		t.Parallel()

		dv := optimus-ide-collabdtest.DeploymentValues(t)

		client, db, owner := optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues: dv,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureTemplateRBAC:    1,
					optimus-ide-collabsdk.FeatureServiceAccounts: 1,
				},
			},
		})

		// Regular user with a workspace.
		workspaceOwnerClient, workspaceOwner := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		_, sharedUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		regularWS := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        workspaceOwner.ID,
			OrganizationID: owner.OrganizationID,
		}).Do().Workspace

		// Service account with a workspace.
		_, saUser := optimus-ide-collabdtest.CreateAnotherUserMutators(t, client, owner.OrganizationID, nil, func(r *optimus-ide-collabsdk.CreateUserRequestWithOrgs) {
			r.ServiceAccount = true
		})
		saWS := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        saUser.ID,
			OrganizationID: owner.OrganizationID,
		}).Do().Workspace

		ctx := testutil.Context(t, testutil.WaitMedium)

		// Share regular user's workspace with sharedUser.
		err := workspaceOwnerClient.UpdateWorkspaceACL(ctx, regularWS.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				sharedUser.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		require.NoError(t, err)

		// Use the owner client (site admin) to share the SA workspace,
		// since the SA can't authenticate via the API.
		err = client.UpdateWorkspaceACL(ctx, saWS.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			UserRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				sharedUser.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		require.NoError(t, err)

		// Switch to service_accounts mode.
		orgAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.ScopedRoleOrgAdmin(owner.OrganizationID))
		_, err = orgAdminClient.PatchWorkspaceSharingSettings(ctx, owner.OrganizationID.String(), optimus-ide-collabsdk.UpdateWorkspaceSharingSettingsRequest{
			ShareableWorkspaceOwners: optimus-ide-collabsdk.ShareableWorkspaceOwnersServiceAccounts,
		})
		require.NoError(t, err)

		// Regular user workspace ACLs should be purged.
		acl, err := workspaceOwnerClient.WorkspaceACL(ctx, regularWS.ID)
		require.NoError(t, err)
		require.Empty(t, acl.Users)

		// Service account workspace ACLs should be preserved.
		acl, err = client.WorkspaceACL(ctx, saWS.ID)
		require.NoError(t, err)
		require.Len(t, acl.Users, 1)
		require.Equal(t, sharedUser.ID, acl.Users[0].ID)
	})
}
