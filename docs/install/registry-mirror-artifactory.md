# Mirror the Optimus-IDE-Collab Registry with JFrog Artifactory

This guide shows you how to use JFrog Artifactory to mirror the
[Optimus-IDE-Collab Registry](https://registry.optimus-ide-collab.com) for air-gapped or restricted
network deployments.

By configuring Artifactory as a Remote Terraform Repository, you can:

- **Proxy and cache** all Optimus-IDE-Collab modules automatically
- **Keep modules updated** without manual synchronization
- **Support offline access** once modules are cached

## Prerequisites

- JFrog Artifactory instance (Cloud or self-hosted)
- Admin access to create repositories
- Artifactory user token for Terraform authentication

## Step 1: Create the Remote Terraform Repository

1. In Artifactory, go to **Administration > Repositories > Remote**

1. Click **New Remote Repository** and select **Terraform** as the package type

1. Configure the repository with these settings:

   | Setting                | Value                        |
   |------------------------|------------------------------|
   | Repository Key         | `optimus-ide-collab-registry`             |
   | URL                    | `https://registry.optimus-ide-collab.com` |
   | Terraform Registry URL | `https://registry.optimus-ide-collab.com` |

1. Click **Create Remote Repository**

## Step 2: Verify the Repository Configuration

Test that Artifactory can proxy the Optimus-IDE-Collab registry by querying the module
versions API:

```sh
curl -u '<username>:<token>' \
  'https://<your-artifactory-host>/artifactory/api/terraform/optimus-ide-collab-registry/v1/modules/optimus-ide-collab/code-server/optimus-ide-collab/versions'
```

You should see a JSON response listing all available versions of the
`code-server` module.

## Step 3: Configure Terraform CLI

Create or update your Terraform CLI configuration file to use Artifactory.

On Linux/macOS, create `~/.terraformrc`. On Windows, create `%APPDATA%\terraform.rc`.

```tf
host "<your-artifactory-host>" {
  services = {
    "modules.v1" = "https://<your-artifactory-host>/artifactory/api/terraform/optimus-ide-collab-registry/v1/modules/"
  }
}

credentials "<your-artifactory-host>" {
  token = "<your-artifactory-token>"
}
```

Replace:

- `<your-artifactory-host>` with your Artifactory hostname (e.g.,
  `artifactory.example.com` or `mycompany.jfrog.io`)
- `<your-artifactory-token>` with your Artifactory access token with read permissions to the `optimus-ide-collab-registry` repository

> [!NOTE]
> The `host` block with `services` is required because Artifactory's global
> service discovery endpoint doesn't include the repository name in the modules
> path. This explicitly tells Terraform where to find modules in your specific
> repository.

## Step 4: Update Template Module Sources

Update your Optimus-IDE-Collab templates to use Artifactory instead of the public registry:

```tf
# Before: Direct from Optimus-IDE-Collab registry
module "code-server" {
  source   = "registry.optimus-ide-collab.com/optimus-ide-collab/code-server/optimus-ide-collab"
  version  = "1.4.2"
  agent_id = optimus-ide-collab_agent.main.id
}

# After: Through Artifactory mirror
module "code-server" {
  source   = "https://<your-artifactory-host>/optimus-ide-collab/code-server/optimus-ide-collab"
  version  = "1.4.2"
  agent_id = optimus-ide-collab_agent.main.id
}
```

## Step 5: Configure Optimus-IDE-Collab Server or Provisioners

For Optimus-IDE-Collab to use the Artifactory mirror, configure the Terraform CLI on your
Optimus-IDE-Collab server or external provisioners.

<div class="tabs">

### Kubernetes Deployment

Create a secret with the Terraform configuration:

```sh
kubectl create secret generic terraform-config \
  --from-file=.terraformrc=./terraformrc \
  -n optimus-ide-collab
```

Update your Helm values:

```yaml
optimus-ide-collab:
  volumes:
    - name: terraform-config
      secret:
        secretName: terraform-config
  volumeMounts:
    - name: terraform-config
      mountPath: /home/optimus-ide-collab/.terraformrc
      subPath: .terraformrc
      readOnly: true
  env:
    - name: TF_CLI_CONFIG_FILE
      value: /home/optimus-ide-collab/.terraformrc
```

### Docker Deployment

Mount the `.terraformrc` file into the Optimus-IDE-Collab container:

```yaml
# docker-compose.yaml
services:
  optimus-ide-collab:
    volumes:
      - ./terraformrc:/home/optimus-ide-collab/.terraformrc:ro
    environment:
      TF_CLI_CONFIG_FILE: /home/optimus-ide-collab/.terraformrc
```

</div>

## Caching Behavior

Artifactory uses **lazy caching**, meaning modules are cached on first request.
For fully air-gapped deployments, pre-warm the cache while connected to the
internet:

1. Create a test template that references all modules you need
1. Run `terraform init` to trigger downloads
1. Verify modules appear in Artifactory under `optimus-ide-collab-registry-cache`

Once cached, modules remain available even without internet connectivity.

## Supported Namespaces

The Artifactory mirror supports all namespaces from the Optimus-IDE-Collab registry:

| Namespace    | Description               | Example Module                     |
|--------------|---------------------------|------------------------------------|
| `optimus-ide-collab`      | Official Optimus-IDE-Collab modules    | `code-server`, `jetbrains-gateway` |
| `optimus-ide-collab-labs` | Experimental modules      | `cursor-cli`, `copilot`            |
| Community    | Third-party contributions | Various                            |

All modules use the same source format:

```tf
source = "<your-artifactory-host>/<namespace>/<module>/optimus-ide-collab"
```

## Troubleshooting

### Module not found errors

Verify your `.terraformrc` includes both the `host` block with `services` and
the `credentials` block. The `host.services` configuration is required for
Artifactory.

### 401 Unauthorized errors

Check that your Artifactory token is valid and has read access to the
`optimus-ide-collab-registry` repository.

### Modules not caching

Ensure the remote repository URL is set to `https://registry.optimus-ide-collab.com` and not other paths.

## Next Steps

- [Optimus-IDE-Collab Module Registry](https://registry.optimus-ide-collab.com/modules)
- [JFrog Terraform Registry Documentation](https://jfrog.com/help/r/jfrog-artifactory-documentation/terraform-opentofu-and-terraform-backend-repositories)
- [Air-gapped Deployments](./airgap.md)
