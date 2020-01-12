package clock

type timePieceOld struct {
	deadline  time.Time
	fire      func() time.Duration
	mock      *mockClock
	heapIndex int
}

const removed = -1

func newTimePiece(m *mockClock, d time.Time) *timePieceOld {
	return &timePieceOld{
		deadline:  d,
		mock:      m,
		heapIndex: removed,
	}
}

func (t timePieceOld) hasStopped() bool {
	return t.heapIndex == removed
}

// pieceHeapOld implements mockTimers with a heap.
type pieceHeapOld []*timePieceOld

func newPieceHeap() *pieceHeapOld {
	res := make(pieceHeapOld, 0, 1024)
	return &res
}

func (h pieceHeapOld) Len() int { return len(h) }

func (h pieceHeapOld) Less(i, j int) bool {
	return h[i].deadline.Before(h[j].deadline)
}

func (h pieceHeapOld) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].heapIndex = i
	h[j].heapIndex = j
}

func (h *pieceHeapOld) Push(x interface{}) {
	n := len(*h)
	t := x.(*timePieceOld)
	t.heapIndex = n
	*h = append(*h, t)
}

func (h *pieceHeapOld) Pop() interface{} {
	old := *h
	n := len(old)
	t := old[n-1]
	t.heapIndex = removed
	*h = old[0 : n-1]
	return t
}

func (h *pieceHeapOld) start(t *timePieceOld) {
	heap.Push(h, t)
}

func (h *pieceHeapOld) stop(t *timePieceOld) {
	if !t.hasStopped() {
		heap.Remove(h, t.heapIndex)
	}
}

func (h *pieceHeapOld) reset(t *timePieceOld) {
	if !t.hasStopped() {
		heap.Fix(h, t.heapIndex)
	} else {
		heap.Push(h, t)
	}
}

func (h pieceHeapOld) next() *timePieceOld {
	if len(h) == 0 {
		return nil
	}
	return h[0]
}
