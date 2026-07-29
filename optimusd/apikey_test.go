package optimus-ide-collabd_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/serpent"
)

func TestTokenCRUD(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	auditor := audit.NewMock()
	numLogs := len(auditor.AuditLogs())
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
	_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
	numLogs++ // add an audit log for user creation

	keys, err := client.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
	require.NoError(t, err)
	require.Empty(t, keys)

	res, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{})
	require.NoError(t, err)
	require.Greater(t, len(res.Key), 2)
	numLogs++ // add an audit log for token creation

	keys, err = client.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
	require.NoError(t, err)
	require.EqualValues(t, len(keys), 1)
	require.Contains(t, res.Key, keys[0].ID)
	// expires_at should default to 30 days
	require.Greater(t, keys[0].ExpiresAt, dbtime.Now().Add(time.Hour*24*6))
	require.Less(t, keys[0].ExpiresAt, dbtime.Now().Add(time.Hour*24*8))
	require.Equal(t, optimus-ide-collabsdk.APIKeyScopeAll, keys[0].Scope)
	require.Len(t, keys[0].AllowList, 1)
	require.Equal(t, "*:*", keys[0].AllowList[0].String())

	// no update

	err = client.DeleteAPIKey(ctx, optimus-ide-collabsdk.Me, keys[0].ID)
	require.NoError(t, err)
	numLogs++ // add an audit log for token deletion
	keys, err = client.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
	require.NoError(t, err)
	require.Empty(t, keys)

	// ensure audit log count is correct
	require.Len(t, auditor.AuditLogs(), numLogs)
	require.Equal(t, database.AuditActionCreate, auditor.AuditLogs()[numLogs-2].Action)
	require.Equal(t, database.AuditActionDelete, auditor.AuditLogs()[numLogs-1].Action)
}

func TestTokensFilterExpired(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	adminClient := optimus-ide-collabdtest.New(t, nil)
	_ = optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

	// Create a token.
	res, err := adminClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 7,
	})
	require.NoError(t, err)
	keyID := strings.Split(res.Key, "-")[0]

	// List tokens without including expired - should see the token.
	keys, err := adminClient.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
	require.NoError(t, err)
	require.Len(t, keys, 1)

	// Expire the token.
	err = adminClient.ExpireAPIKey(ctx, optimus-ide-collabsdk.Me, keyID)
	require.NoError(t, err)

	// List tokens without including expired - should NOT see expired token.
	keys, err = adminClient.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
	require.NoError(t, err)
	require.Empty(t, keys)

	// List tokens WITH including expired - should see expired token.
	keys, err = adminClient.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{
		IncludeExpired: true,
	})
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, keyID, keys[0].ID)
}

func TestTokenScoped(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	client := optimus-ide-collabdtest.New(t, nil)
	_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

	res, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Scope: optimus-ide-collabsdk.APIKeyScopeApplicationConnect,
	})
	require.NoError(t, err)
	require.Greater(t, len(res.Key), 2)

	keys, err := client.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
	require.NoError(t, err)
	require.EqualValues(t, len(keys), 1)
	require.Contains(t, res.Key, keys[0].ID)
	require.Equal(t, keys[0].Scope, optimus-ide-collabsdk.APIKeyScopeApplicationConnect)
	require.Len(t, keys[0].AllowList, 1)
	require.Equal(t, "*:*", keys[0].AllowList[0].String())
}

// Ensure backward-compat: when a token is created using the legacy singular
// scope names ("all" or "application_connect"), the API returns the same
// legacy value in the deprecated singular Scope field while also supporting
// the new multi-scope field.
func TestTokenLegacySingularScopeCompat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		scope  optimus-ide-collabsdk.APIKeyScope
		scopes []optimus-ide-collabsdk.APIKeyScope
	}{
		{
			name:   "all",
			scope:  optimus-ide-collabsdk.APIKeyScopeAll,
			scopes: []optimus-ide-collabsdk.APIKeyScope{optimus-ide-collabsdk.APIKeyScopeOptimus-IDE-CollabAll},
		},
		{
			name:   "application_connect",
			scope:  optimus-ide-collabsdk.APIKeyScopeApplicationConnect,
			scopes: []optimus-ide-collabsdk.APIKeyScope{optimus-ide-collabsdk.APIKeyScopeOptimus-IDE-CollabApplicationConnect},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), testutil.WaitLong)
			defer cancel()
			client := optimus-ide-collabdtest.New(t, nil)
			_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

			// Create with legacy singular scope.
			_, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
				Scope: tc.scope,
			})
			require.NoError(t, err)

			// Read back and ensure the deprecated singular field matches exactly.
			keys, err := client.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
			require.NoError(t, err)
			require.Len(t, keys, 1)
			require.Equal(t, tc.scope, keys[0].Scope)
			require.ElementsMatch(t, keys[0].Scopes, tc.scopes)
			require.Len(t, keys[0].AllowList, 1)
			require.Equal(t, "*:*", keys[0].AllowList[0].String())
		})
	}
}

func TestUserSetTokenDuration(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	client := optimus-ide-collabdtest.New(t, nil)
	_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

	_, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 7,
	})
	require.NoError(t, err)
	keys, err := client.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
	require.NoError(t, err)
	require.Greater(t, keys[0].ExpiresAt, dbtime.Now().Add(time.Hour*6*24))
	require.Less(t, keys[0].ExpiresAt, dbtime.Now().Add(time.Hour*8*24))
}

func TestDefaultTokenDuration(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	client := optimus-ide-collabdtest.New(t, nil)
	_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

	_, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{})
	require.NoError(t, err)
	keys, err := client.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
	require.NoError(t, err)
	require.Greater(t, keys[0].ExpiresAt, dbtime.Now().Add(time.Hour*24*6))
	require.Less(t, keys[0].ExpiresAt, dbtime.Now().Add(time.Hour*24*8))
}

func TestTokenUserSetMaxLifetime(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	dc := optimus-ide-collabdtest.DeploymentValues(t)
	dc.Sessions.MaximumTokenDuration = serpent.Duration(time.Hour * 24 * 7)
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		DeploymentValues: dc,
	})
	_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

	// success
	_, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 6,
	})
	require.NoError(t, err)

	// fail
	_, err = client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 8,
	})
	require.ErrorContains(t, err, "lifetime must be less")
}

func TestTokenAdminSetMaxLifetime(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	dc := optimus-ide-collabdtest.DeploymentValues(t)
	dc.Sessions.MaximumTokenDuration = serpent.Duration(time.Hour * 24 * 7)
	dc.Sessions.MaximumAdminTokenDuration = serpent.Duration(time.Hour * 24 * 14)
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		DeploymentValues: dc,
	})
	adminUser := optimus-ide-collabdtest.CreateFirstUser(t, client)
	nonAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, adminUser.OrganizationID)

	// Admin should be able to create a token with a lifetime longer than the non-admin max.
	_, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 10,
	})
	require.NoError(t, err)

	// Admin should NOT be able to create a token with a lifetime longer than the admin max.
	_, err = client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 15,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "lifetime must be less")

	// Non-admin should NOT be able to create a token with a lifetime longer than the non-admin max.
	_, err = nonAdminClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 8,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "lifetime must be less")

	// Non-admin should be able to create a token with a lifetime shorter than the non-admin max.
	_, err = nonAdminClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 6,
	})
	require.NoError(t, err)
}

func TestTokenAdminSetMaxLifetimeShorter(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	dc := optimus-ide-collabdtest.DeploymentValues(t)
	dc.Sessions.MaximumTokenDuration = serpent.Duration(time.Hour * 24 * 14)
	dc.Sessions.MaximumAdminTokenDuration = serpent.Duration(time.Hour * 24 * 7)
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		DeploymentValues: dc,
	})
	adminUser := optimus-ide-collabdtest.CreateFirstUser(t, client)
	nonAdminClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, adminUser.OrganizationID)

	// Admin should NOT be able to create a token with a lifetime longer than the admin max.
	_, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 8,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "lifetime must be less")

	// Admin should be able to create a token with a lifetime shorter than the admin max.
	_, err = client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 6,
	})
	require.NoError(t, err)

	// Non-admin should be able to create a token with a lifetime longer than the admin max.
	_, err = nonAdminClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 10,
	})
	require.NoError(t, err)

	// Non-admin should NOT be able to create a token with a lifetime longer than the non-admin max.
	_, err = nonAdminClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
		Lifetime: time.Hour * 24 * 15,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "lifetime must be less")
}

func TestTokenCustomDefaultLifetime(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	dc := optimus-ide-collabdtest.DeploymentValues(t)
	dc.Sessions.DefaultTokenDuration = serpent.Duration(time.Hour * 12)
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		DeploymentValues: dc,
	})
	_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

	_, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{})
	require.NoError(t, err)

	tokens, err := client.Tokens(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.TokensFilter{})
	require.NoError(t, err)
	require.Len(t, tokens, 1)
	require.EqualValues(t, dc.Sessions.DefaultTokenDuration.Value().Seconds(), tokens[0].LifetimeSeconds)
}

func TestSessionExpiry(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	dc := optimus-ide-collabdtest.DeploymentValues(t)

	db, pubsub := dbtestutil.NewDB(t)
	adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		DeploymentValues: dc,
		Database:         db,
		Pubsub:           pubsub,
	})
	adminUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

	// This is a hack, but we need the admin account to have a long expiry
	// otherwise the test will flake, so we only update the expiry config after
	// the admin account has been created.
	//
	// We don't support updating the deployment config after startup, but for
	// this test it works because we don't copy the value (and we use pointers).
	dc.Sessions.DefaultDuration = serpent.Duration(time.Second)

	userClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)

	// Find the session cookie, and ensure it has the correct expiry.
	token := userClient.SessionToken()
	apiKey, err := db.GetAPIKeyByID(ctx, strings.Split(token, "-")[0])
	require.NoError(t, err)

	require.EqualValues(t, dc.Sessions.DefaultDuration.Value().Seconds(), apiKey.LifetimeSeconds)
	require.WithinDuration(t, apiKey.CreatedAt.Add(dc.Sessions.DefaultDuration.Value()), apiKey.ExpiresAt, 2*time.Second)

	// Update the session token to be expired so we can test that it is
	// rejected for extra points.
	err = db.UpdateAPIKeyByID(ctx, database.UpdateAPIKeyByIDParams{
		ID:        apiKey.ID,
		LastUsed:  apiKey.LastUsed,
		ExpiresAt: dbtime.Now().Add(-time.Hour),
		IPAddress: apiKey.IPAddress,
	})
	require.NoError(t, err)

	_, err = userClient.User(ctx, optimus-ide-collabsdk.Me)
	require.Error(t, err)
	var sdkErr *optimus-ide-collabsdk.Error
	if assert.ErrorAs(t, err, &sdkErr) {
		require.Equal(t, http.StatusUnauthorized, sdkErr.StatusCode())
		require.Contains(t, sdkErr.Message, "session has expired")
	}
}

// TestSessionCookieMaxAge verifies that the session cookie is a persistent
// cookie (has MaxAge set) rather than a session cookie. Standalone PWAs
// run in their own browser process and mobile OSes purge in-memory
// (session) cookies when that process is killed, so the cookie must be
// persisted to disk.
func TestSessionCookieMaxAge(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	client := optimus-ide-collabdtest.New(t, nil)

	// Create the first user (password-based login).
	req := optimus-ide-collabsdk.CreateFirstUserRequest{
		Email:    "testuser@optimus-ide-collab.com",
		Username: "testuser",
		Password: "SomeSecurePassword!",
	}
	_, err := client.CreateFirstUser(ctx, req)
	require.NoError(t, err)

	// Login via the raw HTTP endpoint so we can inspect the Set-Cookie header.
	loginURL, err := client.URL.Parse("/api/v2/users/login")
	require.NoError(t, err)

	res, err := client.Request(ctx, http.MethodPost, loginURL.String(), optimus-ide-collabsdk.LoginWithPasswordRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusCreated, res.StatusCode)

	oneYear := int((365 * 24 * time.Hour).Seconds())
	var found bool
	for _, cookie := range res.Cookies() {
		if cookie.Name == optimus-ide-collabsdk.SessionTokenCookie {
			// MaxAge should be set to a long value so the browser
			// persists the cookie to disk. The server handles real
			// expiry via the API key's ExpiresAt field.
			require.Equal(t, oneYear, cookie.MaxAge,
				"Session cookie MaxAge should be set to 1 year for disk persistence")
			found = true
		}
	}
	require.True(t, found, "session cookie should be present in login response")
}

func TestAPIKey_OK(t *testing.T) {
	t.Parallel()

	// Given: a deployment with auditing enabled
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	auditor := audit.NewMock()
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
	owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
	auditor.ResetLogs()

	// When: an API key is created
	res, err := client.CreateAPIKey(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)
	require.Greater(t, len(res.Key), 2)

	// Then: an audit log is generated
	als := auditor.AuditLogs()
	require.Len(t, als, 1)
	al := als[0]
	assert.Equal(t, owner.UserID, al.UserID)
	assert.Equal(t, database.AuditActionCreate, al.Action)
	assert.Equal(t, database.ResourceTypeApiKey, al.ResourceType)

	// Then: the diff MUST NOT contain the generated key.
	raw, err := json.Marshal(al)
	require.NoError(t, err)
	require.NotContains(t, res.Key, string(raw))
}

func TestAPIKey_Deleted(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
	user := optimus-ide-collabdtest.CreateFirstUser(t, client)
	_, anotherUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, user.OrganizationID)
	require.NoError(t, client.DeleteUser(context.Background(), anotherUser.ID))

	// Attempt to create an API key for the deleted user. This should fail.
	_, err := client.CreateAPIKey(ctx, anotherUser.Username)
	require.Error(t, err)
	var apiErr *optimus-ide-collabsdk.Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
}

func TestAPIKey_SetDefault(t *testing.T) {
	t.Parallel()

	db, pubsub := dbtestutil.NewDB(t)
	dc := optimus-ide-collabdtest.DeploymentValues(t)
	dc.Sessions.DefaultTokenDuration = serpent.Duration(time.Hour * 12)
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		Database:         db,
		Pubsub:           pubsub,
		DeploymentValues: dc,
	})
	owner := optimus-ide-collabdtest.CreateFirstUser(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	token, err := client.CreateAPIKey(ctx, owner.UserID.String())
	require.NoError(t, err)
	split := strings.Split(token.Key, "-")
	apiKey1, err := db.GetAPIKeyByID(ctx, split[0])
	require.NoError(t, err)
	require.EqualValues(t, dc.Sessions.DefaultTokenDuration.Value().Seconds(), apiKey1.LifetimeSeconds)
}

func TestAPIKey_PrebuildsNotAllowed(t *testing.T) {
	t.Parallel()

	db, pubsub := dbtestutil.NewDB(t)
	dc := optimus-ide-collabdtest.DeploymentValues(t)
	dc.Sessions.DefaultTokenDuration = serpent.Duration(time.Hour * 12)
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		Database:         db,
		Pubsub:           pubsub,
		DeploymentValues: dc,
	})

	setupCtx := testutil.Context(t, testutil.WaitLong)

	// Given: an existing api token for the prebuilds user
	_, prebuildsToken := dbgen.APIKey(t, db, database.APIKey{
		UserID: database.PrebuildsSystemUserID,
	})
	client.SetSessionToken(prebuildsToken)

	// When: the prebuilds user tries to create an API key
	_, err := client.CreateAPIKey(setupCtx, database.PrebuildsSystemUserID.String())
	// Then: denied.
	require.ErrorContains(t, err, httpapi.ResourceForbiddenResponse.Message)

	// When: the prebuilds user tries to create a token
	_, err = client.CreateToken(setupCtx, database.PrebuildsSystemUserID.String(), optimus-ide-collabsdk.CreateTokenRequest{})
	// Then: also denied.
	require.ErrorContains(t, err, httpapi.ResourceForbiddenResponse.Message)
}

//nolint:tparallel,paralleltest // Subtests share the same optimus-ide-collabdtest instance and auditor.
func TestExpireAPIKey(t *testing.T) {
	t.Parallel()

	auditor := audit.NewMock()
	adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
	admin := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)
	memberClient, member := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, admin.OrganizationID)

	t.Run("OwnerCanExpireOwnToken", func(t *testing.T) {
		ctx := testutil.Context(t, testutil.WaitLong)

		// Create a token.
		res, err := adminClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
			Lifetime: time.Hour * 24 * 7,
		})
		require.NoError(t, err)
		keyID := strings.Split(res.Key, "-")[0]

		// Verify the token is not expired.
		key, err := adminClient.APIKeyByID(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)
		require.True(t, key.ExpiresAt.After(dbtime.Now()))

		auditor.ResetLogs()

		// Expire the token.
		err = adminClient.ExpireAPIKey(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)

		// Verify the token is expired.
		key, err = adminClient.APIKeyByID(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)
		require.True(t, key.ExpiresAt.Before(dbtime.Now()))

		// Verify audit log.
		als := auditor.AuditLogs()
		require.Len(t, als, 1)
		require.Equal(t, database.AuditActionWrite, als[0].Action)
		require.Equal(t, database.ResourceTypeApiKey, als[0].ResourceType)
		require.Equal(t, admin.UserID.String(), als[0].UserID.String())
	})

	t.Run("AdminCanExpireOtherUsersToken", func(t *testing.T) {
		ctx := testutil.Context(t, testutil.WaitLong)

		// Create a token for the member.
		res, err := memberClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
			Lifetime: time.Hour * 24 * 7,
		})
		require.NoError(t, err)
		keyID := strings.Split(res.Key, "-")[0]

		// Admin expires the member's token.
		err = adminClient.ExpireAPIKey(ctx, member.ID.String(), keyID)
		require.NoError(t, err)

		// Verify the token is expired.
		key, err := memberClient.APIKeyByID(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)
		require.True(t, key.ExpiresAt.Before(dbtime.Now()))
	})

	t.Run("MemberCannotExpireOtherUsersToken", func(t *testing.T) {
		ctx := testutil.Context(t, testutil.WaitLong)

		// Create a token for the admin.
		res, err := adminClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
			Lifetime: time.Hour * 24 * 7,
		})
		require.NoError(t, err)
		keyID := strings.Split(res.Key, "-")[0]

		// Member attempts to expire admin's token.
		err = memberClient.ExpireAPIKey(ctx, admin.UserID.String(), keyID)
		require.Error(t, err)
		var sdkErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &sdkErr)
		// Members cannot read other users, so they get a 404 Not Found
		// from the authorization layer.
		require.Equal(t, http.StatusNotFound, sdkErr.StatusCode())
	})

	t.Run("NotFound", func(t *testing.T) {
		ctx := testutil.Context(t, testutil.WaitLong)

		// Try to expire a non-existent token.
		err := adminClient.ExpireAPIKey(ctx, optimus-ide-collabsdk.Me, "nonexistent")
		require.Error(t, err)
		var sdkErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &sdkErr)
		require.Equal(t, http.StatusNotFound, sdkErr.StatusCode())
	})

	t.Run("ExpiringAlreadyExpiredTokenSucceeds", func(t *testing.T) {
		ctx := testutil.Context(t, testutil.WaitLong)

		// Create and expire a token.
		res, err := adminClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
			Lifetime: time.Hour * 24 * 7,
		})
		require.NoError(t, err)
		keyID := strings.Split(res.Key, "-")[0]

		// Expire it once.
		err = adminClient.ExpireAPIKey(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)

		// Invariant: make sure it's actually expired
		key, err := adminClient.APIKeyByID(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)
		require.LessOrEqual(t, key.ExpiresAt, dbtime.Now(), "key should be expired")

		// Expire it again - should succeed (idempotent).
		err = adminClient.ExpireAPIKey(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)

		// Token should still be just as expired as before. No more, no less.
		keyAgain, err := adminClient.APIKeyByID(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)
		require.Equal(t, key.ExpiresAt, keyAgain.ExpiresAt, "expiration should be idempotent")
	})

	t.Run("DeletingExpiredTokenSucceeds", func(t *testing.T) {
		ctx := testutil.Context(t, testutil.WaitLong)

		// Create a token.
		res, err := adminClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{
			Lifetime: time.Hour * 24 * 7,
		})
		require.NoError(t, err)
		keyID := strings.Split(res.Key, "-")[0]

		// Expire it first.
		err = adminClient.ExpireAPIKey(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)

		// Verify it's expired.
		key, err := adminClient.APIKeyByID(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)
		require.True(t, key.ExpiresAt.Before(dbtime.Now()))

		// Delete the expired token - should succeed.
		err = adminClient.DeleteAPIKey(ctx, optimus-ide-collabsdk.Me, keyID)
		require.NoError(t, err)

		// Verify it's gone.
		_, err = adminClient.APIKeyByID(ctx, optimus-ide-collabsdk.Me, keyID)
		require.Error(t, err)
		var sdkErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &sdkErr)
		require.Equal(t, http.StatusNotFound, sdkErr.StatusCode())
	})
}
