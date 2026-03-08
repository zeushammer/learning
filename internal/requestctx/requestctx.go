package requestctx

import "context"

// Define a custom type for the context key to prevent collisions

// By defining `type contextKey string`, your key is longer just a string
// it is a `contextKey`. Even if another package has a `type otherKey string`
// with the same value ("unique-trace-id-key-identifer-goes-here"), Go treats
// them as entirely different entities.
type contextKey string

const TraceIDKey contextKey = "unique-trace-id-key-identifer-goes-here"

// GetTraceID safely extracts the ID. If not found, it returns "unknown".
// This prevents the app from crashing, but would be, weird, if it happened.
func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		return id
	}
	return "unknown"
}

// SetTraceID is used by the middle to inject the ID.
func SetTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, TraceIDKey, id)
}
