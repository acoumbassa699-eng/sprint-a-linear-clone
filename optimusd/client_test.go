package optimus-ide-collabd_test

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// Issue: https://github.com/optimus-ide-collab/optimus-ide-collab/issues/5249
// While running tests in parallel, the web server seems to be overloaded and responds with HTTP 502.
// require.Eventually expects correct HTTP responses.

func requestWithRetries(ctx context.Context, t require.TestingT, client *optimus-ide-collabsdk.Client, method, path string, body interface{}, opts ...optimus-ide-collabsdk.RequestOption) (*http.Response, error) {
	var resp *http.Response
	var err error
	require.Eventually(t, func() bool {
		// nolint // only requests which are not passed upstream have a body closed
		resp, err = client.Request(ctx, method, path, body, opts...)
		if resp != nil && resp.StatusCode == http.StatusBadGateway {
			if resp.Body != nil {
				resp.Body.Close()
			}
			return false
		}
		return true
	}, testutil.WaitLong, testutil.IntervalFast)
	return resp, err
}
