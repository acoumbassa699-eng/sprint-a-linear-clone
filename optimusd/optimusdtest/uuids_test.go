package optimus-ide-collabdtest_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
)

func TestDeterministicUUIDGenerator(t *testing.T) {
	t.Parallel()

	ids := optimus-ide-collabdtest.NewDeterministicUUIDGenerator()
	require.Equal(t, ids.ID("g1"), ids.ID("g1"))
	require.NotEqual(t, ids.ID("g1"), ids.ID("g2"))
}
