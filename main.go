package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mochatek/frolang/evaluator"
	"github.com/mochatek/frolang/lexer"
	"github.com/mochatek/frolang/object"
	"github.com/mochatek/frolang/parser"
	"github.com/mochatek/frolang/repl"
)

func main() {
	// If source file path was not passed, then start the REPL
	if len(os.Args) == 1 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	// Read source code from the file into a string
	filePath := os.Args[1]
	if parts := strings.Split(filePath, "."); strings.ToLower(parts[len(parts)-1]) != "fro" {
		fmt.Println("ERROR: Input file is not a valid `.fro` script")
		return
	}
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Print("ERROR: ", err)
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
			fmt.Println("\t" + message + "\n")
		}
	} else {
		env := object.NewEnvironment()
		result := evaluator.Eval(program, env)

		// Show errors if any
		if result != nil {
			fmt.Println(result.Inspect())
		}
	}
}
