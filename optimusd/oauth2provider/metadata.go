package oauth2provider

import (
	"net/http"
	"net/url"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// GetAuthorizationServerMetadata returns an http.HandlerFunc that handles GET /.well-known/oauth-authorization-server
func GetAuthorizationServerMetadata(db database.Store, accessURL *url.URL) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// This is queried on every request rather than cached, for the
		// same reason as the registration endpoint: discovery is not
		// expected to be a hot path, and a flood of requests should be
		// mitigated with rate limiting or firewalling, not a cache.
		//nolint:gocritic // Public discovery endpoint, no authenticated actor to authorize against.
		dcrEnabled, err := db.GetOAuth2DCREnabled(dbauthz.AsSystemOAuth2(ctx))
		if err != nil {
			httpapi.InternalServerError(rw, err)
			return
		}

		metadata := optimus-ide-collabsdk.OAuth2AuthorizationServerMetadata{
			Issuer:                            accessURL.String(),
			AuthorizationEndpoint:             accessURL.JoinPath("/oauth2/authorize").String(),
			TokenEndpoint:                     accessURL.JoinPath("/oauth2/tokens").String(),
			RevocationEndpoint:                accessURL.JoinPath("/oauth2/revoke").String(), // RFC 7009
			ResponseTypesSupported:            []optimus-ide-collabsdk.OAuth2ProviderResponseType{optimus-ide-collabsdk.OAuth2ProviderResponseTypeCode},
			GrantTypesSupported:               []optimus-ide-collabsdk.OAuth2ProviderGrantType{optimus-ide-collabsdk.OAuth2ProviderGrantTypeAuthorizationCode, optimus-ide-collabsdk.OAuth2ProviderGrantTypeRefreshToken},
			CodeChallengeMethodsSupported:     []optimus-ide-collabsdk.OAuth2PKCECodeChallengeMethod{optimus-ide-collabsdk.OAuth2PKCECodeChallengeMethodS256},
			ScopesSupported:                   rbac.ExternalScopeNames(),
			TokenEndpointAuthMethodsSupported: []optimus-ide-collabsdk.OAuth2TokenEndpointAuthMethod{optimus-ide-collabsdk.OAuth2TokenEndpointAuthMethodClientSecretBasic, optimus-ide-collabsdk.OAuth2TokenEndpointAuthMethodClientSecretPost},
		}
		if dcrEnabled {
			metadata.RegistrationEndpoint = accessURL.JoinPath("/oauth2/register").String() // RFC 7591
		}
		httpapi.Write(ctx, rw, http.StatusOK, metadata)
	}
}

// GetProtectedResourceMetadata returns an http.HandlerFunc that handles GET /.well-known/oauth-protected-resource
func GetProtectedResourceMetadata(accessURL *url.URL) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		metadata := optimus-ide-collabsdk.OAuth2ProtectedResourceMetadata{
			Resource:             accessURL.String(),
			AuthorizationServers: []string{accessURL.String()},
			ScopesSupported:      rbac.ExternalScopeNames(),
			// RFC 6750 Bearer Token methods supported as fallback methods in api key middleware
			BearerMethodsSupported: []string{"header", "query"},
		}
		httpapi.Write(ctx, rw, http.StatusOK, metadata)
	}
}
