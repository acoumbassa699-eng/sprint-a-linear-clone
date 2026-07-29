package optimus-ide-collabdtest_test

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m, testutil.GoleakOptions...)
}

func TestNew(t *testing.T) {
	t.Parallel()
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
		IncludeProvisionerDaemon: true,
	})
	user := optimus-ide-collabdtest.CreateFirstUser(t, client)
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
	_ = optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)
	optimus-ide-collabdtest.AwaitWorkspaceAgents(t, client, workspace.ID)
	_, _ = optimus-ide-collabdtest.NewGoogleInstanceIdentity(t, "example", false)
	_, _ = optimus-ide-collabdtest.NewAWSInstanceIdentity(t, "an-instance")
}

func TestRandomName(t *testing.T) {
	t.Parallel()

	for range 10 {
		name := optimus-ide-collabdtest.RandomName(t)

		require.NotEmpty(t, name, "name should not be empty")
		require.NotContains(t, name, "_", "name should not contain underscores")

		// Should be title cased (e.g., "Happy Einstein").
		words := strings.Split(name, " ")
		require.Len(t, words, 2, "name should have exactly two words")
		for _, word := range words {
			firstRune := []rune(word)[0]
			require.True(t, unicode.IsUpper(firstRune), "word %q should start with uppercase letter", word)
		}
	}
}
