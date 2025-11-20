package scheduler

import (
	m "github.com/bivguy/Comp412/models"
)

type DNodeHeap []*m.DependenceNode

func (h DNodeHeap) Len() int { return len(h) }

// implementing a max heap
func (h DNodeHeap) Less(i, j int) bool { return h[i].TotalLatency > h[j].TotalLatency }

func (h DNodeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *DNodeHeap) Push(x any) {
	// Push requires a pointer receiver because it modifies the underlying slice.
	node := x.(*m.DependenceNode)
	*h = append(*h, node)
}

func (h *DNodeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
