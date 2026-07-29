package optimus-ide-collabd_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestMultiOrgFetch(t *testing.T) {
	t.Parallel()
	client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
			},
		},
	})

	ctx := testutil.Context(t, testutil.WaitLong)

	makeOrgs := []string{"foo", "bar", "baz"}
	for _, name := range makeOrgs {
		_, err := client.CreateOrganization(ctx, optimus-ide-collabsdk.CreateOrganizationRequest{
			Name:        name,
			DisplayName: name,
		})
		require.NoError(t, err)
	}

	//nolint:gocritic // using the owner intentionally since only they can make orgs
	myOrgs, err := client.OrganizationsByUser(ctx, optimus-ide-collabsdk.Me)
	require.NoError(t, err)
	require.NotNil(t, myOrgs)
	require.Len(t, myOrgs, len(makeOrgs)+1)

	orgs, err := client.Organizations(ctx)
	require.NoError(t, err)
	require.NotNil(t, orgs)
	require.ElementsMatch(t, myOrgs, orgs)
}

func TestOrganizationsByUser(t *testing.T) {
	t.Parallel()

	t.Run("IsDefault", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		//nolint:gocritic // owner is required to make orgs
		orgs, err := client.OrganizationsByUser(ctx, optimus-ide-collabsdk.Me)
		require.NoError(t, err)
		require.NotNil(t, orgs)
		require.Len(t, orgs, 1)
		require.True(t, orgs[0].IsDefault, "first org is always default")

		// Make an extra org, and it should not be defaulted.
		notDefault, err := client.CreateOrganization(ctx, optimus-ide-collabsdk.CreateOrganizationRequest{
			Name:        "another",
			DisplayName: "Another",
		})
		require.NoError(t, err)
		require.False(t, notDefault.IsDefault, "only 1 default org allowed")
	})

	t.Run("NoMember", func(t *testing.T) {
		t.Parallel()
		client, first := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		other, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID)
		ctx := testutil.Context(t, testutil.WaitLong)

		//nolint:gocritic // owner is required to make orgs
		org, err := client.CreateOrganization(ctx, optimus-ide-collabsdk.CreateOrganizationRequest{
			Name:        "another",
			DisplayName: "Another",
		})
		require.NoError(t, err)

		_, err = other.OrganizationByUserAndName(ctx, optimus-ide-collabsdk.Me, org.Name)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})
}

func TestAddOrganizationMembers(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		ownerClient, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})

		_, user := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)

		ctx := testutil.Context(t, testutil.WaitMedium)
		//nolint:gocritic // must be an owner, only owners can create orgs
		otherOrg, err := ownerClient.CreateOrganization(ctx, optimus-ide-collabsdk.CreateOrganizationRequest{
			Name:        "Other",
			DisplayName: "",
			Description: "",
			Icon:        "",
		})
		require.NoError(t, err, "create another organization")

		inv, root := clitest.New(t, "organization", "members", "add", "-O", otherOrg.ID.String(), user.Username)
		//nolint:gocritic // must be an owner
		clitest.SetupConfig(t, ownerClient, root)

		buf := new(bytes.Buffer)
		inv.Stdout = buf
		err = inv.WithContext(ctx).Run()
		require.NoError(t, err)

		//nolint:gocritic // must be an owner
		members, err := ownerClient.OrganizationMembers(ctx, otherOrg.ID)
		require.NoError(t, err)

		require.Len(t, members, 2)
	})
}

func TestDeleteOrganizationsByUser(t *testing.T) {
	t.Parallel()
	t.Run("Default", func(t *testing.T) {
		t.Parallel()
		client, user := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		// nolint:gocritic // owner used below to delete
		o, err := client.Organization(ctx, user.OrganizationID)
		require.NoError(t, err)

		// nolint:gocritic // only owners can delete orgs
		err = client.DeleteOrganization(ctx, o.ID.String())
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})

	t.Run("DeleteById", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{})

		// nolint:gocritic // only owners can delete orgs
		err := client.DeleteOrganization(ctx, o.ID.String())
		require.NoError(t, err)
	})

	t.Run("DeleteByName", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{})

		// nolint:gocritic // only owners can delete orgs
		err := client.DeleteOrganization(ctx, o.Name)
		require.NoError(t, err)
	})
}

func TestPatchOrganizationsByUser(t *testing.T) {
	t.Parallel()
	t.Run("Conflict", func(t *testing.T) {
		t.Parallel()
		client, user := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		// nolint:gocritic // owner used below as only they can create orgs
		originalOrg, err := client.Organization(ctx, user.OrganizationID)
		require.NoError(t, err)

		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{})

		// nolint:gocritic // owner used above to make the org
		_, err = client.UpdateOrganization(ctx, o.ID.String(), optimus-ide-collabsdk.UpdateOrganizationRequest{
			Name: originalOrg.Name,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusConflict, apiErr.StatusCode())
	})

	t.Run("ReservedName", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		var err error
		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{})

		_, err = client.UpdateOrganization(ctx, o.ID.String(), optimus-ide-collabsdk.UpdateOrganizationRequest{
			Name: optimus-ide-collabsdk.DefaultOrganization,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})

	t.Run("InvalidName", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		var err error
		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{})

		_, err = client.UpdateOrganization(ctx, o.ID.String(), optimus-ide-collabsdk.UpdateOrganizationRequest{
			Name: "something unique but not url safe",
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})

	t.Run("UpdateById", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		var err error
		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{})

		o, err = client.UpdateOrganization(ctx, o.ID.String(), optimus-ide-collabsdk.UpdateOrganizationRequest{
			Name: "new-new-org",
		})
		require.NoError(t, err)
		require.Equal(t, "new-new-org", o.Name)
	})

	t.Run("UpdateByName", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		const displayName = "New Organization"
		var err error
		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{}, func(request *optimus-ide-collabsdk.CreateOrganizationRequest) {
			request.DisplayName = displayName
		})

		o, err = client.UpdateOrganization(ctx, o.Name, optimus-ide-collabsdk.UpdateOrganizationRequest{
			Name: "new-new-org",
		})
		require.NoError(t, err)
		require.Equal(t, "new-new-org", o.Name)
		require.Equal(t, displayName, o.DisplayName) // didn't change
	})

	t.Run("UpdateDisplayName", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		var err error
		const name = "new-org"
		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{}, func(request *optimus-ide-collabsdk.CreateOrganizationRequest) {
			request.Name = name
		})

		const displayName = "The Newest One"
		o, err = client.UpdateOrganization(ctx, o.Name, optimus-ide-collabsdk.UpdateOrganizationRequest{
			DisplayName: "The Newest One",
		})
		require.NoError(t, err)
		require.Equal(t, "new-org", o.Name) // didn't change
		require.Equal(t, displayName, o.DisplayName)
	})

	t.Run("UpdateDescription", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		const displayName = "New Organization"
		var err error
		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{}, func(request *optimus-ide-collabsdk.CreateOrganizationRequest) {
			request.DisplayName = displayName
			request.Name = "new-org"
		})

		const description = "wow, this organization description is so updated!"
		o, err = client.UpdateOrganization(ctx, o.Name, optimus-ide-collabsdk.UpdateOrganizationRequest{
			Description: ptr.Ref(description),
		})

		require.NoError(t, err)
		require.Equal(t, "new-org", o.Name)          // didn't change
		require.Equal(t, displayName, o.DisplayName) // didn't change
		require.Equal(t, description, o.Description)
	})

	t.Run("UpdateIcon", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		const displayName = "New Organization"
		var err error
		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{}, func(request *optimus-ide-collabsdk.CreateOrganizationRequest) {
			request.DisplayName = displayName
			request.Icon = "/emojis/random.png"
			request.Name = "new-org"
		})

		const icon = "/emojis/1f48f-1f3ff.png"
		o, err = client.UpdateOrganization(ctx, o.Name, optimus-ide-collabsdk.UpdateOrganizationRequest{
			Icon: ptr.Ref(icon),
		})

		require.NoError(t, err)
		require.Equal(t, "new-org", o.Name)          // didn't change
		require.Equal(t, displayName, o.DisplayName) // didn't change
		require.Equal(t, icon, o.Icon)
	})

	t.Run("RevokedLicense", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitMedium)

		const displayName = "New Organization"
		o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{}, func(request *optimus-ide-collabsdk.CreateOrganizationRequest) {
			request.DisplayName = displayName
			request.Icon = "/emojis/random.png"
			request.Name = "new-org"
		})

		// Remove the license to block premium functionality.
		licenses, err := client.Licenses(ctx)
		require.NoError(t, err, "get licenses")
		for _, license := range licenses {
			// Should be only 1...
			err := client.DeleteLicense(ctx, license.ID)
			require.NoError(t, err, "delete license")
		}

		// Verify functionality is lost.
		const icon = "/emojis/1f48f-1f3ff.png"
		o, err = client.UpdateOrganization(ctx, o.Name, optimus-ide-collabsdk.UpdateOrganizationRequest{
			Icon: ptr.Ref(icon),
		})
		require.ErrorContains(t, err, "Multiple Organizations is a Premium feature")
	})

	t.Run("DefaultOrgMemberRoles", func(t *testing.T) {
		t.Parallel()

		t.Run("EqualToDefaultAllowedWithoutExperiment", func(t *testing.T) {
			t.Parallel()
			client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
					},
				},
			})
			ctx := testutil.Context(t, testutil.WaitMedium)
			o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{})

			// Writing exactly the deployment default is a no-op and must be allowed.
			//nolint:gocritic // Only owners can update organization settings.
			updated, err := client.UpdateOrganization(ctx, o.ID.String(), optimus-ide-collabsdk.UpdateOrganizationRequest{
				DefaultOrgMemberRoles: ptr.Ref(rbac.DefaultOrgMemberRoles()),
			})
			require.NoError(t, err)
			require.Equal(t, rbac.DefaultOrgMemberRoles(), updated.DefaultOrgMemberRoles)
		})

		t.Run("DeviationRejectedWithoutExperiment", func(t *testing.T) {
			t.Parallel()
			client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
					},
				},
			})
			ctx := testutil.Context(t, testutil.WaitMedium)
			o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{})

			// Empty array represents a Gateway Accounts organization. Without
			// the experiment, this must be rejected.
			//nolint:gocritic // Only owners can update organization settings.
			_, err := client.UpdateOrganization(ctx, o.ID.String(), optimus-ide-collabsdk.UpdateOrganizationRequest{
				DefaultOrgMemberRoles: ptr.Ref([]string{}),
			})
			var apiErr *optimus-ide-collabsdk.Error
			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, http.StatusForbidden, apiErr.StatusCode())
			require.Contains(t, apiErr.Message, "Changing default organization roles is not enabled")
		})

		t.Run("DeviationAllowedWithExperiment", func(t *testing.T) {
			t.Parallel()
			dv := optimus-ide-collabdtest.DeploymentValues(t)
			dv.Experiments = []string{string(optimus-ide-collabsdk.ExperimentMinimumImplicitMember)}
			client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
				Options: &optimus-ide-collabdtest.Options{DeploymentValues: dv},
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
					},
				},
			})
			ctx := testutil.Context(t, testutil.WaitMedium)
			o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{})

			//nolint:gocritic // Only owners can update organization settings.
			updated, err := client.UpdateOrganization(ctx, o.ID.String(), optimus-ide-collabsdk.UpdateOrganizationRequest{
				DefaultOrgMemberRoles: ptr.Ref([]string{}),
			})
			require.NoError(t, err)
			require.Empty(t, updated.DefaultOrgMemberRoles)
		})

		t.Run("NonBuiltInRoleRejected", func(t *testing.T) {
			t.Parallel()
			dv := optimus-ide-collabdtest.DeploymentValues(t)
			dv.Experiments = []string{string(optimus-ide-collabsdk.ExperimentMinimumImplicitMember)}
			client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
				Options: &optimus-ide-collabdtest.Options{DeploymentValues: dv},
				LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
					Features: license.Features{
						optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
					},
				},
			})
			ctx := testutil.Context(t, testutil.WaitMedium)
			o := optimus-ide-collabdenttest.CreateOrganization(t, client, optimus-ide-collabdenttest.CreateOrganizationOptions{})

			// A name that does not resolve via rbac.RoleByName (no such
			// built-in role) must be rejected. This blocks both custom roles
			// and malformed names like "foo:bar" that would otherwise break
			// RoleNameFromString downstream.
			//nolint:gocritic // Only owners can update organization settings.
			_, err := client.UpdateOrganization(ctx, o.ID.String(), optimus-ide-collabsdk.UpdateOrganizationRequest{
				DefaultOrgMemberRoles: ptr.Ref([]string{"not-a-built-in-role"}),
			})
			var apiErr *optimus-ide-collabsdk.Error
			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
			require.Contains(t, apiErr.Message, "Invalid default_org_member_roles entry")
		})
	})
}

func TestPostOrganizationsByUser(t *testing.T) {
	t.Parallel()
	t.Run("Conflict", func(t *testing.T) {
		t.Parallel()
		client, user := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitLong)

		//nolint:gocritic // using owner for below
		org, err := client.Organization(ctx, user.OrganizationID)
		require.NoError(t, err)

		//nolint:gocritic // only owners can create orgs
		_, err = client.CreateOrganization(ctx, optimus-ide-collabsdk.CreateOrganizationRequest{
			Name:        org.Name,
			DisplayName: org.DisplayName,
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusConflict, apiErr.StatusCode())
	})

	t.Run("InvalidName", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitLong)

		//nolint:gocritic // only owners can create orgs
		_, err := client.CreateOrganization(ctx, optimus-ide-collabsdk.CreateOrganizationRequest{
			Name:        "A name which is definitely not url safe",
			DisplayName: "New",
		})
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadRequest, apiErr.StatusCode())
	})

	t.Run("Create", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitLong)

		//nolint:gocritic // only owners can create orgs
		o, err := client.CreateOrganization(ctx, optimus-ide-collabsdk.CreateOrganizationRequest{
			Name:        "new-org",
			DisplayName: "New organization",
			Description: "A new organization to love and cherish forever.",
			Icon:        "/emojis/1f48f-1f3ff.png",
		})
		require.NoError(t, err)
		require.Equal(t, "new-org", o.Name)
		require.Equal(t, "New organization", o.DisplayName)
		require.Equal(t, "A new organization to love and cherish forever.", o.Description)
		require.Equal(t, "/emojis/1f48f-1f3ff.png", o.Icon)
	})

	t.Run("CreateWithoutExplicitDisplayName", func(t *testing.T) {
		t.Parallel()
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitLong)

		//nolint:gocritic // only owners can create orgs
		o, err := client.CreateOrganization(ctx, optimus-ide-collabsdk.CreateOrganizationRequest{
			Name: "new-org",
		})
		require.NoError(t, err)
		require.Equal(t, "new-org", o.Name)
		require.Equal(t, "new-org", o.DisplayName) // should match the given `Name`
	})
}
