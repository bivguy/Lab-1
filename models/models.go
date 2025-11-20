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
	PR     int
	NU     float64
	Active bool
}

func (op Operand) String() string {
	if !op.Active {
		return "[ Inactive Operand ]"
	}
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
	switch op.Opcode {
	// ARITH (two uses, one def)
	case "add", "mult", "sub", "lshift", "rshift": // add rA,rB => rC
		return fmt.Sprintf("%s r%d,r%d => r%d",
			op.Opcode, op.OpOne.VR, op.OpTwo.VR, op.OpThree.VR)

	// LOAD variants
	case "load": // load rAddr => rDst
		return fmt.Sprintf("load r%d => r%d",
			op.OpOne.VR, op.OpThree.VR)

	case "loadI":
		return fmt.Sprintf("loadI %d => r%d",
			op.OpOne.SR, op.OpThree.VR)

	// STORE (two uses, no def)
	case "store": // store rVal => rAddr
		return fmt.Sprintf("store r%d => r%d",
			op.OpOne.VR, op.OpThree.VR)

	// OUTPUT
	case "output": // output => rX
		return fmt.Sprintf("output => %d", op.OpThree.SR)

	// NOP
	case "nop":
		return fmt.Sprintf("nop")
	default:
		return fmt.Sprintf("%s ???", op.Opcode)
	}
}

type Status int

type DependenceNode struct {
	Op           *OperationNode
	Edges        map[int]*DependenceEdge
	ReverseEdges map[int]*DependenceEdge
	TotalLatency int
	Latency      int
	Status       Status
}

type EdgeType int

const (
	DATA          EdgeType = iota // 0
	CONFLICT                      // 1
	SERIALIZATION                 // 2
)

func (e EdgeType) String() string {
	switch e {
	case DATA:
		return "DATA"
	case CONFLICT:
		return "CONFLICT"
	case SERIALIZATION:
		return "SERIALIZATION"
	default:
		return fmt.Sprintf("EdgeType(%d)", e)
	}
}

type DependenceEdge struct {
	// From    *DependenceNode
	To      *DependenceNode
	Type    EdgeType
	Latency int
}
