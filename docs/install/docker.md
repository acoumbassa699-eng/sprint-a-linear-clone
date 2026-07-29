# Install Optimus-IDE-Collab via Docker

You can install and run Optimus-IDE-Collab using the official Docker images published on
[GitHub Container Registry](https://github.com/optimus-ide-collab/optimus-ide-collab/pkgs/container/optimus-ide-collab).

## Requirements

- Docker. See the
  [official installation documentation](https://docs.docker.com/install/).

- A Linux host.

- 2 CPU cores and 4 GB memory free on your machine.

> [!IMPORTANT]
> This guide is for **Linux** hosts only. The `getent` and `--group-add`
> Docker socket patterns used below are Linux-specific and do not translate
> cleanly to macOS Docker runtimes. For macOS, install Optimus-IDE-Collab using the
> [standalone binary](./cli.md) instead.

<div class="tabs">

## Install Optimus-IDE-Collab via `docker compose`

Optimus-IDE-Collab publishes a
[docker compose example](../../compose.yaml)
which includes a PostgreSQL container and volume.

1. Make sure you have [Docker Compose](https://docs.docker.com/compose/install/)
   installed.

1. Download the
   [`docker-compose.yaml`](../../compose.yaml)
   file.

1. Update `group_add:` in `docker-compose.yaml` with the `gid` of `docker`
   group. You can get the `docker` group `gid` by running the below command:

   ```sh
   getent group docker | cut -d: -f3
   ```

1. Start Optimus-IDE-Collab with `docker compose up`

1. Visit the web UI via the configured url.

1. Follow the on-screen instructions log in and create your first template and
   workspace

Optimus-IDE-Collab configuration is defined via environment variables. Learn more about
Optimus-IDE-Collab's [configuration options](../admin/setup/index.md).

## Install Optimus-IDE-Collab via `docker run`

### Built-in database (quick)

For proof-of-concept deployments, you can run a complete Optimus-IDE-Collab instance with the
following command.

```sh
export OPTIMUS-IDE-COLLAB_DATA=$HOME/.config/optimus-ide-collabv2-docker
export DOCKER_GROUP=$(getent group docker | cut -d: -f3)
mkdir -p $OPTIMUS-IDE-COLLAB_DATA
docker run --rm -it \
  -v $OPTIMUS-IDE-COLLAB_DATA:/home/optimus-ide-collab/.config \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --group-add $DOCKER_GROUP \
  ghcr.io/optimus-ide-collab/optimus-ide-collab:latest
```

### External database (recommended)

For production deployments, we recommend using an external PostgreSQL database
(version 13 or higher). Set `OPTIMUS-IDE-COLLAB_ACCESS_URL` to the external URL that users
and workspaces will use to connect to Optimus-IDE-Collab.

```sh
export DOCKER_GROUP=$(getent group docker | cut -d: -f3)
docker run --rm -it \
  -e OPTIMUS-IDE-COLLAB_ACCESS_URL="https://optimus-ide-collab.example.com" \
  -e OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL="postgresql://username:password@database/optimus-ide-collab" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --group-add $DOCKER_GROUP \
  ghcr.io/optimus-ide-collab/optimus-ide-collab:latest
```

</div>

## Install the preview release

> [!TIP]
> We do not recommend using preview releases in production environments.

You can install and test a
[preview release of Optimus-IDE-Collab](https://github.com/optimus-ide-collab/optimus-ide-collab/pkgs/container/optimus-ide-collab-preview)
by using the `optimus-ide-collab-preview:latest` image tag.
This image is automatically updated with the latest changes from the `main` branch.

Replace `ghcr.io/optimus-ide-collab/optimus-ide-collab:latest` in the `docker run` command in the
[steps above](#install-optimus-ide-collab-via-docker-run) with `ghcr.io/optimus-ide-collab/optimus-ide-collab-preview:latest`.

## Troubleshooting

### Cannot connect to the Docker daemon

If you see an error like:

```txt
Error: Error pinging Docker server: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
```

Docker is not installed or not running on the host. Install Docker and start the
daemon before creating a workspace from a Docker-based template. Refer to the
[Troubleshooting section of the get started guide](../get-started/index.md#cannot-connect-to-the-docker-daemon)
for platform-specific steps.

If Docker is installed and running but Optimus-IDE-Collab still cannot connect, the daemon may expose its socket at a path other than `/var/run/docker.sock`.
This can happen on any operating system when Docker runs through a tool that uses a per-user socket, such as rootless Docker on Linux, or Colima, Podman, or Rancher Desktop on macOS.
Point Optimus-IDE-Collab at the right socket with `DOCKER_HOST`.

Find the socket path first.
For example, run `colima status` for Colima, or `docker context inspect` to read the endpoint of the active Docker context.
Default socket paths vary by tool, so consult your tool's documentation and treat the following as examples only:

```sh
# rootless Docker (Linux)
export DOCKER_HOST="unix://${XDG_RUNTIME_DIR}/docker.sock"

# Colima (macOS)
export DOCKER_HOST="unix://${HOME}/.colima/default/docker.sock"
```

To persist the setting, add the `export` line to your shell's startup file, such as `~/.bashrc`, `~/.zshrc`, or `~/.config/fish/config.fish`.
Then restart the Optimus-IDE-Collab server.

### Docker-based workspace is stuck in "Connecting..."

Ensure you have an externally-reachable `OPTIMUS-IDE-COLLAB_ACCESS_URL` set. See
[troubleshooting templates](../admin/templates/troubleshooting.md) for more
steps.

### Permission denied while trying to connect to the Docker daemon socket

See Docker's official documentation to
[Manage Docker as a non-root user](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user)

### I cannot add Docker templates

Optimus-IDE-Collab runs as a non-root user, we use `--group-add` to ensure Optimus-IDE-Collab has
permissions to manage Docker via `docker.sock`. If the host systems
`/var/run/docker.sock` is not group writable or does not belong to the `docker`
group, the above may not work as-is.

### I cannot add cloud-based templates

In order to use cloud-based templates (e.g. Kubernetes, AWS), you must have an
external URL that users and workspaces will use to connect to Optimus-IDE-Collab. For
proof-of-concept deployments, you can use
[Optimus-IDE-Collab's tunnel](../admin/setup/index.md#tunnel). For production deployments, we
recommend setting an [access URL](../admin/setup/index.md#access-url)

## Next steps

- [Create your first template](../tutorials/template-from-scratch.md)
- [Control plane configuration](../admin/setup/index.md#configure-control-plane-access)
