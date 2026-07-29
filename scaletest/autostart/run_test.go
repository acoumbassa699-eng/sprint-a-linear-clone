package autostart_test

import (
	"io"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisionersdk/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scaletest/autostart"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scaletest/createusers"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scaletest/loadtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scaletest/workspacebuild"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestRun(t *testing.T) {
	t.Parallel()
	numUsers := 2
	autoStartDelay := 2 * time.Minute

	// Faking a workspace autostart schedule start time at the optimus-ide-collabd level
	// is difficult and error-prone. This test verifies the setup phase only
	// (creating workspaces, stopping them, and configuring autostart schedules).
	t.Skip("This test takes several minutes to run, and is intended as a manual regression test")

	ctx := testutil.Context(t, time.Minute*3)

	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		IncludeProvisionerDaemon: true,
		AutobuildTicker:          time.NewTicker(time.Second * 1).C,
		DeploymentValues: optimus-ide-collabdtest.DeploymentValues(t, func(dv *optimus-ide-collabsdk.DeploymentValues) {
			dv.Experiments = []string{string(optimus-ide-collabsdk.ExperimentWorkspaceBuildUpdates)}
		}),
	})
	user := optimus-ide-collabdtest.CreateFirstUser(t, client)

	authToken := uuid.NewString()
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, &echo.Responses{
		Parse:         echo.ParseComplete,
		ProvisionPlan: echo.PlanComplete,
		ProvisionGraph: []*proto.Response{
			{
				Type: &proto.Response_Graph{
					Graph: &proto.GraphComplete{
						Resources: []*proto.Resource{
							{
								Name: "example",
								Type: "aws_instance",
								Agents: []*proto.Agent{
									{
										Id:   uuid.NewString(),
										Name: "agent",
										Auth: &proto.Agent_Token{
											Token: authToken,
										},
										Apps: []*proto.App{},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)

	barrier := new(sync.WaitGroup)
	barrier.Add(numUsers)

	// Pre-create channels for each workspace keyed by deterministic name.
	workspaceChannels := make(map[string]chan optimus-ide-collabsdk.WorkspaceBuildUpdate)
	for i := range numUsers {
		id := strconv.Itoa(i)
		workspaceName := loadtestutil.GenerateDeterministicWorkspaceName(id)
		workspaceChannels[workspaceName] = make(chan optimus-ide-collabsdk.WorkspaceBuildUpdate, 16)
	}

	// Start watching all workspace builds.
	deoptimus-ide-collab, err := client.WatchAllWorkspaceBuilds(ctx)
	require.NoError(t, err)
	defer deoptimus-ide-collab.Close()

	// Start the dispatcher goroutine.
	go func() {
		for update := range deoptimus-ide-collab.Chan() {
			if ch, ok := workspaceChannels[update.WorkspaceName]; ok {
				select {
				case ch <- update:
				case <-ctx.Done():
					return
				}
			}
		}
		for _, ch := range workspaceChannels {
			close(ch)
		}
	}()

	eg, runCtx := errgroup.WithContext(ctx)

	runners := make([]*autostart.Runner, 0, numUsers)
	for i := range numUsers {
		id := strconv.Itoa(i)
		workspaceName := loadtestutil.GenerateDeterministicWorkspaceName(id)
		cfg := autostart.Config{
			User: createusers.Config{
				OrganizationID: user.OrganizationID,
			},
			Workspace: workspacebuild.Config{
				OrganizationID: user.OrganizationID,
				Request: optimus-ide-collabsdk.CreateWorkspaceRequest{
					TemplateID: template.ID,
					Name:       workspaceName,
				},
				NoWaitForAgents: true,
			},
			WorkspaceJobTimeout: testutil.WaitMedium,
			AutostartDelay:      autoStartDelay,
			SetupBarrier:        barrier,
			BuildUpdates:        workspaceChannels[workspaceName],
		}
		err := cfg.Validate()
		require.NoError(t, err)

		runner := autostart.NewRunner(client, cfg)
		runners = append(runners, runner)
		eg.Go(func() error {
			return runner.Run(runCtx, strconv.Itoa(i), io.Discard)
		})
	}

	err = eg.Wait()
	require.NoError(t, err)

	users, err := client.Users(ctx, optimus-ide-collabsdk.UsersRequest{})
	require.NoError(t, err)
	require.Len(t, users.Users, 1+numUsers) // owner + created users

	workspaces, err := client.Workspaces(ctx, optimus-ide-collabsdk.WorkspaceFilter{})
	require.NoError(t, err)
	require.Len(t, workspaces.Workspaces, numUsers) // one workspace per user

	// Verify that workspaces have autostart schedules set and are stopped
	// (the test exits after configuring autostart, before it triggers).
	for _, workspace := range workspaces.Workspaces {
		require.NotNil(t, workspace.AutostartSchedule)
		require.Equal(t, optimus-ide-collabsdk.WorkspaceTransitionStop, workspace.LatestBuild.Transition)
		require.Equal(t, optimus-ide-collabsdk.ProvisionerJobSucceeded, workspace.LatestBuild.Job.Status)
	}

	cleanupEg, cleanupCtx := errgroup.WithContext(ctx)
	for i, runner := range runners {
		cleanupEg.Go(func() error {
			return runner.Cleanup(cleanupCtx, strconv.Itoa(i), io.Discard)
		})
	}
	err = cleanupEg.Wait()
	require.NoError(t, err)

	workspaces, err = client.Workspaces(ctx, optimus-ide-collabsdk.WorkspaceFilter{})
	require.NoError(t, err)
	require.Len(t, workspaces.Workspaces, 0)

	users, err = client.Users(ctx, optimus-ide-collabsdk.UsersRequest{})
	require.NoError(t, err)
	require.Len(t, users.Users, 1) // owner
}
