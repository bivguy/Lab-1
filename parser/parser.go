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
}

func New(scanner scanner) *parser {
	return &parser{scanner: scanner, operations: list.New()}
}

func (p *parser) Parse() (*list.List, error) {
	token, err := p.nextOperationToken()

	if err != nil {
		return nil, err
	}

	// once we get a valid lexeme, start building the internal representation
	p.currentOperation.Line = token.LineNumber
	p.currentOperation.Opcode = token.Lexeme

	// calls the corresponding helper function to finish building its operation
	for token.Category != c.EOF {
		switch token.Category {
		case c.MEMOP:
			err = p.finishMemop()
			if err != nil {
				return nil, err
			}

			p.operations.PushFront(p.currentOperation)
			p.currentOperation = m.OperationNode{}
		case c.LOADI:
			err = p.finishLoadI()
			if err != nil {
				return nil, err
			}

			p.operations.PushFront(p.currentOperation)
			p.currentOperation = m.OperationNode{}
		case c.ARITHOP:
			err = p.finishArithop()
			if err != nil {
				return nil, err
			}

			p.operations.PushFront(p.currentOperation)
			p.currentOperation = m.OperationNode{}
		case c.OUTPUT:
			err = p.finishOutput()
			if err != nil {
				return nil, err
			}
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
			return nil, fmt.Errorf("expected a valid opcode category but instead got %v", token)
		}

		token, err = p.nextOperationToken()
		if err != nil {
			p.currentOperation = m.OperationNode{}
			return nil, scannerError(err)
		}
	}

	return p.operations, err

}

// helper function that wraps around the scanner's nextToken so that it ommits comments
func (p *parser) nextToken() (m.Token, error) {
	token, err := p.scanner.NextToken()
	if err != nil {
		return m.Token{}, scannerError(err)
	}

	for token.Category == c.COMMENT {
		token, err = p.scanner.NextToken()
		if err != nil {
			return m.Token{}, scannerError(err)
		}
	}

	return token, err
}

func (p *parser) nextOperationToken() (m.Token, error) {
	token, err := p.nextToken()
	if err != nil {
		return m.Token{}, scannerError(err)
	}

	for token.Category == c.EOL {
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
