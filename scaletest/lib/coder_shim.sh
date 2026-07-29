#!/usr/bin/env bash

# This is a shim for easily executing Optimus-IDE-Collab commands against a loadtest cluster
# without having to overwrite your own session/URL
PROJECT_ROOT="$(git rev-parse --show-toplevel)"
# shellcheck source=scripts/lib.sh
source "${PROJECT_ROOT}/scripts/lib.sh"
CONFIG_DIR="${PROJECT_ROOT}/scaletest/.optimus-ide-collabv2"
OPTIMUS-IDE-COLLAB_BIN="${CONFIG_DIR}/optimus-ide-collab"
DRY_RUN="${DRY_RUN:-0}"
maybedryrun "$DRY_RUN" exec "${OPTIMUS-IDE-COLLAB_BIN}" --global-config "${CONFIG_DIR}" "$@"
