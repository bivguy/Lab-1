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
		graph: func() *DependenceGraph {
			// a=1, b=2, c=3, d=4, e=5, f=6, g=7, h=8, i=9
			a := makeNode(1, "load")  // loadAI rARP,@a => r1
			b := makeNode(2, "add")   // add r1,r1 => r2
			c := makeNode(3, "load")  // loadAI rARP,@b => r3
			d := makeNode(4, "mult")  // mult r2,r3 => r4
			e := makeNode(5, "load")  // loadAI rARP,@c => r5
			f := makeNode(6, "mult")  // mult r4,r5 => r6
			g := makeNode(7, "load")  // loadAI rARP,@d => r7
			h := makeNode(8, "mult")  // mult r6,r7 => r8
			i := makeNode(9, "store") // storeAI r8 => @a

			gph := makeGraph(a, b, c, d, e, f, g, h, i)

			// connect all the DATA edges
			gph.ConnectNodes(b, a, m.DATA) // b -> a
			gph.ConnectNodes(d, b, m.DATA) // d -> b
			gph.ConnectNodes(d, c, m.DATA) // d -> c
			gph.ConnectNodes(f, d, m.DATA) // f -> d
			gph.ConnectNodes(f, e, m.DATA) // f -> e
			gph.ConnectNodes(h, f, m.DATA) // h -> f
			gph.ConnectNodes(h, g, m.DATA) // h -> g
			gph.ConnectNodes(i, h, m.DATA) // i -> h

			// connect all the SERIALIZATION edges
			gph.ConnectNodes(i, a, m.SERIALIZATION)
			gph.ConnectNodes(i, c, m.SERIALIZATION)
			gph.ConnectNodes(i, e, m.SERIALIZATION)
			gph.ConnectNodes(i, g, m.SERIALIZATION)

			return gph
		}(),
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
