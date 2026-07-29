package cli_test

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/pubsub"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
	"github.com/optimus-ide-collab/quartz"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	t.Run("WithParameter", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)
		workspace := optimus-ide-collabdtest.CreateWorkspace(t, member, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)
		inv, root := clitest.New(t, "delete", workspace.Name, "-y")
		clitest.SetupConfig(t, member, root)
		doneChan := make(chan struct{})
		stdout := expecter.NewAttachedToInvocation(t, inv)
		go func() {
			defer close(doneChan)
			err := inv.Run()
			// When running with the race detector on, we sometimes get an EOF.
			if err != nil {
				assert.ErrorIs(t, err, io.EOF)
			}
		}()
		stdout.ExpectMatch(ctx, "has been deleted")
		<-doneChan
	})

	t.Run("Orphan", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, templateAdmin, owner.OrganizationID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, templateAdmin, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, templateAdmin, owner.OrganizationID, version.ID)
		workspace := optimus-ide-collabdtest.CreateWorkspace(t, templateAdmin, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, templateAdmin, workspace.LatestBuild.ID)

		ctx := testutil.Context(t, testutil.WaitShort)
		inv, root := clitest.New(t, "delete", workspace.Name, "-y", "--orphan")
		clitest.SetupConfig(t, templateAdmin, root)

		doneChan := make(chan struct{})
		stdout := expecter.NewAttachedToInvocation(t, inv)
		go func() {
			defer close(doneChan)
			err := inv.WithContext(ctx).Run()
			// When running with the race detector on, we sometimes get an EOF.
			if err != nil {
				assert.ErrorIs(t, err, io.EOF)
			}
		}()
		stdout.ExpectMatch(ctx, "has been deleted")
		testutil.TryReceive(ctx, t, doneChan)

		_, err := client.Workspace(ctx, workspace.ID)
		require.Error(t, err)
		cerr := optimus-ide-collabdtest.SDKError(t, err)
		require.Equal(t, http.StatusGone, cerr.StatusCode())
	})

	// Super orphaned, as the workspace doesn't even have a user.
	// This is not a scenario we should ever get into, as we do not allow users
	// to be deleted if they have workspaces. However issue #7872 shows that
	// it is possible to get into this state. An admin should be able to still
	// force a delete action on the workspace.
	t.Run("OrphanDeletedUser", func(t *testing.T) {
		t.Parallel()
		client, _, api := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		deleteMeClient, deleteMeUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)
		workspace := optimus-ide-collabdtest.CreateWorkspace(t, deleteMeClient, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, deleteMeClient, workspace.LatestBuild.ID)

		// The API checks if the user has any workspaces, so we cannot delete a user
		// this way.
		ctx := testutil.Context(t, testutil.WaitShort)
		err := api.Database.UpdateUserDeletedByID(dbauthz.AsSystemRestricted(ctx), deleteMeUser.ID)
		require.NoError(t, err)

		inv, root := clitest.New(t, "delete", fmt.Sprintf("%s/%s", deleteMeUser.ID, workspace.Name), "-y", "--orphan")

		//nolint:gocritic // Deleting orphaned workspaces requires an admin.
		clitest.SetupConfig(t, client, root)
		doneChan := make(chan struct{})
		stdout := expecter.NewAttachedToInvocation(t, inv)
		go func() {
			defer close(doneChan)
			err := inv.Run()
			// When running with the race detector on, we sometimes get an EOF.
			if err != nil {
				assert.ErrorIs(t, err, io.EOF)
			}
		}()
		stdout.ExpectMatch(ctx, "has been deleted")
		<-doneChan
	})

	t.Run("DifferentUser", func(t *testing.T) {
		t.Parallel()
		adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		adminUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)
		orgID := adminUser.OrganizationID
		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, orgID)
		user, err := client.User(context.Background(), optimus-ide-collabsdk.Me)
		require.NoError(t, err)

		version := optimus-ide-collabdtest.CreateTemplateVersion(t, adminClient, orgID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, adminClient, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, adminClient, orgID, version.ID)
		workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

		ctx := testutil.Context(t, testutil.WaitMedium)
		inv, root := clitest.New(t, "delete", user.Username+"/"+workspace.Name, "-y")
		//nolint:gocritic // This requires an admin.
		clitest.SetupConfig(t, adminClient, root)
		doneChan := make(chan struct{})
		stdout := expecter.NewAttachedToInvocation(t, inv)
		go func() {
			defer close(doneChan)
			err := inv.Run()
			// When running with the race detector on, we sometimes get an EOF.
			if err != nil {
				assert.ErrorIs(t, err, io.EOF)
			}
		}()

		stdout.ExpectMatch(ctx, "has been deleted")
		<-doneChan

		workspace, err = client.Workspace(context.Background(), workspace.ID)
		require.ErrorContains(t, err, "was deleted")
	})

	t.Run("InvalidWorkspaceIdentifier", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		inv, root := clitest.New(t, "delete", "a/b/c", "-y")
		clitest.SetupConfig(t, client, root)
		doneChan := make(chan struct{})
		go func() {
			defer close(doneChan)
			err := inv.Run()
			assert.ErrorContains(t, err, "invalid workspace identifier: \"a/b/c\"")
		}()
		<-doneChan
	})

	t.Run("WarnNoProvisioners", func(t *testing.T) {
		t.Parallel()

		store, ps, db := dbtestutil.NewDBWithSQLDB(t)
		client, closeDaemon := optimus-ide-collabdtest.NewWithProvisionerCloser(t, &optimus-ide-collabdtest.Options{
			Database:                 store,
			Pubsub:                   ps,
			IncludeProvisionerDaemon: true,
		})

		// Given: a user, template, and workspace
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, user.OrganizationID, rbac.RoleTemplateAdmin())
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, templateAdmin, user.OrganizationID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, templateAdmin, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, templateAdmin, user.OrganizationID, version.ID)
		workspace := optimus-ide-collabdtest.CreateWorkspace(t, templateAdmin, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, templateAdmin, workspace.LatestBuild.ID)

		// When: all provisioner daemons disappear
		require.NoError(t, closeDaemon.Close())
		_, err := db.Exec("DELETE FROM provisioner_daemons;")
		require.NoError(t, err)

		// Then: the workspace deletion should warn about no provisioners
		inv, root := clitest.New(t, "delete", workspace.Name, "-y")
		stdout := expecter.NewAttachedToInvocation(t, inv)
		clitest.SetupConfig(t, templateAdmin, root)
		doneChan := make(chan struct{})
		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()
		go func() {
			defer close(doneChan)
			_ = inv.WithContext(ctx).Run()
		}()
		stdout.ExpectMatch(ctx, "there are no provisioners that accept the required tags")
		cancel()
		<-doneChan
	})

	t.Run("Prebuilt workspace delete permissions", func(t *testing.T) {
		t.Parallel()

		// Setup
		db, pb := dbtestutil.NewDB(t, dbtestutil.WithDumpOnFailure())
		client, _ := optimus-ide-collabdtest.NewWithProvisionerCloser(t, &optimus-ide-collabdtest.Options{
			Database:                 db,
			Pubsub:                   pb,
			IncludeProvisionerDaemon: true,
		})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		orgID := owner.OrganizationID

		// Given a template version with a preset and a template
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, orgID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		preset := setupTestDBPreset(t, db, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, orgID, version.ID)

		cases := []struct {
			name                          string
			client                        *optimus-ide-collabsdk.Client
			expectedPrebuiltDeleteErrMsg  string
			expectedWorkspaceDeleteErrMsg string
		}{
			// Users with the OrgAdmin role should be able to delete both normal and prebuilt workspaces
			{
				name: "OrgAdmin",
				client: func() *optimus-ide-collabsdk.Client {
					client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, orgID, rbac.ScopedRoleOrgAdmin(orgID))
					return client
				}(),
			},
			// Users with the TemplateAdmin role should be able to delete prebuilt workspaces, but not normal workspaces
			{
				name: "TemplateAdmin",
				client: func() *optimus-ide-collabsdk.Client {
					client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, orgID, rbac.RoleTemplateAdmin())
					return client
				}(),
				expectedWorkspaceDeleteErrMsg: "unexpected status code 403: You do not have permission to delete this workspace.",
			},
			// Users with the OrgTemplateAdmin role should be able to delete prebuilt workspaces, but not normal workspaces
			{
				name: "OrgTemplateAdmin",
				client: func() *optimus-ide-collabsdk.Client {
					client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, orgID, rbac.ScopedRoleOrgTemplateAdmin(orgID))
					return client
				}(),
				expectedWorkspaceDeleteErrMsg: "unexpected status code 403: You do not have permission to delete this workspace.",
			},
			// Users with the Member role should not be able to delete prebuilt or normal workspaces
			{
				name: "Member",
				client: func() *optimus-ide-collabsdk.Client {
					client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, orgID, rbac.RoleMember())
					return client
				}(),
				expectedPrebuiltDeleteErrMsg:  "unexpected status code 404: Resource not found or you do not have access to this resource",
				expectedWorkspaceDeleteErrMsg: "unexpected status code 404: Resource not found or you do not have access to this resource",
			},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				clock := quartz.NewMock(t)
				ctx := testutil.Context(t, testutil.WaitSuperLong)

				// Create one prebuilt workspace (owned by system user) and one normal workspace (owned by a user)
				// Each workspace is persisted in the DB along with associated workspace jobs and builds.
				dbPrebuiltWorkspace := setupTestDBWorkspace(t, clock, db, pb, orgID, database.PrebuildsSystemUserID, template.ID, version.ID, preset.ID)
				userWorkspaceOwner, err := client.User(context.Background(), "testUser")
				require.NoError(t, err)
				dbUserWorkspace := setupTestDBWorkspace(t, clock, db, pb, orgID, userWorkspaceOwner.ID, template.ID, version.ID, preset.ID)

				assertWorkspaceDelete := func(
					runClient *optimus-ide-collabsdk.Client,
					workspace database.Workspace,
					workspaceOwner string,
					expectedErr string,
				) {
					t.Helper()

					// Attempt to delete the workspace as the test client
					inv, root := clitest.New(t, "delete", workspaceOwner+"/"+workspace.Name, "-y")
					clitest.SetupConfig(t, runClient, root)
					doneChan := make(chan struct{})
					stdout := expecter.NewAttachedToInvocation(t, inv)
					var runErr error
					go func() {
						defer close(doneChan)
						runErr = inv.Run()
					}()

					// Validate the result based on the expected error message
					if expectedErr != "" {
						<-doneChan
						require.Error(t, runErr)
						require.Contains(t, runErr.Error(), expectedErr)
					} else {
						stdout.ExpectMatch(ctx, "has been deleted")
						<-doneChan

						// When running with the race detector on, we sometimes get an EOF.
						if runErr != nil {
							assert.ErrorIs(t, runErr, io.EOF)
						}

						// Verify that the workspace is now marked as deleted
						_, err := client.Workspace(context.Background(), workspace.ID)
						require.ErrorContains(t, err, "was deleted")
					}
				}

				// Ensure at least one prebuilt workspace is reported as running in the database
				testutil.Eventually(ctx, t, func(ctx context.Context) (done bool) {
					running, err := db.GetRunningPrebuiltWorkspaces(ctx)
					if !assert.NoError(t, err) || !assert.GreaterOrEqual(t, len(running), 1) {
						return false
					}
					return true
				}, testutil.IntervalMedium, "running prebuilt workspaces timeout")

				runningWorkspaces, err := db.GetRunningPrebuiltWorkspaces(ctx)
				require.NoError(t, err)
				require.GreaterOrEqual(t, len(runningWorkspaces), 1)

				// Get the full prebuilt workspace object from the DB
				prebuiltWorkspace, err := db.GetWorkspaceByID(ctx, dbPrebuiltWorkspace.ID)
				require.NoError(t, err)

				// Assert the prebuilt workspace deletion
				assertWorkspaceDelete(tc.client, prebuiltWorkspace, "prebuilds", tc.expectedPrebuiltDeleteErrMsg)

				// Get the full user workspace object from the DB
				userWorkspace, err := db.GetWorkspaceByID(ctx, dbUserWorkspace.ID)
				require.NoError(t, err)

				// Assert the user workspace deletion
				assertWorkspaceDelete(tc.client, userWorkspace, userWorkspaceOwner.Username, tc.expectedWorkspaceDeleteErrMsg)
			})
		}
	})
}

func setupTestDBPreset(
	t *testing.T,
	db database.Store,
	templateVersionID uuid.UUID,
) database.TemplateVersionPreset {
	t.Helper()

	preset := dbgen.Preset(t, db, database.InsertPresetParams{
		TemplateVersionID: templateVersionID,
		Name:              "preset-test",
		DesiredInstances: sql.NullInt32{
			Valid: true,
			Int32: 1,
		},
	})
	dbgen.PresetParameter(t, db, database.InsertPresetParametersParams{
		TemplateVersionPresetID: preset.ID,
		Names:                   []string{"test"},
		Values:                  []string{"test"},
	})

	return preset
}

func setupTestDBWorkspace(
	t *testing.T,
	clock quartz.Clock,
	db database.Store,
	ps pubsub.Pubsub,
	orgID uuid.UUID,
	ownerID uuid.UUID,
	templateID uuid.UUID,
	templateVersionID uuid.UUID,
	presetID uuid.UUID,
) database.WorkspaceTable {
	t.Helper()

	workspace := dbgen.Workspace(t, db, database.WorkspaceTable{
		TemplateID:     templateID,
		OrganizationID: orgID,
		OwnerID:        ownerID,
		Deleted:        false,
		CreatedAt:      time.Now().Add(-time.Hour * 2),
	})
	job := dbgen.ProvisionerJob(t, db, ps, database.ProvisionerJob{
		InitiatorID:    ownerID,
		CreatedAt:      time.Now().Add(-time.Hour * 2),
		StartedAt:      sql.NullTime{Time: clock.Now().Add(-time.Hour * 2), Valid: true},
		CompletedAt:    sql.NullTime{Time: clock.Now().Add(-time.Hour), Valid: true},
		OrganizationID: orgID,
	})
	workspaceBuild := dbgen.WorkspaceBuild(t, db, database.WorkspaceBuild{
		WorkspaceID:             workspace.ID,
		InitiatorID:             ownerID,
		TemplateVersionID:       templateVersionID,
		JobID:                   job.ID,
		TemplateVersionPresetID: uuid.NullUUID{UUID: presetID, Valid: true},
		Transition:              database.WorkspaceTransitionStart,
		CreatedAt:               clock.Now(),
	})
	dbgen.WorkspaceBuildParameters(t, db, []database.WorkspaceBuildParameter{
		{
			WorkspaceBuildID: workspaceBuild.ID,
			Name:             "test",
			Value:            "test",
		},
	})

	return workspace
}
