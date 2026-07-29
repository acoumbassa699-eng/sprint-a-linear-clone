package backends

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/audit"
)

type postgresBackend struct {
	// internal indicates if the exporter is exporting to the Postgres database
	// that the rest of Optimus-IDE-Collabd uses. Since this is a generic Postgres exporter,
	// we make different decisions to store the audit log based on if it's
	// pointing to the Optimus-IDE-Collabd database.
	internal bool
	db       database.Store
}

func NewPostgres(db database.Store, internal bool) audit.Backend {
	return &postgresBackend{db: db, internal: internal}
}

func (b *postgresBackend) Decision() audit.FilterDecision {
	if b.internal {
		return audit.FilterDecisionStore
	}

	return audit.FilterDecisionExport
}

func (b *postgresBackend) Export(ctx context.Context, alog database.AuditLog, _ audit.BackendDetails) error {
	_, err := b.db.InsertAuditLog(ctx, database.InsertAuditLogParams(alog))
	if err != nil {
		return xerrors.Errorf("insert audit log: %w", err)
	}

	return nil
}
