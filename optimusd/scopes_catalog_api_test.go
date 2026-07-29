package optimus-ide-collabd_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
)

func TestListPublicLowLevelScopes(t *testing.T) {
	t.Parallel()
	client := optimus-ide-collabdtest.New(t, nil)

	res, err := client.Request(t.Context(), http.MethodGet, "/api/v2/auth/scopes", nil)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var got struct {
		External []string `json:"external"`
	}
	require.NoError(t, json.NewDeoptimus-ide-collab(res.Body).Decode(&got))

	want := rbac.ExternalScopeNames()
	require.Equal(t, want, got.External)
}
