package scanner

import "errors"

// This function operates under the state that 'l' and 's' have already been consumed
func (s *scanner) lshiftHelper() (syntacticCategory, error) {
	category := INVALID
	c, _ := s.next()
	if c != 'h' {
		return category, errors.New("invalid lshift instruction: letter 'h' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'i' {
		return category, errors.New("invalid lshift instruction: letter 'i' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'f' {
		return category, errors.New("invalid lshift instruction: letter 'f' expected but found " + string(c))

	}
	c, _ = s.next()
	if c != 't' {
		return category, errors.New("invalid lshift instruction: letter 't' expected but found " + string(c))
	}
	return ARITHOP, nil
}

func (s *scanner) outputHelper() (syntacticCategory, error) {
	category := INVALID
	c, _ := s.next()
	if c != 'u' {
		return category, errors.New("invalid output instruction: letter 'u' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 't' {
		return category, errors.New("invalid output instruction: letter 't' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'p' {
		return category, errors.New("invalid output instruction: letter 'p' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 'u' {
		return category, errors.New("invalid output instruction: letter 'u' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 't' {
		return category, errors.New("invalid output instruction: letter 't' expected but found " + string(c))
	}

	return OUTPUT, nil
}

func (s *scanner) storeHelper() (syntacticCategory, error) {
	category := INVALID
	c, _ := s.next()
	if c != 'o' {
		return category, errors.New("invalid store instruction: letter 'o' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'r' {
		return category, errors.New("invalid store instruction: letter 'r' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'e' {
		return category, errors.New("invalid store instruction: letter 'e' expected but found " + string(c))
	}
	return MEMOP, nil
}

// ALL ATHROP HELPS ARE LISTED BELOW

func (s *scanner) subHelper() (syntacticCategory, error) {
	category := INVALID
	c, _ := s.next()
	if c != 'b' {
		return category, errors.New("invalid sub instruction: letter 'b' expected but found " + string(c))
	}
	return ARITHOP, nil
}

func (s *scanner) addHelper() (syntacticCategory, error) {
	category := INVALID

	c, _ := s.next()
	if c != 'd' {
		return category, errors.New("invalid add instruction: letter 'd' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 'd' {
		return category, errors.New("invalid add instruction: letter 'd' expected but found " + string(c))
	}

	return ARITHOP, nil
}

func (s *scanner) multHelper() (syntacticCategory, error) {
	category := INVALID

	c, _ := s.next()
	if c != 'u' {
		return category, errors.New("invalid mult instruction: letter 'u' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 'l' {
		return category, errors.New("invalid mult instruction: letter 'l' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 't' {
		return category, errors.New("invalid mult instruction: letter 't' expected but found " + string(c))
	}

	return ARITHOP, nil
}

func (s *scanner) rshiftHelper() (syntacticCategory, error) {
	category := INVALID
	c, _ := s.next()
	if c != 'h' {
		return category, errors.New("invalid rshift instruction: letter 'h' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'i' {
		return category, errors.New("invalid rshift instruction: letter 'i' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 'f' {
		return category, errors.New("invalid rshift instruction: letter 'f' expected but found " + string(c))
	}
	c, _ = s.next()
	if c != 't' {
		return category, errors.New("invalid rshift instruction: letter 't' expected but found " + string(c))
	}
	return ARITHOP, nil
}

func (s *scanner) loadHelper() (syntacticCategory, error) {
	category := INVALID

	c, _ := s.next()
	if c != 'a' {
		return category, errors.New("invalid load instruction: letter 'a' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 'd' {
		return category, errors.New("invalid load instruction: letter 'd' expected but found " + string(c))
	}

	// check if the next letter is an I, indicating that the instruction is loadI
	c, err := s.next() // THIS IS GIVING A BUG RN BECAUSE OF INDEXING
	if err != nil {
		// if there is an error, it means we reached the end of the line, so the instruction must be load
		s.curIdx--
		return MEMOP, nil
	}

	if c == 'I' {
		category = LOADI
	} else {
		// if the next letter is not an I, it must be a space, indicating that the instruction is load
		if c != ' ' && c != '\t' {
			return category, errors.New("invalid load instruction: letter 'I' or whitespace expected but found " + string(c))
		}
		s.curIdx-- // step back one character since we read one character too many
		category = MEMOP
	}

	return category, nil

}

func (s *scanner) nopHelper() (syntacticCategory, error) {
	category := INVALID

	c, _ := s.next()
	if c != 'o' {
		return category, errors.New("invalid nop instruction: letter 'o' expected but found " + string(c))
	}

	c, _ = s.next()
	if c != 'p' {
		return category, errors.New("invalid nop instruction: letter 'p' expected but found " + string(c))
	}

	return NOP, nil
}

func (s *scanner) intoHelper() (syntacticCategory, error) {
	category := INVALID

	c, _ := s.next()
	if c != '>' {
		return category, errors.New("invalid 'into' instruction: letter '>' expected but found " + string(c))
	}

	return INTO, nil
}

func (s *scanner) constantHelper(c byte) (syntacticCategory, error) {
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
	s.curIdx-- // step back one character since we read one character too many

	if c != ' ' && c != '\t' && c != '\n' && c != '\r' && c != '=' {
		return category, errors.New("invalid constant: whitespace or end of line expected but found " + string(c))
	}

	return CONSTANT, nil
}

func (s *scanner) commentHelper() (syntacticCategory, error) {
	var err error

	category := INVALID
	c, err := s.next()
	if err != nil {
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

func (s *scanner) registerHelper(c byte) (syntacticCategory, error) {
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
	s.curIdx-- // step back one character since we read one character too many

	if c != ' ' && c != '\t' && c != '\n' && c != '\r' && c != ',' && c != '=' {
		return category, errors.New("invalid constant: whitespace or end of line expected but found " + string(c))
	}

	return REGISTER, nil
}
