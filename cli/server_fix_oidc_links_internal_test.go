package cli

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"cdr.dev/slog/v3"
	"cdr.dev/slog/v3/sloggers/slogtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/serpent"
)

func fakeOIDCIssuerDiscovery(t *testing.T, issuer string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/.well-known/openid-configuration" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEnoptimus-ide-collab(w).Encode(map[string]string{
			"issuer": issuer,
		})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newDeploymentValues(t *testing.T, issuerURL string, autoRepair bool) *optimus-ide-collabsdk.DeploymentValues {
	t.Helper()
	vals := &optimus-ide-collabsdk.DeploymentValues{}
	require.NoError(t, vals.OIDC.IssuerURL.Set(issuerURL))
	vals.OIDC.AutoRepairLinks = serpent.Bool(autoRepair)
	return vals
}

func TestOIDCAuthLinks(t *testing.T) {
	t.Parallel()

	t.Run("RepairsMismatchedLinks", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitShort)
		logger := testutil.Logger(t)

		const expectedIssuer = "https://accounts.google.com"
		oidcSrv := fakeOIDCIssuerDiscovery(t, expectedIssuer)

		db, _ := dbtestutil.NewDB(t)

		// Correctly linked user.
		correctUser := dbgen.User(t, db, database.User{LoginType: database.LoginTypeOIDC})
		dbgen.UserLink(t, db, database.UserLink{
			UserID:    correctUser.ID,
			LoginType: database.LoginTypeOIDC,
			LinkedID:  expectedIssuer + "||sub-correct",
		})

		// Mismatched user.
		mismatchedUser := dbgen.User(t, db, database.User{LoginType: database.LoginTypeOIDC})
		dbgen.UserLink(t, db, database.UserLink{
			UserID:    mismatchedUser.ID,
			LoginType: database.LoginTypeOIDC,
			LinkedID:  "https://old-issuer.example.com||sub-mismatched",
		})

		vals := newDeploymentValues(t, oidcSrv.URL, true)
		err := oidcAuthLinks(ctx, logger, oidcSrv.Client(), vals, db)
		require.NoError(t, err)

		// Mismatched link should be reset.
		link, err := db.GetUserLinkByUserIDLoginType(ctx, database.GetUserLinkByUserIDLoginTypeParams{
			UserID:    mismatchedUser.ID,
			LoginType: database.LoginTypeOIDC,
		})
		require.NoError(t, err)
		require.Equal(t, "", link.LinkedID)

		// Correct link should be untouched.
		link, err = db.GetUserLinkByUserIDLoginType(ctx, database.GetUserLinkByUserIDLoginTypeParams{
			UserID:    correctUser.ID,
			LoginType: database.LoginTypeOIDC,
		})
		require.NoError(t, err)
		require.Equal(t, expectedIssuer+"||sub-correct", link.LinkedID)
	})

	t.Run("AutoRepairDisabledSkipsReset", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitShort)
		logger := testutil.Logger(t)

		const expectedIssuer = "https://accounts.google.com"
		oidcSrv := fakeOIDCIssuerDiscovery(t, expectedIssuer)

		db, _ := dbtestutil.NewDB(t)

		mismatchedUser := dbgen.User(t, db, database.User{LoginType: database.LoginTypeOIDC})
		dbgen.UserLink(t, db, database.UserLink{
			UserID:    mismatchedUser.ID,
			LoginType: database.LoginTypeOIDC,
			LinkedID:  "https://old-issuer.example.com||sub-mismatched",
		})

		vals := newDeploymentValues(t, oidcSrv.URL, false)
		err := oidcAuthLinks(ctx, logger, oidcSrv.Client(), vals, db)
		require.NoError(t, err)

		// Link must not be modified when auto-repair is off.
		link, err := db.GetUserLinkByUserIDLoginType(ctx, database.GetUserLinkByUserIDLoginTypeParams{
			UserID:    mismatchedUser.ID,
			LoginType: database.LoginTypeOIDC,
		})
		require.NoError(t, err)
		require.Equal(t, "https://old-issuer.example.com||sub-mismatched", link.LinkedID)
	})

	t.Run("DiscoveryFailureNonFatal", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitShort)
		// oidcAuthLinks logs errors but does not return them.
		logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true}).Leveled(slog.LevelDebug)

		// Serve 500 so discovery fails.
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		t.Cleanup(srv.Close)

		db, _ := dbtestutil.NewDB(t)

		vals := newDeploymentValues(t, srv.URL, true)
		err := oidcAuthLinks(ctx, logger, srv.Client(), vals, db)
		require.NoError(t, err, "discovery failure must not be fatal")
	})

	t.Run("AnalyzeFailureNonFatal", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitShort)
		// oidcAuthLinks logs errors but does not return them.
		logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true}).Leveled(slog.LevelDebug)

		const expectedIssuer = "https://accounts.google.com"
		oidcSrv := fakeOIDCIssuerDiscovery(t, expectedIssuer)

		// Use a closed DB so the analysis query fails.
		connectionURL, err := dbtestutil.Open(t)
		require.NoError(t, err)
		sqlDB, err := sql.Open("postgres", connectionURL)
		require.NoError(t, err)
		db := database.New(sqlDB)
		// Close the underlying connection so the query fails.
		sqlDB.Close()

		vals := newDeploymentValues(t, oidcSrv.URL, true)
		err = oidcAuthLinks(ctx, logger, oidcSrv.Client(), vals, db)
		require.NoError(t, err, "analysis failure must not be fatal")
	})
}
