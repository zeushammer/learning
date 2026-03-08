package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"learning/internal/logger"
	"learning/internal/requestctx"
	"learning/internal/response"

	"github.com/google/uuid"
)

// Middleware is a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

func CreateChain(middlewares ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		// We iterate backwards so the first middleware in the slice
		// is the first one to execute in the request flow.
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		ctx := req.Context()

		next.ServeHTTP(w, req)

		logger.Info(ctx, "Method: %s | RequestURI: %s | Duration: %s", req.Method, req.RequestURI, time.Since(start))
	})
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				// convert the recovery 'any' to error
				var err error
				switch t := rvr.(type) {
				case error:
					err = t
				default:
					err = fmt.Errorf("%v", t)
				}

				w.Header().Set("Connection", "close")
				response.SendJSON(w, request, http.StatusInternalServerError, nil, fmt.Errorf("internal server panic"))
				logger.Error(request.Context(), err, "PANIC RECOVERED! Stack: %s", debug.Stack())
			}
		}()
		next.ServeHTTP(w, request)
	})
}

// RequestID supports Distributed Tracing
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		// Check if a Trace ID already exists in the header
		id := request.Header.Get("X-Trace-ID")

		// If it's missing or empty
		if id == "" {
			// Generate a new one
			id = uuid.New().String()
			// Add it to the response header
			w.Header().Set("X-Trace-ID", id)
		}

		// Create a new context with the trace id
		ctx := requestctx.SetTraceID(request.Context(), id)

		// Pass the NEW request to the next handler
		next.ServeHTTP(w, request.WithContext(ctx))
	})
}
