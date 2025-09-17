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

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

var log *zap.Logger

type Handler struct {
	protocol.Server
	Client    protocol.Client
	Documents *DocumentStore
	Env       *object.Environment
	Objs      map[string]uint16
	Items     map[string]uint16
	Npcs      map[string]uint16
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
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "grpgscriptlsp",
			Version: "0.1.0",
		},
	}, nil
}

func (h Handler) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) (err error) {
	h.Documents.Set(params.TextDocument.URI, openParamsToDocuments(params))
	diagnostics := h.validateDocuments(params.TextDocument.Text, ctx)

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

	diagnostics := h.validateDocuments(updatedText, ctx)

	return h.Client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diagnostics,
	})
}

func (h Handler) validateDocuments(text string, ctx context.Context) []protocol.Diagnostic {
	l := lexer.New(text)
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()

	_ = h.Client.LogMessage(ctx, &protocol.LogMessageParams{
		Type:    protocol.MessageTypeInfo,
		Message: fmt.Sprintf("%s", text),
	})

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
		env := object.NewEnclosedEnvinronment(h.Env)

		eval.Eval(program, env)

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

	_ = h.Client.LogMessage(ctx, &protocol.LogMessageParams{
		Type:    protocol.MessageTypeInfo,
		Message: fmt.Sprintf("%v", diags),
	})

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

func openParamsToDocuments(params *protocol.DidOpenTextDocumentParams) *Document {
	return &Document{
		URI:     params.TextDocument.URI,
		Text:    params.TextDocument.Text,
		Version: params.TextDocument.Version,
	}
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
