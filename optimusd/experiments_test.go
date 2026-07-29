package optimus-ide-collabd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func Test_Experiments(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		cfg := optimus-ide-collabdtest.DeploymentValues(t)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues: cfg,
		})
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		experiments, err := client.Experiments(ctx)
		require.NoError(t, err)
		require.NotNil(t, experiments)
		require.Empty(t, experiments)
		require.False(t, experiments.Enabled("foo"))
	})

	t.Run("multiple features", func(t *testing.T) {
		t.Parallel()
		cfg := optimus-ide-collabdtest.DeploymentValues(t)
		cfg.Experiments = []string{"foo", "BAR"}
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues: cfg,
		})
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		experiments, err := client.Experiments(ctx)
		require.NoError(t, err)
		require.NotNil(t, experiments)
		// Should be lower-cased.
		require.ElementsMatch(t, []optimus-ide-collabsdk.Experiment{"foo", "bar"}, experiments)
		require.True(t, experiments.Enabled("foo"))
		require.True(t, experiments.Enabled("bar"))
		require.False(t, experiments.Enabled("baz"))
	})

	t.Run("wildcard", func(t *testing.T) {
		t.Parallel()
		cfg := optimus-ide-collabdtest.DeploymentValues(t)
		cfg.Experiments = []string{"*"}
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues: cfg,
		})
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		experiments, err := client.Experiments(ctx)
		require.NoError(t, err)
		require.NotNil(t, experiments)
		require.ElementsMatch(t, optimus-ide-collabsdk.ExperimentsSafe, experiments)
		for _, ex := range optimus-ide-collabsdk.ExperimentsSafe {
			require.True(t, experiments.Enabled(ex))
		}
		require.False(t, experiments.Enabled("danger"))
	})

	t.Run("alternate wildcard with manual opt-in", func(t *testing.T) {
		t.Parallel()
		cfg := optimus-ide-collabdtest.DeploymentValues(t)
		cfg.Experiments = []string{"*", "dAnGeR"}
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues: cfg,
		})
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		experiments, err := client.Experiments(ctx)
		require.NoError(t, err)
		require.NotNil(t, experiments)
		require.ElementsMatch(t, append(optimus-ide-collabsdk.ExperimentsSafe, "danger"), experiments)
		for _, ex := range optimus-ide-collabsdk.ExperimentsSafe {
			require.True(t, experiments.Enabled(ex))
		}
		require.True(t, experiments.Enabled("danger"))
		require.False(t, experiments.Enabled("herebedragons"))
	})

	t.Run("Unauthorized", func(t *testing.T) {
		t.Parallel()
		cfg := optimus-ide-collabdtest.DeploymentValues(t)
		cfg.Experiments = []string{"*"}
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues: cfg,
		})
		// Explicitly omit creating a user so we're unauthorized.
		// _ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.Experiments(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, httpmw.SignedOutErrorMessage)
	})

	t.Run("available experiments", func(t *testing.T) {
		t.Parallel()
		cfg := optimus-ide-collabdtest.DeploymentValues(t)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues: cfg,
		})
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		experiments, err := client.SafeExperiments(ctx)
		require.NoError(t, err)
		require.NotNil(t, experiments)
		require.ElementsMatch(t, optimus-ide-collabsdk.ExperimentsSafe, experiments.Safe)
	})
}
