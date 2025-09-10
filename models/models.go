package models

import (
	"fmt"
)

type SyntacticCategory int

type Token struct {
	Category   SyntacticCategory
	Lexeme     string
	LineNumber int
}

func (t Token) String() string {
	return fmt.Sprintf("<%v, %q> at line %d", t.Category, t.Lexeme, t.LineNumber)
}

type Operand struct {
	SR int
}

type OperationNode struct {
	Line    int
	Opcode  string
	OpOne   Operand
	OpTwo   Operand
	OpThree Operand
}

func (op OperationNode) String() string {
	return fmt.Sprintf(
		"Line %d: %s %v %v %v",
		op.Line,
		op.Opcode,
		op.OpOne,
		op.OpTwo,
		op.OpThree,
	)
}
