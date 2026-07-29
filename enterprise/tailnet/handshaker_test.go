package tailnet_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/tailnet"
	agpltest "github.com/optimus-ide-collab/optimus-ide-collab/v2/tailnet/test"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestPGCoordinator_ReadyForHandshake_OK(t *testing.T) {
	t.Parallel()

	store, ps := dbtestutil.NewDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitSuperLong)
	defer cancel()
	logger := testutil.Logger(t)
	coord1, err := tailnet.NewPGCoord(ctx, logger.Named("coord1"), ps, store)
	require.NoError(t, err)
	defer coord1.Close()

	agpltest.ReadyForHandshakeTest(ctx, t, coord1)
}

func TestPGCoordinator_ReadyForHandshake_NoPermission(t *testing.T) {
	t.Parallel()

	store, ps := dbtestutil.NewDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitSuperLong)
	defer cancel()
	logger := testutil.Logger(t)
	coord1, err := tailnet.NewPGCoord(ctx, logger.Named("coord1"), ps, store)
	require.NoError(t, err)
	defer coord1.Close()

	agpltest.ReadyForHandshakeNoPermissionTest(ctx, t, coord1)
}
