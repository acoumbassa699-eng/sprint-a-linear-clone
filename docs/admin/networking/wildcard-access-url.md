# Wildcard Access URLs

Wildcard access URLs unlock Optimus-IDE-Collab's full potential for modern development workflows. While optional for basic SSH usage, this feature becomes essential when teams need web applications, development previews, or browser-based tools. **Wildcard access URLs are essential for many development workflows in Optimus-IDE-Collab** - Web IDEs (code-server, VS Code Web, JupyterLab) and some development frameworks work significantly better with subdomain-based access rather than path-based URLs.

## Why configure wildcard access URLs?

### Key benefits

- **Enables port access**: Each application gets a unique subdomain with [port support](https://optimus-ide-collab.com/docs/user-guides/workspace-access/port-forwarding#dashboard) (e.g. `8080--main--myworkspace--john.optimus-ide-collab.example.com`).
- **Enhanced security**: Applications run in isolated subdomains with separate browser security contexts and prevents access to the Optimus-IDE-Collab API from malicious JavaScript
- **Better compatibility**: Most applications are designed to work at the root of a hostname rather than at a subpath, making subdomain access more reliable

### Applications that require subdomain access

The following tools require wildcard access URL:

- **Vite dev server**: Hot module replacement and asset serving issues with path-based routing
- **React dev server**: Similar issues with hot reloading and absolute path references
- **Next.js development server**: Asset serving and routing conflicts with path-based access
- **JupyterLab**: More complex template configuration and security risks when using path-based routing
- **RStudio**: More complex template configuration and security risks when using path-based routing

## Configuration

`OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL` is necessary for [port forwarding](port-forwarding.md#dashboard) via the dashboard or running [optimus-ide-collab_apps](../templates/index.md) on an absolute path.
Set it to a wildcard hostname that resolves to Optimus-IDE-Collab.
The value must contain exactly one `*` at the beginning of the hostname.
Optimus-IDE-Collab replaces `*` with the generated application name, which stays within a single DNS label.

Optimus-IDE-Collab supports the wildcard as a full label or with a suffix in the first label:

| Pattern              | Example generated application hostname           | Required DNS and TLS wildcard |
|----------------------|--------------------------------------------------|-------------------------------|
| `*.apps.example.com` | `8080--main--myworkspace--john.apps.example.com` | `*.apps.example.com`          |
| `*-apps.example.com` | `8080--main--myworkspace--john-apps.example.com` | `*.example.com`               |

For example, use the suffix pattern to keep the Optimus-IDE-Collab dashboard and application hostnames at the same DNS level:

```dotenv
OPTIMUS-IDE-COLLAB_ACCESS_URL=https://apps.example.com
OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL=*-apps.example.com
```

This configuration serves the dashboard from `https://apps.example.com` and a workspace application from a hostname such as `https://8080--main--myworkspace--john-apps.example.com`.

### TLS Certificate Setup

Wildcard access URLs require a TLS certificate that covers the wildcard domain. You have several options:

> [!TIP]
> You can use a single certificate for both the access URL and wildcard access URL.
> For `*.apps.example.com`, the certificate must include `apps.example.com` and `*.apps.example.com`.
> For `*-apps.example.com` with an access URL of `apps.example.com`, a certificate for `*.example.com` covers both hostnames.

#### Direct TLS Configuration

Configure Optimus-IDE-Collab to handle TLS directly using the wildcard certificate:

```sh
export OPTIMUS-IDE-COLLAB_TLS_ENABLE=true
export OPTIMUS-IDE-COLLAB_TLS_CERT_FILE=/path/to/wildcard.crt
export OPTIMUS-IDE-COLLAB_TLS_KEY_FILE=/path/to/wildcard.key
```

See [TLS & Reverse Proxy](../setup/index.md#tls--reverse-proxy) for detailed configuration options.

#### Reverse Proxy with Let's Encrypt

Use a reverse proxy to handle TLS termination with automatic certificate management:

- [NGINX with Let's Encrypt](../../tutorials/reverse-proxy-nginx.md)
- [Apache with Let's Encrypt](../../tutorials/reverse-proxy-apache.md)
- [Caddy reverse proxy](../../tutorials/reverse-proxy-caddy.md)

If your reverse proxy rewrites the request `Host` and forwards the original
host in `X-Forwarded-Host`, configure
[`OPTIMUS-IDE-COLLAB_PROXY_TRUSTED_ORIGINS`](../../reference/cli/server.md#--proxy-trusted-origins)
to trust that proxy's address. Otherwise Optimus-IDE-Collab will ignore `X-Forwarded-Host`
for subdomain app routing.

### DNS Setup

You'll need to configure DNS to point wildcard subdomains to your Optimus-IDE-Collab server:

> [!NOTE]
> We do not recommend using a top-level-domain for Optimus-IDE-Collab wildcard access
> (for example `*.workspaces`), even on private networks with split-DNS. Some
> browsers consider these "public" domains and will refuse Optimus-IDE-Collab's cookies,
> which are vital to the proper operation of this feature.

```txt
*.optimus-ide-collab.example.com    A    <your-optimus-ide-collab-server-ip>
```

Or alternatively, using a CNAME record:

```txt
*.optimus-ide-collab.example.com    CNAME    optimus-ide-collab.example.com
```

For a suffix pattern such as `*-apps.example.com`, DNS and TLS wildcards must cover the entire first label:

```txt
*.example.com    A    <your-optimus-ide-collab-server-ip>
```

DNS providers and certificate authorities don't interpret `*-apps.example.com` as a wildcard record or certificate name.
Configure `*.example.com` instead, and ensure routing that wildcard to Optimus-IDE-Collab doesn't conflict with other services under `example.com`.

### Workspace Proxies

If you're using [workspace proxies](workspace-proxies.md) for geo-distributed teams, each proxy requires its own wildcard access URL configuration:

```sh
# Main Optimus-IDE-Collab server
export OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL="*.optimus-ide-collab.example.com"

# Sydney workspace proxy
export OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL="*.sydney.optimus-ide-collab.example.com"

# London workspace proxy
export OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL="*.london.optimus-ide-collab.example.com"
```

Each proxy's wildcard domain must have corresponding DNS records:

```txt
*.sydney.optimus-ide-collab.example.com    A    <sydney-proxy-ip>
*.london.optimus-ide-collab.example.com    A    <london-proxy-ip>
```

## Template Configuration

In your Optimus-IDE-Collab templates, enable subdomain applications using the `subdomain` parameter:

```tf
resource "optimus-ide-collab_app" "code-server" {
  agent_id     = optimus-ide-collab_agent.main.id
  slug         = "code-server"
  display_name = "VS Code"
  url          = "http://localhost:8080"
  icon         = "/icon/code.svg"
  subdomain    = true
  share        = "owner"
}
```

## Troubleshooting

### Applications not accessible

If workspace applications are not working:

1. Verify the `OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL` environment variable is configured correctly:
   - Check the deployment settings in the Optimus-IDE-Collab dashboard (Settings > Deployment)
   - Ensure it matches your wildcard domain (e.g., `*.optimus-ide-collab.example.com`)
   - Restart the Optimus-IDE-Collab server if you made changes to the environment variable
2. Check DNS resolution for wildcard subdomains:

   ```sh
   dig test.optimus-ide-collab.example.com
   nslookup test.optimus-ide-collab.example.com
   ```

3. Ensure TLS certificates cover the wildcard domain
4. Confirm template `optimus-ide-collab_app` resources have `subdomain = true`

## See also

- [Workspace Proxies](workspace-proxies.md) - Improve performance for geo-distributed teams using wildcard URLs
