package rbac_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
)

func TestIsUnauthorizedError(t *testing.T) {
	t.Parallel()
	t.Run("NotWrapped", func(t *testing.T) {
		t.Parallel()
		errFunc := func() error {
			return rbac.UnauthorizedError{}
		}

		err := errFunc()
		require.True(t, rbac.IsUnauthorizedError(err))
	})

	t.Run("Wrapped", func(t *testing.T) {
		t.Parallel()
		errFunc := func() error {
			return xerrors.Errorf("test error: %w", rbac.UnauthorizedError{})
		}
		err := errFunc()
		require.True(t, rbac.IsUnauthorizedError(err))
	})
}
