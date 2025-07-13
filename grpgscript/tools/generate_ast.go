package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_ast <output directory?")
		os.Exit(64)
	}

	outputDir := os.Args[1]

	defineAst(outputDir, "Expr", []string{
		"Binary   : Left Expr, Operator lex.Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value any",
		"Unary    : Operator lex.Token, Right Expr",
	})
}

func defineAst(outputDir string, baseName string, data []string) {
	file, err := os.Create(outputDir + strings.ToLower(baseName) + ".go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString("package ast\n\n")
	file.WriteString("import \"grpgscript/lex\"\n\n")
	file.WriteString("type Expr interface {}\n")

	for _, dataType := range data {
		split := strings.Split(dataType, ":")
		name := strings.TrimSpace(split[0])
		fields := strings.TrimSpace(split[1])
		defineType(file, baseName, name, fields)
	}
}

func defineType(file *os.File, baseName, name, fields string) {
	file.WriteString("\n")
	fmt.Fprintf(file, "type %s struct {\n", name)

	fieldsList := strings.Split(fields, ", ")
	for _, field := range fieldsList {
		split := strings.Split(field, " ")
		fmt.Fprintf(file, "    %s %s\n", split[0], split[1])
	}

	file.WriteString("}\n")
}
