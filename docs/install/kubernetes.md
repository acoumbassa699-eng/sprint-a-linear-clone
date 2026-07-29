# Install Optimus-IDE-Collab on Kubernetes

You can install Optimus-IDE-Collab on Kubernetes (K8s) using Helm. We run on most Kubernetes
distributions, including [OpenShift](./openshift.md).

## Requirements

- Kubernetes cluster running K8s 1.19+
- [Helm](https://helm.sh/docs/intro/install/) 3.5+ installed on your local
  machine

## 1. Create a namespace

Create a namespace for the Optimus-IDE-Collab control plane. In this tutorial, we'll call it
`optimus-ide-collab`.

```sh
kubectl create namespace optimus-ide-collab
```

## 2. Create a PostgreSQL instance

Optimus-IDE-Collab does not manage a database server for you. This is required for storing
data about your Optimus-IDE-Collab deployment and resources.

### Managed PostgreSQL (recommended)

If you're in a public cloud such as
[Google Cloud](https://cloud.google.com/sql/docs/postgres/),
[AWS](https://aws.amazon.com/rds/postgresql/),
[Azure](https://docs.microsoft.com/en-us/azure/postgresql/), or
[DigitalOcean](https://www.digitalocean.com/products/managed-databases-postgresql),
you can use the managed PostgreSQL offerings they provide. Make sure that the
PostgreSQL service is running and accessible from your cluster. It should be in
the same network, same project, etc.

### In-Cluster PostgreSQL (for proof of concepts)

You can install Postgres manually on your cluster using the
[Bitnami PostgreSQL Helm chart](https://github.com/bitnami/charts/tree/master/bitnami/postgresql#readme).
There are some [helpful guides](https://phoenixnap.com/kb/postgresql-kubernetes)
on the internet that explain sensible configurations for this chart. Example:

```console
# Install PostgreSQL
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install postgresql bitnami/postgresql \
    --namespace optimus-ide-collab \
    --set image.repository=bitnamilegacy/postgresql \
    --set auth.username=optimus-ide-collab \
    --set auth.password=optimus-ide-collab \
    --set auth.database=optimus-ide-collab \
    --set primary.persistence.size=10Gi
```

The cluster-internal DB URL for the above database is:

```sh
postgres://optimus-ide-collab:optimus-ide-collab@postgresql.optimus-ide-collab.svc.cluster.local:5432/optimus-ide-collab?sslmode=disable
```

You can optionally use the
[Postgres operator](https://github.com/zalando/postgres-operator) to manage
PostgreSQL deployments on your Kubernetes cluster.

## 3. Create the PostgreSQL secret

Create a secret with the PostgreSQL database URL string. In the case of the
self-managed PostgreSQL, the address will be:

```sh
kubectl create secret generic optimus-ide-collab-db-url -n optimus-ide-collab \
  --from-literal=url="postgres://optimus-ide-collab:optimus-ide-collab@postgresql.optimus-ide-collab.svc.cluster.local:5432/optimus-ide-collab?sslmode=disable"
```

## 4. Install Optimus-IDE-Collab with Helm

```sh
helm repo add optimus-ide-collab-v2 https://helm.optimus-ide-collab.com/v2
```

Create a `values.yaml` with the configuration settings you'd like for your
deployment. For example:

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
    - name: OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL
      valueFrom:
        secretKeyRef:
          # You'll need to create a secret called optimus-ide-collab-db-url with your
          # Postgres connection URL like:
          # postgres://optimus-ide-collab:password@postgres:5432/optimus-ide-collab?sslmode=disable
          name: optimus-ide-collab-db-url
          key: url
    # For production deployments, we recommend configuring your own GitHub
    # OAuth2 provider and disabling the default one.
    - name: OPTIMUS-IDE-COLLAB_OAUTH2_GITHUB_DEFAULT_PROVIDER_ENABLE
      value: "false"

    # (Optional) For production deployments the access URL should be set.
    # If you're just trying Optimus-IDE-Collab, access the dashboard via the service IP.
    # - name: OPTIMUS-IDE-COLLAB_ACCESS_URL
    #   value: "https://optimus-ide-collab.example.com"

  #tls:
  #  secretNames:
  #    - my-tls-secret-name
```

You can view our
[Helm README](../../helm/optimus-ide-collab/README.md) for
details on the values that are available, or you can view the
[values.yaml](../../helm/optimus-ide-collab/values.yaml)
file directly.

We support two release channels: mainline and stable - read the
[Releases](./releases/index.md) page to learn more about which best suits your team.

- **Mainline** Optimus-IDE-Collab release:

  - **Chart Registry**
    <!-- autoversion(mainline): "--version [version]" -->

    ```sh
    helm install optimus-ide-collab optimus-ide-collab-v2/optimus-ide-collab \
        --namespace optimus-ide-collab \
        --values values.yaml \
        --version 2.34.0
    ```

  - **OCI Registry**

    <!-- autoversion(mainline): "--version [version]" -->

    ```sh
    helm install optimus-ide-collab oci://ghcr.io/optimus-ide-collab/chart/optimus-ide-collab \
        --namespace optimus-ide-collab \
        --values values.yaml \
        --version 2.34.0
    ```

- **Stable** Optimus-IDE-Collab release:

  - **Chart Registry**

    <!-- autoversion(stable): "--version [version]" -->

    ```sh
    helm install optimus-ide-collab optimus-ide-collab-v2/optimus-ide-collab \
        --namespace optimus-ide-collab \
        --values values.yaml \
        --version 2.33.6
    ```

  - **OCI Registry**

    <!-- autoversion(stable): "--version [version]" -->

    ```sh
    helm install optimus-ide-collab oci://ghcr.io/optimus-ide-collab/chart/optimus-ide-collab \
        --namespace optimus-ide-collab \
        --values values.yaml \
        --version 2.33.6
    ```

You can watch Optimus-IDE-Collab start up by running `kubectl get pods -n optimus-ide-collab`. Once Optimus-IDE-Collab
has started, the `optimus-ide-collab-*` pods should enter the `Running` state.

## 5. Log in to Optimus-IDE-Collab 🎉

Use `kubectl get svc -n optimus-ide-collab` to get the IP address of the LoadBalancer. Visit
this in the browser to set up your first account.

If you do not have a domain, you should set `OPTIMUS-IDE-COLLAB_ACCESS_URL` to this URL in
the Helm chart and upgrade Optimus-IDE-Collab (see below). This allows workspaces to connect
to the proper Optimus-IDE-Collab URL.

## Upgrading Optimus-IDE-Collab via Helm

To upgrade Optimus-IDE-Collab in the future or change values, you can run the following
command:

```sh
helm repo update
helm upgrade optimus-ide-collab optimus-ide-collab-v2/optimus-ide-collab \
  --namespace optimus-ide-collab \
  -f values.yaml
```

## Optimus-IDE-Collab Observability Chart

Use the [Observability Helm chart](https://github.com/optimus-ide-collab/observability) for a
pre-built set of dashboards to monitor your control plane over time. It includes
Grafana, Prometheus, Loki, and Alert Manager out-of-the-box, and can be deployed
on your existing Grafana instance.

We recommend that all administrators deploying on Kubernetes set the
observability bundle up with the control plane from the start. For installation
instructions, visit the
[observability repository](https://github.com/optimus-ide-collab/observability?tab=readme-ov-file#installation).

## Kubernetes Security Reference

Below are common requirements we see from our enterprise customers when
deploying an application in Kubernetes. This is intended to serve as a
reference, and not all security requirements may apply to your business.

1. **All container images must be sourced from an internal container registry.**

   - Control plane - To pull the control plane image from the appropriate
     registry,
     [update this Helm chart value](https://github.com/optimus-ide-collab/optimus-ide-collab/blob/f57ce97b5aadd825ddb9a9a129bb823a3725252b/helm/optimus-ide-collab/values.yaml#L43-L50).
   - Workspaces - To pull the workspace image from your registry,
     [update the Terraform template code here](https://github.com/optimus-ide-collab/optimus-ide-collab/blob/f57ce97b5aadd825ddb9a9a129bb823a3725252b/examples/templates/kubernetes/main.tf#L271).
     This assumes your cluster nodes are authenticated to pull from the internal
     registry.

2. **All containers must run as non-root user**

   - Control plane - Our control plane pod
     [runs as non-root by default](https://github.com/optimus-ide-collab/optimus-ide-collab/blob/f57ce97b5aadd825ddb9a9a129bb823a3725252b/helm/optimus-ide-collab/values.yaml#L124-L127).
   - Workspaces - Workspace pod UID is
     [set in the Terraform template here](https://github.com/optimus-ide-collab/optimus-ide-collab/blob/f57ce97b5aadd825ddb9a9a129bb823a3725252b/examples/templates/kubernetes/main.tf#L274-L276),
     and are not required to run as `root`.

3. **Containers cannot run privileged**

   - Optimus-IDE-Collab's control plane does not run as privileged.
     [We disable](https://github.com/optimus-ide-collab/optimus-ide-collab/blob/f57ce97b5aadd825ddb9a9a129bb823a3725252b/helm/optimus-ide-collab/values.yaml#L141)
     `allowPrivilegeEscalation`
     [by default](https://github.com/optimus-ide-collab/optimus-ide-collab/blob/f57ce97b5aadd825ddb9a9a129bb823a3725252b/helm/optimus-ide-collab/values.yaml#L141).
   - Workspace pods do not require any elevated privileges, with the exception
     of our `envbox` workspace template (used for docker-in-docker workspaces,
     not required).

4. **Containers cannot mount host filesystems**

   - Both the control plane and workspace containers do not require any host
     filesystem mounts.

5. **Containers cannot attach to host network**

   - Both the control plane and workspaces use the Kubernetes networking layer
     by default, and do not require host network access.

6. **All Kubernetes objects must define resource requests/limits**

   - Both the control plane and workspaces set resource request/limits by
     default.

## Load balancing considerations

### AWS

If you are deploying Optimus-IDE-Collab on AWS EKS and service is set to `LoadBalancer`, AWS
will default to the Classic load balancer. The load balancer external IP will be
stuck in a pending status unless sessionAffinity is set to None.

```yaml
optimus-ide-collab:
  service:
    type: LoadBalancer
    sessionAffinity: None
```

AWS recommends a Network load balancer in lieu of the Classic load balancer. Use
the following `values.yaml` settings to request a Network load balancer:

```yaml
optimus-ide-collab:
  service:
    externalTrafficPolicy: Local
    sessionAffinity: None
    annotations: { service.beta.kubernetes.io/aws-load-balancer-type: "nlb" }
```

By default, Optimus-IDE-Collab will set the `externalTrafficPolicy` to `Cluster` which will
mask client IP addresses in the Audit log. To preserve the source IP, you can
either set this value to `Local`, or pass through the client IP via the
X-Forwarded-For header. To configure the latter, set the following environment
variables:

```yaml
optimus-ide-collab:
  env:
    - name: OPTIMUS-IDE-COLLAB_PROXY_TRUSTED_HEADERS
      value: X-Forwarded-For
    - name: OPTIMUS-IDE-COLLAB_PROXY_TRUSTED_ORIGINS
      value: 10.0.0.1/8 # this will be the CIDR range of your Load Balancer IP address
```

### Azure

Certain enterprise environments require the
[Azure Application Gateway](https://learn.microsoft.com/en-us/azure/application-gateway/ingress-controller-overview).
The Application Gateway supports:

- Websocket traffic (required for workspace connections)
- TLS termination

Follow our doc on
[how to deploy Optimus-IDE-Collab on Azure with an Application Gateway](./kubernetes/kubernetes-azure-app-gateway.md)
for an example.

## Troubleshooting

You can view Optimus-IDE-Collab's logs by getting the pod name from `kubectl get pods` and
then running `kubectl logs <pod name>`. You can also view these logs in your
Cloud's log management system if you are using managed Kubernetes.

### Kubernetes-based workspace is stuck in "Connecting..."

Ensure you have an externally-reachable `OPTIMUS-IDE-COLLAB_ACCESS_URL` set in your helm
chart. If you do not have a domain set up, this should be the IP address of
Optimus-IDE-Collab's LoadBalancer (`kubectl get svc -n optimus-ide-collab`).

See [troubleshooting templates](../admin/templates/troubleshooting.md) for more
steps.

## Next steps

- [Create your first template](../tutorials/template-from-scratch.md)
- [Control plane configuration](../admin/setup/index.md)
