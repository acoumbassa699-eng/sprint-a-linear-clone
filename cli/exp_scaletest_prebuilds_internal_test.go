//go:build !slim

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scaletest/prebuilds"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func Test_getScaletestPrebuildsTemplates(t *testing.T) {
	t.Parallel()

	client, _, _ := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		IncludeProvisionerDaemon: true,
	})
	user := optimus-ide-collabdtest.CreateFirstUser(t, client)

	makeTemplate := func(t *testing.T, name string) {
		t.Helper()
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID, func(r *optimus-ide-collabsdk.CreateTemplateRequest) {
			r.Name = name
		})
	}

	// The real runner uses a small integer suffix (e.g. "0", "1"), keeping the
	// total name within the 32-character limit enforced by NameValid.
	const (
		scaletestPrebuildName = prebuilds.TemplatePrefix + "0"
		prebuildNoScaletest   = "prebuild-other"
		scaletestNoPrebuild   = "scaletest-other"
		unrelatedTemplate     = "unrelated-template"
	)

	makeTemplate(t, scaletestPrebuildName)
	makeTemplate(t, prebuildNoScaletest)
	makeTemplate(t, scaletestNoPrebuild)
	makeTemplate(t, unrelatedTemplate)

	t.Run("NoFilter", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		got, err := getScaletestPrebuildsTemplates(ctx, client, "")
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, scaletestPrebuildName, got[0].Name)
	})

	t.Run("MatchingTemplate", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		got, err := getScaletestPrebuildsTemplates(ctx, client, scaletestPrebuildName)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, scaletestPrebuildName, got[0].Name)
	})

	t.Run("NonExistentScaletestTemplate", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		got, err := getScaletestPrebuildsTemplates(ctx, client, prebuilds.TemplatePrefix+"99")
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("NonScaletestTemplateReturnsError", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		for _, name := range []string{prebuildNoScaletest, scaletestNoPrebuild, unrelatedTemplate} {
			_, err := getScaletestPrebuildsTemplates(ctx, client, name)
			require.Error(t, err, "expected error for template %q", name)
		}
	})
}
