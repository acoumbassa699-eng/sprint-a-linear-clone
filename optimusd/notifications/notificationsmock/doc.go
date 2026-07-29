// Package notificationsmock contains a mocked implementation of the
// notifications.Enqueuer interface for use in tests.
package notificationsmock

//go:generate go tool mockgen -destination ./notificationsmock.go -package notificationsmock github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications Enqueuer
