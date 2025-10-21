package allocator

import (
	"container/list"
	"math"

	m "github.com/bivguy/Comp412/models"
)

func (a *allocator) getAPR(VR int, NU float64) int {
	l := len(a.freePRStack)
	var x int
	if l > 0 {
		// fmt.Print("about to pop from stack: ", a.freePRStack, ".\n")
		x = a.popStack()
	} else {
		furthestNextUse := -1

		for i := 0; i < a.maxPR; i++ {
			if !a.marks[i] && furthestNextUse < int(a.PRNU[i]) {
				x = i
				furthestNextUse = int(a.PRNU[i])
			}
		}

		// pick an unmarked x to spill
		// fmt.Print("about to spill ", a.curOperationNode.Value, ".\n")
		a.spill(x)
		a.VRToPR[a.PRToVR[x]] = -1
	}

	a.VRToPR[VR] = x
	a.PRToVR[x] = VR
	a.PRNU[x] = NU

	return x
}

func (a *allocator) freeAPR(pr int) {
	vr := a.PRToVR[pr]

	a.VRToPR[vr] = -1
	a.PRToVR[pr] = -1
	a.PRNU[pr] = math.Inf(1)

	a.freePRStack = append(a.freePRStack, pr)
}

// evicts a value from a register and stores its value
func (a *allocator) spill(pr int) {
	loadIInstruction := &m.OperationNode{
		Opcode:  "loadI",
		OpOne:   m.Operand{Active: true, SR: a.memAddress},
		OpThree: m.Operand{Active: true, PR: a.maxPR},
	}

	storeInstruction := &m.OperationNode{
		Opcode:  "store",
		OpOne:   m.Operand{Active: true, VR: a.PRToVR[pr], PR: pr, NU: a.PRNU[pr]},
		OpThree: m.Operand{Active: true, PR: a.maxPR},
	}

	// connect the linked list
	a.IR.InsertBefore(loadIInstruction, a.curOperationNode)
	a.IR.InsertBefore(storeInstruction, a.curOperationNode)

	// update necessary fields
	a.VRToSpillLoc[a.PRToVR[pr]] = a.memAddress
	a.memAddress += 4
}

func (a *allocator) restore(vr int, pr int) {
	// fmt.Print("about to restore ", a.curOperationNode.Value, ".\n")
	loadIInstruction := &m.OperationNode{
		Opcode:  "loadI",
		OpOne:   m.Operand{Active: true, SR: a.VRToSpillLoc[vr]},
		OpThree: m.Operand{Active: true, PR: a.maxPR},
	}

	loadInstruction := &m.OperationNode{
		Opcode:  "load",
		OpOne:   m.Operand{Active: true, PR: a.maxPR},
		OpThree: m.Operand{Active: true, VR: vr, PR: pr, NU: a.PRNU[pr]},
	}

	// connect the linked list
	a.IR.InsertBefore(loadIInstruction, a.curOperationNode)
	a.IR.InsertBefore(loadInstruction, a.curOperationNode)
}

func getMaxLive(IR *list.List, maxVR int) int {
	var maxLive int

	live := make([]bool, maxVR)
	curLive, maxLive := 0, 0

	for node := IR.Back(); node != nil; node = node.Prev() {
		op := node.Value.(*m.OperationNode)
		if op.Opcode == "nop" || op.Opcode == "output" {
			continue
		}
		ops := []*m.Operand{&op.OpOne, &op.OpTwo, &op.OpThree}

		// go through uses
		for i, u := range ops {
			if !u.Active || !isRegister(op.Opcode, i) || isDefinition(op.Opcode, i) {
				continue
			}
			vr := u.VR
			if !live[vr] {
				// if vr >= 0 && vr < len(live) && !live[vr] {
				live[vr] = true
				curLive++
			}
		}

		// go through definitions
		for i, d := range ops {
			if !d.Active || !isRegister(op.Opcode, i) || !isDefinition(op.Opcode, i) {
				continue
			}
			vr := d.VR
			// if vr >= 0 && vr < len(live) && live[vr] {
			live[vr] = false
			curLive--
			// }
		}

		maxLive = max(curLive, maxLive)

	}
	return maxLive
}

func isRegister(opcode string, i int) bool {
	if (opcode == "loadI" && i == 0) || (opcode == "output" && i == 2) {
		return false
	}

	return true
}

func (a *allocator) popStack() int {
	l := len(a.freePRStack)
	x := a.freePRStack[l-1]
	a.freePRStack = a.freePRStack[:l-1]
	// slices.Delete(a.freePRStack, l-1, l)

	return x
}
