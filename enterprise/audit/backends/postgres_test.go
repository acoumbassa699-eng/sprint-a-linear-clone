package backends_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/audit/audittest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/audit/backends"
)

func TestPostgresBackend(t *testing.T) {
	t.Parallel()
	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		var (
			ctx, cancel = context.WithCancel(context.Background())
			db, _       = dbtestutil.NewDB(t)
			pgb         = backends.NewPostgres(db, true)
			alog        = audittest.RandomLog()
		)
		defer cancel()

		err := pgb.Export(ctx, alog, audit.BackendDetails{})
		require.NoError(t, err)

		got, err := db.GetAuditLogsOffset(ctx, database.GetAuditLogsOffsetParams{
			OffsetOpt: 0,
			LimitOpt:  1,
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, alog.ID, got[0].AuditLog.ID)
	})
}
