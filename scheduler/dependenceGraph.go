package scheduler

import (
	"container/list"

	m "github.com/bivguy/Comp412/models"
)

type DependenceGraph struct {
	graph map[int]*m.DependenceNode
}

type EdgeType int

const (
	DATA          EdgeType = iota // 0
	CONFLICT                      // 1
	SERIALIZATION                 // 2
)

// type DependenceNode struct {
// 	op                 *m.OperationNode
// 	edges              map[int]*DependenceNode
// 	conflictEdges      map[int]*DependenceNode
// 	serializationEdges map[int]*DependenceNode
// }

func NewDependenceNode(op *m.OperationNode) *m.DependenceNode {
	return &m.DependenceNode{
		Op:                 op,
		Edges:              make(map[int]*m.DependenceNode),
		ConflictEdges:      make(map[int]*m.DependenceNode),
		SerializationEdges: make(map[int]*m.DependenceNode),
	}
}

func New() *DependenceGraph {
	graph := make(map[int]*m.DependenceNode)

	return &DependenceGraph{
		graph: graph,
	}
}

func (g *DependenceGraph) ConnectNodes(in *m.DependenceNode, out *m.DependenceNode, edgeType EdgeType) {
	line := in.Op.Line
	switch edgeType {
	case DATA:
		// connect the edge where the node is defined to the node where it is used (definition -> use) by mapping the line number to the node
		out.Edges[line] = in
	case CONFLICT:
		out.ConflictEdges[line] = in
	case SERIALIZATION:
		out.SerializationEdges[line] = in
	}
}

func (g *DependenceGraph) CreateDependenceGraph(IR *list.List) map[int]*m.DependenceNode {
	var mostRecentStore *m.DependenceNode
	var mostRecentOutput *m.DependenceNode
	var previousReads []*m.DependenceNode

	for node := IR.Front(); node != nil; node = node.Next() {
		op := node.Value.(*m.OperationNode)
		opCode := op.Opcode

		if opCode == "nop" {
			continue
		}

		// create a node for this operation
		node := NewDependenceNode(op)

		operandList := []*m.Operand{&op.OpOne, &op.OpTwo, &op.OpThree}

		// go through each operand that's defined and add the node to the graph if there is a definition
		for i, o := range operandList {
			// skip if its not active or if its not definiition
			if !o.Active || !isRegister(opCode, i) || !isDefinition(opCode, i) {
				continue
			}

			g.graph[o.VR] = node
		}

		// go through each use and add edges from uses to their definitions
		for i, o := range operandList {
			// skip if its not active or if its a definiition
			if !o.Active || !isRegister(opCode, i) || isDefinition(opCode, i) {
				continue
			}

			// add an edge from the definition node to this use node
			if defNode, exists := g.graph[o.VR]; exists {
				g.ConnectNodes(defNode, node, DATA)
			}
		}

		// load & output need an edge to the most recent store
		if opCode == "load" || opCode == "output" {
			if mostRecentStore != nil {
				g.ConnectNodes(mostRecentStore, node, CONFLICT)
			}
		}

		// output needs an edge to the most recent output
		if opCode == "output" {
			if mostRecentOutput != nil {
				g.ConnectNodes(mostRecentOutput, node, SERIALIZATION)
			}
		}

		// store needs an edge to the most recent store, as well as each previous load & output
		if opCode == "store" {
			if mostRecentStore != nil {
				g.ConnectNodes(mostRecentStore, node, SERIALIZATION)
				// connect this store to all previous reads (store -> read)
				for _, readNode := range previousReads {
					g.ConnectNodes(readNode, node, SERIALIZATION)
				}

				// reset the previous reads
				previousReads = []*m.DependenceNode{}
			}
		}

		// update the recent store/output/reads
		switch opCode {
		case "store":
			mostRecentStore = node
		case "output":
			mostRecentOutput = node
			previousReads = append(previousReads, node)
		case "load":
			previousReads = append(previousReads, node)
		}
	}

	return g.graph
}

func isRegister(opcode string, i int) bool {
	if (opcode == "loadI" && i == 0) || (opcode == "output" && i == 2) {
		return false
	}

	return true
}

func isDefinition(opcode string, i int) bool {
	if opcode == "store" {
		return false
	} else if i == 2 {
		return true
	}
	return false
}
