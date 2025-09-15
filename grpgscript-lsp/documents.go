package grpgscript_lsp

import (
	"sync"

	"go.lsp.dev/protocol"
)

type DocumentStore struct {
	mu        sync.RWMutex
	Documents map[protocol.DocumentURI]*Document
}

type Document struct {
	URI     protocol.DocumentURI
	Text    string
	Version int32
}

func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		Documents: make(map[protocol.DocumentURI]*Document),
	}
}

func (ds *DocumentStore) Set(uri protocol.DocumentURI, doc *Document) {
	ds.mu.Lock()
	ds.Documents[uri] = doc
	ds.mu.Unlock()
}

func (ds *DocumentStore) Get(uri protocol.DocumentURI) (*Document, bool) {
	ds.mu.RLock()
	documents, ok := ds.Documents[uri]
	ds.mu.RUnlock()
	return documents, ok
}
