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

type PerformanceTestCase struct {
	description   string
	input         string
	MaximumTimeMs int64
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
	{
		description: "testing some invalid file",
		input:       "parser_tests/simple_tests/test_6.txt",
		expectedIR: func() *list.List {
			ir := list.New()

			ir.PushBack(m.OperationNode{
				Line:    3,
				Opcode:  "output",
				OpThree: m.Operand{SR: 4},
			})

			return ir
		}(), expectedError: true,
	},
}

var complexTestCases = []TestCase{
	{
		description: "t11",
		input:       "parser_tests/complex_tests/t11.i.txt",
		expectedIR: func() *list.List {
			ir := list.New()

			//   loadI 8 => r1
			ir.PushBack(m.OperationNode{
				Line:    2,
				Opcode:  "loadI",
				OpOne:   m.Operand{SR: 8},
				OpThree: m.Operand{SR: 1},
			})

			//   store r1 => r1
			ir.PushBack(m.OperationNode{
				Line:    3,
				Opcode:  "store",
				OpOne:   m.Operand{SR: 1},
				OpThree: m.Operand{SR: 1},
			})

			//   store r1 => r1
			ir.PushBack(m.OperationNode{
				Line:    4,
				Opcode:  "output",
				OpThree: m.Operand{SR: 4},
			})

			return ir
		}(), expectedError: false,
	},
}

var performanceTestCases = []PerformanceTestCase{
	{
		description:   "Performance Test 1",
		input:         "parser_tests/complex_tests/t128k.i.txt",
		MaximumTimeMs: 200,
	},
}
