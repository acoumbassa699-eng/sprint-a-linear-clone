package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
)

func TestUserList(t *testing.T) {
	t.Parallel()
	t.Run("Table", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())
		inv, root := clitest.New(t, "users", "list")
		clitest.SetupConfig(t, userAdmin, root)
		stdout := expecter.NewAttachedToInvocation(t, inv)
		errC := make(chan error)
		go func() {
			errC <- inv.Run()
		}()
		require.NoError(t, <-errC)
		stdout.ExpectMatch(ctx, "optimus-ide-collab.com")
	})
	t.Run("JSON", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())
		inv, root := clitest.New(t, "users", "list", "-o", "json")
		clitest.SetupConfig(t, userAdmin, root)
		doneChan := make(chan struct{})

		buf := bytes.NewBuffer(nil)
		inv.Stdout = buf
		go func() {
			defer close(doneChan)
			err := inv.Run()
			assert.NoError(t, err)
		}()

		<-doneChan

		var users []optimus-ide-collabsdk.User
		err := json.Unmarshal(buf.Bytes(), &users)
		require.NoError(t, err, "unmarshal JSON output")
		require.Len(t, users, 2)
		for _, u := range users {
			assert.NotEmpty(t, u.ID)
			assert.NotEmpty(t, u.Email)
			assert.NotEmpty(t, u.Username)
			assert.NotEmpty(t, u.Name)
			assert.NotEmpty(t, u.CreatedAt)
			assert.NotEmpty(t, u.Status)
		}
	})
	t.Run("NoURLFileErrorHasHelperText", func(t *testing.T) {
		t.Parallel()

		executable, err := os.Executable()
		require.NoError(t, err)

		inv, _ := clitest.New(t, "users", "list")
		err = inv.Run()
		require.Contains(t, err.Error(), fmt.Sprintf("Try logging in using '%s login <url>'.", executable))
	})
	t.Run("SessionAuthErrorHasHelperText", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, nil)
		inv, root := clitest.New(t, "users", "list")
		clitest.SetupConfig(t, client, root)

		err := inv.Run()

		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Contains(t, err.Error(), "Try logging in using 'optimus-ide-collab login'.")
	})
}

func TestUserShow(t *testing.T) {
	t.Parallel()

	t.Run("Table", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())
		_, otherUser := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		inv, root := clitest.New(t, "users", "show", otherUser.Username)
		clitest.SetupConfig(t, userAdmin, root)
		doneChan := make(chan struct{})
		stdout := expecter.NewAttachedToInvocation(t, inv)
		go func() {
			defer close(doneChan)
			err := inv.Run()
			assert.NoError(t, err)
		}()
		stdout.ExpectMatch(ctx, otherUser.Email)
		<-doneChan
	})

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		client := optimus-ide-collabdtest.New(t, nil)
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		userAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleUserAdmin())
		other, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		otherUser, err := other.User(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err, "fetch other user")
		inv, root := clitest.New(t, "users", "show", otherUser.Username, "-o", "json")
		clitest.SetupConfig(t, userAdmin, root)
		doneChan := make(chan struct{})

		buf := bytes.NewBuffer(nil)
		inv.Stdout = buf
		go func() {
			defer close(doneChan)
			err := inv.Run()
			assert.NoError(t, err)
		}()

		<-doneChan

		var newUser optimus-ide-collabsdk.User
		err = json.Unmarshal(buf.Bytes(), &newUser)
		require.NoError(t, err, "unmarshal JSON output")
		require.Equal(t, otherUser.ID, newUser.ID)
		require.Equal(t, otherUser.Username, newUser.Username)
		require.Equal(t, otherUser.Email, newUser.Email)
		require.Equal(t, otherUser.Name, newUser.Name)
	})
}
