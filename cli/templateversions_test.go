package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
)

func TestTemplateVersions(t *testing.T) {
	t.Parallel()
	t.Run("ListVersions", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		inv, root := clitest.New(t, "templates", "versions", "list", template.Name)
		clitest.SetupConfig(t, member, root)

		stdout := expecter.NewAttachedToInvocation(t, inv)

		errC := make(chan error)
		go func() {
			errC <- inv.Run()
		}()

		require.NoError(t, <-errC)

		stdout.ExpectMatch(ctx, version.Name)
		stdout.ExpectMatch(ctx, version.CreatedBy.Username)
		stdout.ExpectMatch(ctx, "Active")
	})

	t.Run("ListVersionsJSON", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		inv, root := clitest.New(t, "templates", "versions", "list", template.Name, "--output", "json")
		clitest.SetupConfig(t, member, root)

		var stdout bytes.Buffer
		inv.Stdout = &stdout

		require.NoError(t, inv.Run())

		var rows []struct {
			TemplateVersion optimus-ide-collabsdk.TemplateVersion `json:"TemplateVersion"`
			Active          bool                     `json:"active"`
		}
		require.NoError(t, json.Unmarshal(stdout.Bytes(), &rows))
		require.Len(t, rows, 1)
		assert.Equal(t, version.ID, rows[0].TemplateVersion.ID)
		assert.True(t, rows[0].Active)
	})
}

func TestTemplateVersionsPromote(t *testing.T) {
	t.Parallel()

	t.Run("PromoteVersion", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)

		// Create a template with two versions
		version1 := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, completeWithAgent())
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version1.ID)

		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version1.ID)

		version2 := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, completeWithAgent(), func(ctvr *optimus-ide-collabsdk.CreateTemplateVersionRequest) {
			ctvr.TemplateID = template.ID
			ctvr.Name = "2.0.0"
		})
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version2.ID)

		// Ensure version1 is active
		updatedTemplate, err := client.Template(context.Background(), template.ID)
		assert.NoError(t, err)
		assert.Equal(t, version1.ID, updatedTemplate.ActiveVersionID)

		args := []string{
			"templates",
			"versions",
			"promote",
			"--template", template.Name,
			"--template-version", version2.Name,
		}

		inv, root := clitest.New(t, args...)
		//nolint:gocritic // Creating a workspace for another user requires owner permissions.
		clitest.SetupConfig(t, client, root)
		errC := make(chan error)
		go func() {
			errC <- inv.Run()
		}()

		require.NoError(t, <-errC)

		// Verify that version2 is now the active version
		updatedTemplate, err = client.Template(context.Background(), template.ID)
		require.NoError(t, err)
		assert.Equal(t, version2.ID, updatedTemplate.ActiveVersionID)
	})

	t.Run("PromoteNonExistentVersion", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		inv, root := clitest.New(t, "templates", "versions", "promote", "--template", template.Name, "--template-version", "non-existent-version")
		clitest.SetupConfig(t, member, root)

		err := inv.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), "get template version by name")
	})

	t.Run("PromoteVersionInvalidTemplate", func(t *testing.T) {
		t.Parallel()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		inv, root := clitest.New(t, "templates", "versions", "promote", "--template", "non-existent-template", "--template-version", "some-version")
		clitest.SetupConfig(t, member, root)

		err := inv.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), "get template by name")
	})
}
