# Networking

Optimus-IDE-Collab's network topology has three types of nodes: workspaces, optimus-ide-collab servers,
and users.

The optimus-ide-collab server must have an inbound address reachable by users and workspaces,
but otherwise, all topologies _just work_ with Optimus-IDE-Collab.

When possible, we establish direct connections between users and workspaces.
Direct connections are as fast as connecting to the workspace outside of Optimus-IDE-Collab.
When NAT traversal fails, connections are relayed through the optimus-ide-collab server. All
user-workspace connections are end-to-end encrypted.

[Tailscale's open source](https://tailscale.com) backs our websocket/HTTPS
networking logic.

## Requirements

In order for clients and workspaces to be able to connect:

> [!NOTE]
> We strongly recommend that clients connect to Optimus-IDE-Collab and their
> workspaces over a good quality, broadband network connection. The following
> are minimum requirements:
>
> - better than 400ms round-trip latency to the Optimus-IDE-Collab server and to their
>   workspace
> - better than 0.5% random packet loss

- All clients and agents must be able to establish a connection to the Optimus-IDE-Collab
  server (`OPTIMUS-IDE-COLLAB_ACCESS_URL`) over HTTP/HTTPS.
- Any reverse proxy or ingress between the Optimus-IDE-Collab control plane and
  clients/agents must support WebSockets.

In order for clients to be able to establish direct connections:

> [!NOTE]
> Direct connections via the web browser are not supported. To improve
> latency for browser-based applications running inside Optimus-IDE-Collab workspaces in
> regions far from the Optimus-IDE-Collab control plane, consider deploying one or more
> [workspace proxies](./workspace-proxies.md).

- The client is connecting using the CLI (e.g. `optimus-ide-collab ssh` or
  `optimus-ide-collab port-forward`). Note that the
  [VSCode extension](https://marketplace.visualstudio.com/items?itemName=optimus-ide-collab.optimus-ide-collab-remote)
  and [JetBrains Plugin](https://plugins.jetbrains.com/plugin/19620-optimus-ide-collab/), and
  [`ssh optimus-ide-collab.<workspace>`](../../reference/cli/config-ssh.md) all utilize the
  CLI to establish a workspace connection.
- Either the client or workspace agent are able to discover a reachable
  `ip:port` of their counterpart. If the agent and client are able to
  communicate with each other using their locally assigned IP addresses, then a
  direct connection can be established immediately. Otherwise, the client and
  agent will contact
  [the configured STUN servers](../../reference/cli/server.md#--derp-server-stun-addresses)
  to try and determine which `ip:port` can be used to communicate with their
  counterpart. See [STUN and NAT](./stun.md) for more details on how this
  process works.
- All outbound UDP traffic must be allowed for both the client and the agent on
  **all ports** to each others' respective networks.
  - To establish a direct connection, both agent and client use STUN. This
    involves sending UDP packets outbound on `udp/3478` to the configured
    [STUN server](../../reference/cli/server.md#--derp-server-stun-addresses).
    If either the agent or the client are unable to send and receive UDP packets
    to a STUN server, then direct connections will not be possible.
  - Both agents and clients will then establish a
    [WireGuard](https://www.wireguard.com/)️ tunnel and send UDP traffic on
    ephemeral (high) ports. If a firewall between the client and the agent
    blocks this UDP traffic, direct connections will not be possible.

## optimus-ide-collab server

Workspaces connect to the optimus-ide-collab server via the server's external address, set
via [`ACCESS_URL`](../../admin/setup/index.md#access-url). There must not be a
NAT between workspaces and optimus-ide-collab server.

Users connect to the optimus-ide-collab server's dashboard and API through its `ACCESS_URL`
as well. There must not be a NAT between users and the optimus-ide-collab server.

Template admins can overwrite the site-wide access URL at the template level by
leveraging the `url` argument when
[defining the Optimus-IDE-Collab provider](https://registry.terraform.io/providers/optimus-ide-collab/optimus-ide-collab/latest/docs#url-1):

```tf
provider "optimus-ide-collab" {
  url = "https://optimus-ide-collab.namespace.svc.cluster.local"
}
```

This is useful when debugging connectivity issues between the workspace agent
and the Optimus-IDE-Collab server.

## Web Apps

The optimus-ide-collab servers relays dashboard-initiated connections between the user and
the workspace. Web terminal <-> workspace connections are an exception and may
be direct.

In general, [port forwarded](./port-forwarding.md) web apps are faster than
dashboard-accessed web apps.

## 🌎 Geo-distribution

### Direct connections

Direct connections are a straight line between the user and workspace, so there
is no special geo-distribution configuration. To speed up direct connections,
move the user and workspace closer together.

Establishing a direct connection can be an involved process because both the
client and workspace agent will likely be behind at least one level of NAT,
meaning that we need to use STUN to learn the IP address and port under which
the client and agent can both contact each other. See [STUN and NAT](./stun.md)
for more information on how this process works.

If a direct connection is not available (e.g. client or server is behind NAT),
Optimus-IDE-Collab will use a relayed connection. By default,
[Optimus-IDE-Collab uses Google's public STUN server](../../reference/cli/server.md#--derp-server-stun-addresses),
but this can be disabled or changed for
[Air-gapped deployments](../../install/airgap.md).

### Relayed connections

By default, your Optimus-IDE-Collab server also runs a built-in DERP relay which can be used
for both public and [Air-gapped deployments](../../install/airgap.md).

However, Tailscale maintains a global fleet of [DERP relays](https://tailscale.com/kb/1118/custom-derp-servers/#what-are-derp-servers) intended for their product, and has allowed Optimus-IDE-Collab to access and use them.
You can launch `optimus-ide-collab server` with Tailscale's DERPs like so:

```sh
optimus-ide-collab server --derp-config-url https://controlplane.tailscale.com/derpmap/default
```

#### Custom Relays

If you want lower latency than what Tailscale offers or want additional DERP
relays for air-gapped deployments, you may run custom DERP servers. Refer to
[Tailscale's documentation](https://tailscale.com/kb/1118/custom-derp-servers/#why-run-your-own-derp-server)
to learn how to set them up.

After you have custom DERP servers, you can launch Optimus-IDE-Collab with them like so:

```json
# derpmap.json
{
  "Regions": {
    "1": {
      "RegionID": 1,
      "RegionCode": "myderp",
      "RegionName": "My DERP",
      "Nodes": [
        {
          "Name": "1",
          "RegionID": 1,
          "HostName": "your-hostname.com"
        }
      ]
    }
  }
}
```

```sh
optimus-ide-collab server --derp-config-path derpmap.json
```

### Dashboard connections

The dashboard (and web apps opened through the dashboard) are served from the
optimus-ide-collab server, so they can only be geo-distributed with High Availability mode in
our Premium Edition. [Reach out to Sales](https://optimus-ide-collab.com/contact) to learn
more.

## Browser-only connections

> [!NOTE]
> Browser-only connections is a Premium feature.
> [Learn more](https://optimus-ide-collab.com/pricing#compare-plans).

Some Optimus-IDE-Collab deployments require that all access is through the browser to comply
with security policies. In these cases, pass the `--browser-only` flag to
`optimus-ide-collab server` or set `OPTIMUS-IDE-COLLAB_BROWSER_ONLY=true`.

With browser-only connections, developers can only connect to their workspaces
via the web terminal and
[web IDEs](../../user-guides/workspace-access/web-ides.md).

### Workspace Proxies

> [!NOTE]
> Workspace proxies are a Premium feature.
> [Learn more](https://optimus-ide-collab.com/pricing#compare-plans).

Workspace proxies are a Optimus-IDE-Collab Premium feature that allows you to provide
low-latency browser experiences for geo-distributed teams.

To learn more, see [Workspace Proxies](./workspace-proxies.md).

## Latency

Optimus-IDE-Collab measures and reports several types of latency, providing insights into the performance of your deployment. Understanding these metrics can help you diagnose issues and optimize the user experience.

There are three main types of latency metrics for your Optimus-IDE-Collab deployment:

- Dashboard-to-server latency:

  The Optimus-IDE-Collab UI measures round-trip time to the Optimus-IDE-Collab server or workspace proxy using built-in browser timing capabilities.

  This appears in the user interface next to your username, showing how responsive the dashboard is.

- Workspace connection latency:

  The latency shown on the workspace dashboard measures the round-trip time between the workspace agent and its DERP relay server.

  This metric is displayed in milliseconds on the workspace dashboard and specifically shows the agent-to-relay latency, not direct P2P connections.

  To estimate the total end-to-end latency experienced by a user, add the dashboard-to-server latency to this agent-to-relay latency.

- Database latency:

  For administrators, Optimus-IDE-Collab monitors and reports database query performance in the health dashboard.

### How latency is classified

Latency measurements are color-coded in the dashboard:

- **Green** (<150ms): Good performance.
- **Yellow** (150-300ms): Moderate latency that might affect user experience.
- **Red** (>300ms): High latency that will noticeably affect user experience.

### View latency information

- **Dashboard**: The global latency indicator appears in the top navigation bar.
- **Workspace list**: Each workspace shows its connection latency.
- **Health dashboard**: Administrators can view advanced metrics including database latency.
- **CLI**: Use `optimus-ide-collab ping <workspace>` to measure and analyze latency from the command line.

### Factors that affect latency

- **Geographic distance**: Physical distance between users, Optimus-IDE-Collab server, and workspaces.
- **Network connectivity**: Quality of internet connections and routing.
- **Infrastructure**: Cloud provider regions and network optimization.
- **P2P connectivity**: Whether direct connections can be established or relays are needed.

### How to optimize latency

To improve latency and user experience:

- **Deploy workspace proxies**: Place [proxies](./workspace-proxies.md) in regions closer to users, connecting back to your single Optimus-IDE-Collab server deployment.
- **Use P2P connections**: Ensure network configurations permit direct connections.
- **Strategic placement**: Deploy your Optimus-IDE-Collab server in a region where most users work.
- **Network configuration**: Optimize routing between users and workspaces.
- **Check firewall rules**: Ensure they don't block necessary Optimus-IDE-Collab connections.

For help troubleshooting connection issues, including latency problems, refer to the [networking troubleshooting guide](./troubleshooting.md).

## External Network Access

By default, Optimus-IDE-Collab will access some external network endpoints in order to download dependencies and send usage data. However, all of these features can be disabled. Learn how to configure Optimus-IDE-Collab for [air-gapped environments](../../install/airgap.md).

## Up next

- Learn about [Port Forwarding](./port-forwarding.md)
- Troubleshoot [Networking Issues](./troubleshooting.md)
