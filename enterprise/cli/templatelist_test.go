package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestEnterpriseListTemplates(t *testing.T) {
	t.Parallel()

	t.Run("MultiOrg", func(t *testing.T) {
		t.Parallel()

		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations:      1,
					optimus-ide-collabsdk.FeatureExternalProvisionerDaemons: 1,
				},
			},
		})

		// Template in the first organization
		firstVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, nil)
		_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, firstVersion.ID)
		_ = optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, firstVersion.ID)

		secondOrg := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{
			IncludeProvisionerDaemon: true,
		})
		secondVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, client, secondOrg.ID, nil)
		_ = optimus-ide-collabdtest.CreateTemplate(t, client, secondOrg.ID, secondVersion.ID)

		// Create a site wide template admin
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.RoleTemplateAdmin())

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
}
