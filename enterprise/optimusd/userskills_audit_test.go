package optimus-ide-collabd_test

import (
	"encoding/json"
	"fmt"
	"sort"
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

func TestUserSkillAuditDiffTracksContent(t *testing.T) {
	// User skill content is user-authored instruction text, not secret material.
	// The enterprise auditor needs to be used because it writes actual diffs.
	t.Parallel()

	db, ps := dbtestutil.NewDB(t)
	auditor := entaudit.NewAuditor(
		db,
		entaudit.DefaultFilter,
		backends.NewPostgres(db, true),
	)

	ownerClient, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
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
	memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)
	member := optimus-ide-collabsdk.NewExperimentalClient(memberClient)
	ctx := testutil.Context(t, testutil.WaitMedium)

	initialContent := userSkillMarkdown("audit-tracking", "initial", "initial body")
	skill, err := member.CreateUserSkill(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateUserSkillRequest{
		Content: initialContent,
	})
	require.NoError(t, err)

	newContent := userSkillMarkdown("audit-tracking", "after", "new body")
	_, err = member.UpdateUserSkill(ctx, optimus-ide-collabsdk.Me, skill.Name, optimus-ide-collabsdk.UpdateUserSkillRequest{
		Content: newContent,
	})
	require.NoError(t, err)

	rows, err := db.GetAuditLogsOffset(
		dbauthz.AsSystemRestricted(ctx),
		database.GetAuditLogsOffsetParams{
			ResourceType: string(database.ResourceTypeUserSkill),
			LimitOpt:     10,
		},
	)
	require.NoError(t, err)
	require.Len(t, rows, 2, "expected exactly two rows")
	sort.Slice(rows, func(i, j int) bool { return rows[i].AuditLog.Action > rows[j].AuditLog.Action })
	createLog := rows[1].AuditLog
	updateLog := rows[0].AuditLog

	var createDiff audit.Map
	require.NoError(t, json.Unmarshal(createLog.Diff, &createDiff))
	if assert.Contains(t, createDiff, "description", "tracked field missing from create diff") {
		assert.Equal(t, "", createDiff["description"].Old)
		assert.Equal(t, "initial", createDiff["description"].New)
		assert.False(t, createDiff["description"].Secret)
	}
	if assert.Contains(t, createDiff, "content", "content field missing from create diff") {
		assert.False(t, createDiff["content"].Secret)
		assert.Equal(t, "", createDiff["content"].Old)
		assert.Equal(t, initialContent, createDiff["content"].New)
	}

	var updateDiff audit.Map
	require.NoError(t, json.Unmarshal(updateLog.Diff, &updateDiff))
	if assert.Contains(t, updateDiff, "description", "tracked field missing from update diff") {
		assert.Equal(t, "initial", updateDiff["description"].Old)
		assert.Equal(t, "after", updateDiff["description"].New)
		assert.False(t, updateDiff["description"].Secret)
	}
	if assert.Contains(t, updateDiff, "content", "content field missing from update diff") {
		assert.False(t, updateDiff["content"].Secret)
		assert.Equal(t, initialContent, updateDiff["content"].Old)
		assert.Equal(t, newContent, updateDiff["content"].New)
	}
	assert.NotContains(t, updateDiff, "created_at")
	assert.NotContains(t, updateDiff, "updated_at")
}

func userSkillMarkdown(name string, description string, body string) string {
	return fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\n%s\n", name, description, body)
}
