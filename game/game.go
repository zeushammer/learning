// game.go
package game

import (
	"encoding/json"
	"errors"
	"io"
	"learning/internal/logger"
	"learning/internal/models"
	"learning/internal/queue"
	"learning/internal/response"
	"net/http"
	"strconv"
	"time"
)

var matchmakingQueue queue.MatchQueue

func init() {

	useSharding := true

	if useSharding {
		matchmakingQueue = queue.NewShardedQueue(8)
	} else {
		matchmakingQueue = queue.NewPriorityQueue()
	}
}

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

func (gs *Game) StartMatchWorker(q queue.MatchQueue) {

	for {

		req, ok := q.Dequeue()

		if !ok {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		_ = req
		// process(req)
	}
}

var Queue chan models.MatchRequest

type QueueRequest struct {
	PlayerID string `json:"player_id"`
	Rating   int    `json:"rating"`
	Mode     string `json:"mode"`
	Region   string `json:"region"`
}

func (gs *Game) Queue(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	logger.Info(ctx, "Fetching player data...")

	var req QueueRequest

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	matchReq := models.MatchRequest{
		PlayerID: req.PlayerID,
		Rating:   req.Rating,
		Mode:     req.Mode,
		JoinedAt: time.Now(),
	}

	Queue <- matchReq

	w.WriteHeader(http.StatusAccepted)

	// You can now pass this traceID into your DB methods!

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
		io.WriteString(w, "200 OK id: "+request.PathValue("id"))
	}

	io.WriteString(w, "Hello from matchById")

}
