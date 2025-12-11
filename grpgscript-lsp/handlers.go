package grpgscript_lsp

import (
	"context"
	"fmt"
	"grpg/data-go/grpgitem"
	"grpg/data-go/grpgnpc"
	"grpg/data-go/grpgobj"
	"grpgscript/evaluator"
	"grpgscript/lexer"
	"grpgscript/object"
	"grpgscript/parser"
	"strings"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

var log *zap.Logger

var langBuiltinCompletions = []string{
	"println",
	"len",
	"concat",
	"unshift",
	"push",
}

var langBuiltinLabels = map[string]string{
	"println": "println(STRING)",
	"len":     "len(ARRAY|STRING) INT",
	"concat":  "concat(ARRAY, ARRAY) ARRAY",
	"unshift": "unshift(ARRAY, ANY) INT",
	"push":    "push(ARRAY, ANY) INT",
}

type Handler struct {
	protocol.Server
	Client      protocol.Client
	Documents   *DocumentStore
	Env         *object.Environment
	Objs        map[string]uint16
	Items       map[string]uint16
	Npcs        map[string]uint16
	Definitions map[string]BuiltinDefinition
}

func NewHandler(ctx context.Context, server protocol.Server, client protocol.Client, objs []grpgobj.Obj, npcs []grpgnpc.Npc, items []grpgitem.Item, logger *zap.Logger) (Handler, context.Context, error) {
	log = logger
	h := Handler{
		Server:    server,
		Client:    client,
		Documents: NewDocumentStore(),
		Env:       object.NewEnvironment(),
		Objs:      make(map[string]uint16),
		Items:     make(map[string]uint16),
		Npcs:      make(map[string]uint16),
	}

	for _, obj := range objs {
		h.Objs[obj.Name] = obj.ObjId
		h.Env.Set(UppercaseAll(obj.Name), &object.Integer{Value: int64(obj.ObjId)})
	}

	for _, npc := range npcs {
		h.Npcs[npc.Name] = npc.NpcId
		h.Env.Set(UppercaseAll(npc.Name), &object.Integer{Value: int64(npc.NpcId)})
	}

	for _, item := range items {
		h.Objs[item.Name] = item.ItemId
		h.Env.Set(UppercaseAll(item.Name), &object.Integer{Value: int64(item.ItemId)})
	}

	h.Env.Set("FORAGING", &object.Integer{Value: 0})

	h.Definitions = BuildDefinitions()
	MockBuiltins(h.Env, h.Definitions)

	return h, ctx, nil
}

func (h Handler) Initialize(ctx context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	log.Info("GRPGScript LSP Initialized")
	err := h.Client.LogMessage(ctx, &protocol.LogMessageParams{
		Type:    protocol.MessageTypeInfo,
		Message: "GRPGScript LSP Initialized",
	})
	if err != nil {
		return nil, err
	}

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindIncremental,
				Save:      &protocol.SaveOptions{IncludeText: true},
			},
			CompletionProvider: &protocol.CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{},
			},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "grpgscriptlsp",
			Version: "0.1.0",
		},
	}, nil
}

func (h Handler) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) (err error) {
	h.Documents.Set(params.TextDocument.URI, openParamsToDocuments(params, h.Env))
	doc, _ := h.Documents.Get(params.TextDocument.URI)
	diagnostics := h.validateDocuments(params.TextDocument.Text, doc)

	return h.Client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diagnostics,
	})
}

func (h Handler) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) (err error) {
	doc, ok := h.Documents.Get(params.TextDocument.URI)
	if !ok {
		return fmt.Errorf("document not found: %s", params.TextDocument.URI)
	}

	updatedText := h.applyChanges(doc.Text, params.ContentChanges)

	doc.Text = updatedText
	doc.Version = params.TextDocument.Version

	diagnostics := h.validateDocuments(updatedText, doc)

	return h.Client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diagnostics,
	})
}

func (h Handler) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	document, ok := h.Documents.Get(params.TextDocument.URI)
	if !ok {
		return nil, fmt.Errorf("document not found: %s", params.TextDocument.URI)
	}
	line := strings.Split(document.Text, "\n")[params.Position.Line]

	prefix := getPrefixForLine(line, params.Position.Character)

	list := &protocol.CompletionList{
		IsIncomplete: true,
		Items:        make([]protocol.CompletionItem, 0),
	}

	for _, s := range langBuiltinCompletions {
		if strings.HasPrefix(s, prefix) {
			list.Items = append(list.Items, protocol.CompletionItem{
				Label:            langBuiltinLabels[s],
				Kind:             protocol.CompletionItemKindFunction,
				InsertTextFormat: protocol.InsertTextFormatSnippet,
				InsertText:       s + "($0)",
			})
		}
	}

	doc, _ := h.Documents.Get(params.TextDocument.URI)

	for s := range doc.Env.Names {
		if strings.HasPrefix(s, prefix) {
			obj, _ := doc.Env.Get(s)

			label := s
			kind := protocol.CompletionItemKindVariable
			insertText := s
			detail := ""

			if obj.Type() == object.FUNCTION_OBJ || obj.Type() == object.BUILTIN_OBJ {
				kind = protocol.CompletionItemKindFunction
				insertText += "($0)"

				if def, ok := h.Definitions[s]; ok {
					label = def.Label
					if len(def.Types) > 0 && def.Types[len(def.Types)-1] == FUNCTION {
						insertText += " {}"
					}
				}
			} else {
				detail = obj.Inspect()
			}

			list.Items = append(list.Items, protocol.CompletionItem{
				Label:            label,
				Kind:             kind,
				InsertTextFormat: protocol.InsertTextFormatSnippet,
				InsertText:       insertText,
				Detail:           detail,
			})
		}
	}

	return list, nil
}

func (h Handler) validateDocuments(text string, doc *Document) []protocol.Diagnostic {
	l := lexer.New(text)
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()

	diags := make([]protocol.Diagnostic, len(errors))

	for i, err := range errors {
		diags[i] = protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(err.Line),
					Character: uint32(err.Col),
				},
				End: protocol.Position{
					Line:      uint32(err.Line),
					Character: uint32(err.End),
				},
			},
			Severity: protocol.DiagnosticSeverityError,
			Source:   "grpgscriptlsp",
			Message:  err.Msg,
		}
	}

	// unfortunately we can only really eval if the script passes parsing.
	if len(diags) == 0 {
		eval := evaluator.NewEvaluator()
		newEnv := cloneEnv(h.Env)

		eval.Eval(program, newEnv)

		doc.Env = newEnv

		for _, err := range eval.ErrorStore.Errors {
			diags = append(diags, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(err.Position.StartLine),
						Character: uint32(err.Position.Start),
					},
					End: protocol.Position{
						Line:      uint32(err.Position.EndLine),
						Character: uint32(err.Position.End),
					},
				},
				Severity: protocol.DiagnosticSeverityError,
				Source:   "grpgscriptlsp",
				Message:  err.Msg,
			})
		}
	}

	return diags
}

func (h Handler) applyChanges(currText string, changes []protocol.TextDocumentContentChangeEvent) string {
	text := currText

	for _, change := range changes {
		// per documentation: If range and rangeLength are omitted the new text is considered to be the full content of the document.
		text = h.applyRangeChanges(text, change.Range, change.Text)
	}

	return text
}

func (h Handler) applyRangeChanges(text string, rang protocol.Range, changed string) string {
	start := posToOffset(text, rang.Start)
	end := posToOffset(text, rang.End)

	return text[:start] + changed + text[end:]
}

func getPrefixForLine(line string, col uint32) string {
	var startCol uint32 = 0

	if col > 0 {
		for i := col - 1; i > 0; i-- {
			if !IsAlpha(line[i]) {
				startCol = i + 1
				break
			}
		}
	}

	return line[startCol:col]
}

func openParamsToDocuments(params *protocol.DidOpenTextDocumentParams, hEnv *object.Environment) *Document {
	env := cloneEnv(hEnv)

	return &Document{
		URI:     params.TextDocument.URI,
		Text:    params.TextDocument.Text,
		Version: params.TextDocument.Version,
		Env:     env,
	}
}

func cloneEnv(env *object.Environment) *object.Environment {
	newEnv := object.NewEnvironment()

	for s := range env.Names {
		obj, _ := env.Get(s)
		newEnv.Set(s, obj)
	}

	return newEnv
}

func posToOffset(text string, pos protocol.Position) int {
	var line uint32 = 0
	var col uint32 = 0

	for i, r := range text {
		if line == pos.Line && col == pos.Character {
			return i
		}

		if r == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}

	// means pos is eof
	return len(text)
}
