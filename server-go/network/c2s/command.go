package c2s

import (
	"grpg/data-go/gbuf"
	"log"
	"server/scripts"
	"server/shared"
	"strings"
)

type Command struct {
}

func (c *Command) Handle(buf *gbuf.GBuf, game *shared.Game, player *shared.Player, scriptManager *scripts.ScriptManager) {
	cmd, err := buf.ReadString()
	if err != nil {
		log.Printf("failed to read string in command packet\n")
		return
	}
	split := strings.Split(cmd, " ")
	name := split[0]
	args := split[1:]

	script, ok := scriptManager.CommandScripts[name]
	if !ok {
		log.Printf("couldn't find command with name %s\n", name)
		return
	}

	script(scripts.NewCommandCtx(args))
}
