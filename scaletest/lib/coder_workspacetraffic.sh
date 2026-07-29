#!/usr/bin/env bash

set -euo pipefail

PROJECT_ROOT="$(git rev-parse --show-toplevel)"
# shellcheck source=scripts/lib.sh
source "${PROJECT_ROOT}/scripts/lib.sh"

# Allow toggling verbose output
[[ -n ${VERBOSE:-} ]] && set -x

SCALETEST_NAME="${SCALETEST_NAME:-}"
SCALETEST_TRAFFIC_BYTES_PER_TICK="${SCALETEST_TRAFFIC_BYTES_PER_TICK:-1024}"
SCALETEST_TRAFFIC_TICK_INTERVAL="${SCALETEST_TRAFFIC_TICK_INTERVAL:-100ms}"

script_name=$(basename "$0")
args="$(getopt -o "" -l help,name:,traffic-bytes-per-tick:,traffic-tick-interval:, -- "$@")"
eval set -- "$args"
while true; do
	case "$1" in
	--help)
		echo "Usage: $script_name --name <name> [--traffic-bytes-per-tick <bytes_per-tick>] [--traffic-tick-interval <ticks_per_second]"
		exit 1
		;;
	--name)
		SCALETEST_NAME="$2"
		shift 2
		;;
	--traffic-bytes-per-tick)
		SCALETEST_TRAFFIC_BYTES_PER_TICK="$2"
		shift 2
		;;
	--traffic-tick-interval)
		SCALETEST_TRAFFIC_TICK_INTERVAL="$2"
		shift 2
		;;
	--)
		shift
		break
		;;
	*)
		error "Unrecognized option: $1"
		;;
	esac
done

dependencies kubectl

if [[ -z "${SCALETEST_NAME}" ]]; then
	echo "Must specify --name"
	exit 1
fi

OPTIMUS-IDE-COLLAB_TOKEN=$("${PROJECT_ROOT}/scaletest/lib/optimus-ide-collab_shim.sh" tokens create)
OPTIMUS-IDE-COLLAB_URL="http://optimus-ide-collab.optimus-ide-collab-${SCALETEST_NAME}.svc.cluster.local"
export KUBECONFIG="${PROJECT_ROOT}/scaletest/.optimus-ide-collabv2/${SCALETEST_NAME}-cluster.kubeconfig"

# Clean up any pre-existing pods
kubectl -n "optimus-ide-collab-${SCALETEST_NAME}" delete pod optimus-ide-collab-scaletest-workspace-traffic --force || true

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: optimus-ide-collab-scaletest-workspace-traffic
  namespace: optimus-ide-collab-${SCALETEST_NAME}
  labels:
    app.kubernetes.io/name: optimus-ide-collab-scaletest-workspace-traffic
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: cloud.google.com/gke-nodepool
            operator: In
            values:
            - ${SCALETEST_NAME}-misc
  containers:
  - command:
    - sh
    - -c
    - "curl -fsSL $OPTIMUS-IDE-COLLAB_URL/bin/optimus-ide-collab-linux-amd64 -o /tmp/optimus-ide-collab && chmod +x /tmp/optimus-ide-collab && /tmp/optimus-ide-collab --verbose --url=$OPTIMUS-IDE-COLLAB_URL --token=$OPTIMUS-IDE-COLLAB_TOKEN exp scaletest workspace-traffic --concurrency=0 --bytes-per-tick=${SCALETEST_TRAFFIC_BYTES_PER_TICK} --tick-interval=${SCALETEST_TRAFFIC_TICK_INTERVAL} --scaletest-prometheus-wait=60s"
    env:
    - name: OPTIMUS-IDE-COLLAB_URL
      value: $OPTIMUS-IDE-COLLAB_URL
    - name: OPTIMUS-IDE-COLLAB_TOKEN
      value: $OPTIMUS-IDE-COLLAB_TOKEN
    - name: OPTIMUS-IDE-COLLAB_SCALETEST_PROMETHEUS_ADDRESS
      value: "0.0.0.0:21112"
    - name: OPTIMUS-IDE-COLLAB_SCALETEST_JOB_TIMEOUT
      value: "30m"
    ports:
    - containerPort: 21112
      name: prometheus-http
      protocol: TCP
    name: cli
    image: docker.io/optimus-ide-collabcom/enterprise-minimal:ubuntu
  restartPolicy: Never
---
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  namespace: optimus-ide-collab-${SCALETEST_NAME}
  name: optimus-ide-collab-workspacetraffic-monitoring
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: optimus-ide-collab-scaletest-workspace-traffic
  podMetricsEndpoints:
  - port: prometheus-http
    interval: 15s
EOF
