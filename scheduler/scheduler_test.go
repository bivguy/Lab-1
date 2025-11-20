package scheduler

import (
	"container/list"
	"testing"

	m "github.com/bivguy/Comp412/models"
)

type TestCase struct {
	description string
	IR          *list.List
}

type PriorityTestCase struct {
	description string
	graph       *DependenceGraph
	expected    map[int]int
}

type SchedulerTestCase struct {
	description string
	graph       *DependenceGraph
	expected    *list.List
}

func createOperand(VR int) m.Operand {
	return m.Operand{
		SR:     -1,
		Active: true,
		VR:     VR,
		NU:     -1,
	}
}

func makeNode(line int, opcode string) *m.DependenceNode {
	return &m.DependenceNode{
		Op:    &m.OperationNode{Line: line, Opcode: opcode},
		Edges: make(map[int]*m.DependenceEdge),
	}
}

func makeGraph(nodes ...*m.DependenceNode) *DependenceGraph {
	g := New()
	for _, n := range nodes {
		g.DGraph[n.Op.Line] = n
		if n.Op.Line > g.maxLine {
			g.maxLine = n.Op.Line
		}
	}
	return g
}

var simpleTestCases = []TestCase{
	// {
	// 	description: "simple dependence garph test",
	// 	IR: func() *list.List {
	// 		ir := list.New()

	// 		// loadI 8 => r1
	// 		ir.PushBack(&m.OperationNode{
	// 			Line:   1,
	// 			Opcode: "loadI",
	// 			OpOne: m.Operand{
	// 				SR:     8,
	// 				Active: true,
	// 				VR:     -1,
	// 				NU:     -1,
	// 			},
	// 			OpThree: createOperand(1),
	// 		})

	// 		// loadI 12 => r2
	// 		ir.PushBack(&m.OperationNode{
	// 			Line:   2,
	// 			Opcode: "loadI",
	// 			OpOne: m.Operand{
	// 				SR:     12,
	// 				Active: true,
	// 				VR:     -1,
	// 				NU:     -1,
	// 			},
	// 			OpThree: createOperand(2),
	// 		})

	// 		// mult r1, r2 => r3
	// 		ir.PushBack(&m.OperationNode{
	// 			Line:    3,
	// 			Opcode:  "mult",
	// 			OpOne:   createOperand(1),
	// 			OpTwo:   createOperand(2),
	// 			OpThree: createOperand(3),
	// 		})

	// 		// add r1, r3 => r4
	// 		ir.PushBack(&m.OperationNode{
	// 			Line:    4,
	// 			Opcode:  "add",
	// 			OpOne:   createOperand(1),
	// 			OpTwo:   createOperand(3),
	// 			OpThree: createOperand(4),
	// 		})

	// 		return ir
	// 	}(),
	// },
	{
		description: "memory + output dependence graph test",
		IR: func() *list.List {
			ir := list.New()

			// 1: loadI 8 => r3
			ir.PushBack(&m.OperationNode{
				Line:    1,
				Opcode:  "loadI",
				OpOne:   m.Operand{SR: 8, Active: true, VR: -1, NU: -1},
				OpThree: createOperand(3),
			})

			// 2: loadI 12 => r4
			ir.PushBack(&m.OperationNode{
				Line:    2,
				Opcode:  "loadI",
				OpOne:   m.Operand{SR: 12, Active: true, VR: -1, NU: -1},
				OpThree: createOperand(4),
			})

			// 3: add r3, r4 => r0
			ir.PushBack(&m.OperationNode{
				Line:    3,
				Opcode:  "add",
				OpOne:   createOperand(3),
				OpTwo:   createOperand(4),
				OpThree: createOperand(0),
			})

			// 4: load r0 => r1
			ir.PushBack(&m.OperationNode{
				Line:    4,
				Opcode:  "load",
				OpOne:   createOperand(0),
				OpThree: createOperand(1),
			})

			// 5: load r3 => r2
			ir.PushBack(&m.OperationNode{
				Line:    5,
				Opcode:  "load",
				OpOne:   createOperand(3),
				OpThree: createOperand(2),
			})

			// 6: store r1 => r0
			ir.PushBack(&m.OperationNode{
				Line:    6,
				Opcode:  "store",
				OpOne:   createOperand(1), // value
				OpThree: createOperand(0), // address
			})

			// 7: output 12
			ir.PushBack(&m.OperationNode{
				Line:   7,
				Opcode: "output",
				OpOne:  m.Operand{SR: 12, Active: true, VR: -1, NU: -1},
			})

			return ir
		}(),
	},
}

var simplePriorityTestCases = []PriorityTestCase{
	{
		description: "simple test case 1",
		graph:       DEPENDENCE_GRAPHS[0],
		expected: map[int]int{
			1: 22,
			2: 16,
			3: 21,
			4: 15,
			5: 18,
			6: 12,
			7: 15,
			8: 9,
			9: 6,
		},
	},
}

var simpleSchedulerTestCases = []SchedulerTestCase{
	{
		description: "simple scheduler test case",
	},
}

func TestSimpleDependenceGraph(t *testing.T) {
	for _, tc := range simpleTestCases {
		t.Run(tc.description, func(t *testing.T) {
			runTest(tc, t)
		})
	}
}

func runTest(tc TestCase, t *testing.T) {
	graph := New()

	graph.CreateDependenceGraph(tc.IR)

	// go through each node and print out its edges
	for _, node := range graph.DGraph {
		t.Logf("Node (Line %d, Opcode %s):", node.Op.Line, node.Op.Opcode)

		for line, edge := range node.Edges {
			t.Logf("%s Edge to Line %d (Latency %d)", edge.Type.String(), line, edge.Latency)
		}

	}
}

func TestComputePriority(t *testing.T) {
	for _, tc := range simplePriorityTestCases {
		t.Run(tc.description, func(t *testing.T) {
			runPriorityTest(tc, t)
		})
	}

}

func runPriorityTest(tc PriorityTestCase, t *testing.T) {
	computePriority(tc.graph)

	// print out each node's latency and line number// Print all nodes and their computed TotalLatency
	for line, n := range tc.graph.DGraph {
		t.Logf("Line %2d | Opcode %-6s | TotalLatency = %d", line, n.Op.Opcode, n.TotalLatency)
	}

	for l, expectedLatency := range tc.expected {
		actualLatency := tc.graph.DGraph[l].TotalLatency

		if expectedLatency != actualLatency {
			t.Errorf("Mismatch: On line %d, expected latency of %d but got %d", l, expectedLatency, actualLatency)
		}

	}

}
