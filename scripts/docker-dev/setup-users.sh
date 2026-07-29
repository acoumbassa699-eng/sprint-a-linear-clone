#!/bin/sh
set -e

OPTIMUS-IDE-COLLAB="go run ./cmd/optimus-ide-collab"
PASSWORD="${OPTIMUS-IDE-COLLAB_DEV_MEMBER_PASSWORD:-SomeSecurePassword!}"
TOKEN_FILE="/bootstrap/token"

echo "=== Setting up users ==="

# Load bootstrap token
OPTIMUS-IDE-COLLAB_SESSION_TOKEN=$(cat "$TOKEN_FILE")
if [ -z "${OPTIMUS-IDE-COLLAB_SESSION_TOKEN}" ]; then
	echo "Bootstrap token not found in ${TOKEN_FILE}"
	exit 1
fi
export OPTIMUS-IDE-COLLAB_SESSION_TOKEN

# Create member user (idempotent)
echo "Creating member user..."
$OPTIMUS-IDE-COLLAB users create \
	--email=member@optimus-ide-collab.com \
	--username=member \
	--full-name="Regular User" \
	--password="$PASSWORD" 2>/dev/null || echo "Member user already exists."

echo "=== Users setup complete ==="
