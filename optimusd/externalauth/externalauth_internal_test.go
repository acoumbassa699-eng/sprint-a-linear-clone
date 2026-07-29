package externalauth

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/promoauth"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestGitlabDefaults(t *testing.T) {
	t.Parallel()

	// The default cloud setup. Copying this here as hard coded
	// values.
	cloud := func() optimus-ide-collabsdk.ExternalAuthConfig {
		return optimus-ide-collabsdk.ExternalAuthConfig{
			Type:                          string(optimus-ide-collabsdk.EnhancedExternalAuthProviderGitLab),
			ID:                            string(optimus-ide-collabsdk.EnhancedExternalAuthProviderGitLab),
			AuthURL:                       "https://gitlab.com/oauth/authorize",
			TokenURL:                      "https://gitlab.com/oauth/token",
			ValidateURL:                   "https://gitlab.com/oauth/token/info",
			RevokeURL:                     "https://gitlab.com/oauth/revoke",
			DisplayName:                   "GitLab",
			DisplayIcon:                   "/icon/gitlab.svg",
			Regex:                         `^(https?://)?gitlab\.com(/.*)?$`,
			APIBaseURL:                    "https://gitlab.com/api/v4",
			Scopes:                        []string{"write_repository", "read_api"},
			CodeChallengeMethodsSupported: []string{string(promoauth.PKCEChallengeMethodSha256)},
		}
	}

	tests := []struct {
		name           string
		input          optimus-ide-collabsdk.ExternalAuthConfig
		expected       optimus-ide-collabsdk.ExternalAuthConfig
		mutateExpected func(*optimus-ide-collabsdk.ExternalAuthConfig)
	}{
		// Cloud
		{
			name: "OnlyType",
			input: optimus-ide-collabsdk.ExternalAuthConfig{
				Type: string(optimus-ide-collabsdk.EnhancedExternalAuthProviderGitLab),
			},
			expected: cloud(),
		},
		{
			// If someone was to manually configure the gitlab cli.
			name: "CloudByConfig",
			input: optimus-ide-collabsdk.ExternalAuthConfig{
				Type:    string(optimus-ide-collabsdk.EnhancedExternalAuthProviderGitLab),
				AuthURL: "https://gitlab.com/oauth/authorize",
			},
			expected: cloud(),
		},
		{
			// Changing some of the defaults of the cloud option
			name: "CloudWithChanges",
			input: optimus-ide-collabsdk.ExternalAuthConfig{
				Type: string(optimus-ide-collabsdk.EnhancedExternalAuthProviderGitLab),
				// Adding an extra query param intentionally to break simple
				// string comparisons.
				AuthURL:     "https://gitlab.com/oauth/authorize?foo=bar",
				DisplayName: "custom",
				Regex:       ".*",
			},
			expected: cloud(),
			mutateExpected: func(config *optimus-ide-collabsdk.ExternalAuthConfig) {
				config.AuthURL = "https://gitlab.com/oauth/authorize?foo=bar"
				config.DisplayName = "custom"
				config.Regex = ".*"
			},
		},
		// Self-hosted
		{
			// Dynamically figures out the Validate, Token, and Regex fields.
			name: "SelfHostedOnlyAuthURL",
			input: optimus-ide-collabsdk.ExternalAuthConfig{
				Type:    string(optimus-ide-collabsdk.EnhancedExternalAuthProviderGitLab),
				AuthURL: "https://gitlab.company.org/oauth/authorize?foo=bar",
			},
			expected: cloud(),
			mutateExpected: func(config *optimus-ide-collabsdk.ExternalAuthConfig) {
				config.AuthURL = "https://gitlab.company.org/oauth/authorize?foo=bar"
				config.ValidateURL = "https://gitlab.company.org/oauth/token/info"
				config.TokenURL = "https://gitlab.company.org/oauth/token"
				config.RevokeURL = "https://gitlab.company.org/oauth/revoke"
				config.Regex = `^(https?://)?gitlab\.company\.org(/.*)?$`
				config.APIBaseURL = "https://gitlab.company.org/api/v4"
			},
		},
		{
			// Strange values
			name: "RandomValues",
			input: optimus-ide-collabsdk.ExternalAuthConfig{
				Type:                          string(optimus-ide-collabsdk.EnhancedExternalAuthProviderGitLab),
				AuthURL:                       "https://auth.com/auth",
				ValidateURL:                   "https://validate.com/validate",
				TokenURL:                      "https://token.com/token",
				RevokeURL:                     "https://token.com/revoke",
				Regex:                         "random",
				CodeChallengeMethodsSupported: []string{"random"},
			},
			expected: cloud(),
			mutateExpected: func(config *optimus-ide-collabsdk.ExternalAuthConfig) {
				config.AuthURL = "https://auth.com/auth"
				config.ValidateURL = "https://validate.com/validate"
				config.TokenURL = "https://token.com/token"
				config.RevokeURL = "https://token.com/revoke"
				config.Regex = `random`
				config.CodeChallengeMethodsSupported = []string{"random"}
				config.APIBaseURL = "https://auth.com/api/v4"
			},
		},
	}
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			applyDefaultsToConfig(&c.input)
			if c.mutateExpected != nil {
				c.mutateExpected(&c.expected)
			}
			require.Equal(t, c.input, c.expected)
		})
	}
}

func TestIsFailedRefresh(t *testing.T) {
	t.Parallel()

	expiredToken := &oauth2.Token{
		RefreshToken: "refresh-token",
		// isFailedRefresh returns early at the existingToken.Valid()
		// guard if the token is valid. Valid() requires
		// AccessToken != "" AND not expired. This fixture has no
		// AccessToken so Valid() is always false, but we set an
		// expired time as a safety net in case someone later adds
		// an AccessToken field.
		Expiry: time.Now().Add(-time.Hour),
	}

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name: "IncorrectClientCredentials_StatusOK",
			err: &oauth2.RetrieveError{
				Response:  &http.Response{StatusCode: http.StatusOK},
				ErrorCode: "incorrect_client_credentials",
			},
			// StatusOK fallthrough also returns true, so this test
			// documents the combined behavior. See the 403-status
			// variant below for error-code-only isolation.
			expected: true,
		},
		{
			// Uses 403 status (excluded from the status code switch)
			// so the only path to true is the error code switch.
			name: "IncorrectClientCredentials_Status403",
			err: &oauth2.RetrieveError{
				Response:  &http.Response{StatusCode: http.StatusForbidden},
				ErrorCode: "incorrect_client_credentials",
			},
			expected: true,
		},
		{
			name: "InvalidClient_Status401",
			err: &oauth2.RetrieveError{
				Response:  &http.Response{StatusCode: http.StatusUnauthorized},
				ErrorCode: "invalid_client",
			},
			// StatusUnauthorized fallthrough also returns true, so
			// this test documents the combined behavior.
			expected: true,
		},
		{
			// Uses 403 status (excluded from the status code switch)
			// so the only path to true is the error code switch.
			name: "InvalidClient_Status403",
			err: &oauth2.RetrieveError{
				Response:  &http.Response{StatusCode: http.StatusForbidden},
				ErrorCode: "invalid_client",
			},
			expected: true,
		},
		{
			name: "UnknownErrorCode_Status403_Transient",
			err: &oauth2.RetrieveError{
				Response:  &http.Response{StatusCode: http.StatusForbidden},
				ErrorCode: "unknown_code",
			},
			// 403 with unknown error code should be transient (safe
			// default: retry rather than destroy the token).
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isFailedRefresh(expiredToken, tt.err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func Test_bitbucketServerConfigDefaults(t *testing.T) {
	t.Parallel()

	bbType := string(optimus-ide-collabsdk.EnhancedExternalAuthProviderBitBucketServer)
	tests := []struct {
		name     string
		config   *optimus-ide-collabsdk.ExternalAuthConfig
		expected optimus-ide-collabsdk.ExternalAuthConfig
	}{
		{
			// Very few fields are statically defined for Bitbucket Server.
			name: "EmptyBitbucketServer",
			config: &optimus-ide-collabsdk.ExternalAuthConfig{
				Type: bbType,
			},
			expected: optimus-ide-collabsdk.ExternalAuthConfig{
				Type:                          bbType,
				ID:                            bbType,
				DisplayName:                   "Bitbucket Server",
				Scopes:                        []string{"PUBLIC_REPOS", "REPO_READ", "REPO_WRITE"},
				DisplayIcon:                   "/icon/bitbucket.svg",
				CodeChallengeMethodsSupported: []string{string(promoauth.PKCEChallengeMethodNone)},
			},
		},
		{
			// Only the AuthURL is required for defaults to work.
			name: "AuthURL",
			config: &optimus-ide-collabsdk.ExternalAuthConfig{
				Type:    bbType,
				AuthURL: "https://bitbucket.example.com/login/oauth/authorize",
			},
			expected: optimus-ide-collabsdk.ExternalAuthConfig{
				Type:                          bbType,
				ID:                            bbType,
				AuthURL:                       "https://bitbucket.example.com/login/oauth/authorize",
				TokenURL:                      "https://bitbucket.example.com/rest/oauth2/latest/token",
				ValidateURL:                   "https://bitbucket.example.com/rest/api/latest/inbox/pull-requests/count",
				Scopes:                        []string{"PUBLIC_REPOS", "REPO_READ", "REPO_WRITE"},
				Regex:                         `^(https?://)?bitbucket\.example\.com(/.*)?$`,
				DisplayName:                   "Bitbucket Server",
				DisplayIcon:                   "/icon/bitbucket.svg",
				CodeChallengeMethodsSupported: []string{string(promoauth.PKCEChallengeMethodNone)},
			},
		},
		{
			// Ensure backwards compatibility. The type should update to "bitbucket-cloud",
			// but the ID and other fields should remain the same.
			name: "BitbucketLegacy",
			config: &optimus-ide-collabsdk.ExternalAuthConfig{
				Type: "bitbucket",
			},
			expected: optimus-ide-collabsdk.ExternalAuthConfig{
				Type:                          string(optimus-ide-collabsdk.EnhancedExternalAuthProviderBitBucketCloud),
				ID:                            "bitbucket", // Legacy ID remains unchanged
				AuthURL:                       "https://bitbucket.org/site/oauth2/authorize",
				TokenURL:                      "https://bitbucket.org/site/oauth2/access_token",
				ValidateURL:                   "https://api.bitbucket.org/2.0/user",
				DisplayName:                   "BitBucket",
				DisplayIcon:                   "/icon/bitbucket.svg",
				Regex:                         `^(https?://)?bitbucket\.org(/.*)?$`,
				Scopes:                        []string{"account", "repository:write"},
				CodeChallengeMethodsSupported: []string{string(promoauth.PKCEChallengeMethodNone)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			applyDefaultsToConfig(tt.config)
			require.Equal(t, tt.expected, *tt.config)
		})
	}
}

func TestUntyped(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    optimus-ide-collabsdk.ExternalAuthConfig
		expected optimus-ide-collabsdk.ExternalAuthConfig
	}{
		{
			// Unknown Type uses S256 by default.
			name: "RandomValues",
			input: optimus-ide-collabsdk.ExternalAuthConfig{
				Type:        "unknown",
				AuthURL:     "https://auth.com/auth",
				ValidateURL: "https://validate.com/validate",
				TokenURL:    "https://token.com/token",
				RevokeURL:   "https://token.com/revoke",
				Regex:       "random",
			},
			expected: optimus-ide-collabsdk.ExternalAuthConfig{
				ID:                            "unknown",
				Type:                          "unknown",
				DisplayName:                   "unknown",
				DisplayIcon:                   "/emojis/1f511.png",
				AuthURL:                       "https://auth.com/auth",
				ValidateURL:                   "https://validate.com/validate",
				TokenURL:                      "https://token.com/token",
				RevokeURL:                     "https://token.com/revoke",
				Regex:                         `random`,
				CodeChallengeMethodsSupported: []string{string(promoauth.PKCEChallengeMethodSha256)},
			},
		},
	}
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			applyDefaultsToConfig(&c.input)
			require.Equal(t, c.input, c.expected)
		})
	}
}
