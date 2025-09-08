package constants

import "github.com/bivguy/Comp412/models"

const (
	MEMOP    models.SyntacticCategory = iota // 0
	LOADI                                    // 1
	ARITHOP                                  // 2
	OUTPUT                                   // 3
	NOP                                      // 4
	CONSTANT                                 // 5
	REGISTER                                 // 6
	COMMA                                    // 7
	INTO                                     // 8
	EOF                                      // 9
	EOL                                      // 10
	COMMENT                                  // 11; not used explicity in this project, but rather for discarding
	INVALID                                  // 12; not used explicitly in this project, but rather for error handling
)
