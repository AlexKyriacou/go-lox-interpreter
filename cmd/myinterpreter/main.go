package main

//go:generate go run ./../tools/generateAst.go ./

import (
	"fmt"
	"os"
)

var hadError bool = false

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
	if command == "tokenize" {
		for _, token := range tokens {
			fmt.Println(token)
		}

		if hadError {
			os.Exit(65)
		}
	} else if command == "parse" {
		parser := NewParser(tokens)
		expression := parser.parse()

		if hadError {
			os.Exit(65)
		}
		astPrinter := AstPrinter{}
		fmt.Println(astPrinter.print(expression))
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
