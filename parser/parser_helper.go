package parser

import (
	"fmt"
	"math"
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

var memopCategories = []m.SyntacticCategory{c.REGISTER, c.INTO, c.REGISTER, c.EO}
var memopArgs = []category{OPONE, SKIP, OPTHREE, SKIP}

var loadICategories = []m.SyntacticCategory{c.CONSTANT, c.INTO, c.REGISTER, c.EO}
var loadIArgs = []category{OPONE, SKIP, OPTHREE, SKIP}

var arithopCategories = []m.SyntacticCategory{c.REGISTER, c.COMMA, c.REGISTER, c.INTO, c.REGISTER, c.EO}
var arithipArgs = []category{OPONE, SKIP, OPTWO, SKIP, OPTHREE, SKIP}

var outputCategories = []m.SyntacticCategory{c.CONSTANT, c.EO}
var outputArgs = []category{OPTHREE, SKIP}

var nopCategories = []m.SyntacticCategory{c.EO}
var nopArgs = []category{SKIP}

func (p *parser) finishMemop() error {
	err := p.buildCategories(memopCategories, memopArgs)
	if err != nil {
		return err
	}

	return nil
}

func (p *parser) finishLoadI() error {
	err := p.buildCategories(loadICategories, loadIArgs)
	if err != nil {
		return err
	}

	return nil
}

func (p *parser) finishArithop() error {
	err := p.buildCategories(arithopCategories, arithipArgs)
	if err != nil {
		return err
	}

	return nil
}

func (p *parser) finishOutput() error {
	err := p.buildCategories(outputCategories, outputArgs)
	if err != nil {
		return err
	}

	return nil
}

func (p *parser) finishNOP() error {
	err := p.buildCategories(nopCategories, nopArgs)
	if err != nil {
		return err
	}

	return nil
}

func (p *parser) buildCategories(expectedCategories []m.SyntacticCategory, expectedArgs []category) error {
	for i, cat := range expectedCategories {
		token, err := p.scanner.NextToken()
		// scanner error, so we will already be reading the next line at this point
		if err != nil {
			return err
		}
		tokenCat := token.Category
		// check the special case of EO
		if cat == c.EO {
			if tokenCat != c.EOL && tokenCat != c.EOF && tokenCat != c.COMMENT {
				// this makes a new line in the scanner because a parser error does not start a new line in the scanner
				p.scanner.SetNextLine()
				return fmt.Errorf("encountered an error at line %d: expected a token of type EOF or EOF but got one of type %v", token.LineNumber, c.SyntacticCategories[tokenCat])
			}
		} else if tokenCat != cat {
			p.scanner.SetNextLine()
			return parserError(token.Category, cat, token)
		}

		arg := expectedArgs[i]
		err = p.buildOperation(token, arg)
		if err != nil {
			p.scanner.SetNextLine()
			return fmt.Errorf("encountered an error at token %v: %w", token, err)
		}
	}

	return nil
}

// given a token and a category, it adds it to its corresponding place for the operation being built
func (p *parser) buildOperation(token m.Token, arg category) error {
	if arg == SKIP {
		return nil
	}

	lexeme := token.Lexeme
	SR, err := p.sourceRegisterHelper(lexeme)

	if err != nil {
		return err
	}

	var op *m.Operand
	switch arg {
	case OPONE:
		op = &p.currentOperation.OpOne
	case OPTWO:
		op = &p.currentOperation.OpTwo
	case OPTHREE:
		op = &p.currentOperation.OpThree
	}

	op.Active = true

	// set the largest register size
	if lexeme[0] == 'r' {
		if SR > p.largestRegister {
			p.largestRegister = SR
		}
	}

	op.SR = SR
	op.VR = -1
	op.NU = math.Inf(1)

	return nil
}

// converts the constant or register into an integer added to the SR (Source Register)
func (p *parser) sourceRegisterHelper(lexeme string) (int, error) {
	start := 0
	// check to see if its a register or a constant
	if lexeme[0] == 'r' {
		start = 1
	}

	SR, err := strconv.Atoi(lexeme[start:])
	if err != nil {
		return -1, err
	}

	return SR, nil
}

func isRegister(opcode string, cat category) bool {
	if (opcode == "loadI" && cat == OPONE) || (opcode == "output" && cat == OPTHREE) {
		return false
	}

	return true
}
