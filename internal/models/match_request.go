package models

import "time"

type MatchRequest struct {
	PlayerID string
	Rating   int
	Mode     string
	Reqion   string
	JoinedAt time.Time
}
