package parser

import (
	"container/list"
	"os"
	"testing"

	s "github.com/bivguy/Comp412/scanner"
)

func TestSimpleParserTestCases(t *testing.T) {
	for _, tc := range simpleTestCases {
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
	scanner := s.New(file)
	parser := New(scanner)
	IR, err := parser.Parse()
	if err != nil {
		if !tc.expectedError {
			t.Errorf("Unexpected error: %v", err)
		}
	}

	if !compareIR(tc.expectedIR, IR) {
		t.Errorf("IR mismatch in %s", tc.description)
	}

}

func compareIR(expected, actual *list.List) bool {
	if expected.Len() != actual.Len() {
		return false
	}

	e1 := expected.Front()
	e2 := actual.Front()

	for e1 != nil && e2 != nil {
		if e1.Value != e2.Value {
			return false
		}
		e1 = e1.Next()
		e2 = e2.Next()
	}

	return true
}
