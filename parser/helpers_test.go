package parser

import (
	"container/list"

	m "github.com/bivguy/Comp412/models"
)

type TestCase struct {
	description   string
	input         string
	expectedIR    *list.List
	expectedError bool
}

var simpleTestCases = []TestCase{
	{
		description: "simple MEMOP test",
		input:       "parser_tests/simple_tests/test_1.txt",
		expectedIR: func() *list.List {
			ir := list.New()

			ir.PushBack(m.OperationNode{
				Line:    1,
				Opcode:  "store",
				OpOne:   m.Operand{SR: 1},
				OpThree: m.Operand{SR: 2},
			})

			return ir
		}(), expectedError: false,
	},
}
