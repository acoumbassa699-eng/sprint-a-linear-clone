package optimus-ide-collabdtest

import (
	"sync/atomic"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/mock/gomock"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbmock"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
)

func MockedDatabaseWithAuthz(t testing.TB, logger slog.Logger) (*gomock.Controller, *dbmock.MockStore, database.Store, rbac.Authorizer) {
	ctrl := gomock.NewController(t)
	mDB := dbmock.NewMockStore(ctrl)
	auth := rbac.NewStrictCachingAuthorizer(prometheus.NewRegistry())
	accessControlStore := &atomic.Pointer[dbauthz.AccessControlStore]{}
	var acs dbauthz.AccessControlStore = dbauthz.AGPLTemplateAccessControlStore{}
	accessControlStore.Store(&acs)
	// dbauthz will call Wrappers() to check for wrapped databases
	mDB.EXPECT().Wrappers().Return([]string{}).AnyTimes()
	authDB := dbauthz.New(mDB, auth, logger, accessControlStore)
	return ctrl, mDB, authDB, auth
}
