package prebuilds_test

import (
	"io"
	"strconv"
	"sync"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scaletest/prebuilds"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/quartz"
)

func TestRun(t *testing.T) {
	t.Parallel()

	t.Skip("This test takes several minutes to run, and is intended as a manual regression test")

	ctx := testutil.Context(t, testutil.WaitSuperLong*3)

	client, user := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureWorkspacePrebuilds:         1,
				optimus-ide-collabsdk.FeatureExternalProvisionerDaemons: 1,
			},
		},
	})

	// This is a real Terraform provisioner
	_ = optimus-ide-collabdenttest.NewExternalProvisionerDaemonTerraform(t, client, user.OrganizationID, nil)

	numTemplates := 2
	numPresets := 1
	numPresetPrebuilds := 1

	//nolint:gocritic // It's fine to use the owner user to pause prebuilds
	err := client.PutPrebuildsSettings(ctx, optimus-ide-collabsdk.PrebuildsSettings{
		ReconciliationPaused: true,
	})
	require.NoError(t, err)

	setupBarrier := new(sync.WaitGroup)
	setupBarrier.Add(numTemplates)
	creationBarrier := new(sync.WaitGroup)
	creationBarrier.Add(numTemplates)
	deletionSetupBarrier := new(sync.WaitGroup)
	deletionSetupBarrier.Add(1)
	deletionBarrier := new(sync.WaitGroup)
	deletionBarrier.Add(numTemplates)

	metrics := prebuilds.NewMetrics(prometheus.NewRegistry())

	eg, runCtx := errgroup.WithContext(ctx)

	runners := make([]*prebuilds.Runner, 0, numTemplates)
	for i := range numTemplates {
		cfg := prebuilds.Config{
			OrganizationID:            user.OrganizationID,
			NumPresets:                numPresets,
			NumPresetPrebuilds:        numPresetPrebuilds,
			TemplateVersionJobTimeout: testutil.WaitSuperLong * 2,
			PrebuildWorkspaceTimeout:  testutil.WaitSuperLong * 2,
			Metrics:                   metrics,
			SetupBarrier:              setupBarrier,
			CreationBarrier:           creationBarrier,
			DeletionSetupBarrier:      deletionSetupBarrier,
			DeletionBarrier:           deletionBarrier,
			Clock:                     quartz.NewReal(),
		}
		err := cfg.Validate()
		require.NoError(t, err)

		runner := prebuilds.NewRunner(client, cfg)
		runners = append(runners, runner)
		eg.Go(func() error {
			return runner.Run(runCtx, strconv.Itoa(i), io.Discard)
		})
	}

	// Wait for all runners to reach the setup barrier (templates created)
	setupBarrier.Wait()

	// Resume prebuilds to trigger prebuild creation
	err = client.PutPrebuildsSettings(ctx, optimus-ide-collabsdk.PrebuildsSettings{
		ReconciliationPaused: false,
	})
	require.NoError(t, err)

	// Wait for all runners to reach the creation barrier (prebuilds created)
	creationBarrier.Wait()

	//nolint:gocritic // Owner user is fine here as we want to view all workspaces
	workspaces, err := client.Workspaces(ctx, optimus-ide-collabsdk.WorkspaceFilter{})
	require.NoError(t, err)
	expectedWorkspaces := numTemplates * numPresets * numPresetPrebuilds
	require.Equal(t, workspaces.Count, expectedWorkspaces)

	// Pause prebuilds before deletion setup
	err = client.PutPrebuildsSettings(ctx, optimus-ide-collabsdk.PrebuildsSettings{
		ReconciliationPaused: true,
	})
	require.NoError(t, err)

	// Signal runners that prebuilds are paused and they can prepare for deletion
	deletionSetupBarrier.Done()

	// Wait for all runners to reach the deletion barrier (template versions updated to 0 prebuilds)
	deletionBarrier.Wait()

	// Resume prebuilds to trigger prebuild deletion
	err = client.PutPrebuildsSettings(ctx, optimus-ide-collabsdk.PrebuildsSettings{
		ReconciliationPaused: false,
	})
	require.NoError(t, err)

	err = eg.Wait()
	require.NoError(t, err)

	//nolint:gocritic // Owner user is fine here as we want to view all workspaces
	workspaces, err = client.Workspaces(ctx, optimus-ide-collabsdk.WorkspaceFilter{})
	require.NoError(t, err)
	require.Equal(t, workspaces.Count, 0)

	cleanupEg, cleanupCtx := errgroup.WithContext(ctx)
	for i, runner := range runners {
		cleanupEg.Go(func() error {
			return runner.Cleanup(cleanupCtx, strconv.Itoa(i), io.Discard)
		})
	}

	err = cleanupEg.Wait()
	require.NoError(t, err)
}
