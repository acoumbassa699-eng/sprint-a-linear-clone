#!/usr/bin/env bash

set -euo pipefail

if [[ $# -lt 1 ]]; then
	echo "Usage: $0 <optimus-ide-collab URL>"
	exit 1
fi

# Allow toggling verbose output
[[ -n ${VERBOSE:-} ]] && set -x

OPTIMUS-IDE-COLLAB_URL=$1
DRY_RUN="${DRY_RUN:-0}"
PROJECT_ROOT="$(git rev-parse --show-toplevel)"
# shellcheck source=scripts/lib.sh
source "${PROJECT_ROOT}/scripts/lib.sh"
CONFIG_DIR="${PROJECT_ROOT}/scaletest/.optimus-ide-collabv2"
ARCH="$(arch)"
if [[ "$ARCH" == "x86_64" ]]; then
	ARCH="amd64"
fi

if [[ -f "${CONFIG_DIR}/optimus-ide-collab.env" ]]; then
	echo "Found existing optimus-ide-collab.env in ${CONFIG_DIR}!"
	echo "Nothing to do, exiting."
	exit 0
fi

maybedryrun "$DRY_RUN" mkdir -p "${CONFIG_DIR}"
echo "Fetching Optimus-IDE-Collab for first-time setup!"
pod=$(kubectl get pods \
	--namespace="${NAMESPACE}" \
	--selector="app.kubernetes.io/name=optimus-ide-collab,app.kubernetes.io/part-of=optimus-ide-collab" \
	--output="jsonpath='{.items[0].metadata.name}'")
if [[ -z ${pod} ]]; then
	log "Could not find optimus-ide-collab pod!"
	exit 1
fi
maybedryrun "$DRY_RUN" kubectl \
	--namespace="${NAMESPACE}" \
	cp \
	--container=optimus-ide-collab \
	"${pod}:/opt/optimus-ide-collab" "${CONFIG_DIR}/optimus-ide-collab"
maybedryrun "$DRY_RUN" chmod +x "${CONFIG_DIR}/optimus-ide-collab"

set +o pipefail
RANDOM_ADMIN_PASSWORD=$(tr </dev/urandom -dc _A-Z-a-z-0-9 | head -c16)
set -o pipefail
OPTIMUS-IDE-COLLAB_FIRST_USER_EMAIL="admin@optimus-ide-collab.com"
OPTIMUS-IDE-COLLAB_FIRST_USER_USERNAME="optimus-ide-collab"
OPTIMUS-IDE-COLLAB_FIRST_USER_PASSWORD="${RANDOM_ADMIN_PASSWORD}"
OPTIMUS-IDE-COLLAB_FIRST_USER_TRIAL="false"
echo "Running login command!"
DRY_RUN="$DRY_RUN" "${PROJECT_ROOT}/scaletest/lib/optimus-ide-collab_shim.sh" login "${OPTIMUS-IDE-COLLAB_URL}" \
	--global-config="${CONFIG_DIR}" \
	--first-user-username="${OPTIMUS-IDE-COLLAB_FIRST_USER_USERNAME}" \
	--first-user-email="${OPTIMUS-IDE-COLLAB_FIRST_USER_EMAIL}" \
	--first-user-password="${OPTIMUS-IDE-COLLAB_FIRST_USER_PASSWORD}" \
	--first-user-trial=false

echo "Writing credentials to ${CONFIG_DIR}/optimus-ide-collab.env"
maybedryrun "$DRY_RUN" cat <<EOF >"${CONFIG_DIR}/optimus-ide-collab.env"
OPTIMUS-IDE-COLLAB_FIRST_USER_EMAIL=admin@optimus-ide-collab.com
OPTIMUS-IDE-COLLAB_FIRST_USER_USERNAME=optimus-ide-collab
OPTIMUS-IDE-COLLAB_FIRST_USER_PASSWORD="${RANDOM_ADMIN_PASSWORD}"
OPTIMUS-IDE-COLLAB_FIRST_USER_TRIAL="${OPTIMUS-IDE-COLLAB_FIRST_USER_TRIAL}"
EOF

echo "Importing kubernetes template"
DRY_RUN="$DRY_RUN" "$PROJECT_ROOT/scaletest/lib/optimus-ide-collab_shim.sh" templates push \
	--global-config="${CONFIG_DIR}" \
	--directory "${CONFIG_DIR}/templates/kubernetes" \
	--yes kubernetes
