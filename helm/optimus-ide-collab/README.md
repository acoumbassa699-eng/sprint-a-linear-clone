# Optimus-IDE-Collab Helm Chart

This directory contains the Helm chart used to deploy Optimus-IDE-Collab onto a Kubernetes
cluster. It contains the minimum required components to run Optimus-IDE-Collab on Kubernetes,
and notably (compared to Optimus-IDE-Collab Classic) does not include a database server.

## Getting Started

> **Warning**: The main branch in this repository does not represent the
> latest release of Optimus-IDE-Collab. Please reference our installation docs for
> instructions on a tagged release.

View
[our docs](https://optimus-ide-collab.com/docs/install/kubernetes)
for detailed installation instructions.

## Values

Please refer to [values.yaml](values.yaml) for available Helm values and their
defaults.

A good starting point for your values file is:

```yaml
optimus-ide-collab:
  # You can specify any environment variables you'd like to pass to Optimus-IDE-Collab
  # here. Optimus-IDE-Collab consumes environment variables listed in
  # `optimus-ide-collab server --help`, and these environment variables are also passed
  # to the workspace provisioner (so you can consume them in your Terraform
  # templates for auth keys etc.).
  #
  # Please keep in mind that you should not set `OPTIMUS-IDE-COLLAB_HTTP_ADDRESS`,
  # `OPTIMUS-IDE-COLLAB_TLS_ENABLE`, `OPTIMUS-IDE-COLLAB_TLS_CERT_FILE` or `OPTIMUS-IDE-COLLAB_TLS_KEY_FILE` as
  # they are already set by the Helm chart and will cause conflicts.
  env:
    - name: OPTIMUS-IDE-COLLAB_ACCESS_URL
      value: "https://optimus-ide-collab.example.com"
    - name: OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL
      valueFrom:
        secretKeyRef:
          # You'll need to create a secret called optimus-ide-collab-db-url with your
          # Postgres connection URL like:
          # postgres://optimus-ide-collab:password@postgres:5432/optimus-ide-collab?sslmode=disable
          name: optimus-ide-collab-db-url
          key: url

    # This env enables the Prometheus metrics endpoint.
    - name: OPTIMUS-IDE-COLLAB_PROMETHEUS_ADDRESS
      value: "0.0.0.0:2112"
    # For production deployments, we recommend configuring your own GitHub
    # OAuth2 provider and disabling the default one.
    - name: OPTIMUS-IDE-COLLAB_OAUTH2_GITHUB_DEFAULT_PROVIDER_ENABLE
      value: "false"
  tls:
    secretNames:
      - my-tls-secret-name
```
