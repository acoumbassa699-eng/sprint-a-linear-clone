package oauth2provider_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/oauth2provider"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestOAuth2AuthorizationServerMetadata(t *testing.T) {
	t.Parallel()

	client := optimus-ide-collabdtest.New(t, nil)
	serverURL := client.URL

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	// Use a plain HTTP client since this endpoint doesn't require authentication.
	// Add a short readiness wait to avoid rare races with server startup.
	endpoint := serverURL.ResolveReference(&url.URL{Path: "/.well-known/oauth-authorization-server"}).String()
	var metadata optimus-ide-collabsdk.OAuth2AuthorizationServerMetadata
	testutil.RequireEventuallyResponseOK(ctx, t, endpoint, &metadata)

	// Verify the metadata
	require.NotEmpty(t, metadata.Issuer)
	require.NotEmpty(t, metadata.AuthorizationEndpoint)
	require.NotEmpty(t, metadata.TokenEndpoint)
	require.Contains(t, metadata.ResponseTypesSupported, optimus-ide-collabsdk.OAuth2ProviderResponseTypeCode)
	require.Contains(t, metadata.GrantTypesSupported, optimus-ide-collabsdk.OAuth2ProviderGrantTypeAuthorizationCode)
	require.Contains(t, metadata.GrantTypesSupported, optimus-ide-collabsdk.OAuth2ProviderGrantTypeRefreshToken)
	require.Contains(t, metadata.CodeChallengeMethodsSupported, optimus-ide-collabsdk.OAuth2PKCECodeChallengeMethodS256)
	// Supported scopes are published from the curated catalog
	require.Equal(t, rbac.ExternalScopeNames(), metadata.ScopesSupported)
}

// TestGetAuthorizationServerMetadata_DCREnabled is a focused unit test on
// the discovery handler itself, bypassing the full optimus-ide-collabdtest HTTP server.
// It verifies the dynamic-client-registration-enabled gate: registration_endpoint
// is advertised once an admin explicitly enables DCR, omitted when explicitly
// disabled, and omitted by default when the setting has never been
// configured.
func TestGetAuthorizationServerMetadata_DCREnabled(t *testing.T) {
	t.Parallel()

	accessURL, err := url.Parse("https://oauth2-metadata-dcr-test.example.com")
	require.NoError(t, err)

	tests := []struct {
		name string
		// configureDCR is nil for "never configured".
		configureDCR             *bool
		wantRegistrationEndpoint bool
	}{
		{
			name:                     "EnabledAdvertisesRegistrationEndpoint",
			configureDCR:             ptr.Ref(true),
			wantRegistrationEndpoint: true,
		},
		{
			name:                     "DisabledOmitsRegistrationEndpoint",
			configureDCR:             ptr.Ref(false),
			wantRegistrationEndpoint: false,
		},
		{
			name:                     "NeverConfiguredDefaultsToOmitted",
			configureDCR:             nil,
			wantRegistrationEndpoint: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := testutil.Context(t, testutil.WaitLong)

			db, _ := dbtestutil.NewDB(t)
			if tt.configureDCR != nil {
				err := db.UpsertOAuth2DCREnabled(ctx, *tt.configureDCR)
				require.NoError(t, err)
			}

			handler := oauth2provider.GetAuthorizationServerMetadata(db, accessURL)

			r := httptest.NewRequest(http.MethodGet, "/.well-known/oauth-authorization-server", nil).WithContext(ctx)
			rw := httptest.NewRecorder()

			handler.ServeHTTP(rw, r)
			require.Equal(t, http.StatusOK, rw.Code)

			var metadata optimus-ide-collabsdk.OAuth2AuthorizationServerMetadata
			require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &metadata))

			if tt.wantRegistrationEndpoint {
				require.NotEmpty(t, metadata.RegistrationEndpoint)
			} else {
				require.Empty(t, metadata.RegistrationEndpoint)
			}
		})
	}
}

func TestOAuth2ProtectedResourceMetadata(t *testing.T) {
	t.Parallel()

	client := optimus-ide-collabdtest.New(t, nil)
	serverURL := client.URL

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	// Use a plain HTTP client since this endpoint doesn't require authentication.
	// Add a short readiness wait to avoid rare races with server startup.
	endpoint := serverURL.ResolveReference(&url.URL{Path: "/.well-known/oauth-protected-resource"}).String()
	var metadata optimus-ide-collabsdk.OAuth2ProtectedResourceMetadata
	testutil.RequireEventuallyResponseOK(ctx, t, endpoint, &metadata)

	// Verify the metadata
	require.NotEmpty(t, metadata.Resource)
	require.NotEmpty(t, metadata.AuthorizationServers)
	require.Len(t, metadata.AuthorizationServers, 1)
	require.Equal(t, metadata.Resource, metadata.AuthorizationServers[0])
	// RFC 6750 bearer tokens are now supported as fallback methods
	require.Contains(t, metadata.BearerMethodsSupported, "header")
	require.Contains(t, metadata.BearerMethodsSupported, "query")
	// Supported scopes are published from the curated catalog
	require.Equal(t, rbac.ExternalScopeNames(), metadata.ScopesSupported)
}
