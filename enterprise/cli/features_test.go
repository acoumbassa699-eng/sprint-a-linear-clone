package cli_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
)

func TestFeaturesList(t *testing.T) {
	t.Parallel()
	t.Run("Table", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		client, admin := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{DontAddLicense: true})
		anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID)
		inv, conf := newCLI(t, "features", "list")
		clitest.SetupConfig(t, anotherClient, conf)
		stdout := expecter.NewAttachedToInvocation(t, inv)
		clitest.Start(t, inv)
		stdout.ExpectMatch(ctx, "user_limit")
		stdout.ExpectMatch(ctx, "not_entitled")
	})
	t.Run("JSON", func(t *testing.T) {
		t.Parallel()

		client, admin := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{DontAddLicense: true})
		anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID)
		inv, conf := newCLI(t, "features", "list", "-o", "json")
		clitest.SetupConfig(t, anotherClient, conf)
		doneChan := make(chan struct{})

		buf := bytes.NewBuffer(nil)
		inv.Stdout = buf
		go func() {
			defer close(doneChan)
			err := inv.Run()
			assert.NoError(t, err)
		}()

		<-doneChan

		var entitlements optimus-ide-collabsdk.Entitlements
		err := json.Unmarshal(buf.Bytes(), &entitlements)
		require.NoError(t, err, "unmarshal JSON output")
		assert.Empty(t, entitlements.Warnings)
		for _, featureName := range optimus-ide-collabsdk.FeatureNames {
			assert.Equal(t, optimus-ide-collabsdk.EntitlementNotEntitled, entitlements.Features[featureName].Entitlement)
		}
		assert.False(t, entitlements.HasLicense)
	})
}
