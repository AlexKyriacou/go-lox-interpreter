package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate_ast <output directory>")
		os.Exit(64)
	}

	outputDir := os.Args[1]
	err := defineAst(outputDir, "Expr", "(interface{}, error)", []string{
		"Assign : Token name, Expr value",
		"Binary : Expr left, Token operator, Expr right",
		"Call     : Expr callee, Token paren, []Expr arguments",
		"Get      : Expr object, Token name",
		"Grouping : Expr expression",
		"Literal : interface{} value",
		"Logical  : Expr left, Token operator, Expr right",
		"Set      : Expr object, Token name, Expr value",
		"This     : Token keyword",
		"Unary : Token operator, Expr right",
		"Variable : Token name",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	err = defineAst(outputDir, "Stmt", "error", []string{
		"Block : []Stmt statements",
		"Class      : Token name, []Function methods",
		"Expression : Expr expression",
		"Function   : Token name, []Token params, []Stmt body",
		"If         : Expr condition, Stmt thenBranch, Stmt elseBranch",
		"Print      : Expr expression",
		"Return     : Token keyword, Expr value",
		"Var        : Token name, Expr initializer",
	    "While      : Expr condition, Stmt body",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func defineAst(outputDir string, baseName string, returnType string, types []string) error {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("package main\n")
	file.WriteString("\n")
	file.WriteString("type " + baseName + " interface {\n")
	file.WriteString("\tAccept(visitor " + baseName + "Visitor) " + returnType + "\n")
	file.WriteString("}\n")

	defineVisitor(file, baseName, returnType, types)

	for _, astType := range types {
		typeName := strings.TrimSpace(strings.Split(astType, ":")[0])
		fields := strings.TrimSpace(strings.Split(astType, ":")[1])
		defineType(file, typeName, baseName, returnType, fields)
	}

	return nil
}

func defineVisitor(file *os.File, baseName string, returnType string, types []string) {
	file.WriteString("\n")
	file.WriteString("type " + baseName + "Visitor interface {\n")

	for _, astType := range types {
		typeName := strings.TrimSpace(strings.Split(astType, ":")[0])
		file.WriteString("\tVisit" + typeName + baseName + "(" + strings.ToLower(baseName) + " *" + typeName + ") " + returnType + "\n")
	}

	file.WriteString("}\n")
}

func defineType(file *os.File, typeName, baseName, returnType, fieldList string) {
	file.WriteString("\n")
	file.WriteString("type " + typeName + " struct {\n")
	fields := strings.Split(fieldList, ", ")
	for _, field := range fields {
		attrs := strings.Split(field, " ")
		fieldType := strings.TrimSpace(attrs[0])
		fieldName := strings.TrimSpace(attrs[1])
		file.WriteString("\t" + fieldName + " " + fieldType + "\n")
	}
	file.WriteString("}\n")

	file.WriteString("\nfunc (" + strings.ToLower(string(typeName[0])) + " *" + typeName + ") Accept(visitor " + baseName + "Visitor) " + returnType + " {\n")
	file.WriteString("\t return visitor.Visit" + typeName + baseName + "(" + strings.ToLower(string(typeName[0])) + ")\n")
	file.WriteString("}\n")
}
