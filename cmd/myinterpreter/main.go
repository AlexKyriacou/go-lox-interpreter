package main

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

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	scanner := NewScanner(string(fileContents))
	scanner.scanTokens()
	for _, token := range scanner.tokens {
		fmt.Println(token)
	}

	if hadError {
		os.Exit(65)
	}
}

// error reports a provided error mesage at a given line number.
func error(line int, message string) {
	report(line, "", message)
}

// report reports an error message at a given line number
// setting hadError to true.
func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	hadError = true
}
