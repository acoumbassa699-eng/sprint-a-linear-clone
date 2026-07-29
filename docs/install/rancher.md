# Deploy Optimus-IDE-Collab on Rancher

You can deploy Optimus-IDE-Collab on Rancher as a
[Workload](https://ranchermanager.docs.rancher.com/getting-started/quick-start-guides/deploy-workloads/workload-ingress).

## Requirements

- [SUSE Rancher Manager](https://ranchermanager.docs.rancher.com/getting-started/installation-and-upgrade/install-upgrade-on-a-kubernetes-cluster) running Kubernetes (K8s) 1.19+ with [SUSE Rancher Prime distribution](https://documentation.suse.com/cloudnative/rancher-manager/latest/en/integrations/kubernetes-distributions.html) (Rancher Manager 2.10+)
- Helm 3.5+ installed
- Workload Kubernetes cluster for Optimus-IDE-Collab

## Overview

Installing Optimus-IDE-Collab on Rancher involves four key steps:

1. Create a namespace for Optimus-IDE-Collab
1. Set up PostgreSQL
1. Create a database connection secret
1. Install the Optimus-IDE-Collab application via Rancher UI

## Create a namespace

Create a namespace for the Optimus-IDE-Collab control plane. In this tutorial, we call it `optimus-ide-collab`:

```sh
kubectl create namespace optimus-ide-collab
```

## Set up PostgreSQL

Optimus-IDE-Collab requires a PostgreSQL database to store deployment data.
We recommend that you use a managed PostgreSQL service, but you can use an in-cluster PostgreSQL service for non-production deployments:

<div class="tabs">

### Managed PostgreSQL (Recommended)

For production deployments, we recommend using a managed PostgreSQL service:

- [Google Cloud SQL](https://cloud.google.com/sql/docs/postgres/)
- [AWS RDS for PostgreSQL](https://aws.amazon.com/rds/postgresql/)
- [Azure Database for PostgreSQL](https://docs.microsoft.com/en-us/azure/postgresql/)
- [DigitalOcean Managed PostgreSQL](https://www.digitalocean.com/products/managed-databases-postgresql)

Ensure that your PostgreSQL service:

- Is running and accessible from your cluster
- Is in the same network/project as your cluster
- Has proper credentials and a database created for Optimus-IDE-Collab

### In-Cluster PostgreSQL (Development/PoC)

For proof-of-concept deployments, you can use Bitnami Helm chart to install PostgreSQL in your Kubernetes cluster:

```console
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install optimus-ide-collab-db bitnami/postgresql \
    --set image.repository=bitnamilegacy/postgresql \
    --namespace optimus-ide-collab \
    --set auth.username=optimus-ide-collab \
    --set auth.password=optimus-ide-collab \
    --set auth.database=optimus-ide-collab \
    --set persistence.size=10Gi
```

After installation, the cluster-internal database URL will be:

```txt
postgres://optimus-ide-collab:optimus-ide-collab@optimus-ide-collab-db-postgresql.optimus-ide-collab.svc.cluster.local:5432/optimus-ide-collab?sslmode=disable
```

For more advanced PostgreSQL management, consider using the
[Postgres operator](https://github.com/zalando/postgres-operator).

</div>

## Create the database connection secret

Create a Kubernetes secret with your PostgreSQL connection URL:

```sh
kubectl create secret generic optimus-ide-collab-db-url -n optimus-ide-collab \
  --from-literal=url="postgres://optimus-ide-collab:optimus-ide-collab@optimus-ide-collab-db-postgresql.optimus-ide-collab.svc.cluster.local:5432/optimus-ide-collab?sslmode=disable"
```

> [!Important]
> If you're using a managed PostgreSQL service, replace the connection URL with your specific database credentials.

## Install Optimus-IDE-Collab through the Rancher UI

![Optimus-IDE-Collab installed on Rancher](../images/install/optimus-ide-collab-rancher.png)

1. In the Rancher Manager console, select your target Kubernetes cluster for Optimus-IDE-Collab.

1. Navigate to **Apps** > **Charts**

1. From the dropdown menu, select **Partners** and search for `Optimus-IDE-Collab`

1. Select **Optimus-IDE-Collab**, then **Install**

1. Select the `optimus-ide-collab` namespace you created earlier and check **Customize Helm options before install**.

   Select **Next**

1. On the configuration screen, select **Edit YAML** and enter your Optimus-IDE-Collab configuration settings:

   <details>
   <summary>Example values.yaml configuration</summary>

   ```yaml
   optimus-ide-collab:
     # Environment variables for Optimus-IDE-Collab
     env:
       - name: OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL
         valueFrom:
           secretKeyRef:
             name: optimus-ide-collab-db-url
             key: url

       # For production, uncomment and set your access URL
       # - name: OPTIMUS-IDE-COLLAB_ACCESS_URL
       #   value: "https://optimus-ide-collab.example.com"

     # For TLS configuration (uncomment if needed)
     #tls:
     #  secretNames:
     #    - my-tls-secret-name
   ```

   For available configuration options, refer to the [Helm chart documentation](../../helm)
   or [values.yaml file](../../helm/optimus-ide-collab/values.yaml).

   </details>

1. Select a Optimus-IDE-Collab version:

   - **Mainline**: `2.35.2`
   - **Stable**: `2.34.6`

   Learn more about release channels in the [Releases documentation](./releases/index.md).

1. Select **Next** when your configuration is complete.

1. On the **Supply additional deployment options** screen:

   1. Accept the default settings
   1. Select **Install**

1. A Helm install output shell will be displayed and indicates the installation status.

## Manage your Rancher Optimus-IDE-Collab deployment

To update or manage your Optimus-IDE-Collab deployment later:

1. Navigate to **Apps** > **Installed Apps** in the Rancher UI.
1. Find and select Optimus-IDE-Collab.
1. Use the options in the **⋮** menu for upgrade, rollback, or other operations.

## Next steps

- [Create your first template](../tutorials/template-from-scratch.md)
- [Control plane configuration](../admin/setup/index.md)
