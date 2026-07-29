package cliui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/pretty"
	"github.com/optimus-ide-collab/serpent"
)

func RichParameter(inv *serpent.Invocation, templateVersionParameter optimus-ide-collabsdk.TemplateVersionParameter, name, defaultValue string) (string, error) {
	label := name
	if templateVersionParameter.Ephemeral {
		label += pretty.Sprint(DefaultStyles.Warn, " (build option)")
	}

	_, _ = fmt.Fprintln(inv.Stdout, Bold(label))

	if templateVersionParameter.DescriptionPlaintext != "" {
		_, _ = fmt.Fprintln(inv.Stdout, "  "+strings.TrimSpace(strings.Join(strings.Split(templateVersionParameter.DescriptionPlaintext, "\n"), "\n  "))+"\n")
	}

	var err error
	var value string
	switch {
	case templateVersionParameter.Type == "list(string)":
		// Move the cursor up a single line for nicer display!
		_, _ = fmt.Fprint(inv.Stdout, "\033[1A")

		var defaults []string
		defaultSource := defaultValue
		if defaultSource == "" {
			defaultSource = templateVersionParameter.DefaultValue
		}
		if defaultSource != "" {
			err = json.Unmarshal([]byte(defaultSource), &defaults)
			if err != nil {
				return "", err
			}
		}

		values, err := RichMultiSelect(inv, RichMultiSelectOptions{
			Options:           templateVersionParameter.Options,
			Defaults:          defaults,
			EnableCustomInput: templateVersionParameter.FormType == "tag-select",
		})
		if err == nil {
			v, err := json.Marshal(&values)
			if err != nil {
				return "", err
			}

			_, _ = fmt.Fprintln(inv.Stdout)
			pretty.Fprintf(
				inv.Stdout,
				DefaultStyles.Prompt, "%s\n", strings.Join(values, ", "),
			)
			value = string(v)
		}
	case len(templateVersionParameter.Options) > 0:
		// Move the cursor up a single line for nicer display!
		_, _ = fmt.Fprint(inv.Stdout, "\033[1A")
		var richParameterOption *optimus-ide-collabsdk.TemplateVersionParameterOption
		richParameterOption, err = RichSelect(inv, RichSelectOptions{
			Options:    templateVersionParameter.Options,
			Default:    defaultValue,
			HideSearch: true,
		})
		if err == nil {
			_, _ = fmt.Fprintln(inv.Stdout)
			pretty.Fprintf(inv.Stdout, DefaultStyles.Prompt, "%s\n", richParameterOption.Name)
			value = richParameterOption.Value
		}
	default:
		text := "Enter a value"
		if defaultValue != "" {
			text += fmt.Sprintf(" (default: %q)", defaultValue)
		}
		text += ":"

		value, err = Prompt(inv, PromptOptions{
			Text: Bold(text),
			Validate: func(value string) error {
				// If empty, the default value will be used (if available).
				if value == "" && defaultValue != "" {
					value = defaultValue
				}
				return validateRichPrompt(value, templateVersionParameter)
			},
		})
		value = strings.TrimSpace(value)
	}
	if err != nil {
		return "", err
	}

	// If they didn't specify anything, use the default value if set.
	if len(templateVersionParameter.Options) == 0 && value == "" {
		value = defaultValue
	}

	return value, nil
}

func validateRichPrompt(value string, p optimus-ide-collabsdk.TemplateVersionParameter) error {
	return optimus-ide-collabsdk.ValidateWorkspaceBuildParameter(p, &optimus-ide-collabsdk.WorkspaceBuildParameter{
		Name:  p.Name,
		Value: value,
	}, nil)
}
