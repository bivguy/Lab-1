package allocator

import (
	m "github.com/bivguy/Comp412/models"
)

// evicts a value from a register and stores its value
func (a *allocator) spill(pr int) {
	loadIInstruction := &m.OperationNode{
		Opcode: "loadI",
		OpOne:  m.Operand{Active: true, SR: a.memAddress},
		OpThree: m.Operand{Active: true,
			PR: a.maxPR,
		},
	}

	storeInstruction := &m.OperationNode{
		Opcode: "store",
		OpOne:  m.Operand{Active: true, VR: a.PRToVR[pr], PR: pr, NU: a.PRNU[pr]},
		OpThree: m.Operand{Active: true,
			PR: a.maxPR,
		},
	}

	// connect the linked list
	a.IR.InsertBefore(loadIInstruction, a.curOperationNode)
	a.IR.InsertBefore(storeInstruction, a.curOperationNode)

	// update necessary fields
	a.VRToSpillLoc[a.PRToVR[pr]] = a.memAddress
	a.memAddress += 4
}

func (a *allocator) restore(vr int, pr int) {
	constant, ok := a.VRToConstant[vr]
	if ok {
		loadIInstruction := &m.OperationNode{
			Opcode:  "loadI",
			OpOne:   m.Operand{Active: true, SR: constant},
			OpThree: m.Operand{Active: true, PR: pr},
		}

		a.IR.InsertBefore(loadIInstruction, a.curOperationNode)
		return
	}

	// fmt.Print("about to restore ", a.curOperationNode.Value, ".\n")
	loadIInstruction := &m.OperationNode{
		Opcode:  "loadI",
		OpOne:   m.Operand{Active: true, SR: a.VRToSpillLoc[vr]},
		OpThree: m.Operand{Active: true, PR: a.maxPR},
	}

	loadInstruction := &m.OperationNode{
		Opcode: "load",
		OpOne:  m.Operand{Active: true, PR: a.maxPR},

		OpThree: m.Operand{Active: true, VR: vr, PR: pr, NU: a.PRNU[pr]},
	}

	// connect the linked list
	a.IR.InsertBefore(loadIInstruction, a.curOperationNode)
	a.IR.InsertBefore(loadInstruction, a.curOperationNode)
}

func (a *allocator) deletePrevNode() {
	if !a.deletePreviousNode {
		return
	}

	a.IR.Remove(a.curOperationNode.Prev())
	a.deletePreviousNode = false
}
