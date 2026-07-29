package optimus-ide-collabd_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestOAuth2ProviderSettings(t *testing.T) {
	t.Parallel()

	t.Run("DefaultDisabled", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
		ctx := testutil.Context(t, testutil.WaitShort)

		settings, err := client.OAuth2ProviderSettings(ctx)
		require.NoError(t, err)
		require.NotNil(t, settings.DynamicClientRegistrationEnabled, "GET must always return a concrete value")
		require.False(t, *settings.DynamicClientRegistrationEnabled)
	})

	t.Run("RoundTrip", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
		ctx := testutil.Context(t, testutil.WaitShort)

		updated, err := client.PutOAuth2ProviderSettings(ctx, optimus-ide-collabsdk.OAuth2ProviderSettings{
			DynamicClientRegistrationEnabled: ptr.Ref(false),
		})
		require.NoError(t, err)
		require.NotNil(t, updated.DynamicClientRegistrationEnabled)
		require.False(t, *updated.DynamicClientRegistrationEnabled)

		settings, err := client.OAuth2ProviderSettings(ctx)
		require.NoError(t, err)
		require.NotNil(t, settings.DynamicClientRegistrationEnabled)
		require.False(t, *settings.DynamicClientRegistrationEnabled)

		updated, err = client.PutOAuth2ProviderSettings(ctx, optimus-ide-collabsdk.OAuth2ProviderSettings{
			DynamicClientRegistrationEnabled: ptr.Ref(true),
		})
		require.NoError(t, err)
		require.NotNil(t, updated.DynamicClientRegistrationEnabled)
		require.True(t, *updated.DynamicClientRegistrationEnabled)
	})

	t.Run("OmittedFieldLeavesValueUnchanged", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
		ctx := testutil.Context(t, testutil.WaitShort)

		// Establish a known, non-default value first so a subsequent
		// omitted-field PUT has something to preserve or wrongly clear.
		_, err := client.PutOAuth2ProviderSettings(ctx, optimus-ide-collabsdk.OAuth2ProviderSettings{
			DynamicClientRegistrationEnabled: ptr.Ref(true),
		})
		require.NoError(t, err)

		// A PUT with the field omitted (nil) must leave the current value
		// alone rather than decode to false and silently disable it. This
		// guards the fix for https://github.com/optimus-ide-collab/optimus-ide-collab/pull/27316#issuecomment-5086163278:
		// once a second field lands in this struct, an older client that only
		// knows about this field would otherwise always send its zero value.
		updated, err := client.PutOAuth2ProviderSettings(ctx, optimus-ide-collabsdk.OAuth2ProviderSettings{
			DynamicClientRegistrationEnabled: nil,
		})
		require.NoError(t, err)
		require.NotNil(t, updated.DynamicClientRegistrationEnabled, "response must always reflect the resolved value, never echo back nil")
		require.True(t, *updated.DynamicClientRegistrationEnabled, "omitted field must not have reset the value to false")

		settings, err := client.OAuth2ProviderSettings(ctx)
		require.NoError(t, err)
		require.NotNil(t, settings.DynamicClientRegistrationEnabled)
		require.True(t, *settings.DynamicClientRegistrationEnabled, "value must remain unchanged in the database")
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			do   func(ctx context.Context, client *optimus-ide-collabsdk.Client) error
		}{
			{
				name: "Get",
				do: func(ctx context.Context, client *optimus-ide-collabsdk.Client) error {
					_, err := client.OAuth2ProviderSettings(ctx)
					return err
				},
			},
			{
				name: "Put",
				do: func(ctx context.Context, client *optimus-ide-collabsdk.Client) error {
					_, err := client.PutOAuth2ProviderSettings(ctx, optimus-ide-collabsdk.OAuth2ProviderSettings{
						DynamicClientRegistrationEnabled: ptr.Ref(false),
					})
					return err
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				client := optimus-ide-collabdtest.New(t, nil)
				firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client)
				anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, firstUser.OrganizationID)
				ctx := testutil.Context(t, testutil.WaitShort)

				err := tt.do(ctx, anotherClient)
				var sdkError *optimus-ide-collabsdk.Error
				require.Error(t, err)
				require.ErrorAsf(t, err, &sdkError, "error should be of type *optimus-ide-collabsdk.Error")
				require.Equal(t, http.StatusForbidden, sdkError.StatusCode())
			})
		}
	})
}
