package optimus-ide-collabd_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestOrganizationByUserAndName(t *testing.T) {
	t.Parallel()
	t.Run("NoExist", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		optimus-ide-collabdtest.CreateFirstUser(t, client)
		ctx := testutil.Context(t, testutil.WaitLong)

		_, err := client.OrganizationByUserAndName(ctx, optimus-ide-collabsdk.Me, "nothing")
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		ctx := testutil.Context(t, testutil.WaitLong)

		org, err := client.Organization(ctx, user.OrganizationID)
		require.NoError(t, err)
		_, err = client.OrganizationByUserAndName(ctx, optimus-ide-collabsdk.Me, org.Name)
		require.NoError(t, err)
	})
}
