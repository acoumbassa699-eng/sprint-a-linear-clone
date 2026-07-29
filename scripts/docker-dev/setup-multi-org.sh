#!/bin/sh
set -e

OPTIMUS-IDE-COLLAB="go run ./enterprise/cmd/optimus-ide-collab"
TOKEN_FILE="/bootstrap/token"
LICENSE_FILE="/license.txt"
ORG_NAME="${ORG_NAME:-second-organization}"

echo "=== Multi-Organization Setup ==="

# Load bootstrap token
OPTIMUS-IDE-COLLAB_SESSION_TOKEN=$(cat "$TOKEN_FILE")
if [ -z "${OPTIMUS-IDE-COLLAB_SESSION_TOKEN}" ]; then
	echo "Bootstrap token not found in ${TOKEN_FILE}"
	exit 1
fi
export OPTIMUS-IDE-COLLAB_SESSION_TOKEN

# Check if a license has not yet been added
LICENSES=$($OPTIMUS-IDE-COLLAB license list | tail -n +2)
if [ -z "${LICENSES}" ]; then
	echo "No existing license found."
	if [ ! -f "${LICENSE_FILE}" ]; then
		echo "License required, set OPTIMUS-IDE-COLLAB_DEV_LICENSE_FILE=path/to/license.txt"
		exit 1
	fi
	echo "Adding license..."
	$OPTIMUS-IDE-COLLAB license add --file "${LICENSE_FILE}"
fi

# Create second organization if it doesn't exist.
if ! $OPTIMUS-IDE-COLLAB organizations show "$ORG_NAME" >/dev/null 2>&1; then
	echo "Creating organization '$ORG_NAME'..."
	$OPTIMUS-IDE-COLLAB organizations create -y "$ORG_NAME"
else
	echo "Organization '$ORG_NAME' already exists."
fi

# Add member user to the organization.
echo "Adding member user to organization '$ORG_NAME'..."
$OPTIMUS-IDE-COLLAB organizations members add member --org "$ORG_NAME" 2>/dev/null ||
	echo "Member already in organization or failed to add."

echo "=== Multi-org setup complete ==="
