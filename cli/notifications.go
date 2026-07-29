package cli

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) notifications() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "notifications",
		Short: "Manage Optimus-IDE-Collab notifications",
		Long: "Administrators can use these commands to change notification settings.\n" + FormatExamples(
			Example{
				Description: "Pause Optimus-IDE-Collab notifications. Administrators can temporarily stop notifiers from dispatching messages in case of the target outage (for example: unavailable SMTP server or Webhook not responding)",
				Command:     "optimus-ide-collab notifications pause",
			},
			Example{
				Description: "Resume Optimus-IDE-Collab notifications",
				Command:     "optimus-ide-collab notifications resume",
			},
			Example{
				Description: "Send a test notification. Administrators can use this to verify the notification target settings",
				Command:     "optimus-ide-collab notifications test",
			},
			Example{
				Description: "Send a custom notification to the requesting user. Sending notifications targeting other users or groups is currently not supported",
				Command:     "optimus-ide-collab notifications custom \"Custom Title\" \"Custom Message\"",
			},
		),
		Aliases: []string{"notification"},
		Handler: func(inv *serpent.Invocation) error {
			return inv.Command.HelpHandler(inv)
		},
		Children: []*serpent.Command{
			r.pauseNotifications(),
			r.resumeNotifications(),
			r.testNotifications(),
			r.customNotifications(),
		},
	}
	return cmd
}

func (r *RootCmd) pauseNotifications() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "pause",
		Short: "Pause notifications",
		Middleware: serpent.Chain(
			serpent.RequireNArgs(0),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			err = client.PutNotificationsSettings(inv.Context(), optimus-ide-collabsdk.NotificationsSettings{
				NotifierPaused: true,
			})
			if err != nil {
				return xerrors.Errorf("unable to pause notifications: %w", err)
			}

			_, _ = fmt.Fprintln(inv.Stderr, "Notifications are now paused.")
			return nil
		},
	}
	return cmd
}

func (r *RootCmd) resumeNotifications() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "resume",
		Short: "Resume notifications",
		Middleware: serpent.Chain(
			serpent.RequireNArgs(0),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			err = client.PutNotificationsSettings(inv.Context(), optimus-ide-collabsdk.NotificationsSettings{
				NotifierPaused: false,
			})
			if err != nil {
				return xerrors.Errorf("unable to resume notifications: %w", err)
			}

			_, _ = fmt.Fprintln(inv.Stderr, "Notifications are now resumed.")
			return nil
		},
	}
	return cmd
}

func (r *RootCmd) testNotifications() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "test",
		Short: "Send a test notification",
		Middleware: serpent.Chain(
			serpent.RequireNArgs(0),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}

			if err := client.PostTestNotification(inv.Context()); err != nil {
				return xerrors.Errorf("unable to post test notification: %w", err)
			}

			_, _ = fmt.Fprintln(inv.Stderr, "A test notification has been sent. If you don't receive the notification, check Optimus-IDE-Collab's logs for any errors.")
			return nil
		},
	}
	return cmd
}

func (r *RootCmd) customNotifications() *serpent.Command {
	cmd := &serpent.Command{
		Use:   "custom <title> <message>",
		Short: "Send a custom notification",
		Middleware: serpent.Chain(
			serpent.RequireNArgs(2),
		),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}
			err = client.PostCustomNotification(inv.Context(), optimus-ide-collabsdk.CustomNotificationRequest{
				Content: &optimus-ide-collabsdk.CustomNotificationContent{
					Title:   inv.Args[0],
					Message: inv.Args[1],
				},
			})
			if err != nil {
				return xerrors.Errorf("unable to post custom notification: %w", err)
			}

			_, _ = fmt.Fprintln(inv.Stderr, "A custom notification has been sent.")
			return nil
		},
	}
	return cmd
}
