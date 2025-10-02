package main

import (
	"container/list"
	"flag"
	"fmt"
	"os"

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

			renamer := renamer.New(6, IR)

			renamedIR := renamer.Rename()
			fmt.Println("about to print")
			fmt.Println(PrettyPrintIR(renamedIR))
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
		op := e.Value.(m.OperationNode)
		result += fmt.Sprintf("%s\n\n", op)
	}
	return result
}
