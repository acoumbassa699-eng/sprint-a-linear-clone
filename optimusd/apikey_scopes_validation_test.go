package optimus-ide-collabd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestTokenCreation_ScopeValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		scope   optimus-ide-collabsdk.APIKeyScope
		wantErr bool
	}{
		{name: "AllowsPublicLowLevelScope", scope: "workspace:read", wantErr: false},
		{name: "RejectsInternalOnlyScope", scope: "debug_info:read", wantErr: true},
		{name: "AllowsLegacyScopes", scope: "application_connect", wantErr: false},
		{name: "AllowsLegacyScopes2", scope: "all", wantErr: false},
		{name: "AllowsCanonicalSpecialScope", scope: "optimus-ide-collab:all", wantErr: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := optimus-ide-collabdtest.New(t, nil)
			_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

			ctx, cancel := context.WithTimeout(t.Context(), testutil.WaitShort)
			defer cancel()

			resp, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{Scope: tc.scope})
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, resp.Key)

			// Fetch and verify the stored scopes match expectation.
			keys, err := client.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
			require.NoError(t, err)
			require.Len(t, keys, 1)

			// Normalize legacy singular scopes to canonical optimus-ide-collab:* values.
			expected := tc.scope
			switch tc.scope {
			case optimus-ide-collabsdk.APIKeyScopeAll:
				expected = optimus-ide-collabsdk.APIKeyScopeOptimus-IDE-CollabAll
			case optimus-ide-collabsdk.APIKeyScopeApplicationConnect:
				expected = optimus-ide-collabsdk.APIKeyScopeOptimus-IDE-CollabApplicationConnect
			}

			require.Contains(t, keys[0].Scopes, expected)
		})
	}
}
