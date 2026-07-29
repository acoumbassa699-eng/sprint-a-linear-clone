package chatprovider_test

import (
	"runtime"
	"strings"
	"testing"

	"charm.land/fantasy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/buildinfo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd/chatprovider"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd/chattest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestUserAgent(t *testing.T) {
	t.Parallel()
	ua := chatprovider.UserAgent()

	// Must start with "optimus-ide-collab-agents/" so LLM providers can
	// identify traffic from Optimus-IDE-Collab.
	require.True(t, strings.HasPrefix(ua, "optimus-ide-collab-agents/"),
		"User-Agent should start with 'optimus-ide-collab-agents/', got %q", ua)

	// Must contain the build version.
	assert.Contains(t, ua, buildinfo.Version())

	// Must contain OS/arch.
	assert.Contains(t, ua, runtime.GOOS+"/"+runtime.GOARCH)
}

func TestModelFromConfig_UserAgent(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t, testutil.WaitShort)

	expectedUA := chatprovider.UserAgent()
	called := make(chan struct{})
	serverURL := chattest.NewOpenAI(t, func(req *chattest.OpenAIRequest) chattest.OpenAIResponse {
		assert.Equal(t, expectedUA, req.Header.Get("User-Agent"))
		close(called)
		return chattest.OpenAINonStreamingResponse("hello")
	})

	keys := chatprovider.ProviderAPIKeys{
		ByProvider:        map[string]string{"openai": "test-key"},
		BaseURLByProvider: map[string]string{"openai": serverURL},
	}

	model, err := chatprovider.ModelFromConfig("openai", "gpt-4", keys, expectedUA, nil, nil)
	require.NoError(t, err)

	// Make a real call so Fantasy sends an HTTP request to the
	// fake server, which asserts the User-Agent header.
	_, err = model.Generate(ctx, fantasy.Call{
		Prompt: []fantasy.Message{
			{
				Role: fantasy.MessageRoleUser,
				Content: []fantasy.MessagePart{
					fantasy.TextPart{Text: "hello"},
				},
			},
		},
	})
	require.NoError(t, err)
	_ = testutil.TryReceive(ctx, t, called)
}
