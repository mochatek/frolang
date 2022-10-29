package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mochatek/frolang/evaluator"
	"github.com/mochatek/frolang/lexer"
	"github.com/mochatek/frolang/object"
	"github.com/mochatek/frolang/parser"
	"github.com/mochatek/frolang/repl"
)

var RESET = "\033[0m"
var RED = "\033[31m"
var GREEN = "\033[32m"

func main() {
	// Windows doesn't natively support color in cmd
	if runtime.GOOS == "windows" {
		RESET = ""
		RED = ""
		GREEN = ""
	}

	// If source file path was not passed, then start the REPL
	if len(os.Args) == 1 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	// Read source code from the file into a string
	filePath := os.Args[1]
	if parts := strings.Split(filePath, "."); strings.ToLower(parts[len(parts)-1]) != "fro" {
		fmt.Printf("%sSCRIPT ERROR: %s is not a valid FroLang script.\n\tFile extension should be: .fro%s\n", RED, filePath, RESET)
		return
	}
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("%sSCRIPT ERROR: %s%s\n", RED, err, RESET)
		return
	}
	sourceCode := string(contentBytes)

	// Parse the program
	lex := lexer.New(sourceCode)
	par := parser.New(lex)
	program := par.ParseProgram()

	// Evaluate the AST if there was no errors. Else show errors
	if len(par.Errors()) != 0 {
		for _, message := range par.Errors() {
			fmt.Printf("%sPARSE ERROR: %s%s\n", RED, message, RESET)
		}
	} else {
		env := object.NewEnvironment()
		result := evaluator.Eval(program, env)

		// Show errors/result if any
		if result != nil {
			if result.Type() == object.ERROR_OBJ {
				fmt.Printf("%s%s%s\n", RED, result.Inspect(), RESET)
			} else {
				fmt.Printf("%s%s%s\n", GREEN, result.Inspect(), RESET)
			}
		}
	}
}
