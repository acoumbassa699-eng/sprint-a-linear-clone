package chat

import (
	"context"
	"io"

	"github.com/google/uuid"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type chatClient interface {
	SetLogger(logger slog.Logger)
	SetLogBodies(logBodies bool)
	CreateChat(ctx context.Context, req optimus-ide-collabsdk.CreateChatRequest) (optimus-ide-collabsdk.Chat, error)
	StreamChat(ctx context.Context, chatID uuid.UUID, opts *optimus-ide-collabsdk.StreamChatOptions) (<-chan optimus-ide-collabsdk.ChatStreamEvent, io.Closer, error)
	CreateChatMessage(ctx context.Context, chatID uuid.UUID, req optimus-ide-collabsdk.CreateChatMessageRequest) (optimus-ide-collabsdk.CreateChatMessageResponse, error)
	UpdateChat(ctx context.Context, chatID uuid.UUID, req optimus-ide-collabsdk.UpdateChatRequest) error
}

var _ chatClient = (*optimus-ide-collabsdk.ExperimentalClient)(nil)

type chatModelConfigClient interface {
	ListChatModelConfigs(ctx context.Context) ([]optimus-ide-collabsdk.ChatModelConfig, error)
	CreateChatModelConfig(ctx context.Context, req optimus-ide-collabsdk.CreateChatModelConfigRequest) (optimus-ide-collabsdk.ChatModelConfig, error)
}

var _ chatModelConfigClient = (*optimus-ide-collabsdk.ExperimentalClient)(nil)
