package grpgscript_lsp

import (
	"context"
	"fmt"
	"grpgscript/lexer"
	"grpgscript/parser"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

var log *zap.Logger

type Handler struct {
	protocol.Server
	Client    protocol.Client
	Documents *DocumentStore
}

func NewHandler(ctx context.Context, server protocol.Server, client protocol.Client, logger *zap.Logger) (Handler, context.Context, error) {
	log = logger
	return Handler{
		Server:    server,
		Client:    client,
		Documents: NewDocumentStore(),
	}, ctx, nil
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
	diagnostics := h.validateDocuments(params.TextDocument.URI, params.TextDocument.Text, ctx)

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

	diagnostics := h.validateDocuments(params.TextDocument.URI, updatedText, ctx)

	return h.Client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diagnostics,
	})
}

func (h Handler) validateDocuments(uri protocol.DocumentURI, text string, ctx context.Context) []protocol.Diagnostic {
	l := lexer.New(text)
	p := parser.New(l)
	_ = p.ParseProgram()

	errors := p.Errors()
	diags := make([]protocol.Diagnostic, len(errors))

	_ = h.Client.LogMessage(ctx, &protocol.LogMessageParams{
		Type:    protocol.MessageTypeInfo,
		Message: fmt.Sprintf("%v", errors),
	})

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
		if change.RangeLength == 0 && isZeroRange(change.Range) {
			text = change.Text
			return text
		} else {
			text = h.applyRangeChanges(text, change.Range, change.Text)
		}
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

func isZeroRange(rang protocol.Range) bool {
	return rang.Start.Line == 0 && rang.Start.Character == 0 && rang.End.Line == 0 && rang.End.Character == 0
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
