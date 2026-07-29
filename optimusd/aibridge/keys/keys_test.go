package keys_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/aibridge/keys"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/apikey"
)

func TestNew(t *testing.T) {
	t.Parallel()

	params, key, err := keys.New("test-key")
	require.NoError(t, err)
	require.Len(t, key, keys.KeyLength)
	require.Len(t, params.SecretPrefix, keys.KeyPrefixLength)
	require.Equal(t, key[:keys.KeyPrefixLength], params.SecretPrefix)
	require.True(t, apikey.ValidateHash(params.HashedSecret, key))
	require.False(t, apikey.ValidateHash(params.HashedSecret, key[keys.KeyPrefixLength:]))
}
