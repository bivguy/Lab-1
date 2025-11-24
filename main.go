package main

import (
	"container/list"
	"flag"
	"fmt"
	"os"
	"strings"

	m "github.com/bivguy/Comp412/models"
	"github.com/bivguy/Comp412/parser"
	"github.com/bivguy/Comp412/renamer"
	"github.com/bivguy/Comp412/scheduler"

	"github.com/bivguy/Comp412/scanner"
)

func main() {
	hFlag := flag.Bool("h", false, "Display help")

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
	// open the file
	path := args[0]
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to open file: %v\n", err)
		helpMessage()
		return
	}
	defer file.Close()

	// scan and parse
	scanner := scanner.New(file)
	parser := parser.New(scanner)
	IR, err := parser.Parse()
	if parser.ErrorFound || err != nil {
		fmt.Println("Parse found errors")
		return
	}

	// rename
	largestRegister := parser.GetLargestRegister()
	renamer := renamer.New(largestRegister, IR)
	renamedIR := renamer.Rename()

	if *xFlag {
		fmt.Println(renameIR(renamedIR))
		return
	}

	// if neither the xFlag or the hFlag are provided, schedule on the input file passed in
	scheduler := scheduler.NewSchedule(renamedIR)
	scheduler.PrintSchedule()
}

// helpMessage prints to the command line all the possible commands for the 412fe applications
func helpMessage() {
	fmt.Println("Usage: 412fe [flags]")
	fmt.Println()
	fmt.Println("412fe is the front end for the compiler project.")
	fmt.Println()
	fmt.Println("Flags:")

	fmt.Println("  -h \t\t Display this help message.")
	fmt.Println("  -x <filename>\t scans and parse the input block. It should then perform renaming the code in the input block and print the results to the standard output stream.")

	fmt.Println("No flag provided: behaves as if -p <filename> was specified.")
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

		fmt.Fprintf(&b, op.String()+"\n")
	}

	return b.String()
}

func printSchedule(scheduledBlocks [][]*m.DependenceNode) {
	var b strings.Builder

	for _, blocks := range scheduledBlocks {
		// build each block
		for i, block := range blocks {
			opString := block.Op.String()
			if i == 0 {
				fmt.Fprintf(&b, "[ "+opString)
			} else {
				fmt.Fprintf(&b, " ;  "+opString)
			}
		}

		fmt.Fprintf(&b, " ]\n")
	}

	fmt.Println(b.String())
}

// below is the renamed output for ex1.txt to add two numbers

// loadI 314 => r2
// loadI 0 => r4
// load r4 => r3
// add r2,r3 => r0
// loadI 0 => r1
// store r0 => r1
// output => 0
