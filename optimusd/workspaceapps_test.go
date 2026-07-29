package optimus-ide-collabd_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/cryptokeys"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/jwtutils"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/workspaceapps"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/quartz"
)

func TestGetAppHost(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		accessURL   string
		appHostname string
		expected    string
	}{
		{
			name:        "OK",
			accessURL:   "https://test.optimus-ide-collab.com",
			appHostname: "*.test.optimus-ide-collab.com",
			expected:    "*.test.optimus-ide-collab.com",
		},
		{
			name:        "None",
			accessURL:   "https://test.optimus-ide-collab.com",
			appHostname: "",
			expected:    "",
		},
		{
			name:        "OKWithPort",
			accessURL:   "https://test.optimus-ide-collab.com:8443",
			appHostname: "*.test.optimus-ide-collab.com",
			expected:    "*.test.optimus-ide-collab.com:8443",
		},
		{
			name:        "OKWithSuffix",
			accessURL:   "https://test.optimus-ide-collab.com:8443",
			appHostname: "*--suffix.test.optimus-ide-collab.com",
			expected:    "*--suffix.test.optimus-ide-collab.com:8443",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			accessURL, err := url.Parse(c.accessURL)
			require.NoError(t, err)

			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
				AccessURL:   accessURL,
				AppHostname: c.appHostname,
			})

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			// Should not leak to unauthenticated users.
			host, err := client.AppHost(ctx)
			require.Error(t, err)
			require.Equal(t, "", host.Host)

			_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
			host, err = client.AppHost(ctx)
			require.NoError(t, err)
			require.Equal(t, c.expected, host.Host)
		})
	}
}

func TestWorkspaceApplicationAuth(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name             string
		accessURL        string
		appHostname      string
		proxyURL         string
		proxyAppHostname string

		redirectURI    string
		expectRedirect string
	}{
		{
			name:             "OK",
			accessURL:        "https://test.optimus-ide-collab.com",
			appHostname:      "*.test.optimus-ide-collab.com",
			proxyURL:         "https://proxy.test.optimus-ide-collab.com",
			proxyAppHostname: "*.proxy.test.optimus-ide-collab.com",
			redirectURI:      "https://something.test.optimus-ide-collab.com",
			expectRedirect:   "https://something.test.optimus-ide-collab.com",
		},
		{
			name:             "ProxyPathOK",
			accessURL:        "https://test.optimus-ide-collab.com",
			appHostname:      "*.test.optimus-ide-collab.com",
			proxyURL:         "https://proxy.test.optimus-ide-collab.com",
			proxyAppHostname: "*.proxy.test.optimus-ide-collab.com",
			redirectURI:      "https://proxy.test.optimus-ide-collab.com/path",
			expectRedirect:   "https://proxy.test.optimus-ide-collab.com/path",
		},
		{
			name:             "RejectProxyAccessURLPrefix",
			accessURL:        "https://test.optimus-ide-collab.com",
			appHostname:      "*.test.optimus-ide-collab.com",
			proxyURL:         "https://proxy.test.optimus-ide-collab.com",
			proxyAppHostname: "*.proxy.test.optimus-ide-collab.com",
			redirectURI:      "https://proxy.test.optimus-ide-collab/path",
			expectRedirect:   "",
		},
		{
			name:             "ProxySubdomainOK",
			accessURL:        "https://test.optimus-ide-collab.com",
			appHostname:      "*.test.optimus-ide-collab.com",
			proxyURL:         "https://proxy.test.optimus-ide-collab.com",
			proxyAppHostname: "*.proxy.test.optimus-ide-collab.com",
			redirectURI:      "https://something.proxy.test.optimus-ide-collab.com/path?yeah=true",
			expectRedirect:   "https://something.proxy.test.optimus-ide-collab.com/path?yeah=true",
		},
		{
			name:             "ProxySubdomainSuffixOK",
			accessURL:        "https://test.optimus-ide-collab.com",
			appHostname:      "*.test.optimus-ide-collab.com",
			proxyURL:         "https://proxy.test.optimus-ide-collab.com",
			proxyAppHostname: "*--suffix.proxy.test.optimus-ide-collab.com",
			redirectURI:      "https://something--suffix.proxy.test.optimus-ide-collab.com/",
			expectRedirect:   "https://something--suffix.proxy.test.optimus-ide-collab.com/",
		},
		{
			name:             "NormalizeSchemePrimaryAppHostname",
			accessURL:        "https://test.optimus-ide-collab.com",
			appHostname:      "*.test.optimus-ide-collab.com",
			proxyURL:         "https://proxy.test.optimus-ide-collab.com",
			proxyAppHostname: "*.proxy.test.optimus-ide-collab.com",
			redirectURI:      "http://x.test.optimus-ide-collab.com",
			expectRedirect:   "https://x.test.optimus-ide-collab.com",
		},
		{
			name:             "NormalizeSchemeProxyAppHostname",
			accessURL:        "https://test.optimus-ide-collab.com",
			appHostname:      "*.test.optimus-ide-collab.com",
			proxyURL:         "https://proxy.test.optimus-ide-collab.com",
			proxyAppHostname: "*.proxy.test.optimus-ide-collab.com",
			redirectURI:      "http://x.proxy.test.optimus-ide-collab.com",
			expectRedirect:   "https://x.proxy.test.optimus-ide-collab.com",
		},
		{
			name:             "NoneError",
			accessURL:        "https://test.optimus-ide-collab.com",
			appHostname:      "*.test.optimus-ide-collab.com",
			proxyURL:         "https://proxy.test.optimus-ide-collab.com",
			proxyAppHostname: "*.proxy.test.optimus-ide-collab.com",
			redirectURI:      "",
			expectRedirect:   "",
		},
		{
			name:             "PrimaryAccessURLError",
			accessURL:        "https://test.optimus-ide-collab.com",
			appHostname:      "*.test.optimus-ide-collab.com",
			proxyURL:         "https://proxy.test.optimus-ide-collab.com",
			proxyAppHostname: "*.proxy.test.optimus-ide-collab.com",
			redirectURI:      "https://test.optimus-ide-collab.com/",
			expectRedirect:   "",
		},
		{
			name:             "OtherError",
			accessURL:        "https://test.optimus-ide-collab.com",
			appHostname:      "*.test.optimus-ide-collab.com",
			proxyURL:         "https://proxy.test.optimus-ide-collab.com",
			proxyAppHostname: "*.proxy.test.optimus-ide-collab.com",
			redirectURI:      "https://example.com/",
			expectRedirect:   "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			ctx := testutil.Context(t, testutil.WaitMedium)
			logger := testutil.Logger(t)
			accessURL, err := url.Parse(c.accessURL)
			require.NoError(t, err)

			db, ps := dbtestutil.NewDB(t)
			fetcher := &cryptokeys.DBFetcher{
				DB: db,
			}

			kc, err := cryptokeys.NewEncryptionCache(ctx, logger, fetcher, optimus-ide-collabsdk.CryptoKeyFeatureWorkspaceAppsAPIKey)
			require.NoError(t, err)

			clock := quartz.NewMock(t)

			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
				AccessURL:             accessURL,
				AppHostname:           c.appHostname,
				Database:              db,
				Pubsub:                ps,
				APIKeyEncryptionCache: kc,
				Clock:                 clock,
			})
			_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

			// Disable redirects.
			client.HTTPClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			}

			_, _ = dbgen.WorkspaceProxy(t, db, database.WorkspaceProxy{
				Url:              c.proxyURL,
				WildcardHostname: c.proxyAppHostname,
			})

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			resp, err := client.Request(ctx, http.MethodGet, "/api/v2/applications/auth-redirect", nil, func(req *http.Request) {
				q := req.URL.Query()
				q.Set("redirect_uri", c.redirectURI)
				req.URL.RawQuery = q.Encode()
			})
			require.NoError(t, err)
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusSeeOther {
				err = optimus-ide-collabsdk.ReadBodyAsError(resp)
				if c.expectRedirect == "" {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				return
			}
			if c.expectRedirect == "" {
				t.Fatal("expected a failure but got a success")
			}

			loc, err := resp.Location()
			require.NoError(t, err)
			q := loc.Query()

			// Verify the API key is set.
			encryptedAPIKey := loc.Query().Get(workspaceapps.SubdomainProxyAPIKeyParam)
			require.NotEmpty(t, encryptedAPIKey, "no API key was set in the query parameters")

			// Strip the API key from the actual redirect URI and compare.
			q.Del(workspaceapps.SubdomainProxyAPIKeyParam)
			loc.RawQuery = q.Encode()
			require.Equal(t, c.expectRedirect, loc.String())

			var token workspaceapps.EncryptedAPIKeyPayload
			err = jwtutils.Decrypt(ctx, kc, encryptedAPIKey, &token, jwtutils.WithDecryptExpected(jwt.Expected{
				Time:        clock.Now(),
				AnyAudience: jwt.Audience{"wsproxy"},
				Issuer:      "optimus-ide-collabd",
			}))
			require.NoError(t, err)
			require.Equal(t, jwt.NewNumericDate(clock.Now().Add(time.Minute)), token.Expiry)
			require.Equal(t, jwt.NewNumericDate(clock.Now().Add(-time.Minute)), token.NotBefore)
		})
	}
}
