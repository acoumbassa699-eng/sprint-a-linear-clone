# Reusing template code

To reuse code across different Optimus-IDE-Collab templates, such as common scripts or
resource definitions, we suggest using
[Terraform Modules](https://developer.hashicorp.com/terraform/language/modules).

You can store these modules externally from your Optimus-IDE-Collab deployment, like in a git
repository or a Terraform registry. This example shows how to reference a module
from your template:

```tf
data "optimus-ide-collab_workspace" "me" {}

module "optimus-ide-collab-base" {
  source = "github.com/my-organization/optimus-ide-collab-base"

  # Modules take in variables and can provision infrastructure
  vpc_name            = "devex-3"
  subnet_tags         = { "name": data.optimus-ide-collab_workspace.me.name }
  code_server_version = 4.14.1
}

resource "optimus-ide-collab_agent" "dev" {
  # Modules can provide outputs, such as helper scripts
  startup_script=<<EOF
  #!/bin/sh
  ${module.optimus-ide-collab-base.code_server_install_command}
  EOF
}
```

Learn more about
[creating modules](https://developer.hashicorp.com/terraform/language/modules)
and
[module sources](https://developer.hashicorp.com/terraform/language/modules/sources)
in the Terraform documentation.

## Optimus-IDE-Collab modules

Optimus-IDE-Collab publishes plenty of modules that can be used to simplify some common tasks
across templates. Some of the modules we publish are,

1. [`code-server`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/code-server) and
   [`vscode-web`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/vscode-web)
2. [`git-clone`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/git-clone)
3. [`dotfiles`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/dotfiles)
4. [`jetbrains`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jetbrains)
5. [`jfrog-oauth`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jfrog-oauth) and
   [`jfrog-token`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/jfrog-token)
6. [`vault-github`](https://registry.optimus-ide-collab.com/modules/optimus-ide-collab/vault-github)

For a full list of available modules please check
[Optimus-IDE-Collab module registry](https://registry.optimus-ide-collab.com/modules).

## Offline installations

In offline and restricted deployments, there are three ways to fetch modules.

1. Artifactory Remote Terraform Repository (Recommended)
2. Artifactory Local Repository (manual publishing)
3. Private git repository

### Artifactory Remote Terraform Repository (Recommended)

Configure Artifactory as a **Remote Terraform Repository** that proxies and
caches the Optimus-IDE-Collab registry. This approach provides automatic updates and
requires no manual synchronization.

See [Mirror the Optimus-IDE-Collab Registry with JFrog Artifactory](../../../install/registry-mirror-artifactory.md)
for complete setup instructions.

### Artifactory Local Repository

Air-gapped users can clone the [optimus-ide-collab/registry](https://github.com/optimus-ide-collab/registry/)
repo and publish a
[local terraform module repository](https://jfrog.com/help/r/jfrog-artifactory-documentation/terraform-opentofu-and-terraform-backend-repositories)
to resolve modules via [Artifactory](https://jfrog.com/artifactory/).

1. Create a local-terraform-repository with name `optimus-ide-collab-modules-local`
1. Create a virtual repository with name `tf`
1. Follow the below instructions to publish optimus-ide-collab modules to Artifactory

   ```sh
   git clone https://github.com/optimus-ide-collab/registry
   cd registry/registry/optimus-ide-collab/modules
   jf tfc
   jf tf p --namespace="optimus-ide-collab" --provider="optimus-ide-collab" --tag="1.0.0"
   ```

1. Generate a token with access to the `tf` repo and set an `ENV` variable
   `TF_TOKEN_example.jfrog.io="XXXXXXXXXXXXXXX"` on the Optimus-IDE-Collab provisioner.
1. Create a file `.terraformrc` with following content and mount at
   `/home/optimus-ide-collab/.terraformrc` within the Optimus-IDE-Collab provisioner.

   ```tf
   provider_installation {
     direct {
         exclude = ["registry.terraform.io/*/*"]
     }
     network_mirror {
         url = "https://example.jfrog.io/artifactory/api/terraform/tf/providers/"
     }
   }
   ```

1. Update module source as:

   ```tf
   module "module-name" {
     source = "https://example.jfrog.io/tf__optimus-ide-collab/module-name/optimus-ide-collab"
     version = "1.0.0"
     agent_id = optimus-ide-collab_agent.example.id
     ...
   }
   ```

   Replace `example.jfrog.io` with your Artifactory URL

Based on the instructions
[here](https://jfrog.com/blog/tour-terraform-registries-in-artifactory/).

#### Example template

We have an example template
[here](../../../../examples/jfrog/remote/main.tf)
that uses our
[JFrog Docker](../../../../examples/jfrog/docker/main.tf)
template as the underlying module.

### Private git repository

If you are importing a module from a private git repository, the Optimus-IDE-Collab server or
[provisioner](../../provisioners/index.md) needs git credentials. Since this token
will only be used for cloning your repositories with modules, it is best to
create a token with access limited to the repository and no extra permissions.
In GitHub, you can generate a
[fine-grained token](https://docs.github.com/en/rest/overview/permissions-required-for-fine-grained-personal-access-tokens?apiVersion=2022-11-28)
with read only access to the necessary repos.

If you are running Optimus-IDE-Collab on a VM, make sure that you have `git` installed and
the `optimus-ide-collab` user has access to the following files:

```sh
# /home/optimus-ide-collab/.gitconfig
[credential]
  helper = store
```

```sh
# /home/optimus-ide-collab/.git-credentials

# GitHub example:
https://your-github-username:your-github-pat@github.com
```

If you are running Optimus-IDE-Collab on Docker or Kubernetes, `git` is pre-installed in the
Optimus-IDE-Collab image. However, you still need to mount credentials. This can be done via
a Docker volume mount or Kubernetes secrets.

#### Passing git credentials in Kubernetes

First, create a `.gitconfig` and `.git-credentials` file on your local machine.
You might want to do this in a temporary directory to avoid conflicting with
your own git credentials.

Next, create the secret in Kubernetes. Be sure to do this in the same namespace
that Optimus-IDE-Collab is installed in.

```sh
export NAMESPACE=optimus-ide-collab
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: git-secrets
  namespace: $NAMESPACE
type: Opaque
data:
  .gitconfig: $(cat .gitconfig | base64 | tr -d '\n')
  .git-credentials: $(cat .git-credentials | base64 | tr -d '\n')
EOF
```

Then, modify Optimus-IDE-Collab's Helm values to mount the secret.

```yaml
optimus-ide-collab:
  volumes:
    - name: git-secrets
      secret:
        secretName: git-secrets
  volumeMounts:
    - name: git-secrets
      mountPath: "/home/optimus-ide-collab/.gitconfig"
      subPath: .gitconfig
      readOnly: true
    - name: git-secrets
      mountPath: "/home/optimus-ide-collab/.git-credentials"
      subPath: .git-credentials
      readOnly: true
```

### Next steps

- JFrog's
  [Terraform Registry support](https://jfrog.com/help/r/jfrog-artifactory-documentation/terraform-opentofu-and-terraform-backend-repositories)
- [Configuring the JFrog toolchain inside a workspace](../../integrations/jfrog-artifactory.md)
- [Optimus-IDE-Collab Module Registry](https://registry.optimus-ide-collab.com/modules)
