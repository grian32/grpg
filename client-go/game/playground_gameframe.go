package game

import (
	"client/shared"
	"client/util"
	"fmt"
	"image/color"

	gebiten_ui "github.com/grian32/gebiten-ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type PgGameframe struct {
	WorldImage *ebiten.Image
	Font16     *gebiten_ui.GFont
	Font20     *gebiten_ui.GFont
	Font24     *gebiten_ui.GFont

	Textures           map[uint16]*ebiten.Image
	ItemOutlineTexture *ebiten.Image
	SkillIcons         map[shared.Skill]*gebiten_ui.GHoverTexture
	GameframeRight     *ebiten.Image
	GameframeBottom    *ebiten.Image
	SkillsButton       *gebiten_ui.GTextureButton
	InventoryButton    *gebiten_ui.GTextureButton

	CurrActionString string

	Player       *shared.LocalPlayer
	Game         *shared.Game
	InputHandler *PgInputHandler
}

func NewPgGameframe(
	worldImage *ebiten.Image,
	game *shared.Game,
	inputHandler *PgInputHandler,
	font16 *gebiten_ui.GFont,
	font20 *gebiten_ui.GFont,
	font24 *gebiten_ui.GFont,
	textures map[uint16]*ebiten.Image,
	otherTex map[string]*ebiten.Image,
) *PgGameframe {
	g := &PgGameframe{
		WorldImage:         worldImage,
		Game:               game,
		Player:             game.Player,
		InputHandler:       inputHandler,
		Font16:             font16,
		Font20:             font20,
		Font24:             font24,
		Textures:           textures,
		ItemOutlineTexture: otherTex["item_outline"],
		SkillIcons:         make(map[shared.Skill]*gebiten_ui.GHoverTexture),
		CurrActionString:   "Current Action: None :(",
	}

	hoverTex := otherTex["hover_tex"]
	foragingIconTex := otherTex["foraging_icon"]

	g.SkillIcons[shared.Foraging] = gebiten_ui.NewHoverTexture(RightGameframeX+TileSize, TileSize, RightGameframeX+(TileSize*5), foragingIconTex, g.Game.SkillHoverMsgs[shared.Foraging], hoverTex, font16, color.White)
	g.GameframeRight = otherTex["gameframe_right"]
	g.GameframeBottom = otherTex["gameframe_bottom"]

	g.InventoryButton = gebiten_ui.NewTextureButton(RightGameframeX+InvButtonXOffset, 0, otherTex["inv_button"], func() {
		g.Game.GameframeContainerRenderType = shared.Inventory
	})

	g.SkillsButton = gebiten_ui.NewTextureButton(RightGameframeX+SkillsButtonXOffset, 0, otherTex["skills_button"], func() {
		g.Game.GameframeContainerRenderType = shared.Skills
	})
	return g
}

func (g *PgGameframe) Update() {
	facingCoord := g.Game.Player.GetFacingCoord()
	trackedObj, objExists := g.Game.TrackedObjs[facingCoord]
	trackedNpc, npcExists := g.Game.NpcsByPos[facingCoord]
	if objExists {
		g.CurrActionString = "Current Action: " + trackedObj.DataObj.InteractText
	} else if npcExists {
		g.CurrActionString = "Talk to " + trackedNpc.NpcData.Name
	} else {
		g.CurrActionString = "Current Action: None :("
	}

	for _, si := range g.SkillIcons {
		si.Update()
	}

	g.InventoryButton.Update()
	g.SkillsButton.Update()
}

func (g *PgGameframe) Draw() {
	util.DrawImage(g.WorldImage, g.GameframeRight, RightGameframeX, 0)
	if g.InputHandler.IsTypingCommand {
		g.Font16.Draw(g.WorldImage, "Command: "+g.InputHandler.CommandString, 0, CommandY, color.White)
	}

	// TODO: i think i can move render type out of game?
	if g.Game.GameframeContainerRenderType == shared.Inventory {
		var currItemRealPosX int32 = RightGameframeX + TileSize
		var currItemRealPosY int32 = TileSize

		for idx, item := range g.Player.Inventory {
			if item.ItemId == 0 {
				continue
			}

			data := g.Game.Items[item.ItemId]
			tex := g.Textures[data.Texture]
			util.DrawImage(g.WorldImage, tex, currItemRealPosX, currItemRealPosY)

			g.Font16.Draw(g.WorldImage, fmt.Sprintf("%d", item.Count), float64(currItemRealPosX+ItemCountXOffset), float64(currItemRealPosY+ItemCountYOffset), color.White)

			if idx == g.Game.OutlineInvSpot {
				util.DrawImage(g.WorldImage, g.ItemOutlineTexture, currItemRealPosX, currItemRealPosY)
			}

			currItemRealPosX += TileSize

			if (idx+1)%ItemsPerRow == 0 {
				currItemRealPosY += TileSize
				currItemRealPosX = RightGameframeX + TileSize
			}
		}
	} else if g.Game.GameframeContainerRenderType == shared.Skills {
		for _, si := range g.SkillIcons {
			si.Draw(g.WorldImage)
		}
		for i := shared.Foraging; i <= shared.Foraging; i++ {
			// TODO: maybe string can be pre computed by packet here?
			// TODO: magic constants when i do this det
			g.Font16.Draw(g.WorldImage, fmt.Sprintf("%d", g.Game.Skills[i].Level), RightGameframeX+64+32, 64+48, util.Yellow)
		}
	}

	util.DrawImage(g.WorldImage, g.GameframeBottom, 0, RightGameframeX)

	talkbox := g.Game.Talkbox
	// x is offset from 0, y has offset added, to be placed in the right spot
	g.Font20.Draw(g.WorldImage, g.CurrActionString, CurrActionX, RightGameframeX+CurrNameActionYOffset, color.White)
	if talkbox.Active {
		g.Font24.Draw(g.WorldImage, talkbox.CurrentName, CurrNameX, RightGameframeX+CurrNameActionYOffset, color.White)
		var currY float64 = CurrMessageY
		for _, s := range talkbox.CurrentMessage {
			g.Font24.Draw(g.WorldImage, s, CurrMessageX, currY, color.White)
			currY += 30
		}
	}
	g.InventoryButton.Draw(g.WorldImage)
	g.SkillsButton.Draw(g.WorldImage)

	playerCoords := fmt.Sprintf("X: %d, Y: %d, Facing: %s", g.Player.X, g.Player.Y, g.Player.Facing.String())
	g.Font24.Draw(g.WorldImage, playerCoords, RightGameframeX, DebugCoordsY, color.White)
}
