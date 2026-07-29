package optimus-ide-collabd

import "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd"

// ChatDaemonForTest returns the background chat processor for test harnesses.
func (api *API) ChatDaemonForTest() *chatd.Server {
	return api.chatDaemon
}
