package content

import "server/scripts"

func init() {
	scripts.SpawnNpc(scripts.GRPG_GUIDE, 3, 1, 0)

	scripts.OnTalkNpc(scripts.GRPG_GUIDE, func(ctx *scripts.NpcTalkCtx) {
		ctx.ClearDialogueQueue()

		ctx.TalkNpc("Hello, and welcome to GRPG!\nYou can advance my dialogue by pressing Space.")
		ctx.TalkNpc("You can move around with WASD,\nand interact with Objects & NPCs by pressing Q.")

		ctx.StartDialogue()
	})
}
