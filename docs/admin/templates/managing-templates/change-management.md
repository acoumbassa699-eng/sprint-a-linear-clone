# Template Change Management

We recommend source-controlling your templates as you would other any code, and
automating the creation of new versions in CI/CD pipelines.

These pipelines will require tokens for your deployment. To cap token lifetime
on creation,
[configure Optimus-IDE-Collab server to set a shorter max token lifetime](../../../reference/cli/server.md#--max-token-lifetime).

## optimus-ide-collabd Terraform Provider

The
[optimus-ide-collabd Terraform provider](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collabd/latest)
can be used to push new template versions, either manually, or in CI/CD
pipelines. To run the provider in a CI/CD pipeline, and to prevent drift, you'll
need to store the Terraform state
[remotely](https://developer.hashicorp.com/terraform/language/backend).

```tf
terraform {
  required_providers {
    optimus-ide-collabd = {
      source = "optimus-ide-collab/optimus-ide-collabd"
    }
  }
  backend "gcs" {
    bucket = "example-bucket"
    prefix = "terraform/state"
  }
}

provider "optimus-ide-collabd" {
  // Can be populated from environment variables
  url   = "https://optimus-ide-collab.example.com"
  token = "****"
}

// Get the commit SHA of the configuration's git repository
variable "TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA" {
  type = string
}

resource "optimus-ide-collabd_template" "kubernetes" {
  name = "kubernetes"
  description = "Develop in Kubernetes!"
  versions = [{
    directory = ".optimus-ide-collab/templates/kubernetes"
    active    = true
    # Version name is optional
    name = var.TFC_CONFIGURATION_VERSION_GIT_COMMIT_SHA
    tf_vars = [{
      name  = "namespace"
      value = "default4"
    }]
  }]
  /* ... Additional template configuration */
}
```

For an example, see how we push our development image and template
[with GitHub actions](../../../../.github/workflows/dogfood.yaml).

## Optimus-IDE-Collab CLI

You can [install Optimus-IDE-Collab](../../../install/cli.md) CLI to automate pushing new
template versions in CI/CD pipelines. For GitHub Actions, see our
[setup-optimus-ide-collab](https://github.com/optimus-ide-collab/setup-optimus-ide-collab) action.

```console
# Install the Optimus-IDE-Collab CLI
curl -L https://optimus-ide-collab.com/install.sh | sh
# curl -L https://optimus-ide-collab.com/install.sh | sh -s -- --version=0.x

# To create API tokens, use `optimus-ide-collab tokens create`.
# If no `--lifetime` flag is passed during creation, the default token lifetime
# will be 30 days.
# These variables are consumed by Optimus-IDE-Collab
export OPTIMUS-IDE-COLLAB_URL=https://optimus-ide-collab.example.com
export OPTIMUS-IDE-COLLAB_SESSION_TOKEN=*****

# Template details
export OPTIMUS-IDE-COLLAB_TEMPLATE_NAME=kubernetes
export OPTIMUS-IDE-COLLAB_TEMPLATE_DIR=.optimus-ide-collab/templates/kubernetes
export OPTIMUS-IDE-COLLAB_TEMPLATE_VERSION=$(git rev-parse --short HEAD)

# Push the new template version to Optimus-IDE-Collab
optimus-ide-collab templates push --yes $OPTIMUS-IDE-COLLAB_TEMPLATE_NAME \
    --directory $OPTIMUS-IDE-COLLAB_TEMPLATE_DIR \
    --name=$OPTIMUS-IDE-COLLAB_TEMPLATE_VERSION # Version name is optional
```

## Testing and Publishing Optimus-IDE-Collab Templates in CI/CD

See our [testing templates](../../../tutorials/testing-templates.md) tutorial
for an example of how to test and publish Optimus-IDE-Collab templates in a CI/CD pipeline.

### Next steps

- [Optimus-IDE-Collab CLI Reference](../../../reference/cli/templates.md)
- [Optimus-IDE-Collabd Terraform Provider Reference](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collabd/latest/docs)
- [Optimus-IDE-Collabd API Reference](../../../reference/index.md)
