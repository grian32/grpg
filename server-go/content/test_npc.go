package content

import "server/scripts"

func init() {
	scripts.SpawnNpc(scripts.TEST, 1, 1)

	scripts.OnTalkNpc(scripts.TEST, func(ctx *scripts.NpcTalkCtx) {
		ctx.ClearDialogueQueue()

		ctx.TalkPlayer("hello, test")
		ctx.TalkNpc("...")
		ctx.TalkPlayer("C U")

		ctx.StartDialogue()
	})
}
