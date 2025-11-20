package scheduler

import (
	"container/list"
	"fmt"

	m "github.com/bivguy/Comp412/models"
)

type DependenceGraph struct {
	graph           map[int]*m.DependenceNode
	DGraph          map[int]*m.DependenceNode
	leafNodes       []*m.DependenceNode
	maxLine         int
	maxTotalLatency int
}

func NewDependenceNode(op *m.OperationNode) *m.DependenceNode {
	return &m.DependenceNode{
		Op:           op,
		Edges:        make(map[int]*m.DependenceEdge),
		ReverseEdges: make(map[int]*m.DependenceEdge),
		Status:       NOT_READY,
	}
}

func New() *DependenceGraph {
	graph := make(map[int]*m.DependenceNode)
	DGraph := make(map[int]*m.DependenceNode)

	return &DependenceGraph{
		graph:  graph,
		DGraph: DGraph,
	}
}

func computeLatency(node *m.DependenceNode, edgeType m.EdgeType) int {
	if edgeType == m.SERIALIZATION {
		return 1
	}

	switch node.Op.Opcode {
	case "load", "store":
		return 6
	case "mult":
		return 3
	case "add", "sub", "lshift", "rshift", "loadI", "output", "nop":
		return 1
	default:
		return -1
	}
}

// TODO: add in the reverse edges
func (g *DependenceGraph) ConnectNodes(in *m.DependenceNode, out *m.DependenceNode, edgeType m.EdgeType) {
	// line := in.Op.Line
	edge := &m.DependenceEdge{
		To:      out,
		Type:    edgeType,
		Latency: computeLatency(in, edgeType),
	}

	// connect the edge where the node is defined to the node where it is used (definition -> use) by mapping the line number to the node
	in.Edges[out.Op.Line] = edge

	reverseEdge := &m.DependenceEdge{
		To:      in,
		Type:    edgeType,
		Latency: computeLatency(out, edgeType),
	}

	// connect this to the opposite node as well
	out.ReverseEdges[in.Op.Line] = reverseEdge

	// fmt.Printf("line %d has edge to line %d \n", in.Op.Line, out.Op.Line)
}

func (g *DependenceGraph) CreateDependenceGraph(IR *list.List) map[int]*m.DependenceNode {
	var mostRecentStore *m.DependenceNode
	var mostRecentOutput *m.DependenceNode
	var previousReads []*m.DependenceNode

	var line int
	for node := IR.Front(); node != nil; node = node.Next() {
		op := node.Value.(*m.OperationNode)
		opCode := op.Opcode

		if opCode == "nop" {
			continue
		}

		line = op.Line

		// create a node for this operation
		node := NewDependenceNode(op)

		operandList := []*m.Operand{&op.OpOne, &op.OpTwo, &op.OpThree}

		g.DGraph[line] = node

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
				g.ConnectNodes(defNode, node, m.DATA)
			}
		}

		// load & output need an edge to the most recent store
		if opCode == "load" || opCode == "output" {
			if mostRecentStore != nil {
				g.ConnectNodes(mostRecentStore, node, m.CONFLICT)
			}
		}

		// output needs an edge to the most recent output
		if opCode == "output" {
			if mostRecentOutput != nil {
				g.ConnectNodes(mostRecentOutput, node, m.SERIALIZATION)
			}
		}

		// store needs an edge to the most recent store, as well as each previous load & output
		if opCode == "store" {
			if mostRecentStore != nil {
				fmt.Println("Connecting store at line", line, "to most recent store at line", mostRecentStore.Op.Line)
				g.ConnectNodes(mostRecentStore, node, m.SERIALIZATION)
			}

			// connect this store to all previous reads (store -> read)
			// TODO: optimize this later to remove redundant serialization edges
			for _, readNode := range previousReads {
				fmt.Println("Connecting store at line", line, "to previous read at line", readNode.Op.Line)
				g.ConnectNodes(readNode, node, m.SERIALIZATION)
			}

			// reset the previous reads
			previousReads = []*m.DependenceNode{}
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

	g.maxLine = line

	return g.DGraph
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
