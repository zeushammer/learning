package models

import "time"

type Match struct {
	ID      string
	Players []Player
}

type MatchRequest struct {
	PlayerID string
	Rating   int
	Mode     string
	Reqion   string
	JoinedAt time.Time
}

type PriorityQueueItem struct {
	Request  MatchRequest
	Priority int
	Index    int
}

type MatchPriorityQueue []*PriorityQueueItem
