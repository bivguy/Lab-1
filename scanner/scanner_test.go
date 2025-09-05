package scanner

import (
	"os"
	"testing"
)

type TestCase struct {
	description    string
	input          string
	expectedTokens []token
	expectedError  bool
}

var simpleTestCases = []TestCase{
	{
		description: "Valid scanner input with load",
		input:       "scanner_tests/simple_tests/scanner_test1.txt",
		expectedTokens: []token{
			{category: MEMOP, lexeme: "load"},
			{category: EOL, lexeme: "\n"},
			{category: EOF, lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with lshift and load",
		input:       "scanner_tests/simple_tests/scanner_test2.txt",
		expectedTokens: []token{
			{category: ARITHOP, lexeme: "lshift"},
			{category: EOL, lexeme: "\n"},
			{category: ARITHOP, lexeme: "lshift"},
			{category: MEMOP, lexeme: "load"},
			{category: EOF, lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with lshift, load, and loadI",
		input:       "scanner_tests/simple_tests/scanner_test3.txt",
		expectedTokens: []token{
			{category: ARITHOP, lexeme: "lshift"},
			{category: MEMOP, lexeme: "load"},
			{category: EOL, lexeme: "\n"},
			{category: LOADI, lexeme: "loadI"},
			{category: MEMOP, lexeme: "load"},
			{category: LOADI, lexeme: "loadI"},
			{category: EOL, lexeme: "\n"},
			{category: EOL, lexeme: "\n"},
			{category: ARITHOP, lexeme: "lshift"},
			{category: EOF, lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with store, constant, and sub",
		input:       "scanner_tests/simple_tests/scanner_test4.txt",
		expectedTokens: []token{
			{category: MEMOP, lexeme: "store"},
			{category: CONSTANT, lexeme: "132"},
			{category: EOL, lexeme: "\n"},
			{category: ARITHOP, lexeme: "sub"},
			{category: CONSTANT, lexeme: "142"},
			{category: CONSTANT, lexeme: "2"},
			{category: EOL, lexeme: "\n"},
			{category: EOF, lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with add, comma, register, and into",
		input:       "scanner_tests/simple_tests/scanner_test5.txt",
		expectedTokens: []token{
			{category: ARITHOP, lexeme: "add"},
			{category: REGISTER, lexeme: "r0"},
			{category: COMMA, lexeme: ","},
			{category: REGISTER, lexeme: "r1"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r2"},
			{category: EOF, lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with store register, into, output, and constant",
		input:       "scanner_tests/simple_tests/scanner_test6.txt",
		expectedTokens: []token{
			{category: MEMOP, lexeme: "store"},
			{category: REGISTER, lexeme: "r02"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r13204"},
			{category: EOL, lexeme: "\n"},
			{category: OUTPUT, lexeme: "output"},
			{category: CONSTANT, lexeme: "0"},
			{category: EOF, lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "Valid scanner input with sub, register, into, and comments",
		input:       "scanner_tests/simple_tests/scanner_test7.txt",
		expectedTokens: []token{
			{category: COMMENT, lexeme: "//"},
			{category: ARITHOP, lexeme: "sub"},
			{category: REGISTER, lexeme: "r3"},
			{category: COMMA, lexeme: ","},
			{category: REGISTER, lexeme: "r4"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r5"},
			{category: COMMENT, lexeme: "//"},
			{category: COMMENT, lexeme: "//"},
			{category: EOF, lexeme: ""},
		},
		expectedError: false,
	},
}

var complexTestCases = []TestCase{
	{
		description: "t2.i",
		input:       "scanner_tests/complex_tests/t2.i.txt",
		expectedTokens: []token{
			{category: COMMENT, lexeme: "//"},
			{category: COMMENT, lexeme: "//"},
			{category: COMMENT, lexeme: "//"},
			{category: COMMENT, lexeme: "//"},
			{category: COMMENT, lexeme: "//"},
			// loadI 27  => r1
			{category: LOADI, lexeme: "loadI"},
			{category: CONSTANT, lexeme: "27"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r1"},
			{category: EOL, lexeme: "\n"},
			// loadI 27=>r1
			{category: LOADI, lexeme: "loadI"},
			{category: CONSTANT, lexeme: "27"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r1"},
			{category: EOL, lexeme: "\n"},
			// load  r1 => r2
			{category: MEMOP, lexeme: "load"},
			{category: REGISTER, lexeme: "r1"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r2"},
			{category: EOL, lexeme: "\n"},
			// load  r1 => r2
			{category: MEMOP, lexeme: "load"},
			{category: REGISTER, lexeme: "r1"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r2"},
			{category: EOL, lexeme: "\n"},
			// load  r1 =>r2
			{category: MEMOP, lexeme: "load"},
			{category: REGISTER, lexeme: "r1"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r2"},
			{category: EOL, lexeme: "\n"},
			// store r2 => r4
			{category: MEMOP, lexeme: "store"},
			{category: REGISTER, lexeme: "r2"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r4"},
			{category: EOL, lexeme: "\n"},
			// add   r1,r2 => r3
			{category: ARITHOP, lexeme: "add"},
			{category: REGISTER, lexeme: "r1"},
			{category: COMMA, lexeme: ","},
			{category: REGISTER, lexeme: "r2"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r3"},
			{category: EOL, lexeme: "\n"},
			// sub   r3, r4 => r5
			{category: ARITHOP, lexeme: "sub"},
			{category: REGISTER, lexeme: "r3"},
			{category: COMMA, lexeme: ","},
			{category: REGISTER, lexeme: "r4"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r5"},
			{category: EOL, lexeme: "\n"},
			// mult  r5, r6 => r10
			{category: ARITHOP, lexeme: "mult"},
			{category: REGISTER, lexeme: "r5"},
			{category: COMMA, lexeme: ","},
			{category: REGISTER, lexeme: "r6"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r10"},
			{category: EOL, lexeme: "\n"},
			// lshift  r0, r3 => r2
			{category: ARITHOP, lexeme: "lshift"},
			{category: REGISTER, lexeme: "r0"},
			{category: COMMA, lexeme: ","},
			{category: REGISTER, lexeme: "r3"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r2"},
			{category: EOL, lexeme: "\n"},
			// rshift  r2, r3 => r2
			{category: ARITHOP, lexeme: "rshift"},
			{category: REGISTER, lexeme: "r2"},
			{category: COMMA, lexeme: ","},
			{category: REGISTER, lexeme: "r3"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r2"},
			{category: EOL, lexeme: "\n"},
			// output 1024
			{category: OUTPUT, lexeme: "output"},
			{category: CONSTANT, lexeme: "1024"},
			{category: EOL, lexeme: "\n"},
			// nop
			{category: NOP, lexeme: "nop"},
			{category: EOL, lexeme: "\n"},
			// Final EOF
			{category: EOF, lexeme: ""},
		},
		expectedError: false,
	},
	{
		description: "t17.i",
		input:       "scanner_tests/complex_tests/t17.i.txt",
		expectedTokens: []token{
			// store => r5
			{category: MEMOP, lexeme: "store"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r5"},
			{category: EOL, lexeme: "\n"},
			// store r1 r5
			{category: MEMOP, lexeme: "store"},
			{category: REGISTER, lexeme: "r1"},
			{category: REGISTER, lexeme: "r5"},
			{category: EOL, lexeme: "\n"},
			// store r1 =>
			{category: MEMOP, lexeme: "store"},
			{category: REGISTER, lexeme: "r1"},
			{category: INTO, lexeme: "=>"},
			{category: EOL, lexeme: "\n"},
			// loadI => r1
			{category: LOADI, lexeme: "loadI"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r1"},
			{category: EOL, lexeme: "\n"},
			// loadI 1 r2
			{category: LOADI, lexeme: "loadI"},
			{category: CONSTANT, lexeme: "1"},
			{category: REGISTER, lexeme: "r2"},
			{category: EOL, lexeme: "\n"},
			// loadI 1 =>
			{category: LOADI, lexeme: "loadI"},
			{category: CONSTANT, lexeme: "1"},
			{category: INTO, lexeme: "=>"},
			{category: EOL, lexeme: "\n"},
			// add ,r2=>r3
			{category: ARITHOP, lexeme: "add"},
			{category: COMMA, lexeme: ","},
			{category: REGISTER, lexeme: "r2"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r3"},
			{category: EOL, lexeme: "\n"},
			// add r2=>r3
			{category: ARITHOP, lexeme: "add"},
			{category: REGISTER, lexeme: "r2"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r3"},
			{category: EOL, lexeme: "\n"},
			// add r1 => r3
			{category: ARITHOP, lexeme: "add"},
			{category: REGISTER, lexeme: "r1"},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r3"},
			{category: EOL, lexeme: "\n"},
			// add r1,=>r2
			{category: ARITHOP, lexeme: "add"},
			{category: REGISTER, lexeme: "r1"},
			{category: COMMA, lexeme: ","},
			{category: INTO, lexeme: "=>"},
			{category: REGISTER, lexeme: "r2"},
			{category: EOL, lexeme: "\n"},
			{category: EOL, lexeme: "\n"},
			// output
			{category: OUTPUT, lexeme: "output"},
			{category: EOF, lexeme: ""},
		},
		expectedError: false,
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

func runTest(tc TestCase, t *testing.T) {
	file, err := os.Open(tc.input)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := New(file)
	var tokens []token
	for {
		tok, err := scanner.NextToken()
		if err != nil {
			if !tc.expectedError {
				t.Errorf("Unexpected error: %v", err)
			}
			break
		}
		tokens = append(tokens, tok)
		if tok.category == EOF {
			break
		}
	}

	if len(tokens) != len(tc.expectedTokens) {
		t.Errorf("Expected %d tokens, got %d", len(tc.expectedTokens), len(tokens))
		// return true
	}

	for i, expected := range tc.expectedTokens {
		if tokens[i] != expected {
			t.Errorf("Token %d - expected %+v, got %+v", i, expected, tokens[i])
		}
	}
	// return false
}
