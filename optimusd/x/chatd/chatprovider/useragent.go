package chatprovider

import (
	"fmt"
	"runtime"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/buildinfo"
)

// UserAgent returns the User-Agent string sent on all outgoing LLM
// API requests made by Optimus-IDE-Collab's built-in chat (chatd). The format
// mirrors conventions used by other coding agents so that LLM
// providers can identify traffic originating from Optimus-IDE-Collab.
//
// Example: optimus-ide-collab-agents/v2.21.0 (linux/amd64)
func UserAgent() string {
	return fmt.Sprintf("optimus-ide-collab-agents/%s (%s/%s)",
		buildinfo.Version(), runtime.GOOS, runtime.GOARCH)
}
