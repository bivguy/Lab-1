package models

type SyntacticCategory int

type Token struct {
	Category   SyntacticCategory
	Lexeme     string
	LineNumber int
}
