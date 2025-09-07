package main

import (
	"flag"
	"fmt"
	"os"

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
	} else if *rFlag != "" {

	} else if *pFlag != "" {

	} else if *sFlag != "" { // opening a file and outputting all the results of the scanner
		file, err := os.Open(*sFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
			return
		}
		defer file.Close()
		scanner := scanner.New(file)
		for {
			tok, err := scanner.NextToken()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unexpected error: %v\n", err)
				break
			}
			if tok.Category == EOF {
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
