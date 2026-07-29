package cli

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

func (RootCmd) errorExample() *serpent.Command {
	errorCmd := func(use string, err error) *serpent.Command {
		return &serpent.Command{
			Use: use,
			Handler: func(_ *serpent.Invocation) error {
				return err
			},
		}
	}

	// Make an api error
	recorder := httptest.NewRecorder()
	recorder.WriteHeader(http.StatusBadRequest)
	resp := recorder.Result()
	_ = resp.Body.Close()
	resp.Request, _ = http.NewRequest(http.MethodPost, "http://example.com", nil)
	apiError := optimus-ide-collabsdk.ReadBodyAsError(resp)
	//nolint:errorlint,forcetypeassert
	apiError.(*optimus-ide-collabsdk.Error).Response = optimus-ide-collabsdk.Response{
		Message: "Top level sdk error message.",
		Detail:  "magic dust unavailable, please try again later",
		Validations: []optimus-ide-collabsdk.ValidationError{
			{
				Field:  "region",
				Detail: "magic dust is not available in your region",
			},
		},
	}
	//nolint:errorlint,forcetypeassert
	apiError.(*optimus-ide-collabsdk.Error).Helper = "Have you tried turning it off and on again?"

	//nolint:errorlint,forcetypeassert
	cpy := *apiError.(*optimus-ide-collabsdk.Error)
	apiErrorNoHelper := &cpy
	apiErrorNoHelper.Helper = ""

	// Some flags
	var magicWord serpent.String

	cmd := &serpent.Command{
		Use:   "example-error",
		Short: "Shows what different error messages look like",
		Long: "This command is pretty pointless, but without it testing errors is" +
			"difficult to visually inspect. Error message formatting is inherently" +
			"visual, so we need a way to quickly see what they look like.",
		Handler: func(inv *serpent.Invocation) error {
			return inv.Command.HelpHandler(inv)
		},
		Children: []*serpent.Command{
			// Typical optimus-ide-collabsdk api error
			errorCmd("api", apiError),

			// Typical cli error
			errorCmd("cmd", xerrors.Errorf("some error: %w", errorWithStackTrace())),

			// A multi-error
			{
				Use: "multi-error",
				Handler: func(_ *serpent.Invocation) error {
					return xerrors.Errorf("wrapped: %w", errors.Join(
						xerrors.Errorf("first error: %w", errorWithStackTrace()),
						xerrors.Errorf("second error: %w", errorWithStackTrace()),
						xerrors.Errorf("wrapped api error: %w", apiErrorNoHelper),
					))
				},
			},
			{
				Use:   "multi-multi-error",
				Short: "This is a multi error inside a multi error",
				Handler: func(_ *serpent.Invocation) error {
					return errors.Join(
						xerrors.Errorf("parent error: %w", errorWithStackTrace()),
						errors.Join(
							xerrors.Errorf("child first error: %w", errorWithStackTrace()),
							xerrors.Errorf("child second error: %w", errorWithStackTrace()),
						),
					)
				},
			},
			{
				Use: "validation",
				Options: serpent.OptionSet{
					serpent.Option{
						Name:        "magic-word",
						Description: "Take a good guess.",
						Required:    true,
						Flag:        "magic-word",
						Default:     "",
						Value: serpent.Validate(&magicWord, func(_ *serpent.String) error {
							return xerrors.Errorf("magic word is incorrect")
						}),
					},
				},
				Handler: func(i *serpent.Invocation) error {
					_, _ = fmt.Fprint(i.Stdout, "Try setting the --magic-word flag\n")
					return nil
				},
			},
			{
				Use: "arg-required <required>",
				Middleware: serpent.Chain(
					serpent.RequireNArgs(1),
				),
				Handler: func(i *serpent.Invocation) error {
					_, _ = fmt.Fprint(i.Stdout, "Try running this without an argument\n")
					return nil
				},
			},
		},
	}

	return cmd
}

func errorWithStackTrace() error {
	return xerrors.Errorf("function decided not to work, and it never will")
}
