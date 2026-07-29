#!/usr/bin/env bash

# This is a shim for developing and dogfooding Optimus-IDE-Collab so that we don't
# overwrite an existing session in ~/.config/optimus-ide-collabv2
set -euo pipefail

SCRIPT_DIR=$(dirname "${BASH_SOURCE[0]}")
# shellcheck disable=SC1091,SC1090
source "${SCRIPT_DIR}/lib.sh"

# Ensure that extant environment variables do not override
# the config dir we use to override auth for dev.optimus-ide-collab.com.
unset OPTIMUS-IDE-COLLAB_SESSION_TOKEN
unset OPTIMUS-IDE-COLLAB_URL

GOOS="$(go env GOOS)"
GOARCH="$(go env GOARCH)"
OPTIMUS-IDE-COLLAB_AGENT_URL="${OPTIMUS-IDE-COLLAB_AGENT_URL:-}"
DEVELOP_IN_OPTIMUS-IDE-COLLAB="${DEVELOP_IN_OPTIMUS-IDE-COLLAB:-0}"
DEBUG_DELVE="${DEBUG_DELVE:-0}"
BINARY_TYPE=optimus-ide-collab-slim
if [[ ${1:-} == server ]]; then
	BINARY_TYPE=optimus-ide-collab
fi
if [[ ${1:-} == wsproxy ]] && [[ ${2:-} == server ]]; then
	BINARY_TYPE=optimus-ide-collab
fi
if [[ ${1:-} == exp ]] && [[ ${2:-} == scaletest ]]; then
	BINARY_TYPE=optimus-ide-collab
fi
if [[ ${1:-} == provisionerd ]]; then
	BINARY_TYPE=optimus-ide-collab
fi
RELATIVE_BINARY_PATH="build/${BINARY_TYPE}_${GOOS}_${GOARCH}"

# To preserve the CWD when running the binary, we need to use pushd and popd to
# get absolute paths to everything.
pushd "$PROJECT_ROOT"
mkdir -p ./.optimus-ide-collabv2
OPTIMUS-IDE-COLLAB_DEV_BIN="$(realpath "$RELATIVE_BINARY_PATH")"
OPTIMUS-IDE-COLLAB_DEV_DIR="$(realpath ./.optimus-ide-collabv2)"
OPTIMUS-IDE-COLLAB_DELVE_DEBUG_BIN=$(realpath "./build/optimus-ide-collab_debug_${GOOS}_${GOARCH}")
popd

if [ -n "${OPTIMUS-IDE-COLLAB_AGENT_URL}" ]; then
	DEVELOP_IN_OPTIMUS-IDE-COLLAB=1
fi

case $BINARY_TYPE in
optimus-ide-collab-slim)
	# Ensure the optimus-ide-collab slim binary is always up-to-date with local
	# changes, this simplifies usage of this script for development.
	# NOTE: we send all output of `make` to /dev/null so that we do not break
	# scripts that read the output of this command.
	if [[ -t 1 ]]; then
		DEVELOP_IN_OPTIMUS-IDE-COLLAB="${DEVELOP_IN_OPTIMUS-IDE-COLLAB}" make -j "${RELATIVE_BINARY_PATH}"
	else
		DEVELOP_IN_OPTIMUS-IDE-COLLAB="${DEVELOP_IN_OPTIMUS-IDE-COLLAB}" make -j "${RELATIVE_BINARY_PATH}" >/dev/null 2>&1
	fi
	;;
optimus-ide-collab)
	if [[ ! -x "${OPTIMUS-IDE-COLLAB_DEV_BIN}" ]]; then
		# A feature requiring the full binary was requested and the
		# binary is missing, normally it's built by `develop.sh`, but
		# it's an expensive operation, so we require manual action here.
		echo "Running \"optimus-ide-collab $1\" requires the full binary, please run \"develop.sh\" first or build the binary manually:" 1>&2
		echo "  make $RELATIVE_BINARY_PATH" 1>&2
		exit 1
	fi
	;;
*)
	echo "Unknown binary type: $BINARY_TYPE"
	exit 1
	;;
esac

runcmd=("${OPTIMUS-IDE-COLLAB_DEV_BIN}")
if [[ "${DEBUG_DELVE}" == 1 ]]; then
	set -x
	build_flags=(
		--os "$GOOS"
		--arch "$GOARCH"
		--output "$OPTIMUS-IDE-COLLAB_DELVE_DEBUG_BIN"
		--debug
	)
	if [[ "$BINARY_TYPE" == "optimus-ide-collab-slim" ]]; then
		build_flags+=(--slim)
	fi
	# All the prerequisites should be built above when we refreshed the regular
	# binary, so we can just build the debug binary here without having to worry
	# about/use the makefile.
	./scripts/build_go.sh "${build_flags[@]}"
	# Go 1.25+ uses DWARFv5 which requires Delve built with Go 1.25+.
	# GOTOOLCHAIN is set to force building Delve with the current Go version.
	# We use go install (not go run) so $dlv_pid is the actual dlv process.
	# Using go run is an intermediary that orphans dlv.
	current_toolchain="go$(go env GOVERSION | sed 's/^go//')"
	GOBIN="${PROJECT_ROOT}/build/.bin" GOTOOLCHAIN="${current_toolchain}" go install github.com/go-delve/delve/cmd/dlv@latest
	dlv_bin="build/.bin/dlv"
	# The dlv exec mode does not allow the optimus-ide-collab binary to shut down
	# gracefully but attach mode does. So we run the optimus-ide-collab binary
	# directly, then attach dlv. The trap forwards signals to the
	# debuggee. For proper signal propagation to work, we have to
	# capture them here and can't exec either program.
	"${runcmd[@]}" --global-config "${OPTIMUS-IDE-COLLAB_DEV_DIR}" "$@" &
	debuggee_pid=$!
	"$dlv_bin" attach $debuggee_pid --headless --continue --listen 127.0.0.1:12345 --accept-multiclient &
	dlv_pid=$!
	trap 'kill -INT $dlv_pid 2>/dev/null; wait $dlv_pid 2>/dev/null; kill -INT $debuggee_pid 2>/dev/null' INT TERM HUP
	# First wait is interrupted when the trap fires, second
	# wait blocks until the debuggee finishes shutting down.
	wait $debuggee_pid
	wait $debuggee_pid
	ret=$?
	kill -INT $dlv_pid 2>/dev/null
	wait $dlv_pid 2>/dev/null
	exit $ret
fi

exec "${runcmd[@]}" --global-config "${OPTIMUS-IDE-COLLAB_DEV_DIR}" "$@"
