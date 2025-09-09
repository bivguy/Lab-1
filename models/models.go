package models

import (
	"container/list"
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

type IntermediateRepresentation *list.List
