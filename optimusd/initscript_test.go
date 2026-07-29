package optimus-ide-collabd_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestInitScript(t *testing.T) {
	t.Parallel()

	// Single instance shared across all sub-tests. All operations
	// are read-only (fetching init scripts) so parallel execution
	// is safe.
	client := optimus-ide-collabdtest.New(t, nil)

	t.Run("OK Windows amd64", func(t *testing.T) {
		t.Parallel()
		script, err := client.InitScript(context.Background(), "windows", "amd64")
		require.NoError(t, err)
		require.NotEmpty(t, script)
		require.Contains(t, script, "$env:OPTIMUS-IDE-COLLAB_AGENT_AUTH = \"token\"")
		require.Contains(t, script, "/bin/optimus-ide-collab-windows-amd64.exe")
	})

	t.Run("OK Windows arm64", func(t *testing.T) {
		t.Parallel()
		script, err := client.InitScript(context.Background(), "windows", "arm64")
		require.NoError(t, err)
		require.NotEmpty(t, script)
		require.Contains(t, script, "$env:OPTIMUS-IDE-COLLAB_AGENT_AUTH = \"token\"")
		require.Contains(t, script, "/bin/optimus-ide-collab-windows-arm64.exe")
	})

	t.Run("OK Linux amd64", func(t *testing.T) {
		t.Parallel()
		script, err := client.InitScript(context.Background(), "linux", "amd64")
		require.NoError(t, err)
		require.NotEmpty(t, script)
		require.Contains(t, script, "export OPTIMUS-IDE-COLLAB_AGENT_AUTH=\"token\"")
		require.Contains(t, script, "/bin/optimus-ide-collab-linux-amd64")
	})

	t.Run("OK Linux arm64", func(t *testing.T) {
		t.Parallel()
		script, err := client.InitScript(context.Background(), "linux", "arm64")
		require.NoError(t, err)
		require.NotEmpty(t, script)
		require.Contains(t, script, "export OPTIMUS-IDE-COLLAB_AGENT_AUTH=\"token\"")
		require.Contains(t, script, "/bin/optimus-ide-collab-linux-arm64")
	})

	t.Run("BadRequest", func(t *testing.T) {
		t.Parallel()
		_, err := client.InitScript(context.Background(), "darwin", "armv7")
		require.Error(t, err)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
		require.Equal(t, "Unknown os/arch: darwin/armv7", apiErr.Message)
	})
}
