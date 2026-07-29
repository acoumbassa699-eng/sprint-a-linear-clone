package optimus-ide-collabdtest

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/db2sdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/userpassword"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// UsersPagination creates a set of users for testing pagination.  It can be
// used to test paginating both users and group members.
func UsersPagination(
	ctx context.Context,
	t *testing.T,
	client *optimus-ide-collabsdk.Client,
	setup func(users []optimus-ide-collabsdk.User),
	fetch func(req optimus-ide-collabsdk.UsersRequest) ([]optimus-ide-collabsdk.ReducedUser, int),
) {
	t.Helper()

	firstUser, err := client.User(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err, "fetch me")

	count := 10
	users := make([]optimus-ide-collabsdk.User, count)
	orgID := firstUser.OrganizationIDs[0]
	users[0] = firstUser
	for i := range count - 1 {
		_, user := CreateAnotherUserMutators(t, client, orgID, nil, func(r *optimus-ide-collabsdk.CreateUserRequestWithOrgs) {
			if i < 5 {
				r.Name = fmt.Sprintf("before%d", i)
			} else {
				r.Name = fmt.Sprintf("after%d", i)
			}
		})
		users[i+1] = user
	}

	slices.SortFunc(users, func(a, b optimus-ide-collabsdk.User) int {
		return slice.Ascending(strings.ToLower(a.Username), strings.ToLower(b.Username))
	})

	if setup != nil {
		setup(users)
	}

	gotUsers, gotCount := fetch(optimus-ide-collabsdk.UsersRequest{})
	require.Len(t, gotUsers, count)
	require.Equal(t, gotCount, count)

	gotUsers, gotCount = fetch(optimus-ide-collabsdk.UsersRequest{
		Pagination: optimus-ide-collabsdk.Pagination{
			Limit: 1,
		},
	})
	require.Len(t, gotUsers, 1)
	require.Equal(t, gotCount, count)

	gotUsers, gotCount = fetch(optimus-ide-collabsdk.UsersRequest{
		Pagination: optimus-ide-collabsdk.Pagination{
			Offset: 1,
		},
	})
	require.Len(t, gotUsers, count-1)
	require.Equal(t, gotCount, count)

	gotUsers, gotCount = fetch(optimus-ide-collabsdk.UsersRequest{
		Pagination: optimus-ide-collabsdk.Pagination{
			Limit:  1,
			Offset: 1,
		},
	})
	require.Len(t, gotUsers, 1)
	require.Equal(t, gotCount, count)

	// If offset is higher than the count postgres returns an empty array
	// and not an ErrNoRows error.
	gotUsers, gotCount = fetch(optimus-ide-collabsdk.UsersRequest{
		Pagination: optimus-ide-collabsdk.Pagination{
			Offset: count + 1,
		},
	})
	require.Len(t, gotUsers, 0)
	require.Equal(t, gotCount, 0)

	// Check that AfterID works.
	gotUsers, gotCount = fetch(optimus-ide-collabsdk.UsersRequest{
		Pagination: optimus-ide-collabsdk.Pagination{
			AfterID: users[5].ID,
		},
	})
	require.NoError(t, err)
	require.Len(t, gotUsers, 4)
	require.Equal(t, gotCount, 4)

	// Check we can paginate a filtered response.
	gotUsers, gotCount = fetch(optimus-ide-collabsdk.UsersRequest{
		SearchQuery: "name:after",
		Pagination: optimus-ide-collabsdk.Pagination{
			Limit:  1,
			Offset: 1,
		},
	})
	require.NoError(t, err)
	require.Len(t, gotUsers, 1)
	require.Equal(t, gotCount, 4)
	require.Contains(t, gotUsers[0].Name, "after")
}

type UsersFilterOptions struct {
	CreateServiceAccounts bool
}

// UsersFilter creates a set of users to run various filters against for
// testing.  It can be used to test filtering both users and group members.
func UsersFilter(
	setupCtx context.Context,
	t *testing.T,
	client *optimus-ide-collabsdk.Client,
	db database.Store,
	options *UsersFilterOptions,
	setup func(users []optimus-ide-collabsdk.User),
	fetch func(ctx context.Context, req optimus-ide-collabsdk.UsersRequest) []optimus-ide-collabsdk.ReducedUser,
) {
	t.Helper()

	if options == nil {
		options = &UsersFilterOptions{}
	}

	firstUser, err := client.User(setupCtx, optimus-ide-collabsdk.Me)
	require.NoError(t, err, "fetch me")

	// Noon on Jan 18 is the "now" for this test for last_seen timestamps.
	// All these values are equal
	// 2023-01-18T12:00:00Z (UTC)
	// 2023-01-18T07:00:00-05:00 (America/New_York)
	// 2023-01-18T13:00:00+01:00 (Europe/Madrid)
	// 2023-01-16T00:00:00+12:00 (Asia/Anadyr)
	lastSeenNow := time.Date(2023, 1, 18, 12, 0, 0, 0, time.UTC)
	users := make([]optimus-ide-collabsdk.User, 0)
	users = append(users, firstUser)
	orgID := firstUser.OrganizationIDs[0]
	githubIDs := make(map[int]uuid.UUID)
	for i := range 15 {
		roles := []rbac.RoleIdentifier{}
		if i%2 == 0 {
			roles = append(roles, rbac.RoleTemplateAdmin(), rbac.RoleUserAdmin())
		}
		if i%3 == 0 {
			roles = append(roles, rbac.RoleAuditor())
		}
		userClient, userData := CreateAnotherUserMutators(t, client, orgID, roles, func(r *optimus-ide-collabsdk.CreateUserRequestWithOrgs) {
			switch {
			case i%7 == 0:
				r.UserLoginType = optimus-ide-collabsdk.LoginTypeGithub
				r.Password = ""
			case i%6 == 0:
				r.UserLoginType = optimus-ide-collabsdk.LoginTypeOIDC
				r.Password = ""
			default:
				r.UserLoginType = optimus-ide-collabsdk.LoginTypePassword
			}
		})

		// Set the last seen for each user to a unique day
		// nolint:gocritic // Setting up unit test data.
		_, err := db.UpdateUserLastSeenAt(dbauthz.AsSystemRestricted(setupCtx), database.UpdateUserLastSeenAtParams{
			ID:         userData.ID,
			LastSeenAt: lastSeenNow.Add(-1 * time.Hour * 24 * time.Duration(i)),
			UpdatedAt:  time.Now(),
		})
		require.NoError(t, err, "set a last seen")

		// Set a github user ID for github login types.
		if i%7 == 0 {
			// nolint:gocritic // Setting up unit test data.
			err = db.UpdateUserGithubComUserID(dbauthz.AsSystemRestricted(setupCtx), database.UpdateUserGithubComUserIDParams{
				ID: userData.ID,
				GithubComUserID: sql.NullInt64{
					Int64: int64(i),
					Valid: true,
				},
			})
			require.NoError(t, err)
			githubIDs[i] = userData.ID
		}

		user, err := userClient.User(setupCtx, optimus-ide-collabsdk.Me)
		require.NoError(t, err, "fetch me")

		if i%4 == 0 {
			user, err = client.UpdateUserStatus(setupCtx, user.ID.String(), optimus-ide-collabsdk.UserStatusSuspended)
			require.NoError(t, err, "suspend user")
		}

		if i%5 == 0 {
			user, err = client.UpdateUserProfile(setupCtx, user.ID.String(), optimus-ide-collabsdk.UpdateUserProfileRequest{
				Username: strings.ToUpper(user.Username),
			})
			require.NoError(t, err, "update username to uppercase")
		}

		users = append(users, user)
	}

	// Add some service accounts.
	if options.CreateServiceAccounts {
		for range 3 {
			_, user := CreateAnotherUserMutators(t, client, orgID, nil, func(r *optimus-ide-collabsdk.CreateUserRequestWithOrgs) {
				r.ServiceAccount = true
			})
			users = append(users, user)
		}
	}

	hashedPassword, err := userpassword.Hash("SomeStrongPassword!")
	require.NoError(t, err)

	// Add users with different creation dates for testing date filters
	for i := range 3 {
		// nolint:gocritic // Setting up unit test data.
		user1, err := db.InsertUser(dbauthz.AsSystemRestricted(setupCtx), database.InsertUserParams{
			ID:               uuid.New(),
			Email:            fmt.Sprintf("before%d@optimus-ide-collab.com", i),
			Username:         fmt.Sprintf("before%d", i),
			Name:             fmt.Sprintf("Test User %d", i),
			HashedPassword:   []byte(hashedPassword),
			LoginType:        database.LoginTypeNone,
			Status:           string(optimus-ide-collabsdk.UserStatusActive),
			RBACRoles:        []string{optimus-ide-collabsdk.RoleMember},
			CreatedAt:        dbtime.Time(time.Date(2022, 12, 15+i, 12, 0, 0, 0, time.UTC)),
			UpdatedAt:        dbtime.Time(time.Date(2022, 12, 15+i, 12, 0, 0, 0, time.UTC)),
			IsServiceAccount: false,
		})
		require.NoError(t, err)
		// nolint:gocritic // Setting up unit test data.
		_, err = db.InsertOrganizationMember(dbauthz.AsSystemRestricted(setupCtx), database.InsertOrganizationMemberParams{
			OrganizationID: orgID,
			UserID:         user1.ID,
			CreatedAt:      dbtime.Now(),
			UpdatedAt:      dbtime.Now(),
			Roles:          []string{},
		})
		require.NoError(t, err)

		// The expected timestamps must be parsed from strings to compare equal during `ElementsMatch`
		sdkUser1 := db2sdk.User(user1, []uuid.UUID{orgID})
		sdkUser1.CreatedAt, err = time.Parse(time.RFC3339, sdkUser1.CreatedAt.Format(time.RFC3339))
		require.NoError(t, err)
		sdkUser1.UpdatedAt, err = time.Parse(time.RFC3339, sdkUser1.UpdatedAt.Format(time.RFC3339))
		require.NoError(t, err)
		sdkUser1.LastSeenAt, err = time.Parse(time.RFC3339, sdkUser1.LastSeenAt.Format(time.RFC3339))
		require.NoError(t, err)
		users = append(users, sdkUser1)

		// nolint:gocritic // Setting up unit test data.
		user2, err := db.InsertUser(dbauthz.AsSystemRestricted(setupCtx), database.InsertUserParams{
			ID:               uuid.New(),
			Email:            fmt.Sprintf("during%d@optimus-ide-collab.com", i),
			Username:         fmt.Sprintf("during%d", i),
			Name:             "",
			HashedPassword:   []byte(hashedPassword),
			LoginType:        database.LoginTypeNone,
			Status:           string(optimus-ide-collabsdk.UserStatusActive),
			RBACRoles:        []string{optimus-ide-collabsdk.RoleOwner},
			CreatedAt:        dbtime.Time(time.Date(2023, 1, 15+i, 12, 0, 0, 0, time.UTC)),
			UpdatedAt:        dbtime.Time(time.Date(2023, 1, 15+i, 12, 0, 0, 0, time.UTC)),
			IsServiceAccount: false,
		})
		require.NoError(t, err)
		// nolint:gocritic // Setting up unit test data.
		_, err = db.InsertOrganizationMember(dbauthz.AsSystemRestricted(setupCtx), database.InsertOrganizationMemberParams{
			OrganizationID: orgID,
			UserID:         user2.ID,
			CreatedAt:      dbtime.Now(),
			UpdatedAt:      dbtime.Now(),
			Roles:          []string{},
		})
		require.NoError(t, err)

		sdkUser2 := db2sdk.User(user2, []uuid.UUID{orgID})
		sdkUser2.CreatedAt, err = time.Parse(time.RFC3339, sdkUser2.CreatedAt.Format(time.RFC3339))
		require.NoError(t, err)
		sdkUser2.UpdatedAt, err = time.Parse(time.RFC3339, sdkUser2.UpdatedAt.Format(time.RFC3339))
		require.NoError(t, err)
		sdkUser2.LastSeenAt, err = time.Parse(time.RFC3339, sdkUser2.LastSeenAt.Format(time.RFC3339))
		require.NoError(t, err)
		users = append(users, sdkUser2)

		// nolint:gocritic // Setting up unit test data.
		user3, err := db.InsertUser(dbauthz.AsSystemRestricted(setupCtx), database.InsertUserParams{
			ID:               uuid.New(),
			Email:            fmt.Sprintf("after%d@optimus-ide-collab.com", i),
			Username:         fmt.Sprintf("after%d", i),
			Name:             "",
			HashedPassword:   []byte(hashedPassword),
			LoginType:        database.LoginTypeNone,
			Status:           string(optimus-ide-collabsdk.UserStatusActive),
			RBACRoles:        []string{optimus-ide-collabsdk.RoleOwner},
			CreatedAt:        dbtime.Time(time.Date(2023, 2, 15+i, 12, 0, 0, 0, time.UTC)),
			UpdatedAt:        dbtime.Time(time.Date(2023, 2, 15+i, 12, 0, 0, 0, time.UTC)),
			IsServiceAccount: false,
		})
		require.NoError(t, err)
		// nolint:gocritic // Setting up unit test data.
		_, err = db.InsertOrganizationMember(dbauthz.AsSystemRestricted(setupCtx), database.InsertOrganizationMemberParams{
			OrganizationID: orgID,
			UserID:         user3.ID,
			CreatedAt:      dbtime.Now(),
			UpdatedAt:      dbtime.Now(),
			Roles:          []string{},
		})
		require.NoError(t, err)

		sdkUser3 := db2sdk.User(user3, []uuid.UUID{orgID})
		sdkUser3.CreatedAt, err = time.Parse(time.RFC3339, sdkUser3.CreatedAt.Format(time.RFC3339))
		require.NoError(t, err)
		sdkUser3.UpdatedAt, err = time.Parse(time.RFC3339, sdkUser3.UpdatedAt.Format(time.RFC3339))
		require.NoError(t, err)
		sdkUser3.LastSeenAt, err = time.Parse(time.RFC3339, sdkUser3.LastSeenAt.Format(time.RFC3339))
		require.NoError(t, err)
		users = append(users, sdkUser3)
	}

	if setup != nil {
		setup(users)
	}

	// --- Setup done ---
	testCases := []struct {
		Name   string
		Filter optimus-ide-collabsdk.UsersRequest
		// If FilterF is true, we include it in the expected results
		FilterF func(f optimus-ide-collabsdk.UsersRequest, user optimus-ide-collabsdk.User) bool
	}{
		{
			Name: "All",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Status: optimus-ide-collabsdk.UserStatusSuspended + "," + optimus-ide-collabsdk.UserStatusActive,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, _ optimus-ide-collabsdk.User) bool {
				return true
			},
		},
		{
			Name: "Active",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Status: optimus-ide-collabsdk.UserStatusActive,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.Status == optimus-ide-collabsdk.UserStatusActive
			},
		},
		{
			Name: "GithubComUserID",
			Filter: optimus-ide-collabsdk.UsersRequest{
				SearchQuery: "github_com_user_id:7",
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.ID == githubIDs[7]
			},
		},
		{
			Name: "ActiveUppercase",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Status: "ACTIVE",
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.Status == optimus-ide-collabsdk.UserStatusActive
			},
		},
		{
			Name: "Suspended",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Status: optimus-ide-collabsdk.UserStatusSuspended,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.Status == optimus-ide-collabsdk.UserStatusSuspended
			},
		},
		{
			Name: "NameContains",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Search: "a",
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return (strings.ContainsAny(u.Username, "aA") || strings.ContainsAny(u.Email, "aA") || strings.ContainsAny(u.Name, "aA"))
			},
		},
		{
			Name: "DisplayNameSearch",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Search: "user",
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				const term = "user"
				return strings.Contains(strings.ToLower(u.Username), term) ||
					strings.Contains(strings.ToLower(u.Email), term) ||
					strings.Contains(strings.ToLower(u.Name), term)
			},
		},
		{
			Name: "NameAndSearch",
			Filter: optimus-ide-collabsdk.UsersRequest{
				SearchQuery: "name:Test search:before1",
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.Username == "before1"
			},
		},
		{
			Name: "NameNoMatch",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Search: "nonexistent",
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, _ optimus-ide-collabsdk.User) bool {
				return false
			},
		},
		{
			Name: "Admins",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Role:   optimus-ide-collabsdk.RoleOwner,
				Status: optimus-ide-collabsdk.UserStatusSuspended + "," + optimus-ide-collabsdk.UserStatusActive,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				for _, r := range u.Roles {
					if r.Name == optimus-ide-collabsdk.RoleOwner {
						return true
					}
				}
				return false
			},
		},
		{
			Name: "AdminsUppercase",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Role:   "OWNER",
				Status: optimus-ide-collabsdk.UserStatusSuspended + "," + optimus-ide-collabsdk.UserStatusActive,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				for _, r := range u.Roles {
					if r.Name == optimus-ide-collabsdk.RoleOwner {
						return true
					}
				}
				return false
			},
		},
		{
			Name: "Members",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Role:   optimus-ide-collabsdk.RoleMember,
				Status: optimus-ide-collabsdk.UserStatusSuspended + "," + optimus-ide-collabsdk.UserStatusActive,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, _ optimus-ide-collabsdk.User) bool {
				return true
			},
		},
		{
			Name: "SearchQuery",
			Filter: optimus-ide-collabsdk.UsersRequest{
				SearchQuery: "i role:owner status:active",
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				for _, r := range u.Roles {
					if r.Name == optimus-ide-collabsdk.RoleOwner {
						return (strings.ContainsAny(u.Username, "iI") || strings.ContainsAny(u.Email, "iI") || strings.ContainsAny(u.Name, "iI")) &&
							u.Status == optimus-ide-collabsdk.UserStatusActive
					}
				}
				return false
			},
		},
		{
			Name: "SearchQueryInsensitive",
			Filter: optimus-ide-collabsdk.UsersRequest{
				SearchQuery: "i Role:Owner STATUS:Active",
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				for _, r := range u.Roles {
					if r.Name == optimus-ide-collabsdk.RoleOwner {
						return (strings.ContainsAny(u.Username, "iI") || strings.ContainsAny(u.Email, "iI") || strings.ContainsAny(u.Name, "iI")) &&
							u.Status == optimus-ide-collabsdk.UserStatusActive
					}
				}
				return false
			},
		},
		{
			Name: "LastSeenBeforeNow",
			Filter: optimus-ide-collabsdk.UsersRequest{
				SearchQuery: `last_seen_before:"2023-01-16T00:00:00+12:00"`,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.LastSeenAt.Before(lastSeenNow)
			},
		},
		{
			Name: "LastSeenLastWeek",
			Filter: optimus-ide-collabsdk.UsersRequest{
				SearchQuery: `last_seen_before:"2023-01-14T23:59:59Z" last_seen_after:"2023-01-08T00:00:00Z"`,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				start := time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC)
				end := time.Date(2023, 1, 14, 23, 59, 59, 0, time.UTC)
				return u.LastSeenAt.Before(end) && u.LastSeenAt.After(start)
			},
		},
		{
			Name: "CreatedAtBefore",
			Filter: optimus-ide-collabsdk.UsersRequest{
				SearchQuery: `created_before:"2023-01-31T23:59:59Z"`,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				end := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)
				return u.CreatedAt.Before(end)
			},
		},
		{
			Name: "CreatedAtAfter",
			Filter: optimus-ide-collabsdk.UsersRequest{
				SearchQuery: `created_after:"2023-01-01T00:00:00Z"`,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				return u.CreatedAt.After(start)
			},
		},
		{
			Name: "CreatedAtRange",
			Filter: optimus-ide-collabsdk.UsersRequest{
				SearchQuery: `created_after:"2023-01-01T00:00:00Z" created_before:"2023-01-31T23:59:59Z"`,
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				end := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)
				return u.CreatedAt.After(start) && u.CreatedAt.Before(end)
			},
		},
		{
			Name: "LoginTypeNone",
			Filter: optimus-ide-collabsdk.UsersRequest{
				LoginType: []optimus-ide-collabsdk.LoginType{optimus-ide-collabsdk.LoginTypeNone},
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.LoginType == optimus-ide-collabsdk.LoginTypeNone
			},
		},
		{
			Name: "LoginTypeOIDC",
			Filter: optimus-ide-collabsdk.UsersRequest{
				LoginType: []optimus-ide-collabsdk.LoginType{optimus-ide-collabsdk.LoginTypeOIDC},
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.LoginType == optimus-ide-collabsdk.LoginTypeOIDC
			},
		},
		{
			Name: "LoginTypeMultiple",
			Filter: optimus-ide-collabsdk.UsersRequest{
				LoginType: []optimus-ide-collabsdk.LoginType{optimus-ide-collabsdk.LoginTypeNone, optimus-ide-collabsdk.LoginTypeGithub},
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.LoginType == optimus-ide-collabsdk.LoginTypeNone || u.LoginType == optimus-ide-collabsdk.LoginTypeGithub
			},
		},
		{
			Name: "DormantUserWithLoginTypeNone",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Status:    optimus-ide-collabsdk.UserStatusSuspended,
				LoginType: []optimus-ide-collabsdk.LoginType{optimus-ide-collabsdk.LoginTypeNone},
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.Status == optimus-ide-collabsdk.UserStatusSuspended && u.LoginType == optimus-ide-collabsdk.LoginTypeNone
			},
		},
		{
			Name: "IsServiceAccount",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Search: "service_account:true",
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return u.IsServiceAccount
			},
		},
		{
			Name: "IsNotServiceAccount",
			Filter: optimus-ide-collabsdk.UsersRequest{
				Search: "service_account:false",
			},
			FilterF: func(_ optimus-ide-collabsdk.UsersRequest, u optimus-ide-collabsdk.User) bool {
				return !u.IsServiceAccount
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			testCtx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
			defer cancel()

			got := fetch(testCtx, c.Filter)
			exp := make([]optimus-ide-collabsdk.ReducedUser, 0)
			for _, made := range users {
				match := c.FilterF(c.Filter, made)
				if match {
					exp = append(exp, made.ReducedUser)
				}
			}

			require.ElementsMatch(t, exp, got, "expected users returned")
		})
	}
}
