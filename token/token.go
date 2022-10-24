package token

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
	Location string
}

// Identifiers and Literals
const (
	IDENTIFIER = "IDENTIFIER"
	INTEGER    = "INTEGER"
	FLOAT      = "FLOAT"
	TRUE       = "TRUE"
	FALSE      = "FALSE"
	STRING     = "STRING"
)

// Arithmetic Operators
const (
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"
	ASSIGN   = "="
)

// Comparison Operators
const (
	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	LT_EQ  = "<="
	GT     = ">"
	GT_EQ  = ">="
)

// Logical Operators
const (
	AND = "&"
	OR  = "|"
)

// Parentheses, Braces and Special characters
const (
	L_PAREN   = "("
	R_PAREN   = ")"
	L_BRACE   = "{"
	R_BRACE   = "}"
	L_BRACKET = "["
	R_BRACKET = "]"
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	O_COMMENT = "/*"
	C_COMMENT = "*/"
)

// Keywords
const (
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	WHILE    = "WHILE"
	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"
	IN       = "in"
)

// Others
const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

var Keywords = map[string]TokenType{
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"in":     IN,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"while":  WHILE,
	"fn":     FUNCTION,
	"return": RETURN,
}

// Helper function to lookup a word in keyword dictionary
// If word id present, then return its corresponding type
// Else, return identifier type
func LookUpKeywords(word string) TokenType {
	if keywordType, ok := Keywords[word]; ok {
		return keywordType
	} else {
		return IDENTIFIER
	}
}
