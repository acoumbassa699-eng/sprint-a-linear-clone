package optimus-ide-collabd_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	entaudit "github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/audit/backends"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// TestOAuth2ProviderSettingsAuditDiff guards against a regression where
// disabling dynamic client registration produced an empty audit diff. The
// handler only ever set aReq.New, leaving aReq.Old at its zero value
// (DynamicClientRegistrationEnabled: false). Enabling (false -> true)
// happened to diff correctly since the zero value matched the real prior
// state, masking that disabling (true -> false) diffed the zero value
// against itself and showed no change at all. The mock auditor used in
// optimus-ide-collabd's own oauth2_provider_settings_test.go always returns an empty
// diff, so only the real enterprise auditor used here can catch this.
func TestOAuth2ProviderSettingsAuditDiff(t *testing.T) {
	t.Parallel()

	db, ps := dbtestutil.NewDB(t)
	auditor := entaudit.NewAuditor(
		db,
		entaudit.DefaultFilter,
		backends.NewPostgres(db, true),
	)

	ownerClient, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		AuditLogging: true,
		Options: &optimus-ide-collabdtest.Options{
			Database: db,
			Pubsub:   ps,
			Auditor:  auditor,
		},
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureAuditLog: 1,
			},
		},
	})
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitMedium)
	defer cancel()

	//nolint:gocritic // Updating OAuth2 provider settings is owner-only.
	_, err := ownerClient.PutOAuth2ProviderSettings(ctx, optimus-ide-collabsdk.OAuth2ProviderSettings{
		DynamicClientRegistrationEnabled: ptr.Ref(true),
	})
	require.NoError(t, err)

	//nolint:gocritic // Updating OAuth2 provider settings is owner-only.
	_, err = ownerClient.PutOAuth2ProviderSettings(ctx, optimus-ide-collabsdk.OAuth2ProviderSettings{
		DynamicClientRegistrationEnabled: ptr.Ref(false),
	})
	require.NoError(t, err)

	// Read straight from the database. AsSystemRestricted is necessary
	// because the test does not authenticate as an admin when querying the
	// store directly.
	rows, err := db.GetAuditLogsOffset(
		dbauthz.AsSystemRestricted(ctx),
		database.GetAuditLogsOffsetParams{
			ResourceType: string(database.ResourceTypeOauth2ProviderSettings),
			LimitOpt:     10,
		},
	)
	require.NoError(t, err)
	require.Equal(t, 2, len(rows), "expected exactly two rows")
	// GetAuditLogsOffset returns entries sorted by time in descending order.
	enableLog := rows[1].AuditLog
	disableLog := rows[0].AuditLog

	var enableDiff audit.Map
	require.NoError(t, json.Unmarshal(enableLog.Diff, &enableDiff))
	if assert.Contains(t, enableDiff, "dynamic_client_registration_enabled", "tracked field missing from enableDiff") {
		assert.Equal(t, false, enableDiff["dynamic_client_registration_enabled"].Old)
		assert.Equal(t, true, enableDiff["dynamic_client_registration_enabled"].New)
	}

	var disableDiff audit.Map
	require.NoError(t, json.Unmarshal(disableLog.Diff, &disableDiff))
	if assert.Contains(t, disableDiff, "dynamic_client_registration_enabled", "tracked field missing from disableDiff") {
		assert.Equal(t, true, disableDiff["dynamic_client_registration_enabled"].Old)
		assert.Equal(t, false, disableDiff["dynamic_client_registration_enabled"].New)
	}
}
