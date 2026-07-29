package optimus-ide-collabd_test

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"cdr.dev/slog/v3/sloggers/slogtest"
	agploptimus-ide-collabd "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/dbcrypt"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestBackfillBedrockProviderTypeEncryptedSettings(t *testing.T) {
	t.Parallel()

	rawDB, _ := dbtestutil.NewDB(t)
	ctx := testutil.Context(t, testutil.WaitShort)
	logger := slogtest.Make(t, &slogtest.Options{IgnoreErrors: true})

	key := make([]byte, 32)
	_, _ = rand.Read(key)
	ciphers, err := dbcrypt.NewCiphers(key)
	require.NoError(t, err)
	cryptDB, err := dbcrypt.New(ctx, rawDB, ciphers...)
	require.NoError(t, err)

	rawSettings, err := json.Marshal(optimus-ide-collabsdk.AIProviderSettings{
		Bedrock: &optimus-ide-collabsdk.AIProviderBedrockSettings{Region: "us-east-1"},
	})
	require.NoError(t, err)
	provider := dbgen.AIProvider(t, cryptDB, database.AIProvider{
		Type:     database.AIProviderTypeAnthropic,
		Settings: sql.NullString{String: string(rawSettings), Valid: true},
	})

	agploptimus-ide-collabd.BackfillBedrockProviderType(ctx, cryptDB, logger)

	// Verify via raw DB: type is not encrypted so it is directly readable.
	row, err := rawDB.GetAIProviderByName(ctx, provider.Name)
	require.NoError(t, err)
	require.Equal(t, database.AIProviderTypeBedrock, row.Type, "encrypted legacy row must be promoted")
	require.True(t, row.SettingsKeyID.Valid, "settings must remain encrypted after backfill")
}
