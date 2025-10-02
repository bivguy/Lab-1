package renamer

import (
	"container/list"
	"math"

	m "github.com/bivguy/Comp412/models"
)

type renamer struct {
	SRToVR []int
	LU     []int // SWITCH TO SLICES AT SOME POINT
	maxSR  int
	index  int
	IR     *list.List
}

func New(maxSR int, IR *list.List) *renamer {
	SRToVR := make([]int, maxSR+1)
	LU := make([]int, maxSR+1)

	for i := 0; i <= maxSR; i++ {
		SRToVR[i] = -1
		LU[i] = -1
	}

	return &renamer{
		SRToVR: SRToVR,
		LU:     LU,
		maxSR:  maxSR,
		index:  IR.Len(),
		IR:     IR,
	}
}

func (r *renamer) Rename() *list.List {
	vrName := 0
	// var liveMap map[]
	// go through the IR in reverse order
	for node := r.IR.Back(); node != nil; node = node.Prev() {
		op := node.Value.(m.OperationNode)

		if op.Opcode == "nop" {
			continue
		}

		operandList := []*m.Operand{&op.OpOne, &op.OpTwo, &op.OpThree}

		// go through each operand that is defined
		for i, o := range operandList {
			// skip if its not active or if its not definiition
			if !o.Active || !isRegister(op.Opcode, i) || !isDefinition(op.Opcode, i) {
				r.index--
				continue
			}

			if r.SRToVR[o.SR] == -1 {
				r.SRToVR[o.SR] = vrName
				vrName++
			}

			o.VR = r.SRToVR[o.SR]
			o.NU = r.LU[o.SR]

			r.SRToVR[o.SR] = -1
			r.LU[o.SR] = int(math.Inf(1))
		}

		// go through each operand that is used
		for i, o := range operandList {
			// skip if its not active, valid, or if its a definiition
			if !o.Active || !isRegister(op.Opcode, i) || isDefinition(op.Opcode, i) {
				continue
			}

			if r.SRToVR[o.SR] == -1 {
				r.SRToVR[o.SR] = vrName
				vrName++
			}

			o.VR = r.SRToVR[o.SR]
			o.NU = r.LU[o.SR]
		}

		// go through each operand that is used
		for i, o := range operandList {
			// skip if its not active or if its a definiition
			if !o.Active || !isRegister(op.Opcode, i) || isDefinition(op.Opcode, i) {
				continue
			}

			r.LU[o.SR] = r.index
		}

		r.index--
	}

	return r.IR
}

func isDefinition(opcode string, i int) bool {
	if opcode == "store" {
		return false
	} else if i == 2 {
		return true
	}
	return false
}

func isRegister(opcode string, i int) bool {
	if opcode == "loadI" && i == 0 {
		return false
	}

	return true
}
