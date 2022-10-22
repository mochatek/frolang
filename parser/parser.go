package parser

import (
	"fmt"
	"strconv"

	"github.com/mochatek/frolang/ast"
	"github.com/mochatek/frolang/lexer"
	"github.com/mochatek/frolang/token"
)

type (
	prefixParser func() ast.Expression
	infixParser  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer         *lexer.Lexer
	curToken      token.Token
	peekToken     token.Token
	prefixParsers map[token.TokenType]prefixParser
	infixParsers  map[token.TokenType]infixParser
	errors        []string
}

// Precedence scores
const (
	_ int = iota
	LOWEST
	EQUALS
	LESS_GREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

// Operator precedence
var precedenceMap = map[token.TokenType]int{
	token.EQ:        EQUALS,
	token.NOT_EQ:    EQUALS,
	token.AND:       EQUALS,
	token.OR:        EQUALS,
	token.IN:        EQUALS,
	token.LT:        LESS_GREATER,
	token.LT_EQ:     LESS_GREATER,
	token.GT:        LESS_GREATER,
	token.GT_EQ:     LESS_GREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.ASTERISK:  PRODUCT,
	token.SLASH:     PRODUCT,
	token.L_PAREN:   CALL,
	token.L_BRACKET: INDEX,
}

// Constructor function for parser
// Init parser fields before we start using
// Register parser functions for prefix and infix tokens
func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}
	parser.scanToken()
	parser.scanToken()

	parser.prefixParsers = make(map[token.TokenType]prefixParser)
	parser.infixParsers = make(map[token.TokenType]infixParser)

	parser.registerPrefixParser(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefixParser(token.NUMBER, parser.parseNumberLiteral)
	parser.registerPrefixParser(token.STRING, parser.parseStringLiteral)
	parser.registerPrefixParser(token.TRUE, parser.parseBooleanLiteral)
	parser.registerPrefixParser(token.FALSE, parser.parseBooleanLiteral)
	parser.registerPrefixParser(token.FUNCTION, parser.parseFunctionLiteral)
	parser.registerPrefixParser(token.L_BRACKET, parser.parseArrayLiteral)
	parser.registerPrefixParser(token.L_BRACE, parser.parseHashLiteral)
	parser.registerPrefixParser(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefixParser(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefixParser(token.L_PAREN, parser.parseGroupedExpression)
	parser.registerPrefixParser(token.IF, parser.parseIfExpression)
	parser.registerPrefixParser(token.FOR, parser.parseForExpression)

	parser.registerInfixParser(token.PLUS, parser.parseInfixExpression)
	parser.registerInfixParser(token.MINUS, parser.parseInfixExpression)
	parser.registerInfixParser(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfixParser(token.SLASH, parser.parseInfixExpression)
	parser.registerInfixParser(token.EQ, parser.parseInfixExpression)
	parser.registerInfixParser(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfixParser(token.LT, parser.parseInfixExpression)
	parser.registerInfixParser(token.LT_EQ, parser.parseInfixExpression)
	parser.registerInfixParser(token.GT, parser.parseInfixExpression)
	parser.registerInfixParser(token.GT_EQ, parser.parseInfixExpression)
	parser.registerInfixParser(token.AND, parser.parseInfixExpression)
	parser.registerInfixParser(token.OR, parser.parseInfixExpression)
	parser.registerInfixParser(token.IN, parser.parseInfixExpression)
	parser.registerInfixParser(token.L_PAREN, parser.parseCallExpression)
	parser.registerInfixParser(token.L_BRACKET, parser.parseIndexExpression)

	return parser
}

// Registers a prefix parser function for a token
func (parser *Parser) registerPrefixParser(tokenType token.TokenType, parserFunction prefixParser) {
	parser.prefixParsers[tokenType] = parserFunction
}

// Registers an infix parser function for a token
func (parser *Parser) registerInfixParser(tokenType token.TokenType, parserFunction infixParser) {
	parser.infixParsers[tokenType] = parserFunction
}

// Compare current token's type to what is been supplied
func (parser *Parser) curTokenIs(tokenType token.TokenType) bool {
	return parser.curToken.Type == tokenType
}

// Compare peek token's type to what is been supplied
func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return parser.peekToken.Type == tokenType
}

// Returns the precedence score of current token
func (parser *Parser) curPrecedence() int {
	if precedence, ok := precedenceMap[parser.curToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

// Returns the precedence score of peek token
func (parser *Parser) peekPrecedence() int {
	if precedence, ok := precedenceMap[parser.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

// Advances current and peek token
func (parser *Parser) scanToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lexer.ReadToken()
}

// Asserts peek token's type is same as what is expected
// Adds peek error if assertion fails
func (parser *Parser) expectPeek(expectedType token.TokenType) bool {
	if parser.peekTokenIs(expectedType) {
		parser.scanToken()
		return true
	} else {
		parser.peekError(expectedType)
		return false
	}
}

// Returns list of errors discovered while parsing
func (parser *Parser) Errors() []string {
	return parser.errors
}

// Create and add peek error to error list
func (parser *Parser) peekError(expectedType token.TokenType) {
	message := fmt.Sprintf("Expected next token to be %s, got %s instead", expectedType, parser.peekToken.Type)
	parser.errors = append(parser.errors, message)
}

// PROGRAM => STATEMENT[]
// Program is actually a set of statements
// Hence, to parse a program, we need to parse every statement until EOF
// A parser function constructs the abstract syntax tree (AST) for the statement
// Append parsed AST to `Statement` array, if parsing was successful
func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for parser.curToken.Type != token.EOF {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.scanToken()
	}
	return program
}

// STATEMENT => LET / RETURN / EXPRESSION
// Applies parse function to the statement based on current token's type
func (parser *Parser) parseStatement() ast.Statement {
	switch parser.curToken.Type {
	case token.O_COMMENT:
		return parser.parseComment()
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

// /* COMMENT */
// Example: /* This is a comment */
func (parser *Parser) parseComment() ast.Statement {
	for !parser.curTokenIs(token.C_COMMENT) && !parser.curTokenIs(token.EOF) {
		parser.scanToken()
	}
	return nil
}

// LET IDENTIFIER = EXPRESSION
// Example: let language = "FroLang"
func (parser *Parser) parseLetStatement() *ast.LetStatement {
	letStatement := &ast.LetStatement{Token: parser.curToken}
	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}
	letStatement.Name = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}
	parser.scanToken()
	letStatement.Value = parser.parseExpression(LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.scanToken()
	}
	return letStatement
}

// RETURN EXPRESSION
// Example: return 0
func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{Token: parser.curToken}
	parser.scanToken()
	returnStatement.ReturnValue = parser.parseExpression(LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.scanToken()
	}
	return returnStatement
}

// EXPRESSION
// In FroLang, every expression is represented as an expression statement
// The Expression field contains the actual expression
func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expressionStatement := &ast.ExpressionStatement{Token: parser.curToken}
	expressionStatement.Expression = parser.parseExpression(LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.scanToken()
	}
	return expressionStatement
}

// BLOCK_STATEMENT => { STATEMENT[] }
// A block statement is a set of statements enclosed within braces
// Example: { let version = 1; print(version); }
func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStatement := &ast.BlockStatement{Token: parser.curToken}
	blockStatement.Statements = []ast.Statement{}
	parser.scanToken()
	for !parser.curTokenIs(token.R_BRACE) && !parser.curTokenIs(token.EOF) {
		statement := parser.parseStatement()
		if statement != nil {
			blockStatement.Statements = append(blockStatement.Statements, statement)
		}
		parser.scanToken()
	}
	if !parser.curTokenIs(token.R_BRACE) {
		return nil
	} else {
		return blockStatement
	}
}

// EXPRESSION
// Parses an expression using Pratt Parsing
func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParsers[parser.curToken.Type]
	if prefix == nil {
		message := fmt.Sprintf("No prefix parse function registered for %s", parser.curToken.Type)
		parser.errors = append(parser.errors, message)
		return nil
	}
	leftExpression := prefix()

	for !parser.peekTokenIs(token.SEMICOLON) && parser.peekPrecedence() > precedence {
		infix := parser.infixParsers[parser.peekToken.Type]
		if infix == nil {
			return leftExpression
		}
		parser.scanToken()
		leftExpression = infix(leftExpression)
	}
	return leftExpression
}

// PREFIX_EXPRESSION => OPERATOR OPERAND
// Example: -5, !true
func (parser *Parser) parsePrefixExpression() ast.Expression {
	prefixExpression := &ast.PrefixExpression{Token: parser.curToken, Operator: parser.curToken.Literal}
	parser.scanToken()
	prefixExpression.Right = parser.parseExpression(PREFIX)
	return prefixExpression
}

// INFIX_EXPRESSION => OPERAND OPERATOR OPERAND
// Example: 1 + 2
func (parser *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression {
	infixExpression := &ast.InfixExpression{Token: parser.curToken, Left: leftExpression, Operator: parser.curToken.Literal}
	precedence := parser.curPrecedence()
	parser.scanToken()
	infixExpression.Right = parser.parseExpression(precedence)
	return infixExpression
}

// GROUPED_EXPRESSION => ( EXPRESSION )
// A grouped expression is an expression enclosed within parentheses
// Grouped expression will have higher precedence as per our precedence map
// Example: (1 + 2) * 3
func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.scanToken()
	groupedExpression := parser.parseExpression(LOWEST)
	if !parser.expectPeek(token.R_PAREN) {
		return nil
	}
	return groupedExpression
}

// CALL_EXPRESSION => EXPRESSION( ARGUMENT, ARGUMENT, .. )
// Example: print(1, !true)
func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{Token: parser.curToken, Function: function}
	callExpression.Arguments = parser.parseExpressionList(token.R_PAREN)
	return callExpression
}

// IF( CONDITION ) { CONSEQUENCE } <ELSE { ALTERNATE }>
// Else part is optional
// Example: if(age >= 18) { "Adult" } else { "Minor" }
func (parser *Parser) parseIfExpression() ast.Expression {
	ifExpression := &ast.IfExpression{Token: parser.curToken}
	if !parser.expectPeek(token.L_PAREN) {
		return nil
	}
	parser.scanToken()
	ifExpression.Condition = parser.parseExpression(LOWEST)
	if !parser.expectPeek(token.R_PAREN) {
		return nil
	}
	if !parser.expectPeek(token.L_BRACE) {
		return nil
	}
	ifExpression.Consequence = parser.parseBlockStatement()
	if parser.peekTokenIs(token.ELSE) {
		parser.scanToken()
		if !parser.expectPeek(token.L_BRACE) {
			return nil
		}
		ifExpression.Alternate = parser.parseBlockStatement()
	}
	return ifExpression
}

// FOR(ELEMENT IN ITERABLE) { BODY }
// Example: for(num in [1, 2, 3]) { print(num) }
func (parser *Parser) parseForExpression() ast.Expression {
	forExpression := &ast.ForExpression{Token: parser.curToken}
	if !parser.expectPeek(token.L_PAREN) {
		return nil
	}
	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}
	forExpression.Element = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
	if !parser.expectPeek(token.IN) {
		return nil
	}
	parser.scanToken()
	forExpression.Iterable = parser.parseExpression(LOWEST)
	if !parser.expectPeek(token.R_PAREN) {
		return nil
	}
	if !parser.expectPeek(token.L_BRACE) {
		return nil
	}
	forExpression.Body = parser.parseBlockStatement()
	return forExpression
}

// IDENTIFIER
// Identifiers are variable names
// Example: age, first_name
func (parser *Parser) parseIdentifier() ast.Expression {
	identifier := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
	return identifier
}

// NUMBER
// Example: 10
func (parser *Parser) parseNumberLiteral() ast.Expression {
	numberLiteral := &ast.NumberLiteral{Token: parser.curToken}
	value, ok := strconv.Atoi(parser.curToken.Literal)
	if ok != nil {
		message := fmt.Sprintf("Could not parse %q as number", parser.curToken.Literal)
		parser.errors = append(parser.errors, message)
		return nil
	} else {
		numberLiteral.Value = value
		return numberLiteral
	}
}

// STRING
// Example: "FroLang"
func (parser *Parser) parseStringLiteral() ast.Expression {
	stringLiteral := &ast.StringLiteral{Token: parser.curToken, Value: parser.curToken.Literal}
	return stringLiteral
}

// BOOLEAN
// Example: true, false
func (parser *Parser) parseBooleanLiteral() ast.Expression {
	booleanLiteral := &ast.BooleanLiteral{Token: parser.curToken, Value: parser.curTokenIs(token.TRUE)}
	return booleanLiteral
}

// FN( PARAMETER, PARAMETER, ... ) { BODY }
// Example: fn(a, b) { a + b }
func (parser *Parser) parseFunctionLiteral() ast.Expression {
	functionLiteral := &ast.FunctionLiteral{Token: parser.curToken}
	if !parser.expectPeek(token.L_PAREN) {
		return nil
	}
	functionLiteral.Parameters = parser.parseFunctionParameters()
	if !parser.expectPeek(token.L_BRACE) {
		return nil
	}
	functionLiteral.Body = parser.parseBlockStatement()
	return functionLiteral
}

// ARRAY => [ ELEMENT, ELEMENT, ... ]
// Example: [1, "FroLang", true]
func (parser *Parser) parseArrayLiteral() ast.Expression {
	arrayLiteral := &ast.ArrayLiteral{Token: parser.curToken}
	arrayLiteral.Elements = parser.parseExpressionList(token.R_BRACKET)
	return arrayLiteral
}

// HASH => { KEY: VALUE }
// Example: {"language": "FroLang", "version": 1}
func (parser *Parser) parseHashLiteral() ast.Expression {
	hashLiteral := &ast.HashLiteral{Token: parser.curToken}
	hashLiteral.Pairs = make(map[ast.Expression]ast.Expression)
	for !parser.peekTokenIs(token.R_BRACE) {
		parser.scanToken()
		key := parser.parseExpression(LOWEST)
		if !parser.expectPeek(token.COLON) {
			return nil
		}
		parser.scanToken()
		value := parser.parseExpression(LOWEST)
		hashLiteral.Pairs[key] = value
		if !parser.peekTokenIs(token.R_BRACE) && !parser.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !parser.expectPeek(token.R_BRACE) {
		return nil
	}
	return hashLiteral
}

// ITERABLE[INDEX]
// Example: versions[0]
func (parser *Parser) parseIndexExpression(array ast.Expression) ast.Expression {
	indexExpression := &ast.IndexExpression{Token: parser.curToken, Array: array}
	parser.scanToken()
	indexExpression.Index = parser.parseExpression(LOWEST)
	if !parser.expectPeek(token.R_BRACKET) {
		return nil
	}
	return indexExpression
}

// ( EXPRESSION, EXPRESSION )
// Example: (1, true)
func (parser *Parser) parseExpressionList(endToken token.TokenType) []ast.Expression {
	arguments := []ast.Expression{}
	if parser.peekTokenIs(endToken) {
		parser.scanToken()
		return arguments
	}
	parser.scanToken()
	arguments = append(arguments, parser.parseExpression(LOWEST))
	for parser.peekTokenIs(token.COMMA) {
		parser.scanToken()
		parser.scanToken()
		arguments = append(arguments, parser.parseExpression(LOWEST))
	}
	if !parser.expectPeek(endToken) {
		return nil
	}
	return arguments
}

// ( IDENTIFIER, IDENTIFIER )
// Example: (language, version)
func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if parser.peekTokenIs(token.R_PAREN) {
		parser.scanToken()
		return identifiers
	}
	parser.scanToken()
	identifier := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
	identifiers = append(identifiers, identifier)
	for parser.peekTokenIs(token.COMMA) {
		parser.scanToken()
		parser.scanToken()
		identifier := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
		identifiers = append(identifiers, identifier)
	}
	if !parser.expectPeek(token.R_PAREN) {
		return nil
	}
	return identifiers
}
