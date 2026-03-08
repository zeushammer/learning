package queue

import (
	"hash/fnv"
	"learning/internal/models"
)

type ShardedQueue struct {
	shards []MatchQueue
}

func NewShardedQueue(shardCount int) *ShardedQueue {

	shards := make([]MatchQueue, shardCount)

	for i := range shards {
		shards[i] = NewPriorityQueue()
	}

	return &ShardedQueue{
		shards: shards,
	}
}

func (sq *ShardedQueue) shardFor(playerID string) int {
	h := fnv.New32a()
	h.Write([]byte(playerID))

	return int(h.Sum32()) % len(sq.shards)
}

func (sq *ShardedQueue) Enqueue(req models.MatchRequest) {

	idx := sq.shardFor(req.PlayerID)

	sq.shards[idx].Enqueue(req)
}

func (sq *ShardedQueue) Dequeue() (models.MatchRequest, bool) {

	for _, shard := range sq.shards {

		if req, ok := shard.Dequeue(); ok {
			return req, true
		}
	}

	return models.MatchRequest{}, false
}

func (sq *ShardedQueue) Len() int {

	total := 0

	for _, shard := range sq.shards {
		total += shard.Len()
	}

	return total
}
