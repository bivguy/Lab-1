package main

import (
	"container/list"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bivguy/Comp412/allocator"
	c "github.com/bivguy/Comp412/constants"
	m "github.com/bivguy/Comp412/models"
	"github.com/bivguy/Comp412/parser"
	"github.com/bivguy/Comp412/renamer"

	"github.com/bivguy/Comp412/scanner"
)

func main() {
	hFlag := flag.Bool("h", false, "Display help")
	sFlag := flag.Bool("s", false, "Display the Scanner Output")
	pFlag := flag.Bool("p", false, "Display the Parser Output")
	rFlag := flag.Bool("r", false, "Display the Intermediate Representation Output")

	xFlag := flag.Bool("x", false, "Displays the Renamed Intermediate Representation Output")

	flag.Parse()

	// politely report that only a single flag should be passed in
	if flag.NFlag() > 1 {
		fmt.Fprintln(os.Stderr, "Only one flag should be passed at a time; using highest priority (-h, -r, -p, -s).\n")
	}

	if *hFlag {
		helpMessage()
		return
	}

	// filename
	args := flag.Args()

	if len(args) >= 2 {
		if k, err := strconv.Atoi(args[0]); err == nil {
			filename := args[1]
			runAllocatorMode(k, filename)
			return
		}
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: missing <filename>")
		helpMessage()
		return
	}

	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "ERROR: Attempt to open more than one input file.")
		helpMessage()
		return
	}

	path := args[0]
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to open file: %v\n", err)
		helpMessage()
		return
	}
	defer file.Close()
	scanner := scanner.New(file)
	parser := parser.New(scanner)

	if *sFlag || *pFlag || *rFlag || *xFlag {
		if *rFlag {
			IR, err := parser.Parse()
			// we only print the IR there is no error found
			if !parser.ErrorFound && err == nil {
				fmt.Println(PrettyPrintIR(IR))
			} else {
				fmt.Println("\nDue to the syntax error, run terminates.")
			}
		} else if *pFlag {
			IR, err := parser.Parse()
			if parser.ErrorFound || err != nil {
				fmt.Println("Parse found errors")
			} else {
				fmt.Printf("Parse succeeded. Processed %d operations.\n", IR.Len())
			}
		} else if *sFlag {
			for {
				tok, err := scanner.NextToken()
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR %d: Unexpected error: %v\n", tok.LineNumber, err)
				}
				scanner.PrintToken(tok)
				if tok.Category == c.EOF {
					break
				}
			}
		} else if *xFlag {
			IR, err := parser.Parse()
			if parser.ErrorFound || err != nil {
				fmt.Println("Parse found errors")
				return
			}

			largestRegister := parser.GetLargestRegister()

			renamer := renamer.New(largestRegister, IR)

			renamedIR := renamer.Rename()
			fmt.Println(renameIR(renamedIR))
		}
	} else {
		// default behavior is of pflag
		IR, err := parser.Parse()
		if parser.ErrorFound || err != nil {
			fmt.Println("Parse found errors")
		} else {
			fmt.Printf("Parse succeeded. Processed %d operations.\n", IR.Len())
		}
	}
}

// helpMessage prints to the command line all the possible commands for the 412fe applications
func helpMessage() {
	fmt.Println("Usage: 412fe [flags]")
	fmt.Println()
	fmt.Println("412fe is the front end for the compiler project.")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -h \t\t Display this help message.")
	fmt.Println("  -s <filename>\t Read the specified file, scan it, and print a list of tokens.")
	fmt.Println("  -p <filename>\t Read the specified file, scan it, parse it, and report success or failure.")
	fmt.Println("  -r <filename>\t Read the specified file, scan it, parse it, and print the intermediate representation.")
	fmt.Println("No flag provided: behaves as if -p <filename> was specified.")
}

func PrettyPrintIR(ir *list.List) string {
	if ir == nil || ir.Len() == 0 {
		return "[empty IR]"
	}

	result := ""
	for e := ir.Front(); e != nil; e = e.Next() {
		op := e.Value.(*m.OperationNode)
		result += fmt.Sprintf("%s\n\n", op)
	}
	return result
}

func renameIR(ir *list.List) string {
	var b strings.Builder

	for e := ir.Front(); e != nil; e = e.Next() {
		var op *m.OperationNode
		switch v := e.Value.(type) {
		case *m.OperationNode:
			op = v
		case m.OperationNode:
			tmp := v
			op = &tmp
		default:
			continue
		}

		switch op.Opcode {
		// ARITH (two uses, one def)
		case "add", "mult", "sub", "lshift", "rshift": // add rA,rB => rC
			fmt.Fprintf(&b, "%s r%d,r%d => r%d\n",
				op.Opcode, op.OpOne.VR, op.OpTwo.VR, op.OpThree.VR)

		// LOAD variants
		case "load": // load rAddr => rDst
			fmt.Fprintf(&b, "load r%d => r%d\n",
				op.OpOne.VR, op.OpThree.VR)

		case "loadI":
			fmt.Fprintf(&b, "loadI %d => r%d\n",
				op.OpOne.SR, op.OpThree.VR)

		// STORE (two uses, no def)
		case "store": // store rVal => rAddr
			fmt.Fprintf(&b, "store r%d => r%d\n",
				op.OpOne.VR, op.OpThree.VR)

		// OUTPUT
		case "output": // output => rX
			fmt.Fprintf(&b, "output => %d\n", op.OpThree.SR)

		// NOP
		case "nop":
			fmt.Fprintf(&b, "nop\n")
		// case "lshift":

		// case "rshift":

		default:
			fmt.Fprintf(&b, "%s ???\n", op.Opcode)
		}
	}

	return b.String()
}

func allocIR(ir *list.List) string {
	var b strings.Builder

	for e := ir.Front(); e != nil; e = e.Next() {
		var op *m.OperationNode
		switch v := e.Value.(type) {
		case *m.OperationNode:
			op = v
		case m.OperationNode:
			tmp := v
			op = &tmp
		default:
			continue
		}

		switch op.Opcode {
		// ARITH (two uses, one def)
		case "add", "mult", "sub", "lshift", "rshift": // add rA,rB => rC
			fmt.Fprintf(&b, "%s r%d,r%d => r%d\n",
				op.Opcode, op.OpOne.PR, op.OpTwo.PR, op.OpThree.PR)

		// LOAD variants
		case "load": // load rAddr => rDst
			fmt.Fprintf(&b, "load r%d => r%d\n",
				op.OpOne.PR, op.OpThree.PR)

		case "loadI":
			fmt.Fprintf(&b, "loadI %d => r%d\n",
				op.OpOne.SR, op.OpThree.PR)

		// STORE (two uses, no def)
		case "store": // store rVal => rAddr
			fmt.Fprintf(&b, "store r%d => r%d\n",
				op.OpOne.PR, op.OpThree.PR)

		// OUTPUT
		case "output": // output => rX
			fmt.Fprintf(&b, "output => %d\n", op.OpThree.SR)

		// NOP
		case "nop":
			fmt.Fprintf(&b, "nop\n")

		default:
			fmt.Fprintf(&b, "%s ???\n", op.Opcode)
		}
	}

	return b.String()
}

func runAllocatorMode(k int, filename string) {
	// Validate k and file
	if k < 3 || k > 64 {
		fmt.Fprintf(os.Stderr, "ERROR: k must be in [3, 64], got %d\n", k)
		os.Exit(1)
	}
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cannot open file '%s': %v\n", filename, err)
		os.Exit(1)
	}
	defer f.Close()

	// scan and parse
	sc := scanner.New(f)
	ps := parser.New(sc)
	IR, err := ps.Parse()
	if ps.ErrorFound || err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: parse failed for '%s'\n", filename)
		os.Exit(1)
	}

	// Rename
	largestRegister := ps.GetLargestRegister()
	renamer := renamer.New(largestRegister, IR)
	IR = renamer.Rename()

	alloc := allocator.New(renamer.SRToVR, renamer.LU, IR, renamer.MaxVR, k)
	IR = alloc.Allocate()

	PRToVR := alloc.PRToVR
	VRToPR := alloc.VRToPR

	for i := 0; i < len(PRToVR); i++ {
		if PRToVR[i] != -1 && i != VRToPR[PRToVR[i]] {
			fmt.Print("INVALID MAPPING\n")
		}
	}

	fmt.Println(allocIR(IR))
}
