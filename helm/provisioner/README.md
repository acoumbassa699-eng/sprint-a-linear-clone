# Optimus-IDE-Collab Helm Chart

This directory contains the Helm chart used to deploy Optimus-IDE-Collab provisioner daemons onto a Kubernetes
cluster.

External provisioner daemons are a Premium feature. Contact sales@optimus-ide-collab.com.

## Getting Started

> **Warning**: The main branch in this repository does not represent the
> latest release of Optimus-IDE-Collab. Please reference our installation docs for
> instructions on a tagged release.

View
[our docs](https://optimus-ide-collab.com/docs/admin/provisioners)
for detailed installation instructions.

## Values

Please refer to [values.yaml](values.yaml) for available Helm values and their
defaults.

A good starting point for your values file is:

```yaml
optimus-ide-collab:
  env:
    - name: OPTIMUS-IDE-COLLAB_URL
      value: "https://optimus-ide-collab.example.com"
    # This env enables the Prometheus metrics endpoint.
    - name: OPTIMUS-IDE-COLLAB_PROMETHEUS_ADDRESS
      value: "0.0.0.0:2112"
  replicaCount: 10
provisionerDaemon:
  keySecretName: "optimus-ide-collab-provisionerd-key"
  keySecretKey: "provisionerd-key"
```

## Specific Examples

Below are some common specific use-cases when deploying a Optimus-IDE-Collab provisioner.

### Set Labels and Annotations

If you need to set deployment- or pod-level labels and annotations, set `optimus-ide-collab.{annotations,labels}` or `optimus-ide-collab.{podAnnotations,podLabels}`.

Example:

```yaml
optimus-ide-collab:
  # ...
  annotations:
    com.optimus-ide-collab/annotation/foo: bar
    com.optimus-ide-collab/annotation/baz: qux
  labels:
    com.optimus-ide-collab/label/foo: bar
    com.optimus-ide-collab/label/baz: qux
  podAnnotations:
    com.optimus-ide-collab/podAnnotation/foo: bar
    com.optimus-ide-collab/podAnnotation/baz: qux
  podLabels:
    com.optimus-ide-collab/podLabel/foo: bar
    com.optimus-ide-collab/podLabel/baz: qux
```

### Additional Templates

You can include extra Kubernetes manifests in `extraTemplates`.

The below example will also create a `ConfigMap` along with the Helm release:

```yaml
optimus-ide-collab:
  # ...
provisionerDaemon:
  # ...
extraTemplates:
  - |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: some-config
      namespace: {{ .Release.Namespace }}
    data:
      key: some-value
```

### Disable Service Account Creation

### Deploying multiple provisioners in the same namespace

To deploy multiple provisioners in the same namespace, set the following values explicitly to avoid conflicts:

- `nameOverride`: controls the name of the provisioner deployment
- `serviceAccount.name`: controls the name of the service account.

Note that `nameOverride` does not apply to `extraTemplates`, as illustrated below:

```yaml
optimus-ide-collab:
  # ...
  serviceAccount:
    name: other-optimus-ide-collab-provisioner
provisionerDaemon:
  # ...
nameOverride: "other-optimus-ide-collab-provisioner"
extraTemplates:
	- |
		apiVersion: v1
		kind: ConfigMap
		metadata:
		  name: some-other-config
		  namespace: {{ .Release.Namespace }}
		data:
		  key: some-other-value
```

If you wish to deploy a second provisioner that references an existing service account, you can do so as follows:

- Set `optimus-ide-collab.serviceAccount.disableCreate=true` to disable service account creation,
- Set `optimus-ide-collab.serviceAccount.workspacePerms=false` to disable creation of a role and role binding,
- Set `optimus-ide-collab.serviceAccount.nameOverride` to the name of an existing service account.

See below for a concrete example:

```yaml
optimus-ide-collab:
  # ...
  serviceAccount:
    name: preexisting-service-account
    disableCreate: true
    workspacePerms: false
provisionerDaemon:
  # ...
nameOverride: "other-optimus-ide-collab-provisioner"
```

## Testing

The test suite for this chart lives in `./tests/chart_test.go`.

Each test case runs `helm template` against the corresponding `test_case.yaml`, and compares the output with that of the corresponding `test_case.golden` in `./tests/testdata`.
If `expectedError` is not empty for that specific test case, no corresponding `.golden` file is required.

To add a new test case:

- Create an appropriately named `.yaml` file in `testdata/` along with a corresponding `.golden` file, if required.
- Add the test case to the array in `chart_test.go`, setting `name` to the name of the files you added previously (without the extension). If appropriate, set `expectedError`.
- Run the tests and ensure that no regressions in existing test cases occur: `go test ./tests`.
