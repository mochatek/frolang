package ast

import (
	"strings"

	"github.com/mochatek/frolang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Node
	Statements []Statement
}

func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (program *Program) String() string {
	var str strings.Builder
	for _, statement := range program.Statements {
		str.WriteString(statement.String())
	}
	return str.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (letStatement *LetStatement) statementNode()       {}
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }
func (letStatement *LetStatement) String() string {
	var str strings.Builder
	str.WriteString(letStatement.TokenLiteral())
	str.WriteString(" ")
	str.WriteString(letStatement.Name.String())
	str.WriteString(" = ")
	if letStatement.Value != nil {
		str.WriteString(letStatement.Value.String())
	}
	return str.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (returnStatement *ReturnStatement) statementNode()       {}
func (returnStatement *ReturnStatement) TokenLiteral() string { return returnStatement.Token.Literal }
func (returnStatement *ReturnStatement) String() string {
	var str strings.Builder
	str.WriteString(returnStatement.TokenLiteral())
	str.WriteString(" ")
	if returnStatement.ReturnValue != nil {
		str.WriteString(returnStatement.ReturnValue.String())
	}
	return str.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expressionStatement *ExpressionStatement) statementNode() {}
func (expressionStatement *ExpressionStatement) TokenLiteral() string {
	return expressionStatement.Token.Literal
}
func (expressionStatement *ExpressionStatement) String() string {
	var str strings.Builder
	if expressionStatement.Expression != nil {
		str.WriteString(expressionStatement.Expression.String())
	}
	return str.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (blockStatement *BlockStatement) statementNode()       {}
func (blockStatement *BlockStatement) TokenLiteral() string { return blockStatement.Token.Literal }
func (blockStatement *BlockStatement) String() string {
	var str strings.Builder
	str.WriteString("{")
	for _, statement := range blockStatement.Statements {
		str.WriteString("\n")
		str.WriteString(statement.String())
	}
	str.WriteString("\n}")
	return str.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (prefixExpression *PrefixExpression) expressionNode() {}
func (prefixExpression *PrefixExpression) TokenLiteral() string {
	return prefixExpression.Token.Literal
}
func (prefixExpression *PrefixExpression) String() string {
	var str strings.Builder
	str.WriteString(prefixExpression.Operator)
	str.WriteString(prefixExpression.Right.String())
	return str.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (infixExpression *InfixExpression) expressionNode()      {}
func (infixExpression *InfixExpression) TokenLiteral() string { return infixExpression.Token.Literal }
func (infixExpression *InfixExpression) String() string {
	var str strings.Builder
	str.WriteString(infixExpression.Left.String())
	str.WriteString(" ")
	str.WriteString(infixExpression.Operator)
	str.WriteString(" ")
	str.WriteString(infixExpression.Right.String())
	return str.String()
}

type TryExpression struct {
	Token   token.Token
	Try     *BlockStatement
	Catch   *BlockStatement
	Error   *Identifier
	Finally *BlockStatement
}

func (tryExpression *TryExpression) expressionNode()      {}
func (tryExpression *TryExpression) TokenLiteral() string { return tryExpression.Token.Literal }
func (tryExpression *TryExpression) String() string {
	var str strings.Builder
	str.WriteString("try ")
	str.WriteString(tryExpression.Try.String())
	str.WriteString(" catch(")
	str.WriteString(tryExpression.Error.String())
	str.WriteString(") ")
	str.WriteString(tryExpression.Catch.String())
	if tryExpression.Finally != nil {
		str.WriteString(" finally ")
		str.WriteString(tryExpression.Finally.String())
	}
	return str.String()
}

type AssignExpression struct {
	Token    token.Token
	Variable *Identifier
	Value    Expression
}

func (assignExpression *AssignExpression) expressionNode() {}
func (assignExpression *AssignExpression) TokenLiteral() string {
	return assignExpression.Token.Literal
}
func (assignExpression *AssignExpression) String() string {
	var str strings.Builder
	str.WriteString(assignExpression.Variable.String())
	str.WriteString(" = ")
	str.WriteString(assignExpression.Value.String())
	return str.String()
}

type IndexExpression struct {
	Token token.Token
	Array Expression
	Index Expression
}

func (indexExpression *IndexExpression) expressionNode()      {}
func (indexExpression *IndexExpression) TokenLiteral() string { return indexExpression.Token.Literal }
func (indexExpression *IndexExpression) String() string {
	var str strings.Builder
	str.WriteString(indexExpression.Array.String())
	str.WriteString("[")
	str.WriteString(indexExpression.Index.String())
	str.WriteString("]")
	return str.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternate   *BlockStatement
}

func (ifExpression *IfExpression) expressionNode()      {}
func (ifExpression *IfExpression) TokenLiteral() string { return ifExpression.Token.Literal }
func (ifExpression *IfExpression) String() string {
	var str strings.Builder
	str.WriteString(ifExpression.TokenLiteral())
	str.WriteString("(")
	str.WriteString(ifExpression.Condition.String())
	str.WriteString(") ")
	str.WriteString(ifExpression.Consequence.String())
	if ifExpression.Alternate != nil {
		str.WriteString(" else ")
		str.WriteString(ifExpression.Alternate.String())
	}
	return str.String()
}

type ForExpression struct {
	Token    token.Token
	Element  *Identifier
	Iterator Expression
	Body     *BlockStatement
}

func (forExpression *ForExpression) expressionNode()      {}
func (forExpression *ForExpression) TokenLiteral() string { return forExpression.Token.Literal }
func (forExpression *ForExpression) String() string {
	var str strings.Builder
	str.WriteString(forExpression.TokenLiteral())
	str.WriteString("(")
	str.WriteString(forExpression.Element.String())
	str.WriteString(" in ")
	str.WriteString(forExpression.Iterator.String())
	str.WriteString(") ")
	str.WriteString(forExpression.Body.String())
	return str.String()
}

type WhileExpression struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (whileExpression *WhileExpression) expressionNode()      {}
func (whileExpression *WhileExpression) TokenLiteral() string { return whileExpression.Token.Literal }
func (whileExpression *WhileExpression) String() string {
	var str strings.Builder
	str.WriteString(whileExpression.TokenLiteral())
	str.WriteString("(")
	str.WriteString(whileExpression.Condition.String())
	str.WriteString(") ")
	str.WriteString(whileExpression.Body.String())
	return str.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (callExpression *CallExpression) expressionNode()      {}
func (callExpression *CallExpression) TokenLiteral() string { return callExpression.Token.Literal }
func (callExpression *CallExpression) String() string {
	var str strings.Builder
	str.WriteString(callExpression.Function.String())
	str.WriteString("(")
	arguments := []string{}
	for _, argument := range callExpression.Arguments {
		arguments = append(arguments, argument.String())
	}
	str.WriteString(strings.Join(arguments, ", "))
	str.WriteString(")")
	return str.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }
func (identifier *Identifier) String() string       { return identifier.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (integerLiteral *IntegerLiteral) expressionNode()      {}
func (integerLiteral *IntegerLiteral) TokenLiteral() string { return integerLiteral.Token.Literal }
func (integerLiteral *IntegerLiteral) String() string       { return integerLiteral.TokenLiteral() }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (floatLiteral *FloatLiteral) expressionNode()      {}
func (floatLiteral *FloatLiteral) TokenLiteral() string { return floatLiteral.Token.Literal }
func (floatLiteral *FloatLiteral) String() string       { return floatLiteral.TokenLiteral() }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (booleanLiteral *BooleanLiteral) expressionNode()      {}
func (booleanLiteral *BooleanLiteral) TokenLiteral() string { return booleanLiteral.Token.Literal }
func (booleanLiteral *BooleanLiteral) String() string       { return booleanLiteral.TokenLiteral() }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (stringLiteral *StringLiteral) expressionNode()      {}
func (stringLiteral *StringLiteral) TokenLiteral() string { return stringLiteral.Token.Literal }
func (stringLiteral *StringLiteral) String() string       { return stringLiteral.TokenLiteral() }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (arrayLiteral *ArrayLiteral) expressionNode()      {}
func (arrayLiteral *ArrayLiteral) TokenLiteral() string { return arrayLiteral.Token.Literal }
func (arrayLiteral *ArrayLiteral) String() string {
	var str strings.Builder
	str.WriteString("[")
	elements := []string{}
	for _, element := range arrayLiteral.Elements {
		elements = append(elements, element.String())
	}
	str.WriteString(strings.Join(elements, ", "))
	str.WriteString("]")
	return str.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hashLiteral *HashLiteral) expressionNode()      {}
func (hashLiteral *HashLiteral) TokenLiteral() string { return hashLiteral.Token.Literal }
func (hashLiteral *HashLiteral) String() string {
	var str strings.Builder
	str.WriteString("{")
	pairs := []string{}
	for key, value := range hashLiteral.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	str.WriteString(strings.Join(pairs, ", "))
	str.WriteString("}")
	return str.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Name       string
	Parameters []*Identifier
	Body       *BlockStatement
}

func (functionLiteral *FunctionLiteral) expressionNode()      {}
func (functionLiteral *FunctionLiteral) TokenLiteral() string { return functionLiteral.Token.Literal }
func (functionLiteral *FunctionLiteral) String() string {
	var str strings.Builder
	str.WriteString(functionLiteral.Name)
	str.WriteString("fn(")
	parameters := []string{}
	for _, parameter := range functionLiteral.Parameters {
		parameters = append(parameters, parameter.String())
	}
	str.WriteString(strings.Join(parameters, ", "))
	str.WriteString(") ")
	str.WriteString(functionLiteral.Body.String())
	return str.String()
}
