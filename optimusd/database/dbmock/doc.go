// package dbmock contains a mocked implementation of the database.Store interface for use in tests
package dbmock

//go:generate go tool mockgen -destination ./dbmock.go -package dbmock github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database Store
