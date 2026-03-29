package game

import (
	"client/shared"
	"client/util"
	"fmt"
	"image/color"

	gebiten_ui "github.com/grian32/gebiten-ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type RenderType byte

const (
	Inventory RenderType = iota
	Skills
	Equipment
)

type PgGameframe struct {
	Font16 *gebiten_ui.GFont
	Font20 *gebiten_ui.GFont
	Font24 *gebiten_ui.GFont

	Textures           map[uint16]*ebiten.Image
	ItemOutlineTexture *ebiten.Image
	SkillIcons         map[shared.Skill]*gebiten_ui.GHoverTexture
	GameframeRight     *ebiten.Image
	GameframeBottom    *ebiten.Image

	// TODO: maybe move helmet etc into an array based on an  enum somewhere when i get further into the implementation of this..
	EquipmentFrame  *ebiten.Image
	HelmetFrame     *ebiten.Image
	ChestplateFrame *ebiten.Image
	LeggingsFrame   *ebiten.Image
	RingFrame       *ebiten.Image
	WeaponFrame     *ebiten.Image

	SkillsButton    *gebiten_ui.GTextureButton
	InventoryButton *gebiten_ui.GTextureButton
	EquipmentButton *gebiten_ui.GTextureButton

	CurrActionString string

	Player       *shared.LocalPlayer
	Game         *shared.Game
	InputHandler *PgInputHandler

	ContainerRenderType RenderType
	OutlineInvSpot      int
}

func NewPgGameframe(
	game *shared.Game,
	inputHandler *PgInputHandler,
	font16 *gebiten_ui.GFont,
	font20 *gebiten_ui.GFont,
	font24 *gebiten_ui.GFont,
	textures map[uint16]*ebiten.Image,
	otherTex map[string]*ebiten.Image,
) *PgGameframe {
	g := &PgGameframe{
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
		OutlineInvSpot:     -1,
	}

	hoverTex := otherTex["hover_tex"]
	foragingIconTex := otherTex["foraging_icon"]

	g.SkillIcons[shared.Foraging] = gebiten_ui.NewHoverTexture(RightGameframeX+TileSize, TileSize, RightGameframeX+(TileSize*5), foragingIconTex, g.Game.SkillHoverMsgs[shared.Foraging], hoverTex, font16, color.White)
	g.GameframeRight = otherTex["gameframe_right"]
	g.GameframeBottom = otherTex["gameframe_bottom"]

	g.InventoryButton = gebiten_ui.NewTextureButton(RightGameframeX+InvButtonXOffset, 0, otherTex["inv_button"], func() {
		g.ContainerRenderType = Inventory
	})

	g.SkillsButton = gebiten_ui.NewTextureButton(RightGameframeX+SkillsButtonXOffset, 0, otherTex["skills_button"], func() {
		g.ContainerRenderType = Skills
	})

	g.EquipmentButton = gebiten_ui.NewTextureButton(RightGameframeX+EquipmentButtonXOffset, 0, otherTex["equipment_button"], func() {
		g.ContainerRenderType = Equipment
	})

	g.ContainerRenderType = Inventory

	g.EquipmentFrame = otherTex["equipment_outline"]
	g.HelmetFrame = otherTex["helmet_outline"]
	g.ChestplateFrame = otherTex["chestplate_outline"]
	g.LeggingsFrame = otherTex["legs_outline"]
	g.RingFrame = otherTex["ring_outline"]
	g.WeaponFrame = otherTex["wep_outline"]

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
	g.EquipmentButton.Update()

	g.InputHandler.UpdateItemMove(g.ContainerRenderType, &g.OutlineInvSpot)
}

func (g *PgGameframe) Draw(screen *ebiten.Image) {
	util.DrawImage(screen, g.GameframeRight, RightGameframeX, 0)
	if g.InputHandler.IsTypingCommand {
		g.Font16.Draw(screen, "Command: "+g.InputHandler.CommandString, 0, CommandY, color.White)
	}

	// TODO: prob abstracting these into their own functions would be wise
	if g.ContainerRenderType == Inventory {
		var currItemRealPosX int32 = RightGameframeX + TileSize
		var currItemRealPosY int32 = TileSize

		for idx, item := range g.Player.Inventory {
			// you still want to advance the rendering pos since after inv moving is implemented you'll have empty spots and it'll render wrongly
			if item.ItemId != 0 {
				g.DrawItem(item.ItemId, screen, currItemRealPosX, currItemRealPosY)
				g.Font16.Draw(screen, fmt.Sprintf("%d", item.Count), float64(currItemRealPosX+ItemCountXOffset), float64(currItemRealPosY+ItemCountYOffset), color.White)

				if idx == g.OutlineInvSpot {
					util.DrawImage(screen, g.ItemOutlineTexture, currItemRealPosX, currItemRealPosY)
				}
			}

			currItemRealPosX += TileSize
			if (idx+1)%ItemsPerRow == 0 {
				currItemRealPosY += TileSize
				currItemRealPosX = RightGameframeX + TileSize
			}
		}
	} else if g.ContainerRenderType == Skills {
		for _, si := range g.SkillIcons {
			si.Draw(screen)
		}
		for i := shared.Foraging; i <= shared.Foraging; i++ {
			// TODO: maybe string can be pre computed by packet here?
			// TODO: magic constants when i do this det
			g.Font16.Draw(screen, fmt.Sprintf("%d", g.Game.Skills[i].Level), RightGameframeX+64+32, 64+48, util.Yellow)
		}
	} else if g.ContainerRenderType == Equipment {
		util.DrawImage(screen, g.EquipmentFrame, RightGameframeX, 0)
		g.DrawItemOrTex(g.Player.Equipment[shared.HELMET], screen, g.HelmetFrame, RightGameframeX+EquipmentMidOffsetX, HelmetOffsetY)
		g.DrawItemOrTex(g.Player.Equipment[shared.CHESTPLATE], screen, g.ChestplateFrame, RightGameframeX+EquipmentMidOffsetX, EquipmentMidOffsetY)
		g.DrawItemOrTex(g.Player.Equipment[shared.LEGGINGS], screen, g.LeggingsFrame, RightGameframeX+EquipmentMidOffsetX, LeggingsOffsetY)
		g.DrawItemOrTex(g.Player.Equipment[shared.WEAPON], screen, g.WeaponFrame, RightGameframeX+WeaponOffsetX, EquipmentMidOffsetY)
		g.DrawItemOrTex(g.Player.Equipment[shared.RING], screen, g.RingFrame, RightGameframeX+RingOffsetX, EquipmentMidOffsetY)
	}

	util.DrawImage(screen, g.GameframeBottom, 0, RightGameframeX)

	talkbox := g.Game.Talkbox
	// x is offset from 0, y has offset added, to be placed in the right spot
	g.Font20.Draw(screen, g.CurrActionString, CurrActionX, RightGameframeX+CurrNameActionYOffset, color.White)
	if talkbox.Active {
		g.Font24.Draw(screen, talkbox.CurrentName, CurrNameX, RightGameframeX+CurrNameActionYOffset, color.White)
		var currY float64 = CurrMessageY
		for _, s := range talkbox.CurrentMessage {
			g.Font24.Draw(screen, s, CurrMessageX, currY, color.White)
			currY += 30
		}
	}
	g.InventoryButton.Draw(screen)
	g.SkillsButton.Draw(screen)
	g.EquipmentButton.Draw(screen)

	if g.Game.DebugMode {
		playerCoords := fmt.Sprintf("X: %d, Y: %d, Facing: %s", g.Player.X, g.Player.Y, g.Player.Facing.String())
		g.Font24.Draw(screen, playerCoords, RightGameframeX, DebugCoordsY, color.White)
	}
}

func (g *PgGameframe) DrawItem(id uint16, screen *ebiten.Image, x, y int32) {
	data := g.Game.Items[id]
	tex := g.Textures[data.Texture]
	util.DrawImage(screen, tex, x, y)
}

func (g *PgGameframe) DrawItemOrTex(id uint16, screen *ebiten.Image, elseTex *ebiten.Image, x, y int32) {
	if id != 0 {
		g.DrawItem(id, screen, x, y)
	} else {
		util.DrawImage(screen, elseTex, x, y)
	}
}
