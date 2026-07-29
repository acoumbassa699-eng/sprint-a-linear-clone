package cli_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
	"github.com/optimus-ide-collab/pretty"
)

func TestTemplateDelete(t *testing.T) {
	t.Parallel()

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		logger := testutil.Logger(t)
		ctx := testutil.Context(t, testutil.WaitMedium)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		inv, root := clitest.New(t, "templates", "delete", template.Name)

		clitest.SetupConfig(t, templateAdmin, root)
		stdout := expecter.NewAttachedToInvocation(t, inv)
		stdin := testutil.NewWriterAttachedToInvocation(t, logger.Named("stdin"), inv)

		execDone := make(chan error)
		go func() {
			execDone <- inv.Run()
		}()

		stdout.ExpectMatch(ctx, fmt.Sprintf("Delete these templates: %s?", pretty.Sprint(cliui.DefaultStyles.Code, template.Name)))
		stdin.WriteLine("yes")

		require.NoError(t, <-execDone)

		_, err := client.Template(context.Background(), template.ID)
		require.Error(t, err, "template should not exist")
	})

	t.Run("Multiple --yes", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())
		templates := []optimus-ide-collabsdk.Template{}
		templateNames := []string{}
		for i := 0; i < 3; i++ {
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
			_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)
			templates = append(templates, template)
			templateNames = append(templateNames, template.Name)
		}

		inv, root := clitest.New(t, append([]string{"templates", "delete", "--yes"}, templateNames...)...)
		clitest.SetupConfig(t, templateAdmin, root)
		require.NoError(t, inv.Run())

		for _, template := range templates {
			_, err := client.Template(context.Background(), template.ID)
			require.Error(t, err, "template should not exist")
		}
	})

	t.Run("Multiple prompted", func(t *testing.T) {
		t.Parallel()

		logger := testutil.Logger(t)
		ctx := testutil.Context(t, testutil.WaitMedium)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())
		templates := []optimus-ide-collabsdk.Template{}
		templateNames := []string{}
		for i := 0; i < 3; i++ {
			version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
			_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
			template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)
			templates = append(templates, template)
			templateNames = append(templateNames, template.Name)
		}

		inv, root := clitest.New(t, append([]string{"templates", "delete"}, templateNames...)...)
		clitest.SetupConfig(t, templateAdmin, root)
		stdout := expecter.NewAttachedToInvocation(t, inv)
		stdin := testutil.NewWriterAttachedToInvocation(t, logger.Named("stdin"), inv)

		execDone := make(chan error)
		go func() {
			execDone <- inv.Run()
		}()

		stdout.ExpectMatch(ctx,
			fmt.Sprintf("Delete these templates: %s?",
				pretty.Sprint(cliui.DefaultStyles.Code, strings.Join(templateNames, ", "))))
		stdin.WriteLine("yes")

		require.NoError(t, <-execDone)

		for _, template := range templates {
			_, err := client.Template(context.Background(), template.ID)
			require.Error(t, err, "template should not exist")
		}
	})

	t.Run("Selector", func(t *testing.T) {
		t.Parallel()

		logger := testutil.Logger(t)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		inv, root := clitest.New(t, "templates", "delete")
		clitest.SetupConfig(t, templateAdmin, root)

		stdin := testutil.NewWriterAttachedToInvocation(t, logger.Named("stdin"), inv)

		execDone := make(chan error)
		go func() {
			execDone <- inv.Run()
		}()

		stdin.WriteLine("yes")
		require.NoError(t, <-execDone)

		_, err := client.Template(context.Background(), template.ID)
		require.Error(t, err, "template should not exist")
	})
}
