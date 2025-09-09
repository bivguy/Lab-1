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
	operations       m.IntermediateRepresentation
}

type scanner interface {
	NextToken() (models.Token, error)
}

type intermediateRepresentation interface {
	AddOperation()
}

func New(scanner scanner) *parser {
	return &parser{scanner: scanner, operations: list.New()}
}

func (p *parser) Parse() (m.IntermediateRepresentation, error) {
	token, err := p.nextToken()

	if err != nil {
		return nil, err
	}

	// once we get a valid lexeme, start building the internal representation
	p.currentOperation.Line = token.LineNumber
	p.currentOperation.Opcode = token.Lexeme

	for token.Category != c.EOF {
		switch token.Category {
		case c.MEMOP:
			err = p.finishMemop()
		case c.LOADI:
		case c.ARITHOP:
		case c.OUTPUT:
		case c.NOP:
		default:
			p.currentOperation = m.OperationNode{}
			return nil, fmt.Errorf("expected a valid opcode category but instead got %v", token)
		}
		if err != nil {
			p.currentOperation = m.OperationNode{}
			return nil, err
		}
		token, err = p.scanner.NextToken()

		if err != nil {
			p.currentOperation = m.OperationNode{}
			return nil, scannerError(err)
		}
	}
	return nil, err

}

// helper function that wraps around the scanner's nextToken so that it ommits commas
func (p *parser) nextToken() (m.Token, error) {
	var token m.Token
	var err error

	for token.Category != c.COMMENT {
		token, err = p.scanner.NextToken()
		if err != nil {
			return m.Token{}, scannerError(err)
		}
	}

	return token, err
}

// indicates that there is an error with the scanner's nextToken
func scannerError(err error) error {
	return fmt.Errorf("encountered in retrieving next token: %w", err)
}

func tokenError(expected m.SyntacticCategory, recieved m.SyntacticCategory, token m.Token) error {
	return fmt.Errorf("encountered an error at line %d: expected a token of type %v but got one of type %v", token.LineNumber, c.SyntacticCategories[expected], c.SyntacticCategories[recieved])

}
