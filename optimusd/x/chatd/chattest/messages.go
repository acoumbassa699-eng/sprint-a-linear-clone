package chattest

import (
	"encoding/json"

	"github.com/sqlc-dev/pqtype"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// ChatMessageWithParts returns a database chat message whose content is the
// JSON encoding of the provided SDK message parts.
func ChatMessageWithParts(parts []optimus-ide-collabsdk.ChatMessagePart) database.ChatMessage {
	raw, _ := json.Marshal(parts)
	return database.ChatMessage{
		Content: pqtype.NullRawMessage{RawMessage: raw, Valid: true},
	}
}
