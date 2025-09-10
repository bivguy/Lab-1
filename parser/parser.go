package parser

import (
	"container/list"
	"fmt"

	c "github.com/bivguy/Comp412/constants"
	"github.com/bivguy/Comp412/models"
	m "github.com/bivguy/Comp412/models"
)

type parser struct {
	scanner          scanner
	currentOperation m.OperationNode
	operations       *list.List
}

type scanner interface {
	NextToken() (models.Token, error)
	SetNextLine()
}

func New(scanner scanner) *parser {
	return &parser{scanner: scanner, operations: list.New()}
}

func (p *parser) Parse() (*list.List, error) {
	token := p.nextOperationToken()

	// once we get a valid lexeme, start building the internal representation
	p.currentOperation.Line = token.LineNumber
	p.currentOperation.Opcode = token.Lexeme

	// calls the corresponding helper function to finish building its operation
	for token.Category != c.EOF {
		var err error
		switch token.Category {
		case c.MEMOP:
			err = p.finishMemop()

			p.operations.PushFront(p.currentOperation)
			p.currentOperation = m.OperationNode{}
		case c.LOADI:
			err = p.finishLoadI()

			p.operations.PushFront(p.currentOperation)
			p.currentOperation = m.OperationNode{}
		case c.ARITHOP:
			err = p.finishArithop()

			p.operations.PushFront(p.currentOperation)
			p.currentOperation = m.OperationNode{}
		case c.OUTPUT:
			err = p.finishOutput()
			p.operations.PushFront(p.currentOperation)
			p.currentOperation = m.OperationNode{}
		case c.NOP:
			err = p.finishNOP()
			if err != nil {
				return nil, err
			}
			p.operations.PushFront(p.currentOperation)
			p.currentOperation = m.OperationNode{}
		default:
			p.currentOperation = m.OperationNode{}
			err = fmt.Errorf("expected a valid opcode category but instead got %v", token)
		}

		if err != nil {
			fmt.Printf("%v\n", err)
		}

		token = p.nextOperationToken()
	}

	return p.operations, nil

}

// helper function that gets the next token from the scanner that did not return any errors
func (p *parser) nextCorrectToken() m.Token {
	token, err := p.scanner.NextToken()

	// keep printing out the errors of the scanner if they occur
	for err != nil {
		fmt.Printf("error: %v\n", scannerError(err))
		token, err = p.scanner.NextToken()
	}

	return token
}

// skips all comments and EOL tokens and returns a token that can actually be used in the parser
func (p *parser) nextOperationToken() m.Token {
	token := p.nextCorrectToken()

	for token.Category == c.EOL {
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
