package content

import "server/scripts"

func init() {
	scripts.SpawnNpc(scripts.TEST, 3, 3, 2) // so 2 tiles in each direction from x 1 y 1
	scripts.SpawnNpc(scripts.TEST, 1, 1, 2) // so 2 tiles in each direction from x 1 y 1

	scripts.OnTalkNpc(scripts.TEST, func(ctx *scripts.NpcTalkCtx) {
		ctx.ClearDialogueQueue()

		ctx.TalkPlayer("hello, test")
		ctx.TalkNpc("...")
		ctx.TalkPlayer("C U")

		ctx.StartDialogue()
	})
}
