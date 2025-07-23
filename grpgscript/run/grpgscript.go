package run

import (
	"cmp"
	"grpgscript/evaluator"
	"grpgscript/lexer"
	"grpgscript/object"
	"grpgscript/parser"
	"io"
	"log"
	"os"
)

func RunFile(path string) {
	f, err1 := os.Open(path)
	bytes, err2 := io.ReadAll(f)

	if err := cmp.Or(err1, err2); err != nil {
		log.Fatalf("Failed to run file with path %s %v", path, err)
	}

	Run(string(bytes))
}

func Run(str string) {
	l := lexer.New(str)
	p := parser.New(l)
	env := object.NewEnvironment()

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		PrintParserErrors(os.Stdout, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		if evaluated.Type() != object.NULL_OBJ {
			io.WriteString(os.Stdout, evaluated.Inspect()+"\n")
		}
	}
}

// TODO: dupe from repl
func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
