package cli

import (
	"strings"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/pretty"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) publickey() *serpent.Command {
	var reset bool
	cmd := &serpent.Command{
		Use:     "publickey",
		Aliases: []string{"pubkey"},
		Short:   "Output your Optimus-IDE-Collab public key used for Git operations",
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}
			if reset {
				// Confirm prompt if using --reset. We don't want to accidentally
				// reset our public key.
				_, err := cliui.Prompt(inv, cliui.PromptOptions{
					Text: "Confirm regenerate a new sshkey for your workspaces? This will require updating the key " +
						"on any services it is registered with. This action cannot be reverted.",
					IsConfirm: true,
				})
				if err != nil {
					return err
				}

				// Reset the public key, let the retrieve re-read it.
				_, err = client.RegenerateGitSSHKey(inv.Context(), optimus-ide-collabsdk.Me)
				if err != nil {
					return err
				}
			}

			key, err := client.GitSSHKey(inv.Context(), optimus-ide-collabsdk.Me)
			if err != nil {
				return xerrors.Errorf("create optimus-ide-collabsdk client: %w", err)
			}

			cliui.Info(inv.Stdout,
				"This is your public key for using "+pretty.Sprint(cliui.DefaultStyles.Field, "git")+" in "+
					"Optimus-IDE-Collab. All clones with SSH will be authenticated automatically 🪄.",
			)
			cliui.Info(inv.Stdout, pretty.Sprint(cliui.DefaultStyles.Code, strings.TrimSpace(key.PublicKey))+"\n")
			cliui.Info(inv.Stdout, "Add to GitHub and GitLab:")
			cliui.Info(inv.Stdout, "> https://github.com/settings/ssh/new")
			cliui.Info(inv.Stdout, "> https://gitlab.com/-/profile/keys")

			return nil
		},
	}

	cmd.Options = serpent.OptionSet{
		{
			Flag:        "reset",
			Description: "Regenerate your public key. This will require updating the key on any services it's registered with.",
			Value:       serpent.BoolOf(&reset),
		},
		cliui.SkipPromptOption(),
	}

	return cmd
}
