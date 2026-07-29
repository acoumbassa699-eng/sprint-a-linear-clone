package optimus-ide-collabd_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"

	"cdr.dev/slog/v3"
	"cdr.dev/slog/v3/sloggers/slogtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/agent"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/agent/agenttest"
	agploptimus-ide-collabd "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd"
	agplaudit "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbmock"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/pubsub"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/entitlements"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	agplprebuilds "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/prebuilds"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/namesgenerator"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	natspubsub "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/nats"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/workspacesdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/prebuilds"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/dbcrypt"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/replicasync"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisionersdk/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/tailnet/tailnettest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/retry"
	"github.com/optimus-ide-collab/serpent"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m, testutil.GoleakOptions...)
}

func TestEntitlements(t *testing.T) {
	t.Parallel()
	t.Run("NoLicense", func(t *testing.T) {
		t.Parallel()
		adminClient, _, api, adminUser := optimus-ide-collabdenttest.NewWithAPI(t, &optimus-ide-collabdenttest.Options{
			DontAddLicense: true,
		})
		anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)
		res, err := anotherClient.Entitlements(context.Background())
		require.NoError(t, err)
		require.False(t, res.HasLicense)
		require.Empty(t, res.Warnings)

		// Ensure the entitlements are the same reference
		require.Equal(t, fmt.Sprintf("%p", api.Entitlements), fmt.Sprintf("%p", api.AGPL.Entitlements))
	})
	t.Run("FullLicense", func(t *testing.T) {
		t.Parallel()
		adminClient, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			AuditLogging:   true,
			DontAddLicense: true,
		})
		// Enable all features
		features := make(license.Features)
		for _, feature := range optimus-ide-collabsdk.FeatureNames {
			features[feature] = 1
		}
		features[optimus-ide-collabsdk.FeatureUserLimit] = 100
		optimus-ide-collabdenttest.AddLicense(t, adminClient, optimus-ide-collabdenttest.LicenseOptions{
			Features: features,
			Addons:   []optimus-ide-collabsdk.Addon{optimus-ide-collabsdk.AddonAIGovernance},
			GraceAt:  time.Now().Add(59 * 24 * time.Hour),
		})
		res, err := adminClient.Entitlements(context.Background()) //nolint:gocritic // adding another user would put us over user limit
		require.NoError(t, err)
		assert.True(t, res.HasLicense)
		ul := res.Features[optimus-ide-collabsdk.FeatureUserLimit]
		assert.Equal(t, optimus-ide-collabsdk.EntitlementEntitled, ul.Entitlement)
		if assert.NotNil(t, ul.Limit) {
			assert.Equal(t, int64(100), *ul.Limit)
		}
		if assert.NotNil(t, ul.Actual) {
			assert.Equal(t, int64(1), *ul.Actual)
		}
		assert.True(t, ul.Enabled)
		al := res.Features[optimus-ide-collabsdk.FeatureAuditLog]
		assert.Equal(t, optimus-ide-collabsdk.EntitlementEntitled, al.Entitlement)
		assert.True(t, al.Enabled)
		assert.Nil(t, al.Limit)
		assert.Nil(t, al.Actual)
		assert.Empty(t, res.Warnings)
	})

	// TestEntitlements/MultiplePrebuildsLicenseUpdates verifies that uploading
	// multiple licenses with prebuilds enabled doesn't cause a panic from
	// duplicate Prometheus metric registration. This was a bug where the new
	// reconciler's metrics were registered before the old reconciler was stopped.
	t.Run("MultiplePrebuildsLicenseUpdates", func(t *testing.T) {
		t.Parallel()
		adminClient, _, api, _ := optimus-ide-collabdenttest.NewWithAPI(t, &optimus-ide-collabdenttest.Options{
			DontAddLicense: true,
		})

		// Add first license with prebuilds to initialize the reconciler
		features := license.Features{
			optimus-ide-collabsdk.FeatureUserLimit:          100,
			optimus-ide-collabsdk.FeatureWorkspacePrebuilds: 1,
		}
		license1 := optimus-ide-collabdenttest.AddLicense(t, adminClient, optimus-ide-collabdenttest.LicenseOptions{
			Features: features,
		})
		res, err := adminClient.Entitlements(context.Background())
		require.NoError(t, err)
		require.True(t, res.HasLicense)
		require.Equal(t, optimus-ide-collabsdk.EntitlementEntitled, res.Features[optimus-ide-collabsdk.FeatureWorkspacePrebuilds].Entitlement)

		// Verify the reconciler was set up
		reconciler1 := api.AGPL.PrebuildsReconciler.Load()
		require.NotNil(t, reconciler1)

		// Delete the license to disable prebuilds, then add a new one.
		// This tests the enabled -> disabled -> enabled transition.
		err = adminClient.DeleteLicense(context.Background(), license1.ID)
		require.NoError(t, err)

		optimus-ide-collabdenttest.AddLicense(t, adminClient, optimus-ide-collabdenttest.LicenseOptions{
			Features: features,
		})
		res, err = adminClient.Entitlements(context.Background())
		require.NoError(t, err)
		require.True(t, res.HasLicense)
		require.Equal(t, optimus-ide-collabsdk.EntitlementEntitled, res.Features[optimus-ide-collabsdk.FeatureWorkspacePrebuilds].Entitlement)

		// Verify a new reconciler was created
		reconciler2 := api.AGPL.PrebuildsReconciler.Load()
		require.NotNil(t, reconciler2)
	})
	t.Run("FullLicenseToNone", func(t *testing.T) {
		t.Parallel()
		adminClient, adminUser := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			AuditLogging:   true,
			DontAddLicense: true,
		})
		anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)
		license := optimus-ide-collabdenttest.AddLicense(t, adminClient, optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureUserLimit: 100,
				optimus-ide-collabsdk.FeatureAuditLog:  1,
			},
		})
		res, err := anotherClient.Entitlements(context.Background())
		require.NoError(t, err)
		assert.True(t, res.HasLicense)
		al := res.Features[optimus-ide-collabsdk.FeatureAuditLog]
		assert.Equal(t, optimus-ide-collabsdk.EntitlementEntitled, al.Entitlement)
		assert.True(t, al.Enabled)

		err = adminClient.DeleteLicense(context.Background(), license.ID)
		require.NoError(t, err)

		res, err = anotherClient.Entitlements(context.Background())
		require.NoError(t, err)
		assert.False(t, res.HasLicense)
		al = res.Features[optimus-ide-collabsdk.FeatureAuditLog]
		assert.Equal(t, optimus-ide-collabsdk.EntitlementNotEntitled, al.Entitlement)
		assert.False(t, al.Enabled)
	})
	t.Run("Pubsub", func(t *testing.T) {
		t.Parallel()
		adminClient, _, api, adminUser := optimus-ide-collabdenttest.NewWithAPI(t, &optimus-ide-collabdenttest.Options{DontAddLicense: true})
		anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)
		entitlements, err := anotherClient.Entitlements(context.Background())
		require.NoError(t, err)
		require.False(t, entitlements.HasLicense)
		ctx := testDBAuthzRole(context.Background())
		_, err = api.Database.InsertLicense(ctx, database.InsertLicenseParams{
			UploadedAt: dbtime.Now(),
			Exp:        dbtime.Now().AddDate(1, 0, 0),
			JWT: optimus-ide-collabdenttest.GenerateLicense(t, optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAuditLog: 1,
				},
			}),
		})
		require.NoError(t, err)
		err = api.ReplicaSyncPubsub.Publish(optimus-ide-collabd.PubsubEventLicenses, []byte{})
		require.NoError(t, err)
		require.Eventually(t, func() bool {
			entitlements, err := anotherClient.Entitlements(context.Background())
			assert.NoError(t, err)
			return entitlements.HasLicense
		}, testutil.WaitShort, testutil.IntervalFast)
	})
	t.Run("Resync", func(t *testing.T) {
		t.Parallel()
		adminClient, _, api, adminUser := optimus-ide-collabdenttest.NewWithAPI(t, &optimus-ide-collabdenttest.Options{
			EntitlementsUpdateInterval: 25 * time.Millisecond,
			DontAddLicense:             true,
		})
		anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)
		entitlements, err := anotherClient.Entitlements(context.Background())
		require.NoError(t, err)
		require.False(t, entitlements.HasLicense)
		// Valid
		ctx := context.Background()
		_, err = api.Database.InsertLicense(testDBAuthzRole(ctx), database.InsertLicenseParams{
			UploadedAt: dbtime.Now(),
			Exp:        dbtime.Now().AddDate(1, 0, 0),
			JWT: optimus-ide-collabdenttest.GenerateLicense(t, optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAuditLog: 1,
				},
			}),
		})
		require.NoError(t, err)
		// Expired
		_, err = api.Database.InsertLicense(testDBAuthzRole(ctx), database.InsertLicenseParams{
			UploadedAt: dbtime.Now(),
			Exp:        dbtime.Now().AddDate(-1, 0, 0),
			JWT: optimus-ide-collabdenttest.GenerateLicense(t, optimus-ide-collabdenttest.LicenseOptions{
				ExpiresAt: dbtime.Now().AddDate(-1, 0, 0),
			}),
		})
		require.NoError(t, err)
		// Invalid
		_, err = api.Database.InsertLicense(testDBAuthzRole(ctx), database.InsertLicenseParams{
			UploadedAt: dbtime.Now(),
			Exp:        dbtime.Now().AddDate(1, 0, 0),
			JWT:        "invalid",
		})
		require.NoError(t, err)
		require.Eventually(t, func() bool {
			entitlements, err := anotherClient.Entitlements(context.Background())
			assert.NoError(t, err)
			return entitlements.HasLicense
		}, testutil.WaitShort, testutil.IntervalFast)
	})
}

func TestEntitlements_HeaderWarnings(t *testing.T) {
	t.Parallel()
	t.Run("ExistForAdmin", func(t *testing.T) {
		t.Parallel()
		adminClient, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			AuditLogging: true,
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				AllFeatures: false,
			},
		})
		//nolint:gocritic // This isn't actually bypassing any RBAC checks
		res, err := adminClient.Request(context.Background(), http.MethodGet, "/api/v2/users/me", nil)
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NotEmpty(t, res.Header.Values(optimus-ide-collabsdk.EntitlementsWarningHeader))
	})
	t.Run("NoneForNormalUser", func(t *testing.T) {
		t.Parallel()
		adminClient, adminUser := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			AuditLogging: true,
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				AllFeatures: false,
			},
		})
		anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)
		res, err := anotherClient.Request(context.Background(), http.MethodGet, "/api/v2/users/me", nil)
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Empty(t, res.Header.Values(optimus-ide-collabsdk.EntitlementsWarningHeader))
	})
}

func TestEntitlements_Prebuilds(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		featureEnabled  bool
		expectedEnabled bool
	}{
		{
			name:            "Feature enabled",
			featureEnabled:  true,
			expectedEnabled: true,
		},
		{
			name:            "Feature disabled",
			featureEnabled:  false,
			expectedEnabled: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var prebuildsEntitled int64
			if tc.featureEnabled {
				prebuildsEntitled = 1
			}

			_, _, api, _ := optimus-ide-collabdenttest.NewWithAPI(t, &optimus-ide-collabdenttest.Options{
				Options: &optimus-ide-collabdtest.Options{
					DeploymentValues: optimus-ide-collabdtest.DeploymentValues(t),
				},

				EntitlementsUpdateInterval: time.Second,
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureWorkspacePrebuilds: prebuildsEntitled,
					},
				},
			})

			// The entitlements will need to refresh before the reconciler is set.
			require.Eventually(t, func() bool {
				return api.AGPL.PrebuildsReconciler.Load() != nil
			}, testutil.WaitSuperLong, testutil.IntervalFast)

			reconciler := api.AGPL.PrebuildsReconciler.Load()
			claimer := api.AGPL.PrebuildsClaimer.Load()
			require.NotNil(t, reconciler)
			require.NotNil(t, claimer)

			if tc.expectedEnabled {
				require.IsType(t, &prebuilds.StoreReconciler{}, *reconciler)
				require.IsType(t, &prebuilds.EnterpriseClaimer{}, *claimer)
			} else {
				require.Equal(t, &agplprebuilds.DefaultReconciler, reconciler)
				require.Equal(t, &agplprebuilds.DefaultClaimer, claimer)
			}
		})
	}
}

func TestAuditLogging(t *testing.T) {
	t.Parallel()
	t.Run("Enabled", func(t *testing.T) {
		t.Parallel()
		db, _ := dbtestutil.NewDB(t)
		_, _, api, _ := optimus-ide-collabdenttest.NewWithAPI(t, &optimus-ide-collabdenttest.Options{
			AuditLogging: true,
			Options: &optimus-ide-collabdtest.Options{
				Auditor: audit.NewAuditor(db, audit.DefaultFilter),
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAuditLog: 1,
				},
			},
		})
		db, _ = dbtestutil.NewDB(t)
		auditor := *api.AGPL.Auditor.Load()
		ea := audit.NewAuditor(db, audit.DefaultFilter)
		t.Logf("%T = %T", auditor, ea)
		assert.EqualValues(t, reflect.ValueOf(ea).Type(), reflect.ValueOf(auditor).Type())
	})
	t.Run("Disabled", func(t *testing.T) {
		t.Parallel()
		_, _, api, _ := optimus-ide-collabdenttest.NewWithAPI(t, &optimus-ide-collabdenttest.Options{DontAddLicense: true})
		auditor := *api.AGPL.Auditor.Load()
		ea := agplaudit.NewNop()
		t.Logf("%T = %T", auditor, ea)
		assert.Equal(t, reflect.ValueOf(ea).Type(), reflect.ValueOf(auditor).Type())
	})
	// The AGPL code runs with a fake auditor that doesn't represent the real implementation.
	// We do a simple test to ensure that basic flows function.
	t.Run("FullBuild", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitLong)
		client, user := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			DontAddLicense: true,
		})
		r := setupWorkspaceAgent(t, client, user, 0)
		conn, err := workspacesdk.New(client).DialAgent(ctx, r.sdkAgent.ID, nil) //nolint:gocritic // RBAC is not the purpose of this test
		require.NoError(t, err)
		defer conn.Close()
		connected := conn.AwaitReachable(ctx)
		require.True(t, connected)
		_ = r.agent.Close() // close first so we don't drop error logs from outdated build
		build := optimus-ide-collabdtest.CreateWorkspaceBuild(t, client, r.workspace, database.WorkspaceTransitionStop)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, build.ID)
	})
}

func TestExternalTokenEncryption(t *testing.T) {
	t.Parallel()

	t.Run("Enabled", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)
		db, ps := dbtestutil.NewDB(t)
		ciphers, err := dbcrypt.NewCiphers(bytes.Repeat([]byte("a"), 32))
		require.NoError(t, err)
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			EntitlementsUpdateInterval: 25 * time.Millisecond,
			ExternalTokenEncryption:    ciphers,
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureExternalTokenEncryption: 1,
				},
			},
			Options: &optimus-ide-collabdtest.Options{
				Database: db,
				Pubsub:   ps,
			},
		})
		keys, err := db.GetDBCryptKeys(ctx)
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, ciphers[0].HexDigest(), keys[0].ActiveKeyDigest.String)

		require.Eventually(t, func() bool {
			entitlements, err := client.Entitlements(context.Background())
			assert.NoError(t, err)
			feature := entitlements.Features[optimus-ide-collabsdk.FeatureExternalTokenEncryption]
			entitled := feature.Entitlement == optimus-ide-collabsdk.EntitlementEntitled
			var warningExists bool
			for _, warning := range entitlements.Warnings {
				if strings.Contains(warning, optimus-ide-collabsdk.FeatureExternalTokenEncryption.Humanize()) {
					warningExists = true
					break
				}
			}
			t.Logf("feature: %+v, warnings: %+v, errors: %+v", feature, entitlements.Warnings, entitlements.Errors)
			return feature.Enabled && entitled && !warningExists
		}, testutil.WaitShort, testutil.IntervalFast)
	})

	t.Run("Disabled", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)
		db, ps := dbtestutil.NewDB(t)
		ciphers, err := dbcrypt.NewCiphers()
		require.NoError(t, err)
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			DontAddLicense:             true,
			EntitlementsUpdateInterval: 25 * time.Millisecond,
			ExternalTokenEncryption:    ciphers,
			Options: &optimus-ide-collabdtest.Options{
				Database: db,
				Pubsub:   ps,
			},
		})
		keys, err := db.GetDBCryptKeys(ctx)
		require.NoError(t, err)
		require.Empty(t, keys)

		require.Eventually(t, func() bool {
			entitlements, err := client.Entitlements(context.Background())
			assert.NoError(t, err)
			feature := entitlements.Features[optimus-ide-collabsdk.FeatureExternalTokenEncryption]
			entitled := feature.Entitlement == optimus-ide-collabsdk.EntitlementEntitled
			var warningExists bool
			for _, warning := range entitlements.Warnings {
				if strings.Contains(warning, optimus-ide-collabsdk.FeatureExternalTokenEncryption.Humanize()) {
					warningExists = true
					break
				}
			}
			t.Logf("feature: %+v, warnings: %+v, errors: %+v", feature, entitlements.Warnings, entitlements.Errors)
			return !feature.Enabled && !entitled && !warningExists
		}, testutil.WaitShort, testutil.IntervalFast)
	})

	t.Run("PreviouslyEnabledButMissingFromLicense", func(t *testing.T) {
		// If this test fails, it potentially means that a customer who has
		// actively been using this feature is now unable _start optimus-ide-collabd_
		// because of a licensing issue. This should never happen.
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)
		db, ps := dbtestutil.NewDB(t)
		ciphers, err := dbcrypt.NewCiphers(bytes.Repeat([]byte("a"), 32))
		require.NoError(t, err)

		dbc, err := dbcrypt.New(ctx, db, ciphers...) // should insert key
		require.NoError(t, err)

		keys, err := dbc.GetDBCryptKeys(ctx)
		require.NoError(t, err)
		require.Len(t, keys, 1)

		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			DontAddLicense:             true,
			EntitlementsUpdateInterval: 25 * time.Millisecond,
			ExternalTokenEncryption:    ciphers,
			Options: &optimus-ide-collabdtest.Options{
				Database: db,
				Pubsub:   ps,
			},
		})

		require.Eventually(t, func() bool {
			entitlements, err := client.Entitlements(context.Background())
			assert.NoError(t, err)
			feature := entitlements.Features[optimus-ide-collabsdk.FeatureExternalTokenEncryption]
			entitled := feature.Entitlement == optimus-ide-collabsdk.EntitlementEntitled
			var warningExists bool
			for _, warning := range entitlements.Warnings {
				if strings.Contains(warning, optimus-ide-collabsdk.FeatureExternalTokenEncryption.Humanize()) {
					warningExists = true
					break
				}
			}
			t.Logf("feature: %+v, warnings: %+v, errors: %+v", feature, entitlements.Warnings, entitlements.Errors)
			return feature.Enabled && !entitled && warningExists
		}, testutil.WaitShort, testutil.IntervalFast)
	})
}

func TestMultiReplica_EmptyRelayAddress(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)
	db, ps := dbtestutil.NewDB(t)
	logger := testutil.Logger(t)

	_, _ = optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		EntitlementsUpdateInterval: 25 * time.Millisecond,
		ReplicaSyncUpdateInterval:  25 * time.Millisecond,
		Options: &optimus-ide-collabdtest.Options{
			Logger:   &logger,
			Database: db,
			Pubsub:   ps,
		},
	})

	mgr, err := replicasync.New(ctx, logger, db, ps, &replicasync.Options{
		ID:             uuid.New(),
		RelayAddress:   "",
		RegionID:       999,
		UpdateInterval: testutil.IntervalFast,
	})
	require.NoError(t, err)
	defer mgr.Close()

	// Send a bunch of updates to see if the optimus-ide-collabd will log errors.
	{
		ctx, cancel := context.WithTimeout(ctx, testutil.IntervalMedium)
		for r := retry.New(testutil.IntervalFast, testutil.IntervalFast); r.Wait(ctx); {
			require.NoError(t, mgr.PublishUpdate())
		}
		cancel()
	}
}

func TestMultiReplica_EmptyRelayAddress_DisabledDERP(t *testing.T) {
	t.Parallel()

	derpMap, _ := tailnettest.RunDERPAndSTUN(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpapi.Write(context.Background(), w, http.StatusOK, derpMap)
	}))
	t.Cleanup(srv.Close)

	ctx := testutil.Context(t, testutil.WaitLong)
	db, ps := dbtestutil.NewDB(t)
	logger := testutil.Logger(t)

	dv := optimus-ide-collabdtest.DeploymentValues(t)
	dv.DERP.Server.Enable = serpent.Bool(false)
	dv.DERP.Config.URL = serpent.String(srv.URL)

	_, _ = optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		EntitlementsUpdateInterval: 25 * time.Millisecond,
		ReplicaSyncUpdateInterval:  25 * time.Millisecond,
		Options: &optimus-ide-collabdtest.Options{
			Logger:           &logger,
			Database:         db,
			Pubsub:           ps,
			DeploymentValues: dv,
		},
	})

	mgr, err := replicasync.New(ctx, logger, db, ps, &replicasync.Options{
		ID:             uuid.New(),
		RelayAddress:   "",
		RegionID:       999,
		UpdateInterval: testutil.IntervalFast,
	})
	require.NoError(t, err)
	defer mgr.Close()

	// Send a bunch of updates to see if the optimus-ide-collabd will log errors.
	{
		ctx, cancel := context.WithTimeout(ctx, testutil.IntervalMedium)
		for r := retry.New(testutil.IntervalFast, testutil.IntervalFast); r.Wait(ctx); {
			require.NoError(t, mgr.PublishUpdate())
		}
		cancel()
	}
}

func TestMultiReplica_NATSPubsubPeers(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)
	logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true}).Leveled(slog.LevelDebug)
	db, pgPubsub := dbtestutil.NewDB(t)
	clusterToken := "shared-token"

	natsA, err := natspubsub.New(ctx, logger.Named("nats-a"), natspubsub.Options{
		ClusterHost:      "127.0.0.1",
		ClusterPort:      natsserver.RANDOM_PORT,
		ClusterAuthToken: clusterToken,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = natsA.Close() })

	dv := optimus-ide-collabdtest.DeploymentValues(t)
	dv.Experiments = []string{string(optimus-ide-collabsdk.ExperimentNATSPubsub)}
	_, _ = optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		EntitlementsUpdateInterval: 25 * time.Millisecond,
		ReplicaSyncUpdateInterval:  25 * time.Millisecond,
		Options: &optimus-ide-collabdtest.Options{
			Logger:            &logger,
			Database:          db,
			Pubsub:            natsA,
			ReplicaSyncPubsub: pgPubsub.(*pubsub.PGPubsub),
			DeploymentValues:  dv,
		},
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureHighAvailability: 1,
			},
		},
	})

	natsB, err := natspubsub.New(ctx, logger.Named("nats-b"), natspubsub.Options{
		ClusterHost:      "127.0.0.1",
		ClusterPort:      natsserver.RANDOM_PORT,
		ClusterAuthToken: clusterToken,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = natsB.Close() })

	mgr, err := replicasync.New(ctx, logger.Named("replica-b"), db, pgPubsub, &replicasync.Options{
		ID:             uuid.New(),
		ClusterHost:    "127.0.0.1",
		RegionID:       12345,
		UpdateInterval: testutil.IntervalFast,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = mgr.Close() })
	require.NotNil(t, natsB.Server.ClusterAddr())
	// nolint: gosec // nats listens on TCP ports
	mgr.SetSelfNATSPort(int32(natsB.Server.ClusterAddr().Port))

	subject := "nats.replica"
	messages := make(chan []byte, 1)
	cancel, err := natsB.Subscribe(subject, func(_ context.Context, msg []byte) {
		messages <- msg
	})
	require.NoError(t, err)
	defer cancel()

	payload := []byte("from-replicasync-peers")
	var publishErr error
	var flushErr error
	var updateErr error
	require.Eventually(t, func() bool {
		updateErr = mgr.PublishUpdate()
		if updateErr != nil {
			return false
		}
		publishErr = natsA.Publish(subject, payload)
		if publishErr != nil {
			return false
		}
		flushErr = natsA.Flush()
		if flushErr != nil {
			return false
		}
		select {
		case got := <-messages:
			return string(got) == string(payload)
		default:
			return false
		}
	}, testutil.WaitShort, testutil.IntervalFast)
	require.NoError(t, updateErr)
	require.NoError(t, publishErr)
	require.NoError(t, flushErr)
}

func TestSCIMDisabled(t *testing.T) {
	t.Parallel()

	cli, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{})

	checkPaths := []string{
		"/scim/v2",
		"/scim/v2/",
		"/scim/v2/users",
		"/scim/v2/Users",
		"/scim/v2/Users/",
		"/scim/v2/random/path/that/is/long",
		"/scim/v2/random/path/that/is/long.txt",
	}

	client := &http.Client{}
	for _, p := range checkPaths {
		t.Run(p, func(t *testing.T) {
			t.Parallel()

			u, err := cli.URL.Parse(p)
			require.NoError(t, err)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusNotFound, resp.StatusCode)

			var apiError optimus-ide-collabsdk.Response
			err = json.NewDeoptimus-ide-collab(resp.Body).Decode(&apiError)
			require.NoError(t, err)

			require.Contains(t, apiError.Message, "SCIM is disabled")
		})
	}
}

func TestManagedAgentLimit(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)

	cli, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		Options: &optimus-ide-collabdtest.Options{
			IncludeProvisionerDaemon: true,
		},
		LicenseOptions: (&optimus-ide-collabdenttest.LicenseOptions{
			FeatureSet: optimus-ide-collabsdk.FeatureSetPremium,
			// Make it expire in the distant future so it doesn't generate
			// expiry warnings.
			GraceAt:   time.Now().Add(time.Hour * 24 * 60),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 90),
		}).ManagedAgentLimit(1),
	})

	// Get entitlements to check that the license is a-ok.
	sdkEntitlements, err := cli.Entitlements(ctx) //nolint:gocritic // we're not testing authz on the entitlements endpoint, so using owner is fine
	require.NoError(t, err)
	require.True(t, sdkEntitlements.HasLicense)
	agentLimit := sdkEntitlements.Features[optimus-ide-collabsdk.FeatureManagedAgentLimit]
	require.True(t, agentLimit.Enabled)
	require.NotNil(t, agentLimit.Limit)
	require.EqualValues(t, 1, *agentLimit.Limit)
	require.Empty(t, sdkEntitlements.Errors)

	// Create a fake provision response that claims there are agents in the
	// template and every built workspace.
	//
	// It's fine that the app ID is only used in a single successful workspace
	// build.
	appID := uuid.NewString()
	echoRes := &echo.Responses{
		Parse:         echo.ParseComplete,
		ProvisionInit: echo.InitComplete,
		ProvisionPlan: []*proto.Response{
			{
				Type: &proto.Response_Plan{
					Plan: &proto.PlanComplete{
						Plan: []byte("{}"),
					},
				},
			},
		},
		ProvisionApply: echo.ApplyComplete,
		ProvisionGraph: []*proto.Response{{
			Type: &proto.Response_Graph{
				Graph: &proto.GraphComplete{
					Resources: []*proto.Resource{{
						Name: "example",
						Type: "aws_instance",
						Agents: []*proto.Agent{{
							Id:   uuid.NewString(),
							Name: "example",
							Auth: &proto.Agent_Token{
								Token: uuid.NewString(),
							},
							Apps: []*proto.App{{
								Id:   appID,
								Slug: "test",
								Url:  "http://localhost:1234",
							}},
						}},
					}},
					AiTasks: []*proto.AITask{{
						Id: uuid.NewString(),
						SidebarApp: &proto.AITaskSidebarApp{
							Id: appID,
						},
					}},
				},
			},
		}},
	}

	// Create two templates, one with AI and one without.
	aiVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, cli, uuid.Nil, echoRes)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, cli, aiVersion.ID)
	aiTemplate := optimus-ide-collabdtest.CreateTemplate(t, cli, uuid.Nil, aiVersion.ID)
	noAiVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, cli, uuid.Nil, nil) // use default responses
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, cli, noAiVersion.ID)
	noAiTemplate := optimus-ide-collabdtest.CreateTemplate(t, cli, uuid.Nil, noAiVersion.ID)

	// Create one AI workspace, which should succeed.
	task, err := cli.CreateTask(ctx, owner.UserID.String(), optimus-ide-collabsdk.CreateTaskRequest{
		Name:                    namesgenerator.UniqueNameWith("-"),
		TemplateVersionID:       aiTemplate.ActiveVersionID,
		TemplateVersionPresetID: uuid.Nil,
		Input:                   "hi",
		DisplayName:             namesgenerator.UniqueName(),
	})
	require.NoError(t, err, "creating task for AI workspace must succeed")
	workspace, err := cli.Workspace(ctx, task.WorkspaceID.UUID)
	require.NoError(t, err, "fetching AI workspace must succeed")
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, cli, workspace.LatestBuild.ID)

	// Create a second AI task, which should succeed even though the limit is
	// breached. Managed agent limits are advisory only and should never block
	// workspace creation.
	task2, err := cli.CreateTask(ctx, owner.UserID.String(), optimus-ide-collabsdk.CreateTaskRequest{
		Name:                    namesgenerator.UniqueNameWith("-"),
		TemplateVersionID:       aiTemplate.ActiveVersionID,
		TemplateVersionPresetID: uuid.Nil,
		Input:                   "hi",
		DisplayName:             namesgenerator.UniqueName(),
	})
	require.NoError(t, err, "creating task beyond managed agent limit must succeed")
	workspace2, err := cli.Workspace(ctx, task2.WorkspaceID.UUID)
	require.NoError(t, err, "fetching AI workspace must succeed")
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, cli, workspace2.LatestBuild.ID)

	// Create a third workspace using the same template, which should succeed.
	workspace = optimus-ide-collabdtest.CreateWorkspace(t, cli, aiTemplate.ID)
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, cli, workspace.LatestBuild.ID)

	// Create a fourth non-AI workspace, which should also succeed.
	workspace = optimus-ide-collabdtest.CreateWorkspace(t, cli, noAiTemplate.ID)
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, cli, workspace.LatestBuild.ID)
}

func TestCheckBuildUsage_NeverBlocksOnManagedAgentLimit(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare entitlements with a managed agent limit.
	entSet := entitlements.New()
	entSet.Modify(func(e *optimus-ide-collabsdk.Entitlements) {
		e.HasLicense = true
		limit := int64(1)
		issuedAt := time.Now().Add(-2 * time.Hour)
		start := time.Now().Add(-time.Hour)
		end := time.Now().Add(time.Hour)
		e.Features[optimus-ide-collabsdk.FeatureManagedAgentLimit] = optimus-ide-collabsdk.Feature{
			Enabled:     true,
			Limit:       &limit,
			UsagePeriod: &optimus-ide-collabsdk.UsagePeriod{IssuedAt: issuedAt, Start: start, End: end},
		}
	})

	// Enterprise API instance with entitlements injected.
	agpl := &agploptimus-ide-collabd.API{
		Options: &agploptimus-ide-collabd.Options{
			Entitlements: entSet,
		},
	}
	eapi := &optimus-ide-collabd.API{
		AGPL:    agpl,
		Options: &optimus-ide-collabd.Options{Options: agpl.Options},
	}

	// Template version that has an AI task.
	tv := &database.TemplateVersion{
		HasAITask:        sql.NullBool{Valid: true, Bool: true},
		HasExternalAgent: sql.NullBool{Valid: true, Bool: false},
	}

	task := &database.Task{
		TemplateVersionID: tv.ID,
	}

	// Mock DB: no calls expected since managed agent limits are
	// advisory only and no longer query the database at build time.
	mDB := dbmock.NewMockStore(ctrl)

	ctx := context.Background()

	// Start transition: should be permitted even though the limit is
	// breached. Managed agent limits are advisory only.
	startResp, err := eapi.CheckBuildUsage(ctx, mDB, tv, task, database.WorkspaceTransitionStart)
	require.NoError(t, err)
	require.True(t, startResp.Permitted)

	// Stop transition: should also be permitted.
	stopResp, err := eapi.CheckBuildUsage(ctx, mDB, tv, task, database.WorkspaceTransitionStop)
	require.NoError(t, err)
	require.True(t, stopResp.Permitted)

	// Delete transition: should also be permitted.
	deleteResp, err := eapi.CheckBuildUsage(ctx, mDB, tv, task, database.WorkspaceTransitionDelete)
	require.NoError(t, err)
	require.True(t, deleteResp.Permitted)
}

func TestCheckBuildUsage_BlocksStartWithoutEntitlement(t *testing.T) {
	t.Parallel()

	managedAgentVersion := &database.TemplateVersion{
		HasAITask:        sql.NullBool{Valid: true, Bool: true},
		HasExternalAgent: sql.NullBool{Valid: true, Bool: false},
	}
	externalAgentVersion := &database.TemplateVersion{
		HasExternalAgent: sql.NullBool{Valid: true, Bool: true},
	}

	tests := []struct {
		name             string
		templateVersion  *database.TemplateVersion
		task             *database.Task
		setupEnts        func(e *optimus-ide-collabsdk.Entitlements)
		message          string
		checkNoTaskStart bool
	}{
		{
			name:            "ManagedAgentFeatureAbsent",
			templateVersion: managedAgentVersion,
			task:            &database.Task{TemplateVersionID: managedAgentVersion.ID},
			setupEnts: func(e *optimus-ide-collabsdk.Entitlements) {
				e.HasLicense = true
				delete(e.Features, optimus-ide-collabsdk.FeatureManagedAgentLimit)
			},
			message:          "not entitled to managed agents",
			checkNoTaskStart: true,
		},
		{
			name:            "ManagedAgentFeatureDisabled",
			templateVersion: managedAgentVersion,
			task:            &database.Task{TemplateVersionID: managedAgentVersion.ID},
			setupEnts: func(e *optimus-ide-collabsdk.Entitlements) {
				e.HasLicense = true
				e.Features[optimus-ide-collabsdk.FeatureManagedAgentLimit] = optimus-ide-collabsdk.Feature{
					Enabled: false,
				}
			},
			message:          "not entitled to managed agents",
			checkNoTaskStart: true,
		},
		{
			name:            "ExternalAgentFeatureAbsent",
			templateVersion: externalAgentVersion,
			setupEnts: func(e *optimus-ide-collabsdk.Entitlements) {
				e.HasLicense = true
				delete(e.Features, optimus-ide-collabsdk.FeatureWorkspaceExternalAgent)
			},
			message: "uses external agents",
		},
		{
			name:            "ExternalAgentFeatureDisabled",
			templateVersion: externalAgentVersion,
			setupEnts: func(e *optimus-ide-collabsdk.Entitlements) {
				e.HasLicense = true
				e.Features[optimus-ide-collabsdk.FeatureWorkspaceExternalAgent] = optimus-ide-collabsdk.Feature{
					Enabled: false,
				}
			},
			message: "uses external agents",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			entSet := entitlements.New()
			entSet.Modify(tc.setupEnts)

			agpl := &agploptimus-ide-collabd.API{
				Options: &agploptimus-ide-collabd.Options{
					Entitlements: entSet,
				},
			}
			eapi := &optimus-ide-collabd.API{
				AGPL:    agpl,
				Options: &optimus-ide-collabd.Options{Options: agpl.Options},
			}

			mDB := dbmock.NewMockStore(ctrl)
			ctx := context.Background()

			resp, err := eapi.CheckBuildUsage(ctx, mDB, tc.templateVersion, tc.task, database.WorkspaceTransitionStart)
			require.NoError(t, err)
			require.False(t, resp.Permitted)
			require.Contains(t, resp.Message, tc.message)

			stopResp, err := eapi.CheckBuildUsage(ctx, mDB, tc.templateVersion, tc.task, database.WorkspaceTransitionStop)
			require.NoError(t, err)
			require.True(t, stopResp.Permitted)

			deleteResp, err := eapi.CheckBuildUsage(ctx, mDB, tc.templateVersion, tc.task, database.WorkspaceTransitionDelete)
			require.NoError(t, err)
			require.True(t, deleteResp.Permitted)

			if tc.checkNoTaskStart {
				noTaskResp, err := eapi.CheckBuildUsage(ctx, mDB, tc.templateVersion, nil, database.WorkspaceTransitionStart)
				require.NoError(t, err)
				require.True(t, noTaskResp.Permitted)
			}
		})
	}
}

// testDBAuthzRole returns a context with a subject that has a role
// with permissions required for test setup.
func testDBAuthzRole(ctx context.Context) context.Context {
	return dbauthz.As(ctx, rbac.Subject{
		ID: uuid.Nil.String(),
		Roles: rbac.Roles([]rbac.Role{
			{
				Identifier:  rbac.RoleIdentifier{Name: "testing"},
				DisplayName: "Unit Tests",
				Site: rbac.Permissions(map[string][]policy.Action{
					rbac.ResourceWildcard.Type: {policy.WildcardSymbol},
				}),
				User:    []rbac.Permission{},
				ByOrgID: map[string]rbac.OrgPermissions{},
			},
		}),
		Scope: rbac.ScopeAll,
	})
}

// restartableListener is a TCP listener that can have all of it's connections
// severed on demand.
type restartableListener struct {
	net.Listener
	mu    sync.Mutex
	conns []net.Conn
}

func (l *restartableListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	l.mu.Lock()
	l.conns = append(l.conns, conn)
	l.mu.Unlock()
	return conn, nil
}

func (l *restartableListener) CloseConnections() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, conn := range l.conns {
		_ = conn.Close()
	}
	l.conns = nil
}

type restartableTestServer struct {
	options *optimus-ide-collabdenttest.Options
	rl      *restartableListener

	mu     sync.Mutex
	api    *optimus-ide-collabd.API
	closer io.Closer
}

func newRestartableTestServer(t *testing.T, options *optimus-ide-collabdenttest.Options) (*optimus-ide-collabsdk.Client, optimus-ide-collabsdk.CreateFirstUserResponse, *restartableTestServer) {
	t.Helper()
	if options == nil {
		options = &optimus-ide-collabdenttest.Options{}
	}

	s := &restartableTestServer{
		options: options,
	}
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		api := s.api
		s.mu.Unlock()

		if api == nil {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("server is not started"))
			return
		}
		api.AGPL.RootHandler.ServeHTTP(w, r)
	}))
	s.rl = &restartableListener{Listener: srv.Listener}
	srv.Listener = s.rl
	srv.Start()
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	require.NoError(t, err, "failed to parse server URL")
	s.options.AccessURL = u

	client, firstUser := s.startWithFirstUser(t)
	client.URL = u
	return client, firstUser, s
}

func (s *restartableTestServer) Stop(t *testing.T) {
	t.Helper()

	s.mu.Lock()
	closer := s.closer
	s.closer = nil
	api := s.api
	s.api = nil
	s.mu.Unlock()

	if closer != nil {
		err := closer.Close()
		require.NoError(t, err)
	}
	if api != nil {
		err := api.Close()
		require.NoError(t, err)
	}

	s.rl.CloseConnections()
}

func (s *restartableTestServer) Start(t *testing.T) {
	t.Helper()
	_, _ = s.startWithFirstUser(t)
}

func (s *restartableTestServer) startWithFirstUser(t *testing.T) (client *optimus-ide-collabsdk.Client, firstUser optimus-ide-collabsdk.CreateFirstUserResponse) {
	t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closer != nil || s.api != nil {
		t.Fatal("server already started, close must be called first")
	}
	// This creates it's own TCP listener unfortunately, but it's not being
	// used in this test.
	client, s.closer, s.api, firstUser = optimus-ide-collabdenttest.NewWithAPI(t, s.options)

	// Never add the first user or license on subsequent restarts.
	s.options.DontAddFirstUser = true
	s.options.DontAddLicense = true

	return client, firstUser
}

// Test_CoordinatorRollingRestart tests that two peers can maintain a connection
// without forgetting about each other when a HA coordinator does a rolling
// restart.
//
// We had a few issues with this in the past:
//  1. We didn't allow clients to maintain their peer ID after a reconnect,
//     which resulted in the other peer thinking the client was a new peer.
//     (This is fixed and independently tested in AGPL code)
//  2. HA coordinators would delete all peers (via FK constraints) when they
//     were closed, which meant tunnels would be deleted and peers would be
//     notified that the other peer was permanently gone.
//     (This is fixed and independently tested above)
//
// This test uses a real server and real clients.
func TestConn_CoordinatorRollingRestart(t *testing.T) {
	t.Parallel()

	// Although DERP will have connection issues until the connection is
	// reestablished, any open connections should be maintained.
	//
	// Direct connections should be able to transmit packets throughout the
	// restart without issue.
	//nolint:paralleltest // Outdated rule
	for _, direct := range []bool{true, false} {
		name := "DERP"
		if direct {
			name = "Direct"
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			store, ps := dbtestutil.NewDB(t)
			dv := optimus-ide-collabdtest.DeploymentValues(t, func(dv *optimus-ide-collabsdk.DeploymentValues) {
				dv.DERP.Config.BlockDirect = serpent.Bool(!direct)
			})
			logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true}).Leveled(slog.LevelDebug)

			// Create two restartable test servers with the same database.
			client1, user, s1 := newRestartableTestServer(t, &optimus-ide-collabdenttest.Options{
				DontAddFirstUser: false,
				DontAddLicense:   false,
				Options: &optimus-ide-collabdtest.Options{
					Logger:                   ptr.Ref(logger.Named("server1")),
					Database:                 store,
					Pubsub:                   ps,
					DeploymentValues:         dv,
					IncludeProvisionerDaemon: true,
				},
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureHighAvailability: 1,
					},
				},
			})
			client2, _, s2 := newRestartableTestServer(t, &optimus-ide-collabdenttest.Options{
				DontAddFirstUser: true,
				DontAddLicense:   true,
				Options: &optimus-ide-collabdtest.Options{
					Logger:           ptr.Ref(logger.Named("server2")),
					Database:         store,
					Pubsub:           ps,
					DeploymentValues: dv,
				},
			})
			client2.SetSessionToken(client1.SessionToken())

			workspace := dbfake.WorkspaceBuild(t, store, database.WorkspaceTable{
				OrganizationID: user.OrganizationID,
				OwnerID:        user.UserID,
			}).WithAgent().Do()

			// Agent connects via the first coordinator.
			_ = agenttest.New(t, client1.URL, workspace.AgentToken, func(o *agent.Options) {
				o.Logger = logger.Named("agent1")
			})
			resources := optimus-ide-collabdtest.NewWorkspaceAgentWaiter(t, client1, workspace.Workspace.ID).Wait()

			agentID := uuid.Nil
			for _, r := range resources {
				for _, a := range r.Agents {
					agentID = a.ID
					break
				}
			}
			require.NotEqual(t, uuid.Nil, agentID)

			// Client connects via the second coordinator.
			ctx := testutil.Context(t, testutil.WaitSuperLong)
			workspaceClient2 := workspacesdk.New(client2)
			conn, err := workspaceClient2.DialAgent(ctx, agentID, &workspacesdk.DialAgentOptions{
				Logger: logger.Named("client"),
			})
			require.NoError(t, err)
			defer conn.Close()

			require.Eventually(t, func() bool {
				_, p2p, _, err := conn.Ping(ctx)
				assert.NoError(t, err)
				return p2p == direct
			}, testutil.WaitShort, testutil.IntervalFast)

			// Open a TCP server and connection to it through the tunnel that
			// should be maintained throughout the restart.
			tcpServerAddr := tcpEchoServer(t)
			tcpConn, err := conn.DialContext(ctx, "tcp", tcpServerAddr)
			require.NoError(t, err)
			defer tcpConn.Close()
			writeReadEcho(t, ctx, tcpConn)

			// Stop the first server.
			logger.Info(ctx, "test: stopping server 1")
			s1.Stop(t)

			// Pings should fail on DERP but succeed on direct connections.
			pingCtx, pingCancel := context.WithTimeout(ctx, 2*time.Second) //nolint:gocritic // it's going to hang and timeout for DERP, so this needs to be short
			defer pingCancel()
			_, p2p, _, err := conn.Ping(pingCtx)
			if direct {
				require.NoError(t, err)
				require.True(t, p2p, "expected direct connection")
			} else {
				require.ErrorIs(t, err, context.DeadlineExceeded)
			}

			// The existing TCP connection should still be working if we're
			// using direct connections.
			if direct {
				writeReadEcho(t, ctx, tcpConn)
			}

			// Start the first server again.
			logger.Info(ctx, "test: starting server 1")
			s1.Start(t)

			// Restart the second server.
			logger.Info(ctx, "test: stopping server 2")
			s2.Stop(t)
			logger.Info(ctx, "test: starting server 2")
			s2.Start(t)

			// Pings should eventually succeed on both DERP and direct
			// connections.
			require.True(t, conn.AwaitReachable(ctx))
			_, p2p, _, err = conn.Ping(ctx)
			require.NoError(t, err)
			require.Equal(t, direct, p2p, "mismatched p2p state")

			// The existing TCP connection should still be working.
			writeReadEcho(t, ctx, tcpConn)
		})
	}
}

func tcpEchoServer(t *testing.T) string {
	tcpListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tcpListener.Close()
	})
	go func() {
		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				return
			}
			t.Cleanup(func() {
				_ = conn.Close()
			})
			go func() {
				defer conn.Close()
				_, _ = io.Copy(conn, conn)
			}()
		}
	}()
	return tcpListener.Addr().String()
}

// nolint:revive // t takes precedence.
func writeReadEcho(t *testing.T, ctx context.Context, conn net.Conn) {
	msg := namesgenerator.UniqueName()

	deadline, ok := ctx.Deadline()
	if ok {
		_ = conn.SetWriteDeadline(deadline)
		defer conn.SetWriteDeadline(time.Time{})
		_ = conn.SetReadDeadline(deadline)
		defer conn.SetReadDeadline(time.Time{})
	}

	// Write a message
	_, err := conn.Write([]byte(msg))
	require.NoError(t, err)

	// Read the message back
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	require.NoError(t, err)
	require.Equal(t, msg, string(buf[:n]))
}
