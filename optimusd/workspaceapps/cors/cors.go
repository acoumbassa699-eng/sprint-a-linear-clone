package cors

import (
	"context"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type contextKeyBehavior struct{}

// WithBehavior sets the CORS behavior for the given context.
func WithBehavior(ctx context.Context, behavior optimus-ide-collabsdk.CORSBehavior) context.Context {
	return context.WithValue(ctx, contextKeyBehavior{}, behavior)
}

// HasBehavior returns true if the given context has the specified CORS behavior.
func HasBehavior(ctx context.Context, behavior optimus-ide-collabsdk.CORSBehavior) bool {
	val := ctx.Value(contextKeyBehavior{})
	b, ok := val.(optimus-ide-collabsdk.CORSBehavior)
	return ok && b == behavior
}
