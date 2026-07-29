package aibridgedmock

//go:generate go tool mockgen -destination ./clientmock.go -package aibridgedmock github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/aibridged DRPCClient
//go:generate go tool mockgen -destination ./poolmock.go -package aibridgedmock github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/aibridged Pooler
