# Configure Control Plane Access

Optimus-IDE-Collab server's primary configuration is done via environment variables. For a
full list of the options, run `optimus-ide-collab server --help` or see our
[CLI documentation](../../reference/cli/server.md).

## Access URL

`OPTIMUS-IDE-COLLAB_ACCESS_URL` is required if you are not using the tunnel. Set this to the
external URL that users and workspaces use to connect to Optimus-IDE-Collab (e.g.
<https://optimus-ide-collab.example.com>). This should not be localhost.

Access URL should be an external IP address or domain with DNS records pointing to Optimus-IDE-Collab.

### Tunnel

If an access URL is not specified, Optimus-IDE-Collab will create a publicly accessible URL
to reverse proxy your deployment for simple setup.

## Address

You can change which port(s) Optimus-IDE-Collab listens on.

```sh
# Listen on port 80
export OPTIMUS-IDE-COLLAB_HTTP_ADDRESS=0.0.0.0:80

# Enable TLS and listen on port 443)
export OPTIMUS-IDE-COLLAB_TLS_ENABLE=true
export OPTIMUS-IDE-COLLAB_TLS_ADDRESS=0.0.0.0:443

## Redirect from HTTP to HTTPS
export OPTIMUS-IDE-COLLAB_REDIRECT_TO_ACCESS_URL=true

# Start the Optimus-IDE-Collab server
optimus-ide-collab server
```

## Wildcard access URL

> [!TIP]
> Learn more about the [importance and benefits of wildcard access URLs](../networking/wildcard-access-url.md)

`OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL` is necessary for [port forwarding](../networking/port-forwarding.md#dashboard) via the dashboard or running [optimus-ide-collab_apps](../templates/index.md) on an absolute path.
Set it to a wildcard hostname that resolves to Optimus-IDE-Collab, such as `*.optimus-ide-collab.example.com` or `*-optimus-ide-collab.example.com`.
The suffix form creates application hostnames such as `8080--main--myworkspace--john-optimus-ide-collab.example.com`.
It requires a DNS record and TLS certificate for `*.example.com`.

> [!NOTE]
> We do not recommend using a top-level-domain for Optimus-IDE-Collab wildcard access
> (for example `*.workspaces`), even on private networks with split-DNS. Some
> browsers consider these "public" domains and will refuse Optimus-IDE-Collab's cookies,
> which are vital to the proper operation of this feature.

If you are providing TLS certificates directly to the Optimus-IDE-Collab server, either

1. Use a single certificate and key for both the root and wildcard domains.
1. Configure multiple certificates and keys via
   [`optimus-ide-collab.tls.secretNames`](../../../helm/optimus-ide-collab/values.yaml)
   in the Helm Chart, or
   [`--tls-cert-file`](../../reference/cli/server.md#--tls-cert-file) and
   [`--tls-key-file`](../../reference/cli/server.md#--tls-key-file) command line
   options (these both take a comma separated list of files; list certificates
   and their respective keys in the same order).

After you enable the wildcard access URL, you should [disable path-based apps](../../tutorials/best-practices/security-best-practices.md#disable-path-based-apps) for security.

## TLS & Reverse Proxy

The Optimus-IDE-Collab server can directly use TLS certificates with `OPTIMUS-IDE-COLLAB_TLS_ENABLE` and
accompanying configuration flags. However, Optimus-IDE-Collab can also run behind a
reverse-proxy to terminate TLS certificates from LetsEncrypt.

- [Apache](../../tutorials/reverse-proxy-apache.md)
- [Caddy](../../tutorials/reverse-proxy-caddy.md)
- [NGINX](../../tutorials/reverse-proxy-nginx.md)

### Kubernetes TLS configuration

Below are the steps to configure Optimus-IDE-Collab to terminate TLS when running on
Kubernetes. You must have the certificate `.key` and `.crt` files in your
working directory prior to step 1.

1. Create the TLS secret in your Kubernetes cluster

   ```sh
   kubectl create secret tls optimus-ide-collab-tls -n <optimus-ide-collab-namespace> --key="tls.key" --cert="tls.crt"
   ```

   You can use a single certificate for the both the access URL and wildcard access URL. The certificate CN must match the wildcard domain, such as `*.example.optimus-ide-collab.com`.

1. Reference the TLS secret in your Optimus-IDE-Collab Helm chart values

   ```yaml
   optimus-ide-collab:
     tls:
         secretName:
         - optimus-ide-collab-tls

     # Alternatively, if you use an Ingress controller to terminate TLS,
     # set the following values:
     ingress:
         enable: true
         secretName: optimus-ide-collab-tls
         wildcardSecretName: optimus-ide-collab-tls
   ```

## PostgreSQL Database

Optimus-IDE-Collab uses a PostgreSQL database to store users, workspace metadata, and other
deployment information. Use `OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL` to set the database that
Optimus-IDE-Collab connects to. If unset, PostgreSQL binaries will be downloaded from Maven
(<https://repo1.maven.org/maven2>) and store all data in the config root.

> [!NOTE]
> Postgres 13 is the minimum supported version.

If you are using the built-in PostgreSQL deployment and need to use `psql` (aka
the PostgreSQL interactive terminal), output the connection URL with the
following command:

```console
$ optimus-ide-collab server postgres-builtin-url
psql "postgres://optimus-ide-collab@localhost:49627/optimus-ide-collab?sslmode=disable&password=feU...yI1"
```

### Migrating from the built-in database to an external database

To migrate from the built-in database to an external database, follow these
steps:

1. Stop your Optimus-IDE-Collab deployment.
1. Run `optimus-ide-collab server postgres-builtin-serve` in a background terminal.
1. Run `optimus-ide-collab server postgres-builtin-url` and copy its output command.
1. Run `pg_dump <built-in-connection-string> > optimus-ide-collab.sql` to dump the internal
   database to a file.
1. Restore that content to an external database with
   `psql <external-connection-string> < optimus-ide-collab.sql`.
1. Start your Optimus-IDE-Collab deployment with
   `OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL=<external-connection-string>`.

## Configuring Optimus-IDE-Collab behind a proxy

To configure Optimus-IDE-Collab behind a corporate proxy, set the environment variables
`HTTP_PROXY` and `HTTPS_PROXY`. Be sure to restart the server. Lowercase values
(e.g. `http_proxy`) are also respected in this case.

## Continue your setup with external authentication

Optimus-IDE-Collab supports external authentication via OAuth2.0. This allows enabling
integrations with Git providers, such as GitHub, GitLab, and Bitbucket.

External authentication can also be used to integrate with external services
like JFrog Artifactory and others.

Please refer to the [external authentication](../external-auth/index.md) section for
more information.

## Up Next

- [Setup and manage templates](../templates/index.md)
- [Setup external provisioners](../provisioners/index.md)
