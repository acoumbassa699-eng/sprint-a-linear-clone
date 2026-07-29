package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestIsValidTemplateParameterOption(t *testing.T) {
	t.Parallel()

	options := []optimus-ide-collabsdk.TemplateVersionParameterOption{
		{Name: "Vim", Value: "vim"},
		{Name: "Emacs", Value: "emacs"},
		{Name: "VS Code", Value: "vscode"},
	}

	t.Run("SingleSelectValid", func(t *testing.T) {
		t.Parallel()
		bp := optimus-ide-collabsdk.WorkspaceBuildParameter{Name: "editor", Value: "vim"}
		tvp := optimus-ide-collabsdk.TemplateVersionParameter{
			Name:    "editor",
			Type:    "string",
			Options: options,
		}
		assert.True(t, isValidTemplateParameterOption(bp, tvp))
	})

	t.Run("SingleSelectInvalid", func(t *testing.T) {
		t.Parallel()
		bp := optimus-ide-collabsdk.WorkspaceBuildParameter{Name: "editor", Value: "notepad"}
		tvp := optimus-ide-collabsdk.TemplateVersionParameter{
			Name:    "editor",
			Type:    "string",
			Options: options,
		}
		assert.False(t, isValidTemplateParameterOption(bp, tvp))
	})

	t.Run("MultiSelectAllValid", func(t *testing.T) {
		t.Parallel()
		bp := optimus-ide-collabsdk.WorkspaceBuildParameter{Name: "editors", Value: `["vim","emacs"]`}
		tvp := optimus-ide-collabsdk.TemplateVersionParameter{
			Name:    "editors",
			Type:    "list(string)",
			Options: options,
		}
		assert.True(t, isValidTemplateParameterOption(bp, tvp))
	})

	t.Run("MultiSelectOneInvalid", func(t *testing.T) {
		t.Parallel()
		bp := optimus-ide-collabsdk.WorkspaceBuildParameter{Name: "editors", Value: `["vim","notepad"]`}
		tvp := optimus-ide-collabsdk.TemplateVersionParameter{
			Name:    "editors",
			Type:    "list(string)",
			Options: options,
		}
		assert.False(t, isValidTemplateParameterOption(bp, tvp))
	})

	t.Run("MultiSelectEmptyArray", func(t *testing.T) {
		t.Parallel()
		bp := optimus-ide-collabsdk.WorkspaceBuildParameter{Name: "editors", Value: `[]`}
		tvp := optimus-ide-collabsdk.TemplateVersionParameter{
			Name:    "editors",
			Type:    "list(string)",
			Options: options,
		}
		assert.True(t, isValidTemplateParameterOption(bp, tvp))
	})

	t.Run("MultiSelectInvalidJSON", func(t *testing.T) {
		t.Parallel()
		bp := optimus-ide-collabsdk.WorkspaceBuildParameter{Name: "editors", Value: `not-json`}
		tvp := optimus-ide-collabsdk.TemplateVersionParameter{
			Name:    "editors",
			Type:    "list(string)",
			Options: options,
		}
		assert.False(t, isValidTemplateParameterOption(bp, tvp))
	})
}
