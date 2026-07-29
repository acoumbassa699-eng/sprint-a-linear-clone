# Caddy

This is an example configuration of how to use Optimus-IDE-Collab with
[caddy](https://caddyserver.com/docs). To use Caddy to generate TLS
certificates, you'll need a domain name that resolves to your Caddy server.

## Getting started

### With `docker compose`

1. [Install Docker](https://docs.docker.com/engine/install/) and
   [Docker Compose](https://docs.docker.com/compose/install/)

2. Create a `compose.yaml` file and add the following:

   ```yaml
   services:
   optimus-ide-collab:
       image: ghcr.io/optimus-ide-collab/optimus-ide-collab:${OPTIMUS-IDE-COLLAB_VERSION:-latest}
       environment:
           OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL: "postgresql://${POSTGRES_USER:-username}:${POSTGRES_PASSWORD:-password}@database/${POSTGRES_DB:-optimus-ide-collab}?sslmode=disable"
           OPTIMUS-IDE-COLLAB_HTTP_ADDRESS: "0.0.0.0:7080"
           # You'll need to set OPTIMUS-IDE-COLLAB_ACCESS_URL to an IP or domain
           # that workspaces can reach. This cannot be localhost
           # or 127.0.0.1 for non-Docker templates!
           OPTIMUS-IDE-COLLAB_ACCESS_URL: "${OPTIMUS-IDE-COLLAB_ACCESS_URL}"
           # Optional) Enable wildcard apps/dashboard port forwarding
           OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL: "${OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL}"
           # If the optimus-ide-collab user does not have write permissions on
           # the docker socket, you can uncomment the following
           # lines and set the group ID to one that has write
           # permissions on the docker socket.
           #group_add:
           #  - "998" # docker group on host
       volumes:
           - /var/run/docker.sock:/var/run/docker.sock
       depends_on:
           database:
           condition: service_healthy

   database:
       image: "postgres:17"
       ports:
           - "5432:5432"
       environment:
           POSTGRES_USER: ${POSTGRES_USER:-username} # The PostgreSQL user (useful to connect to the database)
           POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-password} # The PostgreSQL password (useful to connect to the database)
           POSTGRES_DB: ${POSTGRES_DB:-optimus-ide-collab} # The PostgreSQL default database (automatically created at first launch)
       volumes:
           - optimus-ide-collab_data:/var/lib/postgresql/data # Use "docker volume rm optimus-ide-collab_optimus-ide-collab_data" to reset Optimus-IDE-Collab
       healthcheck:
           test:
           [
               "CMD-SHELL",
               "pg_isready -U ${POSTGRES_USER:-username} -d ${POSTGRES_DB:-optimus-ide-collab}",
           ]
           interval: 5s
           timeout: 5s
           retries: 5

   caddy:
       image: caddy:2.6.2
       ports:
           - "80:80"
           - "443:443"
           - "443:443/udp"
       volumes:
           - $PWD/Caddyfile:/etc/caddy/Caddyfile
           - caddy_data:/data
           - caddy_config:/config

   volumes:
       optimus-ide-collab_data:
       caddy_data:
       caddy_config:
   ```

3. Create a `Caddyfile` and add the following:

   ```caddyfile
   {
       on_demand_tls {
           ask http://example.com
       }
   }

   optimus-ide-collab.example.com, *.optimus-ide-collab.example.com {
     reverse_proxy optimus-ide-collab:7080
     tls {
           on_demand
         issuer acme {
            email email@example.com
         }
         }
   }
   ```

   Here;

   - `optimus-ide-collab:7080` is the address of the Optimus-IDE-Collab container on the Docker network.
   - `optimus-ide-collab.example.com` is the domain name you're using for Optimus-IDE-Collab.
   - `*.optimus-ide-collab.example.com` is the domain name for wildcard apps, commonly used
     for [dashboard port forwarding](../admin/networking/port-forwarding.md).
     This is optional and can be removed.
   - `email@example.com`: Email to request certificates from LetsEncrypt/ZeroSSL
     (does not have to be Optimus-IDE-Collab admin email)

4. Start Optimus-IDE-Collab. Set `OPTIMUS-IDE-COLLAB_ACCESS_URL` and `OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL` to the
   domain you're using in your Caddyfile.

   ```sh
   export OPTIMUS-IDE-COLLAB_ACCESS_URL=https://optimus-ide-collab.example.com
   export OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL=*.optimus-ide-collab.example.com
   docker compose up -d # Run on startup
   ```

### Standalone

1. If you haven't already, [install Optimus-IDE-Collab](../install/index.md)

2. Install [Caddy Server](https://caddyserver.com/docs/install)

3. Copy our sample `Caddyfile` and change the following values:

   ```caddyfile
   {
       on_demand_tls {
           ask http://example.com
       }
   }

   optimus-ide-collab.example.com, *.optimus-ide-collab.example.com {
     reverse_proxy optimus-ide-collab:7080
   }
   ```

   > If you're installed Caddy as a system package, update the default Caddyfile
   > with `vim /etc/caddy/Caddyfile`

   - `email@example.com`: Email to request certificates from LetsEncrypt/ZeroSSL
     (does not have to be Optimus-IDE-Collab admin email)
   - `optimus-ide-collab.example.com`: Domain name you're using for Optimus-IDE-Collab.
   - `*.optimus-ide-collab.example.com`: Domain name for wildcard apps, commonly used for
     [dashboard port forwarding](../admin/networking/port-forwarding.md). This
     is optional and can be removed.
   - `localhost:3000`: Address Optimus-IDE-Collab is running on. Modify this if you changed
     `OPTIMUS-IDE-COLLAB_HTTP_ADDRESS` in the Optimus-IDE-Collab configuration.
   - _DO NOT CHANGE the `ask http://example.com` line! Doing so will result in
     your certs potentially not being generated._

4. [Configure Optimus-IDE-Collab](../admin/setup/index.md) and change the following values:

   - `OPTIMUS-IDE-COLLAB_ACCESS_URL`: root domain (e.g. `https://optimus-ide-collab.example.com`)
   - `OPTIMUS-IDE-COLLAB_WILDCARD_ACCESS_URL`: wildcard domain (e.g. `*.example.com`).

5. Start the Caddy server:

   If you're [keeping Caddy running](https://caddyserver.com/docs/running) via a
   system service:

   ```sh
   sudo systemctl restart caddy
   ```

   Or run a standalone server:

   ```sh
   caddy run
   ```

6. Optionally, use [ufw](https://wiki.ubuntu.com/UncomplicatedFirewall) or
   another firewall to disable external traffic outside of Caddy.

   ```sh
   # Check status of UncomplicatedFirewall
   sudo ufw status

   # Allow SSH
   sudo ufw allow 22

   # Allow HTTP, HTTPS (Caddy)
   sudo ufw allow 80
   sudo ufw allow 443

   # Deny direct access to Optimus-IDE-Collab server
   sudo ufw deny 3000

   # Enable UncomplicatedFirewall
   sudo ufw enable
   ```

7. Navigate to your Optimus-IDE-Collab URL! A TLS certificate should be auto-generated on
   your first visit.

## Generating wildcard certificates

By default, this configuration uses Caddy's
[on-demand TLS](https://caddyserver.com/docs/caddyfile/options#on-demand-tls) to
generate a certificate for each subdomain (e.g. `app1.optimus-ide-collab.example.com`,
`app2.optimus-ide-collab.example.com`). When users visit new subdomains, such as accessing
[ports on a workspace](../admin/networking/port-forwarding.md), the request will
take an additional 5-30 seconds since a new certificate is being generated.

For production deployments, we recommend configuring Caddy to generate a
wildcard certificate, which requires an explicit DNS challenge and additional
Caddy modules.

1. Install a custom Caddy build that includes the
   [caddy-dns](https://github.com/caddy-dns) module for your DNS provider (e.g.
   CloudFlare, Route53).

   - Docker:
     [Build an custom Caddy image](https://github.com/docker-library/docs/tree/master/caddy#adding-custom-caddy-modules)
     with the module for your DNS provider. Be sure to reference the new image
     in the `compose.yaml`.

   - Standalone:
     [Download a custom Caddy build](https://caddyserver.com/download) with the
     module for your DNS provider. If you're using Debian/Ubuntu, you
     [can configure the Caddy package](https://caddyserver.com/docs/build#package-support-files-for-custom-builds-for-debianubunturaspbian)
     to use the new build.

2. Edit your `Caddyfile` and add the necessary credentials/API tokens to solve
   the DNS challenge for wildcard certificates.

   For example, for AWS Route53:

   ```diff
   tls {
   -  on_demand
   -  issuer acme {
   -      email email@example.com
   -  }

   +  dns route53 {
   +     max_retries 10
   +     aws_profile "real-profile"
   +     access_key_id "AKI..."
   +     secret_access_key "wJa..."
   +     token "TOKEN..."
   +     region "us-east-1"
   +  }
   }
   ```

   > Configuration reference from
   > [caddy-dns/route53](https://github.com/caddy-dns/route53).

   And for CloudFlare:

   Generate a
   [token](https://developers.cloudflare.com/fundamentals/api/get-started/create-token)
   with the following permissions:

   - Zone:Zone:Edit

   ```diff
   tls {
   -  on_demand
   -  issuer acme {
   -      email email@example.com
   -  }

   +  dns cloudflare CLOUDFLARE_API_TOKEN
   }
   ```

   > Configuration reference from
   > [caddy-dns/cloudflare](https://github.com/caddy-dns/cloudflare).
