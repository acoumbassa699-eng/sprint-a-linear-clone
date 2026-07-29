package chaterror_test

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd/chaterror"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestTerminalErrorPayloadUsesNormalizedClassification(t *testing.T) {
	t.Parallel()

	classified := chaterror.Classify(
		xerrors.New("azure openai received status 429 from upstream"),
	)
	payload := chaterror.TerminalErrorPayload(classified)

	require.Equal(t, &optimus-ide-collabsdk.ChatError{
		Message:    "Azure OpenAI is rate limiting requests.",
		Detail:     "azure openai received status 429 from upstream",
		Kind:       optimus-ide-collabsdk.ChatErrorKindRateLimit,
		Provider:   "azure",
		Retryable:  true,
		StatusCode: 429,
	}, payload)
}

func TestTerminalErrorPayloadIncludesProviderDetail(t *testing.T) {
	t.Parallel()

	payload := chaterror.TerminalErrorPayload(chaterror.Classify(testProviderError(
		"",
		400,
		nil,
		testProviderResponseDump(`{"error":{"message":"Image exceeds 5 MB maximum."}}`),
	)))

	require.Equal(t, "Image exceeds 5 MB maximum.", payload.Detail)
}

func TestTerminalErrorPayloadNilForEmptyClassification(t *testing.T) {
	t.Parallel()

	require.Nil(t, chaterror.TerminalErrorPayload(chaterror.ClassifiedError{}))
}

func TestStreamRetryPayloadPreservesRetryableMessage(t *testing.T) {
	t.Parallel()

	delay := 3 * time.Second
	classified := chaterror.Classify(xerrors.Errorf(
		"anthropic stream closed before message_stop: %w",
		io.EOF,
	))
	payload := chaterror.StreamRetryPayload(2, delay, classified)

	require.NotNil(t, payload)
	require.Equal(t,
		"Anthropic stream closed unexpectedly before the response completed.",
		payload.Error,
	)
	require.Equal(t, optimus-ide-collabsdk.ChatErrorKindTimeout, payload.Kind)
	require.Equal(t, "anthropic", payload.Provider)
}

func TestStreamRetryPayloadUsesNormalizedClassification(t *testing.T) {
	t.Parallel()

	delay := 3 * time.Second
	startedAt := time.Now()
	payload := chaterror.StreamRetryPayload(2, delay, chaterror.ClassifiedError{
		Message:    "OpenAI returned an unexpected error.",
		Kind:       optimus-ide-collabsdk.ChatErrorKindGeneric,
		Provider:   "openai",
		Retryable:  true,
		StatusCode: 503,
	})

	require.NotNil(t, payload)
	require.Equal(t, 2, payload.Attempt)
	require.Equal(t, delay.Milliseconds(), payload.DelayMs)
	// Retry messages omit the HTTP status code; the status code is
	// surfaced separately in the payload's StatusCode field.
	require.Equal(t, "OpenAI returned an unexpected error.", payload.Error)
	require.Equal(t, optimus-ide-collabsdk.ChatErrorKindGeneric, payload.Kind)
	require.Equal(t, "openai", payload.Provider)
	require.Equal(t, 503, payload.StatusCode)
	require.WithinDuration(t, startedAt.Add(delay), payload.RetryingAt, time.Second)
}
