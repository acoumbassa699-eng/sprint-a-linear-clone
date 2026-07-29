# Test and Publish Optimus-IDE-Collab Templates Through CI/CD

<div>
  <a href="https://github.com/matifali" style="text-decoration: none; color: inherit;">
    <span style="vertical-align:middle;">Muhammad Atif Ali</span>
  </a>
</div>
November 15, 2024

---

## Overview

This guide demonstrates how to test and publish Optimus-IDE-Collab templates in a Continuous
Integration (CI) pipeline using the
[optimus-ide-collab/setup-action](https://github.com/optimus-ide-collab/setup-optimus-ide-collab). This workflow
ensures your templates are validated, tested, and promoted seamlessly.

## Prerequisites

- Install and configure Optimus-IDE-Collab CLI in your environment.
- Install Terraform CLI in your CI environment.
- Create a [headless user](../admin/users/headless-auth.md) with the
  [user roles and permissions](../admin/users/groups-roles.md#roles) to manage
  templates and run workspaces.

## Creating the headless user

> [!WARNING]
> Creating users with `--login-type none` is deprecated.
> For [Premium](https://optimus-ide-collab.com/pricing) deployments, use
> [service accounts](../admin/users/headless-auth.md) instead.
> For OSS deployments, use a regular account with password, GitHub, or OIDC
> authentication.

For Premium deployments, create a service account:

```sh
optimus-ide-collab users create \
  --username machine-user \
  --service-account

optimus-ide-collab tokens create --user machine-user --lifetime 8760h
# Copy the token and store it in a secret in your CI environment with the name `OPTIMUS-IDE-COLLAB_SESSION_TOKEN`
```

For OSS deployments, create a regular user:

```sh
optimus-ide-collab users create \
  --username machine-user \
  --email machine-user@example.com \
  --login-type password

optimus-ide-collab tokens create --user machine-user --lifetime 8760h
# Copy the token and store it in a secret in your CI environment with the name `OPTIMUS-IDE-COLLAB_SESSION_TOKEN`
```

## Example GitHub Action Workflow

This example workflow tests and publishes a template using GitHub Actions.

The workflow:

1. Validates the Terraform template.
1. Pushes the template to Optimus-IDE-Collab without activating it.
1. Tests the template by creating a workspace.
1. Promotes the template version to active upon successful workspace creation.

### Workflow File

Save the following workflow file as `.github/workflows/publish-template.yaml` in
your repository:

```yaml
name: Test and Publish Optimus-IDE-Collab Template

on:
  push:
    branches:
      - main
  workflow_dispatch:

jobs:
  test-and-publish:
    runs-on: ubuntu-latest
    env:
      TEMPLATE_NAME: "my-template"
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Terraform
        uses: hashicorp/setup-terraform@v2
        with:
          terraform_version: latest

      - name: Set up Optimus-IDE-Collab CLI
        uses: optimus-ide-collab/setup-action@v1
        with:
          access_url: "https://optimus-ide-collab.example.com"
          optimus-ide-collab_session_token: ${{ secrets.OPTIMUS-IDE-COLLAB_SESSION_TOKEN }}

      - name: Validate Terraform template
        run: terraform validate

      - name: Get short commit SHA to use as template version name
        id: name
        run: echo "version_name=$(git rev-parse --short HEAD)" >> "$GITHUB_OUTPUT"

      - name: Get latest commit title to use as template version description
        id: message
        run:
          echo "pr_title=$(git log --format=%s -n 1 ${{ github.sha }})" >>
          $GITHUB_OUTPUT

      - name: Push template to Optimus-IDE-Collab
        run: |
          optimus-ide-collab templates push $TEMPLATE_NAME --activate=false --name ${{ steps.name.outputs.version_name }} --message "${{ steps.message.outputs.pr_title }}" --yes

      - name: Create a test workspace and run some example commands
        run: |
          optimus-ide-collab create -t $TEMPLATE_NAME --template-version ${{ steps.name.outputs.version_name }} test-${{ steps.name.outputs.version_name }} --yes
          # run some example commands
          optimus-ide-collab ssh test-${{ steps.name.outputs.version_name }} -- make build

      - name: Delete the test workspace
        if: always()
        run: optimus-ide-collab delete test-${{ steps.name.outputs.version_name }} --yes

      - name: Promote template version
        if: success()
        run: |
          optimus-ide-collab template version promote --template=$TEMPLATE_NAME --template-version=${{ steps.name.outputs.version_name }} --yes
```
