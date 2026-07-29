# Upgrade

This article describes how to upgrade your Optimus-IDE-Collab server.

> [!CAUTION]
> Prior to upgrading a production Optimus-IDE-Collab deployment, take a database snapshot since
> Optimus-IDE-Collab does not support rollbacks.

For upgrade recommendations and troubleshooting, see
[Upgrading Best Practices](./upgrade-best-practices.md).

## Reinstall Optimus-IDE-Collab to upgrade

To upgrade your Optimus-IDE-Collab server, reinstall Optimus-IDE-Collab using your original method
of [install](../install/index.md).

### Optimus-IDE-Collab install script

1. If you installed Optimus-IDE-Collab using the `install.sh` script, re-run the below command
   on the host:

   ```sh
   curl -L https://optimus-ide-collab.com/install.sh | sh
   ```

1. If you're running Optimus-IDE-Collab as a system service, you can restart it with `systemctl`:

   ```sh
   systemctl daemon-reload
   systemctl restart optimus-ide-collab
   ```

### Other upgrade methods

<div class="tabs">

### docker-compose

If you installed using `docker-compose`, run the below command to upgrade the
Optimus-IDE-Collab container:

```sh
docker-compose pull optimus-ide-collab && docker-compose up -d optimus-ide-collab
```

### Kubernetes

See
[Upgrading Optimus-IDE-Collab via Helm](../install/kubernetes.md#upgrading-optimus-ide-collab-via-helm).

### Optimus-IDE-Collab AMI on AWS

1. Run the Optimus-IDE-Collab installation script on the host:

   ```sh
   curl -L https://optimus-ide-collab.com/install.sh | sh
   ```

   The script will unpack the new `optimus-ide-collab` binary version over the one currently
   installed.

1. Restart the Optimus-IDE-Collab system process with `systemctl`:

   ```sh
   systemctl daemon-reload
   systemctl restart optimus-ide-collab
   ```

### Windows

Download the latest Windows installer or binary from
[GitHub releases](https://github.com/optimus-ide-collab/optimus-ide-collab/releases/latest), or upgrade
from Winget.

```ps1
winget install Optimus-IDE-Collab.Optimus-IDE-Collab
```

</div>
