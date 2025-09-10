package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	. "github.com/bivguy/Comp412/constants"
	"github.com/bivguy/Comp412/models"
)

type scanner struct {
	curIdx     int
	startIdx   int
	lineReader *bufio.Reader
	lineText   string
	lineLength int
	lineEnd    bool
	lineNumber int
}

func New(file *os.File) *scanner {
	lineReader := bufio.NewReader(file)

	return &scanner{curIdx: -1, startIdx: -1, lineNumber: 0, lineReader: lineReader, lineEnd: true}
}

func (s *scanner) NextToken() (models.Token, error) {
	category := INVALID
	var lexeme string
	var err error
	// initialize the scanner state for the new line if we need to do so (indicated by lineEnd being true)
	if s.lineEnd {
		hasLine, err := s.initLine()
		// If there is no new line to process and we reached the end of the file, return EOF token
		if !hasLine && err == io.EOF {
			return models.Token{Category: EOF, Lexeme: "", LineNumber: s.lineNumber}, nil
		}

		// If there is no new line to process but we did not reach the end of the file, return an error
		if err != nil && err != io.EOF {
			s.lineEnd = true
			return models.Token{Category: INVALID, Lexeme: "", LineNumber: s.lineNumber}, err
		}
	}

	s.skipWhitespace()

	// extract the character at the beginning of the current lexeme
	c, err := s.next()
	// if there is an error, it means we reached the end of the file
	if err != nil {
		return models.Token{Category: EOF, Lexeme: "", LineNumber: s.lineNumber}, nil
	}
	s.startIdx = s.curIdx // mark the beginning of the current lexeme

	// figure out what to do with the character
	switch {
	// if it starts with an l, the word can be load, loadI, or lshift
	case c == 'l':
		c, err = s.next()
		if err != nil {
			s.lineEnd = true
			return models.Token{Category: INVALID, Lexeme: "", LineNumber: s.lineNumber}, err
		}
		switch {
		// if the next letter is an o, it must be load or loadI
		case c == 'o':
			category, err = s.loadHelper()
		// if the next letter is an s, it must be lshift
		case c == 's':
			category, err = s.lshiftHelper()
		default:
			err = errors.New("invalid instruction: letter 'o' or 's' expected but found " + string(c))
		}
	case c == '\r': // TODO: this operates under the assumption that a \n will follow and that this is only in windows
		category = EOL
		s.lineEnd = true
	case c == '\n':
		category = EOL
		s.lineEnd = true
	// if it starts with an s, indicates store or sub
	case c == 's':
		c, err = s.next()
		if err != nil {
			s.lineEnd = true
			return models.Token{Category: INVALID, Lexeme: "", LineNumber: s.lineNumber}, err
		}
		switch {
		// if the next letter is a t, it must be store
		case c == 't':
			category, err = s.storeHelper()
		// if the next letter is a u, it must be sub
		case c == 'u':
			category, err = s.subHelper()
		default:
			err = fmt.Errorf("invalid instruction: letter 't' or 'u' expected but found %q at line %d", c, s.lineNumber)
		}
	case c == 'a':
		category, err = s.addHelper()
	case c == 'm':
		category, err = s.multHelper()
	case c == 'o':
		category, err = s.outputHelper()
	case c == 'n':
		category, err = s.nopHelper()
	case c == ',':
		category = COMMA
	case c == '=':
		category, err = s.intoHelper()
	case c == '/':
		category, err = s.commentHelper()
		if err != nil {
			s.lineEnd = true
			return models.Token{Category: INVALID, Lexeme: "", LineNumber: s.lineNumber}, err
		}

		return models.Token{Category: category, Lexeme: "//", LineNumber: s.lineNumber}, nil
	// if it starts with an r, it indicates a register or rshift
	case c == 'r':
		c, err = s.next()
		if err != nil {
			s.lineEnd = true
			return models.Token{Category: INVALID, Lexeme: "", LineNumber: s.lineNumber}, err
		}

		switch {
		// if the next letter is an s, it must be rshift
		case c == 's':
			category, err = s.rshiftHelper()
		// if the next letter is a digit, it must be a register
		case c >= '0' && c <= '9':
			category, err = s.registerHelper(c)
		}
	default:
		// check if the starting character is a integer (for constants)
		if c >= '0' && c <= '9' {
			category, err = s.constantHelper(c)
		} else {
			err = fmt.Errorf("unrecognized instruction: %q at line %d", c, s.lineNumber)
		}
	}

	if err != nil {
		s.lineEnd = true
		// return token{Category: category, Lexeme: ""}, err
	}

	lexeme = s.lineText[s.startIdx : s.curIdx+1]
	return models.Token{Category: category, Lexeme: lexeme, LineNumber: s.lineNumber}, err
}

// This function initializes the scanner state for a new line. It returns true if there is a new line to process, false otherwise.
// If there is an error while reading the line, it returns the error. If the end of the file is reached,
func (s *scanner) initLine() (bool, error) {
	line, err := s.lineReader.ReadString('\n')

	s.lineText = line
	s.lineLength = len(s.lineText)
	s.curIdx = -1
	s.startIdx = -1
	s.lineNumber++

	if len(line) == 0 {
		return false, err
	}
	s.lineEnd = false
	return true, err
}

// This function increments the current index and returns the character at that index
func (s *scanner) next() (byte, error) {
	s.curIdx++
	if s.curIdx >= s.lineLength {
		// This signifies we have reached the end of the file prematurely
		return 0, fmt.Errorf("reached end of file prematurely at line %d", s.lineNumber)
	}
	return s.lineText[s.curIdx], nil
}

// This function skips whitespace characters (spaces and tabs) in the current line
func (s *scanner) skipWhitespace() {
	for s.curIdx < s.lineLength {
		c, err := s.next()
		if (c != ' ' && c != '\t') || err != nil {
			s.curIdx--
			break
		}
	}
}

func (s *scanner) PrintToken(token models.Token) {
	fmt.Printf("<%v, %q> at line %d\n", SyntacticCategories[token.Category], token.Lexeme, token.LineNumber)
}

func (s *scanner) SetNextLine() {
	s.lineEnd = true
}
