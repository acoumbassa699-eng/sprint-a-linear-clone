package optimus-ide-collabd_test

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest/oidctest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications/notificationstest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/serpent"
)

func TestFirstUser(t *testing.T) {
	t.Parallel()
	t.Run("BadRequest", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		has, err := client.HasFirstUser(context.Background())
		require.NoError(t, err)
		require.False(t, has)

		_, err = client.CreateFirstUser(ctx, optimus-ide-collabsdk.CreateFirstUserRequest{})
		require.Error(t, err)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.CreateFirstUser(ctx, optimus-ide-collabsdk.CreateFirstUserRequest{
			Email:    "some@email.com",
			Username: "exampleuser",
			Password: "SomeSecurePassword!",
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusConflict, apiErr.StatusCode())
	})

	t.Run("Create", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitShort)
		client := optimus-ide-collabdtest.New(t, nil)
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
		u, err := client.User(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		assert.Equal(t, optimus-ide-collabdtest.FirstUserParams.Name, u.Name)
		assert.Equal(t, optimus-ide-collabdtest.FirstUserParams.Email, u.Email)
		assert.Equal(t, optimus-ide-collabdtest.FirstUserParams.Username, u.Username)
	})

	t.Run("Trial", func(t *testing.T) {
		t.Parallel()
		trialGenerated := make(chan struct{})
		entitlementsRefreshed := make(chan struct{})

		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			TrialGenerator: func(context.Context, optimus-ide-collabsdk.LicensorTrialRequest) error {
				close(trialGenerated)
				return nil
			},
			RefreshEntitlements: func(context.Context) error {
				close(entitlementsRefreshed)
				return nil
			},
		})

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		req := optimus-ide-collabsdk.CreateFirstUserRequest{
			Email:    "testuser@optimus-ide-collab.com",
			Username: "testuser",
			Name:     "Test User",
			Password: "SomeSecurePassword!",
			Trial:    true,
		}
		_, err := client.CreateFirstUser(ctx, req)
		require.NoError(t, err)

		_ = testutil.TryReceive(ctx, t, trialGenerated)
		_ = testutil.TryReceive(ctx, t, entitlementsRefreshed)
	})
}

func TestFirstUser_OnboardingTelemetry(t *testing.T) {
	t.Parallel()

	t.Run("OnboardingInfoFlowsToSnapshot", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitMedium)
		fTelemetry := newFakeTelemetryReporter(ctx, t, 10)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			TelemetryReporter: fTelemetry,
		})

		_, err := client.CreateFirstUser(ctx, optimus-ide-collabsdk.CreateFirstUserRequest{
			Email:    "admin@optimus-ide-collab.com",
			Username: "admin",
			Password: "SomeSecurePassword!",
			OnboardingInfo: &optimus-ide-collabsdk.CreateFirstUserOnboardingInfo{
				NewsletterMarketing: false,
				NewsletterReleases:  true,
			},
		})
		require.NoError(t, err)

		snapshot := testutil.TryReceive(ctx, t, fTelemetry.snapshots)
		require.NotNil(t, snapshot.FirstUserOnboarding)
		require.False(t, snapshot.FirstUserOnboarding.NewsletterMarketing)
		require.True(t, snapshot.FirstUserOnboarding.NewsletterReleases)
	})

	t.Run("NilWhenOnboardingInfoOmitted", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitMedium)
		fTelemetry := newFakeTelemetryReporter(ctx, t, 10)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			TelemetryReporter: fTelemetry,
		})

		_, err := client.CreateFirstUser(ctx, optimus-ide-collabsdk.CreateFirstUserRequest{
			Email:    "admin@optimus-ide-collab.com",
			Username: "admin",
			Password: "SomeSecurePassword!",
			// No OnboardingInfo — simulates old CLI or OIDC flow.
		})
		require.NoError(t, err)

		snapshot := testutil.TryReceive(ctx, t, fTelemetry.snapshots)
		require.Nil(t, snapshot.FirstUserOnboarding)
	})

	t.Run("EmptyOnboardingInfoIsNonNilWithZeroFields", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		fTelemetry := newFakeTelemetryReporter(ctx, t, 10)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			TelemetryReporter: fTelemetry,
		})
		_, err := client.CreateFirstUser(ctx, optimus-ide-collabsdk.CreateFirstUserRequest{
			Email: "admin@optimus-ide-collab.com", Username: "admin",
			Password:       "SomeSecurePassword!",
			OnboardingInfo: &optimus-ide-collabsdk.CreateFirstUserOnboardingInfo{},
		})
		require.NoError(t, err)
		snapshot := testutil.TryReceive(ctx, t, fTelemetry.snapshots)
		require.NotNil(t, snapshot.FirstUserOnboarding,
			"non-nil OnboardingInfo must produce non-nil telemetry")
		require.False(t, snapshot.FirstUserOnboarding.NewsletterMarketing)
		require.False(t, snapshot.FirstUserOnboarding.NewsletterReleases)
	})
}

func TestPostLogin(t *testing.T) {
	t.Parallel()
	t.Run("InvalidUser", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    "my@email.org",
			Password: "password",
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode())
	})

	t.Run("BadPassword", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		req := optimus-ide-collabsdk.CreateFirstUserRequest{
			Email:    "testuser@optimus-ide-collab.com",
			Username: "testuser",
			Password: "SomeSecurePassword!",
		}
		_, err := client.CreateFirstUser(ctx, req)
		require.NoError(t, err)
		_, err = client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    req.Email,
			Password: "badpass",
		})
		numLogs++ // add an audit log for login
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode())

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionLogin, auditor.AuditLogs()[numLogs-1].Action)
	})

	// "hunter2" was the input of the previous hardcoded simulated hash, which
	// an empty stored hash wrongly matched; this is a regression test.
	t.Run("NonexistentUser401", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    "does-not-exist@optimus-ide-collab.com",
			Password: "hunter2",
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode())
		require.Equal(t, "Incorrect email or password.", apiErr.Message)
	})

	// Attempting built-in login as an SSO user returns a 401 to avoid
	// divulging login type.
	t.Run("SSOReturns401", func(t *testing.T) {
		t.Parallel()
		client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		// An SSO user has no password hash stored. Create one directly in the
		// database since the API requires OIDC to be configured. dbgen.User
		// substitutes a random hash for an empty one, so clear it explicitly.
		ssoUser := dbgen.User(t, db, database.User{
			Email:     "sso-user@optimus-ide-collab.com",
			LoginType: database.LoginTypeOIDC,
		})
		//nolint:gocritic // Test setup requires a system context to clear the hash.
		err := db.UpdateUserHashedPassword(dbauthz.AsSystemRestricted(ctx), database.UpdateUserHashedPasswordParams{
			ID:             ssoUser.ID,
			HashedPassword: []byte{},
		})
		require.NoError(t, err)

		anonClient := optimus-ide-collabsdk.New(client.URL)
		_, err = anonClient.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    ssoUser.Email,
			Password: "hunter2",
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode())
		require.Equal(t, "Incorrect email or password.", apiErr.Message)
		// The login type must not be leaked.
		require.NotContains(t, apiErr.Message, string(optimus-ide-collabsdk.LoginTypeOIDC))
	})

	// Regression: the legacy `login_type = 'none'` migration converts these
	// accounts to password auth, but they have no password hash. Converting
	// login type must never let someone authenticate with an empty or guessed
	// password.
	t.Run("ConvertedNoneUserHasNoUsablePassword", func(t *testing.T) {
		t.Parallel()
		client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		// A legacy machine user was created with login_type 'none' and no
		// password. dbgen.User substitutes a random hash for an empty one, so
		// clear it explicitly to match the real account.
		noneUser := dbgen.User(t, db, database.User{
			Email:     "legacy-machine-user@optimus-ide-collab.com",
			LoginType: database.LoginTypeNone,
		})
		//nolint:gocritic // Test setup requires a system context to clear the hash.
		err := db.UpdateUserHashedPassword(dbauthz.AsSystemRestricted(ctx), database.UpdateUserHashedPasswordParams{
			ID:             noneUser.ID,
			HashedPassword: []byte{},
		})
		require.NoError(t, err)

		// Apply the migration's conversion: login_type 'none' -> 'password'.
		//nolint:gocritic // Test setup requires a system context to convert the login type.
		_, err = db.UpdateUserLoginType(dbauthz.AsSystemRestricted(ctx), database.UpdateUserLoginTypeParams{
			NewLoginType: database.LoginTypePassword,
			UserID:       noneUser.ID,
		})
		require.NoError(t, err)

		// Neither an empty password nor a guessed one may authenticate. An empty
		// password is rejected by request validation (400); a non-empty guess
		// fails the hash comparison against the empty stored hash (401). Both must
		// deny access.
		cases := []struct {
			name       string
			password   string
			wantStatus int
		}{
			{"EmptyPassword", "", http.StatusBadRequest},
			{"GuessedPassword", "hunter2", http.StatusUnauthorized},
		}
		for _, tc := range cases {
			anonClient := optimus-ide-collabsdk.New(client.URL)
			_, err := anonClient.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
				Email:    noneUser.Email,
				Password: tc.password,
			})
			var apiErr *optimus-ide-collabsdk.Error
			require.ErrorAs(t, err, &apiErr, "%s must not authenticate", tc.name)
			require.Equal(t, tc.wantStatus, apiErr.StatusCode(), "%s", tc.name)
		}
	})

	t.Run("Suspended", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())
		first := optimus-ide-collabdtest.CreateFirstUser(t, client)
		numLogs++ // add an audit log for create user
		numLogs++ // add an audit log for login

		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID)
		numLogs++ // add an audit log for create user

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		memberUser, err := member.User(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err, "fetch member user")

		_, err = client.UpdateUserStatus(ctx, memberUser.Username, optimus-ide-collabsdk.UserStatusSuspended)
		require.NoError(t, err, "suspend member")
		numLogs++ // add an audit log for update user

		// Test an existing session
		_, err = member.User(ctx, optimus-ide-collabsdk.Me)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode())
		require.Contains(t, apiErr.Message, "Contact an admin")

		// Test a new session
		_, err = client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    memberUser.Email,
			Password: "SomeSecurePassword!",
		})
		numLogs++ // add an audit log for login
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode())
		require.Contains(t, apiErr.Message, "suspended")

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionLogin, auditor.AuditLogs()[numLogs-1].Action)
	})

	t.Run("DisabledPasswordAuth", func(t *testing.T) {
		t.Parallel()

		dc := optimus-ide-collabdtest.DeploymentValues(t)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues: dc,
		})

		first := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		// With a user account.
		const password = "SomeSecurePassword!"
		user, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           "test+user-@optimus-ide-collab.com",
			Username:        "user",
			Password:        password,
			OrganizationIDs: []uuid.UUID{first.OrganizationID},
		})
		require.NoError(t, err)

		dc.DisablePasswordAuth = serpent.Bool(true)

		userClient := optimus-ide-collabsdk.New(client.URL)
		_, err = userClient.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    user.Email,
			Password: password,
		})
		require.Error(t, err)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusForbidden, apiErr.StatusCode())
		require.Contains(t, apiErr.Message, "Password authentication is disabled")
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		req := optimus-ide-collabsdk.CreateFirstUserRequest{
			Email:    "testuser@optimus-ide-collab.com",
			Username: "testuser",
			Password: "SomeSecurePassword!",
		}
		_, err := client.CreateFirstUser(ctx, req)
		require.NoError(t, err)
		numLogs++ // add an audit log for create user
		numLogs++ // add an audit log for login

		_, err = client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		require.NoError(t, err)

		// Login should be case insensitive
		_, err = client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    strings.ToUpper(req.Email),
			Password: req.Password,
		})
		require.NoError(t, err)

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionLogin, auditor.AuditLogs()[numLogs-1].Action)
	})

	t.Run("Lifetime&Expire", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		split := strings.Split(client.SessionToken(), "-")
		key, err := client.APIKeyByID(ctx, owner.UserID.String(), split[0])
		require.NoError(t, err, "fetch login key")
		require.Equal(t, int64(86400), key.LifetimeSeconds, "default should be 86400")

		// tokens have a longer life
		token, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{})
		require.NoError(t, err, "make new token api key")
		split = strings.Split(token.Key, "-")
		apiKey, err := client.APIKeyByID(ctx, owner.UserID.String(), split[0])
		require.NoError(t, err, "fetch api key")

		require.True(t, apiKey.ExpiresAt.After(dbtime.Now().Add(time.Hour*24*6)), "default tokens lasts more than 6 days")
		require.True(t, apiKey.ExpiresAt.Before(dbtime.Now().Add(time.Hour*24*8)), "default tokens lasts less than 8 days")
		require.Greater(t, apiKey.LifetimeSeconds, key.LifetimeSeconds, "token should have longer lifetime")
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	t.Run("Works", func(t *testing.T) {
		t.Parallel()
		client, _, api := optimus-ide-collabdtest.NewWithAPI(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		authz := optimus-ide-collabdtest.AssertRBAC(t, api, client)

		anotherClient, another := optimus-ide-collabdtest.CreateAnotherUser(t, client, user.OrganizationID)
		err := client.DeleteUser(context.Background(), another.ID)
		require.NoError(t, err)
		// Attempt to create a user with the same email and username, and delete them again.
		another, err = client.CreateUserWithOrgs(context.Background(), optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           another.Email,
			Username:        another.Username,
			Password:        "SomeSecurePassword!",
			OrganizationIDs: []uuid.UUID{user.OrganizationID},
		})
		require.NoError(t, err)
		err = client.DeleteUser(context.Background(), another.ID)
		require.NoError(t, err)

		// IMPORTANT: assert that the deleted user's session is no longer valid.
		_, err = anotherClient.User(context.Background(), optimus-ide-collabsdk.Me)
		require.Error(t, err)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode())

		// RBAC checks
		authz.AssertChecked(t, policy.ActionCreate, rbac.ResourceUser)
		authz.AssertChecked(t, policy.ActionDelete, another)
	})
	t.Run("NoPermission", func(t *testing.T) {
		t.Parallel()
		api := optimus-ide-collabdtest.New(t, nil)
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, api)
		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, api, firstUser.OrganizationID)
		err := client.DeleteUser(context.Background(), firstUser.UserID)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})
	t.Run("HasWorkspaces", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdtest.NewWithProvisionerCloser(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		anotherClient, another := optimus-ide-collabdtest.CreateAnotherUser(t, client, user.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		optimus-ide-collabdtest.CreateWorkspace(t, anotherClient, template.ID)
		err := client.DeleteUser(context.Background(), another.ID)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusExpectationFailed, apiErr.StatusCode())
	})
	t.Run("Self", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		err := client.DeleteUser(context.Background(), user.UserID)
		var apiErr *optimus-ide-collabsdk.Error
		require.Error(t, err, "should not be able to delete self")
		require.ErrorAs(t, err, &apiErr, "should be a optimus-ide-collabd error")
		require.Equal(t, http.StatusForbidden, apiErr.StatusCode(), "should be forbidden")
	})
	t.Run("CountCheckIncludesAllWorkspaces", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdtest.NewWithProvisionerCloser(t, nil)
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client)

		// Create a target user who will own a workspace
		targetUserClient, targetUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, firstUser.OrganizationID)

		// Create a User Admin who should not have permission to see the target user's workspace
		userAdminClient, userAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, client, firstUser.OrganizationID)

		// Grant User Admin role to the userAdmin
		userAdmin, err := client.UpdateUserRoles(context.Background(), userAdmin.ID.String(), optimus-ide-collabsdk.UpdateRoles{
			Roles: []string{rbac.RoleUserAdmin().String()},
		})
		require.NoError(t, err)

		// Create a template and workspace owned by the target user
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, firstUser.OrganizationID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, firstUser.OrganizationID, version.ID)
		_ = optimus-ide-collabdtest.CreateWorkspace(t, targetUserClient, template.ID)

		workspaces, err := userAdminClient.Workspaces(context.Background(), optimus-ide-collabsdk.WorkspaceFilter{
			Owner: targetUser.Username,
		})
		require.NoError(t, err)
		require.Len(t, workspaces.Workspaces, 0)

		// Attempt to delete the target user - this should fail because the
		// user has a workspace not visible to the deleting user.
		err = userAdminClient.DeleteUser(context.Background(), targetUser.ID)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusExpectationFailed, apiErr.StatusCode())
		require.Contains(t, apiErr.Message, "has workspaces")
	})
}

func TestNotifyUserStatusChanged(t *testing.T) {
	t.Parallel()

	type expectedNotification struct {
		TemplateID uuid.UUID
		UserID     uuid.UUID
	}

	verifyNotificationDispatched := func(notifyEnq *notificationstest.FakeEnqueuer, expectedNotifications []expectedNotification, member optimus-ide-collabsdk.User, label string) {
		require.Equal(t, len(expectedNotifications), len(notifyEnq.Sent()))

		// Validate that each expected notification is present in notifyEnq.Sent()
		for _, expected := range expectedNotifications {
			found := false
			for _, sent := range notifyEnq.Sent(notificationstest.WithTemplateID(expected.TemplateID)) {
				if sent.TemplateID == expected.TemplateID &&
					sent.UserID == expected.UserID &&
					slices.Contains(sent.Targets, member.ID) &&
					sent.Labels[label] == member.Username {
					found = true

					require.IsType(t, map[string]any{}, sent.Data["user"])
					userData := sent.Data["user"].(map[string]any)
					require.Equal(t, member.ID, userData["id"])
					require.Equal(t, member.Name, userData["name"])
					require.Equal(t, member.Email, userData["email"])

					break
				}
			}
			require.True(t, found, "Expected notification not found: %+v", expected)
		}
	}

	t.Run("Account suspended", func(t *testing.T) {
		t.Parallel()

		notifyEnq := &notificationstest.FakeEnqueuer{}
		adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			NotificationsEnqueuer: notifyEnq,
		})
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, userAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID, rbac.RoleUserAdmin())

		member, err := adminClient.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{firstUser.OrganizationID},
			Email:           "another@user.org",
			Username:        "someone-else",
			Password:        "SomeSecurePassword!",
		})
		require.NoError(t, err)

		notifyEnq.Clear()

		// when
		_, err = adminClient.UpdateUserStatus(context.Background(), member.Username, optimus-ide-collabsdk.UserStatusSuspended)
		require.NoError(t, err)

		// then
		verifyNotificationDispatched(notifyEnq, []expectedNotification{
			{TemplateID: notifications.TemplateUserAccountSuspended, UserID: firstUser.UserID},
			{TemplateID: notifications.TemplateUserAccountSuspended, UserID: userAdmin.ID},
			{TemplateID: notifications.TemplateYourAccountSuspended, UserID: member.ID},
		}, member, "suspended_account_name")
	})

	t.Run("Account reactivated", func(t *testing.T) {
		t.Parallel()

		// given
		notifyEnq := &notificationstest.FakeEnqueuer{}
		adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			NotificationsEnqueuer: notifyEnq,
		})
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, userAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID, rbac.RoleUserAdmin())

		member, err := adminClient.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{firstUser.OrganizationID},
			Email:           "another@user.org",
			Username:        "someone-else",
			Password:        "SomeSecurePassword!",
		})
		require.NoError(t, err)

		_, err = adminClient.UpdateUserStatus(context.Background(), member.Username, optimus-ide-collabsdk.UserStatusSuspended)
		require.NoError(t, err)

		notifyEnq.Clear()

		// when
		_, err = adminClient.UpdateUserStatus(context.Background(), member.Username, optimus-ide-collabsdk.UserStatusActive)
		require.NoError(t, err)

		// then
		verifyNotificationDispatched(notifyEnq, []expectedNotification{
			{TemplateID: notifications.TemplateUserAccountActivated, UserID: firstUser.UserID},
			{TemplateID: notifications.TemplateUserAccountActivated, UserID: userAdmin.ID},
			{TemplateID: notifications.TemplateYourAccountActivated, UserID: member.ID},
		}, member, "activated_account_name")
	})
}

func TestNotifyDeletedUser(t *testing.T) {
	t.Parallel()

	t.Run("OwnerNotified", func(t *testing.T) {
		t.Parallel()

		// given
		notifyEnq := &notificationstest.FakeEnqueuer{}
		adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			NotificationsEnqueuer: notifyEnq,
		})
		firstUserResponse := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		firstUser, err := adminClient.User(ctx, firstUserResponse.UserID.String())
		require.NoError(t, err)

		user, err := adminClient.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{firstUserResponse.OrganizationID},
			Email:           "another@user.org",
			Username:        "someone-else",
			Password:        "SomeSecurePassword!",
		})
		require.NoError(t, err)

		// when
		err = adminClient.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)

		// then
		require.Len(t, notifyEnq.Sent(), 2)
		// notifyEnq.Sent()[0] is create account event
		require.Equal(t, notifications.TemplateUserAccountDeleted, notifyEnq.Sent()[1].TemplateID)
		require.Equal(t, firstUser.ID, notifyEnq.Sent()[1].UserID)
		require.Contains(t, notifyEnq.Sent()[1].Targets, user.ID)
		require.Equal(t, user.Username, notifyEnq.Sent()[1].Labels["deleted_account_name"])
		require.Equal(t, user.Name, notifyEnq.Sent()[1].Labels["deleted_account_user_name"])
		require.Equal(t, firstUser.Name, notifyEnq.Sent()[1].Labels["initiator"])
	})

	t.Run("UserAdminNotified", func(t *testing.T) {
		t.Parallel()

		// given
		notifyEnq := &notificationstest.FakeEnqueuer{}
		adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			NotificationsEnqueuer: notifyEnq,
		})
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, userAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID, rbac.RoleUserAdmin())

		member, err := adminClient.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{firstUser.OrganizationID},
			Email:           "another@user.org",
			Username:        "someone-else",
			Password:        "SomeSecurePassword!",
		})
		require.NoError(t, err)

		// when
		err = adminClient.DeleteUser(context.Background(), member.ID)
		require.NoError(t, err)

		// then
		sent := notifyEnq.Sent()
		require.Len(t, sent, 5)
		// Other notifications:
		// "User admin" account created, "owner" notified
		// "Member" account created, "owner" notified
		// "Member" account created, "user admin" notified

		// "Member" account deleted, "owner" notified
		ownerNotifications := notifyEnq.Sent(func(n *notificationstest.FakeNotification) bool {
			return n.TemplateID == notifications.TemplateUserAccountDeleted &&
				n.UserID == firstUser.UserID &&
				slices.Contains(n.Targets, member.ID) &&
				n.Labels["deleted_account_name"] == member.Username
		})
		require.Len(t, ownerNotifications, 1)

		// "Member" account deleted, "user admin" notified
		adminNotifications := notifyEnq.Sent(func(n *notificationstest.FakeNotification) bool {
			return n.TemplateID == notifications.TemplateUserAccountDeleted &&
				n.UserID == userAdmin.ID &&
				slices.Contains(n.Targets, member.ID) &&
				n.Labels["deleted_account_name"] == member.Username
		})
		require.Len(t, adminNotifications, 1)
	})
}

func TestPostLogout(t *testing.T) {
	t.Parallel()

	// Checks that the cookie is cleared and the API Key is deleted from the database.
	t.Run("Logout", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		numLogs++ // add an audit log for login

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		keyID := strings.Split(client.SessionToken(), "-")[0]
		apiKey, err := client.APIKeyByID(ctx, owner.UserID.String(), keyID)
		require.NoError(t, err)
		require.Equal(t, keyID, apiKey.ID, "API key should exist in the database")

		fullURL, err := client.URL.Parse("/api/v2/users/logout")
		require.NoError(t, err, "Server URL should parse successfully")

		res, err := client.Request(ctx, http.MethodPost, fullURL.String(), nil)
		numLogs++ // add an audit log for logout

		require.NoError(t, err, "/logout request should succeed")
		res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionLogout, auditor.AuditLogs()[numLogs-1].Action)

		cookies := res.Cookies()

		var found bool
		for _, cookie := range cookies {
			if cookie.Name == optimus-ide-collabsdk.SessionTokenCookie {
				require.Equal(t, optimus-ide-collabsdk.SessionTokenCookie, cookie.Name, "Cookie should be the auth cookie")
				require.Equal(t, -1, cookie.MaxAge, "Cookie should be set to delete")
				found = true
			}
		}
		require.True(t, found, "auth cookie should be returned")

		_, err = client.APIKeyByID(ctx, owner.UserID.String(), keyID)
		sdkErr := &optimus-ide-collabsdk.Error{}
		require.ErrorAs(t, err, &sdkErr)
		require.Equal(t, http.StatusUnauthorized, sdkErr.StatusCode(), "Expecting 401")
	})
}

// nolint:bodyclose
func TestPostUsers(t *testing.T) {
	t.Parallel()
	t.Run("NoAuth", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{})
		require.Error(t, err)
	})

	t.Run("Conflicting", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		me, err := client.User(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		_, err = client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           me.Email,
			Username:        me.Username,
			Password:        "MySecurePassword!",
			OrganizationIDs: []uuid.UUID{uuid.New()},
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusConflict, apiErr.StatusCode())
	})

	t.Run("OrganizationNotFound", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{uuid.New()},
			Email:           "another@user.org",
			Username:        "someone-else",
			Password:        "SomeSecurePassword!",
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	t.Run("Create", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client)
		numLogs++ // add an audit log for user create
		numLogs++ // add an audit log for login

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		user, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{firstUser.OrganizationID},
			Email:           "another@user.org",
			Username:        "someone-else",
			Password:        "SomeSecurePassword!",
		})
		require.NoError(t, err)

		// User should default to dormant.
		require.Equal(t, optimus-ide-collabsdk.UserStatusDormant, user.Status)

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionCreate, auditor.AuditLogs()[numLogs-1].Action)
		require.Equal(t, database.AuditActionLogin, auditor.AuditLogs()[numLogs-2].Action)

		require.Len(t, user.OrganizationIDs, 1)
		assert.Equal(t, firstUser.OrganizationID, user.OrganizationIDs[0])
	})

	t.Run("CreateWithStatus", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client)
		numLogs++ // add an audit log for user create
		numLogs++ // add an audit log for login

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		user, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{firstUser.OrganizationID},
			Email:           "another@user.org",
			Username:        "someone-else",
			Password:        "SomeSecurePassword!",
			UserStatus:      ptr.Ref(optimus-ide-collabsdk.UserStatusActive),
		})
		require.NoError(t, err)

		require.Equal(t, optimus-ide-collabsdk.UserStatusActive, user.Status)

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionCreate, auditor.AuditLogs()[numLogs-1].Action)
		require.Equal(t, database.AuditActionLogin, auditor.AuditLogs()[numLogs-2].Action)

		require.Len(t, user.OrganizationIDs, 1)
		assert.Equal(t, firstUser.OrganizationID, user.OrganizationIDs[0])
	})

	t.Run("LastSeenAt", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		client := optimus-ide-collabdtest.New(t, nil)
		firstUserResp := optimus-ide-collabdtest.CreateFirstUser(t, client)

		firstUser, err := client.User(ctx, firstUserResp.UserID.String())
		require.NoError(t, err)

		_, _ = optimus-ide-collabdtest.CreateAnotherUser(t, client, firstUserResp.OrganizationID)

		allUsersRes, err := client.Users(ctx, optimus-ide-collabsdk.UsersRequest{})
		require.NoError(t, err)

		require.Len(t, allUsersRes.Users, 2)

		// We sent the "GET Users" request with the first user, but the second user
		// should be Never since they haven't performed a request.
		for _, user := range allUsersRes.Users {
			if user.ID == firstUser.ID {
				require.WithinDuration(t, firstUser.LastSeenAt, dbtime.Now(), testutil.WaitShort)
			} else {
				require.Zero(t, user.LastSeenAt)
			}
		}
	})

	t.Run("CreateNoneLoginType", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		first := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{first.OrganizationID},
			Email:           "another@user.org",
			Username:        "someone-else",
			Password:        "",
			UserLoginType:   optimus-ide-collabsdk.LoginTypeNone,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
		require.Contains(t, apiErr.Message, "service account")
	})

	t.Run("CreateOIDCLoginType", func(t *testing.T) {
		t.Parallel()
		email := "another@user.org"
		fake := oidctest.NewFakeIDP(t,
			oidctest.WithServing(),
		)
		cfg := fake.OIDCConfig(t, nil, func(cfg *optimus-ide-collabd.OIDCConfig) {
			cfg.AllowSignups = true
		})

		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			OIDCConfig: cfg,
		})
		first := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{first.OrganizationID},
			Email:           email,
			Username:        "someone-else",
			Password:        "",
			UserLoginType:   optimus-ide-collabsdk.LoginTypeOIDC,
		})
		require.NoError(t, err)

		// Try to log in with OIDC.
		userClient, _ := fake.Login(t, client, jwt.MapClaims{
			"email":          email,
			"email_verified": true,
			"sub":            uuid.NewString(),
		})

		found, err := userClient.User(ctx, "me")
		require.NoError(t, err)
		require.Equal(t, found.LoginType, optimus-ide-collabsdk.LoginTypeOIDC)
	})

	t.Run("ServiceAccount/Unlicensed", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		first := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{first.OrganizationID},
			Username:        "service-acct-ok",
			UserLoginType:   optimus-ide-collabsdk.LoginTypeNone,
			ServiceAccount:  true,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusForbidden, apiErr.StatusCode())
		require.Contains(t, apiErr.Message, "Premium feature")
	})

	t.Run("NonServiceAccount/WithoutEmail", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		first := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{first.OrganizationID},
			Username:        "regular-no-email",
			UserLoginType:   optimus-ide-collabsdk.LoginTypePassword,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})
}

func TestNotifyCreatedUser(t *testing.T) {
	t.Parallel()

	t.Run("OwnerNotified", func(t *testing.T) {
		t.Parallel()

		// given
		notifyEnq := &notificationstest.FakeEnqueuer{}
		adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			NotificationsEnqueuer: notifyEnq,
		})
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		// when
		user, err := adminClient.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{firstUser.OrganizationID},
			Email:           "another@user.org",
			Username:        "someone-else",
			Password:        "SomeSecurePassword!",
		})
		require.NoError(t, err)

		// then
		sent := notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateUserAccountCreated))
		require.Len(t, sent, 1)
		require.Equal(t, notifications.TemplateUserAccountCreated, sent[0].TemplateID)
		require.Equal(t, firstUser.UserID, sent[0].UserID)
		require.Contains(t, sent[0].Targets, user.ID)
		require.Equal(t, user.Username, sent[0].Labels["created_account_name"])

		require.IsType(t, map[string]any{}, sent[0].Data["user"])
		userData := sent[0].Data["user"].(map[string]any)
		require.Equal(t, user.ID, userData["id"])
		require.Equal(t, user.Name, userData["name"])
		require.Equal(t, user.Email, userData["email"])
	})

	t.Run("UserAdminNotified", func(t *testing.T) {
		t.Parallel()

		// given
		notifyEnq := &notificationstest.FakeEnqueuer{}
		adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			NotificationsEnqueuer: notifyEnq,
		})
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		userAdmin, err := adminClient.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{firstUser.OrganizationID},
			Email:           "user-admin@user.org",
			Username:        "mr-user-admin",
			Password:        "SomeSecurePassword!",
		})
		require.NoError(t, err)

		_, err = adminClient.UpdateUserRoles(ctx, userAdmin.Username, optimus-ide-collabsdk.UpdateRoles{
			Roles: []string{
				rbac.RoleUserAdmin().String(),
			},
		})
		require.NoError(t, err)

		// when
		member, err := adminClient.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			OrganizationIDs: []uuid.UUID{firstUser.OrganizationID},
			Email:           "another@user.org",
			Username:        "someone-else",
			Password:        "SomeSecurePassword!",
		})
		require.NoError(t, err)

		// then
		sent := notifyEnq.Sent()
		require.Len(t, sent, 3)

		// "User admin" account created, "owner" notified
		ownerNotifiedAboutUserAdmin := notifyEnq.Sent(func(n *notificationstest.FakeNotification) bool {
			return n.TemplateID == notifications.TemplateUserAccountCreated &&
				n.UserID == firstUser.UserID &&
				slices.Contains(n.Targets, userAdmin.ID) &&
				n.Labels["created_account_name"] == userAdmin.Username
		})
		require.Len(t, ownerNotifiedAboutUserAdmin, 1)

		// "Member" account created, "owner" notified
		ownerNotifiedAboutMember := notifyEnq.Sent(func(n *notificationstest.FakeNotification) bool {
			return n.TemplateID == notifications.TemplateUserAccountCreated &&
				n.UserID == firstUser.UserID &&
				slices.Contains(n.Targets, member.ID) &&
				n.Labels["created_account_name"] == member.Username
		})
		require.Len(t, ownerNotifiedAboutMember, 1)

		// "Member" account created, "user admin" notified
		userAdminNotifiedAboutMember := notifyEnq.Sent(func(n *notificationstest.FakeNotification) bool {
			return n.TemplateID == notifications.TemplateUserAccountCreated &&
				n.UserID == userAdmin.ID &&
				slices.Contains(n.Targets, member.ID) &&
				n.Labels["created_account_name"] == member.Username
		})
		require.Len(t, userAdminNotifiedAboutMember, 1)
	})
}

func TestUpdateUserProfile(t *testing.T) {
	t.Parallel()
	t.Run("UserNotFound", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.UpdateUserProfile(ctx, uuid.New().String(), optimus-ide-collabsdk.UpdateUserProfileRequest{
			Username: "newusername",
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		// Right now, we are raising a BAD request error because we don't support a
		// user accessing other users info
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	t.Run("ConflictingUsername", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		existentUser, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           "bruno@optimus-ide-collab.com",
			Username:        "bruno",
			Password:        "SomeSecurePassword!",
			OrganizationIDs: []uuid.UUID{user.OrganizationID},
		})
		require.NoError(t, err)
		_, err = client.UpdateUserProfile(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserProfileRequest{
			Username: existentUser.Username,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusConflict, apiErr.StatusCode())
	})

	t.Run("UpdateSelf", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		optimus-ide-collabdtest.CreateFirstUser(t, client)
		numLogs++ // add an audit log for login

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		me, err := client.User(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)

		userProfile, err := client.UpdateUserProfile(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserProfileRequest{
			Username: me.Username + "1",
			Name:     me.Name + "1",
		})
		numLogs++ // add an audit log for user update

		require.NoError(t, err)
		require.Equal(t, me.Username+"1", userProfile.Username)
		require.Equal(t, me.Name+"1", userProfile.Name)

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionWrite, auditor.AuditLogs()[numLogs-1].Action)
	})

	t.Run("UpdateSelfAsMember_Name", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client)
		numLogs++ // add an audit log for login

		memberClient, memberUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, firstUser.OrganizationID)
		numLogs++ // add an audit log for user creation

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		newName := optimus-ide-collabdtest.RandomName(t)
		userProfile, err := memberClient.UpdateUserProfile(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserProfileRequest{
			Name:     newName,
			Username: memberUser.Username,
		})
		numLogs++ // add an audit log for user update
		numLogs++ // add an audit log for API key creation

		require.NoError(t, err)
		require.Equal(t, memberUser.Username, userProfile.Username)
		require.Equal(t, newName, userProfile.Name)

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionWrite, auditor.AuditLogs()[numLogs-1].Action)
	})

	t.Run("UpdateSelfAsMember_Username", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})

		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client)
		memberClient, memberUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		newUsername := optimus-ide-collabdtest.RandomUsername(t)
		_, err := memberClient.UpdateUserProfile(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserProfileRequest{
			Name:     memberUser.Name,
			Username: newUsername,
		})

		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	t.Run("UpdateMemberAsAdmin_Username", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		adminUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)
		numLogs++ // add an audit log for login

		_, memberUser := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)
		numLogs++ // add an audit log for user creation

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		newUsername := optimus-ide-collabdtest.RandomUsername(t)
		userProfile, err := adminClient.UpdateUserProfile(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserProfileRequest{
			Name:     memberUser.Name,
			Username: newUsername,
		})

		numLogs++ // add an audit log for user update
		numLogs++ // add an audit log for API key creation

		require.NoError(t, err)
		require.Equal(t, newUsername, userProfile.Username)
		require.Equal(t, memberUser.Name, userProfile.Name)

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionWrite, auditor.AuditLogs()[numLogs-1].Action)
	})

	t.Run("InvalidRealUserName", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           "john@optimus-ide-collab.com",
			Username:        "john",
			Password:        "SomeSecurePassword!",
			OrganizationIDs: []uuid.UUID{user.OrganizationID},
		})
		require.NoError(t, err)
		_, err = client.UpdateUserProfile(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserProfileRequest{
			Name: " Mr Bean", // must not have leading space
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})

	t.Run("UpdateAvatar", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		me, err := client.User(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)

		// The first user is a password user, so the avatar is editable.
		const newAvatar = "/emojis/1f600.png"
		userProfile, err := client.UpdateUserProfile(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserProfileRequest{
			Username:  me.Username,
			Name:      me.Name,
			AvatarURL: newAvatar,
		})
		require.NoError(t, err)
		require.Equal(t, newAvatar, userProfile.AvatarURL)
	})

	t.Run("IgnoresAvatarForSSOUser", func(t *testing.T) {
		t.Parallel()
		client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
		// The first user is an owner and can update other users' profiles.
		optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		// Avatars for SSO users are synced from the identity provider on login,
		// so a submitted avatar must be ignored and the existing one preserved.
		ssoUser := dbgen.User(t, db, database.User{
			Email:     "sso-avatar@optimus-ide-collab.com",
			Username:  "sso-avatar",
			LoginType: database.LoginTypeOIDC,
		})

		// dbgen.User does not persist the avatar at creation, so set it directly
		// to emulate an avatar synced from the identity provider.
		const idpAvatar = "https://idp.example.com/avatar.png"
		//nolint:gocritic // Test setup requires a system context to set the avatar.
		ssoUser, err := db.UpdateUserProfile(dbauthz.AsSystemRestricted(ctx), database.UpdateUserProfileParams{
			ID:        ssoUser.ID,
			Email:     ssoUser.Email,
			Name:      ssoUser.Name,
			AvatarURL: idpAvatar,
			Username:  ssoUser.Username,
			UpdatedAt: dbtime.Now(),
		})
		require.NoError(t, err)

		userProfile, err := client.UpdateUserProfile(ctx, ssoUser.ID.String(), optimus-ide-collabsdk.UpdateUserProfileRequest{
			Username:  ssoUser.Username,
			Name:      ssoUser.Name,
			AvatarURL: "/emojis/1f600.png",
		})
		require.NoError(t, err)
		require.Equal(t, idpAvatar, userProfile.AvatarURL)
	})
}

func TestUpdateUserPassword(t *testing.T) {
	t.Parallel()

	t.Run("MemberCantUpdateAdminPassword", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		err := member.UpdateUserPassword(ctx, owner.UserID.String(), optimus-ide-collabsdk.UpdateUserPasswordRequest{
			Password: "newpassword",
		})
		require.Error(t, err, "member should not be able to update admin password")
	})

	t.Run("AdminCanUpdateMemberPassword", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		member, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           "optimus-ide-collab@optimus-ide-collab.com",
			Username:        "optimus-ide-collab",
			Password:        "SomeStrongPassword!",
			OrganizationIDs: []uuid.UUID{owner.OrganizationID},
		})
		require.NoError(t, err, "create member")
		err = client.UpdateUserPassword(ctx, member.ID.String(), optimus-ide-collabsdk.UpdateUserPasswordRequest{
			Password: "SomeNewStrongPassword!",
		})
		require.NoError(t, err, "admin should be able to update member password")
		// Check if the member can login using the new password
		_, err = client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    "optimus-ide-collab@optimus-ide-collab.com",
			Password: "SomeNewStrongPassword!",
		})
		require.NoError(t, err, "member should login successfully with the new password")
	})

	t.Run("UserAdminCanUpdateMemberPassword", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)

		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())
		memberClient, member := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		apikey1, err := memberClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{})
		require.NoError(t, err)
		apikey1ID, _, ok := strings.Cut(apikey1.Key, "-")
		require.True(t, ok)

		apikey2, err := memberClient.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{})
		require.NoError(t, err)
		apikey2ID, _, ok := strings.Cut(apikey2.Key, "-")
		require.True(t, ok)

		// Confirm the token IDs resolve before the reset, so later
		// 404s prove password reset removed them.
		_, err = memberClient.APIKeyByID(ctx, member.ID.String(), apikey1ID)
		require.NoError(t, err)
		_, err = memberClient.APIKeyByID(ctx, member.ID.String(), apikey2ID)
		require.NoError(t, err)

		err = userAdmin.UpdateUserPassword(ctx, member.ID.String(), optimus-ide-collabsdk.UpdateUserPasswordRequest{
			Password: "SomeNewStrongPassword!",
		})
		require.NoError(t, err, "user-admin should be able to update member password")

		// The old session token was deleted with the password reset, so the
		// member client is unauthorized before it logs back in.
		_, err = memberClient.APIKeyByID(ctx, member.ID.String(), apikey1ID)
		require.Error(t, err)
		cerr := optimus-ide-collabdtest.SDKError(t, err)
		require.Equal(t, http.StatusUnauthorized, cerr.StatusCode())

		resp, err := client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    member.Email,
			Password: "SomeNewStrongPassword!",
		})
		require.NoError(t, err, "member should login successfully with the new password")

		memberClient = optimus-ide-collabsdk.New(
			memberClient.URL,
			optimus-ide-collabsdk.WithSessionToken(resp.SessionToken),
			optimus-ide-collabsdk.WithHTTPClient(optimus-ide-collabdtest.NewIsolatedHTTPClient(memberClient.URL)),
		)

		// After re-authentication, the old API keys should be gone from storage.
		_, err = memberClient.APIKeyByID(ctx, member.ID.String(), apikey1ID)
		require.Error(t, err)
		cerr = optimus-ide-collabdtest.SDKError(t, err)
		require.Equal(t, http.StatusNotFound, cerr.StatusCode())

		_, err = memberClient.APIKeyByID(ctx, member.ID.String(), apikey2ID)
		require.Error(t, err)
		cerr = optimus-ide-collabdtest.SDKError(t, err)
		require.Equal(t, http.StatusNotFound, cerr.StatusCode())
	})

	t.Run("AuditorCantUpdateOtherUserPassword", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)

		auditor, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleAuditor())

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		member, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           "optimus-ide-collab@optimus-ide-collab.com",
			Username:        "optimus-ide-collab",
			Password:        "SomeStrongPassword!",
			OrganizationIDs: []uuid.UUID{owner.OrganizationID},
		})
		require.NoError(t, err, "create member")

		err = auditor.UpdateUserPassword(ctx, member.ID.String(), optimus-ide-collabsdk.UpdateUserPasswordRequest{
			Password: "SomeNewStrongPassword!",
		})
		require.Error(t, err, "auditor should not be able to update member password")
		require.ErrorContains(t, err, "unexpected status code 404: Resource not found or you do not have access to this resource")
	})

	t.Run("MemberCanUpdateOwnPassword", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		numLogs++ // add an audit log for user create
		numLogs++ // add an audit log for login

		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		numLogs++ // add an audit log for user create

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		err := member.UpdateUserPassword(ctx, "me", optimus-ide-collabsdk.UpdateUserPasswordRequest{
			OldPassword: "SomeSecurePassword!",
			Password:    "MyNewSecurePassword!",
		})
		numLogs++ // add an audit log for user update

		require.NoError(t, err, "member should be able to update own password")

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionWrite, auditor.AuditLogs()[numLogs-1].Action)
	})

	t.Run("MemberCantUpdateOwnPasswordWithoutOldPassword", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		err := member.UpdateUserPassword(ctx, "me", optimus-ide-collabsdk.UpdateUserPasswordRequest{
			Password: "newpassword",
		})
		require.Error(t, err, "member should not be able to update own password without providing old password")
		require.ErrorContains(t, err, "Old password is required.")
	})

	t.Run("AuditorCantTellIfPasswordIncorrect", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		adminClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})

		adminUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

		auditorClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient,
			adminUser.OrganizationID,
			rbac.RoleAuditor(),
		)

		_, memberUser := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)
		numLogs := len(auditor.AuditLogs())

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		err := auditorClient.UpdateUserPassword(ctx, memberUser.ID.String(), optimus-ide-collabsdk.UpdateUserPasswordRequest{
			Password: "MySecurePassword!",
		})
		numLogs++ // add an audit log for user update

		require.Error(t, err, "auditors shouldn't be able to update passwords")
		var httpErr *optimus-ide-collabsdk.Error
		require.True(t, xerrors.As(err, &httpErr))
		// ensure that the error we get is "not found" and not "bad request"
		require.Equal(t, http.StatusNotFound, httpErr.StatusCode())

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionWrite, auditor.AuditLogs()[numLogs-1].Action)
		require.Equal(t, int32(http.StatusNotFound), auditor.AuditLogs()[numLogs-1].StatusCode)
	})

	t.Run("AdminCantUpdateOwnPasswordWithoutOldPassword", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
		numLogs++ // add an audit log for login

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		err := client.UpdateUserPassword(ctx, "me", optimus-ide-collabsdk.UpdateUserPasswordRequest{
			Password: "MySecurePassword!",
		})
		numLogs++ // add an audit log for user update

		require.Error(t, err, "admin should not be able to update own password without providing old password")
		require.ErrorContains(t, err, "Old password is required.")

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionWrite, auditor.AuditLogs()[numLogs-1].Action)
	})

	t.Run("ValidateUserPassword", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})

		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		resp, err := client.ValidateUserPassword(ctx, optimus-ide-collabsdk.ValidateUserPasswordRequest{
			Password: "MySecurePassword!",
		})

		require.NoError(t, err, "users shoud be able to validate complexity of a potential new password")
		require.True(t, resp.Valid)
	})

	t.Run("ChangingPasswordDeletesKeys", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		ctx := testutil.Context(t, testutil.WaitLong)

		apikey1, err := client.CreateToken(ctx, user.UserID.String(), optimus-ide-collabsdk.CreateTokenRequest{})
		require.NoError(t, err)

		apikey2, err := client.CreateToken(ctx, user.UserID.String(), optimus-ide-collabsdk.CreateTokenRequest{})
		require.NoError(t, err)

		err = client.UpdateUserPassword(ctx, "me", optimus-ide-collabsdk.UpdateUserPasswordRequest{
			OldPassword: "SomeSecurePassword!",
			Password:    "MyNewSecurePassword!",
		})
		require.NoError(t, err)

		// Trying to get an API key should fail since our client's token
		// has been deleted.
		_, err = client.APIKeyByID(ctx, user.UserID.String(), apikey1.Key)
		require.Error(t, err)
		cerr := optimus-ide-collabdtest.SDKError(t, err)
		require.Equal(t, http.StatusUnauthorized, cerr.StatusCode())

		resp, err := client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    optimus-ide-collabdtest.FirstUserParams.Email,
			Password: "MyNewSecurePassword!",
		})
		require.NoError(t, err)

		client.SetSessionToken(resp.SessionToken)

		// Trying to get an API key should fail since all keys are deleted
		// on password change.
		_, err = client.APIKeyByID(ctx, user.UserID.String(), apikey1.Key)
		require.Error(t, err)
		cerr = optimus-ide-collabdtest.SDKError(t, err)
		require.Equal(t, http.StatusNotFound, cerr.StatusCode())

		_, err = client.APIKeyByID(ctx, user.UserID.String(), apikey2.Key)
		require.Error(t, err)
		cerr = optimus-ide-collabdtest.SDKError(t, err)
		require.Equal(t, http.StatusNotFound, cerr.StatusCode())
	})

	t.Run("UserAdminCannotResetOwnerPassword", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		err := userAdmin.UpdateUserPassword(ctx, owner.UserID.String(), optimus-ide-collabsdk.UpdateUserPasswordRequest{
			Password: "SomeNewStrongPassword!",
		})
		require.Error(t, err, "user-admin should not be able to reset owner password")
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
		require.Contains(t, apiErr.Message, "Only owners can change the password of an owner")
	})

	t.Run("OwnerCanResetOwnerPassword", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		anotherOwner, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           "another-owner@optimus-ide-collab.com",
			Username:        "another-owner",
			Password:        "SomeStrongPassword!",
			OrganizationIDs: []uuid.UUID{owner.OrganizationID},
		})
		require.NoError(t, err)
		_, err = client.UpdateUserRoles(ctx, anotherOwner.ID.String(), optimus-ide-collabsdk.UpdateRoles{
			Roles: []string{rbac.RoleOwner().String()},
		})
		require.NoError(t, err)

		err = client.UpdateUserPassword(ctx, anotherOwner.ID.String(), optimus-ide-collabsdk.UpdateUserPasswordRequest{
			Password: "SomeNewStrongPassword!",
		})
		require.NoError(t, err, "owner should be able to reset another owner's password")

		_, err = client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    "another-owner@optimus-ide-collab.com",
			Password: "SomeNewStrongPassword!",
		})
		require.NoError(t, err, "other owner should login with the new password")
	})

	t.Run("PasswordsMustDiffer", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
		ctx := testutil.Context(t, testutil.WaitLong)

		err := client.UpdateUserPassword(ctx, "me", optimus-ide-collabsdk.UpdateUserPasswordRequest{
			Password: optimus-ide-collabdtest.FirstUserParams.Password,
		})
		require.Error(t, err)
		cerr := optimus-ide-collabdtest.SDKError(t, err)
		require.Equal(t, http.StatusBadRequest, cerr.StatusCode())
	})
}

// TestInitialRoles ensures the starting roles for the first user are correct.
func TestInitialRoles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client := optimus-ide-collabdtest.New(t, nil)
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)

	roles, err := client.UserRoles(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)
	require.ElementsMatch(t, roles.Roles, []string{
		optimus-ide-collabsdk.RoleOwner,
	}, "should be a member and admin")

	require.ElementsMatch(t, roles.OrganizationRoles[first.OrganizationID], []string{}, "should be a member")
}

func TestPutUserSuspend(t *testing.T) {
	t.Parallel()

	t.Run("SuspendAnOwner", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		me := optimus-ide-collabdtest.CreateFirstUser(t, client)
		_, user := optimus-ide-collabdtest.CreateAnotherUser(t, client, me.OrganizationID, rbac.RoleOwner())

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		_, err := client.UpdateUserStatus(ctx, user.Username, optimus-ide-collabsdk.UserStatusSuspended)
		require.Error(t, err, "cannot suspend owners")
	})

	t.Run("SuspendAnotherUser", func(t *testing.T) {
		t.Parallel()
		auditor := audit.NewMock()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Auditor: auditor})
		numLogs := len(auditor.AuditLogs())

		me := optimus-ide-collabdtest.CreateFirstUser(t, client)
		numLogs++ // add an audit log for user create
		numLogs++ // add an audit log for login

		_, user := optimus-ide-collabdtest.CreateAnotherUser(t, client, me.OrganizationID)
		numLogs++ // add an audit log for user create

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		user, err := client.UpdateUserStatus(ctx, user.Username, optimus-ide-collabsdk.UserStatusSuspended)
		require.NoError(t, err)
		require.Equal(t, user.Status, optimus-ide-collabsdk.UserStatusSuspended)
		numLogs++ // add an audit log for user update

		require.Len(t, auditor.AuditLogs(), numLogs)
		require.Equal(t, database.AuditActionWrite, auditor.AuditLogs()[numLogs-1].Action)
	})

	t.Run("SuspendItSelf", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		client.User(ctx, optimus-ide-collabsdk.Me)
		_, err := client.UpdateUserStatus(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UserStatusSuspended)

		require.ErrorContains(t, err, "suspend yourself", "cannot suspend yourself")
	})
}

func TestActivateDormantUser(t *testing.T) {
	t.Parallel()
	client := optimus-ide-collabdtest.New(t, nil)

	// Create users
	me := optimus-ide-collabdtest.CreateFirstUser(t, client)
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	anotherUser, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
		Email:           "optimus-ide-collab@optimus-ide-collab.com",
		Username:        "optimus-ide-collab",
		Password:        "SomeStrongPassword!",
		OrganizationIDs: []uuid.UUID{me.OrganizationID},
	})
	require.NoError(t, err)

	// Ensure that new user has dormant account
	require.Equal(t, optimus-ide-collabsdk.UserStatusDormant, anotherUser.Status)

	// Activate user account
	_, err = client.UpdateUserStatus(ctx, anotherUser.Username, optimus-ide-collabsdk.UserStatusActive)
	require.NoError(t, err)

	// Verify if the account is active now
	anotherUser, err = client.User(ctx, anotherUser.Username)
	require.NoError(t, err)
	require.Equal(t, optimus-ide-collabsdk.UserStatusActive, anotherUser.Status)
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	// Single instance shared across all sub-tests. All lookups
	// are read-only against the first user.
	client := optimus-ide-collabdtest.New(t, nil)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client)

	t.Run("ByMe", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		user, err := client.User(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, firstUser.UserID, user.ID)
		require.Equal(t, firstUser.OrganizationID, user.OrganizationIDs[0])
	})

	t.Run("ByID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		user, err := client.User(ctx, firstUser.UserID.String())
		require.NoError(t, err)
		require.Equal(t, firstUser.UserID, user.ID)
		require.Equal(t, firstUser.OrganizationID, user.OrganizationIDs[0])
	})

	t.Run("ByUsername", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		exp, err := client.User(ctx, firstUser.UserID.String())
		require.NoError(t, err)

		user, err := client.User(ctx, exp.Username)
		require.NoError(t, err)
		require.Equal(t, exp.ID, user.ID)
	})
}

func TestGetUsersFilter(t *testing.T) {
	t.Parallel()

	client, _, api := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		IncludeProvisionerDaemon: true,
		OIDCConfig: &optimus-ide-collabd.OIDCConfig{
			AllowSignups: true,
		},
	})
	_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

	setupCtx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	optimus-ide-collabdtest.UsersFilter(setupCtx, t, client, api.Database, nil, nil, func(testCtx context.Context, req optimus-ide-collabsdk.UsersRequest) []optimus-ide-collabsdk.ReducedUser {
		res, err := client.Users(testCtx, req)
		require.NoError(t, err)
		reduced := make([]optimus-ide-collabsdk.ReducedUser, len(res.Users))
		for i, user := range res.Users {
			reduced[i] = user.ReducedUser
		}
		return reduced
	})
}

func TestGetUsersPagination(t *testing.T) {
	t.Parallel()
	client := optimus-ide-collabdtest.New(t, nil)
	_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	optimus-ide-collabdtest.UsersPagination(ctx, t, client, nil, func(req optimus-ide-collabsdk.UsersRequest) ([]optimus-ide-collabsdk.ReducedUser, int) {
		res, err := client.Users(ctx, req)
		require.NoError(t, err)
		reduced := make([]optimus-ide-collabsdk.ReducedUser, len(res.Users))
		for i, user := range res.Users {
			reduced[i] = user.ReducedUser
		}
		return reduced, res.Count
	})
}

func TestPostTokens(t *testing.T) {
	t.Parallel()
	client := optimus-ide-collabdtest.New(t, nil)
	_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	apiKey, err := client.CreateToken(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.CreateTokenRequest{})
	require.NotNil(t, apiKey)
	require.GreaterOrEqual(t, len(apiKey.Key), 2)
	require.NoError(t, err)
}

func TestUserTerminalFont(t *testing.T) {
	t.Parallel()

	// Single instance shared across all sub-tests. Each sub-test
	// creates its own non-admin user for isolation.
	adminClient := optimus-ide-collabdtest.New(t, nil)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

	t.Run("valid font", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		// given
		initial, err := client.GetUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontName(""), initial.TerminalFont)

		// when
		updated, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "light",
			TerminalFont:    "fira-code",
		})
		require.NoError(t, err)

		// then
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, updated.TerminalFont)
	})

	t.Run("unsupported font", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		// given
		initial, err := client.GetUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontName(""), initial.TerminalFont)

		// when
		_, err = client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "light",
			TerminalFont:    "foobar",
		})

		// then
		require.Error(t, err)
	})

	t.Run("undefined font is not ok", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		// given
		initial, err := client.GetUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontName(""), initial.TerminalFont)

		// when
		_, err = client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "light",
			TerminalFont:    "",
		})

		// then
		require.Error(t, err)
	})
}

func TestUserThemeMode(t *testing.T) {
	t.Parallel()

	adminClient := optimus-ide-collabdtest.New(t, nil)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

	t.Run("defaults to empty", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		initial, err := client.GetUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		// A fresh user has never written any theme_* key. The GET handler
		// should return empty strings rather than error out.
		require.Equal(t, optimus-ide-collabsdk.ThemeModeUnset, initial.ThemeMode)
		require.Equal(t, "", initial.ThemeLight)
		require.Equal(t, "", initial.ThemeDark)
	})

	t.Run("sync mode roundtrip", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		updated, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "dark-tritan",
			ThemeMode:       optimus-ide-collabsdk.ThemeModeSync,
			ThemeLight:      "light-tritan",
			ThemeDark:       "dark-tritan",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontGeistMono,
		})
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSync, updated.ThemeMode)
		require.Equal(t, "light-tritan", updated.ThemeLight)
		require.Equal(t, "dark-tritan", updated.ThemeDark)

		// Fetched values should match.
		fetched, err := client.GetUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSync, fetched.ThemeMode)
		require.Equal(t, "light-tritan", fetched.ThemeLight)
		require.Equal(t, "dark-tritan", fetched.ThemeDark)
	})

	t.Run("sync mode accepts any concrete theme per slot", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		updated, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "dark-tritan",
			ThemeMode:       optimus-ide-collabsdk.ThemeModeSync,
			ThemeLight:      "dark-tritan",
			ThemeDark:       "light-tritan",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontGeistMono,
		})
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSync, updated.ThemeMode)
		require.Equal(t, "dark-tritan", updated.ThemeLight)
		require.Equal(t, "light-tritan", updated.ThemeDark)
	})

	t.Run("empty theme_mode is accepted for back-compat", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		// A concrete legacy preference plus an unset mode is enough for
		// modern clients to treat the user as single mode. The server does
		// not write the new fields for old clients because doing so would
		// erase existing sync settings.
		updated, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "dark",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontGeistMono,
		})
		require.NoError(t, err)
		require.Equal(t, "dark", updated.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeUnset, updated.ThemeMode)
		require.Equal(t, "", updated.ThemeLight)
		require.Equal(t, "", updated.ThemeDark)
	})

	t.Run("omitted theme fields preserve sync settings", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		_, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "dark-tritan",
			ThemeMode:       optimus-ide-collabsdk.ThemeModeSync,
			ThemeLight:      "light-tritan",
			ThemeDark:       "dark-tritan",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontGeistMono,
		})
		require.NoError(t, err)

		updated, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "dark",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontFiraCode,
		})
		require.NoError(t, err)
		require.Equal(t, "dark", updated.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSync, updated.ThemeMode)
		require.Equal(t, "light-tritan", updated.ThemeLight)
		require.Equal(t, "dark-tritan", updated.ThemeDark)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, updated.TerminalFont)

		fetched, err := client.GetUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, "dark", fetched.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSync, fetched.ThemeMode)
		require.Equal(t, "light-tritan", fetched.ThemeLight)
		require.Equal(t, "dark-tritan", fetched.ThemeDark)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, fetched.TerminalFont)
	})

	t.Run("single mode with omitted slots preserves sync settings", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		_, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "dark-tritan",
			ThemeMode:       optimus-ide-collabsdk.ThemeModeSync,
			ThemeLight:      "light-tritan",
			ThemeDark:       "dark-tritan",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontGeistMono,
		})
		require.NoError(t, err)

		updated, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "dark",
			ThemeMode:       optimus-ide-collabsdk.ThemeModeSingle,
			TerminalFont:    optimus-ide-collabsdk.TerminalFontFiraCode,
		})
		require.NoError(t, err)
		require.Equal(t, "dark", updated.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSingle, updated.ThemeMode)
		require.Equal(t, "light-tritan", updated.ThemeLight)
		require.Equal(t, "dark-tritan", updated.ThemeDark)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, updated.TerminalFont)

		fetched, err := client.GetUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, "dark", fetched.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSingle, fetched.ThemeMode)
		require.Equal(t, "light-tritan", fetched.ThemeLight)
		require.Equal(t, "dark-tritan", fetched.ThemeDark)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, fetched.TerminalFont)
	})

	t.Run("single mode with explicit slots updates slots", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		updated, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "light-tritan",
			ThemeMode:       optimus-ide-collabsdk.ThemeModeSingle,
			ThemeLight:      "dark-tritan",
			ThemeDark:       "light-protan-deuter",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontFiraCode,
		})
		require.NoError(t, err)
		require.Equal(t, "light-tritan", updated.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSingle, updated.ThemeMode)
		require.Equal(t, "dark-tritan", updated.ThemeLight)
		require.Equal(t, "light-protan-deuter", updated.ThemeDark)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, updated.TerminalFont)

		fetched, err := client.GetUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, "light-tritan", fetched.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSingle, fetched.ThemeMode)
		require.Equal(t, "dark-tritan", fetched.ThemeLight)
		require.Equal(t, "light-protan-deuter", fetched.ThemeDark)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, fetched.TerminalFont)
	})

	t.Run("single mode with one explicit slot updates only that slot", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		_, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "dark-tritan",
			ThemeMode:       optimus-ide-collabsdk.ThemeModeSync,
			ThemeLight:      "light-tritan",
			ThemeDark:       "dark-tritan",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontGeistMono,
		})
		require.NoError(t, err)

		updated, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "light",
			ThemeMode:       optimus-ide-collabsdk.ThemeModeSingle,
			ThemeLight:      "light-protan-deuter",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontFiraCode,
		})
		require.NoError(t, err)
		require.Equal(t, "light", updated.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSingle, updated.ThemeMode)
		require.Equal(t, "light-protan-deuter", updated.ThemeLight)
		require.Equal(t, "dark-tritan", updated.ThemeDark)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, updated.TerminalFont)

		fetched, err := client.GetUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, "light", fetched.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeSingle, fetched.ThemeMode)
		require.Equal(t, "light-protan-deuter", fetched.ThemeLight)
		require.Equal(t, "dark-tritan", fetched.ThemeDark)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, fetched.TerminalFont)
	})

	t.Run("legacy auto with omitted theme_mode clears mode", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		_, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "dark-tritan",
			ThemeMode:       optimus-ide-collabsdk.ThemeModeSync,
			ThemeLight:      "light-tritan",
			ThemeDark:       "dark-tritan",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontGeistMono,
		})
		require.NoError(t, err)

		updated, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "auto",
			TerminalFont:    optimus-ide-collabsdk.TerminalFontFiraCode,
		})
		require.NoError(t, err)
		require.Equal(t, "auto", updated.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeUnset, updated.ThemeMode)
		require.Equal(t, "light-tritan", updated.ThemeLight)
		require.Equal(t, "dark-tritan", updated.ThemeDark)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, updated.TerminalFont)

		fetched, err := client.GetUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, "auto", fetched.ThemePreference)
		require.Equal(t, optimus-ide-collabsdk.ThemeModeUnset, fetched.ThemeMode)
		require.Equal(t, "light-tritan", fetched.ThemeLight)
		require.Equal(t, "dark-tritan", fetched.ThemeDark)
		require.Equal(t, optimus-ide-collabsdk.TerminalFontFiraCode, fetched.TerminalFont)
	})

	t.Run("invalid theme_mode is rejected", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		_, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
			ThemePreference: "dark",
			ThemeMode:       optimus-ide-collabsdk.ThemeMode("wizard"),
			TerminalFont:    optimus-ide-collabsdk.TerminalFontGeistMono,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})

	t.Run("invalid theme slots are rejected", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		for _, tc := range []struct {
			name       string
			themeMode  optimus-ide-collabsdk.ThemeMode
			themeLight string
			themeDark  string
		}{
			{
				name:       "arbitrary light slot",
				themeMode:  optimus-ide-collabsdk.ThemeModeSync,
				themeLight: "../../etc/passwd",
				themeDark:  "dark",
			},
			{
				name:       "arbitrary dark slot",
				themeMode:  optimus-ide-collabsdk.ThemeModeSync,
				themeLight: "light",
				themeDark:  "xss-payload",
			},
			{
				name:       "empty light slot in sync mode",
				themeMode:  optimus-ide-collabsdk.ThemeModeSync,
				themeLight: "",
				themeDark:  "dark",
			},
			{
				name:       "empty dark slot in sync mode",
				themeMode:  optimus-ide-collabsdk.ThemeModeSync,
				themeLight: "light",
				themeDark:  "",
			},
			{
				name:       "arbitrary light slot in single mode",
				themeMode:  optimus-ide-collabsdk.ThemeModeSingle,
				themeLight: "../../etc/passwd",
			},
			{
				name:      "arbitrary dark slot with omitted mode",
				themeDark: "xss-payload",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				_, err := client.UpdateUserAppearanceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest{
					ThemePreference: "dark",
					ThemeMode:       tc.themeMode,
					ThemeLight:      tc.themeLight,
					ThemeDark:       tc.themeDark,
					TerminalFont:    optimus-ide-collabsdk.TerminalFontGeistMono,
				})
				var apiErr *optimus-ide-collabsdk.Error
				require.ErrorAs(t, err, &apiErr)
				require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
			})
		}
	})
}

func TestUserTaskNotificationAlertDismissed(t *testing.T) {
	t.Parallel()

	// Single instance shared across all sub-tests. Each sub-test
	// creates its own non-admin user for isolation.
	adminClient := optimus-ide-collabdtest.New(t, nil)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

	t.Run("defaults to false", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		// When: getting user preference settings for a user
		settings, err := client.GetUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)

		// Then: the task notification alert dismissed should default to false
		require.False(t, settings.TaskNotificationAlertDismissed)
	})

	t.Run("update to true", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		// When: user dismisses the task notification alert
		updated, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			TaskNotificationAlertDismissed: ptr.Ref(true),
		})
		require.NoError(t, err)

		// Then: the setting is updated to true
		require.True(t, updated.TaskNotificationAlertDismissed)
	})

	t.Run("update to false", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		// Given: user has dismissed the task notification alert
		_, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			TaskNotificationAlertDismissed: ptr.Ref(true),
		})
		require.NoError(t, err)

		// When: the task notification alert dismissal is cleared
		// (e.g., when user enables a task notification in the UI settings)
		updated, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			TaskNotificationAlertDismissed: ptr.Ref(false),
		})
		require.NoError(t, err)

		// Then: the setting is updated to false
		require.False(t, updated.TaskNotificationAlertDismissed)
	})
}

func TestThinkingDisplayMode(t *testing.T) {
	t.Parallel()

	adminClient := optimus-ide-collabdtest.New(t, nil)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

	t.Run("defaults to auto", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		settings, err := client.GetUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ThinkingDisplayModeAuto, settings.ThinkingDisplayMode)
	})

	t.Run("round-trips a valid mode", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		updated, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			ThinkingDisplayMode: optimus-ide-collabsdk.ThinkingDisplayModeAlwaysCollapsed,
		})
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ThinkingDisplayModeAlwaysCollapsed, updated.ThinkingDisplayMode)

		settings, err := client.GetUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ThinkingDisplayModeAlwaysCollapsed, settings.ThinkingDisplayMode)
	})

	t.Run("rejects invalid mode", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		_, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			ThinkingDisplayMode: "bogus",
		})
		var sdkErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &sdkErr)
		require.Equal(t, http.StatusBadRequest, sdkErr.StatusCode())
	})

	t.Run("empty mode preserves stored value", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		// Set a non-default mode.
		_, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			ThinkingDisplayMode: optimus-ide-collabsdk.ThinkingDisplayModePreview,
		})
		require.NoError(t, err)

		// Send an update that omits thinking_display_mode (zero value).
		updated, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			TaskNotificationAlertDismissed: ptr.Ref(true),
		})
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ThinkingDisplayModePreview, updated.ThinkingDisplayMode)
	})
}

func TestAgentChatSendShortcutPreference(t *testing.T) {
	t.Parallel()

	adminClient := optimus-ide-collabdtest.New(t, nil)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

	requireValidationField := func(t *testing.T, err error, field string) {
		t.Helper()

		var sdkErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &sdkErr)
		require.Equal(t, http.StatusBadRequest, sdkErr.StatusCode())
		require.Len(t, sdkErr.Validations, 1)
		require.Equal(t, field, sdkErr.Validations[0].Field)
	}

	t.Run("defaults to enter", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		settings, err := client.GetUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.AgentChatSendShortcutEnter, settings.AgentChatSendShortcut)
	})

	t.Run("round-trips shortcut", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		updated, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			AgentChatSendShortcut: optimus-ide-collabsdk.AgentChatSendShortcutModifierEnter,
		})
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.AgentChatSendShortcutModifierEnter, updated.AgentChatSendShortcut)

		settings, err := client.GetUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.AgentChatSendShortcutModifierEnter, settings.AgentChatSendShortcut)
	})

	t.Run("rejects invalid shortcut", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		_, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			AgentChatSendShortcut: optimus-ide-collabsdk.AgentChatSendShortcut("bogus"),
		})
		requireValidationField(t, err, "agent_chat_send_shortcut")
	})

	t.Run("updates preserve stored shortcut", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		_, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			AgentChatSendShortcut: optimus-ide-collabsdk.AgentChatSendShortcutModifierEnter,
			ThinkingDisplayMode:   optimus-ide-collabsdk.ThinkingDisplayModePreview,
		})
		require.NoError(t, err)

		updated, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			ThinkingDisplayMode: optimus-ide-collabsdk.ThinkingDisplayModeAlwaysExpanded,
		})
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ThinkingDisplayModeAlwaysExpanded, updated.ThinkingDisplayMode)
		require.Equal(t, optimus-ide-collabsdk.AgentChatSendShortcutModifierEnter, updated.AgentChatSendShortcut)
	})
}

func TestAgentDisplayModePreferences(t *testing.T) {
	t.Parallel()

	adminClient := optimus-ide-collabdtest.New(t, nil)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, adminClient)

	requireValidationField := func(t *testing.T, err error, field string) {
		t.Helper()

		var sdkErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &sdkErr)
		require.Equal(t, http.StatusBadRequest, sdkErr.StatusCode())
		require.Len(t, sdkErr.Validations, 1)
		require.Equal(t, field, sdkErr.Validations[0].Field)
	}

	t.Run("defaults shell tools to always collapsed", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		settings, err := client.GetUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.AgentDisplayModeAlwaysCollapsed, settings.ShellToolDisplayMode)
		require.Empty(t, settings.CodeDiffDisplayMode)
	})

	t.Run("round-trips shell tool display mode", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		for _, mode := range []optimus-ide-collabsdk.AgentDisplayMode{
			optimus-ide-collabsdk.AgentDisplayModeAlwaysExpanded,
			optimus-ide-collabsdk.AgentDisplayModeAlwaysCollapsed,
		} {
			updated, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
				ShellToolDisplayMode: mode,
			})
			require.NoError(t, err)
			require.Equal(t, mode, updated.ShellToolDisplayMode)

			settings, err := client.GetUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me)
			require.NoError(t, err)
			require.Equal(t, mode, settings.ShellToolDisplayMode)
		}
	})

	t.Run("round-trips code diff display mode", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		for _, mode := range []optimus-ide-collabsdk.AgentDisplayMode{
			optimus-ide-collabsdk.AgentDisplayModeAlwaysExpanded,
			optimus-ide-collabsdk.AgentDisplayModeAlwaysCollapsed,
		} {
			updated, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
				CodeDiffDisplayMode: mode,
			})
			require.NoError(t, err)
			require.Equal(t, mode, updated.CodeDiffDisplayMode)

			settings, err := client.GetUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me)
			require.NoError(t, err)
			require.Equal(t, mode, settings.CodeDiffDisplayMode)
		}
	})

	t.Run("updates preserve stored display modes", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		_, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			ThinkingDisplayMode:  optimus-ide-collabsdk.ThinkingDisplayModePreview,
			ShellToolDisplayMode: optimus-ide-collabsdk.AgentDisplayModeAlwaysCollapsed,
			CodeDiffDisplayMode:  optimus-ide-collabsdk.AgentDisplayModeAlwaysExpanded,
		})
		require.NoError(t, err)

		updated, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
			ShellToolDisplayMode: optimus-ide-collabsdk.AgentDisplayModeAlwaysExpanded,
		})
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ThinkingDisplayModePreview, updated.ThinkingDisplayMode)
		require.Equal(t, optimus-ide-collabsdk.AgentDisplayModeAlwaysExpanded, updated.ShellToolDisplayMode)
		require.Equal(t, optimus-ide-collabsdk.AgentDisplayModeAlwaysExpanded, updated.CodeDiffDisplayMode)

		settings, err := client.GetUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, optimus-ide-collabsdk.ThinkingDisplayModePreview, settings.ThinkingDisplayMode)
		require.Equal(t, optimus-ide-collabsdk.AgentDisplayModeAlwaysExpanded, settings.ShellToolDisplayMode)
		require.Equal(t, optimus-ide-collabsdk.AgentDisplayModeAlwaysExpanded, settings.CodeDiffDisplayMode)
	})

	t.Run("rejects invalid shell tool display mode", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		for _, tt := range []struct {
			name string
			mode optimus-ide-collabsdk.AgentDisplayMode
		}{
			{
				name: "bogus",
				mode: optimus-ide-collabsdk.AgentDisplayMode("bogus"),
			},
			{
				name: "thinking preview",
				mode: optimus-ide-collabsdk.AgentDisplayMode(optimus-ide-collabsdk.ThinkingDisplayModePreview),
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				_, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
					ShellToolDisplayMode: tt.mode,
				})
				requireValidationField(t, err, "shell_tool_display_mode")
			})
		}
	})

	t.Run("rejects invalid code diff display mode", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, firstUser.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
		defer cancel()

		for _, tt := range []struct {
			name string
			mode optimus-ide-collabsdk.AgentDisplayMode
		}{
			{
				name: "bogus",
				mode: optimus-ide-collabsdk.AgentDisplayMode("bogus"),
			},
			{
				name: "thinking preview",
				mode: optimus-ide-collabsdk.AgentDisplayMode(optimus-ide-collabsdk.ThinkingDisplayModePreview),
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				_, err := client.UpdateUserPreferenceSettings(ctx, optimus-ide-collabsdk.Me, optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest{
					CodeDiffDisplayMode: tt.mode,
				})
				requireValidationField(t, err, "code_diff_display_mode")
			})
		}
	})
}

func TestWorkspacesByUser(t *testing.T) {
	t.Parallel()
	t.Run("Empty", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		res, err := client.Workspaces(ctx, optimus-ide-collabsdk.WorkspaceFilter{
			Owner: optimus-ide-collabsdk.Me,
		})
		require.NoError(t, err)
		require.Len(t, res.Workspaces, 0)
	})
	t.Run("Access", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		newUser, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           "test@optimus-ide-collab.com",
			Username:        "someone",
			Password:        "MySecurePassword!",
			OrganizationIDs: []uuid.UUID{user.OrganizationID},
		})
		require.NoError(t, err)
		auth, err := client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
			Email:    newUser.Email,
			Password: "MySecurePassword!",
		})
		require.NoError(t, err)

		newUserClient := optimus-ide-collabsdk.New(client.URL)
		newUserClient.SetSessionToken(auth.SessionToken)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)

		res, err := newUserClient.Workspaces(ctx, optimus-ide-collabsdk.WorkspaceFilter{Owner: optimus-ide-collabsdk.Me})
		require.NoError(t, err)
		require.Len(t, res.Workspaces, 0)

		res, err = client.Workspaces(ctx, optimus-ide-collabsdk.WorkspaceFilter{Owner: optimus-ide-collabsdk.Me})
		require.NoError(t, err)
		require.Len(t, res.Workspaces, 1)
	})
}

func TestDormantUser(t *testing.T) {
	t.Parallel()

	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
	user := optimus-ide-collabdtest.CreateFirstUser(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	// Create a new user
	newUser, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
		Email:           "test@optimus-ide-collab.com",
		Username:        "someone",
		Password:        "MySecurePassword!",
		OrganizationIDs: []uuid.UUID{user.OrganizationID},
	})
	require.NoError(t, err)

	// User should be dormant as they haven't logged in yet
	users, err := client.Users(ctx, optimus-ide-collabsdk.UsersRequest{Search: newUser.Username})
	require.NoError(t, err)
	require.Len(t, users.Users, 1)
	require.Equal(t, optimus-ide-collabsdk.UserStatusDormant, users.Users[0].Status)

	// User logs in now
	_, err = client.LoginWithPassword(ctx, optimus-ide-collabsdk.LoginWithPasswordRequest{
		Email:    newUser.Email,
		Password: "MySecurePassword!",
	})
	require.NoError(t, err)

	// User status should be active now
	users, err = client.Users(ctx, optimus-ide-collabsdk.UsersRequest{Search: newUser.Username})
	require.NoError(t, err)
	require.Len(t, users.Users, 1)
	require.Equal(t, optimus-ide-collabsdk.UserStatusActive, users.Users[0].Status)
}

// TestSuspendedPagination is when the after_id is a suspended record.
// The database query should still return the correct page, as the after_id
// is in a subquery that finds the record regardless of its status.
// This is mainly to confirm the db fake has the same behavior.
func TestSuspendedPagination(t *testing.T) {
	t.Parallel()
	t.Skip("This fails when two users are created at the exact same time. The reason is unknown... See: https://github.com/optimus-ide-collab/optimus-ide-collab/actions/runs/3057047622/jobs/4931863163")
	client := optimus-ide-collabdtest.New(t, nil)
	optimus-ide-collabdtest.CreateFirstUser(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	t.Cleanup(cancel)

	me, err := client.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)
	orgID := me.OrganizationIDs[0]

	total := 10
	users := make([]optimus-ide-collabsdk.User, 0, total)
	// Create users
	for i := 0; i < total; i++ {
		email := fmt.Sprintf("%d@optimus-ide-collab.com", i)
		username := fmt.Sprintf("user%d", i)
		user, err := client.CreateUserWithOrgs(ctx, optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:           email,
			Username:        username,
			Password:        "MySecurePassword!",
			OrganizationIDs: []uuid.UUID{orgID},
		})
		require.NoError(t, err)
		users = append(users, user)
	}
	sortUsers(users)
	deletedUser := users[2]
	expected := users[3:8]
	_, err = client.UpdateUserStatus(ctx, deletedUser.ID.String(), optimus-ide-collabsdk.UserStatusSuspended)
	require.NoError(t, err, "suspend user")

	page, err := client.Users(ctx, optimus-ide-collabsdk.UsersRequest{
		Pagination: optimus-ide-collabsdk.Pagination{
			Limit:   len(expected),
			AfterID: deletedUser.ID,
		},
	})
	require.NoError(t, err)
	require.Equal(t, expected, page.Users, "expected page")
}

func TestUserAutofillParameters(t *testing.T) {
	t.Parallel()
	t.Run("NotSelf", func(t *testing.T) {
		t.Parallel()
		client1, _, api := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{})

		u1 := optimus-ide-collabdtest.CreateFirstUser(t, client1)

		client2, u2 := optimus-ide-collabdtest.CreateAnotherUser(t, client1, u1.OrganizationID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		db := api.Database

		version := dbfake.TemplateVersion(t, db).Seed(database.TemplateVersion{
			CreatedBy:      u1.UserID,
			OrganizationID: u1.OrganizationID,
		}).Params(database.TemplateVersionParameter{
			Name:     "param",
			Required: true,
		}).Do()

		_, err := client2.UserAutofillParameters(
			ctx,
			u1.UserID.String(),
			version.Template.ID,
		)

		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())

		// u1 should be able to read u2's parameters as u1 is site admin.
		_, err = client1.UserAutofillParameters(
			ctx,
			u2.ID.String(),
			version.Template.ID,
		)
		require.NoError(t, err)
	})

	t.Run("FindsParameters", func(t *testing.T) {
		t.Parallel()
		client1, _, api := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{})

		u1 := optimus-ide-collabdtest.CreateFirstUser(t, client1)

		client2, u2 := optimus-ide-collabdtest.CreateAnotherUser(t, client1, u1.OrganizationID)

		db := api.Database

		version := dbfake.TemplateVersion(t, db).Seed(database.TemplateVersion{
			CreatedBy:      u1.UserID,
			OrganizationID: u1.OrganizationID,
		}).Params(database.TemplateVersionParameter{
			Name:     "param",
			Required: true,
		},
			database.TemplateVersionParameter{
				Name:      "param2",
				Ephemeral: true,
			},
		).Do()

		dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        u2.ID,
			TemplateID:     version.Template.ID,
			OrganizationID: u1.OrganizationID,
		}).Params(
			database.WorkspaceBuildParameter{
				Name:  "param",
				Value: "foo",
			},
			database.WorkspaceBuildParameter{
				Name:  "param2",
				Value: "bar",
			},
		).Do()

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		// Use client2 since client1 is site admin, so
		// we don't get good coverage on RBAC working.
		params, err := client2.UserAutofillParameters(
			ctx,
			u2.ID.String(),
			version.Template.ID,
		)
		require.NoError(t, err)

		require.Equal(t, 1, len(params))

		require.Equal(t, "param", params[0].Name)
		require.Equal(t, "foo", params[0].Value)

		// Verify that latest parameter value is returned.
		dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OrganizationID: u1.OrganizationID,
			OwnerID:        u2.ID,
			TemplateID:     version.Template.ID,
		}).Params(
			database.WorkspaceBuildParameter{
				Name:  "param",
				Value: "foo_new",
			},
		).Do()

		params, err = client2.UserAutofillParameters(
			ctx,
			u2.ID.String(),
			version.Template.ID,
		)
		require.NoError(t, err)

		require.Equal(t, 1, len(params))

		require.Equal(t, "param", params[0].Name)
		require.Equal(t, "foo_new", params[0].Value)
	})
}

// TestPaginatedUsers creates a list of users, then tries to paginate through
// them using different page sizes.
func TestPaginatedUsers(t *testing.T) {
	t.Parallel()
	client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
	optimus-ide-collabdtest.CreateFirstUser(t, client)

	// This test takes longer than a long time.
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong*4)
	t.Cleanup(cancel)

	me, err := client.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)

	// When 50 users exist
	total := 50
	allUsers := make([]database.User, total+1)
	allUsers[0] = database.User{
		Email:    me.Email,
		Username: me.Username,
	}
	specialUsers := make([]database.User, total/2)

	eg, _ := errgroup.WithContext(ctx)
	// Create users
	for i := 0; i < total; i++ {
		eg.Go(func() error {
			email := fmt.Sprintf("%d@optimus-ide-collab.com", i)
			username := fmt.Sprintf("user%d", i)
			if i%2 == 0 {
				email = fmt.Sprintf("%d@gmail.com", i)
				username = fmt.Sprintf("specialuser%d", i)
			}
			if i%3 == 0 {
				username = strings.ToUpper(username)
			}

			// We used to use the API to ceate users, but that is slow.
			// Instead, we create them directly in the database now
			// to prevent timeout flakes.
			newUser := dbgen.User(t, db, database.User{
				Email:    email,
				Username: username,
			})
			allUsers[i+1] = newUser
			if i%2 == 0 {
				specialUsers[i/2] = newUser
			}

			return nil
		})
	}
	err = eg.Wait()
	require.NoError(t, err, "create users failed")

	// Sorting the users will sort by username.
	sortDatabaseUsers(allUsers)
	sortDatabaseUsers(specialUsers)

	gmailSearch := func(request optimus-ide-collabsdk.UsersRequest) optimus-ide-collabsdk.UsersRequest {
		request.Search = "gmail"
		return request
	}
	usernameSearch := func(request optimus-ide-collabsdk.UsersRequest) optimus-ide-collabsdk.UsersRequest {
		request.Search = "specialuser"
		return request
	}

	tests := []struct {
		name     string
		limit    int
		allUsers []database.User
		opt      func(request optimus-ide-collabsdk.UsersRequest) optimus-ide-collabsdk.UsersRequest
	}{
		{name: "all users", limit: 10, allUsers: allUsers},
		{name: "all users", limit: 5, allUsers: allUsers},
		{name: "all users", limit: 3, allUsers: allUsers},
		{name: "gmail search", limit: 3, allUsers: specialUsers, opt: gmailSearch},
		{name: "gmail search", limit: 7, allUsers: specialUsers, opt: gmailSearch},
		{name: "username search", limit: 3, allUsers: specialUsers, opt: usernameSearch},
		{name: "username search", limit: 3, allUsers: specialUsers, opt: usernameSearch},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %d", tt.name, tt.limit), func(t *testing.T) {
			t.Parallel()

			// This test takes longer than a long time.
			ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong*2)
			defer cancel()

			assertPagination(ctx, t, client, tt.limit, tt.allUsers, tt.opt)
		})
	}
}

// Assert pagination will page through the list of all users using the given
// limit for each page. The 'allUsers' is the expected full list to compare
// against.
func assertPagination(ctx context.Context, t *testing.T, client *optimus-ide-collabsdk.Client, limit int, allUsers []database.User,
	opt func(request optimus-ide-collabsdk.UsersRequest) optimus-ide-collabsdk.UsersRequest,
) {
	var count int
	if opt == nil {
		opt = func(request optimus-ide-collabsdk.UsersRequest) optimus-ide-collabsdk.UsersRequest {
			return request
		}
	}

	// Check the first page
	page, err := client.Users(ctx, opt(optimus-ide-collabsdk.UsersRequest{
		Pagination: optimus-ide-collabsdk.Pagination{
			Limit: limit,
		},
	}))
	require.NoError(t, err, "first page")
	require.Equalf(t, onlyUsernames(page.Users), onlyUsernames(allUsers[:limit]), "first page, limit=%d", limit)
	count += len(page.Users)

	for {
		if len(page.Users) == 0 {
			break
		}

		afterCursor := page.Users[len(page.Users)-1].ID
		// Assert each page is the next expected page
		// This is using a cursor, and only works if all users created_at
		// is unique.
		page, err = client.Users(ctx, opt(optimus-ide-collabsdk.UsersRequest{
			Pagination: optimus-ide-collabsdk.Pagination{
				Limit:   limit,
				AfterID: afterCursor,
			},
		}))
		require.NoError(t, err, "next cursor page")

		// Also check page by offset
		offsetPage, err := client.Users(ctx, opt(optimus-ide-collabsdk.UsersRequest{
			Pagination: optimus-ide-collabsdk.Pagination{
				Limit:  limit,
				Offset: count,
			},
		}))
		require.NoError(t, err, "next offset page")

		var expected []database.User
		if count+limit > len(allUsers) {
			expected = allUsers[count:]
		} else {
			expected = allUsers[count : count+limit]
		}
		require.Equalf(t, onlyUsernames(page.Users), onlyUsernames(expected), "next users, after=%s, limit=%d", afterCursor, limit)
		require.Equalf(t, onlyUsernames(offsetPage.Users), onlyUsernames(expected), "offset users, offset=%d, limit=%d", count, limit)

		// Also check the before
		prevPage, err := client.Users(ctx, opt(optimus-ide-collabsdk.UsersRequest{
			Pagination: optimus-ide-collabsdk.Pagination{
				Offset: count - limit,
				Limit:  limit,
			},
		}))
		require.NoError(t, err, "prev page")
		require.Equal(t, onlyUsernames(allUsers[count-limit:count]), onlyUsernames(prevPage.Users), "prev users")
		count += len(page.Users)
	}
}

// sortUsers sorts by (created_at, id)
func sortUsers(users []optimus-ide-collabsdk.User) {
	slices.SortFunc(users, func(a, b optimus-ide-collabsdk.User) int {
		return slice.Ascending(strings.ToLower(a.Username), strings.ToLower(b.Username))
	})
}

func sortDatabaseUsers(users []database.User) {
	slices.SortFunc(users, func(a, b database.User) int {
		return slice.Ascending(strings.ToLower(a.Username), strings.ToLower(b.Username))
	})
}

func onlyUsernames[U optimus-ide-collabsdk.User | database.User](users []U) []string {
	var out []string
	for _, u := range users {
		switch u := (any(u)).(type) {
		case optimus-ide-collabsdk.User:
			out = append(out, u.Username)
		case database.User:
			out = append(out, u.Username)
		}
	}
	return out
}

func BenchmarkUsersMe(b *testing.B) {
	client := optimus-ide-collabdtest.New(b, nil)
	_ = optimus-ide-collabdtest.CreateFirstUser(b, client)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.User(ctx, optimus-ide-collabsdk.Me)
		require.NoError(b, err)
	}
}
