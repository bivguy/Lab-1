package renamer

import (
	"container/list"
	"math"

	m "github.com/bivguy/Comp412/models"
)

type renamer struct {
	SRToVR []int
	LU     []float64

	MaxSR int
	MaxVR int
	index int

	IR *list.List
}

func New(maxSR int, IR *list.List) *renamer {
	SRToVR := make([]int, maxSR+1)
	LU := make([]float64, maxSR+1)

	for i := 0; i <= maxSR; i++ {
		SRToVR[i] = -1
		LU[i] = math.Inf(1)
	}

	return &renamer{
		SRToVR: SRToVR,
		LU:     LU,
		MaxSR:  maxSR,
		MaxVR:  0,
		index:  IR.Len(),
		IR:     IR,
	}
}

func (r *renamer) Rename() *list.List {
	vrName := 0
	// var liveMap map[]
	// curLive, maxLive := 0, 0
	// go through the IR in reverse order
	for node := r.IR.Back(); node != nil; node = node.Prev() {
		op := node.Value.(*m.OperationNode)

		if op.Opcode == "nop" || op.Opcode == "output" {
			r.index--
			continue
		}

		operandList := []*m.Operand{&op.OpOne, &op.OpTwo, &op.OpThree}

		// go through each operand that is defined
		for i, o := range operandList {
			// skip if its not active or if its not definiition
			if !o.Active || !isRegister(op.Opcode, i) || !isDefinition(op.Opcode, i) {
				continue
			}

			if r.SRToVR[o.SR] == -1 {
				r.SRToVR[o.SR] = vrName
				vrName++
				// curLive += 1
			}

			o.VR = r.SRToVR[o.SR]
			o.NU = r.LU[o.SR]

			r.SRToVR[o.SR] = -1
			r.LU[o.SR] = math.Inf(1)
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

			r.LU[o.SR] = float64(r.index)
		}

		if r.MaxVR < vrName {
			r.MaxVR = vrName
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
	if (opcode == "loadI" && i == 0) || (opcode == "output" && i == 2) {
		return false
	}

	return true
}
