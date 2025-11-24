package scheduler

import (
	"fmt"

	m "github.com/bivguy/Comp412/models"
)

func computePriority(graph *DependenceGraph) {
	// get the starting node as the node with the highest line number
	startNode, found := graph.DGraph[graph.maxLine]

	if !found {
		// TODO: maybe handle this as an error at some point
		return
	}
	seen := make(map[int]bool)
	if DEBUG_PRIORITY_COMPUTATION {
		fmt.Println("About to do DFS on node of opCode ", startNode.Op.Opcode)
	}
	graph.dfs(startNode, 0, m.DATA, seen)
}

func (g *DependenceGraph) dfs(node *m.DependenceNode, incomingLatency int, edgeType m.EdgeType, seen map[int]bool) {
	curLine := node.Op.Line
	// skip if already seen
	if seen[curLine] {
		if DEBUG_PRIORITY_COMPUTATION {
			fmt.Println("alrady seen line ", curLine)
		}
		return
	}

	seen[curLine] = true
	if DEBUG_PRIORITY_COMPUTATION {
		fmt.Println("doing DFS on line ", curLine, " with", len(node.Edges), " neighbors")
	}
	// fmt.Println("Doing DFS for line ", node.Op.Line)
	curTotalLatency := incomingLatency + computeLatency(node, edgeType)
	// check if this node has already been visited (unvisited means node latency is 0)
	if node.TotalLatency != 0 && node.TotalLatency >= curTotalLatency {
		if DEBUG_PRIORITY_COMPUTATION {
			fmt.Println("already visisted line ", node.Op.Line)
		}
		return
	}

	// means this node is unvisited
	node.TotalLatency = curTotalLatency
	node.Latency = computeLatency(node, edgeType)

	// check if this is a leaf node if it has 0 outgoing edges
	// TOOD: this reverse edge check might be wrong
	if len(node.Edges) == 0 {
		g.leafNodes = append(g.leafNodes, node)
		if DEBUG_PRIORITY_COMPUTATION {
			fmt.Println("adding the leaf node of line ", node.Op.Line)
		}
	}

	// traverse the other nodes
	for nextNodeLine, edge := range node.Edges {
		nextNode := edge.To
		if DEBUG_PRIORITY_COMPUTATION {
			fmt.Println("about to go to node of line ", nextNodeLine)
		}

		if nextNodeLine == node.Op.Line { // no self-loop traversal
			continue
		}

		g.dfs(nextNode, node.TotalLatency, edge.Type, seen)
	}
}
