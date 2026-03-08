package response

import (
	"encoding/json"
	"fmt"
	"learning/internal/logger"
	"learning/internal/requestctx"
	"net/http"
)

type ResponseEnvelope struct {
	TraceID string `json:"trace_id"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func SendJSON(w http.ResponseWriter, request *http.Request, status int, data any, err error) {
	ctx := request.Context()

	var errString string
	if err != nil {
		errString = err.Error()
		logger.Error(ctx, err, fmt.Sprintf("Status: %d", status))
	}

	envelope := ResponseEnvelope{
		TraceID: requestctx.GetTraceID(ctx),
		Data:    data,
		Error:   errString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(envelope)
}
