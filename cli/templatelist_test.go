package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
)

func TestTemplateList(t *testing.T) {
	t.Parallel()
	t.Run("ListTemplates", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())
		firstVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, firstVersion.ID)
		firstTemplate := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, firstVersion.ID)

		secondVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, secondVersion.ID)
		secondTemplate := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, secondVersion.ID)

		inv, root := clitest.New(t, "templates", "list")
		clitest.SetupConfig(t, templateAdmin, root)

		stdout := expecter.NewAttachedToInvocation(t, inv)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()

		errC := make(chan error)
		go func() {
			errC <- inv.WithContext(ctx).Run()
		}()

		// expect that templates are listed alphabetically
		templatesList := []string{firstTemplate.Name, secondTemplate.Name}
		slices.Sort(templatesList)

		require.NoError(t, <-errC)

		for _, name := range templatesList {
			stdout.ExpectMatch(ctx, name)
		}
	})
	t.Run("ListTemplatesJSON", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())
		firstVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, firstVersion.ID)
		_ = optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, firstVersion.ID)

		secondVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, secondVersion.ID)
		_ = optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, secondVersion.ID)

		inv, root := clitest.New(t, "templates", "list", "--output=json")
		clitest.SetupConfig(t, templateAdmin, root)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()

		out := bytes.NewBuffer(nil)
		inv.Stdout = out
		err := inv.WithContext(ctx).Run()
		require.NoError(t, err)

		var templates []optimus-ide-collabsdk.Template
		require.NoError(t, json.Unmarshal(out.Bytes(), &templates))
		require.Len(t, templates, 2)
	})
	t.Run("NoTemplates", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)

		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())

		inv, root := clitest.New(t, "templates", "list")
		clitest.SetupConfig(t, templateAdmin, root)

		stdout := expecter.NewAttachedToInvocation(t, inv)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()

		errC := make(chan error)
		go func() {
			errC <- inv.WithContext(ctx).Run()
		}()

		require.NoError(t, <-errC)

		stdout.ExpectMatch(ctx, "No templates found")
		stdout.ExpectMatch(ctx, "Create one:")
	})
}
