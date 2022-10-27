package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mochatek/frolang/evaluator"
	"github.com/mochatek/frolang/lexer"
	"github.com/mochatek/frolang/object"
	"github.com/mochatek/frolang/parser"
)

const HEADER = "🐸 FroLang v0.1.0 REPL"
const PROMPT = ">> "

const RESET = "\033[0m"
const RED = "\033[31m"
const GREEN = "\033[32m"

// Creates the global environment
// Enters the loop
// Take input statement form user
// Lexer will tokenize the input
// Parser will read tokens through lexer and constructs the program AST
// If there were any parse errors, we will display it
// Else, evaluator will evaluate the program AST and displays the result
// Ask user for next input
// Ctrl + C input will terminate the loop
func Start(in io.Reader, out io.Writer) {
	fmt.Printf("%s%s%s\n", GREEN, HEADER, RESET)
	fmt.Println(strings.Repeat("-", len(HEADER)-2))

	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		code := scanner.Text()
		lex := lexer.New(code)
		par := parser.New(lex)
		program := par.ParseProgram()

		if len(par.Errors()) != 0 {
			for _, message := range par.Errors() {
				io.WriteString(out, fmt.Sprintf("%sPARSE ERROR: %s%s\n", RED, message, RESET))
			}
			continue
		}

		result := evaluator.Eval(program, env)
		if result != nil {
			if result.Type() == object.ERROR_OBJ {
				io.WriteString(out, fmt.Sprintf("%s%s%s\n", RED, result.Inspect(), RESET))
			} else {
				io.WriteString(out, fmt.Sprintf("%s%s%s\n", GREEN, result.Inspect(), RESET))
			}
		}
	}
}
