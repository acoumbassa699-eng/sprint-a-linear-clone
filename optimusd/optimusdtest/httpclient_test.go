package optimus-ide-collabdtest_test

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestNewIsolatedHTTPClient(t *testing.T) {
	t.Parallel()

	client := optimus-ide-collabdtest.NewIsolatedHTTPClient(testutil.MustURL(t, "http://example.com"))
	require.NotNil(t, client.Transport)
	require.NotSame(t, http.DefaultTransport, client.Transport)

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.Nil(t, transport.TLSClientConfig)
}

func TestNewIsolatedHTTPSClient(t *testing.T) {
	t.Parallel()

	client := optimus-ide-collabdtest.NewIsolatedHTTPClient(testutil.MustURL(t, "https://example.com"))
	require.NotSame(t, http.DefaultTransport, client.Transport)

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport.TLSClientConfig)
	require.True(t, transport.TLSClientConfig.InsecureSkipVerify)
	require.Equal(t, uint16(tls.VersionTLS12), transport.TLSClientConfig.MinVersion)
}

func TestNewIsolatedHTTPClientNilURL(t *testing.T) {
	t.Parallel()

	client := optimus-ide-collabdtest.NewIsolatedHTTPClient(nil)
	require.NotNil(t, client.Transport)
	require.NotSame(t, http.DefaultTransport, client.Transport)

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.Nil(t, transport.TLSClientConfig)
}

func TestCreateAnotherUserHTTPClient(t *testing.T) {
	t.Parallel()

	client := optimus-ide-collabdtest.New(t, nil)
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	client.HTTPClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}

	other, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID)

	require.NotSame(t, client.HTTPClient, other.HTTPClient)
	require.Same(t, client.HTTPClient.Transport, other.HTTPClient.Transport)
	require.Nil(t, other.HTTPClient.CheckRedirect)
}

func TestCreateAnotherUserHTTPClientDefaultTransport(t *testing.T) {
	t.Parallel()

	client := optimus-ide-collabdtest.New(t, nil)
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	base := optimus-ide-collabsdk.New(
		client.URL,
		optimus-ide-collabsdk.WithSessionToken(client.SessionToken()),
		optimus-ide-collabsdk.WithHTTPClient(&http.Client{Timeout: time.Second}),
	)

	other, _ := optimus-ide-collabdtest.CreateAnotherUser(t, base, first.OrganizationID)

	require.NotSame(t, base.HTTPClient, other.HTTPClient)
	require.NotNil(t, other.HTTPClient.Transport)
	require.NotSame(t, http.DefaultTransport, other.HTTPClient.Transport)
	require.Equal(t, base.HTTPClient.Timeout, other.HTTPClient.Timeout)
}
