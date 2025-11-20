package scheduler

import (
	m "github.com/bivguy/Comp412/models"
)

func computePriority(graph *DependenceGraph) {
	// get the starting node as the node with the highest line number
	startNode, found := graph.DGraph[graph.maxLine]

	if !found {
		// TODO: maybe handle this as an error at some point
		return
	}

	// fmt.Println("About to do DFS")
	graph.dfs(startNode, 0, m.DATA)
	// fmt.Println("Finished DFS")
}

func (g *DependenceGraph) dfs(node *m.DependenceNode, incomingLatency int, edgeType m.EdgeType) {
	// fmt.Println("Doing DFS for line ", node.Op.Line)
	curTotalLatency := incomingLatency + computeLatency(node, edgeType)
	// check if this node has already been visited (unvisited means node latency is 0)
	if node.TotalLatency != 0 && node.TotalLatency >= curTotalLatency {
		return
	}

	// means this node is unvisited
	node.TotalLatency = curTotalLatency
	node.Latency = computeLatency(node, edgeType)

	// check if this is a leaf node if it has 0 outgoing edges
	if len(node.Edges) == 0 {
		g.leafNodes = append(g.leafNodes, node)
		return
	}

	// traverse the other nodes
	for _, edge := range node.Edges {
		nextNode := edge.To
		if nextNode == node { // defensive: no self-loop traversal
			continue
		}
		g.dfs(nextNode, node.TotalLatency, edge.Type)
	}
}
