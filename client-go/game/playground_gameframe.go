package game

import (
	"client/shared"
	"client/util"
	"fmt"
	"image"
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

type DrawItem struct {
	currX, currY int32
	count        uint16
	invIdx       int
	dynHoverTex  *gebiten_ui.GDynHoverTexture
}

type PgGameframe struct {
	Font16 *gebiten_ui.GFont
	Font20 *gebiten_ui.GFont
	Font24 *gebiten_ui.GFont

	Textures           map[uint16]*ebiten.Image
	ItemOutlineTexture *ebiten.Image
	SkillIcons         map[shared.Skill]*gebiten_ui.GHoverTexture
	GameframeRight     *ebiten.Image
	GameframeBottom    *ebiten.Image
	HealthLeftTexture  *ebiten.Image
	HealthMidTexture   *ebiten.Image
	HealthRightTexture *ebiten.Image

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

	ItemHoverTextures map[uint16]*gebiten_ui.GDynHoverTexture
	HoverTexX         int
	HoverTexY         int
	DrawInvItems      []DrawItem
	DrawEquipItems    [5]*gebiten_ui.GDynHoverTexture

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

	healthTex := otherTex["healthbar_segments"]
	g.HealthLeftTexture = healthTex.SubImage(image.Rect(0, 0, 8, 24)).(*ebiten.Image)
	g.HealthMidTexture = healthTex.SubImage(image.Rect(8, 0, 16, 24)).(*ebiten.Image)
	g.HealthRightTexture = healthTex.SubImage(image.Rect(16, 0, 24, 24)).(*ebiten.Image)

	g.SkillIcons[shared.Foraging] = gebiten_ui.NewHoverTexture(RightGameframeX+TileSize, TileSize, RightGameframeX+(TileSize*5), foragingIconTex, g.Game.SkillHoverMsgs[shared.Foraging], hoverTex, font16, color.White)
	g.GameframeRight = otherTex["gameframe_right"]
	g.GameframeBottom = otherTex["gameframe_bottom"]

	g.InventoryButton = gebiten_ui.NewTextureButton(RightGameframeX+InvButtonXOffset, 0, otherTex["inv_button"], func() {
		g.ContainerRenderType = Inventory
		g.OutlineInvSpot = -1
	})

	g.SkillsButton = gebiten_ui.NewTextureButton(RightGameframeX+SkillsButtonXOffset, 0, otherTex["skills_button"], func() {
		g.ContainerRenderType = Skills
	})

	g.EquipmentButton = gebiten_ui.NewTextureButton(RightGameframeX+EquipmentButtonXOffset, 0, otherTex["equipment_button"], func() {
		g.ContainerRenderType = Equipment
		g.OutlineInvSpot = -1
	})

	g.ContainerRenderType = Inventory

	g.EquipmentFrame = otherTex["equipment_outline"]
	g.HelmetFrame = otherTex["helmet_outline"]
	g.ChestplateFrame = otherTex["chestplate_outline"]
	g.LeggingsFrame = otherTex["legs_outline"]
	g.RingFrame = otherTex["ring_outline"]
	g.WeaponFrame = otherTex["wep_outline"]

	g.ItemHoverTextures = make(map[uint16]*gebiten_ui.GDynHoverTexture)
	for id, item := range game.Items {
		g.ItemHoverTextures[id] = gebiten_ui.NewDynHoverTexture(RightGameframeX+(TileSize*5), textures[item.Texture], new(item.Name), hoverTex, font16, color.White)
	}

	return g
}

func (g *PgGameframe) Update() {
	g.DrawInvItems = make([]DrawItem, 0)
	g.DrawEquipItems = [5]*gebiten_ui.GDynHoverTexture{}
	g.UpdateCurrActionString()

	for _, si := range g.SkillIcons {
		si.Update()
	}

	if g.ContainerRenderType == Inventory {
		g.UpdateInventoryPanel()
	} else if g.ContainerRenderType == Equipment {
		g.UpdateEquipmentPanel()
	}

	g.InventoryButton.Update()
	g.SkillsButton.Update()
	g.EquipmentButton.Update()

	g.InputHandler.UpdateItemMove(g.ContainerRenderType, &g.OutlineInvSpot)
}

func (g *PgGameframe) Draw(screen *ebiten.Image) {
	util.DrawImage(screen, g.GameframeRight, RightGameframeX, 0)

	if g.Player.Health == 1 {
		util.DrawImage(screen, g.HealthLeftTexture, HealthBarX, HealthBarY)
		util.DrawImage(screen, g.HealthRightTexture, HealthBarX+8, HealthBarY)
	} else if g.Player.Health >= 2 {
		util.DrawImage(screen, g.HealthLeftTexture, HealthBarX, HealthBarY)
		var currX int32 = 8
		for _ = range int(g.Player.Health) * 35 / 100 {
			util.DrawImage(screen, g.HealthMidTexture, HealthBarX+currX, HealthBarY)
			currX += 8
		}
		util.DrawImage(screen, g.HealthRightTexture, HealthBarX+currX, HealthBarY)
	}

	if g.InputHandler.IsTypingCommand {
		g.Font16.Draw(screen, "Command: "+g.InputHandler.CommandString, 0, CommandY, color.White)
	}

	if g.ContainerRenderType == Inventory {
		g.DrawInventoryPanel(screen)
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
		g.DrawEquipmentPanel(screen)
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

func (g *PgGameframe) DrawItem(tex *gebiten_ui.GDynHoverTexture, screen *ebiten.Image, x, y, hoverOffsetX, hoverOffsetY int32, outlineSpot int) {
	tex.Draw(screen, float64(x), float64(y), float64(hoverOffsetX), float64(hoverOffsetY))

	if g.OutlineInvSpot == outlineSpot {
		util.DrawImage(screen, g.ItemOutlineTexture, x, y)
	}
}

func (g *PgGameframe) UpdateCurrActionString() {
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
}

func (g *PgGameframe) UpdateInventoryPanel() {
	var currItemRealPosX int32 = RightGameframeX + TileSize
	var currItemRealPosY int32 = TileSize

	for idx, item := range g.Player.Inventory {
		if item.ItemId != 0 {
			g.DrawInvItems = append(g.DrawInvItems, DrawItem{
				currX:       currItemRealPosX,
				currY:       currItemRealPosY,
				count:       item.Count,
				invIdx:      idx,
				dynHoverTex: g.ItemHoverTextures[item.ItemId],
			})
		}

		currItemRealPosX += TileSize
		if (idx+1)%ItemsPerRow == 0 {
			currItemRealPosY += TileSize
			currItemRealPosX = RightGameframeX + TileSize
		}
	}

	for _, item := range g.DrawInvItems {
		item.dynHoverTex.Update(float64(item.currX), float64(item.currY))
	}
}

func (g *PgGameframe) DrawInventoryPanel(screen *ebiten.Image) {
	// reverse order so that hover textures show up in the correct order, if i was drawing it in order then hover tex would show up beneath lower items etc
	for i := len(g.DrawInvItems) - 1; i >= 0; i-- {
		item := g.DrawInvItems[i]
		g.DrawItem(item.dynHoverTex, screen, item.currX, item.currY, 0, 0, item.invIdx)
		if item.count > 1 {
			g.Font16.Draw(screen, fmt.Sprintf("%d", item.count), float64(item.currX+ItemCountXOffset), float64(item.currY+ItemCountYOffset), color.White)
		}
	}
}

func (g *PgGameframe) UpdateEquipmentPanel() {
	e := g.Player.Equipment

	// maybe inefficient access?
	if e[shared.HELMET] != 0 {
		g.DrawEquipItems[shared.HELMET] = g.ItemHoverTextures[e[shared.HELMET]]
	}
	if e[shared.CHESTPLATE] != 0 {
		g.DrawEquipItems[shared.CHESTPLATE] = g.ItemHoverTextures[e[shared.CHESTPLATE]]
	}
	if e[shared.LEGGINGS] != 0 {
		g.DrawEquipItems[shared.LEGGINGS] = g.ItemHoverTextures[e[shared.LEGGINGS]]
	}
	if e[shared.WEAPON] != 0 {
		g.DrawEquipItems[shared.WEAPON] = g.ItemHoverTextures[e[shared.WEAPON]]
	}
	if e[shared.RING] != 0 {
		g.DrawEquipItems[shared.RING] = g.ItemHoverTextures[e[shared.RING]]
	}

	for idx, item := range g.DrawEquipItems {
		if item != nil {
			x, y := getEquipmentPositions(shared.EquipmentType(idx))
			item.Update(x, y)
		}
	}
}

func (g *PgGameframe) DrawEquipmentPanel(screen *ebiten.Image) {
	// reverse order so that hover textures show up in the correct order, if i was drawing it in order then hover tex would show up beneath lower items etc
	for i := 4; i >= 0; i-- {
		item := g.DrawEquipItems[i]
		if item != nil {
			x, y := getEquipmentPositions(shared.EquipmentType(i))
			g.DrawItem(item, screen, int32(x), int32(y), -4, 0, 24+i)
		}
	}
}

func getEquipmentPositions(e shared.EquipmentType) (float64, float64) {
	switch e {
	case shared.HELMET:
		return RightGameframeX + EquipmentMidOffsetX, HelmetOffsetY
	case shared.CHESTPLATE:
		return RightGameframeX + EquipmentMidOffsetX, EquipmentMidOffsetY
	case shared.LEGGINGS:
		return RightGameframeX + EquipmentMidOffsetX, LeggingsOffsetY
	case shared.WEAPON:
		return RightGameframeX + WeaponOffsetX, EquipmentMidOffsetY
	case shared.RING:
		return RightGameframeX + RingOffsetX, EquipmentMidOffsetY
	}

	return 0, 0
}
