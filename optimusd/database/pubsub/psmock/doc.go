// package psmock contains a mocked implementation of the pubsub.Pubsub interface for use in tests
package psmock

//go:generate go tool mockgen -destination ./psmock.go -package psmock github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/pubsub Pubsub
