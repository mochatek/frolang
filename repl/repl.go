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

const HEADER = "🐸 FroLang v0.1 REPL"
const PROMPT = ">> "

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
	fmt.Println(HEADER)
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
				io.WriteString(out, "\t"+message+"\n")
			}
			continue
		}

		result := evaluator.Eval(program, env)
		if result != nil {
			io.WriteString(out, result.Inspect()+"\n")
		}
	}
}
