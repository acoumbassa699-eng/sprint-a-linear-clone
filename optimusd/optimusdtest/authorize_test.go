package optimus-ide-collabdtest_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
)

func TestAuthzRecorder(t *testing.T) {
	t.Parallel()

	t.Run("Authorize", func(t *testing.T) {
		t.Parallel()

		rec := &optimus-ide-collabdtest.RecordingAuthorizer{
			Wrapped: &optimus-ide-collabdtest.FakeAuthorizer{},
		}
		sub := optimus-ide-collabdtest.RandomRBACSubject()
		pairs := fuzzAuthz(t, sub, rec, 10)
		rec.AssertActor(t, sub, pairs...)
		require.NoError(t, rec.AllAsserted(), "all assertions should have been made")
	})

	t.Run("Authorize2Subjects", func(t *testing.T) {
		t.Parallel()

		rec := &optimus-ide-collabdtest.RecordingAuthorizer{
			Wrapped: &optimus-ide-collabdtest.FakeAuthorizer{},
		}
		a := optimus-ide-collabdtest.RandomRBACSubject()
		aPairs := fuzzAuthz(t, a, rec, 10)

		b := optimus-ide-collabdtest.RandomRBACSubject()
		bPairs := fuzzAuthz(t, b, rec, 10)

		rec.AssertActor(t, b, bPairs...)
		rec.AssertActor(t, a, aPairs...)
		require.NoError(t, rec.AllAsserted(), "all assertions should have been made")
	})

	t.Run("Authorize_Prepared", func(t *testing.T) {
		t.Parallel()

		rec := &optimus-ide-collabdtest.RecordingAuthorizer{
			Wrapped: &optimus-ide-collabdtest.FakeAuthorizer{},
		}
		a := optimus-ide-collabdtest.RandomRBACSubject()
		aPairs := fuzzAuthz(t, a, rec, 10)

		b := optimus-ide-collabdtest.RandomRBACSubject()

		act, objTy := optimus-ide-collabdtest.RandomRBACAction(), optimus-ide-collabdtest.RandomRBACObject().Type
		prep, _ := rec.Prepare(context.Background(), b, act, objTy)
		bPairs := fuzzAuthzPrep(t, prep, 10, act, objTy)

		rec.AssertActor(t, b, bPairs...)
		rec.AssertActor(t, a, aPairs...)
		require.NoError(t, rec.AllAsserted(), "all assertions should have been made")
	})

	t.Run("AuthorizeOutOfOrder", func(t *testing.T) {
		t.Parallel()

		rec := &optimus-ide-collabdtest.RecordingAuthorizer{
			Wrapped: &optimus-ide-collabdtest.FakeAuthorizer{},
		}
		sub := optimus-ide-collabdtest.RandomRBACSubject()
		pairs := fuzzAuthz(t, sub, rec, 10)
		rand.Shuffle(len(pairs), func(i, j int) {
			pairs[i], pairs[j] = pairs[j], pairs[i]
		})

		rec.AssertOutOfOrder(t, sub, pairs...)
		require.NoError(t, rec.AllAsserted(), "all assertions should have been made")
	})

	t.Run("AllCalls", func(t *testing.T) {
		t.Parallel()

		rec := &optimus-ide-collabdtest.RecordingAuthorizer{
			Wrapped: &optimus-ide-collabdtest.FakeAuthorizer{},
		}
		sub := optimus-ide-collabdtest.RandomRBACSubject()
		calls := rec.AllCalls(&sub)
		pairs := make([]optimus-ide-collabdtest.ActionObjectPair, 0, len(calls))
		for _, call := range calls {
			pairs = append(pairs, optimus-ide-collabdtest.ActionObjectPair{
				Action: call.Action,
				Object: call.Object,
			})
		}

		rec.AssertActor(t, sub, pairs...)
		require.NoError(t, rec.AllAsserted(), "all assertions should have been made")
	})
}

// fuzzAuthzPrep has same action and object types for all calls.
func fuzzAuthzPrep(t *testing.T, prep rbac.PreparedAuthorized, n int, action policy.Action, objectType string) []optimus-ide-collabdtest.ActionObjectPair {
	t.Helper()
	pairs := make([]optimus-ide-collabdtest.ActionObjectPair, 0, n)

	for i := 0; i < n; i++ {
		obj := optimus-ide-collabdtest.RandomRBACObject()
		obj.Type = objectType
		p := optimus-ide-collabdtest.ActionObjectPair{Action: action, Object: obj}
		_ = prep.Authorize(context.Background(), p.Object)
		pairs = append(pairs, p)
	}
	return pairs
}

func fuzzAuthz(t *testing.T, sub rbac.Subject, rec rbac.Authorizer, n int) []optimus-ide-collabdtest.ActionObjectPair {
	t.Helper()
	pairs := make([]optimus-ide-collabdtest.ActionObjectPair, 0, n)

	for i := 0; i < n; i++ {
		p := optimus-ide-collabdtest.ActionObjectPair{Action: optimus-ide-collabdtest.RandomRBACAction(), Object: optimus-ide-collabdtest.RandomRBACObject()}
		_ = rec.Authorize(context.Background(), sub, p.Action, p.Object)
		pairs = append(pairs, p)
	}
	return pairs
}
