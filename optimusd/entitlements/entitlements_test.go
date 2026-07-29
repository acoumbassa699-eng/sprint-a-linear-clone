package entitlements_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/entitlements"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestModify(t *testing.T) {
	t.Parallel()

	set := entitlements.New()
	require.False(t, set.Enabled(optimus-ide-collabsdk.FeatureMultipleOrganizations))

	set.Modify(func(entitlements *optimus-ide-collabsdk.Entitlements) {
		entitlements.Features[optimus-ide-collabsdk.FeatureMultipleOrganizations] = optimus-ide-collabsdk.Feature{
			Enabled:     true,
			Entitlement: optimus-ide-collabsdk.EntitlementEntitled,
		}
	})
	require.True(t, set.Enabled(optimus-ide-collabsdk.FeatureMultipleOrganizations))
}

func TestAllowRefresh(t *testing.T) {
	t.Parallel()

	now := time.Now()
	set := entitlements.New()
	set.Modify(func(entitlements *optimus-ide-collabsdk.Entitlements) {
		entitlements.RefreshedAt = now
	})

	ok, wait := set.AllowRefresh(now)
	require.False(t, ok)
	require.InDelta(t, time.Minute.Seconds(), wait.Seconds(), 5)

	set.Modify(func(entitlements *optimus-ide-collabsdk.Entitlements) {
		entitlements.RefreshedAt = now.Add(time.Minute * -2)
	})

	ok, wait = set.AllowRefresh(now)
	require.True(t, ok)
	require.Equal(t, time.Duration(0), wait)
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t, testutil.WaitShort)

	set := entitlements.New()
	require.False(t, set.Enabled(optimus-ide-collabsdk.FeatureMultipleOrganizations))
	fetchStarted := make(chan struct{})
	firstDone := make(chan struct{})
	errCh := make(chan error, 2)
	go func() {
		err := set.Update(ctx, func(_ context.Context) (optimus-ide-collabsdk.Entitlements, error) {
			close(fetchStarted)
			select {
			case <-firstDone:
				// OK!
			case <-ctx.Done():
				t.Error("timeout")
				return optimus-ide-collabsdk.Entitlements{}, ctx.Err()
			}
			return optimus-ide-collabsdk.Entitlements{
				Features: map[optimus-ide-collabsdk.FeatureName]optimus-ide-collabsdk.Feature{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: {
						Enabled: true,
					},
				},
			}, nil
		})
		errCh <- err
	}()
	testutil.TryReceive(ctx, t, fetchStarted)
	require.False(t, set.Enabled(optimus-ide-collabsdk.FeatureMultipleOrganizations))
	// start a second update while the first one is in progress
	go func() {
		err := set.Update(ctx, func(_ context.Context) (optimus-ide-collabsdk.Entitlements, error) {
			return optimus-ide-collabsdk.Entitlements{
				Features: map[optimus-ide-collabsdk.FeatureName]optimus-ide-collabsdk.Feature{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: {
						Enabled: true,
					},
					optimus-ide-collabsdk.FeatureAppearance: {
						Enabled: true,
					},
				},
			}, nil
		})
		errCh <- err
	}()
	close(firstDone)
	err := testutil.TryReceive(ctx, t, errCh)
	require.NoError(t, err)
	err = testutil.TryReceive(ctx, t, errCh)
	require.NoError(t, err)
	require.True(t, set.Enabled(optimus-ide-collabsdk.FeatureMultipleOrganizations))
	require.True(t, set.Enabled(optimus-ide-collabsdk.FeatureAppearance))
}

func TestUpdate_LicenseRequiresTelemetry(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t, testutil.WaitShort)
	set := entitlements.New()
	set.Modify(func(entitlements *optimus-ide-collabsdk.Entitlements) {
		entitlements.Errors = []string{"some error"}
		entitlements.Features[optimus-ide-collabsdk.FeatureAppearance] = optimus-ide-collabsdk.Feature{
			Enabled: true,
		}
	})
	err := set.Update(ctx, func(_ context.Context) (optimus-ide-collabsdk.Entitlements, error) {
		return optimus-ide-collabsdk.Entitlements{}, entitlements.ErrLicenseRequiresTelemetry
	})
	require.NoError(t, err)
	require.True(t, set.Enabled(optimus-ide-collabsdk.FeatureAppearance))
	require.Equal(t, []string{entitlements.ErrLicenseRequiresTelemetry.Error()}, set.Errors())
}
