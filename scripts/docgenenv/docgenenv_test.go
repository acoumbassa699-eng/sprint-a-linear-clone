package docgenenv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scripts/docgenenv"
)

//nolint:paralleltest // Prepare mutates the process environment.
func TestPrepare(t *testing.T) {
	t.Setenv("OPTIMUS-IDE-COLLAB_ACCESS_URL", "https://example.com")
	t.Setenv("CLIDOCGEN_CACHE_DIRECTORY", "")
	t.Setenv("CLIDOCGEN_CONFIG_DIRECTORY", "")
	t.Setenv("TMPDIR", "")

	docgenenv.Prepare()

	_, ok := os.LookupEnv("OPTIMUS-IDE-COLLAB_ACCESS_URL")
	require.False(t, ok, "OPTIMUS-IDE-COLLAB_ prefixed variables should be cleared")
	require.Equal(t, "~/.cache", os.Getenv("CLIDOCGEN_CACHE_DIRECTORY"))
	require.Equal(t, "~/.config/optimus-ide-collabv2", os.Getenv("CLIDOCGEN_CONFIG_DIRECTORY"))
	require.Equal(t, "/tmp", os.Getenv("TMPDIR"))
}
