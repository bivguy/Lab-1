package scanner

import (
	"os"
	"testing"
	"time"

	. "github.com/bivguy/Comp412/constants"
	"github.com/bivguy/Comp412/models"
)

type TestCase struct {
	description    string
	input          string
	expectedTokens []models.Token
	expectedError  bool
}

type PerformanceTestCase struct {
	description   string
	input         string
	MaximumTimeMs int64
}

var simpleTestCases = []TestCase{
	{
		description: "Valid scanner input with load",
		input:       "scanner_tests/simple_tests/scanner_test1.txt",
		expectedTokens: []models.Token{
			{Category: MEMOP, Lexeme: "load"},
			{Category: EOL, Lexeme: "\n"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with lshift and load",
		input:       "scanner_tests/simple_tests/scanner_test2.txt",
		expectedTokens: []models.Token{
			{Category: ARITHOP, Lexeme: "lshift"},
			{Category: EOL, Lexeme: "\n"},
			{Category: ARITHOP, Lexeme: "lshift"},
			{Category: MEMOP, Lexeme: "load"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with lshift, load, and loadI",
		input:       "scanner_tests/simple_tests/scanner_test3.txt",
		expectedTokens: []models.Token{
			{Category: ARITHOP, Lexeme: "lshift"},
			{Category: MEMOP, Lexeme: "load"},
			{Category: EOL, Lexeme: "\n"},
			{Category: LOADI, Lexeme: "loadI"},
			{Category: MEMOP, Lexeme: "load"},
			{Category: LOADI, Lexeme: "loadI"},
			{Category: EOL, Lexeme: "\n"},
			{Category: EOL, Lexeme: "\n"},
			{Category: ARITHOP, Lexeme: "lshift"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with store, constant, and sub",
		input:       "scanner_tests/simple_tests/scanner_test4.txt",
		expectedTokens: []models.Token{
			{Category: MEMOP, Lexeme: "store"},
			{Category: CONSTANT, Lexeme: "132"},
			{Category: EOL, Lexeme: "\n"},
			{Category: ARITHOP, Lexeme: "sub"},
			{Category: CONSTANT, Lexeme: "142"},
			{Category: CONSTANT, Lexeme: "2"},
			{Category: EOL, Lexeme: "\n"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with add, comma, register, and into",
		input:       "scanner_tests/simple_tests/scanner_test5.txt",
		expectedTokens: []models.Token{
			{Category: ARITHOP, Lexeme: "add"},
			{Category: REGISTER, Lexeme: "r0"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with store register, into, output, and constant",
		input:       "scanner_tests/simple_tests/scanner_test6.txt",
		expectedTokens: []models.Token{
			{Category: MEMOP, Lexeme: "store"},
			{Category: REGISTER, Lexeme: "r02"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r13204"},
			{Category: EOL, Lexeme: "\n"},
			{Category: OUTPUT, Lexeme: "output"},
			{Category: CONSTANT, Lexeme: "0"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with sub, register, into, and comments",
		input:       "scanner_tests/simple_tests/scanner_test7.txt",
		expectedTokens: []models.Token{
			{Category: COMMENT, Lexeme: "//"},
			{Category: ARITHOP, Lexeme: "sub"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r4"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r5"},
			{Category: COMMENT, Lexeme: "//"},
			{Category: COMMENT, Lexeme: "//"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "scanner input with some invalid tokens",
		input:       "scanner_tests/simple_tests/scanner_test8.txt",
		expectedTokens: []models.Token{
			{Category: INVALID, Lexeme: "storei"},
			{Category: INVALID, Lexeme: "30c"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: true,
	},
	{
		description: "scanner input with some invalid token at the end of the file",
		input:       "scanner_tests/simple_tests/scanner_test9.txt",
		expectedTokens: []models.Token{
			{Category: MEMOP, Lexeme: "load"},
			{Category: INVALID, Lexeme: "a"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: true,
	},
}

var complexTestCases = []TestCase{
	{
		description: "t1.i",
		input:       "scanner_tests/complex_tests/t1.i.txt",
		expectedTokens: []models.Token{
			{Category: COMMENT, Lexeme: "//"},
			{Category: COMMENT, Lexeme: "//"},
			{Category: COMMENT, Lexeme: "//"},
			{Category: COMMENT, Lexeme: "//"},

			{Category: LOADI, Lexeme: "loadI"},
			{Category: INVALID, Lexeme: "10a"},

			{Category: INVALID, Lexeme: "storea"},

			{Category: MEMOP, Lexeme: "load"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r1"},

			{Category: EOL, Lexeme: "\n"},

			{Category: COMMENT, Lexeme: "//"},
			{Category: COMMENT, Lexeme: "//"},

			{Category: INVALID, Lexeme: "addI"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: true,
	},
	{
		description: "t2.i",
		input:       "scanner_tests/complex_tests/t2.i.txt",
		expectedTokens: []models.Token{
			{Category: COMMENT, Lexeme: "//"},
			{Category: COMMENT, Lexeme: "//"},
			{Category: COMMENT, Lexeme: "//"},
			{Category: COMMENT, Lexeme: "//"},
			{Category: COMMENT, Lexeme: "//"},
			// loadI 27  => r1
			{Category: LOADI, Lexeme: "loadI"},
			{Category: CONSTANT, Lexeme: "27"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: EOL, Lexeme: "\n"},
			// loadI 27=>r1
			{Category: LOADI, Lexeme: "loadI"},
			{Category: CONSTANT, Lexeme: "27"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: EOL, Lexeme: "\n"},
			// load  r1 => r2
			{Category: MEMOP, Lexeme: "load"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			// load  r1 => r2
			{Category: MEMOP, Lexeme: "load"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			// load  r1 =>r2
			{Category: MEMOP, Lexeme: "load"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			// store r2 => r4
			{Category: MEMOP, Lexeme: "store"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r4"},
			{Category: EOL, Lexeme: "\n"},
			// add   r1,r2 => r3
			{Category: ARITHOP, Lexeme: "add"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: EOL, Lexeme: "\n"},
			// sub   r3, r4 => r5
			{Category: ARITHOP, Lexeme: "sub"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r4"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r5"},
			{Category: EOL, Lexeme: "\n"},
			// mult  r5, r6 => r10
			{Category: ARITHOP, Lexeme: "mult"},
			{Category: REGISTER, Lexeme: "r5"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r6"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r10"},
			{Category: EOL, Lexeme: "\n"},
			// lshift  r0, r3 => r2
			{Category: ARITHOP, Lexeme: "lshift"},
			{Category: REGISTER, Lexeme: "r0"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			// rshift  r2, r3 => r2
			{Category: ARITHOP, Lexeme: "rshift"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			// output 1024
			{Category: OUTPUT, Lexeme: "output"},
			{Category: CONSTANT, Lexeme: "1024"},
			{Category: EOL, Lexeme: "\n"},
			// nop
			{Category: NOP, Lexeme: "nop"},
			{Category: EOL, Lexeme: "\n"},
			// Final EOF
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "t7.i",
		input:       "scanner_tests/complex_tests/t7.i.txt",
		expectedTokens: []models.Token{
			// loadI 20 => r1
			{Category: LOADI, Lexeme: "loadI"},
			{Category: CONSTANT, Lexeme: "20"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: EOL, Lexeme: "\n"},

			// load  r1 => r2
			{Category: MEMOP, Lexeme: "load"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},

			// loadI r24 => r3
			{Category: LOADI, Lexeme: "loadI"},
			{Category: REGISTER, Lexeme: "r24"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: EOL, Lexeme: "\n"},

			// load  r3 => r4
			{Category: MEMOP, Lexeme: "load"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r4"},
			{Category: EOL, Lexeme: "\n"},

			// add   r2, 3 => r4
			{Category: ARITHOP, Lexeme: "add"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: COMMA, Lexeme: ","},
			{Category: CONSTANT, Lexeme: "3"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r4"},
			{Category: EOL, Lexeme: "\n"},

			// mult  r1, r2 =>5
			{Category: ARITHOP, Lexeme: "mult"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: INTO, Lexeme: "=>"},
			{Category: CONSTANT, Lexeme: "5"},
			{Category: EOL, Lexeme: "\n"},

			// add   r4, => r6
			{Category: ARITHOP, Lexeme: "add"},
			{Category: REGISTER, Lexeme: "r4"},
			{Category: COMMA, Lexeme: ","},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r6"},
			{Category: EOL, Lexeme: "\n"},

			// store r6 =>
			{Category: MEMOP, Lexeme: "store"},
			{Category: REGISTER, Lexeme: "r6"},
			{Category: INTO, Lexeme: "=>"},
			{Category: EOL, Lexeme: "\n"},

			// output 20
			{Category: OUTPUT, Lexeme: "output"},
			{Category: CONSTANT, Lexeme: "20"},

			// EOF
			{Category: EOL, Lexeme: "\n"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "t9.i",
		input:       "scanner_tests/complex_tests/t9.i.txt",
		expectedTokens: []models.Token{
			// loadI 20=>r1
			{Category: LOADI, Lexeme: "loadI"},
			{Category: CONSTANT, Lexeme: "20"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: EOL, Lexeme: "\n"},

			// load r1=>r2
			{Category: MEMOP, Lexeme: "load"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},

			// mult  r1,r2 => r3 r4
			{Category: ARITHOP, Lexeme: "mult"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: REGISTER, Lexeme: "r4"},
			{Category: EOL, Lexeme: "\n"},

			// mult  r1,r2 => r3 a
			{Category: ARITHOP, Lexeme: "mult"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: INVALID, Lexeme: "a "},

			// store r3 => r1
			{Category: MEMOP, Lexeme: "store"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: EOL, Lexeme: "\n"},

			// output 20
			{Category: OUTPUT, Lexeme: "output"},
			{Category: CONSTANT, Lexeme: "20"},
			{Category: EOL, Lexeme: "\n"},

			// EOF
			{Category: EOL, Lexeme: "\n"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: true,
	},
	{
		description: "t12.i",
		input:       "scanner_tests/complex_tests/t12.i.txt",
		expectedTokens: []models.Token{
			//   loadI 27  => r1
			{Category: LOADI, Lexeme: "loadI"},
			{Category: CONSTANT, Lexeme: "27"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: EOL, Lexeme: "\n"},
			//   loadI 27=>r1
			{Category: LOADI, Lexeme: "loadI"},
			{Category: CONSTANT, Lexeme: "27"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: EOL, Lexeme: "\n"},
			//   load  r1 => r2
			{Category: MEMOP, Lexeme: "load"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			//   load  r1 => r2
			{Category: MEMOP, Lexeme: "load"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			//   load  r1 =>r2
			{Category: MEMOP, Lexeme: "load"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			//   store r2 => r4
			{Category: MEMOP, Lexeme: "store"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r4"},
			{Category: EOL, Lexeme: "\n"},
			//   add   r1,r2 => r3
			{Category: ARITHOP, Lexeme: "add"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: EOL, Lexeme: "\n"},
			//   sub   r3, r4 => r5
			{Category: ARITHOP, Lexeme: "sub"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r4"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r5"},
			{Category: EOL, Lexeme: "\n"},
			//   mult  r5, r6 => r10
			{Category: ARITHOP, Lexeme: "mult"},
			{Category: REGISTER, Lexeme: "r5"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r6"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r10"},
			{Category: EOL, Lexeme: "\n"},
			//   lshift  r0, r3 => r2
			{Category: ARITHOP, Lexeme: "lshift"},
			{Category: REGISTER, Lexeme: "r0"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			//   rshift  r2, r3 => r2
			{Category: ARITHOP, Lexeme: "rshift"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			//   output 1024
			{Category: OUTPUT, Lexeme: "output"},
			{Category: CONSTANT, Lexeme: "1024"},
			{Category: EOL, Lexeme: "\n"},
			//   nop
			{Category: NOP, Lexeme: "nop"},
			{Category: EOL, Lexeme: "\n"},
			// Final EOF
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "t17.i",
		input:       "scanner_tests/complex_tests/t17.i.txt",
		expectedTokens: []models.Token{
			// store => r5
			{Category: MEMOP, Lexeme: "store"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r5"},
			{Category: EOL, Lexeme: "\n"},
			// store r1 r5
			{Category: MEMOP, Lexeme: "store"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: REGISTER, Lexeme: "r5"},
			{Category: EOL, Lexeme: "\n"},
			// store r1 =>
			{Category: MEMOP, Lexeme: "store"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: EOL, Lexeme: "\n"},
			// loadI => r1
			{Category: LOADI, Lexeme: "loadI"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: EOL, Lexeme: "\n"},
			// loadI 1 r2
			{Category: LOADI, Lexeme: "loadI"},
			{Category: CONSTANT, Lexeme: "1"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			// loadI 1 =>
			{Category: LOADI, Lexeme: "loadI"},
			{Category: CONSTANT, Lexeme: "1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: EOL, Lexeme: "\n"},
			// add ,r2=>r3
			{Category: ARITHOP, Lexeme: "add"},
			{Category: COMMA, Lexeme: ","},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: EOL, Lexeme: "\n"},
			// add r2=>r3
			{Category: ARITHOP, Lexeme: "add"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: EOL, Lexeme: "\n"},
			// add r1 => r3
			{Category: ARITHOP, Lexeme: "add"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r3"},
			{Category: EOL, Lexeme: "\n"},
			// add r1,=>r2
			{Category: ARITHOP, Lexeme: "add"},
			{Category: REGISTER, Lexeme: "r1"},
			{Category: COMMA, Lexeme: ","},
			{Category: INTO, Lexeme: "=>"},
			{Category: REGISTER, Lexeme: "r2"},
			{Category: EOL, Lexeme: "\n"},
			{Category: EOL, Lexeme: "\n"},
			// output
			{Category: OUTPUT, Lexeme: "output"},
			{Category: EOF, Lexeme: ""},
		},
		expectedError: false,
	},
}

var performanceTestCases = []PerformanceTestCase{
	{
		description:   "Performance Test 1",
		input:         "scanner_tests/complex_tests/t128k.i.txt",
		MaximumTimeMs: 200,
	},
}

func TestSimpleScannerTestCases(t *testing.T) {
	for _, tc := range simpleTestCases {
		t.Run(tc.description, func(t *testing.T) {
			runTest(tc, t)
		})
	}
}

func TestComplexScannerTestCases(t *testing.T) {
	for _, tc := range complexTestCases {
		t.Run(tc.description, func(t *testing.T) {
			runTest(tc, t)
		})
	}
}

func TestScannerPerformance(t *testing.T) {
	for _, tc := range performanceTestCases {
		t.Run(tc.description, func(t *testing.T) {
			start := time.Now()

			file, err := os.Open(tc.input)
			if err != nil {
				t.Fatalf("Failed to open file: %v", err)
			}
			defer file.Close()

			scanner := New(file)

			for {
				tok, err := scanner.NextToken()
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					break
				}
				if tok.Category == EOF {
					break
				}
			}

			duration := time.Since(start)
			if duration.Milliseconds() > tc.MaximumTimeMs {
				t.Errorf("Performance test failed: took %d ms, expected maximum %d ms", duration.Milliseconds(), tc.MaximumTimeMs)
			}
		})
	}

}

func runTest(tc TestCase, t *testing.T) {
	file, err := os.Open(tc.input)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := New(file)
	var tokens []models.Token
	for {
		tok, err := scanner.NextToken()
		if err != nil {
			if !tc.expectedError {
				t.Errorf("Unexpected error: %v", err)
			}
		}
		tokens = append(tokens, tok)
		if tok.Category == EOF {
			break
		}
	}

	if len(tokens) != len(tc.expectedTokens) {
		t.Errorf("Expected %d tokens, got %d", len(tc.expectedTokens), len(tokens))
		// return true
	}

	for i, expected := range tc.expectedTokens {
		curToken := tokens[i]
		if curToken.Category != expected.Category {
			t.Errorf("Token %d - expected category %+v, got %+v", i, expected.Category, curToken.Category)
		}

		if curToken.Lexeme != expected.Lexeme {
			t.Errorf("Token %d - expected lexeme %+v, got %+v", i, expected.Lexeme, curToken.Lexeme)
		}
	}
}
