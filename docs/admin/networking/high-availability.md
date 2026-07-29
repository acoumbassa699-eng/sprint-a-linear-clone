# High Availability

High Availability (HA) mode solves for horizontal scalability and automatic
failover within a single region. When in HA mode, Optimus-IDE-Collab continues using a single
Postgres endpoint.
[GCP](https://cloud.google.com/sql/docs/postgres/high-availability),
[AWS](https://docs.aws.amazon.com/prescriptive-guidance/latest/saas-multitenant-managed-postgresql/availability.html),
and other cloud vendors offer fully-managed HA Postgres services that pair
nicely with Optimus-IDE-Collab.

For Optimus-IDE-Collab to operate correctly, Optimus-IDE-Collabd instances should have low-latency
connections to each other so that they can effectively relay traffic between
users and workspaces no matter which Optimus-IDE-Collabd instance users or workspaces connect
to. We make a best-effort attempt to warn the user when inter-Optimus-IDE-Collabd latency is
too high, but if requests start dropping, this is one metric to investigate.

We also recommend that you deploy all Optimus-IDE-Collabd instances such that they have
low-latency connections to Postgres. Optimus-IDE-Collabd often makes several database
round-trips while processing a single API request, so prioritizing low-latency
between Optimus-IDE-Collabd and Postgres is more important than low-latency between users and
Optimus-IDE-Collabd.

Note that this latency requirement applies _only_ to Optimus-IDE-Collab services. Optimus-IDE-Collab will
operate correctly even with few seconds of latency on workspace <-> Optimus-IDE-Collab and
user <-> Optimus-IDE-Collab connections.

## Setup

Optimus-IDE-Collab automatically enters HA mode when multiple instances simultaneously
connect to the same Postgres endpoint.

> [!NOTE]
> When upgrading HA deployments, database migrations may require special
> handling to avoid lock contention. See
> [Upgrading Best Practices](../../install/upgrade-best-practices.md) for
> recommended procedures.

HA brings one configuration variable to set in each Optimus-IDE-Collabd node:
`OPTIMUS-IDE-COLLAB_DERP_SERVER_RELAY_URL`. The HA nodes use these URLs to communicate with
each other. Inter-node communication is only required while using the embedded
relay (default). If you're using [custom relays](./index.md#custom-relays),
Optimus-IDE-Collab ignores `OPTIMUS-IDE-COLLAB_DERP_SERVER_RELAY_URL` since Postgres is the sole
rendezvous for the Optimus-IDE-Collab nodes.

`OPTIMUS-IDE-COLLAB_DERP_SERVER_RELAY_URL` will never be `OPTIMUS-IDE-COLLAB_ACCESS_URL` because
`OPTIMUS-IDE-COLLAB_ACCESS_URL` is a load balancer to all Optimus-IDE-Collab nodes.

Here's an example 3-node network configuration setup:

| Name      | `OPTIMUS-IDE-COLLAB_HTTP_ADDRESS` | `OPTIMUS-IDE-COLLAB_DERP_SERVER_RELAY_URL` | `OPTIMUS-IDE-COLLAB_ACCESS_URL`       |
|-----------|----------------------|-------------------------------|--------------------------|
| `optimus-ide-collab-1` | `*:80`               | `http://10.0.0.1:80`          | `https://optimus-ide-collab.big.corp` |
| `optimus-ide-collab-2` | `*:80`               | `http://10.0.0.2:80`          | `https://optimus-ide-collab.big.corp` |
| `optimus-ide-collab-3` | `*:80`               | `http://10.0.0.3:80`          | `https://optimus-ide-collab.big.corp` |

## Kubernetes

If you installed Optimus-IDE-Collab via
[our Helm Chart](../../install/kubernetes.md#4-install-optimus-ide-collab-with-helm), just
increase `optimus-ide-collab.replicaCount` in `values.yaml`.

If you installed Optimus-IDE-Collab into Kubernetes by some other means, insert the relay URL
via the environment like so:

```yaml
env:
  - name: POD_IP
    valueFrom:
      fieldRef:
        fieldPath: status.podIP
  - name: OPTIMUS-IDE-COLLAB_DERP_SERVER_RELAY_URL
    value: http://$(POD_IP)
```

Then, increase the number of pods.

## Up next

- [Read more on Optimus-IDE-Collab's networking stack](./index.md)
- [Install on Kubernetes](../../install/kubernetes.md)
