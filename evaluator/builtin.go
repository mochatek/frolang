package evaluator

import (
	"fmt"
	"strings"

	"github.com/mochatek/frolang/object"
)

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

// Separate Dictionary to support builtin methods
var builtins = map[string]object.Object{
	"print": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
		items := []string{}
		for _, argument := range arguments {
			items = append(items, argument.Inspect())
		}
		fmt.Println(strings.Join(items, " "))
		return nil
	}},
	"type": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
		if len(arguments) != 1 {
			return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
		}
		return &object.String{Value: string(arguments[0].Type())}
	}},
	"str": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
		if len(arguments) != 1 {
			return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
		}
		return &object.String{Value: arguments[0].Inspect()}
	}},
	"len": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
		if len(arguments) != 1 {
			return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
		}
		switch arg := arguments[0].(type) {
		case *object.String:
			return &object.Number{Value: len(arg.Value)}
		case *object.Array:
			return &object.Number{Value: len(arg.Elements)}
		case *object.Hash:
			return &object.Number{Value: len(arg.Pairs)}
		default:
			return newError("Cannot calculate len for argument of type %s", arguments[0].Type())
		}
	}},
	"reversed": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
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
	}},
	"slice": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
		if len(arguments) != 3 {
			return newError("Wrong number of arguments. Got=%d want=3", len(arguments))
		}
		if arguments[0].Type() != object.ARRAY_OBJ && arguments[0].Type() != object.STRING_OBJ {
			return newError("Cannot perform slice on argument of type %s", arguments[0].Type())
		}
		iterable := arguments[0].(object.Iterable)
		if arguments[1].Type() != arguments[2].Type() || arguments[1].Type() != object.NUMBER_OBJ {
			return newError("Start and End values should be NUMBER. Got=%s, %s", arguments[1].Type(), arguments[2].Type())
		}

		length := _len(iterable)
		start := arguments[1].(*object.Number).Value
		end := min(arguments[2].(*object.Number).Value, length)
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
	}},
	"range": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
		if len(arguments) != 2 {
			return newError("Wrong number of arguments. Got=%d want=2", len(arguments))
		}
		if arguments[0].Type() != arguments[1].Type() || arguments[0].Type() != object.NUMBER_OBJ {
			return newError("Argument to range must be NUMBER. Got %s", arguments[0].Type())
		}
		start := arguments[0].(*object.Number).Value
		end := arguments[1].(*object.Number).Value
		if end < start {
			return newError("Need (end >= start). Got start=%d end=%d", start, end)
		}
		elements := make([]object.Object, end-start, end-start)
		for idx, _ := range elements {
			elements[idx] = &object.Number{Value: start}
			start++
		}
		return &object.Array{Elements: elements}
	}},
	"lower": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
		if len(arguments) != 1 {
			return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
		}
		if arguments[0].Type() != object.STRING_OBJ {
			return newError("Argument to lower must be STRING. Got %s", arguments[0].Type())
		}
		str := arguments[0].(*object.String)
		return &object.String{Value: strings.ToLower(str.Value)}
	}},
	"upper": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
		if len(arguments) != 1 {
			return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
		}
		if arguments[0].Type() != object.STRING_OBJ {
			return newError("Argument to upper must be STRING. Got %s", arguments[0].Type())
		}
		str := arguments[0].(*object.String)
		return &object.String{Value: strings.ToUpper(str.Value)}
	}},
	"split": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
		if len(arguments) != 1 {
			return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
		}
		if arguments[0].Type() != object.STRING_OBJ {
			return newError("Argument to split must be STRING. Got %s", arguments[0].Type())
		}
		str := arguments[0].(*object.String)
		array := str.Iter()
		return &array
	}},
	"join": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
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
	}},
	"push": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
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
	}},
	"pop": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
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
	}},
	"unshift": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
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
	}},
	"shift": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
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
	}},
	"keys": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
		if len(arguments) != 1 {
			return newError("Wrong number of arguments. Got=%d want=1", len(arguments))
		}
		if arguments[0].Type() != object.HASH_OBJ {
			return newError("Argument to keys must be HASH. Got %s", arguments[0].Type())
		}
		hash := arguments[0].(*object.Hash)
		array := hash.Iter()
		return &array
	}},
	"values": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
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
	}},
	"delete": &object.Builtin{Fn: func(arguments ...object.Object) object.Object {
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
	}},
}
