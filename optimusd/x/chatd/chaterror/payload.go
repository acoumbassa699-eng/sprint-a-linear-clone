package chaterror

import (
	"time"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TerminalErrorPayload(classified ClassifiedError) *optimus-ide-collabsdk.ChatError {
	if classified.Message == "" {
		return nil
	}
	return &optimus-ide-collabsdk.ChatError{
		Message:    classified.Message,
		Detail:     classified.Detail,
		Kind:       classified.Kind,
		Provider:   classified.Provider,
		Retryable:  classified.Retryable,
		StatusCode: classified.StatusCode,
	}
}

func StreamRetryPayload(
	attempt int,
	delay time.Duration,
	classified ClassifiedError,
) *optimus-ide-collabsdk.ChatStreamRetry {
	if classified.Message == "" {
		return nil
	}
	return &optimus-ide-collabsdk.ChatStreamRetry{
		Attempt:    attempt,
		DelayMs:    delay.Milliseconds(),
		Error:      retryMessage(classified),
		Kind:       classified.Kind,
		Provider:   classified.Provider,
		StatusCode: classified.StatusCode,
		RetryingAt: time.Now().Add(delay),
	}
}
