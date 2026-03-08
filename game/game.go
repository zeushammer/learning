// game.go
package game

import (
	"errors"
	"io"
	"learning/internal/logger"
	"learning/internal/response"
	"net/http"
	"strconv"
	"time"
)

type Game struct{}

func New() *Game {
	g := &Game{}
	return g
}

func (gs *Game) Player(w http.ResponseWriter, req *http.Request) {
	// todo: use WithCancelCause instead
	// rednafi.com/go/context-cancellation-cause
	ctx := req.Context()

	logger.Info(req.Context(), "server: player handler started")
	defer logger.Info(ctx, "server: player handler ended")

	select {
	case <-time.After(10 * time.Second):
		io.WriteString(w, "Hello from player")
	case <-ctx.Done():
		err := ctx.Err()
		logger.Error(ctx, err, "")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (gs *Game) Queue(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	logger.Info(ctx, "Fetching player data...")
	// You can now pass this traceID into your DB methods!

	io.WriteString(w, "Hello from queue")
}

func (gs *Game) QueuePlayer(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	logger.Info(ctx, "Adding Player ID...")

	response.SendJSON(w, request, http.StatusOK, nil, nil)
}

func (gs *Game) MatchById(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger.Info(ctx, "handling get task at %s\n", request.URL.Path)

	if id, err := strconv.Atoi(request.PathValue("id")); err != nil {
		response.SendJSON(w, request, http.StatusNotFound, nil, errors.New("user not found"))
		return
	} else {
		response.SendJSON(w, request, http.StatusOK, id, nil)
		io.WriteString(w, "200 OK id: "+string(id))
	}

	io.WriteString(w, "Hello from matchById")

}
