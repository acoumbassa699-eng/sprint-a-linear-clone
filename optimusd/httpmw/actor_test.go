package httpmw_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestRequireAPIKeyOrWorkspaceProxyAuth(t *testing.T) {
	t.Parallel()

	t.Run("None", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		rw := httptest.NewRecorder()

		httpmw.RequireAPIKeyOrWorkspaceProxyAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("should not have been called")
		})).ServeHTTP(rw, r)

		require.Equal(t, http.StatusUnauthorized, rw.Code)
	})

	t.Run("APIKey", func(t *testing.T) {
		t.Parallel()

		var (
			db, _    = dbtestutil.NewDB(t)
			user     = dbgen.User(t, db, database.User{})
			_, token = dbgen.APIKey(t, db, database.APIKey{
				UserID:    user.ID,
				ExpiresAt: dbtime.Now().AddDate(0, 0, 1),
			})

			r  = httptest.NewRequest("GET", "/", nil)
			rw = httptest.NewRecorder()
		)
		r.Header.Set(optimus-ide-collabsdk.SessionTokenHeader, token)

		var called atomic.Int64
		httpmw.ExtractAPIKeyMW(httpmw.ExtractAPIKeyConfig{
			DB:              db,
			RedirectToLogin: false,
		})(
			httpmw.RequireAPIKeyOrWorkspaceProxyAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called.Add(1)
				rw.WriteHeader(http.StatusOK)
			}))).
			ServeHTTP(rw, r)

		res := rw.Result()
		defer res.Body.Close()
		dump, err := httputil.DumpResponse(res, true)
		require.NoError(t, err)
		t.Log(string(dump))

		require.Equal(t, http.StatusOK, rw.Code)
		require.Equal(t, int64(1), called.Load())
	})

	t.Run("WorkspaceProxy", func(t *testing.T) {
		t.Parallel()

		var (
			db, _        = dbtestutil.NewDB(t)
			user         = dbgen.User(t, db, database.User{})
			_, userToken = dbgen.APIKey(t, db, database.APIKey{
				UserID:    user.ID,
				ExpiresAt: dbtime.Now().AddDate(0, 0, 1),
			})
			proxy, proxyToken = dbgen.WorkspaceProxy(t, db, database.WorkspaceProxy{})

			r  = httptest.NewRequest("GET", "/", nil)
			rw = httptest.NewRecorder()
		)
		r.Header.Set(optimus-ide-collabsdk.SessionTokenHeader, userToken)
		r.Header.Set(httpmw.WorkspaceProxyAuthTokenHeader, fmt.Sprintf("%s:%s", proxy.ID, proxyToken))

		httpmw.ExtractAPIKeyMW(httpmw.ExtractAPIKeyConfig{
			DB:              db,
			RedirectToLogin: false,
		})(
			httpmw.ExtractWorkspaceProxy(httpmw.ExtractWorkspaceProxyConfig{
				DB: db,
			})(
				httpmw.RequireAPIKeyOrWorkspaceProxyAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rw.WriteHeader(http.StatusOK)
				})))).
			ServeHTTP(rw, r)

		res := rw.Result()
		defer res.Body.Close()
		dump, err := httputil.DumpResponse(res, true)
		require.NoError(t, err)
		t.Log(string(dump))

		require.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run("Both", func(t *testing.T) {
		t.Parallel()

		var (
			db, _        = dbtestutil.NewDB(t)
			proxy, token = dbgen.WorkspaceProxy(t, db, database.WorkspaceProxy{})

			r  = httptest.NewRequest("GET", "/", nil)
			rw = httptest.NewRecorder()
		)
		r.Header.Set(httpmw.WorkspaceProxyAuthTokenHeader, fmt.Sprintf("%s:%s", proxy.ID, token))

		var called atomic.Int64
		httpmw.ExtractWorkspaceProxy(httpmw.ExtractWorkspaceProxyConfig{
			DB: db,
		})(
			httpmw.RequireAPIKeyOrWorkspaceProxyAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called.Add(1)
				rw.WriteHeader(http.StatusOK)
			}))).
			ServeHTTP(rw, r)

		res := rw.Result()
		defer res.Body.Close()
		dump, err := httputil.DumpResponse(res, true)
		require.NoError(t, err)
		t.Log(string(dump))

		require.Equal(t, http.StatusOK, rw.Code)
		require.Equal(t, int64(1), called.Load())
	})
}
