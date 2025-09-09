package parser

import (
	"fmt"
	"strconv"

	c "github.com/bivguy/Comp412/constants"
	m "github.com/bivguy/Comp412/models"
)

type category int

const (
	OPCODE category = iota // 0
	OPONE
	OPTWO
	OPTHREE
	SKIP
)

var memopCategories = []m.SyntacticCategory{c.REGISTER, c.INTO, c.REGISTER, c.EOL}
var memopArgs = []category{OPONE, SKIP, OPTHREE, SKIP}

var loadICategories = []m.SyntacticCategory{c.LOADI, c.CONSTANT, c.INTO, c.REGISTER}
var loadIArgs = []category{OPONE, SKIP, OPTHREE, SKIP}

func (p *parser) finishMemop() error {
	err := p.checkCategories(memopCategories, memopArgs)
	if err != nil {
		return err
	}

	return nil
}

func (p *parser) checkCategories(expectedCategories []m.SyntacticCategory, expectedArgs []category) error {
	for i, cat := range expectedCategories {
		token, err := p.nextToken()
		if err != nil {
			return err
		}
		if token.Category != cat {
			return tokenError(token.Category, cat, token)
		}

		arg := expectedArgs[i]
		err = p.buildOperation(token, arg)
		if err != nil {
			return fmt.Errorf("encountered an error at token %v: %w", token, err)
		}
	}

	return nil
}

func (p *parser) buildOperation(token m.Token, arg category) error {
	if arg == SKIP {
		return nil
	}

	lexeme := token.Lexeme
	SR, err := p.sourceRegisterHelper(lexeme)

	if err != nil {
		return err
	}

	switch arg {
	case OPONE:
		p.currentOperation.OpOne.SR = SR
	case OPTWO:
		p.currentOperation.OpTwo.SR = SR
	case OPTHREE:
		p.currentOperation.OpThree.SR = SR
	}

	return nil
}

// converts the constant or register into an integer added to the SR (Source Register)
func (p *parser) sourceRegisterHelper(lexeme string) (int, error) {
	start := 0
	// check to see if its a register or a constant
	if lexeme[0] == 'r' {
		start = 1
	}

	SR, err := strconv.Atoi(lexeme[start : len(lexeme)+1])
	if err != nil {
		return -1, err
	}

	return SR, nil
}
