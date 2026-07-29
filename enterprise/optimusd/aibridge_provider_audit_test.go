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
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	entaudit "github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/audit/backends"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// TestAIProviderAuditDiff exercises the full HTTP -> enterprise auditor
// -> Postgres write path for AI provider updates. The mock auditor used
// elsewhere returns an empty diff, so this is the only place that
// proves changed properties land in the audit_logs row.
func TestAIProviderAuditDiff(t *testing.T) {
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

	//nolint:gocritic // Owner role is the audience for this endpoint.
	provider, err := ownerClient.CreateAIProvider(ctx, optimus-ide-collabsdk.CreateAIProviderRequest{
		Type:        optimus-ide-collabsdk.AIProviderTypeOpenAI,
		Name:        "audit-target",
		DisplayName: "Audit Target",
		Enabled:     true,
		BaseURL:     "https://api.openai.com/v1",
	})
	require.NoError(t, err)

	newDisplay := "Renamed"
	newURL := "https://api.openai.com/v2"
	disabled := false
	_, err = ownerClient.UpdateAIProvider(ctx, provider.Name, optimus-ide-collabsdk.UpdateAIProviderRequest{
		DisplayName: &newDisplay,
		BaseURL:     &newURL,
		Enabled:     &disabled,
	})
	require.NoError(t, err)

	rows, err := db.GetAuditLogsOffset(
		dbauthz.AsSystemRestricted(ctx),
		database.GetAuditLogsOffsetParams{
			ResourceType: string(database.ResourceTypeAIProvider),
			LimitOpt:     10,
		},
	)
	require.NoError(t, err)
	require.Len(t, rows, 2, "expected one create and one update audit row")

	// GetAuditLogsOffset returns entries sorted by time in descending order.
	updateLog := rows[0].AuditLog
	require.Equal(t, database.AuditActionWrite, updateLog.Action)

	var updateDiff audit.Map
	require.NoError(t, json.Unmarshal(updateLog.Diff, &updateDiff))

	if assert.Contains(t, updateDiff, "display_name", "display_name missing from diff") {
		assert.Equal(t, "Audit Target", updateDiff["display_name"].Old)
		assert.Equal(t, newDisplay, updateDiff["display_name"].New)
	}
	if assert.Contains(t, updateDiff, "base_url", "base_url missing from diff") {
		assert.Equal(t, "https://api.openai.com/v1", updateDiff["base_url"].Old)
		assert.Equal(t, newURL, updateDiff["base_url"].New)
	}
	if assert.Contains(t, updateDiff, "enabled", "enabled missing from diff") {
		assert.Equal(t, true, updateDiff["enabled"].Old)
		assert.Equal(t, false, updateDiff["enabled"].New)
	}
}
