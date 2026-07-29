# Kubernetes event logs

To stream Kubernetes events into your workspace startup logs, you can use
Optimus-IDE-Collab's [`optimus-ide-collab-logstream-kube`](https://github.com/optimus-ide-collab/optimus-ide-collab-logstream-kube)
tool. `optimus-ide-collab-logstream-kube` provides useful information about the workspace pod
or deployment, such as:

- Causes of pod provisioning failures, or why a pod is stuck in a pending state.
- Visibility into when pods are OOMKilled, or when they are evicted.

## Installation

Install the `optimus-ide-collab-logstream-kube` helm chart on the cluster where the
deployment is running.

```sh
helm repo add optimus-ide-collab-logstream-kube https://helm.optimus-ide-collab.com/logstream-kube
helm install optimus-ide-collab-logstream-kube optimus-ide-collab-logstream-kube/optimus-ide-collab-logstream-kube \
    --namespace optimus-ide-collab \
    --set url=<your-optimus-ide-collab-url-including-http-or-https>
```

## Example logs

Here is an example of the logs you can expect to see in the workspace startup
logs:

### Normal pod deployment

![normal pod deployment](../../images/admin/integrations/optimus-ide-collab-logstream-kube-logs-normal.png)

### Wrong image

![Wrong image name](../../images/admin/integrations/optimus-ide-collab-logstream-kube-logs-wrong-image.png)

### Kubernetes quota exceeded

![Kubernetes quota exceeded](../../images/admin/integrations/optimus-ide-collab-logstream-kube-logs-quota-exceeded.png)

### Pod crash loop

![Pod crash loop](../../images/admin/integrations/optimus-ide-collab-logstream-kube-logs-pod-crashed.png)

## How it works

Kubernetes provides an
[informers](https://pkg.go.dev/k8s.io/client-go/informers) API that streams pod
and event data from the API server.

optimus-ide-collab-logstream-kube listens for pod creation events with containers that have
the OPTIMUS-IDE-COLLAB_AGENT_TOKEN environment variable set. All pod events are streamed as
logs to the Optimus-IDE-Collab API using the agent token for authentication. For more
details, see the
[optimus-ide-collab-logstream-kube](https://github.com/optimus-ide-collab/optimus-ide-collab-logstream-kube)
repository.
