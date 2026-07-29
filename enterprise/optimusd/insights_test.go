package optimus-ide-collabd_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestTemplateInsightsWithTemplateAdminACL(t *testing.T) {
	t.Parallel()

	y, m, d := time.Now().UTC().Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	type test struct {
		interval optimus-ide-collabsdk.InsightsReportInterval
	}

	tests := []test{
		{optimus-ide-collabsdk.InsightsReportIntervalDay},
		{""},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("with interval=%q", tt.interval), func(t *testing.T) {
			t.Parallel()

			client, admin := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
				},
			}})
			templateAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID, rbac.RoleTemplateAdmin())

			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, admin.OrganizationID, nil)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, admin.OrganizationID, version.ID)

			regular, regularUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID)

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
			defer cancel()

			err := templateAdminClient.UpdateTemplateACL(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateACL{
				UserPerms: map[string]optimus-ide-collabsdk.TemplateRole{
					regularUser.ID.String(): optimus-ide-collabsdk.TemplateRoleAdmin,
				},
			})
			require.NoError(t, err)

			_, err = regular.TemplateInsights(ctx, optimus-ide-collabsdk.TemplateInsightsRequest{
				StartTime:   today.AddDate(0, 0, -1),
				EndTime:     today,
				TemplateIDs: []uuid.UUID{template.ID},
			})
			require.NoError(t, err)
		})
	}
}

func TestTemplateInsightsWithRole(t *testing.T) {
	t.Parallel()

	y, m, d := time.Now().UTC().Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	type test struct {
		interval optimus-ide-collabsdk.InsightsReportInterval
		role     rbac.RoleIdentifier
		allowed  bool
	}

	tests := []test{
		{optimus-ide-collabsdk.InsightsReportIntervalDay, rbac.RoleTemplateAdmin(), true},
		{"", rbac.RoleTemplateAdmin(), true},
		{optimus-ide-collabsdk.InsightsReportIntervalDay, rbac.RoleAuditor(), true},
		{"", rbac.RoleAuditor(), true},
		{optimus-ide-collabsdk.InsightsReportIntervalDay, rbac.RoleUserAdmin(), false},
		{"", rbac.RoleUserAdmin(), false},
		{optimus-ide-collabsdk.InsightsReportIntervalDay, rbac.RoleMember(), false},
		{"", rbac.RoleMember(), false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("with interval=%q role=%q", tt.interval, tt.role), func(t *testing.T) {
			t.Parallel()

			client, admin := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureTemplateRBAC: 1,
				},
			}})
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, admin.OrganizationID, nil)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, admin.OrganizationID, version.ID)

			aud, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID, tt.role)

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
			defer cancel()

			_, err := aud.TemplateInsights(ctx, optimus-ide-collabsdk.TemplateInsightsRequest{
				StartTime:   today.AddDate(0, 0, -1),
				EndTime:     today,
				TemplateIDs: []uuid.UUID{template.ID},
			})
			if tt.allowed {
				require.NoError(t, err)
			} else {
				var sdkErr *optimus-ide-collabsdk.Error
				require.ErrorAs(t, err, &sdkErr)
				require.Equal(t, sdkErr.StatusCode(), http.StatusNotFound)
			}
		})
	}
}
