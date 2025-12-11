module grpgscript-lsp

go 1.24.4

require (
	go.lsp.dev/jsonrpc2 v0.10.0
	go.lsp.dev/protocol v0.12.0
	go.uber.org/multierr v1.8.0
	go.uber.org/zap v1.21.0
	grpg/data-go v0.0.0
	grpgscript v0.0.0
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/segmentio/asm v1.1.3 // indirect
	github.com/segmentio/encoding v0.3.4 // indirect
	go.lsp.dev/pkg v0.0.0-20210717090340-384b27a52fb2 // indirect
	go.lsp.dev/uri v0.3.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace grpgscript => ../grpgscript

replace grpg/data-go => ../data-go
