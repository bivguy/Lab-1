package scheduler

import (
	"container/heap"
	"container/list"
	"fmt"
	"strings"

	m "github.com/bivguy/Comp412/models"
)

const DEBUG_DEPENDENCE_GRAPH = false
const DEBUG_PRIORITY_COMPUTATION = false
const DEBUG_SCHEDULING = false

const (
	NOT_READY m.Status = iota // 0
	READY                     // 1
	ACTIVE                    // 2
	RETIRED                   // 3
)

type operationBlock struct {
	operationOne *m.DependenceNode
	operationTwo *m.DependenceNode
}

type scheduler struct {
	IR *list.List
}

func NewSchedule(IR *list.List) *scheduler {
	return &scheduler{IR: IR}
}

func (s *scheduler) Schedule() []*operationBlock {
	var schedule []*operationBlock
	// create a dependence graph
	graph := New()
	graph.CreateDependenceGraph(s.IR)

	computePriority(graph)

	if DEBUG_DEPENDENCE_GRAPH {
		// go through each node and print out its edges
		for _, node := range graph.DGraph {
			fmt.Println("Node (Line ", node.Op.Line, ", , Opcode ", node.Op.Opcode, "): with priority value ", node.TotalLatency)

			for line, edge := range node.Edges {
				fmt.Println(edge.Type.String(), " Edge to Line ", line, " (Latency ", edge.Latency, ")")
			}
		}
	}
	var ready DNodeHeap
	heap.Init(&ready)
	for _, leafNode := range graph.leafNodes {
		heap.Push(&ready, leafNode)
		leafNode.Status = READY
	}

	// create a NOP operation to be used
	op := &m.OperationNode{
		Opcode: "nop",
	}
	notOp := &m.DependenceNode{
		Op:           op,
		Edges:        make(map[int]*m.DependenceEdge),
		ReverseEdges: make(map[int]*m.DependenceEdge),
		Status:       RETIRED,
	}

	// TODO: make a list of how many of each node's nighbors we've retired so its faster
	active := make(map[int][]*m.DependenceNode)
	cycle := 1

	for len(ready) > 0 || len(active) > 0 {
		if DEBUG_SCHEDULING {
			fmt.Println("Length of ready is ", len(ready), " and length of active is ", len(active))
			if len(active) == 1 && len(ready) == 0 {
				for k, v := range active {
					fmt.Println("active key is ", k, " with value length ", len(v), " and opcode ", v[0].Op.Opcode, " at line ", v[0].Op.Line)
				}
			}
		}
		skipped := []*m.DependenceNode{}
		opBlock := &operationBlock{
			operationOne: notOp,
			operationTwo: notOp,
		}
		numIssued := 0

		for numIssued < 2 && len(ready) > 0 {
			dn := heap.Pop(&ready).(*m.DependenceNode)

			// check where this should go in the opBlock
			opCode := dn.Op.Opcode
			switch opCode {
			// only one output is allowed per cycle (in either slot)
			case "output":
				// check if either slot is taken
				if opBlock.operationOne.Op.Opcode != "nop" && opBlock.operationTwo.Op.Opcode != "nop" {
					// both slots are taken, skip this one for now
					skipped = append(skipped, dn)
					continue
				}
				opBlock.operationOne = dn
				numIssued = 2 // force an exit
			// mult can only go in slot two
			case "mult":
				if opBlock.operationTwo.Op.Opcode != "nop" {
					skipped = append(skipped, dn)
					continue
				}
				opBlock.operationTwo = dn
				numIssued += 1
			// load or store can only go in slot one
			case "load", "store":
				if opBlock.operationOne.Op.Opcode != "nop" {
					skipped = append(skipped, dn)
					continue
				}
				opBlock.operationOne = dn
				numIssued += 1
			// all other operations can go in either slot
			default:
				if opBlock.operationTwo.Op.Opcode == "nop" {
					opBlock.operationTwo = dn
					numIssued += 1
				} else if opBlock.operationOne.Op.Opcode == "nop" {
					opBlock.operationOne = dn
					numIssued += 1
				} else {
					fmt.Println("ERROR: both slots taken when they shouldn't be")
					numIssued = 2
					continue
				}
			}

			// pick an operation from each functional unit
			removeIndex := cycle + computeLatency(dn, m.DATA)
			if DEBUG_SCHEDULING {
				if removeIndex <= cycle {
					fmt.Println("ERROR: remove index ", removeIndex, " is less than or equal to current cycle ", cycle)
					fmt.Println("with op ", dn.Op.Opcode, " at line ", dn.Op.Line, " with latency ", dn.Latency)
				}
			}
			// check if we need to make a new slice here
			_, ok := active[removeIndex]
			if !ok {
				active[removeIndex] = []*m.DependenceNode{}
			}

			// add this to the active list corresponding to when it is removed
			active[removeIndex] = append(active[removeIndex], dn)
			dn.Status = ACTIVE
			if DEBUG_SCHEDULING {
				fmt.Println("at the ready node ", dn.Op, " at line ", dn.Op.Line)
				fmt.Println("we want to remove at cycle ", removeIndex)
			}

		}

		schedule = append(schedule, opBlock)
		// push back the skipped nodes
		for _, skippedNode := range skipped {
			heap.Push(&ready, skippedNode)
		}
		cycle += 1
		// if cycle > 50 {
		// 	return nil
		// }
		if DEBUG_SCHEDULING {
			fmt.Println("cycle is now", cycle)
		}
		retiredOps := active[cycle]
		for _, retiredOp := range retiredOps {
			retiredOp.Status = RETIRED
		}

		if len(retiredOps) > 0 {
			if DEBUG_SCHEDULING {
				fmt.Println("deleted the ops of length ", len(retiredOps))
			}
		}

		// remove the active ops
		delete(active, cycle)

		// find each op in the active list that retires
		for _, retiredOp := range retiredOps {
			// skip if it's a NOP
			if retiredOp.Op.Opcode == "nop" {
				continue
			}
			if DEBUG_SCHEDULING {
				fmt.Println("about to retired the op ", retiredOp.Op.Opcode, " of line ", retiredOp.Op.Line)
			}
			// check for each op that that relies on this retired op
			for _, d := range retiredOp.ReverseEdges {
				// skip nodes already added to ready
				if d.To.Status != NOT_READY {
					continue
				}
				addNeighbor := true
				// TODO: if a node needs multiple VRs to be ready, we might need a more thorough check
				// check if all of the dependences for this outgoing node are retited
				for _, node := range d.To.Edges {
					if node.To.Status != RETIRED {
						addNeighbor = false
						break
					}
				}

				if addNeighbor {
					heap.Push(&ready, d.To)
					d.To.Status = READY
				}
			}
		}

		// go through each multicycle op in active and see if there are any early releases we can do
		// TODO: this part can be optimized
		// for _, ops := range active {
		// 	for _, dn := range ops {
		// 		// examine each load and store
		// 		opCode := dn.Op.Opcode
		// 		if opCode == "load" || opCode == "store" {
		// 			// check ops that depend on this current operation an early release
		// 			for _, edge := range dn.ReverseEdges {
		// 				// skip all nodes that aren't connected by a serial edge
		// 				if edge.Type != m.SERIALIZATION {
		// 					continue
		// 				}
		// 				// an early release candidate
		// 				candidate := edge.To
		// 				// now we must check that all of this node's other dependencies are also satisfied
		// 				valid := true
		// 				for _, candidateEdge := range candidate.Edges {
		// 					// if the node connected via serial edge is not ready or if the node connected by other edge is not retired, we can't do an early release
		// 					if (candidateEdge.Type == m.SERIALIZATION && candidateEdge.To.Status == NOT_READY) ||
		// 						(candidateEdge.Type != m.SERIALIZATION && candidateEdge.To.Status != RETIRED) {
		// 						valid = false
		// 						break
		// 					}
		// 				}

		// 				// this node can be set as ready via early release
		// 				if valid {
		// 					heap.Push(&ready, candidate)
		// 					candidate.Status = READY
		// 				}
		// 			}
		// 		}

		// 	}
		// }
	}

	if DEBUG_SCHEDULING {
		fmt.Println("length of scheduler: ", len(schedule))
	}
	return schedule
}

func (s *scheduler) PrintSchedule() {
	scheduledBlocks := s.Schedule()
	var b strings.Builder

	for _, block := range scheduledBlocks {
		// build each block

		opString := "[ " + block.operationOne.Op.String() + " ; " + block.operationTwo.Op.String() + " ]\n"

		fmt.Fprintf(&b, opString)
	}

	fmt.Println(b.String())
}
