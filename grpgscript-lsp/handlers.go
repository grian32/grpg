package grpgscript_lsp

import (
	"context"
	"fmt"

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
	_ = h.Client.LogMessage(ctx, &protocol.LogMessageParams{
		Type:    protocol.MessageTypeInfo,
		Message: "GRPGScript LSP Initialized",
	})

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
	diagnostics := h.validateDocuments(params.TextDocument.URI, params.TextDocument.Text)

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

	updatedText := h.applyChanges(doc.Text, params.ContentChanges, ctx)

	doc.Text = updatedText
	doc.Version = params.TextDocument.Version

	diagnostics := h.validateDocuments(params.TextDocument.URI, updatedText)

	return h.Client.PublishDiagnostics(ctx, &protocol.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diagnostics,
	})
}

func (h Handler) validateDocuments(uri protocol.DocumentURI, text string) []protocol.Diagnostic {
	return []protocol.Diagnostic{}
}

func (h Handler) applyChanges(currText string, changes []protocol.TextDocumentContentChangeEvent, ctx context.Context) string {
	text := currText

	for _, change := range changes {
		// if not included = 0
		if change.RangeLength == 0 && isZeroRange(change.Range) {
			text = change.Text
		} else {
			h.applyRangeChanges(text, change.Range, change.Text, ctx)
		}
	}

	return text
}

func (h Handler) applyRangeChanges(text string, rang protocol.Range, changed string, ctx context.Context) string {
	_ = h.Client.LogMessage(ctx, &protocol.LogMessageParams{
		Type:    protocol.MessageTypeInfo,
		Message: fmt.Sprintf("file: %s; range: %v, changed: %s", text, rang, changed),
	})
	return ""
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
