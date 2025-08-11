package scripts

import (
	"errors"
	"fmt"
	"grpgscript/ast"
	"grpgscript/evaluator"
	"grpgscript/lexer"
	"grpgscript/object"
	"grpgscript/parser"
	"grpgscript/perf"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ScriptManager TODO: dubious name
type ScriptManager struct {
	InteractScripts map[uint16]*ast.BlockStatement
}

func NewScriptManager() *ScriptManager {
	return &ScriptManager{
		InteractScripts: make(map[uint16]*ast.BlockStatement),
	}
}

func (s *ScriptManager) LoadScripts(path string) error {
	// TODO: add my game stuff to env
	env := object.NewEnvironment()
	AddListeners(env, s)
	entries, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".grpgscript") {
			fullPath := filepath.Join(path, entry.Name())

			file, err := os.Open(fullPath)
			if err != nil {
				return err
			}
			bytes, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			l := lexer.New(string(bytes))
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				fmt.Printf("Found errors parsing script: %s\n", fullPath)
				for _, msg := range p.Errors() {
					fmt.Println(msg)
				}
				return errors.New("errors parsing scripts")
			}

			perf.ConstFold(program)

			evaluator.Eval(program, env)
		}
	}

	return nil
}
