package logger

import (
	"context"
	"learning/internal/requestctx"
	"log"
)

// Info logs a message along with the TraceID found in context
func Info(ctx context.Context, msg string, args ...any) {
	trace := requestctx.GetTraceID(ctx)

	// Format: [INFO][Trace] your message
	prefix := "[INFO][" + trace + "] "
	log.Printf(prefix+msg, args...)
}

// Error logs and error along with the Trace ID
func Error(ctx context.Context, err error, msg string, args ...any) {
	trace := requestctx.GetTraceID(ctx)

	prefix := "[ERROR][" + trace + "] "
	log.Printf(prefix+err.Error()+msg, args...)
}
