#!/bin/sh
set -e

OPTIMUS-IDE-COLLAB="go run ./cmd/optimus-ide-collab"
PASSWORD="${OPTIMUS-IDE-COLLAB_DEV_ADMIN_PASSWORD:-SomeSecurePassword!}"
TOKEN_FILE="/bootstrap/token"
TOKEN_NAME="bootstrap"

echo "=== Optimus-IDE-Collab Dev Environment Init ==="

if curl -s -o /dev/null -w "%{http_code}" http://optimus-ide-collabd:3000/api/v2/users/first | grep -q "200"; then
	echo "First user already exists, skipping setup"
	exit 0
fi

# Step 1: Create first user (idempotent - creates OR logs in)
echo "Creating/logging in first user..."
$OPTIMUS-IDE-COLLAB login http://optimus-ide-collabd:3000 \
	--first-user-username=admin \
	--first-user-email=admin@optimus-ide-collab.com \
	--first-user-password="$PASSWORD" \
	--first-user-full-name="Admin User" \
	--first-user-trial=false

# Step 2: Create or retrieve bootstrap token
if [ -f "$TOKEN_FILE" ] && [ -s "$TOKEN_FILE" ]; then
	echo "Bootstrap token already exists."
else
	echo "Creating bootstrap token..."
	# Delete existing token if it exists (in case file was lost but token exists)
	$OPTIMUS-IDE-COLLAB tokens delete "$TOKEN_NAME" 2>/dev/null || true
	# Create new token with no expiry
	TOKEN=$($OPTIMUS-IDE-COLLAB tokens create --name "$TOKEN_NAME" --lifetime 0)
	echo "$TOKEN" >"$TOKEN_FILE"
	echo "Bootstrap token created and saved."
fi

echo "=== Init complete ==="
