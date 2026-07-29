package chatloop

import (
	"context"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type messagePartPublisherKey struct{}

// WithMessagePartPublisher returns a context carrying the streaming
// message-part publisher so tools can stream intermediate output (e.g.
// advisor advice deltas) while they execute. ExecuteLocalTools injects
// the publisher before running tools.
func WithMessagePartPublisher(
	ctx context.Context,
	publish func(optimus-ide-collabsdk.ChatMessageRole, optimus-ide-collabsdk.ChatMessagePart),
) context.Context {
	if publish == nil {
		return ctx
	}
	return context.WithValue(ctx, messagePartPublisherKey{}, publish)
}

// MessagePartPublisherFromContext returns the publisher injected by
// ExecuteLocalTools, or nil when absent.
func MessagePartPublisherFromContext(
	ctx context.Context,
) func(optimus-ide-collabsdk.ChatMessageRole, optimus-ide-collabsdk.ChatMessagePart) {
	publish, _ := ctx.Value(messagePartPublisherKey{}).(func(optimus-ide-collabsdk.ChatMessageRole, optimus-ide-collabsdk.ChatMessagePart))
	return publish
}
