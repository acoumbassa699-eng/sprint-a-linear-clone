package cli_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestOAuth2ProviderDCR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		command     string
		expectValue bool
		expectMsg   string
	}{
		{
			name:        "Enable",
			command:     "enable",
			expectValue: true,
			expectMsg:   "Dynamic client registration is now enabled.",
		},
		{
			name:        "Disable",
			command:     "disable",
			expectValue: false,
			expectMsg:   "Dynamic client registration is now disabled.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := optimus-ide-collabdtest.New(t, nil)
			_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

			inv, root := clitest.New(t, "oauth2-provider", "dcr", tt.command)
			clitest.SetupConfig(t, client, root)

			var buf bytes.Buffer
			inv.Stderr = &buf
			err := inv.Run()
			require.NoError(t, err)
			assert.Contains(t, buf.String(), tt.expectMsg)

			ctx := testutil.Context(t, testutil.WaitShort)
			settings, err := client.OAuth2ProviderSettings(ctx)
			require.NoError(t, err)
			require.NotNil(t, settings.DynamicClientRegistrationEnabled, "GET must always return a concrete value")
			require.Equal(t, tt.expectValue, *settings.DynamicClientRegistrationEnabled)
		})
	}
}

func TestOAuth2ProviderDCR_RegularUser(t *testing.T) {
	t.Parallel()

	client := optimus-ide-collabdtest.New(t, nil)
	owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
	anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

	inv, root := clitest.New(t, "oauth2-provider", "dcr", "enable")
	clitest.SetupConfig(t, anotherClient, root)

	var buf bytes.Buffer
	inv.Stderr = &buf
	err := inv.Run()
	var sdkError *optimus-ide-collabsdk.Error
	require.Error(t, err)
	require.ErrorAsf(t, err, &sdkError, "error should be of type *optimus-ide-collabsdk.Error")
	assert.Equal(t, http.StatusForbidden, sdkError.StatusCode())
}
