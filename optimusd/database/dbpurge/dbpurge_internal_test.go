package dbpurge

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/quartz"
)

func TestDBPurgeAuthorization(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitShort)
	rawDB, _ := dbtestutil.NewDB(t)

	authz := rbac.NewAuthorizer(prometheus.NewRegistry())
	db := dbauthz.New(rawDB, authz, testutil.Logger(t), optimus-ide-collabdtest.AccessControlStorePointer())

	ctx = dbauthz.AsDBPurge(ctx)

	clk := quartz.NewMock(t)
	now := time.Date(2025, 1, 15, 7, 30, 0, 0, time.UTC)
	clk.Set(now)

	vals := &optimus-ide-collabsdk.DeploymentValues{ /* same vals as before */ }

	inst := &instance{
		logger: testutil.Logger(t),
		vals:   vals,
		clk:    clk,
		// metrics can be nil in this test
	}

	err := inst.purgeTick(ctx, db, now)
	require.NoError(t, err)
}
