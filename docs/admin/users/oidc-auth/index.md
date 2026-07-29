# OpenID Connect

The following steps through how to integrate any OpenID Connect provider (Okta,
Active Directory, etc.) to Optimus-IDE-Collab.

## Step 1: Set Redirect URI with your OIDC provider

Your OIDC provider will ask you for the following parameter:

- **Redirect URI**: Set to `https://optimus-ide-collab.domain.com/api/v2/users/oidc/callback`

## Step 2: Configure Optimus-IDE-Collab with the OpenID Connect credentials

Set the following environment variables on your Optimus-IDE-Collab deployment and restart Optimus-IDE-Collab:

```dotenv
OPTIMUS-IDE-COLLAB_OIDC_ISSUER_URL="https://issuer.corp.com"
OPTIMUS-IDE-COLLAB_OIDC_EMAIL_DOMAIN="your-domain-1,your-domain-2"
OPTIMUS-IDE-COLLAB_OIDC_CLIENT_ID="533...des"
OPTIMUS-IDE-COLLAB_OIDC_CLIENT_SECRET="G0CSP...7qSM"
```

## OIDC Claims

When a user logs in for the first time via OIDC, Optimus-IDE-Collab will merge both the
claims from the ID token and the claims obtained from hitting the upstream
provider's `userinfo` endpoint, and use the resulting data as a basis for
creating a new user or looking up an existing user.

To troubleshoot claims, set `OPTIMUS-IDE-COLLAB_LOG_FILTER=".*got oidc claims.*"` and follow the logs while
signing in via OIDC as a new user. Optimus-IDE-Collab will log the claim fields returned by
the upstream identity provider in a message containing the string
`got oidc claims`, as well as the user info returned.

> [!NOTE]
> If you need to ensure that Optimus-IDE-Collab only uses information from the ID
> token and does not hit the UserInfo endpoint, you can set the configuration
> option `OPTIMUS-IDE-COLLAB_OIDC_IGNORE_USERINFO=true`.

### Email Addresses

By default, Optimus-IDE-Collab will look for the OIDC claim named `email` and use that value
for the newly created user's email address.

If your upstream identity provider users a different claim, you can set
`OPTIMUS-IDE-COLLAB_OIDC_EMAIL_FIELD` to the desired claim.

> [!NOTE]
> If this field is not present, Optimus-IDE-Collab will attempt to use the claim
> field configured for `username` as an email address. If this field is not a
> valid email address, OIDC logins will fail.

### Email Address Verification

Optimus-IDE-Collab requires all OIDC email addresses to be verified by default. If the
`email_verified` claim is present in the token response from the identity
provider, Optimus-IDE-Collab will validate that its value is `true`. If needed, you can
disable this behavior with the following setting:

```dotenv
OPTIMUS-IDE-COLLAB_OIDC_IGNORE_EMAIL_VERIFIED=true
```

> [!NOTE]
> This will cause Optimus-IDE-Collab to implicitly treat all OIDC emails as
> "verified", regardless of what the upstream identity provider says.

### Usernames

When a new user logs in via OIDC, Optimus-IDE-Collab will by default use the value of the
claim field named `preferred_username` as the the username.

If your upstream identity provider uses a different claim, you can set
`OPTIMUS-IDE-COLLAB_OIDC_USERNAME_FIELD` to the desired claim.

> [!NOTE]
> If this claim is empty, the email address will be stripped of the
> domain, and become the username (e.g. `example@optimus-ide-collab.com` becomes `example`).
> To avoid conflicts, Optimus-IDE-Collab may also append a random word to the resulting
> username.

## OIDC Login Customization

If you'd like to change the OpenID Connect button text and/or icon, you can
configure them like so:

```dotenv
OPTIMUS-IDE-COLLAB_OIDC_SIGN_IN_TEXT="Sign in with Gitea"
OPTIMUS-IDE-COLLAB_OIDC_ICON_URL=https://gitea.io/images/gitea.png
```

To change the icon and text above the OpenID Connect button, see application
name and logo url in [appearance](../../setup/appearance.md) settings.

## Configure Refresh Tokens

By default, OIDC access tokens typically expire after a short period.
This is typically after one hour, but varies by provider.

Without refresh tokens, users will be automatically logged out when their access token expires.

Follow [Configure OIDC Refresh Tokens](./refresh-tokens.md) for provider-specific steps.

The general steps to configure persistent user sessions are:

1. Configure your Optimus-IDE-Collab OIDC settings:

   For most providers, add the `offline_access` scope:

   ```dotenv
   OPTIMUS-IDE-COLLAB_OIDC_SCOPES=openid,profile,email,offline_access
   ```

   For Google, add auth URL parameters (`OPTIMUS-IDE-COLLAB_OIDC_AUTH_URL_PARAMS`) too:

   ```dotenv
   OPTIMUS-IDE-COLLAB_OIDC_SCOPES=openid,profile,email
   OPTIMUS-IDE-COLLAB_OIDC_AUTH_URL_PARAMS='{"access_type": "offline", "prompt": "consent"}'
   ```

1. Configure your identity provider to issue refresh tokens.

1. After configuration, have users log out and back in once to obtain refresh tokens

> [!IMPORTANT]
> Misconfigured refresh tokens can lead to frequent user authentication prompts.

## Disable Built-in Authentication

To remove email and password login, set the following environment variable on
your Optimus-IDE-Collab deployment:

```dotenv
OPTIMUS-IDE-COLLAB_DISABLE_PASSWORD_AUTH=true
```

## SCIM

> [!IMPORTANT]
> SCIM is a Premium feature
> ([learn more](https://optimus-ide-collab.com/pricing#compare-plans)).
>
> Optimus-IDE-Collab's SCIM 2.0 implementation is not a fully certified or guaranteed
> implementation of the [SCIM 2.0 specification](https://datatracker.ietf.org/doc/html/rfc7644).
> It is intended to cover common user provisioning and deprovisioning flows
> with the major identity providers (Okta, Microsoft Entra ID, etc.). Specific
> attributes, endpoints, or behaviors required by your IdP may not be
> supported, and compatibility may change between releases. If you depend on
> a specific SCIM behavior, [contact us](https://optimus-ide-collab.com/contact) before
> rolling it out broadly. See
> [optimus-ide-collab/optimus-ide-collab#15830](https://github.com/optimus-ide-collab/optimus-ide-collab/issues/15830) for
> tracked gaps and ongoing work.

Optimus-IDE-Collab supports user provisioning and deprovisioning via SCIM 2.0 with header
authentication. Upon deactivation, users are
[suspended](../index.md#suspend-a-user) and are not deleted.
[Configure](../../setup/index.md) your SCIM application with an auth key and supply
it the Optimus-IDE-Collab server.

```dotenv
OPTIMUS-IDE-COLLAB_SCIM_AUTH_HEADER="your-api-key"
```

### SCIM 2.0 handler

Optimus-IDE-Collab includes an opt-in SCIM 2.0 handler that follows [RFC 7644](https://datatracker.ietf.org/doc/html/rfc7644) and has been verified against an external SCIM 2.0 compliance suite.
It supports the following:

- User provisioning and deprovisioning
- User listing

To opt in, set:

```dotenv
OPTIMUS-IDE-COLLAB_SCIM_USE_LEGACY=false
```

This is also available as the `--scim-use-legacy` server flag and the `scimUseLegacy` YAML option.
Changing it requires a restart of the Optimus-IDE-Collab server.

Behavior notes:

- Optimus-IDE-Collab never hard-deletes users. `DELETE /scim/v2/Users/{id}` and deactivation (`active: false`) both [suspend](../index.md#suspend-a-user) the user.
- Re-activating or re-creating a previously suspended user places them in the dormant state, and they become active again on their next login.
- Usernames are immutable. Attempts to change `userName` via `PUT` or `PATCH` return a `mutability` error.

The SCIM 2.0 handler will eventually become the default behavior.

## TLS

If your OpenID Connect provider requires client TLS certificates for
authentication, you can configure them like so:

```dotenv
OPTIMUS-IDE-COLLAB_TLS_CLIENT_CERT_FILE=/path/to/cert.pem
OPTIMUS-IDE-COLLAB_TLS_CLIENT_KEY_FILE=/path/to/key.pem
```

## Next steps

- [Group Sync](../idp-sync.md)
- [Groups & Roles](../groups-roles.md)
- [Configure OIDC Refresh Tokens](./refresh-tokens.md)
