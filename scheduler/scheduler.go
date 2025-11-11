package scheduler

import "container/list"

type scheduler struct {
	IR *list.List
}

func NewSchedule(IR *list.List) *scheduler {
	return &scheduler{}
}

func (s *scheduler) Schedule() {
	// create a dependence graph
	graph := New()
	graph.CreateDependenceGraph(s.IR)

	// compute the priorities
	computePriority(graph)

}
