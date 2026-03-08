package queue

import "learning/internal/models"

type MatchQueue interface {
	Enqueue(models.MatchRequest)
	Dequeue() (models.MatchRequest, bool)
	Len() int
}
