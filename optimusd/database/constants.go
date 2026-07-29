package database

import (
	"github.com/google/uuid"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// PrebuildsSystemUserID mirrors optimus-ide-collabsdk.PrebuildsSystemUserID, parsed
// for use as a uuid.UUID. Both must agree; tests pin the value to the
// optimus-ide-collabsdk constant so the two cannot drift.
var PrebuildsSystemUserID = uuid.MustParse(optimus-ide-collabsdk.PrebuildsSystemUserID)
