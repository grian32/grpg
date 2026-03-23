package game

import (
	"client/util"
	"fmt"
	"image/color"
	"log"

	"client/shared"

	gebitenui "github.com/grian32/gebiten-ui"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WorldImageSize  = 1024
	RightGameframeX = 768

	ChunkSize     = 16
	TilesPerChunk = 256
	TileSize      = 64 // pixels per tile, this is used in both ui and actual world stuff

	ExclamAnimTickInterval = 20
	ExclamAnimFrameCount   = 2
	ExclamBobOffset        = -4

	InvButtonXOffset    = TileSize + 16
	SkillsButtonXOffset = TileSize*2 + 32

	CameraOffsetTiles    = 4
	CameraBoundaryTiles  = 12
	CameraMinOffsetTiles = 9
	CameraPanSpeed       = 16.0

	CommandY = 740

	ItemCountXOffset = 6
	ItemCountYOffset = 4
	ItemsPerRow      = 4

	CurrActionX           = 110
	CurrNameActionYOffset = 28 + 3
	CurrNameX             = CurrActionX + 332
	CurrMessageX          = 90
	CurrMessageY          = 840

	DebugCoordsY = 800

	AssetsDir = "../../grpg-assets/"
)

type Playground struct {
	Font16   *gebitenui.GFont
	Font18   *gebitenui.GFont
	Font20   *gebitenui.GFont
	Font24   *gebitenui.GFont
	Textures map[uint16]*ebiten.Image
	Game     *shared.Game

	Camera         *PgCamera
	World          *PgWorld
	InputHandler   *PgInputHandler
	PlayerRenderer *PgPlayerRenderer

	GameframeRight     *ebiten.Image
	GameframeBottom    *ebiten.Image
	SkillIcons         map[shared.Skill]*gebitenui.GHoverTexture
	InventoryButton    *gebitenui.GTextureButton
	SkillsButton       *gebitenui.GTextureButton
	WorldImage         *ebiten.Image
	ItemOutlineTexture *ebiten.Image
	CurrActionString   string

	Ticks uint32
}

func (p *Playground) Setup() {
	// need to update this to independent sizes when the time comes
	font16, err := gebitenui.NewGFont(AssetsDir+"font.ttf", 16)
	if err != nil {
		log.Fatalf("failed loading font: %v\n\n", err)
	}
	font18, err := gebitenui.NewGFont(AssetsDir+"font.ttf", 18)
	if err != nil {
		log.Fatalf("failed loading font: %v\n\n", err)
	}
	font20, err := gebitenui.NewGFont(AssetsDir+"font.ttf", 20)
	if err != nil {
		log.Fatalf("failed loading font: %v\n\n", err)
	}
	font24, err := gebitenui.NewGFont(AssetsDir+"font.ttf", 24)
	if err != nil {
		log.Fatalf("failed loading font: %v\n\n", err)
	}
	p.Font16 = font16
	p.Font18 = font18
	p.Font20 = font20
	p.Font24 = font24

	otherTex := loadTex(AssetsDir + "assets/other.grpgtex")
	p.Textures = loadTextures(AssetsDir + "assets/textures.grpgtex")
	p.Game.Objs = loadObjs(AssetsDir + "assets/objs.grpgobj")
	p.Game.Tiles = loadTiles(AssetsDir + "assets/tiles.grpgtile")
	p.Game.Items = loadItems(AssetsDir + "assets/items.grpgitem")
	p.Game.Npcs = loadNpcs(AssetsDir + "assets/npcs.grpgnpc")

	p.WorldImage = ebiten.NewImage(WorldImageSize, WorldImageSize)

	p.CurrActionString = "Current Action: None :("
	p.Camera = NewPgCamera(p.Game.Player)
	p.World = NewPgWorld(p.Game, p.WorldImage, p.Textures, otherTex["exclam"], p.Font16)
	p.InputHandler = NewPgInputHandler(p.Game)
	p.PlayerRenderer = NewPgPlayerRenderer(p.WorldImage, otherTex, p.Game, p.Font16)

	p.GameframeRight = otherTex["gameframe_right"]
	p.GameframeBottom = otherTex["gameframe_bottom"]

	p.SkillIcons = make(map[shared.Skill]*gebitenui.GHoverTexture)

	// TODO: refactor this into its own function at some point when i add more skills, it'll probably end up being rather manual unfortunately, don't think there's much i can do
	hoverTex := otherTex["hover_tex"]
	foragingIconTex := otherTex["foraging_icon"]
	p.ItemOutlineTexture = otherTex["item_outline"]

	p.SkillIcons[shared.Foraging] = gebitenui.NewHoverTexture(RightGameframeX+TileSize, TileSize, RightGameframeX+(TileSize*5), foragingIconTex, p.Game.SkillHoverMsgs[shared.Foraging], hoverTex, font16, color.White)

	p.InventoryButton = gebitenui.NewTextureButton(RightGameframeX+InvButtonXOffset, 0, otherTex["inv_button"], func() {
		p.Game.GameframeContainerRenderType = shared.Inventory
	})

	p.SkillsButton = gebitenui.NewTextureButton(RightGameframeX+SkillsButtonXOffset, 0, otherTex["skills_button"], func() {
		p.Game.GameframeContainerRenderType = shared.Skills
	})
}

func (p *Playground) Cleanup() {
	// would usually dispose here but gc will take care of it since this is the last scene they're on
	// NOTE: will require disposal if i start switching scenes back to login or something else
}

func (p *Playground) Update() error {
	player := p.Game.Player
	p.InputHandler.Update()

	crossedZone := player.PrevX/ChunkSize != player.ChunkX || player.PrevY/ChunkSize != player.ChunkY

	//pass crossed zone here as im already computing it for camera
	player.Update(p.Game, crossedZone)

	for _, rp := range p.Game.OtherPlayers {
		rp.Update(p.Game)
	}

	updateCurrActionString(p)
	p.Camera.Update(crossedZone)
	p.InventoryButton.Update()
	p.SkillsButton.Update()
	// TODO: inefficient?
	p.World.Update(p.Ticks)

	for _, si := range p.SkillIcons {
		si.Update()
	}

	p.Ticks++
	return nil
}

func (p *Playground) Draw(screen *ebiten.Image) {
	p.WorldImage.Clear()

	p.World.Draw()
	p.PlayerRenderer.Draw()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-p.Camera.CameraTarget.X, -p.Camera.CameraTarget.Y)
	op.Filter = ebiten.FilterNearest

	screen.DrawImage(p.WorldImage, op)

	drawGameFrame(p, screen)
}

func updateCurrActionString(p *Playground) {
	facingCoord := p.Game.Player.GetFacingCoord()
	trackedObj, objExists := p.Game.TrackedObjs[facingCoord]
	trackedNpc, npcExists := p.Game.NpcsByPos[facingCoord]
	if objExists {
		p.CurrActionString = "Current Action: " + trackedObj.DataObj.InteractText
	} else if npcExists {
		p.CurrActionString = "Talk to " + trackedNpc.NpcData.Name
	} else {
		p.CurrActionString = "Current Action: None :("
	}
}

func drawGameFrame(p *Playground, screen *ebiten.Image) {
	player := p.Game.Player
	util.DrawImage(screen, p.GameframeRight, RightGameframeX, 0)
	if p.InputHandler.IsTypingCommand {
		p.Font16.Draw(screen, "Command: "+p.InputHandler.CommandString, 0, CommandY, color.White)
	}

	if p.Game.GameframeContainerRenderType == shared.Inventory {
		var currItemRealPosX int32 = RightGameframeX + TileSize
		var currItemRealPosY int32 = TileSize

		for idx, item := range p.Game.Player.Inventory {
			if item.ItemId == 0 {
				continue
			}

			data := p.Game.Items[item.ItemId]
			tex := p.Textures[data.Texture]
			util.DrawImage(screen, tex, currItemRealPosX, currItemRealPosY)

			p.Font16.Draw(screen, fmt.Sprintf("%d", item.Count), float64(currItemRealPosX+ItemCountXOffset), float64(currItemRealPosY+ItemCountYOffset), color.White)

			if idx == p.Game.OutlineInvSpot {
				util.DrawImage(screen, p.ItemOutlineTexture, currItemRealPosX, currItemRealPosY)
			}

			currItemRealPosX += TileSize

			if (idx+1)%ItemsPerRow == 0 {
				currItemRealPosY += TileSize
				currItemRealPosX = RightGameframeX + TileSize
			}
		}
	} else if p.Game.GameframeContainerRenderType == shared.Skills {
		for _, si := range p.SkillIcons {
			si.Draw(screen)
		}
		for i := shared.Foraging; i <= shared.Foraging; i++ {
			// TODO: maybe string can be pre computed by packet here?
			// TODO: magic constants when i do this det
			p.Font16.Draw(screen, fmt.Sprintf("%d", p.Game.Skills[i].Level), RightGameframeX+64+32, 64+48, util.Yellow)
		}
	}

	util.DrawImage(screen, p.GameframeBottom, 0, RightGameframeX)

	talkbox := p.Game.Talkbox
	// x is offset from 0, y has offset added, to be placed in the right spot
	p.Font20.Draw(screen, p.CurrActionString, CurrActionX, RightGameframeX+CurrNameActionYOffset, color.White)
	if talkbox.Active {
		p.Font24.Draw(screen, talkbox.CurrentName, CurrNameX, RightGameframeX+CurrNameActionYOffset, color.White)
		var currY float64 = CurrMessageY
		for _, s := range talkbox.CurrentMessage {
			p.Font24.Draw(screen, s, CurrMessageX, currY, color.White)
			currY += 30
		}
	}
	p.InventoryButton.Draw(screen)
	p.SkillsButton.Draw(screen)

	playerCoords := fmt.Sprintf("X: %d, Y: %d, Facing: %s", player.X, player.Y, player.Facing.String())
	p.Font24.Draw(screen, playerCoords, RightGameframeX, DebugCoordsY, color.White)
}
