package scripts

import (
	"errors"
	"fmt"
	"grpg/data-go/grpgitem"
	"grpg/data-go/grpgobj"
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
	Env             *object.Environment
	InteractScripts map[uint16]*ast.BlockStatement
}

func NewScriptManager() *ScriptManager {
	return &ScriptManager{
		InteractScripts: make(map[uint16]*ast.BlockStatement),
		Env:             object.NewEnvironment(),
	}
}

func (s *ScriptManager) LoadObjConstants(objs []grpgobj.Obj) {
	for _, obj := range objs {
		s.Env.Set(uppercaseAll(obj.Name), &object.Integer{Value: int64(obj.ObjId)})
	}
}

func (s *ScriptManager) LoadItemConstants(items []grpgitem.Item) {
	for _, item := range items {
		s.Env.Set(uppercaseAll(item.Name), &object.Integer{Value: int64(item.ItemId)})
	}
}

func (s *ScriptManager) LoadScripts(path string) error {
	env := object.NewEnclosedEnvinronment(s.Env)
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

func uppercaseAll(str string) string {
	chars := []int32(str)

	for i, b := range str {
		if b >= 'a' && b <= 'z' {
			chars[i] = b - 32
		}
	}
	return string(chars)
}
