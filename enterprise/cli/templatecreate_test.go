package cli_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestTemplateCreate(t *testing.T) {
	t.Parallel()

	t.Run("RequireActiveVersion", func(t *testing.T) {
		t.Parallel()

		client, user := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAccessControl: 1,
				},
			},
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
		})
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, user.OrganizationID, rbac.RoleTemplateAdmin())

		source := clitest.CreateTemplateVersionSource(t, &echo.Responses{
			Parse:          echo.ParseComplete,
			ProvisionApply: echo.ApplyComplete,
		})

		inv, conf := newCLI(t, "templates",
			"create", "new-template",
			"--directory", source,
			"--test.provisioner", string(database.ProvisionerTypeEcho),
			"--require-active-version",
			"-y",
		)

		clitest.SetupConfig(t, templateAdmin, conf)

		err := inv.Run()
		require.NoError(t, err)

		ctx := testutil.Context(t, testutil.WaitMedium)
		template, err := templateAdmin.TemplateByName(ctx, user.OrganizationID, "new-template")
		require.NoError(t, err)
		require.True(t, template.RequireActiveVersion)
	})

	t.Run("WorkspaceCleanup", func(t *testing.T) {
		t.Parallel()

		client, user := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAdvancedTemplateScheduling: 1,
				},
			},
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
		})
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, user.OrganizationID, rbac.RoleTemplateAdmin())

		source := clitest.CreateTemplateVersionSource(t, &echo.Responses{
			Parse:          echo.ParseComplete,
			ProvisionApply: echo.ApplyComplete,
		})

		const (
			expectedFailureTTL           = time.Hour * 3
			expectedDormancyThreshold    = time.Hour * 4
			expectedDormancyAutoDeletion = time.Minute * 10
		)

		inv, conf := newCLI(t, "templates",
			"create", "new-template",
			"--directory", source,
			"--test.provisioner", string(database.ProvisionerTypeEcho),
			"--failure-ttl="+expectedFailureTTL.String(),
			"--dormancy-threshold="+expectedDormancyThreshold.String(),
			"--dormancy-auto-deletion="+expectedDormancyAutoDeletion.String(),
			"-y",
			"--",
		)

		clitest.SetupConfig(t, templateAdmin, conf)

		err := inv.Run()
		require.NoError(t, err)

		ctx := testutil.Context(t, testutil.WaitMedium)
		template, err := templateAdmin.TemplateByName(ctx, user.OrganizationID, "new-template")
		require.NoError(t, err)
		require.Equal(t, expectedFailureTTL.Milliseconds(), template.FailureTTLMillis)
		require.Equal(t, expectedDormancyThreshold.Milliseconds(), template.TimeTilDormantMillis)
		require.Equal(t, expectedDormancyAutoDeletion.Milliseconds(), template.TimeTilDormantAutoDeleteMillis)
	})

	t.Run("NotEntitled", func(t *testing.T) {
		t.Parallel()

		client, admin := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{},
			},
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
		})
		templateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID, rbac.RoleTemplateAdmin())

		inv, conf := newCLI(t, "templates",
			"create", "new-template",
			"--require-active-version",
			"-y",
		)

		clitest.SetupConfig(t, templateAdmin, conf)

		err := inv.Run()
		require.Error(t, err)
		require.Contains(t, err.Error(), "your license is not entitled to use enterprise access control, so you cannot set --require-active-version")
	})

	// Create a template in a second organization via custom role
	t.Run("SecondOrganization", func(t *testing.T) {
		t.Parallel()

		ownerClient, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				// This only affects the first org.
				IncludeProvisionerDaemon: false,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAccessControl:              1,
					optimus-ide-collabsdk.FeatureCustomRoles:                1,
					optimus-ide-collabsdk.FeatureExternalProvisionerDaemons: 1,
					optimus-ide-collabsdk.FeatureMultipleOrganizations:      1,
				},
			},
		})

		// Create the second organization
		secondOrg := optimus-ide-collabdenttest.CreateOrganization(t, ownerClient, optimus-ide-collabdenttest.CreateOrganizationOptions{
			IncludeProvisionerDaemon: true,
		})

		ctx := testutil.Context(t, testutil.WaitMedium)

		//nolint:gocritic // owner required to make custom roles
		orgTemplateAdminRole, err := ownerClient.CreateOrganizationRole(ctx, optimus-ide-collabsdk.Role{
			Name:           "org-template-admin",
			OrganizationID: secondOrg.ID.String(),
			OrganizationPermissions: optimus-ide-collabsdk.CreatePermissions(map[optimus-ide-collabsdk.RBACResource][]optimus-ide-collabsdk.RBACAction{
				optimus-ide-collabsdk.ResourceTemplate: optimus-ide-collabsdk.RBACResourceActions[optimus-ide-collabsdk.ResourceTemplate],
			}),
		})
		require.NoError(t, err, "create admin role")

		orgTemplateAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, secondOrg.ID, rbac.RoleIdentifier{
			Name:           orgTemplateAdminRole.Name,
			OrganizationID: secondOrg.ID,
		})

		source := clitest.CreateTemplateVersionSource(t, &echo.Responses{
			Parse:          echo.ParseComplete,
			ProvisionApply: echo.ApplyComplete,
		})

		const templateName = "new-template"
		inv, conf := newCLI(t, "templates",
			"push", templateName,
			"--directory", source,
			"--test.provisioner", string(database.ProvisionerTypeEcho),
			"-y",
		)

		clitest.SetupConfig(t, orgTemplateAdmin, conf)

		err = inv.Run()
		require.NoError(t, err)

		ctx = testutil.Context(t, testutil.WaitMedium)
		template, err := orgTemplateAdmin.TemplateByName(ctx, secondOrg.ID, templateName)
		require.NoError(t, err)
		require.Equal(t, template.OrganizationID, secondOrg.ID)
	})
}
