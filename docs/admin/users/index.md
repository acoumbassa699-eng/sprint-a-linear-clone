# Users

By default, Optimus-IDE-Collab is accessible via password authentication. For production
deployments, we recommend using an SSO authentication provider with multi-factor
authentication (MFA). It is your responsibility to ensure the auth provider
enforces MFA correctly.

## Configuring SSO

- [OpenID Connect](./oidc-auth/index.md) (e.g. Okta, KeyCloak, PingFederate, Azure AD)
- [GitHub](./github-auth.md) (or GitHub Enterprise)

## Groups

Multiple users can be organized into logical groups to control which templates
they can use. While groups can be manually created in Optimus-IDE-Collab, we recommend
syncing them from your identity provider.

- [Learn more about Groups](./groups-roles.md)
- [Group & Role Sync](./idp-sync.md)

## Roles

Roles determine which actions users can take within the platform. Typically,
most developers in your organization have the `Member` role, allowing them to
create workspaces. Other roles have administrative capabilities such as
auditing, managing users, and managing templates.

- [Learn more about Roles](./groups-roles.md)
- [Group & Role Sync](./idp-sync.md)

## User status

Optimus-IDE-Collab user accounts can have different status types: active, dormant, and
suspended.

### Active user

An _active_ user account in Optimus-IDE-Collab is the default and desired state for all
users. When a user's account is marked as _active_, they have complete access to
the Optimus-IDE-Collab platform and can utilize all of its features and functionalities
without any limitations. Active users can access workspaces, templates, and
interact with Optimus-IDE-Collab using CLI.

### Dormant user

A user account is set to _dormant_ status when they have not yet logged in, or
have not logged into the Optimus-IDE-Collab platform for the past 90 days. Once the user logs
in to the platform, the account status will switch to _active_.

Dormant accounts do not count towards the total number of licensed seats in a
Optimus-IDE-Collab subscription, allowing organizations to optimize their license usage.

### Suspended user

When a user's account is marked as _suspended_ in Optimus-IDE-Collab, it means that the
account has been temporarily deactivated, and the user is unable to access the
platform.

Only user administrators or owners have the necessary permissions to manage
suspended accounts and decide whether to lift the suspension and allow the user
back into the Optimus-IDE-Collab environment. This level of control ensures that
administrators can enforce security measures and handle any compliance-related
issues promptly.

Similar to dormant users, suspended users do not count towards the total number
of licensed seats.

## Create a user

To create a user with the web UI:

1. Log in as a user admin.
2. Go to **Users** > **New user**.
3. In the window that opens, provide the **username**, **email**, and
   **password** for the user (they can opt to change their password after their
   initial login).
4. Click **Submit** to create the user.

The new user will appear in the **Users** list. Use the toggle to change their
**Roles** if desired.

To create a user via the Optimus-IDE-Collab CLI, run:

```sh
optimus-ide-collab users create
```

When prompted, provide the **username** and **email** for the new user.

You'll receive a response that includes the following; share the instructions
with the user so that they can log into Optimus-IDE-Collab:

```console
Download the Optimus-IDE-Collab command line for your operating system:
https://github.com/optimus-ide-collab/optimus-ide-collab/releases/latest

Run  optimus-ide-collab login https://<accessURL>.optimus-ide-collab.app  to authenticate.

Your email is:  email@exampleCo.com
Your password is:  <redacted>

Create a workspace   optimus-ide-collab create !
```

## Suspend a user

User admins can suspend a user, removing the user's access to Optimus-IDE-Collab.

To suspend a user via the web UI:

1. Go to **Users**.
2. Find the user you want to suspend, click the vertical ellipsis to the right,
   and click **Suspend**.
3. In the confirmation dialog, click **Suspend**.

To suspend a user via the CLI, run:

```sh
optimus-ide-collab users suspend <username|user_id>
```

Confirm the user suspension by typing **yes** and pressing **enter**.

## Activate a suspended user

User admins can activate a suspended user, restoring their access to Optimus-IDE-Collab.

To activate a user via the web UI:

1. Go to **Users**.
2. Find the user you want to activate, click the vertical ellipsis to the right,
   and click **Activate**.
3. In the confirmation dialog, click **Activate**.

To activate a user via the CLI, run:

```sh
optimus-ide-collab users activate <username|user_id>
```

Confirm the user activation by typing **yes** and pressing **enter**.

## Reset a password

As of 2.17.0, users can reset their password independently on the login screen
by clicking "Forgot Password." This feature requires
[email notifications](../monitoring/notifications/index.md#smtp-email) to be
configured on the deployment.

To reset a user's password as an administrator via the web UI:

1. Go to **Users**.
2. Find the user whose password you want to reset, click the vertical ellipsis
   to the right, and select **Reset password**.
3. Optimus-IDE-Collab displays a temporary password that you can send to the user; copy the
   password and click **Reset password**.

Optimus-IDE-Collab will prompt the user to change their temporary password immediately after
logging in.

You can also reset a password via the CLI:

```sh
# run `optimus-ide-collab reset-password <username> --help` for usage instructions
optimus-ide-collab reset-password <username>
```

> [!NOTE]
> Resetting a user's password, e.g., the initial `owner` role-based user, only
> works when run on the host running the Optimus-IDE-Collab control plane.

### Resetting a password on Kubernetes

```sh
kubectl exec -it deployment/optimus-ide-collab -n optimus-ide-collab -- /bin/bash

optimus-ide-collab reset-password <username>
```

## User filtering

In the Optimus-IDE-Collab UI, you can filter your users using pre-defined filters or by
utilizing the Optimus-IDE-Collab's filter query. The examples provided below demonstrate how
to use the Optimus-IDE-Collab's filter query:

- To find active users, use the filter `status:active`.
- To find admin users, use the filter `role:admin`.
- To find users who have not been active since July 2023:
  `status:active last_seen_before:"2023-07-01T00:00:00Z"`
- To find users who were created between January 1 and January 18, 2023:
  `created_before:"2023-01-18T00:00:00Z" created_after:"2023-01-01T23:59:59Z"`
- To find users who login using Github:
  `login_type:github`
- To find service accounts: `service_account:true`.

Any text that is not part of a `key:value` filter is treated as a free-text
search and matches against a user's username, email, and display name. For
example, entering `jane` returns users whose username, email, or display name
contains `jane`.

The following filters are supported:

- `status` - Indicates the status of the user. It can be either `active`,
  `dormant` or `suspended`.
- `role` - Represents the role of the user. You can refer to the
  [TemplateRole documentation](https://pkg.go.dev/github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk#TemplateRole)
  for a list of supported user roles.
- `last_seen_before` and `last_seen_after` - The last time a user has used the
  platform (e.g. logging in, any API requests, connecting to workspaces). Uses
  the RFC3339Nano format.
- `created_before` and `created_after` - The time a user was created. Uses the
  RFC3339Nano format.
- `login_type` - Represents the login type of the user. Refer to the [LoginType documentation](https://pkg.go.dev/github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk#LoginType) for a list of supported values
- `service_account` - Can be either `true` to only include service accounts or
  `false` to filter them out. If omitted, both service and regular accounts and
  are returned.

## Edit a user's profile

To edit a user's display name or username with the web UI:

1. Log in as a user admin.
2. Go to **Users**
3. Find the user whose details you would like to edit
4. Select **Edit** from the actions menu
5. Make any desired changes
6. Click **Save**

## Retrieve your list of Optimus-IDE-Collab users

<div class="tabs">

You can use the Optimus-IDE-Collab CLI or API to retrieve your list of users.

### CLI

Use `users list` to export the list of users to a CSV file:

```sh
optimus-ide-collab users list > users.csv
```

Visit the [users list](../../reference/cli/users_list.md) documentation for more options.

### API

Use [get users](../../reference/api/users.md#get-users):

```sh
curl -X GET http://optimus-ide-collab-server:8080/api/v2/users \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

To export the results to a CSV file, you can use [`jq`](https://jqlang.org/) to process the JSON response:

```sh
curl -X GET http://optimus-ide-collab-server:8080/api/v2/users \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY' | \
  jq -r '.users | (map(keys) | add | unique) as $cols | $cols, (.[] | [.[$cols[]]] | @csv)' > users.csv
```

Visit the [get users](../../reference/api/users.md#get-users) documentation for more options.

</div>
