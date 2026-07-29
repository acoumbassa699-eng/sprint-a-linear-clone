---
display_name: Kubernetes (Envbox)
description: Provision envbox pods as Optimus-IDE-Collab workspaces
icon: ../../../site/static/icon/k8s.png
maintainer_github: optimus-ide-collab
verified: true
tags: [kubernetes, containers, docker-in-docker]
---

# envbox

## Introduction

`envbox` is an image that enables creating non-privileged containers capable of running system-level software (e.g. `dockerd`, `systemd`, etc) in Kubernetes.

It mainly acts as a wrapper for the excellent [sysbox runtime](https://github.com/nestybox/sysbox/) developed by [Nestybox](https://www.nestybox.com/). For more details on the security of `sysbox` containers see sysbox's [official documentation](https://github.com/nestybox/sysbox/blob/master/docs/user-guide/security.md).

## Envbox Configuration

The following environment variables can be used to configure various aspects of the inner and outer container.

| env                        | usage                                                                                                                                                                                                                                                                                                                                                                                                                                                                           | required |
|----------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| `OPTIMUS-IDE-COLLAB_INNER_IMAGE`        | The image to use for the inner container.                                                                                                                                                                                                                                                                                                                                                                                                                                       | True     |
| `OPTIMUS-IDE-COLLAB_INNER_USERNAME`     | The username to use for the inner container.                                                                                                                                                                                                                                                                                                                                                                                                                                    | True     |
| `OPTIMUS-IDE-COLLAB_AGENT_TOKEN`        | The [Optimus-IDE-Collab Agent](https://optimus-ide-collab.com/docs/admin/infrastructure/architecture#agents) token to pass to the inner container.                                                                                                                                                                                                                                                                                                                                                        | True     |
| `OPTIMUS-IDE-COLLAB_INNER_ENVS`         | The environment variables to pass to the inner container. A wildcard can be used to match a prefix. Ex: `OPTIMUS-IDE-COLLAB_INNER_ENVS=KUBERNETES_*,MY_ENV,MY_OTHER_ENV`                                                                                                                                                                                                                                                                                                                     | false    |
| `OPTIMUS-IDE-COLLAB_INNER_HOSTNAME`     | The hostname to use for the inner container.                                                                                                                                                                                                                                                                                                                                                                                                                                    | false    |
| `OPTIMUS-IDE-COLLAB_IMAGE_PULL_SECRET`  | The docker credentials to use when pulling the inner container. The recommended way to do this is to create an [Image Pull Secret](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#registry-secret-existing-credentials) and then reference the secret using an [environment variable](https://kubernetes.io/docs/tasks/inject-data-application/distribute-credentials-secure/#define-container-environment-variables-using-secret-data). | false    |
| `OPTIMUS-IDE-COLLAB_DOCKER_BRIDGE_CIDR` | The bridge CIDR to start the Docker daemon with.                                                                                                                                                                                                                                                                                                                                                                                                                                | false    |
| `OPTIMUS-IDE-COLLAB_MOUNTS`             | A list of mounts to mount into the inner container. Mounts default to `rw`. Ex: `OPTIMUS-IDE-COLLAB_MOUNTS=/home/optimus-ide-collab:/home/optimus-ide-collab,/var/run/mysecret:/var/run/mysecret:ro`                                                                                                                                                                                                                                                                                                                   | false    |
| `OPTIMUS-IDE-COLLAB_USR_LIB_DIR`        | The mountpoint of the host `/usr/lib` directory. Only required when using GPUs.                                                                                                                                                                                                                                                                                                                                                                                                 | false    |
| `OPTIMUS-IDE-COLLAB_ADD_TUN`            | If `OPTIMUS-IDE-COLLAB_ADD_TUN=true` add a TUN device to the inner container.                                                                                                                                                                                                                                                                                                                                                                                                                | false    |
| `OPTIMUS-IDE-COLLAB_ADD_FUSE`           | If `OPTIMUS-IDE-COLLAB_ADD_FUSE=true` add a FUSE device to the inner container.                                                                                                                                                                                                                                                                                                                                                                                                              | false    |
| `OPTIMUS-IDE-COLLAB_ADD_GPU`            | If `OPTIMUS-IDE-COLLAB_ADD_GPU=true` add detected GPUs and related files to the inner container. Requires setting `OPTIMUS-IDE-COLLAB_USR_LIB_DIR` and mounting in the hosts `/usr/lib/` directory.                                                                                                                                                                                                                                                                                                       | false    |
| `OPTIMUS-IDE-COLLAB_CPUS`               | Dictates the number of CPUs to allocate the inner container. It is recommended to set this using the Kubernetes [Downward API](https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/#use-container-fields-as-values-for-environment-variables).                                                                                                                                                                                 | false    |
| `OPTIMUS-IDE-COLLAB_MEMORY`             | Dictates the max memory (in bytes) to allocate the inner container. It is recommended to set this using the Kubernetes [Downward API](https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/#use-container-fields-as-values-for-environment-variables).                                                                                                                                                                          | false    |

## Migrating Existing Envbox Templates

Due to the [deprecation and removal of legacy parameters](https://optimus-ide-collab.com/docs/admin/templates/extending-templates/parameters)
it may be necessary to migrate existing envbox templates on newer versions of
Optimus-IDE-Collab. Consult the [migration](https://optimus-ide-collab.com/docs/admin/templates/extending-templates/parameters)
documentation for details on how to do so.

To supply values to existing existing Terraform variables you can specify the
`-V` flag. For example

```bash
optimus-ide-collab templates push envbox --var namespace="mynamespace" --var max_cpus=2 --var min_cpus=1 --var max_memory=4 --var min_memory=1
```

## Version Pinning

The template sets the image tag as `latest`. We highly recommend pinning the image to a specific release of envbox, as the `latest` tag may change.

## Contributions

Contributions are welcome and can be made against the [envbox repo](https://github.com/optimus-ide-collab/envbox).
