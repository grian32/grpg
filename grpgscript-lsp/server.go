package grpgscript_lsp

import (
	"context"
	"errors"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgitem"
	"grpg/data-go/grpgnpc"
	"grpg/data-go/grpgobj"
	"io"
	"net/http"
	"os"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var InvalidHeaderError = errors.New("invalid header")

// StartServer starts the language server.
// It reads from stdin and writes to stdout.
func StartServer(logger *zap.Logger) {
	conn := jsonrpc2.NewConn(jsonrpc2.NewStream(&readWriteCloser{
		reader: os.Stdin,
		writer: os.Stdout,
	}))

	client := protocol.ClientDispatcher(conn, logger)
	server := protocol.ServerDispatcher(conn, logger)

	// TODO: update this later when its an actual URL and keep cache of it
	assets := "http://51.83.129.212:4022/assets/"

	objBuf := getAndRead(assets+"objs.grpgobj", "objs", logger)
	npcBuf := getAndRead(assets+"npcs.grpgnpc", "npcs", logger)
	itemsBuf := getAndRead(assets+"items.grpgitem", "items", logger)

	objHeader, err := grpgobj.ReadHeader(objBuf)
	panicLogErr("reading grpgobj header", err, logger)
	if objHeader.Magic != [8]byte([]byte("GRPGOBJ\x00")[:]) {
		panicLogErr("grpgobj file doesn't having GRPGOBJ header", InvalidHeaderError, logger)
	}

	npcHeader, err := grpgnpc.ReadHeader(npcBuf)
	panicLogErr("reading grpgnpc header", err, logger)
	if npcHeader.Magic != [8]byte([]byte("GRPGNPC\x00")[:]) {
		panicLogErr("grpgnpc file doesn't having GRPGNPC header", InvalidHeaderError, logger)
	}

	itemHeader, err := grpgitem.ReadHeader(itemsBuf)
	panicLogErr("reading grpgitem header", err, logger)
	if itemHeader.Magic != [8]byte([]byte("GRPGITEM")[:]) {
		panicLogErr("grpgitem file doesn't having GRPGITEM header", InvalidHeaderError, logger)
	}

	objs, err := grpgobj.ReadObjs(objBuf)
	panicLogErr("reading grpgobj objs", err, logger)

	npcs, err := grpgnpc.ReadNpcs(npcBuf)
	panicLogErr("reading grpgnpc npcs", err, logger)

	items, err := grpgitem.ReadItems(itemsBuf)
	panicLogErr("reading grpgitem items", err, logger)

	handler, ctx, err := NewHandler(
		context.Background(),
		server,
		client,
		objs,
		npcs,
		items,
		logger,
	)

	if err != nil {
		logger.Sugar().Fatalf("while initializing handler: %v", err)
	}

	conn.Go(ctx, protocol.ServerHandler(
		handler, jsonrpc2.MethodNotFoundHandler,
	))
	<-conn.Done()
}

func getAndRead(url, name string, logger *zap.Logger) *gbuf.GBuf {
	req, err := http.Get(url)
	panicLogErr("fetching"+name, err, logger)
	data, err := io.ReadAll(req.Body)
	panicLogErr("reading"+name, err, logger)

	return gbuf.NewGBuf(data)
}

func panicLogErr(name string, err error, logger *zap.Logger) {
	if err != nil {
		logger.Sugar().Fatalf("while %s: %v", name, err)
	}
}

type readWriteCloser struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

func (r *readWriteCloser) Read(b []byte) (int, error) {
	n, err := r.reader.Read(b)
	return n, err
}

func (r *readWriteCloser) Write(b []byte) (int, error) {
	return r.writer.Write(b)
}

func (r *readWriteCloser) Close() error {
	return multierr.Append(r.reader.Close(), r.writer.Close())
}
