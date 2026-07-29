package chatd_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestChatPersonalModelOverrideKey(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"chat_personal_model_override:root",
		chatd.ChatPersonalModelOverrideKey(optimus-ide-collabsdk.ChatPersonalModelOverrideContextRoot),
	)
}

func TestParseChatPersonalModelOverride(t *testing.T) {
	t.Parallel()

	modelConfigID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	tests := []struct {
		name        string
		raw         string
		defaultMode optimus-ide-collabsdk.ChatPersonalModelOverrideMode
		want        chatd.ParsedChatPersonalModelOverride
	}{
		{
			name:        "EmptyUsesDefault",
			raw:         "",
			defaultMode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
			want: chatd.ParsedChatPersonalModelOverride{
				Mode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
			},
		},
		{
			name:        "ChatDefault",
			raw:         string(optimus-ide-collabsdk.ChatPersonalModelOverrideModeChatDefault),
			defaultMode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
			want: chatd.ParsedChatPersonalModelOverride{
				Mode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeChatDefault,
			},
		},
		{
			name:        "DeploymentDefault",
			raw:         string(optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault),
			defaultMode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeChatDefault,
			want: chatd.ParsedChatPersonalModelOverride{
				Mode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
			},
		},
		{
			name:        "Model",
			raw:         "model:" + modelConfigID.String(),
			defaultMode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
			want: chatd.ParsedChatPersonalModelOverride{
				Mode:          optimus-ide-collabsdk.ChatPersonalModelOverrideModeModel,
				ModelConfigID: modelConfigID,
			},
		},
		{
			name:        "ModelWithReasoningEffort",
			raw:         "model:" + modelConfigID.String() + ":high",
			defaultMode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
			want: chatd.ParsedChatPersonalModelOverride{
				Mode:            optimus-ide-collabsdk.ChatPersonalModelOverrideModeModel,
				ModelConfigID:   modelConfigID,
				ReasoningEffort: ptr.Ref("high"),
			},
		},
		{
			name:        "ModelWithEmptyReasoningEffort",
			raw:         "model:" + modelConfigID.String() + ":",
			defaultMode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
			want: chatd.ParsedChatPersonalModelOverride{
				Mode:      optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
				Malformed: true,
			},
		},
		{
			name:        "InvalidModelUUID",
			raw:         "model:not-a-uuid",
			defaultMode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
			want: chatd.ParsedChatPersonalModelOverride{
				Mode:      optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
				Malformed: true,
			},
		},
		{
			name:        "UnknownValue",
			raw:         "unknown",
			defaultMode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeChatDefault,
			want: chatd.ParsedChatPersonalModelOverride{
				Mode:      optimus-ide-collabsdk.ChatPersonalModelOverrideModeChatDefault,
				Malformed: true,
			},
		},
		{
			name:        "OuterWhitespace",
			raw:         " \tmodel:" + modelConfigID.String() + "\n",
			defaultMode: optimus-ide-collabsdk.ChatPersonalModelOverrideModeDeploymentDefault,
			want: chatd.ParsedChatPersonalModelOverride{
				Mode:          optimus-ide-collabsdk.ChatPersonalModelOverrideModeModel,
				ModelConfigID: modelConfigID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := chatd.ParseChatPersonalModelOverride(tt.raw, tt.defaultMode)
			require.Equal(t, tt.want, got)
		})
	}
}
