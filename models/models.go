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
	SR     int
	VR     int
	NU     int
	Active bool
}

func (op Operand) String() string {
	return fmt.Sprintf("[SR=%d, VR=%d, NU=%d, Active=%t]", op.SR, op.VR, op.NU, op.Active)
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
		"Line %d: OpCode: %s; Operand One: %v; Operand Two: %v; Operand Three: %v",
		op.Line,
		op.Opcode,
		op.OpOne,
		op.OpTwo,
		op.OpThree,
	)
}
