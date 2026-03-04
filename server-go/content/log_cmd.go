package content

import (
	"log"
	"server/scripts"
	"strings"
)

func init() {
	scripts.OnCommand("log", func(ctx *scripts.CommandCtx) {
		log.Printf("cmd log: %s", strings.Join(ctx.Args(), " "))
	})
}
