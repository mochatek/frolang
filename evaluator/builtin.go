package evaluator

import (
	"fmt"
	"strings"

	"github.com/mochatek/frolang/object"
)

const RESET = "\033[0m"
const GREEN = "\033[32m"

// Separate Dictionary to support builtin methods
var builtins = map[string]object.Object{
	"print":    &object.Builtin{Fn: print},
	"type":     &object.Builtin{Fn: typeOf},
	"str":      &object.Builtin{Fn: str},
	"len":      &object.Builtin{Fn: length},
	"reversed": &object.Builtin{Fn: reversed},
	"slice":    &object.Builtin{Fn: slice},
	"range":    &object.Builtin{Fn: rangeOf},
	"lower":    &object.Builtin{Fn: lower},
	"upper":    &object.Builtin{Fn: upper},
	"split":    &object.Builtin{Fn: split},
	"join":     &object.Builtin{Fn: join},
	"push":     &object.Builtin{Fn: push},
	"pop":      &object.Builtin{Fn: pop},
	"unshift":  &object.Builtin{Fn: unShift},
	"shift":    &object.Builtin{Fn: shift},
	"keys":     &object.Builtin{Fn: keys},
	"values":   &object.Builtin{Fn: values},
	"delete":   &object.Builtin{Fn: delete},
}

// Print arguments to stdOut
func print(arguments ...object.Object) object.Object {
	items := []string{}
	for _, argument := range arguments {
		items = append(items, argument.Inspect())
	}
	fmt.Println(GREEN, strings.Join(items, " "), RESET)
	return nil
}

// Returns the type of an identifier
func typeOf(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	return &object.String{Value: string(arguments[0].Type())}
}

// Returns the stringified form of any value
func str(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	return &object.String{Value: arguments[0].Inspect()}
}

// Returns the length of an iterable
func length(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	switch arg := arguments[0].(type) {
	case *object.String:
		return &object.Integer{Value: len(arg.Value)}
	case *object.Array:
		return &object.Integer{Value: len(arg.Elements)}
	case *object.Hash:
		return &object.Integer{Value: len(arg.Pairs)}
	default:
		return newError("Cannot calculate len for argument of type %s", arguments[0].Type())
	}
}

// Returns the reversed form of an array/string
func reversed(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	switch arg := arguments[0].(type) {
	case *object.String:
		runes := []rune(arg.Value)
		length := len(arg.Value)
		for i, j := 0, length-1; i < length/2; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return &object.String{Value: string(runes)}
	case *object.Array:
		length := len(arg.Elements)
		elements := make([]object.Object, length, length)
		copy(elements, arg.Elements)
		for i, j := 0, length-1; i < length/2; i, j = i+1, j-1 {
			elements[i], elements[j] = elements[j], elements[i]
		}
		return &object.Array{Elements: elements}
	default:
		return newError("Cannot reverse value for argument of type %s", arguments[0].Type())
	}
}

// Returns a slice from an array/string
// End index is exclusive
func slice(arguments ...object.Object) object.Object {
	if len(arguments) != 3 {
		return newError("Wrong number of arguments. Got=%d want=3", len(arguments))
	}
	if arguments[0].Type() != object.ARRAY_OBJ && arguments[0].Type() != object.STRING_OBJ {
		return newError("Cannot perform slice on argument of type %s", arguments[0].Type())
	}
	iterable := arguments[0].(object.Iterable)
	if arguments[1].Type() != arguments[2].Type() || arguments[1].Type() != object.INTEGER_OBJ {
		return newError("Start and End values should be INTEGERS. Got=%s, %s", arguments[1].Type(), arguments[2].Type())
	}

	length := _len(iterable)
	start := arguments[1].(*object.Integer).Value
	end := min(arguments[2].(*object.Integer).Value, length)
	if 0 > start || start > length || start > end {
		return newError("For slicing, (0 <= start < length) and (start <= end). Got start=%d, end=%d", start, end)
	}
	var sliced object.Object
	switch arg := iterable.(type) {
	case *object.String:
		sliced = &object.String{Value: string([]rune(arg.Value)[start:end])}
	case *object.Array:
		sliced = &object.Array{Elements: arg.Elements[start:end]}
	}
	return sliced
}

// Returns an array of integers ranging from start and end values
// End index is exclusive
func rangeOf(arguments ...object.Object) object.Object {
	if len(arguments) != 2 {
		return newError("Wrong number of arguments. Got=%d want=2", len(arguments))
	}
	if arguments[0].Type() != arguments[1].Type() || arguments[0].Type() != object.INTEGER_OBJ {
		return newError("Argument to range must be INTEGERS. Got %s", arguments[0].Type())
	}
	start := arguments[0].(*object.Integer).Value
	end := arguments[1].(*object.Integer).Value
	if end < start {
		return newError("Need (end >= start). Got start=%d end=%d", start, end)
	}
	elements := make([]object.Object, end-start, end-start)
	for idx, _ := range elements {
		elements[idx] = &object.Integer{Value: start}
		start++
	}
	return &object.Array{Elements: elements}
}

// Returns the lower case form of a string
func lower(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	if arguments[0].Type() != object.STRING_OBJ {
		return newError("Argument to lower must be STRING. Got %s", arguments[0].Type())
	}
	str := arguments[0].(*object.String)
	return &object.String{Value: strings.ToLower(str.Value)}
}

// Returns the upper case form of a string
func upper(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	if arguments[0].Type() != object.STRING_OBJ {
		return newError("Argument to upper must be STRING. Got %s", arguments[0].Type())
	}
	str := arguments[0].(*object.String)
	return &object.String{Value: strings.ToUpper(str.Value)}
}

// Returns an array of characters in a string
func split(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	if arguments[0].Type() != object.STRING_OBJ {
		return newError("Argument to split must be STRING. Got %s", arguments[0].Type())
	}
	str := arguments[0].(*object.String)
	array := str.Iter()
	return &array
}

// Combine elements in an array to form a string and return it
// Separating character will be comma, if not supplied
func join(arguments ...object.Object) object.Object {
	if 1 > len(arguments) || len(arguments) > 2 {
		return newError("Wrong number of arguments. Got=%d want=(min:1, max: 2)", len(arguments))
	}
	if arguments[0].Type() != object.ARRAY_OBJ {
		return newError("First argument to join must be ARRAY. Got %s", arguments[0].Type())
	}
	array := arguments[0].(*object.Array)
	separator := ", "
	if len(arguments) == 2 {
		if arguments[1].Type() != object.STRING_OBJ {
			return newError("Separator to join must be STRING. Got %s", arguments[0].Type())
		}
		separator = arguments[1].(*object.String).Value
	}
	stringArray := make([]string, len(array.Elements), len(array.Elements))
	for idx, element := range array.Elements {
		stringArray[idx] = element.Inspect()
	}
	return &object.String{Value: strings.Join(stringArray, separator)}
}

// Add elements to the end of an array and return it
func push(arguments ...object.Object) object.Object {
	if len(arguments) < 2 {
		return newError("Wrong number of arguments. Got=%d want=minimum 2", len(arguments))
	}
	if arguments[0].Type() != object.ARRAY_OBJ {
		return newError("First argument to push must be ARRAY. Got %s", arguments[0].Type())
	}
	array := arguments[0].(*object.Array)
	length := len(array.Elements)
	newElements := make([]object.Object, length, length)
	copy(newElements, array.Elements)
	newElements = append(newElements, arguments[1:]...)
	return &object.Array{Elements: newElements}
}

// Remove last element from an array and return it
func pop(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	if arguments[0].Type() != object.ARRAY_OBJ {
		return newError("Argument to pop must be ARRAY. Got %s", arguments[0].Type())
	}
	array := arguments[0].(*object.Array)
	length := len(array.Elements)
	if length == 0 {
		return newError("Cannot pop from an empty array")
	}
	newElements := make([]object.Object, length-1, length-1)
	copy(newElements, array.Elements)
	return &object.Array{Elements: newElements}
}

// Add elements to the beginning of an array and return it
func unShift(arguments ...object.Object) object.Object {
	if len(arguments) < 2 {
		return newError("Wrong number of arguments. Got=%d want=minimum 2", len(arguments))
	}
	if arguments[0].Type() != object.ARRAY_OBJ {
		return newError("First argument to unshift must be ARRAY. Got %s", arguments[0].Type())
	}
	array := arguments[0].(*object.Array)
	length := len(arguments[1:])
	newElements := make([]object.Object, length, length)
	copy(newElements, arguments[1:])
	newElements = append(newElements, array.Elements...)
	return &object.Array{Elements: newElements}
}

// Remove first element from an array and return it
func shift(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	if arguments[0].Type() != object.ARRAY_OBJ {
		return newError("Argument to shift must be ARRAY. Got %s", arguments[0].Type())
	}
	array := arguments[0].(*object.Array)
	length := len(array.Elements)
	if length == 0 {
		return newError("Cannot pop from an empty array")
	}
	newElements := make([]object.Object, length-1, length-1)
	copy(newElements, array.Elements[1:])
	return &object.Array{Elements: newElements}
}

// Returns an array of keys in a hash
func keys(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	if arguments[0].Type() != object.HASH_OBJ {
		return newError("Argument to keys must be HASH. Got %s", arguments[0].Type())
	}
	hash := arguments[0].(*object.Hash)
	array := hash.Iter()
	return &array
}

// Returns an array of values in a hash
func values(arguments ...object.Object) object.Object {
	if len(arguments) != 1 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	if arguments[0].Type() != object.HASH_OBJ {
		return newError("Argument to values must be HASH. Got %s", arguments[0].Type())
	}
	hash := arguments[0].(*object.Hash)
	array := object.Array{}
	for _, pair := range hash.Pairs {
		array.Elements = append(array.Elements, pair.Value)
	}
	return &array
}

// Removes a key-value pair form a hash and return it
func delete(arguments ...object.Object) object.Object {
	if len(arguments) != 2 {
		return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
	}
	if arguments[0].Type() != object.HASH_OBJ {
		return newError("First argument to delete must be HASH. Got %s", arguments[0].Type())
	}
	hash := arguments[0].(*object.Hash)
	if deleteKey, ok := arguments[1].(object.Hashable); ok {
		newHashPairs := make(map[object.HashKey]object.HashPair)
		for key, value := range hash.Pairs {
			if key != deleteKey.HashKey() {
				newHashPairs[key] = value
			}
		}
		return &object.Hash{Pairs: newHashPairs}
	}
	return newError("Key of type %s cannot be hashed", arguments[1].Type())
}

// Helper function to calculate minimum of two numbers
func min(num1, num2 int) int {
	if num1 < num2 {
		return num1
	}
	return num2
}

// Helper function to calculate the length of string/array object
func _len(iterable object.Iterable) int {
	var length int
	switch obj := iterable.(type) {
	case *object.String:
		return len(obj.Value)
	case *object.Array:
		return len(obj.Elements)
	}
	return length
}
