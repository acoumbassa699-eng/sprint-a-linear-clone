package chatdebug_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd/chatdebug"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// toStrings converts a typed string slice to []string for comparison.
func toStrings[T ~string](values []T) []string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = string(v)
	}
	return out
}

// TestTypesMatchSDK verifies that every chatdebug constant has a
// corresponding optimus-ide-collabsdk constant with the same string value.
// If this test fails you probably added a constant to one package
// but forgot to update the other.
func TestTypesMatchSDK(t *testing.T) {
	t.Parallel()

	t.Run("RunKind", func(t *testing.T) {
		t.Parallel()
		require.ElementsMatch(t,
			toStrings(chatdebug.AllRunKinds),
			toStrings(optimus-ide-collabsdk.AllChatDebugRunKinds),
			"chatdebug.AllRunKinds and optimus-ide-collabsdk.AllChatDebugRunKinds have diverged",
		)
	})

	t.Run("Status", func(t *testing.T) {
		t.Parallel()
		require.ElementsMatch(t,
			toStrings(chatdebug.AllStatuses),
			toStrings(optimus-ide-collabsdk.AllChatDebugStatuses),
			"chatdebug.AllStatuses and optimus-ide-collabsdk.AllChatDebugStatuses have diverged",
		)
	})

	t.Run("Operation", func(t *testing.T) {
		t.Parallel()
		require.ElementsMatch(t,
			toStrings(chatdebug.AllOperations),
			toStrings(optimus-ide-collabsdk.AllChatDebugStepOperations),
			"chatdebug.AllOperations and optimus-ide-collabsdk.AllChatDebugStepOperations have diverged",
		)
	})
}
