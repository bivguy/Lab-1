package main

import "flag"

func main() {
	hFlag := flag.String("h", "", "Display help")

	sFlag := flag.String("s", "", "Display the Scanner Output")
	pFlag := flag.String("p", "", "Display the Parser Output")
	rFlag := flag.String("r", "", "Display the Intermediate Representation Output")

}
