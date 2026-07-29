package cli_test

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// TestStart also tests restart since the tests are virtually identical.
func TestStart(t *testing.T) {
	t.Parallel()

	t.Run("RequireActiveVersion", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitMedium)
		ownerClient, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAccessControl:              1,
					optimus-ide-collabsdk.FeatureTemplateRBAC:               1,
					optimus-ide-collabsdk.FeatureAdvancedTemplateScheduling: 1,
				},
			},
		})
		templateAdminClient, templateAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.RoleTemplateAdmin())

		// Create an initial version.
		oldVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, templateAdminClient, owner.OrganizationID, nil)
		// Create a template that mandates the promoted version.
		// This should be enforced for everyone except template admins.
		template := optimus-ide-collabdtest.CreateTemplate(t, templateAdminClient, owner.OrganizationID, oldVersion.ID)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, templateAdminClient, oldVersion.ID)
		require.Equal(t, oldVersion.ID, template.ActiveVersionID)
		template = optimus-ide-collabdtest.UpdateTemplateMeta(t, templateAdminClient, template.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
			RequireActiveVersion: ptr.Ref(true),
		})
		require.True(t, template.RequireActiveVersion)

		// Create a new version that we will promote.
		activeVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, templateAdminClient, owner.OrganizationID, nil, func(ctvr *optimus-ide-collabsdk.CreateTemplateVersionRequest) {
			ctvr.TemplateID = template.ID
		})
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, templateAdminClient, activeVersion.ID)
		err := templateAdminClient.UpdateActiveTemplateVersion(ctx, template.ID, optimus-ide-collabsdk.UpdateActiveTemplateVersion{
			ID: activeVersion.ID,
		})
		require.NoError(t, err)

		templateACLAdminClient, templateACLAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)
		templateGroupACLAdminClient, templateGroupACLAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)
		memberClient, member := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)

		// Create a group so we can also test group template admin ownership.
		// Add the user who gains template admin via group membership.
		group := optimus-ide-collabdtest.CreateGroup(t, ownerClient, owner.OrganizationID, "test", templateGroupACLAdmin)

		// Update the template for both users and groups.
		err = ownerClient.UpdateTemplateACL(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateACL{
			UserPerms: map[string]optimus-ide-collabsdk.TemplateRole{
				templateACLAdmin.ID.String(): optimus-ide-collabsdk.TemplateRoleAdmin,
			},
			GroupPerms: map[string]optimus-ide-collabsdk.TemplateRole{
				group.ID.String(): optimus-ide-collabsdk.TemplateRoleAdmin,
			},
		})
		require.NoError(t, err)

		type testcase struct {
			Name            string
			Client          *optimus-ide-collabsdk.Client
			WorkspaceOwner  uuid.UUID
			ExpectedVersion uuid.UUID
		}

		// All users should be updated to the active version when
		// require_active_version is set, matching web UI behavior.
		cases := []testcase{
			{
				Name:            "OwnerUpdates",
				Client:          ownerClient,
				WorkspaceOwner:  owner.UserID,
				ExpectedVersion: activeVersion.ID,
			},
			{
				Name:            "TemplateAdminUpdates",
				Client:          templateAdminClient,
				WorkspaceOwner:  templateAdmin.ID,
				ExpectedVersion: activeVersion.ID,
			},
			{
				Name:            "TemplateACLAdminUpdates",
				Client:          templateACLAdminClient,
				WorkspaceOwner:  templateACLAdmin.ID,
				ExpectedVersion: activeVersion.ID,
			},
			{
				Name:            "TemplateGroupACLAdminUpdates",
				Client:          templateGroupACLAdminClient,
				WorkspaceOwner:  templateGroupACLAdmin.ID,
				ExpectedVersion: activeVersion.ID,
			},
			{
				Name:            "MemberUpdates",
				Client:          memberClient,
				WorkspaceOwner:  member.ID,
				ExpectedVersion: activeVersion.ID,
			},
		}

		for _, cmd := range []string{"start", "restart"} {
			t.Run(cmd, func(t *testing.T) {
				t.Parallel()
				for _, c := range cases {
					t.Run(c.Name, func(t *testing.T) {
						t.Parallel()

						// Instantiate a new context for each subtest since
						// they can potentially be lengthy.
						ctx := testutil.Context(t, testutil.WaitMedium)
						// Create the workspace using the admin since we want
						// to force the old version.
						ws, err := ownerClient.CreateWorkspace(ctx, owner.OrganizationID, c.WorkspaceOwner.String(), optimus-ide-collabsdk.CreateWorkspaceRequest{
							TemplateVersionID: oldVersion.ID,
							Name:              optimus-ide-collabdtest.RandomUsername(t),
							AutomaticUpdates:  optimus-ide-collabsdk.AutomaticUpdatesNever,
						})
						require.NoError(t, err)
						optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, c.Client, ws.LatestBuild.ID)

						initialTemplateVersion := ws.LatestBuild.TemplateVersionID

						if cmd == "start" {
							// Stop the workspace so that we can start it.
							optimus-ide-collabdtest.MustTransitionWorkspace(t, c.Client, ws.ID, optimus-ide-collabsdk.WorkspaceTransitionStart, optimus-ide-collabsdk.WorkspaceTransitionStop)
						}
						// Start the workspace. Every test permutation should
						// pass.
						var buf bytes.Buffer
						inv, conf := newCLI(t, cmd, ws.Name, "-y")
						inv.Stdout = &buf
						clitest.SetupConfig(t, c.Client, conf)
						err = inv.Run()
						require.NoError(t, err)

						ws = optimus-ide-collabdtest.MustWorkspace(t, c.Client, ws.ID)
						require.Equal(t, c.ExpectedVersion, ws.LatestBuild.TemplateVersionID)
						// The CLI should proactively use the active version
						// without hitting the 403→retry path.
						if initialTemplateVersion != ws.LatestBuild.TemplateVersionID {
							require.NotContains(t, buf.String(), "Unable to start the workspace with the template version from the last build")
							require.NotContains(t, buf.String(), "Unable to restart the workspace with the template version from the last build")
						}
					})
				}
			})
		}
	})

	t.Run("StartActivatesDormant", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitMedium)
		ownerClient, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAdvancedTemplateScheduling: 1,
				},
			},
		})

		version := optimus-ide-collabdtest.CreateTemplateVersion(t, ownerClient, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, ownerClient, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, ownerClient, owner.OrganizationID, version.ID)

		memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)
		workspace := optimus-ide-collabdtest.CreateWorkspace(t, memberClient, template.ID)
		_ = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, memberClient, workspace.LatestBuild.ID)
		_ = optimus-ide-collabdtest.MustTransitionWorkspace(t, memberClient, workspace.ID, optimus-ide-collabsdk.WorkspaceTransitionStart, optimus-ide-collabsdk.WorkspaceTransitionStop)
		err := memberClient.UpdateWorkspaceDormancy(ctx, workspace.ID, optimus-ide-collabsdk.UpdateWorkspaceDormancy{
			Dormant: true,
		})
		require.NoError(t, err)

		inv, root := newCLI(t, "start", workspace.Name)
		clitest.SetupConfig(t, memberClient, root)

		var buf bytes.Buffer
		inv.Stdout = &buf

		err = inv.Run()
		require.NoError(t, err)
		require.Contains(t, buf.String(), "Activating dormant workspace...")

		workspace = optimus-ide-collabdtest.MustWorkspace(t, memberClient, workspace.ID)
		require.Equal(t, optimus-ide-collabsdk.WorkspaceTransitionStart, workspace.LatestBuild.Transition)
	})
}
