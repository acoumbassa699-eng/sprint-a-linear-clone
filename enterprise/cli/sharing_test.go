package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestSharingShare(t *testing.T) {
	t.Parallel()

	t.Run("ShareWithGroups_Simple", func(t *testing.T) {
		t.Parallel()

		var (
			client, db, orgOwner = optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
					},
				},
			})
			workspaceOwnerClient, workspaceOwner = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID, rbac.ScopedRoleOrgAuditor(orgOwner.OrganizationID))
			workspace                            = dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
				OwnerID:        workspaceOwner.ID,
				OrganizationID: orgOwner.OrganizationID,
			}).Do().Workspace
			_, orgMember = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID)
		)

		ctx := testutil.Context(t, testutil.WaitMedium)

		group, err := createGroupWithMembers(ctx, client, orgOwner.OrganizationID, "new-group", []uuid.UUID{orgMember.ID})
		require.NoError(t, err)

		inv, root := clitest.New(t, "sharing", "share", workspace.Name, "--group", group.Name)
		clitest.SetupConfig(t, workspaceOwnerClient, root)

		out := new(bytes.Buffer)
		inv.Stdout = out
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)

		acl, err := workspaceOwnerClient.WorkspaceACL(inv.Context(), workspace.ID)
		require.NoError(t, err)
		assert.Len(t, acl.Groups, 1)
		assert.Equal(t, acl.Groups[0].Group.ID, group.ID)
		assert.Equal(t, acl.Groups[0].Role, optimus-ide-collabsdk.WorkspaceRoleUse)

		found := false
		for _, line := range strings.Split(out.String(), "\n") {
			found = strings.Contains(line, group.Name) && strings.Contains(line, string(optimus-ide-collabsdk.WorkspaceRoleUse))
			if found {
				break
			}
		}
		assert.True(t, found, "Expected to find group name %s and role %s in output: %s", group.Name, optimus-ide-collabsdk.WorkspaceRoleUse, out.String())
	})

	t.Run("ShareWithGroups_Multiple", func(t *testing.T) {
		t.Parallel()

		var (
			client, db, orgOwner = optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
					},
				},
			})

			workspaceOwnerClient, workspaceOwner = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID, rbac.ScopedRoleOrgAuditor(orgOwner.OrganizationID))
			workspace                            = dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
				OwnerID:        workspaceOwner.ID,
				OrganizationID: orgOwner.OrganizationID,
			}).Do().Workspace

			_, wibbleMember = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID)
			_, wobbleMember = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID)
		)

		ctx := testutil.Context(t, testutil.WaitMedium)

		wibbleGroup, err := createGroupWithMembers(ctx, client, orgOwner.OrganizationID, "wibble", []uuid.UUID{wibbleMember.ID})
		require.NoError(t, err)

		wobbleGroup, err := createGroupWithMembers(ctx, client, orgOwner.OrganizationID, "wobble", []uuid.UUID{wobbleMember.ID})
		require.NoError(t, err)

		inv, root := clitest.New(t, "sharing", "share", workspace.Name,
			fmt.Sprintf("--group=%s,%s", wibbleGroup.Name, wobbleGroup.Name))
		clitest.SetupConfig(t, workspaceOwnerClient, root)

		out := new(bytes.Buffer)
		inv.Stdout = out
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)

		acl, err := workspaceOwnerClient.WorkspaceACL(inv.Context(), workspace.ID)
		require.NoError(t, err)
		assert.Len(t, acl.Groups, 2)

		type workspaceGroup []optimus-ide-collabsdk.WorkspaceGroup
		assert.NotEqual(t, -1, slices.IndexFunc(workspaceGroup(acl.Groups), func(g optimus-ide-collabsdk.WorkspaceGroup) bool {
			return g.Group.ID == wibbleGroup.ID
		}))
		assert.NotEqual(t, -1, slices.IndexFunc(workspaceGroup(acl.Groups), func(g optimus-ide-collabsdk.WorkspaceGroup) bool {
			return g.Group.ID == wobbleGroup.ID
		}))

		t.Run("ShareWithGroups_Role", func(t *testing.T) {
			t.Parallel()

			var (
				client, db, orgOwner = optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
					LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
						Features: license.Features{
							optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
						},
					},
				})
				workspaceOwnerClient, workspaceOwner = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID, rbac.ScopedRoleOrgAuditor(orgOwner.OrganizationID))
				workspace                            = dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
					OwnerID:        workspaceOwner.ID,
					OrganizationID: orgOwner.OrganizationID,
				}).Do().Workspace
				_, orgMember = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID)
			)

			ctx := testutil.Context(t, testutil.WaitMedium)

			group, err := createGroupWithMembers(ctx, client, orgOwner.OrganizationID, "new-group", []uuid.UUID{orgMember.ID})
			require.NoError(t, err)

			inv, root := clitest.New(t, "sharing", "share", workspace.Name, "--group", fmt.Sprintf("%s:admin", group.Name))
			clitest.SetupConfig(t, workspaceOwnerClient, root)

			out := new(bytes.Buffer)
			inv.Stdout = out
			err = inv.WithContext(ctx).Run()
			require.NoError(t, err)

			acl, err := workspaceOwnerClient.WorkspaceACL(inv.Context(), workspace.ID)
			require.NoError(t, err)
			assert.Len(t, acl.Groups, 1)
			assert.Equal(t, acl.Groups[0].Group.ID, group.ID)
			assert.Equal(t, acl.Groups[0].Role, optimus-ide-collabsdk.WorkspaceRoleAdmin)

			found := false
			for _, line := range strings.Split(out.String(), "\n") {
				found = strings.Contains(line, group.Name) && strings.Contains(line, string(optimus-ide-collabsdk.WorkspaceRoleAdmin))
				if found {
					break
				}
			}
			assert.True(t, found, "Expected to find group name %s and role %s in output: %s", group.Name, optimus-ide-collabsdk.WorkspaceRoleAdmin, out.String())
		})
	})
}

func TestSharingStatus(t *testing.T) {
	t.Parallel()

	t.Run("ListSharedUsers", func(t *testing.T) {
		t.Parallel()

		var (
			client, db, orgOwner = optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
					},
				},
			})
			workspaceOwnerClient, workspaceOwner = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID, rbac.ScopedRoleOrgAuditor(orgOwner.OrganizationID))
			workspace                            = dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
				OwnerID:        workspaceOwner.ID,
				OrganizationID: orgOwner.OrganizationID,
			}).Do().Workspace
			_, orgMember = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID)
			ctx          = testutil.Context(t, testutil.WaitMedium)
		)

		group, err := createGroupWithMembers(ctx, client, orgOwner.OrganizationID, "new-group", []uuid.UUID{orgMember.ID})
		require.NoError(t, err)

		err = client.UpdateWorkspaceACL(ctx, workspace.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			GroupRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				group.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		require.NoError(t, err)

		inv, root := clitest.New(t, "sharing", "status", workspace.Name)
		clitest.SetupConfig(t, workspaceOwnerClient, root)

		out := new(bytes.Buffer)
		inv.Stdout = out
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)

		// The ACL endpoint omits group member rosters to avoid leaking member
		// PII, so the output lists the group itself rather than its members.
		found := false
		for _, line := range strings.Split(out.String(), "\n") {
			if strings.Contains(line, group.Name) && strings.Contains(line, string(optimus-ide-collabsdk.WorkspaceRoleUse)) {
				found = true
				break
			}
		}
		assert.True(t, found, "expected to find group %s with role %s in the output: %s", group.Name, optimus-ide-collabsdk.WorkspaceRoleUse, out.String())
		assert.NotContains(t, out.String(), orgMember.Username, "group member roster must not be exposed in sharing status output")
	})
}

func TestSharingRemove(t *testing.T) {
	t.Parallel()

	t.Run("RemoveSharedGroup_Single", func(t *testing.T) {
		t.Parallel()

		var (
			client, db, orgOwner = optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
					},
				},
			})
			workspaceOwnerClient, workspaceOwner = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID, rbac.ScopedRoleOrgAuditor(orgOwner.OrganizationID))
			workspace                            = dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
				OwnerID:        workspaceOwner.ID,
				OrganizationID: orgOwner.OrganizationID,
			}).Do().Workspace
			_, groupUser1 = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID)
			_, groupUser2 = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID)
		)

		ctx := testutil.Context(t, testutil.WaitMedium)

		group1, err := createGroupWithMembers(ctx, client, orgOwner.OrganizationID, "group-1", []uuid.UUID{groupUser1.ID, groupUser2.ID})
		require.NoError(t, err)

		group2, err := createGroupWithMembers(ctx, client, orgOwner.OrganizationID, "group-2", []uuid.UUID{groupUser1.ID, groupUser2.ID})
		require.NoError(t, err)

		// Share the workspace with a user to later remove
		err = client.UpdateWorkspaceACL(ctx, workspace.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			GroupRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				group1.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
				group2.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		require.NoError(t, err)

		inv, root := clitest.New(t,
			"sharing",
			"remove",
			workspace.Name,
			"--group", group1.Name,
		)
		clitest.SetupConfig(t, workspaceOwnerClient, root)

		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)

		acl, err := workspaceOwnerClient.WorkspaceACL(inv.Context(), workspace.ID)
		require.NoError(t, err)

		removedGroup1 := true
		removedGroup2 := true
		for _, group := range acl.Groups {
			if group.ID == group1.ID {
				removedGroup1 = false
				continue
			}

			if group.ID == group2.ID {
				removedGroup2 = false
				continue
			}
		}
		assert.True(t, removedGroup1)
		assert.False(t, removedGroup2)
	})

	t.Run("RemoveSharedGroup_Multiple", func(t *testing.T) {
		t.Parallel()

		var (
			client, db, orgOwner = optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
					},
				},
			})
			workspaceOwnerClient, workspaceOwner = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID, rbac.ScopedRoleOrgAuditor(orgOwner.OrganizationID))
			workspace                            = dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
				OwnerID:        workspaceOwner.ID,
				OrganizationID: orgOwner.OrganizationID,
			}).Do().Workspace
			_, groupUser1 = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID)
			_, groupUser2 = optimus-ide-collabdtest.CreateAnotherUser(t, client, orgOwner.OrganizationID)
		)

		ctx := testutil.Context(t, testutil.WaitMedium)

		group1, err := createGroupWithMembers(ctx, client, orgOwner.OrganizationID, "group-1", []uuid.UUID{groupUser1.ID, groupUser2.ID})
		require.NoError(t, err)

		group2, err := createGroupWithMembers(ctx, client, orgOwner.OrganizationID, "group-2", []uuid.UUID{groupUser1.ID, groupUser2.ID})
		require.NoError(t, err)

		// Share the workspace with a user to later remove
		err = client.UpdateWorkspaceACL(ctx, workspace.ID, optimus-ide-collabsdk.UpdateWorkspaceACL{
			GroupRoles: map[string]optimus-ide-collabsdk.WorkspaceRole{
				group1.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
				group2.ID.String(): optimus-ide-collabsdk.WorkspaceRoleUse,
			},
		})
		require.NoError(t, err)

		inv, root := clitest.New(t,
			"sharing",
			"remove",
			workspace.Name,
			fmt.Sprintf("--group=%s,%s", group1.Name, group2.Name),
		)
		clitest.SetupConfig(t, workspaceOwnerClient, root)

		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)

		acl, err := workspaceOwnerClient.WorkspaceACL(inv.Context(), workspace.ID)
		require.NoError(t, err)

		removedGroup1 := true
		removedGroup2 := true
		for _, group := range acl.Groups {
			if group.ID == group1.ID {
				removedGroup1 = false
				continue
			}

			if group.ID == group2.ID {
				removedGroup2 = false
				continue
			}
		}
		assert.True(t, removedGroup1)
		assert.True(t, removedGroup2)
	})
}

func createGroupWithMembers(ctx context.Context, client *optimus-ide-collabsdk.Client, orgID uuid.UUID, name string, memberIDs []uuid.UUID) (optimus-ide-collabsdk.Group, error) {
	group, err := client.CreateGroup(ctx, orgID, optimus-ide-collabsdk.CreateGroupRequest{
		Name:        name,
		DisplayName: name,
	})
	if err != nil {
		return optimus-ide-collabsdk.Group{}, err
	}

	ids := make([]string, len(memberIDs))
	for i, id := range memberIDs {
		ids[i] = id.String()
	}

	return client.PatchGroup(ctx, group.ID, optimus-ide-collabsdk.PatchGroupRequest{
		AddUsers: ids,
	})
}
