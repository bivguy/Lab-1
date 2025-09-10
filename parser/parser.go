package parser

import (
	"container/list"
	"fmt"
	"os"

	c "github.com/bivguy/Comp412/constants"
	"github.com/bivguy/Comp412/models"
	m "github.com/bivguy/Comp412/models"
)

type parser struct {
	scanner          scanner
	currentOperation m.OperationNode
	operations       *list.List
	ErrorFound       bool
}

type scanner interface {
	NextToken() (models.Token, error)
	SetNextLine()
	GetCurrentLine() int
}

func New(scanner scanner) *parser {
	return &parser{scanner: scanner, operations: list.New()}
}

func (p *parser) Parse() (*list.List, error) {
	token := p.nextOperationToken()

	// calls the corresponding helper function to finish building its operation
	for token.Category != c.EOF {
		// once we get a valid lexeme, start building the internal representation
		p.currentOperation.Line = token.LineNumber
		p.currentOperation.Opcode = token.Lexeme
		var err error
		switch token.Category {
		case c.MEMOP:
			err = p.finishMemop()
		case c.LOADI:
			err = p.finishLoadI()
		case c.ARITHOP:
			err = p.finishArithop()
		case c.OUTPUT:
			err = p.finishOutput()
		case c.NOP:
			err = p.finishNOP()
		default:
			p.currentOperation = m.OperationNode{}
			err = fmt.Errorf("expected a valid opcode category but instead got %v", token)
		}

		if err != nil {
			wrappedErr := fmt.Errorf("ERROR %d: %w", p.scanner.GetCurrentLine(), err)
			fmt.Fprintln(os.Stderr, wrappedErr)
			p.ErrorFound = true
		} else {
			p.operations.PushBack(p.currentOperation)
		}

		p.currentOperation = m.OperationNode{}

		token = p.nextOperationToken()
	}

	return p.operations, nil
}

// helper function that gets the next token from the scanner that did not return any errors
func (p *parser) nextCorrectToken() m.Token {
	token, err := p.scanner.NextToken()

	// keep printing out the errors of the scanner if they occur
	for err != nil {
		p.ErrorFound = true
		wrappedErr := fmt.Errorf("ERROR %d: %w", p.scanner.GetCurrentLine(), err)
		fmt.Fprintln(os.Stderr, wrappedErr)
		token, err = p.scanner.NextToken()
	}

	return token
}

// TODO: create some helper function that makes the default operand SR values to be -1

// skips all comments and EOL tokens and returns a token that can actually be used in the parser
func (p *parser) nextOperationToken() m.Token {
	token := p.nextCorrectToken()

	for token.Category == c.EOL || token.Category == c.COMMENT {
		token = p.nextCorrectToken()
	}

	return token
}

// indicates that there is an error with the scanner's nextToken
func scannerError(err error) error {
	return fmt.Errorf("scanner encountered in retrieving next token: %w", err)
}

func parserError(expected m.SyntacticCategory, recieved m.SyntacticCategory, token m.Token) error {
	return fmt.Errorf("parser encountered an error at line %d: expected a token of type %v but got one of type %v", token.LineNumber, c.SyntacticCategories[expected], c.SyntacticCategories[recieved])
}
