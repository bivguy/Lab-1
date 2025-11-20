package scheduler

import (
	"container/heap"
	"container/list"

	m "github.com/bivguy/Comp412/models"
)

const (
	NOT_READY m.Status = iota // 0
	READY                     // 1
	ACTIVE                    // 2
	RETIRED                   // 3
)

type scheduler struct {
	IR *list.List
}

func NewSchedule(IR *list.List) *scheduler {
	return &scheduler{IR: IR}
}

func (s *scheduler) Schedule() {
	// create a dependence graph
	graph := New()
	graph.CreateDependenceGraph(s.IR)

	// compute the priorities
	computePriority(graph)

	// ready list represents the leaf nodes initially
	var ready DNodeHeap
	heap.Init(&ready)
	for _, leafNode := range graph.leafNodes {
		heap.Push(&ready, leafNode)
		leafNode.Status = READY
	}

	active := make(map[int][]*m.DependenceNode)
	cycle := 1

	for len(ready) > 0 || len(active) > 0 {
		if len(ready) > 1 {
			dn := heap.Pop(&ready).(*m.DependenceNode)
			// pick an operation from each functional unit
			removeIndex := cycle + dn.Latency

			// check if we need to make a new slice here
			_, ok := active[removeIndex]
			if !ok {
				active[removeIndex] = []*m.DependenceNode{}
			}

			// add this to the active list corresponding to when it is removed
			active[removeIndex] = append(active[removeIndex], dn)
			dn.Status = ACTIVE
		}

		cycle += 1

		retiredOps := active[cycle]
		for _, retiredOp := range retiredOps {
			retiredOp.Status = RETIRED
		}

		// remove the active ops
		delete(active, cycle)
		// find each op in the active list that retires
		for _, retiredOp := range retiredOps {
			// check for each op that that relies on this retired op
			for _, d := range retiredOp.ReverseEdges {
				// TODO: check if this applies to all edges or just specific edges
				// skip nodes already added to ready
				if d.To.Status != NOT_READY {
					continue
				}
				// TODO: if a node needs multiple VRs to be ready, we might need a more thorough check
				heap.Push(&ready, d.To)
				d.To.Status = READY
			}
		}

	}
}
