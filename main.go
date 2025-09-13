package main

import (
	"container/list"
	"flag"
	"fmt"
	"os"

	c "github.com/bivguy/Comp412/constants"
	m "github.com/bivguy/Comp412/models"
	"github.com/bivguy/Comp412/parser"

	"github.com/bivguy/Comp412/scanner"
)

func main() {
	hFlag := flag.Bool("h", false, "Display help")

	sFlag := flag.String("s", "", "Display the Scanner Output")
	pFlag := flag.String("p", "", "Display the Parser Output")
	rFlag := flag.String("r", "", "Display the Intermediate Representation Output")
	flag.Parse()

	// politely report that only a single flag should be passed in
	if flag.NFlag() > 1 {

	}

	if *hFlag {
		helpMessage()
	} else if *rFlag != "" || *pFlag != "" {

		path := *rFlag
		if path == "" {
			path = *pFlag
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to open file: %v\n", err)
			return
		}
		defer file.Close()
		scanner := scanner.New(file)
		parser := parser.New(scanner)
		IR, err := parser.Parse()

		// print the results of the parser's Intermediate Representation in a human readable format
		if *rFlag != "" {
			// we only print the IR there is no error found
			if !parser.ErrorFound {
				fmt.Println(PrettyPrintIR(IR))
			}
		} else { // *p flag
			if parser.ErrorFound || err != nil {
				fmt.Println("Parse found errors")
			} else {
				fmt.Printf("Parse succeeded. Processed %d operations.\n", IR.Len())
			}

		}

	} else if *sFlag != "" { // opening a file and outputting all the results of the scanner
		file, err := os.Open(*sFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to open file: %v\n", err)
			return
		}
		defer file.Close()
		scanner := scanner.New(file)
		for {
			tok, err := scanner.NextToken()
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR %d: Unexpected error: %v\n", tok.LineNumber, err)
				break
			}
			scanner.PrintToken(tok)
			if tok.Category == c.EOF {
				break
			}
		}

	} else { // pflag is the default behavior

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
}

func PrettyPrintIR(ir *list.List) string {
	if ir == nil || ir.Len() == 0 {
		return "[empty IR]"
	}

	result := ""
	for e := ir.Front(); e != nil; e = e.Next() {
		op := e.Value.(m.OperationNode)
		result += fmt.Sprintf("%s\n", op)
	}
	return result
}
