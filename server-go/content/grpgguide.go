package content

import (
	"server/constants"
	"server/scripts"
)

func init() {
	scripts.SpawnNpc(constants.GRPG_GUIDE, 3, 1, 0)

	scripts.OnTalkNpc(constants.GRPG_GUIDE, func(ctx *scripts.NpcTalkCtx) {
		ctx.ClearDialogueQueue()

		ctx.TalkNpc("Hello, and welcome to GRPG!\nYou can advance my dialogue by pressing Space.")
		ctx.TalkNpc("You can move around with WASD,\nand interact with Objects & NPCs by pressing Q.")

		ctx.StartDialogue()

		if (ctx.GetPlayerVar(constants.SHOULD_SHOW_TUTORIAL_INDICATOR) == 0) {
			ctx.SetPlayerVar(constants.SHOULD_SHOW_TUTORIAL_INDICATOR, 1)
		}
	})
}
