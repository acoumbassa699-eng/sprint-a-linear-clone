package optimus-ide-collabsdk_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestAPIAllowListTarget_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	all := optimus-ide-collabsdk.AllowAllTarget()
	b, err := json.Marshal(all)
	require.NoError(t, err)
	require.JSONEq(t, `"*:*"`, string(b))
	var rt optimus-ide-collabsdk.APIAllowListTarget
	require.NoError(t, json.Unmarshal(b, &rt))
	require.Equal(t, optimus-ide-collabsdk.ResourceWildcard, rt.Type)
	require.Equal(t, policy.WildcardSymbol, rt.ID)

	ty := optimus-ide-collabsdk.AllowTypeTarget(optimus-ide-collabsdk.ResourceWorkspace)
	b, err = json.Marshal(ty)
	require.NoError(t, err)
	require.JSONEq(t, `"workspace:*"`, string(b))
	require.NoError(t, json.Unmarshal(b, &rt))
	require.Equal(t, optimus-ide-collabsdk.ResourceWorkspace, rt.Type)
	require.Equal(t, policy.WildcardSymbol, rt.ID)

	id := uuid.New()
	res := optimus-ide-collabsdk.AllowResourceTarget(optimus-ide-collabsdk.ResourceTemplate, id)
	b, err = json.Marshal(res)
	require.NoError(t, err)
	exp := `"template:` + id.String() + `"`
	require.JSONEq(t, exp, string(b))
}
