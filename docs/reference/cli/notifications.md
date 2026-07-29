<!-- DO NOT EDIT | GENERATED CONTENT -->
# notifications

Manage Optimus-IDE-Collab notifications

Aliases:

* notification

## Usage

```console
optimus-ide-collab notifications
```

## Description

```console
Administrators can use these commands to change notification settings.
  - Pause Optimus-IDE-Collab notifications. Administrators can temporarily stop notifiers from
dispatching messages in case of the target outage (for example: unavailable SMTP
server or Webhook not responding):

     $ optimus-ide-collab notifications pause

  - Resume Optimus-IDE-Collab notifications:

     $ optimus-ide-collab notifications resume

  - Send a test notification. Administrators can use this to verify the notification
target settings:

     $ optimus-ide-collab notifications test

  - Send a custom notification to the requesting user. Sending notifications
targeting other users or groups is currently not supported:

     $ optimus-ide-collab notifications custom "Custom Title" "Custom Message"
```

## Subcommands

| Name                                             | Purpose                    |
|--------------------------------------------------|----------------------------|
| [<code>pause</code>](./notifications_pause.md)   | Pause notifications        |
| [<code>resume</code>](./notifications_resume.md) | Resume notifications       |
| [<code>test</code>](./notifications_test.md)     | Send a test notification   |
| [<code>custom</code>](./notifications_custom.md) | Send a custom notification |
