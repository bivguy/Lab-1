package allocator

import (
	"container/list"
	"math"

	m "github.com/bivguy/Comp412/models"
)

const RESERVEDREGISTER = 32768
const INVALIDREGISTER = 32767

type allocator struct {
	SRToVR []int
	LU     []float64

	VRToSpillLoc []int       // a map from a virtual register to its spill location in memory
	VRToConstant map[int]int // a map from a virtual register to its rematerialization location

	VRToPR      []int
	PRToVR      []int
	PRNU        []float64 // a map from a physical register to the operation that next uses the value that is currently in that physical register
	marks       []bool
	freePRStack []int

	deletePreviousNode bool

	curOperationNode *list.Element

	memAddress int
	maxPR      int
	maxVR      int

	IR *list.List
}

func New(SRToVR []int, LU []float64, IR *list.List, maxVR int, maxPR int, VRToConstant map[int]int) *allocator {
	VRToPR := make([]int, maxVR)
	VRToSpillLoc := make([]int, maxVR)

	//  vr := 0; vr <= maxVR; vr++ {
	for vr := range maxVR {
		VRToPR[vr] = -1
		VRToSpillLoc[vr] = INVALIDREGISTER
	}

	if maxPR < getMaxLive(IR, maxVR) {
		maxPR -= 1
	}

	PRToVR := make([]int, maxPR)
	PRNU := make([]float64, maxPR)
	freePRStack := []int{}
	marks := make([]bool, maxPR)

	for pr := range maxPR {
		PRToVR[pr] = -1
		PRNU[pr] = math.Inf(1)
		marks[pr] = false
	}

	for pr := maxPR - 1; pr >= 0; pr-- {
		freePRStack = append(freePRStack, pr)
	}

	a := &allocator{
		SRToVR: SRToVR,
		LU:     LU,

		VRToSpillLoc: VRToSpillLoc,
		VRToConstant: VRToConstant,
		VRToPR:       VRToPR,
		PRToVR:       PRToVR,
		PRNU:         PRNU,
		marks:        marks,
		freePRStack:  freePRStack,

		deletePreviousNode: false,

		maxPR:      maxPR,
		maxVR:      maxVR,
		memAddress: RESERVEDREGISTER,

		IR: IR,
	}

	return a
}

func (a *allocator) Allocate() *list.List {
	// iterate over the block
	for node := a.IR.Front(); node != nil; node = node.Next() {
		a.deletePrevNode()

		op := node.Value.(*m.OperationNode)
		a.curOperationNode = node
		// TODO: check if this is necessary or not
		if op.Opcode == "nop" || op.Opcode == "output" {
			continue
		}

		// clear the mark in each PR
		for i := range a.marks {
			a.marks[i] = false
		}

		operandList := []*m.Operand{&op.OpOne, &op.OpTwo, &op.OpThree}

		// go through each use, allocating uses
		for i, u := range operandList {
			// skip if it's a definition since its not a use
			if isDefinition(op.Opcode, i) || !u.Active || !isRegister(op.Opcode, i) {
				continue
			}

			pr := a.VRToPR[u.VR]
			// fmt.Println("at op ", op)

			if pr == -1 {
				u.PR = a.getAPR(u.VR, u.NU)
				// fmt.Println("Does not exist for operand ", u, ". Now, it is ", u.PR)
				_, ok := a.VRToConstant[u.VR]
				// restore: only restore if its in the spill address or in the constants
				if a.VRToSpillLoc[u.VR] != INVALIDREGISTER || ok {
					a.restore(u.VR, u.PR)
				}

			} else {
				// fmt.Println(" exist for operand ", u, ", as ", u.PR)
				u.PR = pr
			}
			// set the mark in u.PR
			a.marks[u.PR] = true

			// *** early-free LAST USE right away ***
			// if u.NU == math.Inf(1) && a.PRToVR[u.PR] != -1 {
			// 	a.freeAPR(u.PR)
			// 	a.marks[u.PR] = false
			// }
		}

		// go through each use, checking for last use
		for i, u := range operandList {
			// skip if it's a definition since its not a use
			// TODO: may have to add in more checks (only ones with valid registers)
			if isDefinition(op.Opcode, i) || !u.Active || !isRegister(op.Opcode, i) {
				continue
			}

			if u.NU == math.Inf(1) && a.PRToVR[u.PR] != -1 {
				// fmt.Println("freeing the PR: ", u.PR)
				a.freeAPR(u.PR)
			}
		}

		// reset marks: clear the mark in each PR
		for i := range a.marks {
			a.marks[i] = false
		}

		// allocate defs
		for i, d := range operandList {
			// skip if it's not a definition
			if !isDefinition(op.Opcode, i) || !d.Active || !isRegister(op.Opcode, i) {
				continue
			}

			d.PR = a.getAPR(d.VR, d.NU)
			// fmt.Println("Got the definition pr: ", d.PR, " at def ", d)
			// set the mark in D.PR
			a.marks[d.PR] = true
			// a.PRNU[d.PR] = d.NU??
		}
		// fmt.Println()
		// fmt.Println()
	}

	return a.IR
}

func isDefinition(opcode string, i int) bool {
	if opcode == "store" {
		return false
	} else if i == 2 {
		return true
	}
	return false
}
