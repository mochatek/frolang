package evaluator

import (
	"fmt"

	"github.com/mochatek/frolang/ast"
	"github.com/mochatek/frolang/object"
	"github.com/mochatek/frolang/token"
)

// Constants to save memory
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

// Function to create error object
func newError(format string, rest ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, rest...)}
}

// Function to check whether the supplied object is an error or not
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// Function to evaluate AST to object
// Based on the node's type, call the appropriate evaluator and return the resultant object
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.AssignExpression:
		return evalAssignExpression(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	}
	return nil
}

// Evaluates each statement of the program and returns the final result
// If any of the statement was return statement, then return its return value as final result
// Similarly if we encounter an error object, return the result there itself
// In both cases no further statements will be evaluated
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

// Evaluates the value assigned to an identifier.
// If the evaluation was successful, then set the variable in environment
// If evaluated object was error, then directly return it
func evalLetStatement(LetStatement *ast.LetStatement, env *object.Environment) object.Object {
	value := Eval(LetStatement.Value, env)
	if isError(value) {
		return value
	}
	env.Set(LetStatement.Name.Value, value)
	return nil
}

// Evaluates the return value of a return statement
// If evaluated object was error, then directly return it
func evalReturnStatement(returnStatement *ast.ReturnStatement, env *object.Environment) object.Object {
	returnValue := Eval(returnStatement.ReturnValue, env)
	if isError(returnValue) {
		return returnValue
	}
	return &object.ReturnValue{Value: returnValue}
}

// Evaluates a block statement
// Evaluate each statement in the block
// Return error immediately if any statement evaluated to error
// Return the result immediately if we encounter return statement
// Otherwise return the final result as in parseProgram
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if isError(result) {
			return result
		}
		if result != nil && result.Type() == object.RETURN_OBJ {
			return result
		}
	}
	return result
}

// Evaluates an prefix expression
// If right operand was evaluated to error object, then return it directly
// If the operator is a valid prefix operator, then perform that operation on the right operand and return result
// Otherwise return unknown operator error
func evalPrefixExpression(prefixExpression *ast.PrefixExpression, env *object.Environment) object.Object {
	operand := Eval(prefixExpression.Right, env)
	if isError(operand) {
		return operand
	}
	operator := prefixExpression.Operator

	switch operator {
	case token.MINUS:
		return evalMinusExpression(operand)
	case token.BANG:
		return evalBangExpression(operand)
	default:
		return newError("Unknown operator: %s%s", operator, operand.Type())
	}
}

// Evaluates an infix expression
// If left or right operand was evaluated to error object, then return it directly
// Else perform the operation on the operands and return the result
func evalInfixExpression(infixExpression *ast.InfixExpression, env *object.Environment) object.Object {
	leftOperand := Eval(infixExpression.Left, env)
	if isError(leftOperand) {
		return leftOperand
	}
	rightOperand := Eval(infixExpression.Right, env)
	if isError(rightOperand) {
		return rightOperand
	}
	operator := infixExpression.Operator
	return evalInfixOperation(leftOperand, operator, rightOperand)
}

// Evaluated assignment expression
// Return error if variable is not defined before
// Else, evaluate the value
// If value evaluated to error, then return it
// Else, update value of that variable in env and return the value
func evalAssignExpression(assignExpression *ast.AssignExpression, env *object.Environment) object.Object {
	variable := assignExpression.Variable
	if _, ok := env.Get(variable.Value); !ok {
		return newError("Identifier: %s is not defined at %s", variable.Value, variable.Token.Location)
	}
	value := Eval(assignExpression.Value, env)
	if isError(value) {
		return value
	}
	return env.Update(variable.Value, value)
}

// Evaluates a if expression
// First evaluated the condition
// If evaluated object was error, then directly return it
// If it is true, then return the evaluated result of consequence
// Else if alternate was defined, return its evaluated result
// Otherwise return NULL
func evalIfExpression(ifExpression *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ifExpression.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTrue(condition) {
		return Eval(ifExpression.Consequence, env)
	} else if ifExpression.Alternate != nil {
		return Eval(ifExpression.Alternate, env)
	} else {
		return NULL
	}
}

// Evaluates a for expression
// If object is not iterable, then return error
// Else, provision a local environment
// Get the elements from the iterable object
// Repeatedly evaluate the body length(element) times
// Return error immediately if body evaluates to error or returnValue
// Before each iteration, set the element in the local environment
func evalForExpression(forExpression *ast.ForExpression, env *object.Environment) object.Object {
	iterObject := Eval(forExpression.Iterator, env)
	iterable, ok := iterObject.(object.Iterable)
	if !ok {
		return newError("%s: is not iterable", iterObject.Type())
	}
	elementName := forExpression.Element.Value
	localEnv := object.NewEnclosedEnvironment(env)
	array := iterable.Iter().Elements
	for _, item := range array {
		localEnv.Set(elementName, item)
		result := Eval(forExpression.Body, localEnv)
		if isError(result) {
			return result
		} else if result != nil && result.Type() == object.RETURN_OBJ {
			return result
		}
	}
	return nil
}

// Provision a local environment and start an infinite loop
// Evaluate the condition
// If condition evaluated to an error, then return it immediately
// If condition returned true, then execute body
// Return error immediately if body evaluates to error or returnValue
// If condition returned false, then break from loop
func evalWhileExpression(whileExpression *ast.WhileExpression, env *object.Environment) object.Object {
	localEnv := object.NewEnclosedEnvironment(env)
	for {
		condition := Eval(whileExpression.Condition, localEnv)
		if isError(condition) {
			return condition
		}
		if isTrue(condition) {
			result := Eval(whileExpression.Body, localEnv)
			if isError(result) {
				return result
			} else if result != nil && result.Type() == object.RETURN_OBJ {
				return result
			}
		} else {
			break
		}
	}
	return nil
}

// If left operand and index evaluates to error, then return that error immediately
// Otherwise, based on left and index type, call appropriate evaluator
// Return error if operand is not compatible for index operation
func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(node.Array, env)
	if isError(left) {
		return left
	}
	index := Eval(node.Index, env)
	if isError(index) {
		return index
	}

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("Index operation not supported for: %s[%s]", left.Type(), index.Type())
	}
}

// Return index-th element from the array
// If index exceeded array length, then return NULL
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := len(arrayObject.Elements) - 1

	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}

// Return index-th character from the staring
// If index exceeded string length, then return NULL
func evalStringIndexExpression(str, index object.Object) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := len(strObject.Value) - 1

	if idx < 0 || idx > max {
		return NULL
	}
	return &object.String{Value: string(strObject.Value[idx])}
}

// If index is not hash-able object, return error
// Otherwise, get hash the index and get value from hashPair
// If value was got, then return it. Else, return NULL
func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("Key: %s cannot be hashed", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

// Evaluate the function expression. In case of error, return it
// Otherwise, evaluate the argument list
// If there was only 1 valid argument and it evaluated to error, then return the error
// Otherwise, apply the function on the arguments to get the return value
func evalCallExpression(functionCall *ast.CallExpression, env *object.Environment) object.Object {
	function := Eval(functionCall.Function, env)
	if isError(function) {
		return function
	}

	arguments := evalExpressions(functionCall.Arguments, env)
	if len(arguments) == 1 && isError(arguments[0]) {
		return arguments[0]
	}

	return applyFunction(function, arguments)
}

// Evaluates an array of expressions
// Returns array of evaluated objects as result
// In case of error, returns a single element array with the error object
func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, expression := range expressions {
		evaluated := Eval(expression, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

// If function is user defined
// Then get the local environment for it with all of its argument values set to the parameter identifiers
// Evaluate that function body on this local environment
// Determine the return value and return the result (explicit/implicit return)
// If it was builtin function then call it with the arguments and return the result
// Otherwise return error
func applyFunction(function object.Object, arguments []object.Object) object.Object {
	switch function := function.(type) {
	case *object.Function:
		enclosedEnv := getEnclosedFunctionEnv(function, arguments)
		evaluated := Eval(function.Body, enclosedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return function.Fn(arguments...)
	default:
		return newError("%s: not a function", function.Type())
	}
}

// Creates a local environment for function execution
// The outer of this local env will point to the env in which that function was called
// Sets all the function parameters in this local env, with values as passed in argument list
// Returns the local environment
func getEnclosedFunctionEnv(function *object.Function, arguments []object.Object) *object.Environment {
	enclosedEnv := object.NewEnclosedEnvironment(function.Env)
	for index, parameter := range function.Parameters {
		enclosedEnv.Set(parameter.Value, arguments[index])
	}
	return enclosedEnv
}

// If the value returned was return value object, It means there was an explicit return statement in body
// In that case, return the value of that return object
// Otherwise, return the result itself
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// If the operator is a valid infix operator, then perform that operation on the operands and return result
// Otherwise return unknown operator error
func evalInfixOperation(leftOperand object.Object, operator string, rightOperand object.Object) object.Object {
	switch {
	case operator == token.AND:
		return nativeToBooleanObject(isTrue(leftOperand) && isTrue(rightOperand))
	case operator == token.OR:
		return nativeToBooleanObject(isTrue(leftOperand) || isTrue(rightOperand))
	case operator == token.IN:
		return evalInExpression(leftOperand, rightOperand)
	case (leftOperand.Type() == object.INTEGER_OBJ || leftOperand.Type() == object.FLOAT_OBJ) && (rightOperand.Type() == object.INTEGER_OBJ || rightOperand.Type() == object.FLOAT_OBJ):
		return evalArithmeticExpression(leftOperand, operator, rightOperand)
	case leftOperand.Type() == object.STRING_OBJ && rightOperand.Type() == object.STRING_OBJ:
		return evalStringOperation(leftOperand, operator, rightOperand)
	case operator == token.EQ:
		return nativeToBooleanObject(leftOperand == rightOperand)
	case operator == token.NOT_EQ:
		return nativeToBooleanObject(leftOperand != rightOperand)
	case leftOperand.Type() != rightOperand.Type():
		return newError("Type mismatch: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	default:
		return newError("Unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
}

// If operand is number, do a minus operation and return the result
// Else, return invalid operand error
func evalMinusExpression(operand object.Object) object.Object {
	if operand.Type() == object.INTEGER_OBJ {
		value := operand.(*object.Integer).Value
		return &object.Integer{Value: -value}
	} else if operand.Type() == object.FLOAT_OBJ {
		value := operand.(*object.Float).Value
		return &object.Float{Value: -value}
	} else {
		return newError("Invalid operand: -%s", operand.Type())
	}
}

// Evaluate the operand to boolean and return the negated result
func evalBangExpression(operand object.Object) object.Object {
	return nativeToBooleanObject(!isTrue(operand))
}

// Check left and right operands, perform the appropriate arithmetic operation and return the result
func evalArithmeticExpression(leftOperand object.Object, operator string, rightOperand object.Object) object.Object {
	var resultObject object.Object
	switch {
	case leftOperand.Type() == object.INTEGER_OBJ && rightOperand.Type() == object.INTEGER_OBJ:
		resultObject = evalIntOperation(leftOperand.(*object.Integer), operator, rightOperand.(*object.Integer))
	case leftOperand.Type() == object.FLOAT_OBJ && rightOperand.Type() == object.FLOAT_OBJ:
		resultObject = evalFloatOperation(leftOperand.(*object.Float), operator, rightOperand.(*object.Float))
	case leftOperand.Type() == object.INTEGER_OBJ && rightOperand.Type() == object.FLOAT_OBJ:
		resultObject = evalIntFloatOperation(leftOperand.(*object.Integer), operator, rightOperand.(*object.Float))
	case leftOperand.Type() == object.FLOAT_OBJ && rightOperand.Type() == object.INTEGER_OBJ:
		resultObject = evalFloatIntOperation(leftOperand.(*object.Float), operator, rightOperand.(*object.Integer))
	}
	return resultObject
}

// Return the result of arithmetic operation between two integer operands
func evalIntOperation(leftOperand *object.Integer, operator string, rightOperand *object.Integer) object.Object {
	leftValue := leftOperand.Value
	rightValue := rightOperand.Value

	switch operator {
	case token.PLUS:
		return &object.Integer{Value: leftValue + rightValue}
	case token.MINUS:
		return &object.Integer{Value: leftValue - rightValue}
	case token.ASTERISK:
		return &object.Integer{Value: leftValue * rightValue}
	case token.SLASH:
		return &object.Integer{Value: leftValue / rightValue}
	case token.EQ:
		return nativeToBooleanObject(leftValue == rightValue)
	case token.NOT_EQ:
		return nativeToBooleanObject(leftValue != rightValue)
	case token.LT:
		return nativeToBooleanObject(leftValue < rightValue)
	case token.LT_EQ:
		return nativeToBooleanObject(leftValue <= rightValue)
	case token.GT:
		return nativeToBooleanObject(leftValue > rightValue)
	case token.GT_EQ:
		return nativeToBooleanObject(leftValue >= rightValue)
	default:
		return newError("Unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
}

// Return the result of arithmetic operation between two float operands
func evalFloatOperation(leftOperand *object.Float, operator string, rightOperand *object.Float) object.Object {
	leftValue := leftOperand.Value
	rightValue := rightOperand.Value

	switch operator {
	case token.PLUS:
		return &object.Float{Value: leftValue + rightValue}
	case token.MINUS:
		return &object.Float{Value: leftValue - rightValue}
	case token.ASTERISK:
		return &object.Float{Value: leftValue * rightValue}
	case token.SLASH:
		return &object.Float{Value: leftValue / rightValue}
	case token.EQ:
		return nativeToBooleanObject(leftValue == rightValue)
	case token.NOT_EQ:
		return nativeToBooleanObject(leftValue != rightValue)
	case token.LT:
		return nativeToBooleanObject(leftValue < rightValue)
	case token.LT_EQ:
		return nativeToBooleanObject(leftValue <= rightValue)
	case token.GT:
		return nativeToBooleanObject(leftValue > rightValue)
	case token.GT_EQ:
		return nativeToBooleanObject(leftValue >= rightValue)
	default:
		return newError("Unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
}

// Return the result of arithmetic operation between int & float operands
func evalIntFloatOperation(leftOperand *object.Integer, operator string, rightOperand *object.Float) object.Object {
	leftValue := float64(leftOperand.Value)
	rightValue := rightOperand.Value

	switch operator {
	case token.PLUS:
		return &object.Float{Value: leftValue + rightValue}
	case token.MINUS:
		return &object.Float{Value: leftValue - rightValue}
	case token.ASTERISK:
		return &object.Float{Value: leftValue * rightValue}
	case token.SLASH:
		return &object.Float{Value: leftValue / rightValue}
	case token.EQ:
		return nativeToBooleanObject(leftValue == rightValue)
	case token.NOT_EQ:
		return nativeToBooleanObject(leftValue != rightValue)
	case token.LT:
		return nativeToBooleanObject(leftValue < rightValue)
	case token.LT_EQ:
		return nativeToBooleanObject(leftValue <= rightValue)
	case token.GT:
		return nativeToBooleanObject(leftValue > rightValue)
	case token.GT_EQ:
		return nativeToBooleanObject(leftValue >= rightValue)
	default:
		return newError("Unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
}

// Return the result of arithmetic operation between float & int operands
func evalFloatIntOperation(leftOperand *object.Float, operator string, rightOperand *object.Integer) object.Object {
	leftValue := leftOperand.Value
	rightValue := float64(rightOperand.Value)

	switch operator {
	case token.PLUS:
		return &object.Float{Value: leftValue + rightValue}
	case token.MINUS:
		return &object.Float{Value: leftValue - rightValue}
	case token.ASTERISK:
		return &object.Float{Value: leftValue * rightValue}
	case token.SLASH:
		return &object.Float{Value: leftValue / rightValue}
	case token.EQ:
		return nativeToBooleanObject(leftValue == rightValue)
	case token.NOT_EQ:
		return nativeToBooleanObject(leftValue != rightValue)
	case token.LT:
		return nativeToBooleanObject(leftValue < rightValue)
	case token.LT_EQ:
		return nativeToBooleanObject(leftValue <= rightValue)
	case token.GT:
		return nativeToBooleanObject(leftValue > rightValue)
	case token.GT_EQ:
		return nativeToBooleanObject(leftValue >= rightValue)
	default:
		return newError("Unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
}

// If operator is a valid string operator, then perform the operation and return the result
// Else return unknown operator error
func evalStringOperation(leftOperand object.Object, operator string, rightOperand object.Object) object.Object {
	leftValue := leftOperand.(*object.String).Value
	rightValue := rightOperand.(*object.String).Value

	switch operator {
	case token.PLUS:
		return &object.String{Value: leftValue + rightValue}
	case token.EQ:
		return nativeToBooleanObject(leftValue == rightValue)
	case token.NOT_EQ:
		return nativeToBooleanObject(leftValue != rightValue)
	default:
		return newError("Unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
}

// If rightOperand is not iterable, then return invalid operand error
// If it is hash, then see if leftOperand is hash-able object
// If so, then get the hash key and return presence of the key in hash pairs
// Otherwise, loop through the iterator and evaluate each element == leftOperand
// If it evaluates to true, then return true
func evalInExpression(leftOperand object.Object, rightOperand object.Object) object.Object {
	if iterable, ok := rightOperand.(object.Iterable); ok {
		if hash, ok := iterable.(*object.Hash); ok {
			if key, ok := leftOperand.(object.Hashable); ok {
				if _, exist := hash.Pairs[key.HashKey()]; exist {
					return TRUE
				}
			}
			return FALSE
		}
		for _, element := range iterable.Iter().Elements {
			if evalInfixOperation(leftOperand, token.EQ, element) == TRUE {
				return TRUE
			}
		}
		return FALSE
	}
	return newError("Invalid operand: in %s", rightOperand.Type())
}

// Evaluate all the array elements
// If there was only 1 valid argument and it evaluated to error, then return the err
// Else, create and return Array object
func evalArrayLiteral(array *ast.ArrayLiteral, env *object.Environment) object.Object {
	elements := evalExpressions(array.Elements, env)
	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}
	return &object.Array{Elements: elements}
}

// Create a map - internal data structure for hash
// Loop through each key, value
// If key was evaluated to error/ it is not hash-able, then return error
// Evaluate the value. Return error if it resulted in error
// Otherwise, hash the key and get hashKey
// Add the key, value objects as hash-pair into the map, with hashKey as its key
// Return the hash object
func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("Key: %s cannot be hashed", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

// If identifier is set in environment chain, then return it
// Else, check in built-ins and return it, if present
// Otherwise, return unknown identifier error
func evalIdentifier(identifier *ast.Identifier, env *object.Environment) object.Object {
	if value, ok := env.Get(identifier.Value); ok {
		return value
	}
	if builtin, ok := builtins[identifier.Value]; ok {
		return builtin
	}
	return newError("Identifier: %s not found at %s", identifier.Value, identifier.Token.Location)
}

// Convert boolean value to boolean object
// Useful for reference comparison
func nativeToBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

// Check whether object is having truthy value or not
func isTrue(obj object.Object) bool {
	switch variable := obj.(type) {
	case *object.Boolean:
		return variable.Value
	case *object.Integer:
		if variable.Value != 0 {
			return true
		}
	case *object.Float:
		if variable.Value != 0 {
			return true
		}
	case *object.String:
		if len(variable.Value) > 0 {
			return true
		}
	case *object.Array:
		if len(variable.Elements) > 0 {
			return true
		}
	case *object.Hash:
		if len(variable.Pairs) > 0 {
			return true
		}
	}
	return false
}
