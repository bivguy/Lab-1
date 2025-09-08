package scanner

import (
	"errors"

	. "github.com/bivguy/Comp412/constants"
	"github.com/bivguy/Comp412/models"
)

// This function operates under the state that 'l' and 's' have already been consumed
func (s *scanner) lshiftHelper() (models.SyntacticCategory, error) {
	category := INVALID
	c, err := s.next()
	if c != 'h' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid lshift instruction: letter 'h' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'i' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid lshift instruction: letter 'i' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'f' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid lshift instruction: letter 'f' expected but found " + string(c))

	}
	c, _ = s.next()
	if c != 't' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid lshift instruction: letter 't' expected but found " + string(c))
	}

	if !s.checkValidEnding() {
		return category, errors.New("invalid lshift instruction: unexpected character found after 'lshift'")
	}
	return ARITHOP, nil
}

func (s *scanner) outputHelper() (models.SyntacticCategory, error) {
	category := INVALID
	c, err := s.next()
	if c != 'u' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid output instruction: letter 'u' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 't' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid output instruction: letter 't' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'p' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid output instruction: letter 'p' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 'u' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid output instruction: letter 'u' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 't' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid output instruction: letter 't' expected but found " + string(c))
	}

	if !s.checkValidEnding() {
		return category, errors.New("invalid output instruction: unexpected character found after 'output'")
	}

	return OUTPUT, nil
}

func (s *scanner) storeHelper() (models.SyntacticCategory, error) {
	category := INVALID
	c, err := s.next()
	if c != 'o' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid store instruction: letter 'o' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'r' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid store instruction: letter 'r' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'e' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid store instruction: letter 'e' expected but found " + string(c))
	}

	if !s.checkValidEnding() {
		return category, errors.New("invalid store instruction: unexpected character found after 'store'")
	}
	return MEMOP, nil
}

// ALL ATHROP HELPS ARE LISTED BELOW

func (s *scanner) subHelper() (models.SyntacticCategory, error) {
	category := INVALID
	c, err := s.next()
	if c != 'b' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid sub instruction: letter 'b' expected but found " + string(c))
	}

	if !s.checkValidEnding() {
		return category, errors.New("invalid sub instruction: unexpected character found after 'sub'")
	}
	return ARITHOP, nil
}

func (s *scanner) addHelper() (models.SyntacticCategory, error) {
	category := INVALID

	c, err := s.next()
	if c != 'd' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid add instruction: letter 'd' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 'd' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid add instruction: letter 'd' expected but found " + string(c))
	}

	if !s.checkValidEnding() {
		return category, errors.New("invalid add instruction: unexpected character found after 'add'")
	}

	return ARITHOP, nil
}

func (s *scanner) multHelper() (models.SyntacticCategory, error) {
	category := INVALID

	c, err := s.next()
	if c != 'u' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid mult instruction: letter 'u' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 'l' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid mult instruction: letter 'l' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 't' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid mult instruction: letter 't' expected but found " + string(c))
	}

	if !s.checkValidEnding() {
		return category, errors.New("invalid mult instruction: unexpected character found after 'mult'")
	}

	return ARITHOP, nil
}

func (s *scanner) rshiftHelper() (models.SyntacticCategory, error) {
	category := INVALID
	c, err := s.next()
	if c != 'h' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid rshift instruction: letter 'h' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'i' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid rshift instruction: letter 'i' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'f' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid rshift instruction: letter 'f' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 't' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid rshift instruction: letter 't' expected but found " + string(c))
	}

	if !s.checkValidEnding() {
		return category, errors.New("invalid rshift instruction: unexpected character found after 'rshift'")
	}
	return ARITHOP, nil
}

func (s *scanner) loadHelper() (models.SyntacticCategory, error) {
	category := INVALID

	c, err := s.next()
	if c != 'a' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid load instruction: letter 'a' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 'd' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid load instruction: letter 'd' expected but found " + string(c))
	}

	// check if the next letter is an I, indicating that the instruction is loadI
	c, err = s.next() // THIS IS GIVING A BUG RN BECAUSE OF INDEXING
	if err != nil {
		// if there is an error, it means we reached the end of the line, so the instruction must be load
		s.curIdx--
		return MEMOP, nil
	}

	if c == 'I' {
		if !s.checkValidEnding() {
			return category, errors.New("invalid loadI instruction: unexpected character found after 'loadI'")
		}
		category = LOADI
	} else {
		// if the next letter is not an I, it must be a valid ending character
		s.curIdx--
		if !(s.checkValidEnding()) {
			return category, errors.New("invalid load instruction: invalid character after load: " + string(c))
		}
		category = MEMOP
	}

	return category, nil

}

func (s *scanner) nopHelper() (models.SyntacticCategory, error) {
	category := INVALID

	c, err := s.next()
	if c != 'o' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid nop instruction: letter 'o' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 'p' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid nop instruction: letter 'p' expected but found " + string(c))
	}

	if !s.checkValidEnding() {
		return category, errors.New("invalid nop instruction: unexpected character found after 'nop'")
	}

	return NOP, nil
}

func (s *scanner) intoHelper() (models.SyntacticCategory, error) {
	category := INVALID

	c, err := s.next()
	if c != '>' {
		if err != nil {
			s.curIdx--
		}
		return category, errors.New("invalid 'into' instruction: letter '>' expected but found " + string(c))
	}

	return INTO, nil
}

func (s *scanner) constantHelper(c byte) (models.SyntacticCategory, error) {
	var err error

	category := INVALID
	for c >= '0' && c <= '9' {
		c, err = s.next()
		// meaning we reached the end of the line
		if err != nil {
			s.curIdx--
			return CONSTANT, nil
		}
	}

	if c != ' ' && c != '\t' && c != '\n' && c != '\r' && c != '=' && c != '/' {
		return category, errors.New("invalid constant: whitespace or end of line expected but found " + string(c))
	}

	s.curIdx-- // step back one character since we read one character too many

	return CONSTANT, nil
}

func (s *scanner) commentHelper() (models.SyntacticCategory, error) {
	var err error

	category := INVALID
	c, err := s.next()
	if err != nil {
		s.curIdx--
		// if there is an error, comment is invalid
		return category, errors.New("invalid comment: expected another '/' but found end of file")
	}
	if c != '/' {
		return category, errors.New("invalid comment: expected another '/' but found " + string(c))
	}

	// Valid comment, skip to the next line when we want the next token
	s.lineEnd = true

	return COMMENT, nil
}

func (s *scanner) registerHelper(c byte) (models.SyntacticCategory, error) {
	var err error

	category := INVALID
	for c >= '0' && c <= '9' {
		c, err = s.next()
		// meaning we reached the end of the file
		if err != nil {
			s.curIdx--
			return REGISTER, nil
		}
	}

	if c != ' ' && c != '\t' && c != '\n' && c != '\r' && c != ',' && c != '=' && c != '/' {
		return category, errors.New("invalid constant: whitespace or end of line expected but found " + string(c))
	}

	s.curIdx-- // step back one character since we read one character too many

	return REGISTER, nil
}

func (s *scanner) checkValidEnding() bool {
	if s.curIdx >= (len(s.lineText) - 1) {
		return true
	}

	c, err := s.next()

	// valid characters that can come after some token
	// error means we reached the end of the file, while is also valid
	if err != nil || c == '/' || c == ' ' || c == '\t' || c == '\n' || c == '\r' {
		s.curIdx--
		return true
	}

	return false
}
