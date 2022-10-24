# 🐸 FroLang v0.0.1
FroLang is a toy language created for the sole purpose of learning how to build an interpreter and to sharpen my Go skills.

<p align="center">
  <img src="https://github.com/mochatek/frolang/blob/master/logo.png" alt="Logo" />
</p>

__How does it work?__

FroLang interpreter uses an approach called _Tree-Walking_, which parses the source code, builds an abstract syntax tree (AST), and then evaluates this tree. These are the steps involved:

1. __Lexical Analysis__: from source code (free text) to Tokens/Lexemes;
2. __Parsing__: uses the generated tokens to create an Abstract Syntax Tree;
3. __AST Construction__: a structural representation of the source code with its precedence level, associativity, etc;
4. __Evaluation__: runs through the AST evaluating all expressions.

## Installation and Usage
You can follow any of these methods to install and use FroLang:

1. Running from this project:
    - Clone this repo
    - Install [Go](https://go.dev/dl/)
    - Run `go run main.go` for the _REPL_
    - Run `go run main.go fro_script_path` to run a valid _.fro_ script
2. Using _go get_ utility
3. Download the compiled binary from [Releases](https://github.com/mochatek/frolang/releases)

# Features
- [Variables](#variables)
- [Comments](#comments)
- [Primitive Data Types](#primitive-data-types)
  - [Integer](#integer)
  - [Float](#float)
  - [String](#string)
  - [Boolean](#boolean)
- [Container Types](#container-types)
  - [Array](#array)
  - [Hash](#hash)
- [Functions](#functions)
- [Operators](#operators)
  - [Arithmetic operators](#arithmetic-operators)
  - [String operators](#string-operators)
  - [Conditional operators](#conditional-operators)
  - [Logical operators](#logical-operators)
  - [Presence operators](#presence-operators)
- [Conditionals](#conditionals)
- [Loops](#loops)
- [Builtin Methods](#builtin-methods)
- [To-Do](#to-do)

## Variables
- Declare variables using `let` keyword
- Variable name should only contain letters and underscore
- Variable names are case sensitive

**Example**
```js
let language_name = "FroLang";
```

## Comments
In FroLang, you can create single/multi-line comment using `/* */`

**Example**
```js
/* This is a comment */
```

## Primitive Data Types
Following primitive types are available in FroLang:

### Integer
- Whole-valued positive, negative number or zero
- Truthy value: Non zero value

**Example**
```js
let rank = 1;
```

### Float
- Positive or negative whole number with a decimal point
- Truthy value: Non zero value

**Example**
```js
let version = 1.1;
```

### String
- Sequence of characters enclosed in double quotes
- You can access individual character by their index and index starts from 0
- Strings in FroLang are immutable
- Truthy value: Non empty string

**Example**
```js
let author = "MochaTek";
let firstCharacter = author[0];
```

### Boolean
- Represents truth value; ie, true / false
- Truthy value: true

**Example**
```js
let completed = false;
```

## Container Types
Following container types are available in FroLang:

### Array
- Represents a collection of elements
- Elements of a FroLang array can be of completely different types
- Multi-dimensional arrays are also supported
- Array elements are ordered by their index and index starts from 0
- Arrays in FroLang are immutable
- Truthy value: Non empty array (contains at least 1 element)

**Example**
```js
let items = [1, 2.5, true, "Go", [1, 2]];
let lastItem = items[4][1];
```

### Hash
- Represents dictionary that can store key-value pairs
- Keys of a hash must be of primitive type (hash-able)
- Keys of a hash is unordered
- Values can be of any type
- Retrieve value from a hash using the key as the index
- Truthy value: Non empty hash (contains at least 1 key)

**Example**
```js
let passwordDict = {"gmail": 123, "fb": 456};
let fbPassword = passwordDict["fb"];
```

## Functions
- Functions in FroLang are fist class citizens
- Functions are created using `fn` keyword
- FroLang as of now doesn't support default arguments
- Functions in FroLang does create `closures`
- Functions in froLang implicitly returns the value of last statement
- You can explicitly return from anywhere within the body using `return` keyword

**Example**
```js
let speak = fn(prefix) {
    let sep = ">>";
    return fn(message) { prefix + sep + message }
};

print(speak("Bot")("Hello World"));
```

## Operators
Following operators are supported by FroLang:

### Arithmetic operators
| Operator | Description | Operands | Example |
|-|-|-|-|
|__+__|Sum|integer/float|`let sum = 3 + 1;`|
|__-__|Subtract|integer/float|`let diff = 3 - 1;`|
|__*__|Multiply|integer/float|`let prod = 3 * 2;`|
|__/__|Divide|integer/float|`let quot = 6 / 2;`|
> 💡In case of arithmetic operation, if any of the operand is having float value, then the result of the operation will also be a float value 

### String operators
| Operator | Description | Operands | Example |
|-|-|-|-|
|__+__|Concatenate|string|`let msg = "Mocha" + "Tek";`|

### Conditional operators
| Operator | Description | Operands | Example |
|-|-|-|-|
|__<__|Less than|int/float|`let res = 3 < 4;`|
|__>__|Greater than|int/float|`let res = 3 > 4;`|
|__<=__|Less than or equal|int/float|`let res = 3 <= 4;`|
|__>=__|Less than or equal|int/float|`let res = 3 >= 4;`|
|__==__|Equality|any|`let res = 3 == 4;`|
|__!=__|Inequality|any|`let res = 3 != 4;`|

> 💡== and != returns boolean value. It compares the value of operands in case of primitive types, whereas it compares the reference in case of containers. Therefore, [1, 2] will not be equal to [1, 2]

> 💡Comparison like: (2.0 == 2) will evaluate to true, but (2.1 == 2) will not

### Logical operators
| Operator | Description | Operands | Example |
|-|-|-|-|
|__!__|Not - negates the truth value|any|`let res = !true;`|
|__&__|And - evaluate to true if both operands are true|any|`let res = true & true;`|
|__\|__|Or - evaluate to true if any operand is true|any|`let res = true \| false;`|

### Presence operators
| Operator | Description | Operands | Example |
|-|-|-|-|
|__in__|Check if an element exists in a sequence or not|string/array/hash|`let res = "a" in "FroLang";`|

> 💡In case of hash, the _in_ operator looks for the key rather than value as in string/array

## Conditionals
- FroLang only has if and else. It doesn't have any elif or else if like in other languages
- In FroLang, you can use `if - else` as an expression to mimic a ternary operation

**Example**
```js
let name = "Ronaldinho";
let goals = 5;
let rank = if(goals >= 3) { 1 } else { 2 };

if(rank == 1) {
    print(name + " won 🥇");
} else {
    if(goals) {
        print(name, " won 🥈");
    } else {
        print(name + " won 🥉");
    }
};
```

## Loops
- FroLang as of now, only has the `for` loop for iteration
- You can use the `in` operator to iterate through each element of a sequence
- In case of hash, iterating element is the key
- For looping _n_ times, you can use `range(start, end)` to create a sequence of length: n

**Example**
```js
let string = "FroLang";
for (char in string) {
    print(char);
}

let array = ["Pen", 5];
for (element in array) {
    print(element);
}

let hash = {"fb": 123, true: "Valid!"};
for (key in hash) {
    print(key, hash[key]);
}

/* n = 2 */
for (i in range(1, 3)) {
    print("count", i);
}
```

## Builtin Methods
|Method|Description|Example|
|-|-|-|
|_print(...args)_|Prints arguments to stdout separated by space|`print("Hello ", "World")`|
|_type(arg)_|Returns the type of the argument|`type(1)`|
|_str(arg)_|Returns the stringified form of the argument|`str([1, 2])`|
|_len(iterable)_|Returns the length of a string/array/hash|`len("FroLang")`|
|_reversed(str_or_array)_|Reverse the order of elements in a string/array|`reversed("FroLang")`|
|_slice(str_or_array, start, end)_|Returns a slice from start to end index of a string/array. End index is exclusive|`slice("MochaTek", 0, 5)`|
|_range(start, end)_|Returns an integer array with elements ranging from start to end. End is exclusive|`range(0, 5)`|
|_lower(str)_|Returns the lower case representation of a string|`lower("HeLlO")`|
|_upper(str)_|Returns the upper case representation of a string|`upper("HeLlO")`|
|_split(str)_|Returns an array with characters of a string as elements|`split("FroLang")`|
|_join(array, sep=", ")_|Returns a string created by combing array elements separated by _sep_, which is _", "_ by default|`join(["F", "r", "o", "L", "a", "n", "g"], "")`|
|_push(array, ...elements)_|Returns a new array with elements inserted at the end|`push([1, 2], 3, 4)`|
|_pop(array)_|Returns a new array with the last element removed|`pop([1, 2, 3])`|
|_unshift(array, ...elements)_|Returns a new array with elements inserted at the beginning|`unshift([3, 4], 1, 2)`|
|_shift(array)_|Returns a new array with the first element removed|`shift([1, 2, 3])`|
|_keys(hash)_|Returns an array of keys in a hash|`keys({1: "one", "two": 2})`|
|_values(hash)_|Returns an array of values in a hash|`values({1: "one", "two": 2})`|
|_delete_(hash, key)_|Returns a new hash with the key-value pair removed for the supplied key|`delete({1: "one", "two": 2}, 1)`|

## To-Do
- [ ] Assignment expression `version = version + 1`
- [ ] Controls statements: `break, continue`
- [ ] While loop
- [ ] Error handling: `try - catch`
- [ ] Environment variables
- [ ] Modules
- [ ] StdLib: `datetime, fileIO`
- [ ] Help
- [ ] Example programs
- [ ] Docker image
- [ ] Compiler

# Reference

[Writing An Interpreter In Go](https://interpreterbook.com/) by Thorsten Ball
