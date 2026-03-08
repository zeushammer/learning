package queue

// todo: implement a priority queue and a generic queue

// Use for:
// Queue[Player]
// Queue[MatchRequest]

type Queue[T any] struct {
	items []T
}

// Push adds an item to the `back` of the line
func (q *Queue[T]) Push(v T) {
	q.items = append(q.items, v)
}

// Pop removes the item from the `front` of the line (First-In, First-Out or FIFO)
func (q *Queue[T]) Pop() (T, bool) {
	// empty slice
	if len(q.items) == 0 {
		// Since `T` could be an
		// `int` (zero value `0`),
		// a `string` (zero value ""),
		// or a pointer (zero value `nil`)
		// declaring `var zero T` might be the idiomatic way to return the "nothing" value
		var zero T
		return zero, false
	}

	v := q.items[0]
	// slicing a slice, only moves the pointer forward.
	// If T is a large struct or a pointer, the first element
	// will stay in memory because the underlying array still references it.
	var zero T
	q.items[0] = zero
	// zero'ing it out, signals to the Garbage collector the reference is no longer

	q.items = q.items[1:]
	return v, true
}
