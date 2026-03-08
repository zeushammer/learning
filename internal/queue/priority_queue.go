package queue

import (
	"container/heap"
	"sync"

	"learning/internal/models"
)

type pqItem struct {
	req      models.MatchRequest
	priority int
	index    int
}

type priorityHeap []*pqItem

func (h priorityHeap) Len() int { return len(h) }

func (h priorityHeap) Less(i, j int) bool {
	return h[i].priority > h[j].priority
}

func (h priorityHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *priorityHeap) Push(x interface{}) {
	item := x.(*pqItem)
	item.index = len(*h)
	*h = append(*h, item)
}

func (h *priorityHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return item
}

type PriorityQueue struct {
	h  priorityHeap
	mu sync.Mutex
}

func NewPriorityQueue() *PriorityQueue {
	h := make(priorityHeap, 0)
	heap.Init(&h)

	return &PriorityQueue{
		h: h,
	}
}

func (pq *PriorityQueue) Enqueue(req models.MatchRequest) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	item := &pqItem{
		req:      req,
		priority: req.Rating,
	}

	heap.Push(&pq.h, item)
}

func (pq *PriorityQueue) Dequeue() (models.MatchRequest, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if pq.h.Len() == 0 {
		return models.MatchRequest{}, false
	}

	item := heap.Pop(&pq.h).(*pqItem)
	return item.req, true
}

func (pq *PriorityQueue) Len() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	return pq.h.Len()
}
