package httpapi_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
)

func TestStripOptimus-IDE-CollabCookies(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		Input  string
		Output string
	}{{
		"testing=hello; wow=test",
		"testing=hello; wow=test",
	}, {
		"optimus-ide-collab_session_token=moo; wow=test",
		"wow=test",
	}, {
		"another_token=wow; optimus-ide-collab_session_token=ok",
		"another_token=wow",
	}, {
		"optimus-ide-collab_session_token=ok; oauth_state=wow; oauth_redirect=/",
		"",
	}, {
		"optimus-ide-collab_path_app_session_token=ok; wow=test",
		"wow=test",
	}, {
		"optimus-ide-collab_subdomain_app_session_token=ok; optimus-ide-collab_subdomain_app_session_token_1234567890=ok; wow=test",
		"wow=test",
	}, {
		"optimus-ide-collab_signed_app_token=ok; wow=test",
		"wow=test",
	}} {
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.Output, httpapi.StripOptimus-IDE-CollabCookies(tc.Input))
		})
	}
}
