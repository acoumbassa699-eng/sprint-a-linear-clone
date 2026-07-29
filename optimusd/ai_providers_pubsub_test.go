package optimus-ide-collabd_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	optimus-ide-collabpubsub "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/pubsub"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// TestAIProvidersChangedPubsub asserts that the CRUD handlers publish
// on AIProvidersChangedChannel for the operations that affect the
// runtime provider set. Subscribers (aibridged, aibridgeproxyd) depend
// on these notifications to trigger their pool reload.
//
// The handlers publish best-effort and the payload is empty, so we
// assert "at least one event per mutation" via a counter.
func TestAIProvidersChangedPubsub(t *testing.T) {
	t.Parallel()

	client, _, api := optimus-ide-collabdtest.NewWithAPI(t, nil)
	_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
	ctx := testutil.Context(t, testutil.WaitLong)

	var count atomic.Int64
	unsubscribe, err := api.Pubsub.Subscribe(optimus-ide-collabpubsub.AIProvidersChangedChannel, func(_ context.Context, _ []byte) {
		count.Add(1)
	})
	require.NoError(t, err)
	t.Cleanup(unsubscribe)

	// Create.
	req := optimus-ide-collabsdk.CreateAIProviderRequest{
		Type:    optimus-ide-collabsdk.AIProviderTypeOpenAI,
		Name:    "pubsub-openai",
		Enabled: true,
		BaseURL: "https://api.openai.com/v1/",
		APIKeys: []string{"k1"},
	}
	//nolint:gocritic // Owner role is the audience for this endpoint.
	created, err := client.CreateAIProvider(ctx, req)
	require.NoError(t, err)
	testutil.Eventually(ctx, t, func(_ context.Context) bool { return count.Load() >= 1 }, testutil.IntervalFast)

	// Update.
	newKey := "k2"
	_, err = client.UpdateAIProvider(ctx, created.ID.String(), optimus-ide-collabsdk.UpdateAIProviderRequest{
		APIKeys: &[]optimus-ide-collabsdk.AIProviderKeyMutation{{APIKey: &newKey}},
	})
	require.NoError(t, err)
	testutil.Eventually(ctx, t, func(_ context.Context) bool { return count.Load() >= 2 }, testutil.IntervalFast)

	// Delete.
	err = client.DeleteAIProvider(ctx, created.ID.String())
	require.NoError(t, err)
	testutil.Eventually(ctx, t, func(_ context.Context) bool { return count.Load() >= 3 }, testutil.IntervalFast)
}
