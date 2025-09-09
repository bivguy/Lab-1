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
	{
		description: "simple LOADI test",
		input:       "parser_tests/simple_tests/test_2.txt",
		expectedIR: func() *list.List {
			ir := list.New()

			ir.PushBack(m.OperationNode{
				Line:    1,
				Opcode:  "loadI",
				OpOne:   m.Operand{SR: 17},
				OpThree: m.Operand{SR: 1},
			})

			return ir
		}(), expectedError: false,
	},
	{
		description: "simple ARITHOP add test",
		input:       "parser_tests/simple_tests/test_3.txt",
		expectedIR: func() *list.List {
			ir := list.New()

			ir.PushBack(m.OperationNode{
				Line:    1,
				Opcode:  "add",
				OpOne:   m.Operand{SR: 1},
				OpTwo:   m.Operand{SR: 1},
				OpThree: m.Operand{SR: 2},
			})

			return ir
		}(), expectedError: false,
	},
	{
		description: "simple OUTPUT test",
		input:       "parser_tests/simple_tests/test_4.txt",
		expectedIR: func() *list.List {
			ir := list.New()

			ir.PushBack(m.OperationNode{
				Line:    1,
				Opcode:  "output",
				OpThree: m.Operand{SR: 2214},
			})

			return ir
		}(), expectedError: false,
	},
	{
		description: "simple NOP test",
		input:       "parser_tests/simple_tests/test_5.txt",
		expectedIR: func() *list.List {
			ir := list.New()

			ir.PushBack(m.OperationNode{
				Line:   1,
				Opcode: "nop",
			})

			return ir
		}(), expectedError: false,
	},
}
