package httpmw_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/proxyhealth"
)

func TestCSPFrameAncestors(t *testing.T) {
	t.Parallel()

	t.Run("DefaultSelf", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		rw := httptest.NewRecorder()

		httpmw.CSPHeaders(false, func() []*proxyhealth.ProxyHost {
			return nil
		}, nil)(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
		})).ServeHTTP(rw, r)

		csp := rw.Header().Get("Content-Security-Policy")
		require.Contains(t, csp, "frame-ancestors 'self'")
	})

	t.Run("OverrideViaStaticAdditions", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		rw := httptest.NewRecorder()

		httpmw.CSPHeaders(false, func() []*proxyhealth.ProxyHost {
			return nil
		}, map[httpmw.CSPFetchDirective][]string{
			httpmw.CSPFrameAncestors: {"https://example.com"},
		})(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
		})).ServeHTTP(rw, r)

		csp := rw.Header().Get("Content-Security-Policy")
		require.Contains(t, csp, "frame-ancestors https://example.com")
		require.NotContains(t, csp, "frame-ancestors 'self'")
	})

	t.Run("OmitWhenEmpty", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		rw := httptest.NewRecorder()

		httpmw.CSPHeaders(false, func() []*proxyhealth.ProxyHost {
			return nil
		}, map[httpmw.CSPFetchDirective][]string{
			httpmw.CSPFrameAncestors: {},
		})(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
		})).ServeHTTP(rw, r)

		csp := rw.Header().Get("Content-Security-Policy")
		require.NotContains(t, csp, "frame-ancestors")
	})
}

func TestCSP(t *testing.T) {
	t.Parallel()

	proxyHosts := []*proxyhealth.ProxyHost{
		{
			Host:    "test.com",
			AppHost: "*.test.com",
		},
		{
			Host:    "optimus-ide-collab.com",
			AppHost: "*.optimus-ide-collab.com",
		},
		{
			// Host is not added because it duplicates the host header.
			Host:    "example.com",
			AppHost: "*.optimus-ide-collab2.com",
		},
	}
	expectedMedia := []string{"media.com", "media2.com"}

	expected := []string{
		"frame-src 'self' *.test.com *.optimus-ide-collab.com *.optimus-ide-collab2.com",
		"media-src 'self' " + strings.Join(expectedMedia, " "),
		strings.Join([]string{
			"connect-src", "'self'",
			// Added from host header.
			"wss://example.com", "ws://example.com",
			// Added via proxy hosts.
			"wss://test.com", "ws://test.com", "https://test.com", "http://test.com",
			"wss://optimus-ide-collab.com", "ws://optimus-ide-collab.com", "https://optimus-ide-collab.com", "http://optimus-ide-collab.com",
		}, " "),
	}

	// When the host is empty, it uses example.com.
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()

	httpmw.CSPHeaders(false, func() []*proxyhealth.ProxyHost {
		return proxyHosts
	}, map[httpmw.CSPFetchDirective][]string{
		httpmw.CSPDirectiveMediaSrc: expectedMedia,
	})(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})).ServeHTTP(rw, r)

	require.NotEmpty(t, rw.Header().Get("Content-Security-Policy"), "Content-Security-Policy header should not be empty")
	for _, e := range expected {
		require.Contains(t, rw.Header().Get("Content-Security-Policy"), e)
	}
}
