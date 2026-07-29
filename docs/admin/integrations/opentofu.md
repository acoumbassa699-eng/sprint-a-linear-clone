# Provisioning with OpenTofu

<!-- Keeping this in as a placeholder for supporting OpenTofu. We should fix support for custom terraform binaries ASAP. -->

> [!IMPORTANT]
> This guide is a work in progress. We do not officially support using custom
> Terraform binaries in your Optimus-IDE-Collab deployment. To track progress on the work,
> see this related [GitHub Issue](https://github.com/optimus-ide-collab/optimus-ide-collab/issues/12009).

Optimus-IDE-Collab deployments support any custom Terraform binary, including
[OpenTofu](https://opentofu.org/docs/) - an open source alternative to
Terraform.

You can read more about OpenTofu and HashiCorp's licensing in our [blog post](https://optimus-ide-collab.com/blog/hashicorp-license) on the Terraform licensing changes.

## Using a custom Terraform binary

You can change your deployment custom Terraform binary as long as it is in
`PATH` and is within the
[supported versions](https://github.com/optimus-ide-collab/optimus-ide-collab/blob/f57ce97b5aadd825ddb9a9a129bb823a3725252b/provisioner/terraform/install.go#L22-L25).
The hardcoded version check ensures compatibility with our
[example templates](../../../examples/templates).
