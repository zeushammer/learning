package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"learning/internal/models"
)

var Queue chan models.MatchRequest

type QueueRequest struct {
	PlayerID string `json:"player_id"`
	Rating   int    `json:"rating"`
	Mode     string `json:"mode"`
	Region   string `json:"region"`
}

func QueueHandler(w http.ResponseWriter, r *http.Request) {

	var req QueueRequest

	err := json.NewDecoder(r.Body).Decode(&req)
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
}
