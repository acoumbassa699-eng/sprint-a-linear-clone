package optimus-ide-collabd_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestAddMember(t *testing.T) {
	t.Parallel()

	owner := optimus-ide-collabdtest.New(t, nil)
	first := optimus-ide-collabdtest.CreateFirstUser(t, owner)
	_, user := optimus-ide-collabdtest.CreateAnotherUser(t, owner, first.OrganizationID)

	t.Run("AlreadyMember", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		// Add user to org, even though they already exist
		// nolint:gocritic // must be an owner to see the user
		_, err := owner.PostOrganizationMember(ctx, first.OrganizationID, user.Username)
		require.ErrorContains(t, err, "already an organization member")

		org, err := owner.Organization(ctx, first.OrganizationID)
		require.NoError(t, err)

		member, err := owner.OrganizationMember(ctx, org.Name, user.Username)
		require.NoError(t, err)
		require.Equal(t, member.UserID, user.ID)
	})

	t.Run("Me", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)

		member, err := owner.OrganizationMember(ctx, first.OrganizationID.String(), optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.Equal(t, member.UserID, first.UserID)
	})
}

func TestDeleteMember(t *testing.T) {
	t.Parallel()

	t.Run("Allowed", func(t *testing.T) {
		t.Parallel()
		owner := optimus-ide-collabdtest.New(t, nil)
		first := optimus-ide-collabdtest.CreateFirstUser(t, owner)
		_, user := optimus-ide-collabdtest.CreateAnotherUser(t, owner, first.OrganizationID)

		ctx := testutil.Context(t, testutil.WaitMedium)
		// Deleting members from the default org is not allowed.
		// If this behavior changes, and we allow deleting members from the default org,
		// this test should be updated to check there is no error.
		// nolint:gocritic // must be an owner to see the user
		err := owner.DeleteOrganizationMember(ctx, first.OrganizationID, user.Username)
		require.NoError(t, err)
	})
}

func TestListMembers(t *testing.T) {
	t.Parallel()

	client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
	owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
	_, orgMember := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
	_, orgAdmin := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
	anotherOrg := dbgen.Organization(t, db, database.Organization{})
	anotherUser := dbgen.User(t, db, database.User{
		GithubComUserID: sql.NullInt64{Valid: true, Int64: 12345},
	})
	_ = dbgen.OrganizationMember(t, db, database.OrganizationMember{
		OrganizationID: anotherOrg.ID,
		UserID:         anotherUser.ID,
	})

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)
		members, err := client.OrganizationMembers(ctx, owner.OrganizationID)
		require.NoError(t, err)
		require.Len(t, members, 3)
		require.ElementsMatch(t,
			[]uuid.UUID{owner.UserID, orgMember.ID, orgAdmin.ID},
			slice.List(members, onlyIDs))
	})

	t.Run("UserID", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)
		members, err := client.OrganizationMembers(ctx, owner.OrganizationID, optimus-ide-collabsdk.OrganizationMembersQueryOptionUserID(orgMember.ID))
		require.NoError(t, err)
		require.Len(t, members, 1)
		require.ElementsMatch(t,
			[]uuid.UUID{orgMember.ID},
			slice.List(members, onlyIDs))
	})

	t.Run("IncludeSystem", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)
		members, err := client.OrganizationMembers(ctx, owner.OrganizationID, optimus-ide-collabsdk.OrganizationMembersQueryOptionIncludeSystem())
		require.NoError(t, err)
		require.Len(t, members, 4)
		require.ElementsMatch(t,
			[]uuid.UUID{owner.UserID, orgMember.ID, orgAdmin.ID, database.PrebuildsSystemUserID},
			slice.List(members, onlyIDs))
	})

	t.Run("GithubUserID", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)
		members, err := client.OrganizationMembers(ctx, anotherOrg.ID, optimus-ide-collabsdk.OrganizationMembersQueryOptionGithubUserID(anotherUser.GithubComUserID.Int64))
		require.NoError(t, err)
		require.Len(t, members, 1)
		require.ElementsMatch(t,
			[]uuid.UUID{anotherUser.ID},
			slice.List(members, onlyIDs))
	})
}

func TestGetOrgMembersFilter(t *testing.T) {
	t.Parallel()

	client, _, api := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		IncludeProvisionerDaemon: true,
		OIDCConfig: &optimus-ide-collabd.OIDCConfig{
			AllowSignups: true,
		},
	})
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)

	setupCtx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	optimus-ide-collabdtest.UsersFilter(setupCtx, t, client, api.Database, nil, nil, func(testCtx context.Context, req optimus-ide-collabsdk.UsersRequest) []optimus-ide-collabsdk.ReducedUser {
		res, err := client.OrganizationMembersPaginated(testCtx, first.OrganizationID, req)
		require.NoError(t, err)
		reduced := make([]optimus-ide-collabsdk.ReducedUser, len(res.Members))
		for i, user := range res.Members {
			reduced[i] = orgMemberToReducedUser(user)
		}
		return reduced
	})
}

func TestGetOrgMembersPagination(t *testing.T) {
	t.Parallel()
	client := optimus-ide-collabdtest.New(t, nil)
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	optimus-ide-collabdtest.UsersPagination(ctx, t, client, nil, func(req optimus-ide-collabsdk.UsersRequest) ([]optimus-ide-collabsdk.ReducedUser, int) {
		res, err := client.OrganizationMembersPaginated(ctx, first.OrganizationID, req)
		require.NoError(t, err)
		reduced := make([]optimus-ide-collabsdk.ReducedUser, len(res.Members))
		for i, user := range res.Members {
			reduced[i] = orgMemberToReducedUser(user)
		}
		return reduced, res.Count
	})
}

func onlyIDs(u optimus-ide-collabsdk.OrganizationMemberWithUserData) uuid.UUID {
	return u.UserID
}

func orgMemberToReducedUser(user optimus-ide-collabsdk.OrganizationMemberWithUserData) optimus-ide-collabsdk.ReducedUser {
	return optimus-ide-collabsdk.ReducedUser{
		MinimalUser: optimus-ide-collabsdk.MinimalUser{
			ID:        user.UserID,
			Username:  user.Username,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
		},
		Email:            user.Email,
		CreatedAt:        user.UserCreatedAt,
		UpdatedAt:        user.UserUpdatedAt,
		LastSeenAt:       user.LastSeenAt,
		Status:           user.Status,
		IsServiceAccount: user.IsServiceAccount,
		LoginType:        user.LoginType,
	}
}
