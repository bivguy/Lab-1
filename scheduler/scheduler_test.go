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

func createOperand(VR int) m.Operand {
	return m.Operand{
		SR:     -1,
		Active: true,
		VR:     VR,
		NU:     -1,
	}
}

var simpleTestCases = []TestCase{
	{
		description: "simple dependence garph test",
		IR: func() *list.List {
			ir := list.New()

			// loadI 8 => r1
			ir.PushBack(&m.OperationNode{
				Line:   1,
				Opcode: "loadI",
				OpOne: m.Operand{
					SR:     8,
					Active: true,
					VR:     -1,
					NU:     -1,
				},
				OpThree: createOperand(1),
			})

			// loadI 12 => r2
			ir.PushBack(&m.OperationNode{
				Line:   2,
				Opcode: "loadI",
				OpOne: m.Operand{
					SR:     12,
					Active: true,
					VR:     -1,
					NU:     -1,
				},
				OpThree: createOperand(2),
			})

			// mult r1, r2 => r3
			ir.PushBack(&m.OperationNode{
				Line:    3,
				Opcode:  "mult",
				OpOne:   createOperand(1),
				OpTwo:   createOperand(2),
				OpThree: createOperand(3),
			})

			// add r1, r3 => r4
			ir.PushBack(&m.OperationNode{
				Line:    4,
				Opcode:  "add",
				OpOne:   createOperand(1),
				OpTwo:   createOperand(3),
				OpThree: createOperand(4),
			})

			return ir
		}(),
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

	dependenceGraph := graph.CreateDependenceGraph(tc.IR)

	// go through each node and print out its edges
	for _, node := range dependenceGraph {
		t.Logf("Node (Line %d, Opcode %s):", node.Op.Line, node.Op.Opcode)

		for line, depNode := range node.Edges {
			t.Logf("  DATA Edge to Line %d (Opcode %s)", line, depNode.Op.Opcode)
		}
		for line, depNode := range node.ConflictEdges {
			t.Logf("  CONFLICT Edge to Line %d (Opcode %s)", line, depNode.Op.Opcode)
		}
		for line, depNode := range node.SerializationEdges {
			t.Logf("  SERIALIZATION Edge to Line %d (Opcode %s)", line, depNode.Op.Opcode)
		}
	}
}
