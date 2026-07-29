package optimus-ide-collabd_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestWorkspaceBuildAfterTemplateAccessRevokedFails(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)
	ownerClient, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		Options: &optimus-ide-collabdtest.Options{
			IncludeProvisionerDaemon: true,
		},
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
			},
		},
	})
	memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)

	// Given: a member owns a workspace created from a template they can access.
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, ownerClient, owner.OrganizationID, nil)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, ownerClient, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, ownerClient, owner.OrganizationID, version.ID)

	workspace, err := memberClient.CreateUserWorkspace(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateWorkspaceRequest{
		TemplateID:        template.ID,
		Name:              testutil.GetRandomNameHyphenated(t),
		AutomaticUpdates:  optimus-ide-collabsdk.AutomaticUpdatesNever,
		AutostartSchedule: ptr.Ref(""),
	})
	require.NoError(t, err)
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, memberClient, workspace.LatestBuild.ID)

	stopBuild, err := memberClient.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
	})
	require.NoError(t, err)
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, memberClient, stopBuild.ID)

	// When: the template ACL is changed so the member can no longer access it.
	err = ownerClient.UpdateTemplateACL(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateACL{
		GroupPerms: map[string]optimus-ide-collabsdk.TemplateRole{
			owner.OrganizationID.String(): optimus-ide-collabsdk.TemplateRoleDeleted,
		},
	})
	require.NoError(t, err)

	_, err = memberClient.Template(ctx, template.ID)
	require.Error(t, err)
	var tplErr *optimus-ide-collabsdk.Error
	require.ErrorAs(t, err, &tplErr)
	require.Equal(t, http.StatusNotFound, tplErr.StatusCode())
	// Then: starting the existing workspace fails because the build path still
	// requires template access.
	_, err = memberClient.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
	})
	require.Error(t, err)
	apiErr, ok := optimus-ide-collabsdk.AsError(err)
	require.True(t, ok)
	require.Equal(t, http.StatusInternalServerError, apiErr.StatusCode())
	require.Contains(t, apiErr.Message, "failed to fetch template")
}

func TestWorkspaceBuild(t *testing.T) {
	t.Parallel()

	// Only use this context for setup. Use a separate context for subtests!
	setupCtx := testutil.Context(t, testutil.WaitMedium)
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

	// For this test we create two templates:
	// tplA will be used to test creation of new workspaces.
	// tplB will be used to test builds on existing workspaces.
	// This is done to enable parallelization of the sub-tests without them interfering with each other.
	// Both templates mandate the promoted version.
	// This should be enforced for everyone except template admins.
	tplAv1 := optimus-ide-collabdtest.CreateTemplateVersion(t, ownerClient, owner.OrganizationID, nil)
	tplA := optimus-ide-collabdtest.CreateTemplate(t, ownerClient, owner.OrganizationID, tplAv1.ID)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, ownerClient, tplAv1.ID)
	require.Equal(t, tplAv1.ID, tplA.ActiveVersionID)
	tplA = optimus-ide-collabdtest.UpdateTemplateMeta(t, ownerClient, tplA.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
		RequireActiveVersion: ptr.Ref(true),
	})
	require.True(t, tplA.RequireActiveVersion)
	tplAv2 := optimus-ide-collabdtest.CreateTemplateVersion(t, ownerClient, owner.OrganizationID, nil, func(ctvr *optimus-ide-collabsdk.CreateTemplateVersionRequest) {
		ctvr.TemplateID = tplA.ID
	})
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, ownerClient, tplAv2.ID)
	optimus-ide-collabdtest.UpdateActiveTemplateVersion(t, ownerClient, tplA.ID, tplAv2.ID)

	tplBv1 := optimus-ide-collabdtest.CreateTemplateVersion(t, ownerClient, owner.OrganizationID, nil)
	tplB := optimus-ide-collabdtest.CreateTemplate(t, ownerClient, owner.OrganizationID, tplBv1.ID)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, ownerClient, tplBv1.ID)
	require.Equal(t, tplBv1.ID, tplB.ActiveVersionID)
	tplB = optimus-ide-collabdtest.UpdateTemplateMeta(t, ownerClient, tplB.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
		RequireActiveVersion: ptr.Ref(true),
	})
	require.True(t, tplB.RequireActiveVersion)

	templateAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.RoleTemplateAdmin())
	templateACLAdminClient, templateACLAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)
	templateGroupACLAdminClient, templateGroupACLAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)
	memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)

	// Create a group so we can also test group template admin ownership.
	// Add the user who gains template admin via group membership.
	group := optimus-ide-collabdtest.CreateGroup(t, ownerClient, owner.OrganizationID, "test", templateGroupACLAdmin)

	// Update the template for both users and groups.
	//nolint:gocritic // test setup
	for _, tpl := range []optimus-ide-collabsdk.Template{tplA, tplB} {
		err := ownerClient.UpdateTemplateACL(setupCtx, tpl.ID, optimus-ide-collabsdk.UpdateTemplateACL{
			UserPerms: map[string]optimus-ide-collabsdk.TemplateRole{
				templateACLAdmin.ID.String(): optimus-ide-collabsdk.TemplateRoleAdmin,
			},
			GroupPerms: map[string]optimus-ide-collabsdk.TemplateRole{
				group.ID.String(): optimus-ide-collabsdk.TemplateRoleAdmin,
			},
		})
		require.NoError(t, err, "updating template ACL for template %q", tpl.ID)
	}

	type testcase struct {
		Name               string
		Client             *optimus-ide-collabsdk.Client
		ExpectedStatusCode int
	}

	cases := []testcase{
		{
			Name:               "OwnerOK",
			Client:             ownerClient,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "TemplateAdminOK",
			Client:             templateAdminClient,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "TemplateACLAdminOK",
			Client:             templateACLAdminClient,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "TemplateGroupACLAdminOK",
			Client:             templateGroupACLAdminClient,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:               "MemberFailsToCreate",
			Client:             memberClient,
			ExpectedStatusCode: http.StatusForbidden,
		},
	}

	// Create pre-existing workspaces for each of the test cases.
	var extantWorkspaces []optimus-ide-collabsdk.Workspace
	for _, c := range cases {
		extantWs, err := c.Client.CreateUserWorkspace(setupCtx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateWorkspaceRequest{
			TemplateVersionID: tplB.ActiveVersionID,
			Name:              testutil.GetRandomNameHyphenated(t),
			AutomaticUpdates:  optimus-ide-collabsdk.AutomaticUpdatesNever,
		})
		require.NoError(t, err, "setup workspace for case %q", c.Name)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, c.Client, extantWs.LatestBuild.ID)
		extantWorkspaces = append(extantWorkspaces, extantWs)
	}

	// Create a new version of template B and promote it to be the active version.
	tplBv2 := optimus-ide-collabdtest.CreateTemplateVersion(t, ownerClient, owner.OrganizationID, nil, func(ctvr *optimus-ide-collabsdk.CreateTemplateVersionRequest) {
		ctvr.TemplateID = tplB.ID
	})
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, ownerClient, tplBv2.ID)
	optimus-ide-collabdtest.UpdateActiveTemplateVersion(t, ownerClient, tplB.ID, tplBv2.ID)

	t.Run("NewWorkspace", func(t *testing.T) {
		t.Parallel()

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				t.Parallel()
				ctx := testutil.Context(t, testutil.WaitMedium)
				ws, err := c.Client.CreateUserWorkspace(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateWorkspaceRequest{
					TemplateVersionID: tplAv1.ID,
					Name:              testutil.GetRandomNameHyphenated(t),
					AutomaticUpdates:  optimus-ide-collabsdk.AutomaticUpdatesNever,
				})
				if c.ExpectedStatusCode == http.StatusOK {
					require.NoError(t, err)
					require.Equal(t, tplAv1.ID, ws.LatestBuild.TemplateVersionID, "workspace did not use expected version for case %q", c.Name)
				} else {
					require.Error(t, err)
					cerr, ok := optimus-ide-collabsdk.AsError(err)
					require.True(t, ok)
					require.Equal(t, c.ExpectedStatusCode, cerr.StatusCode())
				}
			})
		}
	})

	t.Run("ExistingWorkspace", func(t *testing.T) {
		t.Parallel()

		for idx, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				t.Parallel()
				ctx := testutil.Context(t, testutil.WaitMedium)
				// Stopping the workspace must always succeed.
				wb, err := c.Client.CreateWorkspaceBuild(ctx, extantWorkspaces[idx].ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
					Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
				})
				require.NoError(t, err, "stopping workspace for case %q", c.Name)
				optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, c.Client, wb.ID)

				// Attempt to start the workspace with the given version.
				wb, err = c.Client.CreateWorkspaceBuild(ctx, extantWorkspaces[idx].ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
					Transition:        optimus-ide-collabsdk.WorkspaceTransitionStart,
					TemplateVersionID: tplBv1.ID,
				})
				if c.ExpectedStatusCode == http.StatusOK {
					require.NoError(t, err, "starting workspace for case %q", c.Name)
					optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, c.Client, wb.ID)
					require.Equal(t, tplBv1.ID, wb.TemplateVersionID, "workspace did not use expected version for case %q", c.Name)
				} else {
					require.Error(t, err, "starting workspace for case %q", c.Name)
					cerr, ok := optimus-ide-collabsdk.AsError(err)
					require.True(t, ok)
					require.Equal(t, c.ExpectedStatusCode, cerr.StatusCode())
				}
			})
		}
	})
}
