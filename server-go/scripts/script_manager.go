package scripts

import (
	"errors"
	"grpg/data-go/grpgitem"
	"grpg/data-go/grpgnpc"
	"grpg/data-go/grpgobj"
	"grpgscript/ast"
	"grpgscript/evaluator"
	"grpgscript/lexer"
	"grpgscript/object"
	"grpgscript/parser"
	"grpgscript/perf"
	"log"
	"os"
	"path/filepath"
	"server/shared"
	"strings"
)

// ScriptManager TODO: dubious name
type ScriptManager struct {
	Env             *object.Environment
	InteractScripts map[uint16]*ast.BlockStatement
	NpcTalkScripts  map[uint16]*ast.BlockStatement
	TimedScripts    map[uint32][]TimedScript
}

func (s *ScriptManager) AddTimedScript(tick uint32, script TimedScript) {
	_, ok := s.TimedScripts[tick]
	if !ok {
		s.TimedScripts[tick] = []TimedScript{script}
	} else {
		s.TimedScripts[tick] = append(s.TimedScripts[tick], script)
	}
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

func (s *ScriptManager) LoadNpcConstants(npcs map[uint16]*grpgnpc.Npc) {
	for _, npc := range npcs {
		s.Env.Set(uppercaseAll(npc.Name), &object.Integer{Value: int64(npc.NpcId)})
	}
}

func (s *ScriptManager) LoadScripts(path string, game *shared.Game, npcs map[uint16]*grpgnpc.Npc) error {
	env := object.NewEnclosedEnvinronment(s.Env)
	AddListeners(env, s)
	AddGlobals(env, game, npcs)

	entries, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".grpgscript") {
			fullPath := filepath.Join(path, entry.Name())

			bytes, err := os.ReadFile(fullPath)
			if err != nil {
				return err
			}

			l := lexer.New(string(bytes))
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				log.Printf("Found errors parsing script: %s\n", fullPath)
				for _, msg := range p.Errors() {
					log.Println(msg)
				}
				return errors.New("errors parsing scripts")
			}

			perf.ConstFold(program)

			obj := evaluator.Eval(program, env)
			if obj != nil && obj.Type() == object.ERROR_OBJ {
				log.Printf("script %s errored %s", fullPath, obj.Inspect())
			}
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
