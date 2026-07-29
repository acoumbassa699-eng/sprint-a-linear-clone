# Google authentication (OIDC)

This guide shows how to configure Optimus-IDE-Collab to authenticate users with Google using OpenID Connect (OIDC).

## Prerequisites

- A Google Cloud project with the OAuth consent screen configured
- Permission to create OAuth 2.0 Client IDs in Google Cloud

## Step 1: Create an OAuth client in Google Cloud

1. Open Google Cloud Console → APIs & Services → Credentials → Create Credentials → OAuth client ID.
2. Application type: Web application.
3. Authorized redirect URIs: add your Optimus-IDE-Collab callback URL:
   - `https://optimus-ide-collab.example.com/api/v2/users/oidc/callback`
4. Save and note the Client ID and Client secret.

## Step 2: Configure Optimus-IDE-Collab OIDC for Google

Set the following environment variables on your Optimus-IDE-Collab deployment and restart Optimus-IDE-Collab:

```dotenv
OPTIMUS-IDE-COLLAB_OIDC_ISSUER_URL=https://accounts.google.com
OPTIMUS-IDE-COLLAB_OIDC_CLIENT_ID=<client id>
OPTIMUS-IDE-COLLAB_OIDC_CLIENT_SECRET=<client secret>
# Restrict to one or more email domains (comma-separated)
OPTIMUS-IDE-COLLAB_OIDC_EMAIL_DOMAIN="example.com"
# Standard OIDC scopes for Google
OPTIMUS-IDE-COLLAB_OIDC_SCOPES=openid,profile,email
# Optional: customize the login button
OPTIMUS-IDE-COLLAB_OIDC_SIGN_IN_TEXT="Sign in with Google"
OPTIMUS-IDE-COLLAB_OIDC_ICON_URL=/icon/google.svg
```

> [!NOTE]
> The redirect URI must exactly match what you configured in Google Cloud.

## Enable refresh tokens (recommended)

Google uses auth URL parameters to issue refresh tokens. Configure:

```dotenv
# Keep standard scopes
OPTIMUS-IDE-COLLAB_OIDC_SCOPES=openid,profile,email
# Add Google-specific auth URL params
OPTIMUS-IDE-COLLAB_OIDC_AUTH_URL_PARAMS='{"access_type": "offline", "prompt": "consent"}'
```

After changing settings, users must log out and back in once to obtain refresh tokens.

Learn more in [Configure OIDC refresh tokens](./refresh-tokens.md).

## Troubleshooting

- "invalid redirect_uri": ensure the redirect URI in Google Cloud matches `https://<your-optimus-ide-collab-host>/api/v2/users/oidc/callback`.
- Domain restriction: if users from unexpected domains can log in, verify `OPTIMUS-IDE-COLLAB_OIDC_EMAIL_DOMAIN`.
- Claims: to inspect claims returned by Google, see guidance in the [OIDC overview](./index.md#oidc-claims).

## See also

- [OIDC overview](./index.md)
- [Configure OIDC refresh tokens](./refresh-tokens.md)
