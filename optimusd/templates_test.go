package optimus-ide-collabd_test

import (
	"context"
	"database/sql"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/agent/agenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications/notificationstest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/schedule"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/workspacesdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// TestTemplatesListSingleAuthorizePrepare guards against reintroducing the
// redundant OPA partial evaluation the GET /api/v2/templates handler used to
// perform. The handler called AuthorizeSQLFilter to build a prepared
// ResourceTemplate authorizer, but the dbauthz GetAuthorizedTemplates wrapper
// ignored it and re-prepared inside GetTemplatesWithFilter, so every request
// ran partial evaluation twice. Partial-evaluation cost scales with the
// number of organization-scoped roles the subject carries (see #21890), so the
// duplicate prepare doubled an already expensive operation. A single list
// request must prepare the ResourceTemplate authorizer exactly once.
func TestTemplatesListSingleAuthorizePrepare(t *testing.T) {
	t.Parallel()

	authz := &optimus-ide-collabdtest.RecordingAuthorizer{Wrapped: rbac.NewStrictCachingAuthorizer(prometheus.NewRegistry())}
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		IncludeProvisionerDaemon: true,
		Authorizer:               authz,
	})
	owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

	ctx := testutil.Context(t, testutil.WaitLong)

	// Reset immediately before the measured request so setup prepares (template
	// creation, version jobs) are excluded. Counts are keyed by subject ID, so
	// background work under system subjects is ignored.
	authz.Reset()
	templates, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{})
	require.NoError(t, err)
	require.Len(t, templates, 1)

	count := authz.PrepareCount(owner.UserID.String(), policy.ActionRead, rbac.ResourceTemplate.Type)
	require.Equal(t, 1, count,
		"GET /templates must prepare the ResourceTemplate authorizer exactly once; a higher count means a redundant partial evaluation was reintroduced")
}

func TestTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)

		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.Template(ctx, template.ID)
		require.NoError(t, err)
	})
}

func TestPostTemplateByOrganization(t *testing.T) {
	t.Parallel()
	t.Run("Create", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		ownerClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true, Auditor: auditor})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)

		// Use org scoped template admin
		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.ScopedRoleOrgTemplateAdmin(owner.OrganizationID))
		// By default, everyone in the org can read the template.
		user, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)
		auditor.ResetLogs()

		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)

		expected := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.ActivityBumpMillis = ptr.Ref((3 * time.Hour).Milliseconds())
			ctr.TimeTilAutostopNotifyMillis = ptr.Ref((5 * time.Minute).Milliseconds())
		})
		assert.Equal(t, (3 * time.Hour).Milliseconds(), expected.ActivityBumpMillis)
		assert.Equal(t, (5 * time.Minute).Milliseconds(), expected.TimeTilAutostopNotifyMillis)

		ctx := testutil.Context(t, testutil.WaitLong)

		got, err := user.Template(ctx, expected.ID)
		require.NoError(t, err)

		assert.Equal(t, expected.Name, got.Name)
		assert.Equal(t, expected.Description, got.Description)
		assert.Equal(t, expected.ActivityBumpMillis, got.ActivityBumpMillis)
		assert.Equal(t, expected.TimeTilAutostopNotifyMillis, got.TimeTilAutostopNotifyMillis)
		assert.Equal(t, expected.UseClassicParameterFlow, false) // Current default is false

		require.Len(t, auditor.AuditLogs(), 3)
		assert.Equal(t, database.AuditActionCreate, auditor.AuditLogs()[0].Action)
		assert.Equal(t, database.AuditActionWrite, auditor.AuditLogs()[1].Action)
		assert.Equal(t, database.AuditActionCreate, auditor.AuditLogs()[2].Action)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		t.Parallel()
		ownerClient := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)
		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.ScopedRoleOrgTemplateAdmin(owner.OrganizationID))

		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.CreateTemplate(ctx, owner.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
			Name:      template.Name,
			VersionID: version.ID,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusConflict, apiErr.StatusCode())
	})

	t.Run("ReservedName", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

		ctx := testutil.Context(t, testutil.WaitShort)

		_, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
			Name:      "new",
			VersionID: version.ID,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})

	t.Run("DefaultTTLTooLow", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

		ctx := testutil.Context(t, testutil.WaitLong)
		_, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
			Name:             "testing",
			VersionID:        version.ID,
			DefaultTTLMillis: ptr.Ref(int64(-1)),
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
		require.Contains(t, err.Error(), "default_ttl_ms: Must be a positive integer")
	})

	t.Run("TimeTilAutostopNotifyTooLow", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

		ctx := testutil.Context(t, testutil.WaitLong)
		_, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
			Name:                        "testing",
			VersionID:                   version.ID,
			TimeTilAutostopNotifyMillis: ptr.Ref(int64(30_000)),
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
		require.Len(t, apiErr.Validations, 1)
		assert.Equal(t, "time_til_autostop_notify_ms", apiErr.Validations[0].Field)
		assert.Equal(t, "Must be 0 (disabled) or at least one minute.", apiErr.Validations[0].Detail)
	})

	t.Run("TimeTilAutostopNotifyExactlyOneMinute", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

		ctx := testutil.Context(t, testutil.WaitLong)
		got, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
			Name:                        "testing",
			VersionID:                   version.ID,
			TimeTilAutostopNotifyMillis: ptr.Ref(time.Minute.Milliseconds()),
		})
		require.NoError(t, err)
		assert.Equal(t, time.Minute.Milliseconds(), got.TimeTilAutostopNotifyMillis)
	})

	t.Run("TimeTilAutostopNotifyNegative", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

		ctx := testutil.Context(t, testutil.WaitLong)
		_, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
			Name:                        "testing",
			VersionID:                   version.ID,
			TimeTilAutostopNotifyMillis: ptr.Ref(int64(-1)),
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
		require.Len(t, apiErr.Validations, 1)
		assert.Equal(t, "time_til_autostop_notify_ms", apiErr.Validations[0].Field)
		assert.Equal(t, "Must be a positive integer.", apiErr.Validations[0].Detail)
	})

	t.Run("NoDefaultTTL", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

		ctx := testutil.Context(t, testutil.WaitLong)
		got, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
			Name:             "testing",
			VersionID:        version.ID,
			DefaultTTLMillis: ptr.Ref(int64(0)),
		})
		require.NoError(t, err)
		require.Zero(t, got.DefaultTTLMillis)
	})

	t.Run("DisableEveryone", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true, Auditor: auditor})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		user, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		expected := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.DisableEveryoneGroupAccess = true
		})

		ctx := testutil.Context(t, testutil.WaitLong)
		_, err := user.Template(ctx, expected.ID)

		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	t.Run("Unauthorized", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)

		ctx := testutil.Context(t, testutil.WaitLong)
		_, err := client.CreateTemplate(ctx, uuid.New(), optimus-ide-collabsdk.CreateTemplateRequest{
			Name:      "test",
			VersionID: uuid.New(),
		})

		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode())
		require.Contains(t, err.Error(), "Try logging in using 'optimus-ide-collab login'.")
	})

	t.Run("AllowUserScheduling", func(t *testing.T) {
		t.Parallel()

		t.Run("OK", func(t *testing.T) {
			t.Parallel()

			var setCalled atomic.Int64
			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
				TemplateScheduleStore: schedule.MockTemplateScheduleStore{
					SetFn: func(ctx context.Context, db database.Store, template database.Template, options schedule.TemplateScheduleOptions) (database.Template, error) {
						setCalled.Add(1)
						require.False(t, options.UserAutostartEnabled)
						require.False(t, options.UserAutostopEnabled)
						template.AllowUserAutostart = options.UserAutostartEnabled
						template.AllowUserAutostop = options.UserAutostopEnabled
						return template, nil
					},
				},
			})
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			got, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
				Name:               "testing",
				VersionID:          version.ID,
				AllowUserAutostart: ptr.Ref(false),
				AllowUserAutostop:  ptr.Ref(false),
			})
			require.NoError(t, err)

			require.EqualValues(t, 1, setCalled.Load())
			require.False(t, got.AllowUserAutostart)
			require.False(t, got.AllowUserAutostop)
		})

		t.Run("IgnoredUnlicensed", func(t *testing.T) {
			t.Parallel()

			client := optimus-ide-collabdtest.New(t, nil)
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			got, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
				Name:               "testing",
				VersionID:          version.ID,
				AllowUserAutostart: ptr.Ref(false),
				AllowUserAutostop:  ptr.Ref(false),
			})
			require.NoError(t, err)
			// ignored and use AGPL defaults
			require.True(t, got.AllowUserAutostart)
			require.True(t, got.AllowUserAutostop)
		})
	})

	t.Run("NoVersion", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
			Name:      "test",
			VersionID: uuid.New(),
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	t.Run("AutostopRequirement", func(t *testing.T) {
		t.Parallel()

		t.Run("None", func(t *testing.T) {
			t.Parallel()

			var setCalled atomic.Int64
			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
				TemplateScheduleStore: schedule.MockTemplateScheduleStore{
					SetFn: func(ctx context.Context, db database.Store, template database.Template, options schedule.TemplateScheduleOptions) (database.Template, error) {
						setCalled.Add(1)
						assert.Zero(t, options.AutostopRequirement.DaysOfWeek)
						assert.Zero(t, options.AutostopRequirement.Weeks)

						err := db.UpdateTemplateScheduleByID(ctx, database.UpdateTemplateScheduleByIDParams{
							ID:                            template.ID,
							UpdatedAt:                     dbtime.Now(),
							AllowUserAutostart:            options.UserAutostartEnabled,
							AllowUserAutostop:             options.UserAutostopEnabled,
							DefaultTTL:                    int64(options.DefaultTTL),
							ActivityBump:                  int64(options.ActivityBump),
							AutostopRequirementDaysOfWeek: int16(options.AutostopRequirement.DaysOfWeek),
							AutostopRequirementWeeks:      options.AutostopRequirement.Weeks,
							FailureTTL:                    int64(options.FailureTTL),
							TimeTilDormant:                int64(options.TimeTilDormant),
							TimeTilDormantAutoDelete:      int64(options.TimeTilDormantAutoDelete),
						})
						if !assert.NoError(t, err) {
							return database.Template{}, err
						}

						return db.GetTemplateByID(ctx, template.ID)
					},
				},
			})
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			got, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
				Name:                "testing",
				VersionID:           version.ID,
				AutostopRequirement: nil,
			})
			require.NoError(t, err)

			require.EqualValues(t, 1, setCalled.Load())
			require.Empty(t, got.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 1, got.AutostopRequirement.Weeks)
		})

		t.Run("OK", func(t *testing.T) {
			t.Parallel()

			var setCalled atomic.Int64
			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
				TemplateScheduleStore: schedule.MockTemplateScheduleStore{
					SetFn: func(ctx context.Context, db database.Store, template database.Template, options schedule.TemplateScheduleOptions) (database.Template, error) {
						setCalled.Add(1)
						assert.EqualValues(t, 0b00110000, options.AutostopRequirement.DaysOfWeek)
						assert.EqualValues(t, 2, options.AutostopRequirement.Weeks)

						err := db.UpdateTemplateScheduleByID(ctx, database.UpdateTemplateScheduleByIDParams{
							ID:                            template.ID,
							UpdatedAt:                     dbtime.Now(),
							AllowUserAutostart:            options.UserAutostartEnabled,
							AllowUserAutostop:             options.UserAutostopEnabled,
							DefaultTTL:                    int64(options.DefaultTTL),
							ActivityBump:                  int64(options.ActivityBump),
							AutostopRequirementDaysOfWeek: int16(options.AutostopRequirement.DaysOfWeek),
							AutostopRequirementWeeks:      options.AutostopRequirement.Weeks,
							FailureTTL:                    int64(options.FailureTTL),
							TimeTilDormant:                int64(options.TimeTilDormant),
							TimeTilDormantAutoDelete:      int64(options.TimeTilDormantAutoDelete),
						})
						if !assert.NoError(t, err) {
							return database.Template{}, err
						}

						return db.GetTemplateByID(ctx, template.ID)
					},
				},
			})
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			got, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
				Name:      "testing",
				VersionID: version.ID,
				AutostopRequirement: &optimus-ide-collabsdk.TemplateAutostopRequirement{
					// wrong order
					DaysOfWeek: []string{"saturday", "friday"},
					Weeks:      2,
				},
			})
			require.NoError(t, err)

			require.EqualValues(t, 1, setCalled.Load())
			require.Equal(t, []string{"friday", "saturday"}, got.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 2, got.AutostopRequirement.Weeks)

			got, err = client.Template(ctx, got.ID)
			require.NoError(t, err)
			require.Equal(t, []string{"friday", "saturday"}, got.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 2, got.AutostopRequirement.Weeks)
		})

		t.Run("IgnoredUnlicensed", func(t *testing.T) {
			t.Parallel()

			client := optimus-ide-collabdtest.New(t, nil)
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			got, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
				Name:      "testing",
				VersionID: version.ID,
				AutostopRequirement: &optimus-ide-collabsdk.TemplateAutostopRequirement{
					DaysOfWeek: []string{"friday", "saturday"},
					Weeks:      2,
				},
			})
			require.NoError(t, err)
			// ignored and use AGPL defaults
			require.Empty(t, got.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 1, got.AutostopRequirement.Weeks)
		})
	})

	t.Run("MaxPortShareLevel", func(t *testing.T) {
		t.Parallel()

		t.Run("OK", func(t *testing.T) {
			client := optimus-ide-collabdtest.New(t, nil)
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			got, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
				Name:      "testing",
				VersionID: version.ID,
			})
			require.NoError(t, err)
			require.Equal(t, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic, got.MaxPortShareLevel)
		})

		t.Run("EnterpriseLevelError", func(t *testing.T) {
			client := optimus-ide-collabdtest.New(t, nil)
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			_, err := client.CreateTemplate(ctx, user.OrganizationID, optimus-ide-collabsdk.CreateTemplateRequest{
				Name:              "testing",
				VersionID:         version.ID,
				MaxPortShareLevel: ptr.Ref(optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic),
			})
			var apiErr *optimus-ide-collabsdk.Error
			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
		})
	})
}

func TestTemplates(t *testing.T) {
	t.Parallel()

	t.Run("ListEmpty", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx := testutil.Context(t, testutil.WaitLong)

		templates, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{})
		require.NoError(t, err)
		require.NotNil(t, templates)
		require.Len(t, templates, 0)
	})

	// Should return only non-deprecated templates by default
	t.Run("ListMultiple non-deprecated", func(t *testing.T) {
		t.Parallel()

		owner, db := optimus-ide-collabdtest.NewWithDatabase(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: false})
		user := optimus-ide-collabdtest.CreateFirstUser(t, owner)
		client, tplAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, owner, user.OrganizationID, rbac.RoleTemplateAdmin())
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		version2 := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		foo := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "foo"
		})
		bar := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version2.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "bar"
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		// Deprecate bar template
		deprecationMessage := "Some deprecated message"
		err := db.UpdateTemplateAccessControlByID(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubject(tplAdmin)), database.UpdateTemplateAccessControlByIDParams{
			ID:                   bar.ID,
			RequireActiveVersion: false,
			Deprecated:           deprecationMessage,
		})
		require.NoError(t, err)

		updatedBar, err := client.Template(ctx, bar.ID)
		require.NoError(t, err)
		require.True(t, updatedBar.Deprecated)
		require.Equal(t, deprecationMessage, updatedBar.DeprecationMessage)

		// Should return only the non-deprecated template (foo)
		templates, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{})
		require.NoError(t, err)
		require.Len(t, templates, 1)

		require.Equal(t, foo.ID, templates[0].ID)
		require.False(t, templates[0].Deprecated)
		require.Empty(t, templates[0].DeprecationMessage)
	})

	// Should return only deprecated templates when filtering by deprecated:true
	t.Run("ListMultiple deprecated:true", func(t *testing.T) {
		t.Parallel()

		owner, db := optimus-ide-collabdtest.NewWithDatabase(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: false})
		user := optimus-ide-collabdtest.CreateFirstUser(t, owner)
		client, tplAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, owner, user.OrganizationID, rbac.RoleTemplateAdmin())
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		version2 := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		foo := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "foo"
		})
		bar := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version2.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "bar"
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		// Deprecate foo and bar templates
		deprecationMessage := "Some deprecated message"
		err := db.UpdateTemplateAccessControlByID(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubject(tplAdmin)), database.UpdateTemplateAccessControlByIDParams{
			ID:                   foo.ID,
			RequireActiveVersion: false,
			Deprecated:           deprecationMessage,
		})
		require.NoError(t, err)
		err = db.UpdateTemplateAccessControlByID(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubject(tplAdmin)), database.UpdateTemplateAccessControlByIDParams{
			ID:                   bar.ID,
			RequireActiveVersion: false,
			Deprecated:           deprecationMessage,
		})
		require.NoError(t, err)

		// Should have deprecation message set
		updatedFoo, err := client.Template(ctx, foo.ID)
		require.NoError(t, err)
		require.True(t, updatedFoo.Deprecated)
		require.Equal(t, deprecationMessage, updatedFoo.DeprecationMessage)

		updatedBar, err := client.Template(ctx, bar.ID)
		require.NoError(t, err)
		require.True(t, updatedBar.Deprecated)
		require.Equal(t, deprecationMessage, updatedBar.DeprecationMessage)

		// Should return only the deprecated templates (foo and bar)
		templates, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
			SearchQuery: "deprecated:true",
		})
		require.NoError(t, err)
		require.Len(t, templates, 2)

		// Make sure all the deprecated templates are returned
		expectedTemplates := map[uuid.UUID]optimus-ide-collabsdk.Template{
			updatedFoo.ID: updatedFoo,
			updatedBar.ID: updatedBar,
		}
		actualTemplates := map[uuid.UUID]optimus-ide-collabsdk.Template{}
		for _, template := range templates {
			actualTemplates[template.ID] = template
		}

		require.Equal(t, len(expectedTemplates), len(actualTemplates))
		for id, expectedTemplate := range expectedTemplates {
			actualTemplate, ok := actualTemplates[id]
			require.True(t, ok)
			require.Equal(t, expectedTemplate.ID, actualTemplate.ID)
			require.Equal(t, true, actualTemplate.Deprecated)
			require.Equal(t, expectedTemplate.DeprecationMessage, actualTemplate.DeprecationMessage)
		}
	})

	// Should return only non-deprecated templates when filtering by deprecated:false
	t.Run("ListMultiple deprecated:false", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		version2 := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		foo := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "foo"
		})
		bar := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version2.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "bar"
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		// Should return only the non-deprecated templates
		templates, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
			SearchQuery: "deprecated:false",
		})
		require.NoError(t, err)
		require.Len(t, templates, 2)

		// Make sure all the non-deprecated templates are returned
		expectedTemplates := map[uuid.UUID]optimus-ide-collabsdk.Template{
			foo.ID: foo,
			bar.ID: bar,
		}
		actualTemplates := map[uuid.UUID]optimus-ide-collabsdk.Template{}
		for _, template := range templates {
			actualTemplates[template.ID] = template
		}

		require.Equal(t, len(expectedTemplates), len(actualTemplates))
		for id, expectedTemplate := range expectedTemplates {
			actualTemplate, ok := actualTemplates[id]
			require.True(t, ok)
			require.Equal(t, expectedTemplate.ID, actualTemplate.ID)
			require.Equal(t, false, actualTemplate.Deprecated)
			require.Equal(t, expectedTemplate.DeprecationMessage, actualTemplate.DeprecationMessage)
		}
	})

	// Should return a re-enabled template in the default (non-deprecated) list
	t.Run("ListMultiple re-enabled template", func(t *testing.T) {
		t.Parallel()

		owner, db := optimus-ide-collabdtest.NewWithDatabase(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: false})
		user := optimus-ide-collabdtest.CreateFirstUser(t, owner)
		client, tplAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, owner, user.OrganizationID, rbac.RoleTemplateAdmin())
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		version2 := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		foo := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "foo"
		})
		bar := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version2.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "bar"
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		// Deprecate bar template
		deprecationMessage := "Some deprecated message"
		err := db.UpdateTemplateAccessControlByID(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubject(tplAdmin)), database.UpdateTemplateAccessControlByIDParams{
			ID:                   bar.ID,
			RequireActiveVersion: false,
			Deprecated:           deprecationMessage,
		})
		require.NoError(t, err)

		updatedBar, err := client.Template(ctx, bar.ID)
		require.NoError(t, err)
		require.True(t, updatedBar.Deprecated)
		require.Equal(t, deprecationMessage, updatedBar.DeprecationMessage)

		// Re-enable bar template
		err = db.UpdateTemplateAccessControlByID(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubject(tplAdmin)), database.UpdateTemplateAccessControlByIDParams{
			ID:                   bar.ID,
			RequireActiveVersion: false,
			Deprecated:           "",
		})
		require.NoError(t, err)

		reEnabledBar, err := client.Template(ctx, bar.ID)
		require.NoError(t, err)
		require.False(t, reEnabledBar.Deprecated)
		require.Empty(t, reEnabledBar.DeprecationMessage)

		// Should return only the non-deprecated templates (foo and bar)
		templates, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{})
		require.NoError(t, err)
		require.Len(t, templates, 2)

		// Make sure all the non-deprecated templates are returned
		expectedTemplates := map[uuid.UUID]optimus-ide-collabsdk.Template{
			foo.ID: foo,
			bar.ID: bar,
		}
		actualTemplates := map[uuid.UUID]optimus-ide-collabsdk.Template{}
		for _, template := range templates {
			actualTemplates[template.ID] = template
		}

		require.Equal(t, len(expectedTemplates), len(actualTemplates))
		for id, expectedTemplate := range expectedTemplates {
			actualTemplate, ok := actualTemplates[id]
			require.True(t, ok)
			require.Equal(t, expectedTemplate.ID, actualTemplate.ID)
			require.Equal(t, false, actualTemplate.Deprecated)
			require.Equal(t, expectedTemplate.DeprecationMessage, actualTemplate.DeprecationMessage)
		}
	})
}

func TestTemplatesByOrganization(t *testing.T) {
	t.Parallel()
	t.Run("ListEmpty", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx := testutil.Context(t, testutil.WaitLong)

		templates, err := client.TemplatesByOrganization(ctx, user.OrganizationID)
		require.NoError(t, err)
		require.NotNil(t, templates)
		require.Len(t, templates, 0)
	})

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)

		ctx := testutil.Context(t, testutil.WaitLong)

		templates, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
			OrganizationID: user.OrganizationID,
		})
		require.NoError(t, err)
		require.Len(t, templates, 1)
	})
	t.Run("ListMultiple", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		version2 := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		foo := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "foobar"
		})
		bar := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version2.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "barbaz"
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		templates, err := client.TemplatesByOrganization(ctx, user.OrganizationID)
		require.NoError(t, err)
		require.Len(t, templates, 2)

		// Listing all should match
		templates, err = client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{})
		require.NoError(t, err)
		require.Len(t, templates, 2)

		org, err := client.Organization(ctx, user.OrganizationID)
		require.NoError(t, err)
		for _, tmpl := range templates {
			require.Equal(t, tmpl.OrganizationID, user.OrganizationID, "organization ID")
			require.Equal(t, tmpl.OrganizationName, org.Name, "organization name")
			require.Equal(t, tmpl.OrganizationDisplayName, org.DisplayName, "organization display name")
			require.Equal(t, tmpl.OrganizationIcon, org.Icon, "organization display name")
		}

		// Check fuzzy name matching
		templates, err = client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
			FuzzyName: "bar",
		})
		require.NoError(t, err)
		require.Len(t, templates, 2)

		templates, err = client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
			FuzzyName: "foo",
		})
		require.NoError(t, err)
		require.Len(t, templates, 1)
		require.Equal(t, foo.ID, templates[0].ID)

		templates, err = client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
			FuzzyName: "baz",
		})
		require.NoError(t, err)
		require.Len(t, templates, 1)
		require.Equal(t, bar.ID, templates[0].ID)
	})

	// Should return only non-deprecated templates by default
	t.Run("ListMultiple non-deprecated", func(t *testing.T) {
		t.Parallel()

		owner, db := optimus-ide-collabdtest.NewWithDatabase(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: false})
		user := optimus-ide-collabdtest.CreateFirstUser(t, owner)
		client, tplAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, owner, user.OrganizationID, rbac.RoleTemplateAdmin())
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		version2 := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		foo := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "foo"
		})
		bar := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version2.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "bar"
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		// Deprecate bar template
		deprecationMessage := "Some deprecated message"
		err := db.UpdateTemplateAccessControlByID(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubject(tplAdmin)), database.UpdateTemplateAccessControlByIDParams{
			ID:                   bar.ID,
			RequireActiveVersion: false,
			Deprecated:           deprecationMessage,
		})
		require.NoError(t, err)

		updatedBar, err := client.Template(ctx, bar.ID)
		require.NoError(t, err)
		require.True(t, updatedBar.Deprecated)
		require.Equal(t, deprecationMessage, updatedBar.DeprecationMessage)

		// Should return only the non-deprecated template (foo)
		templates, err := client.TemplatesByOrganization(ctx, user.OrganizationID)
		require.NoError(t, err)
		require.Len(t, templates, 1)

		require.Equal(t, foo.ID, templates[0].ID)
		require.False(t, templates[0].Deprecated)
		require.Empty(t, templates[0].DeprecationMessage)
	})

	t.Run("ListByAuthor", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		adminAlpha, adminAlphaData := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())
		adminBravo, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())
		adminCharlie, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())

		versionA := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		versionB := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		versionC := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		foo := optimus-ide-collabdtest.CreateTemplate(t, adminAlpha, owner.OrganizationID, versionA.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "foo"
		})
		bar := optimus-ide-collabdtest.CreateTemplate(t, adminBravo, owner.OrganizationID, versionB.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "bar"
		})
		_ = optimus-ide-collabdtest.CreateTemplate(t, adminCharlie, owner.OrganizationID, versionC.ID, func(request *optimus-ide-collabsdk.CreateTemplateRequest) {
			request.Name = "baz"
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		// List alpha
		alpha, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
			AuthorUsername: adminAlphaData.Username,
		})
		require.NoError(t, err)
		require.Len(t, alpha, 1)
		require.Equal(t, foo.ID, alpha[0].ID)

		// List bravo
		bravo, err := adminBravo.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
			AuthorUsername: optimus-ide-collabsdk.Me,
		})
		require.NoError(t, err)
		require.Len(t, bravo, 1)
		require.Equal(t, bar.ID, bravo[0].ID)
	})
}

func TestTemplateByOrganizationAndName(t *testing.T) {
	t.Parallel()
	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.TemplateByName(ctx, user.OrganizationID, "something")
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	t.Run("Found", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)

		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.TemplateByName(ctx, user.OrganizationID, template.Name)
		require.NoError(t, err)
	})
}

func TestPatchTemplateMeta(t *testing.T) {
	t.Parallel()

	t.Run("Modified", func(t *testing.T) {
		t.Parallel()

		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true, Auditor: auditor})
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		assert.Equal(t, (1 * time.Hour).Milliseconds(), template.ActivityBumpMillis)

		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			Name:                         ptr.Ref("new-template-name"),
			DisplayName:                  ptr.Ref("Displayed Name 456"),
			Description:                  ptr.Ref("lorem ipsum dolor sit amet et cetera"),
			Icon:                         ptr.Ref("/icon/new-icon.png"),
			DefaultTTLMillis:             ptr.Ref(12 * time.Hour.Milliseconds()),
			ActivityBumpMillis:           ptr.Ref(3 * time.Hour.Milliseconds()),
			TimeTilAutostopNotifyMillis:  ptr.Ref(5 * time.Minute.Milliseconds()),
			AllowUserCancelWorkspaceJobs: ptr.Ref(false),
		}
		// It is unfortunate we need to sleep, but the test can fail if the
		// updatedAt is too close together.
		time.Sleep(time.Millisecond * 5)

		ctx := testutil.Context(t, testutil.WaitLong)

		updated, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		assert.Greater(t, updated.UpdatedAt, template.UpdatedAt)
		assert.Equal(t, *req.Name, updated.Name)
		assert.Equal(t, *req.DisplayName, updated.DisplayName)
		assert.Equal(t, *req.Description, updated.Description)
		assert.Equal(t, *req.Icon, updated.Icon)
		assert.Equal(t, *req.DefaultTTLMillis, updated.DefaultTTLMillis)
		assert.Equal(t, *req.ActivityBumpMillis, updated.ActivityBumpMillis)
		assert.Equal(t, *req.TimeTilAutostopNotifyMillis, updated.TimeTilAutostopNotifyMillis)
		assert.False(t, *req.AllowUserCancelWorkspaceJobs)

		// Extra paranoid: did it _really_ happen?
		updated, err = client.Template(ctx, template.ID)
		require.NoError(t, err)
		assert.Greater(t, updated.UpdatedAt, template.UpdatedAt)
		assert.Equal(t, *req.Name, updated.Name)
		assert.Equal(t, *req.DisplayName, updated.DisplayName)
		assert.Equal(t, *req.Description, updated.Description)
		assert.Equal(t, *req.Icon, updated.Icon)
		assert.Equal(t, *req.DefaultTTLMillis, updated.DefaultTTLMillis)
		assert.Equal(t, *req.ActivityBumpMillis, updated.ActivityBumpMillis)
		assert.Equal(t, *req.TimeTilAutostopNotifyMillis, updated.TimeTilAutostopNotifyMillis)
		assert.False(t, *req.AllowUserCancelWorkspaceJobs)

		require.Len(t, auditor.AuditLogs(), 5)
		assert.Equal(t, database.AuditActionWrite, auditor.AuditLogs()[4].Action)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		t.Parallel()

		ownerClient := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)
		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.ScopedRoleOrgTemplateAdmin(owner.OrganizationID))

		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		version2 := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)
		template2 := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version2.ID)

		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.UpdateTemplateMeta(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
			Name: &template2.Name,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusConflict, apiErr.StatusCode())
	})

	t.Run("AGPL_Deprecated", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: false})
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		// It is unfortunate we need to sleep, but the test can fail if the
		// updatedAt is too close together.
		time.Sleep(time.Millisecond * 5)

		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			DeprecationMessage: ptr.Ref("APGL cannot deprecate"),
		}

		ctx := testutil.Context(t, testutil.WaitLong)

		updated, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		assert.Greater(t, updated.UpdatedAt, template.UpdatedAt)
		// AGPL cannot deprecate, expect no change
		assert.False(t, updated.Deprecated)
		assert.Empty(t, updated.DeprecationMessage)
	})

	// AGPL cannot deprecate, but it can be unset
	t.Run("AGPL_Unset_Deprecated", func(t *testing.T) {
		t.Parallel()

		owner, db := optimus-ide-collabdtest.NewWithDatabase(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: false})
		user := optimus-ide-collabdtest.CreateFirstUser(t, owner)
		client, tplAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, owner, user.OrganizationID, rbac.RoleTemplateAdmin())
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		// It is unfortunate we need to sleep, but the test can fail if the
		// updatedAt is too close together.
		time.Sleep(time.Millisecond * 5)

		ctx := testutil.Context(t, testutil.WaitLong)

		// nolint:gocritic // Setting up unit test data
		err := db.UpdateTemplateAccessControlByID(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubject(tplAdmin)), database.UpdateTemplateAccessControlByIDParams{
			ID:                   template.ID,
			RequireActiveVersion: false,
			Deprecated:           "Some deprecated message",
		})
		require.NoError(t, err)

		// Check that it is deprecated
		got, err := client.Template(ctx, template.ID)
		require.NoError(t, err)
		require.NotEmpty(t, got.DeprecationMessage, "template is deprecated to start")
		require.True(t, got.Deprecated, "template is deprecated to start")

		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			DeprecationMessage: ptr.Ref(""),
		}

		updated, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		assert.Greater(t, updated.UpdatedAt, template.UpdatedAt)
		assert.False(t, updated.Deprecated)
		assert.Empty(t, updated.DeprecationMessage)
	})

	t.Run("AGPL_MaxPortShareLevel", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: false})
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		require.Equal(t, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic, template.MaxPortShareLevel)

		var level optimus-ide-collabsdk.WorkspaceAgentPortShareLevel = optimus-ide-collabsdk.WorkspaceAgentPortShareLevelAuthenticated
		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			MaxPortShareLevel: &level,
		}

		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		// AGPL cannot change max port sharing level
		require.ErrorContains(t, err, "port sharing level is an enterprise feature")

		// Ensure the same value port share level is a no-op
		level = optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic
		_, err = client.UpdateTemplateMeta(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
			Name:              ptr.Ref(optimus-ide-collabdtest.RandomUsername(t)),
			MaxPortShareLevel: &level,
		})
		require.NoError(t, err)
	})

	t.Run("NoDefaultTTL", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.DefaultTTLMillis = ptr.Ref(24 * time.Hour.Milliseconds())
		})
		// It is unfortunate we need to sleep, but the test can fail if the
		// updatedAt is too close together.
		time.Sleep(time.Millisecond * 5)

		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			DefaultTTLMillis: ptr.Ref(int64(0)),
		}

		// We're too fast! Sleep so we can be sure that updatedAt is greater
		time.Sleep(time.Millisecond * 5)

		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)

		// Extra paranoid: did it _really_ happen?
		updated, err := client.Template(ctx, template.ID)
		require.NoError(t, err)
		assert.Greater(t, updated.UpdatedAt, template.UpdatedAt)
		assert.Equal(t, *req.DefaultTTLMillis, updated.DefaultTTLMillis)
		assert.Empty(t, updated.DeprecationMessage)
		assert.False(t, updated.Deprecated)
	})

	t.Run("DefaultTTLTooLow", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.DefaultTTLMillis = ptr.Ref(24 * time.Hour.Milliseconds())
		})
		// It is unfortunate we need to sleep, but the test can fail if the
		// updatedAt is too close together.
		time.Sleep(time.Millisecond * 5)

		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			DefaultTTLMillis: ptr.Ref(int64(-1)),
		}

		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.ErrorContains(t, err, "default_ttl_ms: Must be a positive integer")

		// Ensure no update occurred
		updated, err := client.Template(ctx, template.ID)
		require.NoError(t, err)
		assert.Equal(t, updated.UpdatedAt, template.UpdatedAt)
		assert.Equal(t, updated.DefaultTTLMillis, template.DefaultTTLMillis)
		assert.Empty(t, updated.DeprecationMessage)
		assert.False(t, updated.Deprecated)
	})

	t.Run("CleanupTTLs", func(t *testing.T) {
		t.Parallel()

		const (
			failureTTL               = 7 * 24 * time.Hour
			inactivityTTL            = 180 * 24 * time.Hour
			timeTilDormantAutoDelete = 360 * 24 * time.Hour
		)

		t.Run("OK", func(t *testing.T) {
			t.Parallel()

			var setCalled atomic.Int64
			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
				TemplateScheduleStore: schedule.MockTemplateScheduleStore{
					SetFn: func(ctx context.Context, db database.Store, template database.Template, options schedule.TemplateScheduleOptions) (database.Template, error) {
						if setCalled.Add(1) == 2 {
							require.Equal(t, failureTTL, options.FailureTTL)
							require.Equal(t, inactivityTTL, options.TimeTilDormant)
							require.Equal(t, timeTilDormantAutoDelete, options.TimeTilDormantAutoDelete)
						}
						template.FailureTTL = int64(options.FailureTTL)
						template.TimeTilDormant = int64(options.TimeTilDormant)
						template.TimeTilDormantAutoDelete = int64(options.TimeTilDormantAutoDelete)
						return template, nil
					},
				},
			})
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
				ctr.FailureTTLMillis = ptr.Ref(0 * time.Hour.Milliseconds())
				ctr.TimeTilDormantMillis = ptr.Ref(0 * time.Hour.Milliseconds())
				ctr.TimeTilDormantAutoDeleteMillis = ptr.Ref(0 * time.Hour.Milliseconds())
			})

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			got, err := client.UpdateTemplateMeta(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
				Name:                           &template.Name,
				DisplayName:                    &template.DisplayName,
				Description:                    &template.Description,
				Icon:                           &template.Icon,
				DefaultTTLMillis:               ptr.Ref(int64(0)),
				AutostopRequirement:            &template.AutostopRequirement,
				AllowUserCancelWorkspaceJobs:   &template.AllowUserCancelWorkspaceJobs,
				FailureTTLMillis:               ptr.Ref(failureTTL.Milliseconds()),
				TimeTilDormantMillis:           ptr.Ref(inactivityTTL.Milliseconds()),
				TimeTilDormantAutoDeleteMillis: ptr.Ref(timeTilDormantAutoDelete.Milliseconds()),
			})
			require.NoError(t, err)

			require.EqualValues(t, 2, setCalled.Load())
			require.Equal(t, failureTTL.Milliseconds(), got.FailureTTLMillis)
			require.Equal(t, inactivityTTL.Milliseconds(), got.TimeTilDormantMillis)
			require.Equal(t, timeTilDormantAutoDelete.Milliseconds(), got.TimeTilDormantAutoDeleteMillis)
		})

		t.Run("IgnoredUnlicensed", func(t *testing.T) {
			t.Parallel()

			client := optimus-ide-collabdtest.New(t, nil)
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
				ctr.FailureTTLMillis = ptr.Ref(0 * time.Hour.Milliseconds())
				ctr.TimeTilDormantMillis = ptr.Ref(0 * time.Hour.Milliseconds())
				ctr.TimeTilDormantAutoDeleteMillis = ptr.Ref(0 * time.Hour.Milliseconds())
			})

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			got, err := client.UpdateTemplateMeta(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
				Name:                           &template.Name,
				DisplayName:                    &template.DisplayName,
				Description:                    &template.Description,
				Icon:                           &template.Icon,
				DefaultTTLMillis:               &template.DefaultTTLMillis,
				AutostopRequirement:            &template.AutostopRequirement,
				AllowUserCancelWorkspaceJobs:   &template.AllowUserCancelWorkspaceJobs,
				FailureTTLMillis:               ptr.Ref(failureTTL.Milliseconds()),
				TimeTilDormantMillis:           ptr.Ref(inactivityTTL.Milliseconds()),
				TimeTilDormantAutoDeleteMillis: ptr.Ref(timeTilDormantAutoDelete.Milliseconds()),
			})
			require.NoError(t, err)
			require.Zero(t, got.FailureTTLMillis)
			require.Zero(t, got.TimeTilDormantMillis)
			require.Zero(t, got.TimeTilDormantAutoDeleteMillis)
			require.Empty(t, got.DeprecationMessage)
			require.False(t, got.Deprecated)
		})
	})

	t.Run("AllowUserScheduling", func(t *testing.T) {
		t.Parallel()

		t.Run("OK", func(t *testing.T) {
			t.Parallel()

			var (
				setCalled      atomic.Int64
				allowAutostart atomic.Bool
				allowAutostop  atomic.Bool
			)
			allowAutostart.Store(true)
			allowAutostop.Store(true)
			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
				TemplateScheduleStore: schedule.MockTemplateScheduleStore{
					SetFn: func(ctx context.Context, db database.Store, template database.Template, options schedule.TemplateScheduleOptions) (database.Template, error) {
						setCalled.Add(1)
						assert.Equal(t, allowAutostart.Load(), options.UserAutostartEnabled)
						assert.Equal(t, allowAutostop.Load(), options.UserAutostopEnabled)

						template.DefaultTTL = int64(options.DefaultTTL)
						template.AllowUserAutostart = options.UserAutostartEnabled
						template.AllowUserAutostop = options.UserAutostopEnabled
						return template, nil
					},
				},
			})
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
				ctr.DefaultTTLMillis = ptr.Ref(24 * time.Hour.Milliseconds())
			})
			require.Equal(t, allowAutostart.Load(), template.AllowUserAutostart)
			require.Equal(t, allowAutostop.Load(), template.AllowUserAutostop)

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			allowAutostart.Store(false)
			allowAutostop.Store(false)
			got, err := client.UpdateTemplateMeta(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
				Name:                         &template.Name,
				DisplayName:                  &template.DisplayName,
				Description:                  &template.Description,
				Icon:                         &template.Icon,
				DefaultTTLMillis:             &template.DefaultTTLMillis,
				AutostopRequirement:          &template.AutostopRequirement,
				AllowUserCancelWorkspaceJobs: &template.AllowUserCancelWorkspaceJobs,
				AllowUserAutostart:           ptr.Ref(allowAutostart.Load()),
				AllowUserAutostop:            ptr.Ref(allowAutostop.Load()),
			})
			require.NoError(t, err)

			require.EqualValues(t, 2, setCalled.Load())
			require.Equal(t, allowAutostart.Load(), got.AllowUserAutostart)
			require.Equal(t, allowAutostop.Load(), got.AllowUserAutostop)
		})

		t.Run("IgnoredUnlicensed", func(t *testing.T) {
			t.Parallel()

			client := optimus-ide-collabdtest.New(t, nil)
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
				ctr.DefaultTTLMillis = ptr.Ref(24 * time.Hour.Milliseconds())
			})

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			got, err := client.UpdateTemplateMeta(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
				Name:                         &template.Name,
				DisplayName:                  &template.DisplayName,
				Description:                  &template.Description,
				Icon:                         &template.Icon,
				DefaultTTLMillis:             ptr.Ref(template.DefaultTTLMillis + 1),
				AutostopRequirement:          &template.AutostopRequirement,
				AllowUserCancelWorkspaceJobs: &template.AllowUserCancelWorkspaceJobs,
				AllowUserAutostart:           ptr.Ref(false),
				AllowUserAutostop:            ptr.Ref(false),
			})
			require.NoError(t, err)
			require.True(t, got.AllowUserAutostart)
			require.True(t, got.AllowUserAutostop)
		})
	})

	t.Run("NotModified", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.Description = "original description"
			ctr.Icon = "/icon/original-icon.png"
			ctr.DefaultTTLMillis = ptr.Ref(24 * time.Hour.Milliseconds())
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			Name:                &template.Name,
			Description:         &template.Description,
			Icon:                &template.Icon,
			DefaultTTLMillis:    &template.DefaultTTLMillis,
			ActivityBumpMillis:  &template.ActivityBumpMillis,
			AutostopRequirement: nil,
			AllowUserAutostart:  &template.AllowUserAutostart,
			AllowUserAutostop:   &template.AllowUserAutostop,
		}
		_, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		updated, err := client.Template(ctx, template.ID)
		require.NoError(t, err)
		assert.Equal(t, template.Name, updated.Name)
		assert.Equal(t, template.Description, updated.Description)
		assert.Equal(t, template.Icon, updated.Icon)
		assert.Equal(t, template.DefaultTTLMillis, updated.DefaultTTLMillis)
		assert.Equal(t, template.ActivityBumpMillis, updated.ActivityBumpMillis)
		assert.Equal(t, template.AllowUserAutostart, updated.AllowUserAutostart)
		assert.Equal(t, template.AllowUserAutostop, updated.AllowUserAutostop)
	})

	t.Run("TimeTilAutostopNotifyPreservedWhenOmitted", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.TimeTilAutostopNotifyMillis = ptr.Ref((5 * time.Minute).Milliseconds())
		})
		require.Equal(t, (5 * time.Minute).Milliseconds(), template.TimeTilAutostopNotifyMillis)

		ctx := testutil.Context(t, testutil.WaitLong)

		// Patch an unrelated field, omitting TimeTilAutostopNotifyMillis.
		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			Description: ptr.Ref("updated description"),
		}
		_, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)

		updated, err := client.Template(ctx, template.ID)
		require.NoError(t, err)
		assert.Equal(t, "updated description", updated.Description)
		assert.Equal(t, template.TimeTilAutostopNotifyMillis, updated.TimeTilAutostopNotifyMillis)
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.Description = "original description"
			ctr.DefaultTTLMillis = ptr.Ref(24 * time.Hour.Milliseconds())
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			DefaultTTLMillis: ptr.Ref(-int64(time.Hour)),
		}
		_, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Contains(t, apiErr.Message, "Invalid request")
		require.Len(t, apiErr.Validations, 1)
		assert.Equal(t, apiErr.Validations[0].Field, "default_ttl_ms")

		updated, err := client.Template(ctx, template.ID)
		require.NoError(t, err)
		assert.WithinDuration(t, template.UpdatedAt, updated.UpdatedAt, time.Minute)
		assert.Equal(t, template.Name, updated.Name)
		assert.Equal(t, template.Description, updated.Description)
		assert.Equal(t, template.Icon, updated.Icon)
		assert.Equal(t, template.DefaultTTLMillis, updated.DefaultTTLMillis)
	})

	t.Run("TimeTilAutostopNotifyInvalid", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)

		ctx := testutil.Context(t, testutil.WaitLong)

		// Sub-minute (non-zero) values are rejected.
		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			TimeTilAutostopNotifyMillis: ptr.Ref(int64(30_000)),
		}
		_, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Contains(t, apiErr.Message, "Invalid request")
		require.Len(t, apiErr.Validations, 1)
		assert.Equal(t, "time_til_autostop_notify_ms", apiErr.Validations[0].Field)
		assert.Equal(t, "Must be 0 (disabled) or at least one minute.", apiErr.Validations[0].Detail)

		// Negative values are rejected.
		req = optimus-ide-collabsdk.UpdateTemplateMeta{
			TimeTilAutostopNotifyMillis: ptr.Ref(int64(-1)),
		}
		_, err = client.UpdateTemplateMeta(ctx, template.ID, req)
		require.ErrorAs(t, err, &apiErr)
		require.Contains(t, apiErr.Message, "Invalid request")
		require.Len(t, apiErr.Validations, 1)
		assert.Equal(t, "time_til_autostop_notify_ms", apiErr.Validations[0].Field)
		assert.Equal(t, "Must be a positive integer.", apiErr.Validations[0].Detail)
	})

	t.Run("TimeTilAutostopNotifyDisableAfterEnable", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.TimeTilAutostopNotifyMillis = ptr.Ref((5 * time.Minute).Milliseconds())
		})
		require.Equal(t, (5 * time.Minute).Milliseconds(), template.TimeTilAutostopNotifyMillis)

		ctx := testutil.Context(t, testutil.WaitLong)

		// Explicitly disable by sending 0, which must not be treated as omitted.
		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			TimeTilAutostopNotifyMillis: ptr.Ref(int64(0)),
		}
		updated, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		assert.Equal(t, int64(0), updated.TimeTilAutostopNotifyMillis)
	})

	t.Run("RemoveIcon", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.Icon = "/icon/code.png"
		})
		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			Icon: ptr.Ref(""),
		}

		ctx := testutil.Context(t, testutil.WaitLong)

		updated, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		assert.Equal(t, updated.Icon, "")
	})

	t.Run("AutostopRequirement", func(t *testing.T) {
		t.Parallel()

		t.Run("OK", func(t *testing.T) {
			t.Parallel()

			var setCalled atomic.Int64
			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
				TemplateScheduleStore: schedule.MockTemplateScheduleStore{
					SetFn: func(ctx context.Context, db database.Store, template database.Template, options schedule.TemplateScheduleOptions) (database.Template, error) {
						if setCalled.Add(1) == 2 {
							assert.EqualValues(t, 0b0110000, options.AutostopRequirement.DaysOfWeek)
							assert.EqualValues(t, 2, options.AutostopRequirement.Weeks)
						}

						err := db.UpdateTemplateScheduleByID(ctx, database.UpdateTemplateScheduleByIDParams{
							ID:                            template.ID,
							UpdatedAt:                     dbtime.Now(),
							AllowUserAutostart:            options.UserAutostartEnabled,
							AllowUserAutostop:             options.UserAutostopEnabled,
							DefaultTTL:                    int64(options.DefaultTTL),
							ActivityBump:                  int64(options.ActivityBump),
							AutostopRequirementDaysOfWeek: int16(options.AutostopRequirement.DaysOfWeek),
							AutostopRequirementWeeks:      options.AutostopRequirement.Weeks,
							FailureTTL:                    int64(options.FailureTTL),
							TimeTilDormant:                int64(options.TimeTilDormant),
							TimeTilDormantAutoDelete:      int64(options.TimeTilDormantAutoDelete),
						})
						if !assert.NoError(t, err) {
							return database.Template{}, err
						}

						return db.GetTemplateByID(ctx, template.ID)
					},
				},
			})
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)

			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
			require.EqualValues(t, 1, setCalled.Load())
			require.Empty(t, template.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 1, template.AutostopRequirement.Weeks)
			req := optimus-ide-collabsdk.UpdateTemplateMeta{
				Name:                         &template.Name,
				DisplayName:                  &template.DisplayName,
				Description:                  &template.Description,
				Icon:                         &template.Icon,
				AllowUserCancelWorkspaceJobs: &template.AllowUserCancelWorkspaceJobs,
				DefaultTTLMillis:             ptr.Ref(time.Hour.Milliseconds()),
				AutostopRequirement: &optimus-ide-collabsdk.TemplateAutostopRequirement{
					// wrong order
					DaysOfWeek: []string{"saturday", "friday"},
					Weeks:      2,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			updated, err := client.UpdateTemplateMeta(ctx, template.ID, req)
			require.NoError(t, err)
			require.EqualValues(t, 2, setCalled.Load())
			require.Equal(t, []string{"friday", "saturday"}, updated.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 2, updated.AutostopRequirement.Weeks)

			template, err = client.Template(ctx, template.ID)
			require.NoError(t, err)
			require.Equal(t, []string{"friday", "saturday"}, template.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 2, template.AutostopRequirement.Weeks)
			require.Empty(t, template.DeprecationMessage)
			require.False(t, template.Deprecated)
		})

		t.Run("Unset", func(t *testing.T) {
			t.Parallel()

			var setCalled atomic.Int64
			client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
				TemplateScheduleStore: schedule.MockTemplateScheduleStore{
					SetFn: func(ctx context.Context, db database.Store, template database.Template, options schedule.TemplateScheduleOptions) (database.Template, error) {
						if setCalled.Add(1) == 2 {
							assert.EqualValues(t, 0, options.AutostopRequirement.DaysOfWeek)
							assert.EqualValues(t, 1, options.AutostopRequirement.Weeks)
						}

						err := db.UpdateTemplateScheduleByID(ctx, database.UpdateTemplateScheduleByIDParams{
							ID:                            template.ID,
							UpdatedAt:                     dbtime.Now(),
							AllowUserAutostart:            options.UserAutostartEnabled,
							AllowUserAutostop:             options.UserAutostopEnabled,
							DefaultTTL:                    int64(options.DefaultTTL),
							ActivityBump:                  int64(options.ActivityBump),
							AutostopRequirementDaysOfWeek: int16(options.AutostopRequirement.DaysOfWeek),
							AutostopRequirementWeeks:      options.AutostopRequirement.Weeks,
							FailureTTL:                    int64(options.FailureTTL),
							TimeTilDormant:                int64(options.TimeTilDormant),
							TimeTilDormantAutoDelete:      int64(options.TimeTilDormantAutoDelete),
						})
						if !assert.NoError(t, err) {
							return database.Template{}, err
						}

						return db.GetTemplateByID(ctx, template.ID)
					},
				},
			})
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)

			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
				ctr.AutostopRequirement = &optimus-ide-collabsdk.TemplateAutostopRequirement{
					// wrong order
					DaysOfWeek: []string{"sunday", "saturday", "friday", "thursday", "wednesday", "tuesday", "monday"},
					Weeks:      2,
				}
			})
			require.EqualValues(t, 1, setCalled.Load())
			require.Equal(t, []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}, template.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 2, template.AutostopRequirement.Weeks)
			req := optimus-ide-collabsdk.UpdateTemplateMeta{
				Name:                         &template.Name,
				DisplayName:                  &template.DisplayName,
				Description:                  &template.Description,
				Icon:                         &template.Icon,
				AllowUserCancelWorkspaceJobs: &template.AllowUserCancelWorkspaceJobs,
				DefaultTTLMillis:             ptr.Ref(time.Hour.Milliseconds()),
				AutostopRequirement: &optimus-ide-collabsdk.TemplateAutostopRequirement{
					DaysOfWeek: []string{},
					Weeks:      0,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			updated, err := client.UpdateTemplateMeta(ctx, template.ID, req)
			require.NoError(t, err)
			require.EqualValues(t, 2, setCalled.Load())
			require.Empty(t, updated.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 1, updated.AutostopRequirement.Weeks)

			template, err = client.Template(ctx, template.ID)
			require.NoError(t, err)
			require.Empty(t, template.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 1, template.AutostopRequirement.Weeks)
		})

		t.Run("EnterpriseOnly", func(t *testing.T) {
			t.Parallel()

			client := optimus-ide-collabdtest.New(t, nil)
			user := optimus-ide-collabdtest.CreateFirstUser(t, client)
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
			require.Empty(t, template.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 1, template.AutostopRequirement.Weeks)
			req := optimus-ide-collabsdk.UpdateTemplateMeta{
				Name:                         &template.Name,
				DisplayName:                  &template.DisplayName,
				Description:                  &template.Description,
				Icon:                         &template.Icon,
				AllowUserCancelWorkspaceJobs: &template.AllowUserCancelWorkspaceJobs,
				DefaultTTLMillis:             ptr.Ref(time.Hour.Milliseconds()),
				AutostopRequirement: &optimus-ide-collabsdk.TemplateAutostopRequirement{
					DaysOfWeek: []string{"monday"},
					Weeks:      2,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			updated, err := client.UpdateTemplateMeta(ctx, template.ID, req)
			require.NoError(t, err)
			require.Empty(t, updated.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 1, updated.AutostopRequirement.Weeks)

			template, err = client.Template(ctx, template.ID)
			require.NoError(t, err)
			require.Empty(t, template.AutostopRequirement.DaysOfWeek)
			require.EqualValues(t, 1, template.AutostopRequirement.Weeks)
			require.Empty(t, template.DeprecationMessage)
			require.False(t, template.Deprecated)
		})
	})

	t.Run("ClassicParameterFlow", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		require.False(t, template.UseClassicParameterFlow, "default is false")

		bTrue := true
		bFalse := false
		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			UseClassicParameterFlow: &bTrue,
		}

		ctx := testutil.Context(t, testutil.WaitLong)

		// set to true
		updated, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		assert.True(t, updated.UseClassicParameterFlow, "expected true")

		req.UseClassicParameterFlow = nil
		_, err = client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)

		updated, err = client.Template(ctx, template.ID)
		require.NoError(t, err)
		assert.True(t, updated.UseClassicParameterFlow, "expected true")

		// back to false
		req.UseClassicParameterFlow = &bFalse
		updated, err = client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		assert.False(t, updated.UseClassicParameterFlow, "expected false")
	})

	t.Run("DisableModuleCache", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		require.False(t, template.DisableModuleCache, "default is false")

		req := optimus-ide-collabsdk.UpdateTemplateMeta{
			DisableModuleCache: ptr.Ref(true),
		}

		ctx := testutil.Context(t, testutil.WaitLong)

		// set to true
		updated, err := client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		assert.True(t, updated.DisableModuleCache, "expected true")

		// Sending DisableModuleCache: nil with no other changes is a true
		// no-op and produces a 304 Not Modified (surfaced as an error by the
		// SDK).
		req.DisableModuleCache = nil
		_, err = client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		updated, err = client.Template(ctx, template.ID)
		require.NoError(t, err)
		assert.True(t, updated.DisableModuleCache, "expected true")

		// back to false
		req.DisableModuleCache = ptr.Ref(false)
		updated, err = client.UpdateTemplateMeta(ctx, template.ID, req)
		require.NoError(t, err)
		assert.False(t, updated.DisableModuleCache, "expected false")
	})

	t.Run("SupportEmptyOrDefaultFields", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)

		displayName := "Test Display Name"
		description := "test-description"
		icon := "/icon/icon.png"
		defaultTTLMillis := 10 * time.Hour.Milliseconds()

		reference := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.DisplayName = displayName
			ctr.Description = description
			ctr.Icon = icon
			ctr.DefaultTTLMillis = ptr.Ref(defaultTTLMillis)
		})
		require.Equal(t, displayName, reference.DisplayName)
		require.Equal(t, description, reference.Description)
		require.Equal(t, icon, reference.Icon)

		restoreReq := optimus-ide-collabsdk.UpdateTemplateMeta{
			DisplayName:      &displayName,
			Description:      &description,
			Icon:             &icon,
			DefaultTTLMillis: ptr.Ref(defaultTTLMillis),
		}

		type expected struct {
			displayName      string
			description      string
			icon             string
			defaultTTLMillis int64
		}

		type testCase struct {
			name     string
			req      optimus-ide-collabsdk.UpdateTemplateMeta
			expected expected
		}

		tests := []testCase{
			{
				name:     "Only update default_ttl_ms",
				req:      optimus-ide-collabsdk.UpdateTemplateMeta{DefaultTTLMillis: ptr.Ref(99 * time.Hour.Milliseconds())},
				expected: expected{displayName: reference.DisplayName, description: reference.Description, icon: reference.Icon, defaultTTLMillis: 99 * time.Hour.Milliseconds()},
			},
			{
				name:     "Clear display name",
				req:      optimus-ide-collabsdk.UpdateTemplateMeta{DisplayName: ptr.Ref("")},
				expected: expected{displayName: "", description: reference.Description, icon: reference.Icon, defaultTTLMillis: defaultTTLMillis},
			},
			{
				name:     "Clear description",
				req:      optimus-ide-collabsdk.UpdateTemplateMeta{Description: ptr.Ref("")},
				expected: expected{displayName: reference.DisplayName, description: "", icon: reference.Icon, defaultTTLMillis: defaultTTLMillis},
			},
			{
				name:     "Clear icon",
				req:      optimus-ide-collabsdk.UpdateTemplateMeta{Icon: ptr.Ref("")},
				expected: expected{displayName: reference.DisplayName, description: reference.Description, icon: "", defaultTTLMillis: defaultTTLMillis},
			},
			// A request whose only field is nil is a true no-op under the new
			// PATCH semantics; the handler returns 304 Not Modified and the
			// template values are preserved.
			{
				name:     "Nil display name is a no-op",
				req:      optimus-ide-collabsdk.UpdateTemplateMeta{DisplayName: nil},
				expected: expected{displayName: reference.DisplayName, description: reference.Description, icon: reference.Icon, defaultTTLMillis: defaultTTLMillis},
			},
			{
				name:     "Nil description is a no-op",
				req:      optimus-ide-collabsdk.UpdateTemplateMeta{Description: nil},
				expected: expected{displayName: reference.DisplayName, description: reference.Description, icon: reference.Icon, defaultTTLMillis: defaultTTLMillis},
			},
			{
				name:     "Nil icon is a no-op",
				req:      optimus-ide-collabsdk.UpdateTemplateMeta{Icon: nil},
				expected: expected{displayName: reference.DisplayName, description: reference.Description, icon: reference.Icon, defaultTTLMillis: defaultTTLMillis},
			},
		}

		for _, tc := range tests {
			//nolint:tparallel,paralleltest
			t.Run(tc.name, func(t *testing.T) {
				defer func() {
					ctx := testutil.Context(t, testutil.WaitLong)
					// Restore reference after each test case. The restore
					// itself can be a no-op (and return an error) when the
					// previous test case was already a no-op; that is
					// expected, so we ignore the error here.
					_, _ = client.UpdateTemplateMeta(ctx, reference.ID, restoreReq)
				}()
				ctx := testutil.Context(t, testutil.WaitLong)
				_, err := client.UpdateTemplateMeta(ctx, reference.ID, tc.req)
				require.NoError(t, err)
				updated, err := client.Template(ctx, reference.ID)
				require.NoError(t, err)
				assert.Equal(t, tc.expected.displayName, updated.DisplayName)
				assert.Equal(t, tc.expected.description, updated.Description)
				assert.Equal(t, tc.expected.icon, updated.Icon)
				assert.Equal(t, tc.expected.defaultTTLMillis, updated.DefaultTTLMillis)
			})
		}
	})

	// EmptyBodyPreservesAllFields ensures the PATCH endpoint treats an empty
	// body as a no-op so that omitted fields do not overwrite existing values.
	t.Run("EmptyBodyPreservesAllFields", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.DisplayName = "Original Display"
			ctr.Description = "Original description"
			ctr.Icon = "/icon/original.png"
			ctr.DefaultTTLMillis = ptr.Ref((24 * time.Hour).Milliseconds())
			ctr.AllowUserCancelWorkspaceJobs = ptr.Ref(true)
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.UpdateTemplateMeta(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateMeta{})
		require.NoError(t, err)

		updated, err := client.Template(ctx, template.ID)
		require.NoError(t, err)
		assert.Equal(t, template.Name, updated.Name)
		assert.Equal(t, template.DisplayName, updated.DisplayName)
		assert.Equal(t, template.Description, updated.Description)
		assert.Equal(t, template.Icon, updated.Icon)
		assert.Equal(t, template.DefaultTTLMillis, updated.DefaultTTLMillis)
		assert.Equal(t, template.AllowUserCancelWorkspaceJobs, updated.AllowUserCancelWorkspaceJobs)
		assert.Equal(t, template.RequireActiveVersion, updated.RequireActiveVersion)
	})

	// PartialUpdatePreservesOtherFields ensures sending a single field on the
	// PATCH body changes only that field and leaves the others alone. This is
	// the headline behavior PLAT-184 enables: previously, omitted booleans
	// were silently overwritten with false because the SDK type used
	// non-pointer booleans.
	t.Run("PartialUpdatePreservesOtherFields", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
			ctr.AllowUserCancelWorkspaceJobs = ptr.Ref(true)
			ctr.DefaultTTLMillis = ptr.Ref((24 * time.Hour).Milliseconds())
		})
		require.True(t, template.AllowUserCancelWorkspaceJobs)
		require.Equal(t, (24 * time.Hour).Milliseconds(), template.DefaultTTLMillis)

		ctx := testutil.Context(t, testutil.WaitLong)

		// Sending only DefaultTTLMillis must not flip AllowUserCancelWorkspaceJobs
		// to false.
		newTTL := (12 * time.Hour).Milliseconds()
		updated, err := client.UpdateTemplateMeta(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
			DefaultTTLMillis: &newTTL,
		})
		require.NoError(t, err)
		assert.Equal(t, newTTL, updated.DefaultTTLMillis)
		assert.True(t, updated.AllowUserCancelWorkspaceJobs, "omitted bool field must not be overwritten")

		// Conversely, sending only AllowUserCancelWorkspaceJobs must not zero
		// out DefaultTTLMillis.
		updated, err = client.UpdateTemplateMeta(ctx, template.ID, optimus-ide-collabsdk.UpdateTemplateMeta{
			AllowUserCancelWorkspaceJobs: ptr.Ref(false),
		})
		require.NoError(t, err)
		assert.False(t, updated.AllowUserCancelWorkspaceJobs)
		assert.Equal(t, newTTL, updated.DefaultTTLMillis, "omitted int64 field must not be overwritten")
	})
}

func TestDeleteTemplate(t *testing.T) {
	t.Parallel()

	t.Run("NoWorkspaces", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true, Auditor: auditor})
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)

		ctx := testutil.Context(t, testutil.WaitLong)

		err := client.DeleteTemplate(ctx, template.ID)
		require.NoError(t, err)

		require.Len(t, auditor.AuditLogs(), 5)
		assert.Equal(t, database.AuditActionDelete, auditor.AuditLogs()[4].Action)
	})

	t.Run("Workspaces", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)

		ctx := testutil.Context(t, testutil.WaitLong)

		err := client.DeleteTemplate(ctx, template.ID)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})

	t.Run("NoPermission", func(t *testing.T) {
		t.Parallel()
		client, db := optimus-ide-collabdtest.NewWithDatabase(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		tpl := dbfake.TemplateVersion(t, db).Seed(database.TemplateVersion{CreatedBy: owner.UserID, OrganizationID: owner.OrganizationID}).Do()

		ctx := testutil.Context(t, testutil.WaitShort)
		err := memberClient.DeleteTemplate(ctx, tpl.Template.ID)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusForbidden, apiErr.StatusCode())
	})

	t.Run("OnlyPrebuilds", func(t *testing.T) {
		t.Parallel()
		client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		tpl := dbfake.TemplateVersion(t, db).
			Seed(database.TemplateVersion{
				CreatedBy:      owner.UserID,
				OrganizationID: owner.OrganizationID,
			}).Do()

		// Create a workspace owned by the prebuilds system user.
		dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        database.PrebuildsSystemUserID,
			OrganizationID: owner.OrganizationID,
			TemplateID:     tpl.Template.ID,
		}).Seed(database.WorkspaceBuild{
			TemplateVersionID: tpl.TemplateVersion.ID,
		}).Do()

		ctx := testutil.Context(t, testutil.WaitLong)

		err := client.DeleteTemplate(ctx, tpl.Template.ID)
		require.NoError(t, err)
	})

	t.Run("PrebuildsAndHumanWorkspaces", func(t *testing.T) {
		t.Parallel()
		client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		tpl := dbfake.TemplateVersion(t, db).
			Seed(database.TemplateVersion{
				CreatedBy:      owner.UserID,
				OrganizationID: owner.OrganizationID,
			}).Do()

		// Create a prebuild workspace.
		dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        database.PrebuildsSystemUserID,
			OrganizationID: owner.OrganizationID,
			TemplateID:     tpl.Template.ID,
		}).Seed(database.WorkspaceBuild{
			TemplateVersionID: tpl.TemplateVersion.ID,
		}).Do()

		// Create a human-owned workspace.
		dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        owner.UserID,
			OrganizationID: owner.OrganizationID,
			TemplateID:     tpl.Template.ID,
		}).Seed(database.WorkspaceBuild{
			TemplateVersionID: tpl.TemplateVersion.ID,
		}).Do()

		ctx := testutil.Context(t, testutil.WaitLong)

		err := client.DeleteTemplate(ctx, tpl.Template.ID)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})

	t.Run("DeletedIsSet", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)

		ctx := testutil.Context(t, testutil.WaitLong)

		// Verify the deleted field is exposed in the SDK and set to false for active templates
		got, err := client.Template(ctx, template.ID)
		require.NoError(t, err)
		require.False(t, got.Deleted)
	})

	t.Run("DeletedIsTrue", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)

		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)

		ctx := testutil.Context(t, testutil.WaitLong)

		err := client.DeleteTemplate(ctx, template.ID)
		require.NoError(t, err)

		// Verify the deleted field is set to true by listing templates with
		// deleted:true filter.
		templates, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
			OrganizationID: user.OrganizationID,
			SearchQuery:    "deleted:true",
		})
		require.NoError(t, err)

		require.Len(t, templates, 1)
		require.Equal(t, template.ID, templates[0].ID)
		require.True(t, templates[0].Deleted)
	})
}

func TestTemplateMetrics(t *testing.T) {
	t.Parallel()

	t.Skip("flaky test: https://github.com/optimus-ide-collab/optimus-ide-collab/issues/6481")

	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		IncludeProvisionerDaemon:    true,
		AgentStatsRefreshInterval:   time.Millisecond * 100,
		MetricsCacheRefreshInterval: time.Millisecond * 100,
	})

	user := optimus-ide-collabdtest.CreateFirstUser(t, client)
	authToken := uuid.NewString()
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, &echo.Responses{
		Parse:          echo.ParseComplete,
		ProvisionPlan:  echo.PlanComplete,
		ProvisionGraph: echo.ProvisionGraphWithAgent(authToken),
	})
	template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
	require.Equal(t, -1, template.ActiveUserCount)
	require.Empty(t, template.BuildTimeStats[optimus-ide-collabsdk.WorkspaceTransitionStart])

	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

	_ = agenttest.New(t, client.URL, authToken)
	resources := optimus-ide-collabdtest.AwaitWorkspaceAgents(t, client, workspace.ID)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	daus, err := client.TemplateDAUs(context.Background(), template.ID, optimus-ide-collabsdk.TimezoneOffsetHour(time.UTC))
	require.NoError(t, err)

	require.Equal(t, &optimus-ide-collabsdk.DAUsResponse{
		Entries: []optimus-ide-collabsdk.DAUEntry{},
	}, daus, "no DAUs when stats are empty")

	res, err := client.Workspaces(ctx, optimus-ide-collabsdk.WorkspaceFilter{})
	require.NoError(t, err)
	assert.Zero(t, res.Workspaces[0].LastUsedAt)

	conn, err := workspacesdk.New(client).
		DialAgent(ctx, resources[0].Agents[0].ID, &workspacesdk.DialAgentOptions{
			Logger: testutil.Logger(t).Named("tailnet"),
		})
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	sshConn, err := conn.SSHClient(ctx)
	require.NoError(t, err)
	_ = sshConn.Close()

	wantDAUs := &optimus-ide-collabsdk.DAUsResponse{
		Entries: []optimus-ide-collabsdk.DAUEntry{
			{
				Date:   time.Now().UTC().Truncate(time.Hour * 24).Format("2006-01-02"),
				Amount: 1,
			},
		},
	}
	require.Eventuallyf(t, func() bool {
		daus, err = client.TemplateDAUs(ctx, template.ID, optimus-ide-collabsdk.TimezoneOffsetHour(time.UTC))
		require.NoError(t, err)
		return len(daus.Entries) > 0
	},
		testutil.WaitShort, testutil.IntervalFast,
		"template daus never loaded",
	)
	gotDAUs, err := client.TemplateDAUs(ctx, template.ID, optimus-ide-collabsdk.TimezoneOffsetHour(time.UTC))
	require.NoError(t, err)
	require.Equal(t, gotDAUs, wantDAUs)

	template, err = client.Template(ctx, template.ID)
	require.NoError(t, err)
	require.Equal(t, 1, template.ActiveUserCount)

	require.Eventuallyf(t, func() bool {
		template, err = client.Template(ctx, template.ID)
		require.NoError(t, err)
		startMs := template.BuildTimeStats[optimus-ide-collabsdk.WorkspaceTransitionStart].P50
		return startMs != nil && *startMs > 1
	},
		testutil.WaitShort, testutil.IntervalFast,
		"BuildTimeStats never loaded",
	)

	res, err = client.Workspaces(ctx, optimus-ide-collabsdk.WorkspaceFilter{})
	require.NoError(t, err)
	assert.WithinDuration(t,
		dbtime.Now(), res.Workspaces[0].LastUsedAt, time.Minute,
	)
}

func TestTemplateNotifications(t *testing.T) {
	t.Parallel()

	t.Run("Delete", func(t *testing.T) {
		t.Parallel()

		t.Run("InitiatorIsNotNotified", func(t *testing.T) {
			t.Parallel()

			// Given: an initiator
			var (
				notifyEnq = &notificationstest.FakeEnqueuer{}
				client    = optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
					IncludeProvisionerDaemon: true,
					NotificationsEnqueuer:    notifyEnq,
				})
				initiator = optimus-ide-collabdtest.CreateFirstUser(t, client)
				version   = optimus-ide-collabdtest.CreateTemplateVersion(t, client, initiator.OrganizationID, nil)
				_         = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
				template  = optimus-ide-collabdtest.CreateTemplate(t, client, initiator.OrganizationID, version.ID)
				ctx       = testutil.Context(t, testutil.WaitLong)
			)

			// When: the template is deleted by the initiator
			err := client.DeleteTemplate(ctx, template.ID)
			require.NoError(t, err)

			// Then: the delete notification is not sent to the initiator.
			deleteNotifications := make([]*notificationstest.FakeNotification, 0)
			for _, n := range notifyEnq.Sent() {
				if n.TemplateID == notifications.TemplateTemplateDeleted {
					deleteNotifications = append(deleteNotifications, n)
				}
			}
			require.Len(t, deleteNotifications, 0)
		})

		t.Run("OnlyOwnersAndAdminsAreNotified", func(t *testing.T) {
			t.Parallel()

			// Given: multiple users with different roles
			var (
				notifyEnq = &notificationstest.FakeEnqueuer{}
				client    = optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
					IncludeProvisionerDaemon: true,
					NotificationsEnqueuer:    notifyEnq,
				})
				initiator = optimus-ide-collabdtest.CreateFirstUser(t, client)
				ctx       = testutil.Context(t, testutil.WaitLong)

				// Setup template
				version  = optimus-ide-collabdtest.CreateTemplateVersion(t, client, initiator.OrganizationID, nil)
				_        = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
				template = optimus-ide-collabdtest.CreateTemplate(t, client, initiator.OrganizationID, version.ID, func(ctr *optimus-ide-collabsdk.CreateTemplateRequest) {
					ctr.DisplayName = "Bobby's Template"
				})
			)

			// Setup users with different roles
			_, owner := optimus-ide-collabdtest.CreateAnotherUser(t, client, initiator.OrganizationID, rbac.RoleOwner())
			_, tmplAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, client, initiator.OrganizationID, rbac.RoleTemplateAdmin())
			optimus-ide-collabdtest.CreateAnotherUser(t, client, initiator.OrganizationID, rbac.RoleMember())
			optimus-ide-collabdtest.CreateAnotherUser(t, client, initiator.OrganizationID, rbac.RoleUserAdmin())
			optimus-ide-collabdtest.CreateAnotherUser(t, client, initiator.OrganizationID, rbac.RoleAuditor())

			// When: the template is deleted by the initiator
			err := client.DeleteTemplate(ctx, template.ID)
			require.NoError(t, err)

			// Then: only owners and template admins should receive the
			// notification.
			shouldBeNotified := []uuid.UUID{owner.ID, tmplAdmin.ID}
			var deleteTemplateNotifications []*notificationstest.FakeNotification
			for _, n := range notifyEnq.Sent() {
				if n.TemplateID == notifications.TemplateTemplateDeleted {
					deleteTemplateNotifications = append(deleteTemplateNotifications, n)
				}
			}
			notifiedUsers := make([]uuid.UUID, 0, len(deleteTemplateNotifications))
			for _, n := range deleteTemplateNotifications {
				notifiedUsers = append(notifiedUsers, n.UserID)
			}
			require.ElementsMatch(t, shouldBeNotified, notifiedUsers)

			// Validate the notification content
			for _, n := range deleteTemplateNotifications {
				require.Equal(t, n.TemplateID, notifications.TemplateTemplateDeleted)
				require.Contains(t, notifiedUsers, n.UserID)
				require.Contains(t, n.Targets, template.ID)
				require.Contains(t, n.Targets, template.OrganizationID)
				require.Equal(t, n.Labels["name"], template.DisplayName)
				require.Equal(t, n.Labels["initiator"], optimus-ide-collabdtest.FirstUserParams.Username)
			}
		})
	})
}

func TestTemplateFilterHasAITask(t *testing.T) {
	t.Parallel()

	db, pubsub := dbtestutil.NewDB(t)
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		Database:                 db,
		Pubsub:                   pubsub,
		IncludeProvisionerDaemon: true,
	})
	user := optimus-ide-collabdtest.CreateFirstUser(t, client)

	jobWithAITask := dbgen.ProvisionerJob(t, db, pubsub, database.ProvisionerJob{
		OrganizationID: user.OrganizationID,
		InitiatorID:    user.UserID,
		Tags:           database.StringMap{},
		Type:           database.ProvisionerJobTypeTemplateVersionImport,
	})
	jobWithoutAITask := dbgen.ProvisionerJob(t, db, pubsub, database.ProvisionerJob{
		OrganizationID: user.OrganizationID,
		InitiatorID:    user.UserID,
		Tags:           database.StringMap{},
		Type:           database.ProvisionerJobTypeTemplateVersionImport,
	})
	versionWithAITask := dbgen.TemplateVersion(t, db, database.TemplateVersion{
		OrganizationID: user.OrganizationID,
		CreatedBy:      user.UserID,
		HasAITask:      sql.NullBool{Bool: true, Valid: true},
		JobID:          jobWithAITask.ID,
	})
	versionWithoutAITask := dbgen.TemplateVersion(t, db, database.TemplateVersion{
		OrganizationID: user.OrganizationID,
		CreatedBy:      user.UserID,
		HasAITask:      sql.NullBool{Bool: false, Valid: true},
		JobID:          jobWithoutAITask.ID,
	})
	templateWithAITask := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, versionWithAITask.ID)
	templateWithoutAITask := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, versionWithoutAITask.ID)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	// Test filtering
	templates, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
		SearchQuery: "has-ai-task:true",
	})
	require.NoError(t, err)
	require.Len(t, templates, 1)
	require.Equal(t, templateWithAITask.ID, templates[0].ID)

	templates, err = client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
		SearchQuery: "has-ai-task:false",
	})
	require.NoError(t, err)
	require.Len(t, templates, 1)
	require.Equal(t, templateWithoutAITask.ID, templates[0].ID)

	templates, err = client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{})
	require.NoError(t, err)
	require.Len(t, templates, 2)
	require.Contains(t, templates, templateWithAITask)
	require.Contains(t, templates, templateWithoutAITask)
}

func TestTemplateFilterHasExternalAgent(t *testing.T) {
	t.Parallel()

	db, pubsub := dbtestutil.NewDB(t)
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		Database:                 db,
		Pubsub:                   pubsub,
		IncludeProvisionerDaemon: true,
	})
	user := optimus-ide-collabdtest.CreateFirstUser(t, client)

	jobWithExternalAgent := dbgen.ProvisionerJob(t, db, pubsub, database.ProvisionerJob{
		OrganizationID: user.OrganizationID,
		InitiatorID:    user.UserID,
		Tags:           database.StringMap{},
		Type:           database.ProvisionerJobTypeTemplateVersionImport,
	})
	jobWithoutExternalAgent := dbgen.ProvisionerJob(t, db, pubsub, database.ProvisionerJob{
		OrganizationID: user.OrganizationID,
		InitiatorID:    user.UserID,
		Tags:           database.StringMap{},
		Type:           database.ProvisionerJobTypeTemplateVersionImport,
	})
	versionWithExternalAgent := dbgen.TemplateVersion(t, db, database.TemplateVersion{
		OrganizationID:   user.OrganizationID,
		CreatedBy:        user.UserID,
		HasExternalAgent: sql.NullBool{Bool: true, Valid: true},
		JobID:            jobWithExternalAgent.ID,
	})
	versionWithoutExternalAgent := dbgen.TemplateVersion(t, db, database.TemplateVersion{
		OrganizationID:   user.OrganizationID,
		CreatedBy:        user.UserID,
		HasExternalAgent: sql.NullBool{Bool: false, Valid: true},
		JobID:            jobWithoutExternalAgent.ID,
	})
	templateWithExternalAgent := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, versionWithExternalAgent.ID)
	templateWithoutExternalAgent := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, versionWithoutExternalAgent.ID)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	templates, err := client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
		SearchQuery: "has_external_agent:true",
	})
	require.NoError(t, err)
	require.Len(t, templates, 1)
	require.Equal(t, templateWithExternalAgent.ID, templates[0].ID)

	templates, err = client.Templates(ctx, optimus-ide-collabsdk.TemplateFilter{
		SearchQuery: "has_external_agent:false",
	})
	require.NoError(t, err)
	require.Len(t, templates, 1)
	require.Equal(t, templateWithoutExternalAgent.ID, templates[0].ID)
}
