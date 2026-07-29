package optimus-ide-collabsdk_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestChatModelProviderOptions_MarshalJSON_UsesPlainProviderPayload(t *testing.T) {
	t.Parallel()

	sendReasoning := true
	thinkingDisplay := "summarized"

	raw, err := json.Marshal(optimus-ide-collabsdk.ChatModelProviderOptions{
		Anthropic: &optimus-ide-collabsdk.ChatModelAnthropicProviderOptions{
			SendReasoning:   &sendReasoning,
			ThinkingDisplay: &thinkingDisplay,
		},
	})
	require.NoError(t, err)
	require.NotContains(t, string(raw), `"type":"anthropic.options"`)
	require.NotContains(t, string(raw), `"data":`)
	require.Contains(t, string(raw), `"send_reasoning":true`)
	require.Contains(t, string(raw), `"thinking_display":"summarized"`)
}

func TestChatModelProviderOptions_UnmarshalJSON_ParsesPlainProviderPayloads(t *testing.T) {
	t.Parallel()

	raw := []byte(`{
		"anthropic": {
			"send_reasoning": true,
			"thinking_display": "summarized"
		}
	}`)

	var decoded optimus-ide-collabsdk.ChatModelProviderOptions
	err := json.Unmarshal(raw, &decoded)
	require.NoError(t, err)
	require.NotNil(t, decoded.Anthropic)
	require.NotNil(t, decoded.Anthropic.SendReasoning)
	require.True(t, *decoded.Anthropic.SendReasoning)
	require.NotNil(t, decoded.Anthropic.ThinkingDisplay)
	require.Equal(t, "summarized", *decoded.Anthropic.ThinkingDisplay)
}

func TestChatUsageLimitExceededFrom(t *testing.T) {
	t.Parallel()

	t.Run("ExtractsTyped409", func(t *testing.T) {
		t.Parallel()

		want := optimus-ide-collabsdk.ChatUsageLimitExceededResponse{
			Response:    optimus-ide-collabsdk.Response{Message: "Chat usage limit exceeded."},
			SpentMicros: 123,
			LimitMicros: 456,
			ResetsAt:    time.Date(2026, time.March, 16, 12, 0, 0, 0, time.UTC),
		}

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/api/experimental/chats", r.URL.Path)
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusConflict)
			require.NoError(t, json.NewEnoptimus-ide-collab(rw).Encode(want))
		}))
		defer srv.Close()

		serverURL, err := url.Parse(srv.URL)
		require.NoError(t, err)

		client := optimus-ide-collabsdk.NewExperimentalClient(optimus-ide-collabsdk.New(serverURL))
		_, err = client.CreateChat(context.Background(), optimus-ide-collabsdk.CreateChatRequest{
			Content: []optimus-ide-collabsdk.ChatInputPart{{
				Type: optimus-ide-collabsdk.ChatInputPartTypeText,
				Text: "hello",
			}},
		})
		require.Error(t, err)

		sdkErr, ok := optimus-ide-collabsdk.AsError(err)
		require.True(t, ok)
		require.Equal(t, http.StatusConflict, sdkErr.StatusCode())
		require.Equal(t, want.Message, sdkErr.Message)

		limitErr := optimus-ide-collabsdk.ChatUsageLimitExceededFrom(err)
		require.NotNil(t, limitErr)
		require.Equal(t, want, *limitErr)
	})

	t.Run("ReturnsNilForNonLimitErrors", func(t *testing.T) {
		t.Parallel()

		require.Nil(t, optimus-ide-collabsdk.ChatUsageLimitExceededFrom(optimus-ide-collabsdk.NewError(http.StatusConflict, optimus-ide-collabsdk.Response{Message: "plain conflict"})))

		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			require.NoError(t, json.NewEnoptimus-ide-collab(rw).Encode(optimus-ide-collabsdk.Response{Message: "Invalid request."}))
		}))
		defer srv.Close()

		serverURL, err := url.Parse(srv.URL)
		require.NoError(t, err)

		client := optimus-ide-collabsdk.NewExperimentalClient(optimus-ide-collabsdk.New(serverURL))
		_, err = client.CreateChat(context.Background(), optimus-ide-collabsdk.CreateChatRequest{
			Content: []optimus-ide-collabsdk.ChatInputPart{{
				Type: optimus-ide-collabsdk.ChatInputPartTypeText,
				Text: "hello",
			}},
		})
		require.Error(t, err)

		sdkErr, ok := optimus-ide-collabsdk.AsError(err)
		require.True(t, ok)
		require.Equal(t, http.StatusBadRequest, sdkErr.StatusCode())
		require.Nil(t, optimus-ide-collabsdk.ChatUsageLimitExceededFrom(err))
	})
}

func TestChatErrorKind_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	terminal := optimus-ide-collabsdk.ChatError{
		Message: "limit reached",
		Kind:    optimus-ide-collabsdk.ChatErrorKindUsageLimit,
	}
	data, err := json.Marshal(terminal)
	require.NoError(t, err)
	require.Contains(t, string(data), `"kind":"usage_limit"`)

	var decodedTerminal optimus-ide-collabsdk.ChatError
	require.NoError(t, json.Unmarshal(data, &decodedTerminal))
	require.Equal(t, optimus-ide-collabsdk.ChatErrorKindUsageLimit, decodedTerminal.Kind)

	retry := optimus-ide-collabsdk.ChatStreamRetry{
		Attempt: 1,
		Error:   "retrying",
		Kind:    optimus-ide-collabsdk.ChatErrorKindUsageLimit,
	}
	data, err = json.Marshal(retry)
	require.NoError(t, err)
	require.Contains(t, string(data), `"kind":"usage_limit"`)

	var decodedRetry optimus-ide-collabsdk.ChatStreamRetry
	require.NoError(t, json.Unmarshal(data, &decodedRetry))
	require.Equal(t, optimus-ide-collabsdk.ChatErrorKindUsageLimit, decodedRetry.Kind)
}

func TestChatStreamEvent_JSONRoundTripIncludesResetTypesAndPartMetadata(t *testing.T) {
	t.Parallel()

	chatID := uuid.New()
	events := []optimus-ide-collabsdk.ChatStreamEvent{
		{Type: optimus-ide-collabsdk.ChatStreamEventTypePreviewReset, ChatID: chatID},
		{Type: optimus-ide-collabsdk.ChatStreamEventTypeHistoryReset, ChatID: chatID},
		{
			Type:   optimus-ide-collabsdk.ChatStreamEventTypeMessagePart,
			ChatID: chatID,
			MessagePart: &optimus-ide-collabsdk.ChatStreamMessagePart{
				Role:              optimus-ide-collabsdk.ChatMessageRoleAssistant,
				Part:              optimus-ide-collabsdk.ChatMessageText("partial"),
				HistoryVersion:    12,
				GenerationAttempt: 3,
				Seq:               4,
			},
		},
	}
	data, err := json.Marshal(events)
	require.NoError(t, err)
	require.Contains(t, string(data), `"type":"preview_reset"`)
	require.Contains(t, string(data), `"type":"history_reset"`)
	require.Contains(t, string(data), `"history_version":12`)
	require.Contains(t, string(data), `"generation_attempt":3`)
	require.Contains(t, string(data), `"seq":4`)

	var decoded []optimus-ide-collabsdk.ChatStreamEvent
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, optimus-ide-collabsdk.ChatStreamEventTypePreviewReset, decoded[0].Type)
	require.Equal(t, optimus-ide-collabsdk.ChatStreamEventTypeHistoryReset, decoded[1].Type)
	require.Equal(t, int64(12), decoded[2].MessagePart.HistoryVersion)
	require.Equal(t, int64(3), decoded[2].MessagePart.GenerationAttempt)
	require.Equal(t, int64(4), decoded[2].MessagePart.Seq)
}

func TestChatMessagePart_StripInternal(t *testing.T) {
	t.Parallel()

	t.Run("StripsProviderMetadata", func(t *testing.T) {
		t.Parallel()
		part := optimus-ide-collabsdk.ChatMessagePart{
			Type:             optimus-ide-collabsdk.ChatMessagePartTypeToolCall,
			ToolCallID:       "call-1",
			ToolName:         "some_tool",
			Args:             json.RawMessage(`{"key":"value"}`),
			ProviderMetadata: json.RawMessage(`{"type":"ephemeral"}`),
		}
		part.StripInternal()
		assert.Nil(t, part.ProviderMetadata)
		// Public fields preserved.
		assert.Equal(t, optimus-ide-collabsdk.ChatMessagePartTypeToolCall, part.Type)
		assert.Equal(t, "call-1", part.ToolCallID)
		assert.Equal(t, "some_tool", part.ToolName)
		assert.JSONEq(t, `{"key":"value"}`, string(part.Args))
	})

	t.Run("StripsFileDataWhenFileIDSet", func(t *testing.T) {
		t.Parallel()
		id := uuid.New()
		part := optimus-ide-collabsdk.ChatMessagePart{
			Type:      optimus-ide-collabsdk.ChatMessagePartTypeFile,
			FileID:    uuid.NullUUID{UUID: id, Valid: true},
			MediaType: "image/png",
			Data:      []byte("binary-payload"),
		}
		part.StripInternal()
		assert.Nil(t, part.Data)
		assert.Equal(t, id, part.FileID.UUID)
		assert.Equal(t, "image/png", part.MediaType)
	})

	t.Run("PreservesDataWhenNoFileID", func(t *testing.T) {
		t.Parallel()
		part := optimus-ide-collabsdk.ChatMessagePart{
			Type:      optimus-ide-collabsdk.ChatMessagePartTypeFile,
			MediaType: "image/png",
			Data:      []byte("inline-data"),
		}
		part.StripInternal()
		assert.Equal(t, []byte("inline-data"), part.Data)
	})

	t.Run("StripsContextFileContent", func(t *testing.T) {
		t.Parallel()
		agentID := uuid.New()
		part := optimus-ide-collabsdk.ChatMessagePart{
			Type:                     optimus-ide-collabsdk.ChatMessagePartTypeContextFile,
			ContextFilePath:          "/home/optimus-ide-collab/AGENTS.md",
			ContextFileContent:       "large content",
			ContextFileAgentID:       uuid.NullUUID{UUID: agentID, Valid: true},
			ContextFileOS:            "linux",
			ContextFileDirectory:     "/home/optimus-ide-collab/project",
			ContextFileSkillMetaFile: "CUSTOM.md",
		}
		part.StripInternal()
		// Internal fields stripped.
		assert.Empty(t, part.ContextFileContent)
		assert.Empty(t, part.ContextFileOS)
		assert.Empty(t, part.ContextFileDirectory)
		assert.Empty(t, part.ContextFileSkillMetaFile)
		// Public fields preserved.
		assert.Equal(t, "/home/optimus-ide-collab/AGENTS.md", part.ContextFilePath)
		assert.Equal(t, agentID, part.ContextFileAgentID.UUID)
		assert.True(t, part.ContextFileAgentID.Valid)
	})

	t.Run("NoopOnCleanPart", func(t *testing.T) {
		t.Parallel()
		part := optimus-ide-collabsdk.ChatMessageText("hello")
		part.StripInternal()
		assert.Equal(t, "hello", part.Text)
		assert.Equal(t, optimus-ide-collabsdk.ChatMessagePartTypeText, part.Type)
	})
}

func TestChatModelReasoningEffortConfigEnumTags(t *testing.T) {
	t.Parallel()

	want := strings.Join(optimus-ide-collabsdk.ChatModelReasoningEffortValues(), ",")
	typ := reflect.TypeOf(optimus-ide-collabsdk.ChatModelReasoningEffortConfig{})
	for _, fieldName := range []string{"Default", "Max"} {
		field, ok := typ.FieldByName(fieldName)
		require.True(t, ok)
		require.Equal(t, want, field.Tag.Get("enum"))
	}
}

// TestChatMessagePartVariantTags validates the `variants` struct tags
// on ChatMessagePart fields. Every field must either declare variant
// membership or be explicitly excluded, and every known part type
// must appear in at least one tag.
//
// If this test fails, edit the variants struct tags on ChatMessagePart
// in optimus-ide-collabsdk/chats.go.
func TestChatMessagePartVariantTags(t *testing.T) {
	t.Parallel()

	const editHint = "edit the variants struct tags on ChatMessagePart in optimus-ide-collabsdk/chats.go"

	// Fields intentionally excluded from all generated variants.
	// If you add a new field to ChatMessagePart, either add a
	// variants tag or add it here with a comment explaining why.
	excludedFields := map[string]string{
		"type":                         "discriminant, added automatically by codegen",
		"provider_metadata":            "internal only, stripped by db2sdk before API responses",
		"context_file_content":         "internal only, stripped before API responses (typescript:\"-\")",
		"context_file_os":              "internal only, used during prompt expansion (typescript:\"-\")",
		"context_file_directory":       "internal only, used during prompt expansion (typescript:\"-\")",
		"skill_dir":                    "internal only, used by read_skill tools (typescript:\"-\")",
		"context_file_skill_meta_file": "internal only, restored on subsequent turns (typescript:\"-\")",
	}
	knownTypes := make(map[optimus-ide-collabsdk.ChatMessagePartType]bool)
	for _, pt := range optimus-ide-collabsdk.AllChatMessagePartTypes() {
		knownTypes[pt] = true
	}

	// Parse all variants tags from the struct and validate them.
	typ := reflect.TypeOf(optimus-ide-collabsdk.ChatMessagePart{})
	coveredTypes := make(map[optimus-ide-collabsdk.ChatMessagePartType]bool)

	for i := range typ.NumField() {
		f := typ.Field(i)
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		jsonName, _, _ := strings.Cut(jsonTag, ",")

		varTag := f.Tag.Get("variants")
		if varTag == "" {
			assert.Contains(t, excludedFields, jsonName,
				"field %s (json:%q) has no variants tag and is not in excludedFields; %s",
				f.Name, jsonName, editHint)
			continue
		}

		assert.NotEqual(t, "type", jsonName,
			"the discriminant field must not have a variants tag; %s", editHint)

		for _, entry := range strings.Split(varTag, ",") {
			typeLit := optimus-ide-collabsdk.ChatMessagePartType(strings.TrimSuffix(entry, "?"))

			assert.True(t, knownTypes[typeLit],
				"field %s variants tag references unknown type %q; %s",
				f.Name, typeLit, editHint)

			coveredTypes[typeLit] = true
		}
	}

	// Every known type must appear in at least one variants tag.
	for pt := range knownTypes {
		assert.True(t, coveredTypes[pt],
			"ChatMessagePartType %q is not referenced by any variants tag; %s", pt, editHint)
	}

	// Enforce the omitempty <-> variants invariant:
	//   required in any variant  => must NOT have omitempty
	//   optional in all variants => MUST have omitempty
	// See the struct comment on ChatMessagePart for rationale.
	t.Run("omitempty must match variant optionality", func(t *testing.T) {
		t.Parallel()

		typ := reflect.TypeOf(optimus-ide-collabsdk.ChatMessagePart{})
		for i := range typ.NumField() {
			f := typ.Field(i)
			varTag := f.Tag.Get("variants")
			if varTag == "" {
				continue
			}

			allOptional := true
			for _, entry := range strings.Split(varTag, ",") {
				if !strings.HasSuffix(entry, "?") {
					allOptional = false
					break
				}
			}

			jsonTag := f.Tag.Get("json")
			hasOmitEmpty := strings.Contains(jsonTag, "omitempty")

			if !allOptional {
				assert.False(t, hasOmitEmpty,
					"field %s is required in at least one variant but has omitempty in its json tag; "+
						"remove omitempty so Go does not silently drop the zero value that TypeScript expects to always be present",
					f.Name)
			} else {
				assert.True(t, hasOmitEmpty,
					"field %s is optional in all variants but is missing omitempty in its json tag; "+
						"add omitempty to avoid sending zero values for fields the frontend does not expect",
					f.Name)
			}
		}
	})
}

func TestChatMessagePart_CreatedAt_JSON(t *testing.T) {
	t.Parallel()

	t.Run("RoundTrips", func(t *testing.T) {
		t.Parallel()
		ts := time.Date(2025, 6, 15, 12, 30, 0, 0, time.UTC)
		part := optimus-ide-collabsdk.ChatMessagePart{
			Type:       optimus-ide-collabsdk.ChatMessagePartTypeToolCall,
			ToolCallID: "tc-1",
			ToolName:   "execute",
			CreatedAt:  &ts,
		}
		data, err := json.Marshal(part)
		require.NoError(t, err)
		require.Contains(t, string(data), `"created_at"`)

		var decoded optimus-ide-collabsdk.ChatMessagePart
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)
		require.NotNil(t, decoded.CreatedAt)
		require.True(t, ts.Equal(*decoded.CreatedAt))
	})

	t.Run("OmittedWhenNil", func(t *testing.T) {
		t.Parallel()
		part := optimus-ide-collabsdk.ChatMessagePart{
			Type:       optimus-ide-collabsdk.ChatMessagePartTypeToolCall,
			ToolCallID: "tc-1",
			ToolName:   "execute",
		}
		data, err := json.Marshal(part)
		require.NoError(t, err)
		require.NotContains(t, string(data), `"created_at"`)
	})
}

func TestChatMessagePart_ReasoningTimestamps_JSON(t *testing.T) {
	t.Parallel()

	t.Run("RoundTrips", func(t *testing.T) {
		t.Parallel()
		startedAt := time.Date(2025, 6, 15, 12, 30, 0, 0, time.UTC)
		completedAt := startedAt.Add(2 * time.Second)
		part := optimus-ide-collabsdk.ChatMessagePart{
			Type:        optimus-ide-collabsdk.ChatMessagePartTypeReasoning,
			Text:        "thinking out loud",
			CreatedAt:   &startedAt,
			CompletedAt: &completedAt,
		}
		data, err := json.Marshal(part)
		require.NoError(t, err)
		require.Contains(t, string(data), `"created_at"`)
		require.Contains(t, string(data), `"completed_at"`)

		var decoded optimus-ide-collabsdk.ChatMessagePart
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)
		require.NotNil(t, decoded.CreatedAt)
		require.NotNil(t, decoded.CompletedAt)
		require.True(t, startedAt.Equal(*decoded.CreatedAt))
		require.True(t, completedAt.Equal(*decoded.CompletedAt))
	})

	t.Run("OmittedWhenNil", func(t *testing.T) {
		t.Parallel()
		part := optimus-ide-collabsdk.ChatMessagePart{
			Type: optimus-ide-collabsdk.ChatMessagePartTypeReasoning,
			Text: "thinking out loud",
		}
		data, err := json.Marshal(part)
		require.NoError(t, err)
		require.NotContains(t, string(data), `"created_at"`)
		require.NotContains(t, string(data), `"completed_at"`)
	})

	t.Run("LegacyCreatedAtWithoutCompletedAt", func(t *testing.T) {
		t.Parallel()
		// CompletedAt is omitted on messages persisted before this
		// feature shipped. Confirm round-trip leaves CompletedAt nil
		// while preserving CreatedAt so legacy data does not break
		// API consumers.
		startedAt := time.Date(2025, 6, 15, 12, 30, 0, 0, time.UTC)
		part := optimus-ide-collabsdk.ChatMessagePart{
			Type:      optimus-ide-collabsdk.ChatMessagePartTypeReasoning,
			Text:      "legacy reasoning",
			CreatedAt: &startedAt,
		}
		data, err := json.Marshal(part)
		require.NoError(t, err)
		require.Contains(t, string(data), `"created_at"`)
		require.NotContains(t, string(data), `"completed_at"`)

		var decoded optimus-ide-collabsdk.ChatMessagePart
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)
		require.NotNil(t, decoded.CreatedAt)
		require.Nil(t, decoded.CompletedAt)
	})
}

func TestModelCostConfig_LegacyNumericJSON(t *testing.T) {
	t.Parallel()

	var decoded optimus-ide-collabsdk.ModelCostConfig
	err := json.Unmarshal([]byte("{\"input_price_per_million_tokens\": 1.5}"), &decoded)
	require.NoError(t, err)
	require.NotNil(t, decoded.InputPricePerMillionTokens)
	require.True(t, decoded.InputPricePerMillionTokens.Equal(decimal.RequireFromString("1.5")))
}

func TestModelCostConfig_QuotedDecimalJSON(t *testing.T) {
	t.Parallel()

	var decoded optimus-ide-collabsdk.ModelCostConfig
	err := json.Unmarshal([]byte("{\"input_price_per_million_tokens\": \"1.5\"}"), &decoded)
	require.NoError(t, err)
	require.NotNil(t, decoded.InputPricePerMillionTokens)
	require.True(t, decoded.InputPricePerMillionTokens.Equal(decimal.RequireFromString("1.5")))
}

func TestModelCostConfig_NilVsZero(t *testing.T) {
	t.Parallel()

	zero := decimal.Zero
	raw, err := json.Marshal(struct {
		Nil  optimus-ide-collabsdk.ModelCostConfig `json:"nil"`
		Zero optimus-ide-collabsdk.ModelCostConfig `json:"zero"`
	}{
		Nil:  optimus-ide-collabsdk.ModelCostConfig{},
		Zero: optimus-ide-collabsdk.ModelCostConfig{InputPricePerMillionTokens: &zero},
	})
	require.NoError(t, err)
	require.Contains(t, string(raw), "\"zero\":{\"input_price_per_million_tokens\":\"0\"}")
	require.Contains(t, string(raw), "\"nil\":{}")
}

func TestChatModelCallConfig_UnmarshalLegacyPricing(t *testing.T) {
	t.Parallel()

	var decoded optimus-ide-collabsdk.ChatModelCallConfig
	err := json.Unmarshal([]byte("{\"input_price_per_million_tokens\": 1.5}"), &decoded)
	require.NoError(t, err)
	require.NotNil(t, decoded.Cost)
	require.NotNil(t, decoded.Cost.InputPricePerMillionTokens)
	require.True(t, decoded.Cost.InputPricePerMillionTokens.Equal(decimal.RequireFromString("1.5")))
}

func TestChatModelCallConfig_UnmarshalStrict(t *testing.T) {
	t.Parallel()

	var decoded optimus-ide-collabsdk.ChatModelCallConfig
	err := decoded.UnmarshalStrict([]byte(`{
		"temperature": 0.5,
		"cost": {"input_price_per_million_tokens": "5"},
		"input_price_per_million_tokens": 1.5,
		"provider_options": {"anthropic": {"thinking": {"budget_tokens": 1024}}}
	}`))
	require.NoError(t, err)
	require.NotNil(t, decoded.Temperature)
	require.True(t, decoded.Cost.InputPricePerMillionTokens.Equal(decimal.RequireFromString("5")))

	err = decoded.UnmarshalStrict([]byte(`{"provider_options": {"anthropic": {"bogus_setting": true}}}`))
	require.ErrorContains(t, err, `unknown field "bogus_setting"`)

	// Trailing data after the first value is rejected, matching json.Unmarshal.
	err = decoded.UnmarshalStrict([]byte(`{"temperature": 0.5} {"bogus_setting": true}`))
	require.ErrorContains(t, err, "trailing data")

	// UnmarshalJSON stays lenient.
	require.NoError(t, json.Unmarshal([]byte(`{"bogus_setting": true}`), &decoded))
}

func TestChatCostSummary_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := optimus-ide-collabsdk.ChatCostSummary{
		TotalCostMicros: 123,
	}
	raw, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded optimus-ide-collabsdk.ChatCostSummary
	err = json.Unmarshal(raw, &decoded)
	require.NoError(t, err)
	require.Equal(t, original.TotalCostMicros, decoded.TotalCostMicros)
}

// TestChat_JSONRoundTrip verifies that every field of optimus-ide-collabsdk.Chat
// survives a JSON marshal/unmarshal cycle. This catches omitempty
// silently eating zero-ish values, struct tag typos, and similar
// serialization bugs in the pubsub path.
func TestChat_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Microsecond)
	prState := "open"
	prTitle := "test PR"
	authorLogin := "testuser"
	avatarURL := "https://example.com/avatar.png"
	baseBranch := "main"
	headBranch := "feature/test"
	prNumber := int32(42)
	commits := int32(3)
	approved := true
	reviewerCount := int32(2)
	refreshedAt := now
	staleAt := now.Add(time.Hour)
	lastError := &optimus-ide-collabsdk.ChatError{
		Message:    "boom",
		Detail:     "provider detail",
		Kind:       optimus-ide-collabsdk.ChatErrorKindGeneric,
		Provider:   "openai",
		Retryable:  true,
		StatusCode: 503,
	}
	prURL := "https://github.com/optimus-ide-collab/optimus-ide-collab/pull/42"
	workspaceID := uuid.New()
	buildID := uuid.New()
	agentID := uuid.New()
	parentChatID := uuid.New()
	rootChatID := uuid.New()

	original := optimus-ide-collabsdk.Chat{
		ID:                uuid.New(),
		OwnerID:           uuid.New(),
		WorkspaceID:       &workspaceID,
		BuildID:           &buildID,
		AgentID:           &agentID,
		ParentChatID:      &parentChatID,
		RootChatID:        &rootChatID,
		LastModelConfigID: uuid.New(),
		Title:             "round-trip-test",
		Status:            optimus-ide-collabsdk.ChatStatusRunning,
		LastError:         lastError,
		CreatedAt:         now,
		UpdatedAt:         now,
		Archived:          true,
		MCPServerIDs:      []uuid.UUID{uuid.New()},
		Labels:            map[string]string{"env": "prod"},
		DiffStatus: &optimus-ide-collabsdk.ChatDiffStatus{
			ChatID:           uuid.New(),
			URL:              &prURL,
			PullRequestState: &prState,
			PullRequestTitle: prTitle,
			PullRequestDraft: true,
			ChangesRequested: true,
			Additions:        10,
			Deletions:        5,
			ChangedFiles:     3,
			AuthorLogin:      &authorLogin,
			AuthorAvatarURL:  &avatarURL,
			BaseBranch:       &baseBranch,
			HeadBranch:       &headBranch,
			PRNumber:         &prNumber,
			Commits:          &commits,
			Approved:         &approved,
			ReviewerCount:    &reviewerCount,
			RefreshedAt:      &refreshedAt,
			StaleAt:          &staleAt,
		},
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded optimus-ide-collabsdk.Chat
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	require.Equal(t, original, decoded)
}

func TestNewDynamicTool(t *testing.T) {
	t.Parallel()

	type testArgs struct {
		Query string `json:"query"`
	}

	t.Run("CorrectSchema", func(t *testing.T) {
		t.Parallel()

		tool := optimus-ide-collabsdk.NewDynamicTool(
			"search", "search things",
			func(_ context.Context, args testArgs, _ optimus-ide-collabsdk.DynamicToolCall) (optimus-ide-collabsdk.DynamicToolResponse, error) {
				return optimus-ide-collabsdk.DynamicToolResponse{Content: args.Query}, nil
			},
		)

		require.Equal(t, "search", tool.Name)
		require.Equal(t, "search things", tool.Description)
		require.Contains(t, string(tool.InputSchema), `"query"`)
		require.Contains(t, string(tool.InputSchema), `"string"`)
	})

	t.Run("HandlerReceivesArgs", func(t *testing.T) {
		t.Parallel()

		var received testArgs
		tool := optimus-ide-collabsdk.NewDynamicTool(
			"search", "search things",
			func(_ context.Context, args testArgs, _ optimus-ide-collabsdk.DynamicToolCall) (optimus-ide-collabsdk.DynamicToolResponse, error) {
				received = args
				return optimus-ide-collabsdk.DynamicToolResponse{Content: "ok"}, nil
			},
		)

		resp, err := tool.Handler(context.Background(), optimus-ide-collabsdk.DynamicToolCall{
			Args: `{"query":"hello"}`,
		})
		require.NoError(t, err)
		require.Equal(t, "ok", resp.Content)
		require.Equal(t, "hello", received.Query)
	})

	t.Run("InvalidJSONArgs", func(t *testing.T) {
		t.Parallel()

		tool := optimus-ide-collabsdk.NewDynamicTool(
			"search", "search things",
			func(_ context.Context, args testArgs, _ optimus-ide-collabsdk.DynamicToolCall) (optimus-ide-collabsdk.DynamicToolResponse, error) {
				return optimus-ide-collabsdk.DynamicToolResponse{Content: "should not reach"}, nil
			},
		)

		resp, err := tool.Handler(context.Background(), optimus-ide-collabsdk.DynamicToolCall{
			Args: "not-json",
		})
		require.NoError(t, err)
		require.True(t, resp.IsError)
		require.Contains(t, resp.Content, "invalid parameters")
	})
}

//nolint:tparallel,paralleltest
func TestParseChatWorkspaceTTL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{"Empty_ReturnsDefault", "", 0, false},
		{"ValidDuration_Hours", "2h", 2 * time.Hour, false},
		{"ValidDuration_HoursAndMinutes", "2h30m", 2*time.Hour + 30*time.Minute, false},
		{"ValidDuration_Minutes", "90m", 90 * time.Minute, false},
		{"Zero", "0s", 0, false},
		{"Negative", "-1h", 0, true},
		{"Invalid", "not-a-duration", 0, true},
		{"LargeDuration", "720h", 720 * time.Hour, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := optimus-ide-collabsdk.ParseChatWorkspaceTTL(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
