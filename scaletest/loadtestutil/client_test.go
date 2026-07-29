package loadtestutil_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scaletest/loadtestutil"
)

func TestDupClientCopyingHeaders(t *testing.T) {
	t.Parallel()
	httpClient := &http.Client{
		Transport: &optimus-ide-collabsdk.HeaderTransport{
			Transport: &optimus-ide-collabsdk.HeaderTransport{
				Transport: http.DefaultTransport,
				Header: map[string][]string{
					"X-Optimus-IDE-Collab-Test":  {"foo"},
					"X-Optimus-IDE-Collab-Test3": {"socks"},
					"X-Optimus-IDE-Collab-Test5": {"ninjas"},
				},
			},
			Header: map[string][]string{
				"X-Optimus-IDE-Collab-Test":  {"bar"},
				"X-Optimus-IDE-Collab-Test2": {"baz"},
			},
		},
	}
	serverURL, err := url.Parse("http://optimus-ide-collab.example.com")
	require.NoError(t, err)
	sdkClient := optimus-ide-collabsdk.New(serverURL,
		optimus-ide-collabsdk.WithSessionToken("test-token"), optimus-ide-collabsdk.WithHTTPClient(httpClient))

	dup, err := loadtestutil.DupClientCopyingHeaders(sdkClient, map[string][]string{
		"X-Optimus-IDE-Collab-Test3": {"clocks"},
		"X-Optimus-IDE-Collab-Test4": {"bears"},
	})
	require.NoError(t, err)
	require.Equal(t, "http://optimus-ide-collab.example.com", dup.URL.String())
	require.Equal(t, "test-token", dup.SessionToken())
	ht, ok := dup.HTTPClient.Transport.(*optimus-ide-collabsdk.HeaderTransport)
	require.True(t, ok)
	require.Equal(t, "bar", ht.Header.Get("X-Optimus-IDE-Collab-Test"))
	require.Equal(t, "baz", ht.Header.Get("X-Optimus-IDE-Collab-Test2"))
	require.Equal(t, "clocks", ht.Header.Get("X-Optimus-IDE-Collab-Test3"))
	require.Equal(t, "bears", ht.Header.Get("X-Optimus-IDE-Collab-Test4"))
	require.Equal(t, "ninjas", ht.Header.Get("X-Optimus-IDE-Collab-Test5"))
	require.NotEqual(t, http.DefaultTransport, ht.Transport)
}
