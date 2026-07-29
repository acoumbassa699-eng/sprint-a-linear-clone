//go:build !linux && !(windows && amd64)

package agent

import (
	"time"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type osListeningPortsGetter struct {
	cacheDuration time.Duration
}

func (*osListeningPortsGetter) GetListeningPorts() ([]optimus-ide-collabsdk.WorkspaceAgentListeningPort, error) {
	// Can't scan for ports on non-linux or non-windows_amd64 systems at the
	// moment. The UI will not show any "no ports found" message to the user, so
	// the user won't suspect a thing.
	return []optimus-ide-collabsdk.WorkspaceAgentListeningPort{}, nil
}
