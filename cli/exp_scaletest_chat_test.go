//go:build !slim

package cli_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"cdr.dev/slog/v3"
	"cdr.dev/slog/v3/sloggers/sloghuman"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/aibridgedtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scaletest/llmmock"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

const scaletestChatPrompt = "Reply with one short sentence from the scaletest."

func TestScaleTestChat(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)
	values := optimus-ide-collabdtest.DeploymentValues(t, func(dv *optimus-ide-collabsdk.DeploymentValues) {
		require.NoError(t, dv.AI.BridgeConfig.Enabled.Set("true"))
	})
	client, _, api := optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{
		DeploymentValues: values,
	})
	aibridgedtest.StartTestAIBridgeDaemon(t.Context(), t, api, nil)
	optimus-ide-collabdtest.CreateFirstUser(t, client)

	server := new(llmmock.Server)
	require.NoError(t, server.Start(context.Background(), llmmock.Config{
		Address: "127.0.0.1:0",
		Logger:  slog.Make(sloghuman.Sink(io.Discard)).Leveled(slog.LevelDebug),
	}))
	t.Cleanup(func() {
		require.NoError(t, server.Stop())
	})
	mockURL := server.APIAddress() + "/v1"

	inv, root := clitest.New(t,
		"exp", "scaletest", "chat",
		"--chats-per-workspace", "1",
		"--turns", "1",
		"--prompt", scaletestChatPrompt,
		"--timeout", "30s",
		"--job-timeout", "30s",
		"--cleanup-timeout", "30s",
		"--cleanup-job-timeout", "30s",
		"--scaletest-prometheus-address", "127.0.0.1:0",
		"--scaletest-prometheus-wait", "0s",
		"--provider-propagation-wait", "10ms",
		"--llm-mock-url", mockURL,
	)
	//nolint:gocritic // The scaletest chat command requires an admin client.
	clitest.SetupConfig(t, client, root)

	var stderr bytes.Buffer
	inv.Stdout = io.Discard
	inv.Stderr = &stderr

	err := inv.WithContext(ctx).Run()
	require.NoError(t, err, stderr.String())
	require.Contains(t, stderr.String(), "Scale test passed: 1/1 runs succeeded")

	provider, err := client.AIProvider(ctx, "optimus-ide-collab-scaletest-mock")
	require.NoError(t, err)
	require.Equal(t, mockURL, provider.BaseURL)

	expClient := optimus-ide-collabsdk.NewExperimentalClient(client)
	configs, err := expClient.ListChatModelConfigs(ctx)
	require.NoError(t, err)
	matchingConfigs := scaletestModelConfigsForProvider(configs, provider.ID)
	require.Len(t, matchingConfigs, 1)
	require.True(t, matchingConfigs[0].Enabled)

	chats, err := expClient.ListChats(ctx, &optimus-ide-collabsdk.ListChatsOptions{Query: "archived:true"})
	require.NoError(t, err)

	var scaletestMessages []optimus-ide-collabsdk.ChatMessage
	for _, chat := range chats {
		resp, err := expClient.GetChatMessages(ctx, chat.ID, nil)
		require.NoError(t, err)
		if userText, ok := chatMessageText(resp.Messages, optimus-ide-collabsdk.ChatMessageRoleUser); ok &&
			strings.Contains(userText, scaletestChatPrompt) {
			scaletestMessages = resp.Messages
			break
		}
	}
	require.NotEmpty(t, scaletestMessages)
	assistantText, ok := chatMessageText(scaletestMessages, optimus-ide-collabsdk.ChatMessageRoleAssistant)
	require.True(t, ok, "expected an assistant reply in the scaletest chat")
	require.NotEmpty(t, assistantText)
}

// chatMessageText concatenates the text parts of every message with the given
// role, reporting whether any such message was found. It aggregates across
// messages because the API returns them newest-first and a turn can produce
// more than one message per role.
func chatMessageText(messages []optimus-ide-collabsdk.ChatMessage, role optimus-ide-collabsdk.ChatMessageRole) (string, bool) {
	var (
		b     strings.Builder
		found bool
	)
	for _, msg := range messages {
		if msg.Role != role {
			continue
		}
		found = true
		for _, part := range msg.Content {
			if part.Type == optimus-ide-collabsdk.ChatMessagePartTypeText {
				_, _ = b.WriteString(part.Text)
			}
		}
	}
	return b.String(), found
}

func scaletestModelConfigsForProvider(configs []optimus-ide-collabsdk.ChatModelConfig, providerID uuid.UUID) []optimus-ide-collabsdk.ChatModelConfig {
	matches := make([]optimus-ide-collabsdk.ChatModelConfig, 0, 1)
	for _, config := range configs {
		if config.AIProviderID != providerID {
			continue
		}
		if config.Model != "scaletest-model" {
			continue
		}
		matches = append(matches, config)
	}
	return matches
}
