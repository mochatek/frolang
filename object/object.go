package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/mochatek/frolang/ast"
)

const (
	NUMBER_OBJ   = "NUMBER"
	STRING_OBJ   = "STRING"
	BOOLEAN_OBJ  = "BOOLEAN"
	ARRAY_OBJ    = "ARRAY"
	HASH_OBJ     = "HASH"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN_VALUE"
	FUNCTION_OBJ = "FUNCTION"
	ERROR_OBJ    = "ERROR"
	BUILTIN_OBJ  = "BUILTIN"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Iterable interface {
	Iter() Array
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type Number struct {
	Value int
}

func (number *Number) Type() ObjectType { return NUMBER_OBJ }
func (number *Number) Inspect() string  { return fmt.Sprintf("%d", number.Value) }
func (number *Number) HashKey() HashKey {
	return HashKey{Type: number.Type(), Value: uint64(number.Value)}
}

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (boolean *Boolean) Inspect() string  { return fmt.Sprintf("%t", boolean.Value) }
func (boolean *Boolean) HashKey() HashKey {
	var value uint64
	if boolean.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: boolean.Type(), Value: value}
}

type String struct {
	Value string
}

func (str *String) Type() ObjectType { return STRING_OBJ }
func (str *String) Inspect() string  { return str.Value }
func (str *String) HashKey() HashKey {
	hash := fnv.New64a()
	hash.Write([]byte(str.Value))
	return HashKey{Type: str.Type(), Value: hash.Sum64()}
}
func (str *String) Iter() Array {
	array := Array{}
	for _, char := range str.Value {
		array.Elements = append(array.Elements, &String{Value: string(char)})
	}
	return array
}

type Array struct {
	Elements []Object
}

func (array *Array) Type() ObjectType { return ARRAY_OBJ }
func (array *Array) Inspect() string {
	var str strings.Builder
	elements := []string{}
	for _, element := range array.Elements {
		elements = append(elements, element.Inspect())
	}
	str.WriteString("[")
	str.WriteString(strings.Join(elements, ", "))
	str.WriteString("]")
	return str.String()
}
func (array *Array) Iter() Array {
	return *array
}

type Null struct{}

func (null *Null) Type() ObjectType { return NULL_OBJ }
func (null *Null) Inspect() string  { return "null" }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (function *Function) Type() ObjectType { return FUNCTION_OBJ }
func (function *Function) Inspect() string {
	var str strings.Builder
	parameters := []string{}
	for _, parameter := range function.Parameters {
		parameters = append(parameters, parameter.String())
	}
	str.WriteString("fn(")
	str.WriteString(strings.Join(parameters, ", "))
	str.WriteString(")")
	str.WriteString(function.Body.String())
	return str.String()
}

type ReturnValue struct {
	Value Object
}

func (returnValue *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (returnValue *ReturnValue) Inspect() string  { return returnValue.Value.Inspect() }

type Error struct {
	Message string
}

func (err *Error) Type() ObjectType { return ERROR_OBJ }
func (err *Error) Inspect() string  { return "EVAL ERROR: " + err.Message }

type builtinFunction func(arguments ...Object) Object

type Builtin struct {
	Fn builtinFunction
}

func (builtin *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (builtin *Builtin) Inspect() string  { return "Builtin function" }

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (hash *Hash) Type() ObjectType { return HASH_OBJ }
func (hash *Hash) Inspect() string {
	var str strings.Builder
	pairs := []string{}
	for _, pair := range hash.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	str.WriteString("{")
	str.WriteString(strings.Join(pairs, ", "))
	str.WriteString("}")
	return str.String()
}
func (hash *Hash) Iter() Array {
	array := Array{}
	for _, pair := range hash.Pairs {
		array.Elements = append(array.Elements, pair.Key)
	}
	return array
}
