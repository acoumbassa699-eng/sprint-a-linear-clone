package optimus-ide-collabd_test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisionersdk/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestPostWorkspaceBuildsOnSuccessRestart(t *testing.T) {
	t.Parallel()

	const paramName = "foo"

	// GIVEN: a running workspace with an existing rich parameter value.
	deploymentValues := optimus-ide-collabdtest.DeploymentValues(t)
	deploymentValues.EnableTerraformDebugMode = true
	db, ps, sqlDB := dbtestutil.NewDBWithSQLDB(t)
	client, _, _ := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		Database:                 db,
		Pubsub:                   ps,
		IncludeProvisionerDaemon: true,
		DeploymentValues:         deploymentValues,
	})
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, first.OrganizationID,
		echoResponsesWithRichParameter(paramName, echoResponseOptions{
			blockStopApply: false,
		}),
	)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, first.OrganizationID, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID, func(request *optimus-ide-collabsdk.CreateWorkspaceRequest) {
		request.RichParameterValues = []optimus-ide-collabsdk.WorkspaceBuildParameter{
			{Name: paramName, Value: "bar"},
		}
	})
	initialBuild := optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, initialBuild.Status)

	// WHEN: a stop build is created with an on_success start build.
	ctx := testutil.Context(t, testutil.WaitLong)
	user, err := client.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)

	stopBuild, err := client.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
		Reason:     optimus-ide-collabsdk.CreateWorkspaceBuildReasonCLI,
		LogLevel:   optimus-ide-collabsdk.ProvisionerLogLevelDebug,
		OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
			Transition:        optimus-ide-collabsdk.WorkspaceTransitionStart,
			TemplateVersionID: template.ActiveVersionID,
			RichParameterValues: []optimus-ide-collabsdk.WorkspaceBuildParameter{
				{Name: paramName, Value: "baz"},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceTransitionStop, stopBuild.Transition)
	require.Equal(t, optimus-ide-collabsdk.BuildReasonCLI, stopBuild.Reason)

	// THEN: the server persists the child start build intent.
	orchestration, err := dbtestutil.GetWorkspaceBuildOrchestrationByParentBuildID(ctx, sqlDB, stopBuild.ID)
	require.NoError(t, err)
	require.Equal(t, "pending", orchestration.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceTransitionStart, optimus-ide-collabsdk.WorkspaceTransition(orchestration.ChildTransition))
	require.True(t, orchestration.ChildTemplateVersionID.Valid)
	require.Equal(t, template.ActiveVersionID, orchestration.ChildTemplateVersionID.UUID)
	require.False(t, orchestration.ChildTemplateVersionPresetID.Valid)
	require.Equal(t, string(optimus-ide-collabsdk.ProvisionerLogLevelDebug), orchestration.ChildLogLevel)
	require.True(t, orchestration.ChildReason.Valid)
	require.Equal(t, optimus-ide-collabsdk.BuildReasonCLI, optimus-ide-collabsdk.BuildReason(orchestration.ChildReason.BuildReason))

	var childRichParameterValues []optimus-ide-collabsdk.WorkspaceBuildParameter
	require.NoError(t, json.Unmarshal(orchestration.ChildRichParameterValues, &childRichParameterValues))
	require.ElementsMatch(t, []optimus-ide-collabsdk.WorkspaceBuildParameter{
		{Name: paramName, Value: "baz"},
	}, childRichParameterValues)

	// THEN: the returned parent stop build completes successfully.
	stopBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, stopBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, stopBuild.Job.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusStopped, stopBuild.Status)

	// THEN: the server creates and completes the child start build.
	var childBuild optimus-ide-collabsdk.WorkspaceBuild
	require.Eventually(t, func() bool {
		childBuild, err = client.WorkspaceBuildByUsernameAndWorkspaceNameAndBuildNumber(
			ctx,
			user.Username,
			workspace.Name,
			strconv.FormatInt(int64(stopBuild.BuildNumber+1), 10),
		)
		return err == nil &&
			childBuild.Transition == optimus-ide-collabsdk.WorkspaceTransitionStart
	}, testutil.WaitMedium, testutil.IntervalFast)

	childBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, childBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, childBuild.Job.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, childBuild.Status)
	require.Equal(t, optimus-ide-collabsdk.BuildReasonCLI, childBuild.Reason)
	require.Equal(t, template.ActiveVersionID, childBuild.TemplateVersionID)

	// THEN: the child build uses the on_success parameter values.
	params, err := client.WorkspaceBuildParameters(ctx, childBuild.ID)
	require.NoError(t, err)
	require.ElementsMatch(t, []optimus-ide-collabsdk.WorkspaceBuildParameter{
		{Name: paramName, Value: "baz"},
	}, params)
}

func TestPostWorkspaceBuildsOnSuccessTemplateVersionPreset(t *testing.T) {
	t.Parallel()

	// GIVEN: a running workspace and a preset on its active template
	// version.
	db, ps, sqlDB := dbtestutil.NewDBWithSQLDB(t)
	client, _, _ := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		Database:                 db,
		Pubsub:                   ps,
		IncludeProvisionerDaemon: true,
	})
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, first.OrganizationID, nil)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, first.OrganizationID, version.ID)
	preset := dbgen.Preset(t, db, database.InsertPresetParams{
		Name:              "on-success-preset",
		TemplateVersionID: version.ID,
	})
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
	initialBuild := optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, initialBuild.Status)

	// WHEN: a stop build is created with an on_success start build
	// that requests the preset.
	ctx := testutil.Context(t, testutil.WaitLong)
	user, err := client.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)

	stopBuild, err := client.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
		OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
			Transition:              optimus-ide-collabsdk.WorkspaceTransitionStart,
			TemplateVersionID:       template.ActiveVersionID,
			TemplateVersionPresetID: preset.ID,
		},
	})
	require.NoError(t, err)

	// THEN: the server persists the child preset intent.
	orchestration, err := dbtestutil.GetWorkspaceBuildOrchestrationByParentBuildID(ctx, sqlDB, stopBuild.ID)
	require.NoError(t, err)
	require.True(t, orchestration.ChildTemplateVersionPresetID.Valid)
	require.Equal(t, preset.ID, orchestration.ChildTemplateVersionPresetID.UUID)

	stopBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, stopBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, stopBuild.Job.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusStopped, stopBuild.Status)

	// THEN: the child start build uses the preset.
	var childBuild optimus-ide-collabsdk.WorkspaceBuild
	require.Eventually(t, func() bool {
		childBuild, err = client.WorkspaceBuildByUsernameAndWorkspaceNameAndBuildNumber(
			ctx,
			user.Username,
			workspace.Name,
			strconv.FormatInt(int64(stopBuild.BuildNumber+1), 10),
		)
		return err == nil &&
			childBuild.Transition == optimus-ide-collabsdk.WorkspaceTransitionStart
	}, testutil.WaitShort, testutil.IntervalFast)

	childBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, childBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, childBuild.Job.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, childBuild.Status)
	require.NotNil(t, childBuild.TemplateVersionPresetID)
	require.Equal(t, preset.ID, *childBuild.TemplateVersionPresetID)
}

func TestPostWorkspaceBuildsOnSuccessUnpinnedChildUsesActiveTemplateVersion(t *testing.T) {
	t.Parallel()

	// GIVEN: a running workspace and a second completed template
	// version that is not active yet.
	db, ps, sqlDB := dbtestutil.NewDBWithSQLDB(t)
	client, provisionerCloser, api := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		Database:                 db,
		Pubsub:                   ps,
		IncludeProvisionerDaemon: true,
	})
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	userClient, user := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID)

	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, first.OrganizationID, nil)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, first.OrganizationID, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, userClient, template.ID)
	initialBuild := optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, userClient, workspace.LatestBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, initialBuild.Status)

	newVersion := optimus-ide-collabdtest.UpdateTemplateVersion(t, client, first.OrganizationID, nil, template.ID)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, newVersion.ID)

	// WHEN: a non-template-admin queues an unpinned on_success child
	// build.

	// Stop the provisioner so the parent build cannot complete before
	// the test updates the active template version.
	require.NoError(t, provisionerCloser.Close())
	ctx := testutil.Context(t, testutil.WaitLong)
	stopBuild, err := userClient.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
		OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
			Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
		},
	})
	require.NoError(t, err)

	// THEN: the child build remains unpinned in the orchestration row.
	orchestration, err := dbtestutil.GetWorkspaceBuildOrchestrationByParentBuildID(ctx, sqlDB, stopBuild.ID)
	require.NoError(t, err)
	require.False(t, orchestration.ChildTemplateVersionID.Valid, "child build should remain unpinned")

	// WHEN: the active version changes before the parent succeeds.
	optimus-ide-collabdtest.UpdateActiveTemplateVersion(t, client, template.ID, newVersion.ID)
	optimus-ide-collabdtest.NewProvisionerDaemon(t, api)

	stopBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, userClient, stopBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, stopBuild.Job.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusStopped, stopBuild.Status)

	// THEN: the child build uses the active version when the
	// orchestrator creates it.
	var childBuild optimus-ide-collabsdk.WorkspaceBuild
	require.Eventually(t, func() bool {
		childBuild, err = userClient.WorkspaceBuildByUsernameAndWorkspaceNameAndBuildNumber(
			ctx,
			user.Username,
			workspace.Name,
			strconv.FormatInt(int64(stopBuild.BuildNumber+1), 10),
		)
		return err == nil &&
			childBuild.Transition == optimus-ide-collabsdk.WorkspaceTransitionStart
	}, testutil.WaitMedium, testutil.IntervalFast)

	childBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, userClient, childBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, childBuild.Job.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, childBuild.Status)
	require.Equal(t, newVersion.ID, childBuild.TemplateVersionID)
}

func TestPostWorkspaceBuildsOnSuccessUnpinnedChildNoParams(t *testing.T) {
	t.Parallel()

	// GIVEN: a running workspace owned by a non-template-admin.
	db, ps, sqlDB := dbtestutil.NewDBWithSQLDB(t)
	client, _, _ := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		Database:                 db,
		Pubsub:                   ps,
		IncludeProvisionerDaemon: true,
	})
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	userClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID)

	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, first.OrganizationID, nil)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, first.OrganizationID, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, userClient, template.ID)
	initialBuild := optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, userClient, workspace.LatestBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, initialBuild.Status)

	// WHEN: the non-template-admin queues a stop build with an unpinned
	// on_success start build that supplies no parameters, reason, or log level.
	ctx := testutil.Context(t, testutil.WaitLong)
	stopBuild, err := userClient.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
		OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
			Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
		},
	})
	// THEN: the request is permitted without template-update privileges,
	// because no durable template version pin is requested.
	require.NoError(t, err)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceTransitionStop, stopBuild.Transition)

	// THEN: the persisted child build intent leaves the optional fields unset.
	orchestration, err := dbtestutil.GetWorkspaceBuildOrchestrationByParentBuildID(ctx, sqlDB, stopBuild.ID)
	require.NoError(t, err)
	require.Equal(t, "pending", orchestration.Status)
	require.Equal(t, workspace.ID, orchestration.WorkspaceID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceTransitionStart, optimus-ide-collabsdk.WorkspaceTransition(orchestration.ChildTransition))
	require.False(t, orchestration.ChildTemplateVersionID.Valid)
	require.False(t, orchestration.ChildTemplateVersionPresetID.Valid)
	require.False(t, orchestration.ChildReason.Valid)
	require.Empty(t, orchestration.ChildLogLevel)

	// THEN: nil parameters are coerced to an empty JSON array, not null, to
	// satisfy the database CHECK constraint.
	require.JSONEq(t, "[]", string(orchestration.ChildRichParameterValues))
	var childRichParameterValues []optimus-ide-collabsdk.WorkspaceBuildParameter
	require.NoError(t, json.Unmarshal(orchestration.ChildRichParameterValues, &childRichParameterValues))
	require.Empty(t, childRichParameterValues)
}

func TestPostWorkspaceBuildsOnSuccessPinnedChildVersionRequiresTemplateUpdate(t *testing.T) {
	t.Parallel()

	// GIVEN: a running workspace owned by a non-template-admin.
	client, _ := optimus-ide-collabdtest.NewWithDatabase(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	userClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID)

	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, first.OrganizationID, nil)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, first.OrganizationID, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, userClient, template.ID)
	initialBuild := optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, userClient, workspace.LatestBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, initialBuild.Status)

	// WHEN: the non-template-admin tries to queue a stop build with a
	// pinned on_success child version.
	ctx := testutil.Context(t, testutil.WaitLong)
	_, err := userClient.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
		OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
			Transition:        optimus-ide-collabsdk.WorkspaceTransitionStart,
			TemplateVersionID: version.ID,
		},
	})
	require.Error(t, err)

	// THEN: the API rejects the durable child version pin and explains the
	// missing template update permission.
	var apiErr *optimus-ide-collabsdk.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusForbidden, apiErr.StatusCode())
	require.Contains(t, apiErr.Response.Detail, "template update permission")

	// THEN: no new workspace build is created.
	builds, err := userClient.WorkspaceBuilds(ctx, optimus-ide-collabsdk.WorkspaceBuildsRequest{WorkspaceID: workspace.ID})
	require.NoError(t, err)
	require.Len(t, builds, 1)
	require.Equal(t, initialBuild.ID, builds[0].ID)
}

func TestPostWorkspaceBuildsOnSuccessValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		request     optimus-ide-collabsdk.CreateWorkspaceBuildRequest
		wantMessage string
	}{
		{
			name: "ParentMustBeStop",
			request: optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
					Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				},
			},
			wantMessage: "OnSuccess is only permitted when stopping a workspace.",
		},
		{
			// The oneof=start struct tag on OnSuccess.Transition rejects this
			// during httpapi.Read, before the explicit check is reached.
			name: "ChildMustBeStart",
			request: optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
				OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
					Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
				},
			},
			wantMessage: "Validation failed.",
		},
		{
			name: "ParentDryRunRejected",
			request: optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
				DryRun:     true,
				OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
					Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				},
			},
			wantMessage: "OnSuccess cannot be set alongside DryRun.",
		},
		{
			name: "ParentOrphanRejected",
			request: optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
				Orphan:     true,
				OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
					Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				},
			},
			wantMessage: "OnSuccess cannot be set alongside Orphan.",
		},
		{
			name: "ParentProvisionerStateRejected",
			request: optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
				Transition:       optimus-ide-collabsdk.WorkspaceTransitionStop,
				ProvisionerState: []byte("state"),
				OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
					Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				},
			},
			wantMessage: "OnSuccess cannot be set alongside ProvisionerState.",
		},
		{
			name: "ChildPresetWithoutVersionRejected",
			request: optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
				OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
					Transition:              optimus-ide-collabsdk.WorkspaceTransitionStart,
					TemplateVersionPresetID: uuid.New(),
				},
			},
			wantMessage: "OnSuccess TemplateVersionPresetID requires TemplateVersionID.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// GIVEN: a running workspace.
			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
			first := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, first.OrganizationID, nil)
			optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, first.OrganizationID, version.ID)
			workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
			optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

			// WHEN: an invalid on_success request is posted.
			_, err := client.CreateWorkspaceBuild(testutil.Context(t, testutil.WaitLong), workspace.ID, tt.request)
			require.Error(t, err)

			// THEN: the API rejects the request before creating a build.
			var apiErr *optimus-ide-collabsdk.Error
			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
			require.Contains(t, apiErr.Message, tt.wantMessage)
		})
	}
}

// Canceling an already-running job resolves the orchestration as
// "failed", not "canceled". This hits the same orchestrator branch as
// TestPostWorkspaceBuildsOnSuccessParentFailed below; despite that
// overlap, the test pins this non-obvious end-to-end behavior.
func TestPostWorkspaceBuildsOnSuccessParentCanceledMidFlight(t *testing.T) {
	t.Parallel()

	// GIVEN: a running workspace whose stop apply will block.
	db, ps, sqlDB := dbtestutil.NewDBWithSQLDB(t)
	client, _, _ := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		Database:                 db,
		Pubsub:                   ps,
		IncludeProvisionerDaemon: true,
	})
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, first.OrganizationID,
		echoResponsesWithRichParameter("foo", echoResponseOptions{
			blockStopApply: true,
		}),
	)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, first.OrganizationID, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
	initialBuild := optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, initialBuild.Status)

	// WHEN: a stop build is created with an on_success start build.
	ctx := testutil.Context(t, testutil.WaitLong)
	user, err := client.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)

	stopBuild, err := client.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
		OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
			Transition:        optimus-ide-collabsdk.WorkspaceTransitionStart,
			TemplateVersionID: template.ActiveVersionID,
		},
	})
	require.NoError(t, err)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceTransitionStop, stopBuild.Transition)
	require.Equal(t, optimus-ide-collabsdk.BuildReasonInitiator, stopBuild.Reason)

	// WHEN: the parent stop build starts running and is canceled.
	require.Eventually(t, func() bool {
		var err error
		stopBuild, err = client.WorkspaceBuild(ctx, stopBuild.ID)
		return err == nil &&
			stopBuild.Job.Status == optimus-ide-collabsdk.ProvisionerJobRunning
	}, testutil.WaitShort, testutil.IntervalFast)

	require.NoError(t, client.CancelWorkspaceBuild(ctx, stopBuild.ID, optimus-ide-collabsdk.CancelWorkspaceBuildParams{}))
	require.Eventually(t, func() bool {
		var err error
		stopBuild, err = client.WorkspaceBuild(ctx, stopBuild.ID)
		if err != nil {
			return false
		}
		return stopBuild.Job.Status == optimus-ide-collabsdk.ProvisionerJobFailed &&
			stopBuild.Job.Error == "canceled"
	}, testutil.WaitShort, testutil.IntervalFast)

	// THEN: the server resolves the orchestration without creating the
	// child start build.
	require.Eventually(t, func() bool {
		orchestration, err := dbtestutil.GetWorkspaceBuildOrchestrationByParentBuildID(ctx, sqlDB, stopBuild.ID)
		return err == nil &&
			orchestration.Status == "failed" &&
			!orchestration.ChildBuildID.Valid &&
			orchestration.Error.Valid &&
			orchestration.Error.String == "parent workspace build failed: canceled"
	}, testutil.WaitShort, testutil.IntervalFast)

	_, err = client.WorkspaceBuildByUsernameAndWorkspaceNameAndBuildNumber(
		ctx,
		user.Username,
		workspace.Name,
		strconv.FormatInt(int64(stopBuild.BuildNumber+1), 10),
	)
	var apiErr *optimus-ide-collabsdk.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
}

func TestPostWorkspaceBuildsOnSuccessParentFailed(t *testing.T) {
	t.Parallel()

	// GIVEN: a running workspace whose stop apply will fail.
	db, ps, sqlDB := dbtestutil.NewDBWithSQLDB(t)
	client, _, _ := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		Database:                 db,
		Pubsub:                   ps,
		IncludeProvisionerDaemon: true,
	})
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, first.OrganizationID,
		echoResponsesWithRichParameter("foo", echoResponseOptions{
			failStopApply: true,
		}),
	)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, first.OrganizationID, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
	initialBuild := optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, initialBuild.Status)

	// WHEN: a stop build is created with an on_success start build.
	ctx := testutil.Context(t, testutil.WaitLong)
	user, err := client.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)

	stopBuild, err := client.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
		OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
			Transition:        optimus-ide-collabsdk.WorkspaceTransitionStart,
			TemplateVersionID: template.ActiveVersionID,
		},
	})
	require.NoError(t, err)

	// WHEN: the parent stop build fails.
	stopBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, stopBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobFailed, stopBuild.Job.Status)

	// THEN: the server resolves the orchestration without creating the
	// child start build.
	require.Eventually(t, func() bool {
		orchestration, err := dbtestutil.GetWorkspaceBuildOrchestrationByParentBuildID(ctx, sqlDB, stopBuild.ID)
		return err == nil &&
			orchestration.Status == "failed" &&
			!orchestration.ChildBuildID.Valid &&
			orchestration.Error.Valid &&
			orchestration.Error.String == "parent workspace build failed: failed!"
	}, testutil.WaitShort, testutil.IntervalFast)

	_, err = client.WorkspaceBuildByUsernameAndWorkspaceNameAndBuildNumber(
		ctx,
		user.Username,
		workspace.Name,
		strconv.FormatInt(int64(stopBuild.BuildNumber+1), 10),
	)
	var apiErr *optimus-ide-collabsdk.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
}

func TestPostWorkspaceBuildsOnSuccessNonRetryableChildBuildFailure(t *testing.T) {
	t.Parallel()

	// GIVEN: a running workspace with a rich parameter value that
	// satisfies the template regex validation.
	db, ps, sqlDB := dbtestutil.NewDBWithSQLDB(t)
	client, _, _ := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		Database:                 db,
		Pubsub:                   ps,
		IncludeProvisionerDaemon: true,
	})
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, first.OrganizationID,
		echoResponsesWithRichParameter("foo", echoResponseOptions{
			validationRegex: "^good$",
		}),
	)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, first.OrganizationID, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID, func(request *optimus-ide-collabsdk.CreateWorkspaceRequest) {
		request.RichParameterValues = []optimus-ide-collabsdk.WorkspaceBuildParameter{
			{Name: "foo", Value: "good"},
		}
	})
	initialBuild := optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, initialBuild.Status)

	// WHEN: a stop build is created with an on_success child build
	// that has an invalid rich parameter value.
	ctx := testutil.Context(t, testutil.WaitLong)
	user, err := client.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)

	stopBuild, err := client.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
		OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
			Transition:        optimus-ide-collabsdk.WorkspaceTransitionStart,
			TemplateVersionID: template.ActiveVersionID,
			// The invalid value triggers a non-retryable child build
			// creation error when the orchestrator processes the row.
			RichParameterValues: []optimus-ide-collabsdk.WorkspaceBuildParameter{
				{Name: "foo", Value: "bad"},
			},
		},
	})
	require.NoError(t, err)

	stopBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, stopBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, stopBuild.Job.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusStopped, stopBuild.Status)

	// THEN: the server marks the orchestration as failed without
	// retrying or creating the child start build.
	var orchestration database.WorkspaceBuildOrchestration
	require.Eventually(t, func() bool {
		orchestration, err = dbtestutil.GetWorkspaceBuildOrchestrationByParentBuildID(ctx, sqlDB, stopBuild.ID)
		return err == nil &&
			orchestration.Status == "failed" &&
			!orchestration.ChildBuildID.Valid &&
			orchestration.AttemptCount == 0 &&
			!orchestration.NextRetryAfter.Valid &&
			orchestration.Error.Valid
	}, testutil.WaitShort, testutil.IntervalFast)
	require.Contains(t, orchestration.Error.String, "Unable to validate parameters")

	_, err = client.WorkspaceBuildByUsernameAndWorkspaceNameAndBuildNumber(
		ctx,
		user.Username,
		workspace.Name,
		strconv.FormatInt(int64(stopBuild.BuildNumber+1), 10),
	)
	var apiErr *optimus-ide-collabsdk.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
}

func TestPostWorkspaceBuildsOnSuccessRetryableChildBuildFailureDoesNotBlockLaterRestart(t *testing.T) {
	t.Parallel()

	// GIVEN: two provisioners, one holding a template import job
	// open to make the child build fail retryably while the other
	// processes workspace builds.
	db, ps, sqlDB := dbtestutil.NewDBWithSQLDB(t)
	client, _, api := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		Database:                 db,
		Pubsub:                   ps,
		IncludeProvisionerDaemon: true,
	})
	optimus-ide-collabdtest.NewProvisionerDaemon(t, api)

	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, first.OrganizationID, nil)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, first.OrganizationID, version.ID)

	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
	startBuild := optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, startBuild.Status)

	blockedVersion := optimus-ide-collabdtest.UpdateTemplateVersion(t, client, first.OrganizationID,
		// Without a PlanComplete response, the echo provisioner will
		// keep the template version import job running.
		&echo.Responses{
			Parse: echo.ParseComplete,
			ProvisionPlan: []*proto.Response{{
				Type: &proto.Response_Log{
					Log: &proto.Log{},
				},
			}},
		}, template.ID,
	)
	optimus-ide-collabdtest.AwaitTemplateVersionJobRunning(t, client, blockedVersion.ID)

	// WHEN: a stop build is created with an on_success child build
	// request that references a template version whose import job is
	// still running, causing child build creation to be retried later.
	ctx := testutil.Context(t, testutil.WaitLong)
	badStopBuild, err := client.CreateWorkspaceBuild(ctx, workspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
		OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
			Transition:        optimus-ide-collabsdk.WorkspaceTransitionStart,
			TemplateVersionID: blockedVersion.ID,
		},
	})
	require.NoError(t, err)

	badStopBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, badStopBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, badStopBuild.Job.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusStopped, badStopBuild.Status)

	// THEN: the orchestrator records a delayed retry after child
	// build creation fails because the requested template version
	// is still importing.
	var badOrchestration database.WorkspaceBuildOrchestration
	require.Eventually(t, func() bool {
		badOrchestration, err = dbtestutil.GetWorkspaceBuildOrchestrationByParentBuildID(ctx, sqlDB, badStopBuild.ID)
		return err == nil &&
			badOrchestration.Status == "pending" &&
			badOrchestration.AttemptCount == 1 &&
			badOrchestration.NextRetryAfter.Valid
	}, testutil.WaitShort, testutil.IntervalFast)
	require.False(t, badOrchestration.ChildBuildID.Valid)
	require.True(t, badOrchestration.Error.Valid)
	require.Contains(t, badOrchestration.Error.String, "template version is running")

	// WHEN: a later restart uses a valid child build request.
	goodWorkspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
	goodStartBuild := optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, goodWorkspace.LatestBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, goodStartBuild.Status)

	user, err := client.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)

	goodStopBuild, err := client.CreateWorkspaceBuild(ctx, goodWorkspace.ID, optimus-ide-collabsdk.CreateWorkspaceBuildRequest{
		Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
		OnSuccess: &optimus-ide-collabsdk.CreateWorkspaceBuildOnSuccessRequest{
			Transition:        optimus-ide-collabsdk.WorkspaceTransitionStart,
			TemplateVersionID: template.ActiveVersionID,
		},
	})
	require.NoError(t, err)

	goodStopBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, goodStopBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, goodStopBuild.Job.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusStopped, goodStopBuild.Status)

	// THEN: the delayed retry row does not block the later orchestration.
	var goodChildBuild optimus-ide-collabsdk.WorkspaceBuild
	require.Eventually(t, func() bool {
		goodChildBuild, err = client.WorkspaceBuildByUsernameAndWorkspaceNameAndBuildNumber(
			ctx,
			user.Username,
			goodWorkspace.Name,
			strconv.FormatInt(int64(goodStopBuild.BuildNumber+1), 10),
		)
		return err == nil &&
			goodChildBuild.Transition == optimus-ide-collabsdk.WorkspaceTransitionStart
	}, testutil.WaitMedium, testutil.IntervalFast)

	goodChildBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, goodChildBuild.ID)
	require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, goodChildBuild.Job.Status)
	require.Equal(t, optimus-ide-collabsdk.WorkspaceStatusRunning, goodChildBuild.Status)
}

type echoResponseOptions struct {
	blockStopApply  bool
	failStopApply   bool
	validationRegex string
}

func echoResponsesWithRichParameter(paramName string, options echoResponseOptions) *echo.Responses {
	validationError := ""
	if options.validationRegex != "" {
		validationError = "invalid parameter value"
	}

	responses := &echo.Responses{
		Parse:         echo.ParseComplete,
		ProvisionInit: echo.InitComplete,
		ProvisionGraph: []*proto.Response{{
			Type: &proto.Response_Graph{
				Graph: &proto.GraphComplete{
					Parameters: []*proto.RichParameter{{
						Name:            paramName,
						Type:            "string",
						DefaultValue:    "bar",
						Mutable:         true,
						FormType:        proto.ParameterFormType_INPUT,
						ValidationRegex: options.validationRegex,
						ValidationError: validationError,
					}},
				},
			},
		}},
		ProvisionPlan:  echo.PlanComplete,
		ProvisionApply: echo.ApplyComplete,
	}
	if options.blockStopApply {
		responses.ProvisionApplyMap = map[proto.WorkspaceTransition][]*proto.Response{
			proto.WorkspaceTransition_START: echo.ApplyComplete,
			proto.WorkspaceTransition_STOP: {{
				Type: &proto.Response_Log{
					Log: &proto.Log{},
				},
			}},
		}
	}
	if options.failStopApply {
		responses.ProvisionApplyMap = map[proto.WorkspaceTransition][]*proto.Response{
			proto.WorkspaceTransition_START: echo.ApplyComplete,
			proto.WorkspaceTransition_STOP:  echo.ApplyFailed,
		}
	}
	return responses
}
