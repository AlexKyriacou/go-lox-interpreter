package main

//go:generate go run ./../tools/generateAst.go ./
//go:generate go fmt

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var hadError bool = false
var hadRuntimeError bool = false

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [command] [file]\n", os.Args[0])
		fmt.Println("Commands:")
		fmt.Println("  tokenize   Tokenize the input file and print the tokens")
		fmt.Println("  parse      Parse the input file and print the AST")
		fmt.Println("  run        Run the input file (default command)")
		fmt.Println("If no file is provided, the interpreter will start in REPL mode.")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) > 2 {
		fmt.Fprintf(os.Stderr, "Too many arguments provided\n")
		os.Exit(1)
	}

	var command, filename string
	if len(args) == 1 {
		filename = args[0]
	} else if len(args) == 2 {
		command = args[0]
		filename = args[1]
	}

	if filename == "" {
		runPrompt()
	} else {
		runFile(command, filename)
	}
}

func runFile(command, filename string) {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tokenize":
		tokenize(fileContents)
	case "parse":
		parse(fileContents)
	case "run", "":
		run(string(fileContents))
		if hadError {
			os.Exit(65)
		}
		if hadRuntimeError {
			os.Exit(70)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func tokenize(fileContents []byte) {
	scanner := NewScanner(string(fileContents))
	tokens := scanner.scanTokens()
	for _, token := range tokens {
		fmt.Println(token)
	}
	if hadError {
		os.Exit(65)
	}
}

func parse(fileContents []byte) {
	scanner := NewScanner(string(fileContents))
	tokens := scanner.scanTokens()
	parser := NewParser(tokens)
	statements := parser.parse()
	if hadError {
		os.Exit(65)
	}
	astPrinter := AstPrinter{}
	for _, stmt := range statements {
		fmt.Println(astPrinter.printStmt(stmt))
	}
}

func runPrompt() {
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := input.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		run(line)
		hadError = false
	}
}

func run(line string) {
	scanner := NewScanner(line)
	tokens := scanner.scanTokens()
	parser := NewParser(tokens)
	interpreter := NewInterpreter()
	statements := parser.parse()
	if hadError {
		return
	}
	var resolver Resolver = *NewResolver(&interpreter)
	resolver.resolveStatements(statements)
	if hadError {
		return
	}
	interpreter.interpret(statements)
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
