package chat

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

const (
	scaletestAIProviderType        = optimus-ide-collabsdk.AIProviderTypeOpenAICompat
	scaletestAIProviderName        = "optimus-ide-collab-scaletest-mock"
	scaletestAIProviderDisplayName = "Scaletest LLM Mock"
	scaletestAIProviderAPIKey      = "optimus-ide-collab-scaletest"
	scaletestModelName             = "scaletest-model"
	scaletestModelDisplayName      = "Scaletest Model"
	scaletestModelContextLimit     = int64(4096)
)

// DefaultProviderPropagationWait is how long to wait after creating or
// updating the mock LLM provider before starting chats. Provider config is
// cached per optimus-ide-collabd replica with a 10 second TTL (see
// optimus-ide-collabd/x/chatd/configcache.go), and a change is only guaranteed to be
// visible everywhere once every replica's cached entry has expired. 15
// seconds comfortably exceeds that TTL.
const DefaultProviderPropagationWait = 15 * time.Second

type scaletestAIProviderAction string

const (
	scaletestAIProviderActionCreated scaletestAIProviderAction = "created"
	scaletestAIProviderActionUpdated scaletestAIProviderAction = "updated"
	scaletestAIProviderActionReused  scaletestAIProviderAction = "reused"
)

// EnsureScaletestModelConfig bootstraps the shared AI provider and model
// config used by chat scaletests. When the provider was created or updated,
// it sleeps for propagationWait so every optimus-ide-collabd replica's cached provider
// config expires before chats start.
func EnsureScaletestModelConfig(ctx context.Context, client *optimus-ide-collabsdk.Client, logger slog.Logger, llmMockURL string, propagationWait time.Duration) (uuid.UUID, error) {
	expClient := optimus-ide-collabsdk.NewExperimentalClient(client)

	logger.Info(ctx, "bootstrapping mock LLM provider", slog.F("llm_mock_url", llmMockURL))

	provider, providerAction, err := ensureScaletestAIProvider(ctx, expClient, llmMockURL)
	if err != nil {
		return uuid.Nil, err
	}

	switch providerAction {
	case scaletestAIProviderActionCreated:
		logger.Info(ctx, "created mock LLM provider",
			slog.F("provider_name", provider.Name),
			slog.F("provider_id", provider.ID),
			slog.F("llm_mock_url", llmMockURL),
		)
	case scaletestAIProviderActionUpdated:
		logger.Info(ctx, "updated mock LLM provider",
			slog.F("provider_name", provider.Name),
			slog.F("provider_id", provider.ID),
			slog.F("llm_mock_url", llmMockURL),
		)
	case scaletestAIProviderActionReused:
		logger.Info(ctx, "reusing mock LLM provider",
			slog.F("provider_name", provider.Name),
			slog.F("provider_id", provider.ID),
		)
	}

	modelConfigID, err := ensureScaletestChatModelConfig(ctx, expClient, logger, provider)
	if err != nil {
		return uuid.Nil, err
	}

	if providerAction != scaletestAIProviderActionReused && propagationWait > 0 {
		logger.Info(ctx, "waiting for mock LLM provider propagation",
			slog.F("provider_name", provider.Name),
			slog.F("wait", propagationWait),
		)
		select {
		case <-ctx.Done():
			return uuid.Nil, ctx.Err()
		case <-time.After(propagationWait):
		}
	}

	return modelConfigID, nil
}

func ensureScaletestChatModelConfig(ctx context.Context, client chatModelConfigClient, logger slog.Logger, provider optimus-ide-collabsdk.AIProvider) (uuid.UUID, error) {
	modelConfigs, err := client.ListChatModelConfigs(ctx)
	if err != nil {
		return uuid.Nil, xerrors.Errorf("list chat model configs: %w", err)
	}

	for i := range modelConfigs {
		matchesProvider := modelConfigs[i].AIProviderID == provider.ID
		matchesModel := modelConfigs[i].Model == scaletestModelName
		if !matchesProvider || !matchesModel {
			continue
		}
		if !modelConfigs[i].Enabled {
			return uuid.Nil, xerrors.Errorf("existing scaletest chat model config %s is disabled; re-enable or delete it before running scaletests", modelConfigs[i].ID)
		}
		modelConfigID := modelConfigs[i].ID
		logger.Info(ctx, "reusing scaletest model config", slog.F("model_config_id", modelConfigID))
		return modelConfigID, nil
	}

	enabled := true
	isDefault := false
	contextLimit := scaletestModelContextLimit
	created, err := client.CreateChatModelConfig(ctx, optimus-ide-collabsdk.CreateChatModelConfigRequest{
		AIProviderID: &provider.ID,
		Model:        scaletestModelName,
		DisplayName:  scaletestModelDisplayName,
		Enabled:      &enabled,
		IsDefault:    &isDefault,
		ContextLimit: &contextLimit,
	})
	if err != nil {
		return uuid.Nil, xerrors.Errorf("create scaletest chat model config: %w", err)
	}
	logger.Info(ctx, "created scaletest model config", slog.F("model_config_id", created.ID))
	return created.ID, nil
}

func ensureScaletestAIProvider(ctx context.Context, client *optimus-ide-collabsdk.ExperimentalClient, llmMockURL string) (optimus-ide-collabsdk.AIProvider, scaletestAIProviderAction, error) {
	provider, err := client.AIProvider(ctx, scaletestAIProviderName)
	if err != nil {
		var sdkErr *optimus-ide-collabsdk.Error
		if !xerrors.As(err, &sdkErr) || sdkErr.StatusCode() != http.StatusNotFound {
			return optimus-ide-collabsdk.AIProvider{}, "", xerrors.Errorf("look up scaletest AI provider: %w", err)
		}

		created, err := client.CreateAIProvider(ctx, optimus-ide-collabsdk.CreateAIProviderRequest{
			Type:        scaletestAIProviderType,
			Name:        scaletestAIProviderName,
			DisplayName: scaletestAIProviderDisplayName,
			Enabled:     true,
			BaseURL:     llmMockURL,
			APIKeys:     []string{scaletestAIProviderAPIKey},
		})
		if err == nil {
			return created, scaletestAIProviderActionCreated, nil
		}

		sdkErr = nil
		if !xerrors.As(err, &sdkErr) || sdkErr.StatusCode() != http.StatusConflict {
			return optimus-ide-collabsdk.AIProvider{}, "", xerrors.Errorf("create scaletest AI provider: %w", err)
		}

		provider, err = client.AIProvider(ctx, scaletestAIProviderName)
		if err != nil {
			return optimus-ide-collabsdk.AIProvider{}, "", xerrors.Errorf("look up scaletest AI provider after conflict: %w", err)
		}
	}

	if provider.Type != scaletestAIProviderType {
		return optimus-ide-collabsdk.AIProvider{}, "", xerrors.Errorf("refusing to use scaletest AI provider %s with type %q", provider.ID, provider.Type)
	}
	if provider.DisplayName != scaletestAIProviderDisplayName {
		return optimus-ide-collabsdk.AIProvider{}, "", xerrors.Errorf("refusing to use scaletest AI provider %s with display name %q", provider.ID, provider.DisplayName)
	}
	if !provider.Enabled {
		return optimus-ide-collabsdk.AIProvider{}, "", xerrors.Errorf("existing scaletest AI provider %s is disabled; re-enable or delete it before running scaletests", provider.ID)
	}

	var update optimus-ide-collabsdk.UpdateAIProviderRequest
	needsUpdate := false
	if provider.BaseURL != llmMockURL {
		update.BaseURL = &llmMockURL
		needsUpdate = true
	}
	if len(provider.APIKeys) == 0 {
		apiKey := scaletestAIProviderAPIKey
		apiKeys := []optimus-ide-collabsdk.AIProviderKeyMutation{{APIKey: &apiKey}}
		update.APIKeys = &apiKeys
		needsUpdate = true
	}
	if !needsUpdate {
		return provider, scaletestAIProviderActionReused, nil
	}

	updated, err := client.UpdateAIProvider(ctx, scaletestAIProviderName, update)
	if err != nil {
		return optimus-ide-collabsdk.AIProvider{}, "", xerrors.Errorf("update scaletest AI provider: %w", err)
	}
	return updated, scaletestAIProviderActionUpdated, nil
}
