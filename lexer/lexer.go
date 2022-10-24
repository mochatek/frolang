package lexer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mochatek/frolang/token"
)

type Lexer struct {
	input        string
	char         byte
	curPosition  int
	peekPosition int
	line         int
	col          int
}

// Constructor function for lexer
// Read once to init lexer fields before we start using it
func New(input string) *Lexer {
	lexer := &Lexer{input: input, line: 1}
	lexer.readChar()
	return lexer
}

// Reads 1 character from input string
// Assign read character to `char`
// Advance position pointers
func (lexer *Lexer) readChar() {
	if lexer.peekPosition >= len(lexer.input) {
		lexer.char = 0 // EOF
	} else {
		lexer.char = lexer.input[lexer.peekPosition]
	}
	lexer.curPosition = lexer.peekPosition
	lexer.peekPosition += 1
	lexer.col += 1
}

// Equate character at peekPosition to what is expected
// Return equated result
func (lexer *Lexer) peekCharIs(expectedChar byte) bool {
	var peekChar byte
	if lexer.peekPosition >= len(lexer.input) {
		peekChar = 0
	} else {
		peekChar = lexer.input[lexer.peekPosition]
	}
	return peekChar == expectedChar
}

// Continue reading characters until assertion on `char` fails
// Returns the read string
func (lexer *Lexer) readAheadIfPeekChar(assert func(char byte) bool) string {
	startIndex := lexer.curPosition
	for assert(lexer.char) {
		lexer.readChar()
	}
	return lexer.input[startIndex:lexer.curPosition]
}

// Read character literal and return it
func (lexer *Lexer) readString() string {
	startIndex := lexer.curPosition + 1
	for {
		lexer.readChar()
		if lexer.char == '"' || lexer.char == 0 {
			break
		}
	}
	return lexer.input[startIndex:lexer.curPosition]
}

// Skip processing whitespace character
// Create token based on `char`
// Advance lexer fields through readChar()
// Return the created token
func (lexer *Lexer) ReadToken() token.Token {
	var tok token.Token
	lexer.skipWhiteSpace()

	location := fmt.Sprintf("%d:%d", lexer.line, lexer.col)

	switch lexer.char {
	case 0:
		tok = createToken(token.EOF, lexer.char, location)
	case '+':
		tok = createToken(token.PLUS, lexer.char, location)
	case '-':
		tok = createToken(token.MINUS, lexer.char, location)
	case '(':
		tok = createToken(token.L_PAREN, lexer.char, location)
	case ')':
		tok = createToken(token.R_PAREN, lexer.char, location)
	case '{':
		tok = createToken(token.L_BRACE, lexer.char, location)
	case '}':
		tok = createToken(token.R_BRACE, lexer.char, location)
	case '[':
		tok = createToken(token.L_BRACKET, lexer.char, location)
	case ']':
		tok = createToken(token.R_BRACKET, lexer.char, location)
	case ',':
		tok = createToken(token.COMMA, lexer.char, location)
	case ';':
		tok = createToken(token.SEMICOLON, lexer.char, location)
	case ':':
		tok = createToken(token.COLON, lexer.char, location)
	case '&':
		tok = createToken(token.AND, lexer.char, location)
	case '|':
		tok = createToken(token.OR, lexer.char, location)
	case '/':
		if lexer.peekCharIs('*') {
			char := lexer.char
			lexer.readChar()
			tok = token.Token{Type: token.O_COMMENT, Literal: string(char) + string(lexer.char), Location: location}
		} else {
			tok = createToken(token.SLASH, lexer.char, location)
		}
	case '*':
		if lexer.peekCharIs('/') {
			char := lexer.char
			lexer.readChar()
			tok = token.Token{Type: token.C_COMMENT, Literal: string(char) + string(lexer.char), Location: location}
		} else {
			tok = createToken(token.ASTERISK, lexer.char, location)
		}
	case '=':
		if lexer.peekCharIs('=') {
			char := lexer.char
			lexer.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(char) + string(lexer.char), Location: location}
		} else {
			tok = createToken(token.ASSIGN, lexer.char, location)
		}
	case '!':
		if lexer.peekCharIs('=') {
			char := lexer.char
			lexer.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(char) + string(lexer.char), Location: location}
		} else {
			tok = createToken(token.BANG, lexer.char, location)
		}
	case '<':
		if lexer.peekCharIs('=') {
			char := lexer.char
			lexer.readChar()
			tok = token.Token{Type: token.LT_EQ, Literal: string(char) + string(lexer.char), Location: location}
		} else {
			tok = createToken(token.LT, lexer.char, location)
		}
	case '>':
		if lexer.peekCharIs('=') {
			char := lexer.char
			lexer.readChar()
			tok = token.Token{Type: token.GT_EQ, Literal: string(char) + string(lexer.char), Location: location}
		} else {
			tok = createToken(token.GT, lexer.char, location)
		}
	case '"':
		tok.Type = token.STRING
		tok.Literal = lexer.readString()
	default:
		if isLetter(lexer.char) {
			word := lexer.readAheadIfPeekChar(isLetter)
			tokenType := resolveType(word) // word is identifier/keyword ?
			tok = token.Token{Type: tokenType, Literal: word, Location: location}
			return tok
		} else if isNumber(lexer.char) {
			number := lexer.readAheadIfPeekChar(isNumber)
			numberType := resolveNumberType(number)
			tok = token.Token{Type: numberType, Literal: number, Location: location}
			return tok
		}
		tok = createToken(token.ILLEGAL, lexer.char, location)
	}

	lexer.readChar()
	return tok
}

// Advance to next character if `char` is whitespace
// Increment line counter if we hit new line character and reset col to 0
func (lexer *Lexer) skipWhiteSpace() {
	for lexer.char != 0 && (lexer.char == ' ' || lexer.char == '\t' || lexer.char == '\r' || lexer.char == '\n') {
		if lexer.char == '\n' {
			lexer.line += 1
			lexer.col = 0
		}
		lexer.readChar()
	}
}

// helper function to create token
func createToken(tokenType token.TokenType, literal byte, location string) token.Token {
	return token.Token{Type: tokenType, Literal: string(literal), Location: location}
}

// Helper function to check for valid character
func isLetter(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}

// Helper function to check for valid digit
func isNumber(char byte) bool {
	return '0' <= char && char <= '9' || char == '.' || char == '-'
}

// Lookup in keyword dictionary to decide whether the supplied string is a keyword/identifier
func resolveType(word string) token.TokenType {
	return token.LookUpKeywords(word)
}

// Helper function to get the appropriate token type for a number string
func resolveNumberType(number string) token.TokenType {
	if _, err := strconv.ParseFloat(number, 64); err != nil {
		return token.ILLEGAL
	}
	if strings.Contains(number, ".") {
		return token.FLOAT
	}
	return token.INTEGER
}
