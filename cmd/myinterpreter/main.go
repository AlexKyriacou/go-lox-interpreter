package main

//go:generate go run ./../tools/generateAst.go ./
//go:generate go fmt

import (
	"fmt"
	"os"
)

var hadError bool = false
var hadRuntimeError bool = false

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(64)
	}

	command := os.Args[1]

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := NewScanner(string(fileContents))
	tokens := scanner.scanTokens()
	parser := NewParser(tokens)
	interpreter := NewInterpreter()
	if command == "tokenize" {
		for _, token := range tokens {
			fmt.Println(token)
		}

		if hadError {
			os.Exit(65)
		}
	} else if command == "run" {
		statements := parser.parse()
		if hadError {
			os.Exit(65)
		}

		var resolver Resolver = *NewResolver(&interpreter)
		resolver.resolveStatements(statements)
		if hadError {
			os.Exit(65)
		}

		interpreter.interpret(statements)
		if hadRuntimeError {
			os.Exit(70)
		}
	} else if command == "parse" {
		statements := parser.parse()
		if hadError {
			os.Exit(65)
		}
		astPrinter := AstPrinter{}
		for _, stmt := range statements {
			fmt.Println(astPrinter.printStmt(stmt))
		}
	} else if command == "evaluate" {
		expr, _ := parser.expression()
		value, err := interpreter.evaluate(expr)
		if err != nil {
			reportRuntimeError(*(err.(*RuntimeError)))
		}
		if hadRuntimeError {
			os.Exit(70)
		}
		fmt.Println(interpreter.stringify(value))
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// report reports an error message at a given line number
// setting hadError to true.
func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	hadError = true
}

func reportRuntimeError(runtimeError RuntimeError) {
	fmt.Fprintf(os.Stderr, "%s\n[line %d]\n", runtimeError.message, runtimeError.token.line)
	hadRuntimeError = true
}
