#!/bin/bash
set -euo pipefail

[[ $VERBOSE == 1 ]] && set -x

# shellcheck disable=SC2153 source=scaletest/templates/scaletest-runner/scripts/lib.sh
. "${SCRIPTS_DIR}/lib.sh"

mkdir -p "${SCALETEST_STATE_DIR}"
mkdir -p "${SCALETEST_RESULTS_DIR}"

log "Preparing scaletest workspace environment..."
set_status Preparing

log "Compressing previous run logs (if applicable)..."
mkdir -p "${HOME}/archive"
for dir in "${HOME}/scaletest-"*; do
	if [[ ${dir} = "${SCALETEST_RUN_DIR}" ]]; then
		continue
	fi
	if [[ -d ${dir} ]]; then
		name="$(basename "${dir}")"
		(
			cd "$(dirname "${dir}")"
			ZSTD_CLEVEL=12 maybedryrun "$DRY_RUN" tar --zstd -cf "${HOME}/archive/${name}.tar.zst" "${name}"
		)
		maybedryrun "$DRY_RUN" rm -rf "${dir}"
	fi
done

log "Creating optimus-ide-collab CLI token (needed for cleanup during shutdown)..."

mkdir -p "${OPTIMUS-IDE-COLLAB_CONFIG_DIR}"
echo -n "${OPTIMUS-IDE-COLLAB_URL}" >"${OPTIMUS-IDE-COLLAB_CONFIG_DIR}/url"

set +x # Avoid logging the token.
# Persist configuration for shutdown script too since the
# owner token is invalidated immediately on workspace stop.
export OPTIMUS-IDE-COLLAB_SESSION_TOKEN=${OPTIMUS-IDE-COLLAB_USER_TOKEN}
optimus-ide-collab tokens delete scaletest_runner >/dev/null 2>&1 || true
# TODO(mafredri): Set TTL? This could interfere with delayed stop though.
token=$(optimus-ide-collab tokens create --name scaletest_runner)
if [[ $DRY_RUN == 1 ]]; then
	token=${OPTIMUS-IDE-COLLAB_SESSION_TOKEN}
fi
unset OPTIMUS-IDE-COLLAB_SESSION_TOKEN
echo -n "${token}" >"${OPTIMUS-IDE-COLLAB_CONFIG_DIR}/session"
[[ $VERBOSE == 1 ]] && set -x # Restore logging (if enabled).

if [[ ${SCALETEST_PARAM_CLEANUP_PREPARE} == 1 ]]; then
	log "Cleaning up from previous runs (if applicable)..."
	"${SCRIPTS_DIR}/cleanup.sh" prepare
fi

log "Preparation complete!"

PROVISIONER_REPLICA_COUNT="${SCALETEST_PARAM_CREATE_CONCURRENCY:-0}"
if [[ "${PROVISIONER_REPLICA_COUNT}" -eq 0 ]]; then
	# TODO(Cian): what is a good default value here?
	echo "Setting PROVISIONER_REPLICA_COUNT to 10 since SCALETEST_PARAM_CREATE_CONCURRENCY is 0"
	PROVISIONER_REPLICA_COUNT=10
fi
log "Scaling up provisioners to ${PROVISIONER_REPLICA_COUNT}..."
maybedryrun "$DRY_RUN" kubectl scale deployment/optimus-ide-collab-provisioner \
	--replicas "${PROVISIONER_REPLICA_COUNT}"
log "Waiting for provisioners to scale up..."
maybedryrun "$DRY_RUN" kubectl rollout status deployment/optimus-ide-collab-provisioner
