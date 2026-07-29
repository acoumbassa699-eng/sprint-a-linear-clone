package cli

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) oauth2Provider() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "oauth2-provider",
		Short: "Manage Optimus-IDE-Collab OAuth2 provider settings",
		Long: "Administrators can use these commands to change OAuth2 provider settings.\n" + FormatExamples(
			Example{
				Description: "Enable dynamic client registration (RFC 7591), allowing OAuth2/MCP clients to self-register without an admin creating an app first",
				Command:     "optimus-ide-collab oauth2-provider dcr enable",
			},
			Example{
				Description: "Disable dynamic client registration. Clients that already registered are unaffected; only new self-registration attempts are rejected",
				Command:     "optimus-ide-collab oauth2-provider dcr disable",
			},
		),
		Handler: func(inv *serpent.Invocation) error {
			return inv.Command.HelpHandler(inv)
		},
		Children: []*serpent.Command{
			r.oauth2ProviderDCR(),
		},
	}
	return cmd
}

func (r *RootCmd) oauth2ProviderDCR() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "dcr",
		Short: "Manage OAuth2 dynamic client registration (RFC 7591)",
		Handler: func(inv *serpent.Invocation) error {
			return inv.Command.HelpHandler(inv)
		},
		Children: []*serpent.Command{
			r.oauth2ProviderDCRToggle(dcrToggleEnable),
			r.oauth2ProviderDCRToggle(dcrToggleDisable),
		},
	}
	return cmd
}

// dcrToggleAction distinguishes the "enable" and "disable" subcommands of
// `optimus-ide-collab oauth2-provider dcr`, which are otherwise identical.
type dcrToggleAction int

const (
	dcrToggleDisable dcrToggleAction = iota
	dcrToggleEnable
)

func (r *RootCmd) oauth2ProviderDCRToggle(action dcrToggleAction) *serpent.Command {
	enabled := action == dcrToggleEnable
	use, short, verb := "disable", "Disable OAuth2 dynamic client registration", "disable"
	if enabled {
		use, short, verb = "enable", "Enable OAuth2 dynamic client registration", "enable"
	}

	cmd := &serpent.Command{
		Use:   use,
		Short: short,
		Middleware: serpent.Chain(
			serpent.RequireNArgs(0),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			_, err = client.PutOAuth2ProviderSettings(inv.Context(), optimus-ide-collabsdk.OAuth2ProviderSettings{
				DynamicClientRegistrationEnabled: ptr.Ref(enabled),
			})
			if err != nil {
				return xerrors.Errorf("unable to %s dynamic client registration: %w", verb, err)
			}

			state := "disabled"
			if enabled {
				state = "enabled"
			}
			_, _ = fmt.Fprintf(inv.Stderr, "Dynamic client registration is now %s.\n", state)
			return nil
		},
	}
	return cmd
}
