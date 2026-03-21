package game

import (
	"client/constants"
	"client/network/c2s"
	"client/util"
	"fmt"
	"grpg/data-go/grpgmap"
	"image"
	"image/color"
	"log"

	"client/shared"

	gebitenui "github.com/grian32/gebiten-ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	WorldImageSize = 1024
	RightGameframeX = 768

	ChunkSize = 16
	TilesPerChunk = 256
	TileSize = 64 // pixels per tile, this is used in both ui and actual world stuff

	ExclamAnimTickInterval = 20
	ExclamAnimFrameCount = 2
	ExclamBobOffset = -4

	InvButtonXOffset = TileSize + 16
	SkillsButtonXOffset = TileSize * 2 + 32

	CameraOffsetTiles = 4
	CameraBoundaryTiles = 12
	CameraMinOffsetTiles = 9
	CameraPanSpeed = 16.0

	CommandY = 740

	ItemCountXOffset = 6
	ItemCountYOffset = 4
	ItemsPerRow = 4

	CurrActionX = 110
	CurrNameActionYOffset = 28 + 3
	CurrNameX = CurrActionX+332
	CurrMessageX = 90
	CurrMessageY = 840

	DebugCoordsY = 800
)

type Playground struct {
	Font16             *gebitenui.GFont
	Font18             *gebitenui.GFont
	Font20             *gebitenui.GFont
	Font24             *gebitenui.GFont

	Camera             *PgCamera

	Game               *shared.Game
	GameframeRight     *ebiten.Image
	GameframeBottom    *ebiten.Image
	ExclamTexture      *ebiten.Image
	SkillIcons         map[shared.Skill]*gebitenui.GHoverTexture
	InventoryButton    *gebitenui.GTextureButton
	SkillsButton       *gebitenui.GTextureButton
	PlayerTextures     map[shared.Direction]*ebiten.Image
	Textures           map[uint16]*ebiten.Image
	Zones              map[util.Vector2I]grpgmap.Zone
	WorldImage         *ebiten.Image
	ItemOutlineTexture *ebiten.Image
	CurrActionString   string
	IsTypingCommand    bool
	CommandString      string
	ExclamYOffset      int32
	Ticks              uint32
}

func (p *Playground) Setup() {
	var assetsDirectory = "../../grpg-assets/"

	// need to update this to independent sizes when the time comes
	font16, err := gebitenui.NewGFont(assetsDirectory+"font.ttf", 16)
	if err != nil {
		log.Fatalf("failed loading font: %v\n\n", err)
	}
	font18, err := gebitenui.NewGFont(assetsDirectory+"font.ttf", 18)
	if err != nil {
		log.Fatalf("failed loading font: %v\n\n", err)
	}
	font20, err := gebitenui.NewGFont(assetsDirectory+"font.ttf", 20)
	if err != nil {
		log.Fatalf("failed loading font: %v\n\n", err)
	}
	font24, err := gebitenui.NewGFont(assetsDirectory+"font.ttf", 24)
	if err != nil {
		log.Fatalf("failed loading font: %v\n\n", err)
	}
	p.Font16 = font16
	p.Font18 = font18
	p.Font20 = font20
	p.Font24 = font24

	p.CurrActionString = "Current Action: None :("
	p.Camera = NewPgCamera(p.Game.Player)

	p.Textures = loadTextures(assetsDirectory + "assets/textures.grpgtex")
	p.Game.Objs = loadObjs(assetsDirectory + "assets/objs.grpgobj")
	p.Game.Tiles = loadTiles(assetsDirectory + "assets/tiles.grpgtile")
	p.Game.Items = loadItems(assetsDirectory + "assets/items.grpgitem")
	p.Game.Npcs = loadNpcs(assetsDirectory + "assets/npcs.grpgnpc")
	p.Zones = loadMaps(assetsDirectory+"maps/", p.Game)

	otherTex := loadTex(assetsDirectory + "assets/other.grpgtex")

	p.GameframeRight = otherTex["gameframe_right"]
	p.GameframeBottom = otherTex["gameframe_bottom"]
	p.ExclamTexture = otherTex["exclam"]
	p.ExclamYOffset = 0

	p.PlayerTextures = make(map[shared.Direction]*ebiten.Image)
	p.PlayerTextures[shared.UP] = otherTex["player_up"]
	p.PlayerTextures[shared.DOWN] = otherTex["player_down"]
	p.PlayerTextures[shared.LEFT] = otherTex["player_left"]
	p.PlayerTextures[shared.RIGHT] = otherTex["player_right"]

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

	p.WorldImage = ebiten.NewImage(WorldImageSize, WorldImageSize)
}

func (p *Playground) Cleanup() {
	// would usually dispose here but gc will take care of it since this is the last scene they're on
	// NOTE: will require disposal if i start switching scenes back to login or something else
}

func (p *Playground) Update() error {
	player := p.Game.Player

	if p.IsTypingCommand {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			p.IsTypingCommand = false
			p.CommandString = ""
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			p.IsTypingCommand = false
			player.SendCmdPacket(p.Game, p.CommandString)
			p.CommandString = ""
		} else if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(p.CommandString) != 0 {
			p.CommandString = p.CommandString[:len(p.CommandString)-1]
		} else {
			p.CommandString = string(ebiten.AppendInputChars([]rune(p.CommandString)))
		}
	} else {
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			player.SendMovePacket(p.Game, player.X, player.Y-1, shared.UP)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			player.SendMovePacket(p.Game, player.X, player.Y+1, shared.DOWN)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			player.SendMovePacket(p.Game, player.X-1, player.Y, shared.LEFT)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			player.SendMovePacket(p.Game, player.X+1, player.Y, shared.RIGHT)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			player.SendInteractPacket(p.Game)
		} else if p.Game.Talkbox.Active && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			shared.SendPacket(p.Game.Conn, &c2s.Continue{})
		} else if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			p.IsTypingCommand = true
		}
	}

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
	if (p.Ticks/ExclamAnimTickInterval)%ExclamAnimFrameCount == 0 {
		p.ExclamYOffset = 0
	} else {
		p.ExclamYOffset = ExclamBobOffset
	}

	for _, si := range p.SkillIcons {
		si.Update()
	}

	p.Ticks++
	return nil
}

func (p *Playground) Draw(screen *ebiten.Image) {
	p.WorldImage.Clear()

	drawWorld(p, p.WorldImage)
	drawOtherPlayers(p, p.WorldImage)
	drawPlayer(p, p.WorldImage)

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

func drawWorld(p *Playground, screen *ebiten.Image) {
	player := p.Game.Player

	mapTiles := p.Zones[util.Vector2I{X: p.Game.Player.ChunkX, Y: p.Game.Player.ChunkY}]

	for i := range TilesPerChunk {
		localX := int32(i % ChunkSize)
		localY := int32(i / ChunkSize)

		dx := localX * p.Game.TileSize
		dy := localY * p.Game.TileSize

		texId := p.Game.Tiles[uint16(mapTiles.Tiles[i])].TexId

		tex := p.Textures[texId]

		util.DrawImage(screen, tex, dx, dy)

		worldPos := util.Vector2I{
			X: localX + (player.ChunkX * ChunkSize),
			Y: localY + (player.ChunkY * ChunkSize),
		}

		obj := mapTiles.Objs[i]
		if obj != 0 {
			trackedObj, ok := p.Game.TrackedObjs[worldPos]

			// fallback pretty much, might not be necessary in the future
			var state uint16 = 0
			if ok {
				state = uint16(trackedObj.State)
			}

			objTexId := p.Game.Objs[uint16(mapTiles.Objs[i])].Textures[state]

			objTex := p.Textures[objTexId]
			util.DrawImage(screen, objTex, dx, dy)
		}
		if trackedNpc, ok := p.Game.NpcsByPos[worldPos]; ok {
			// TODO: maybe don't render if player is standing over?
			npcTexId := trackedNpc.NpcData.TextureId
			util.DrawImage(screen, p.Textures[npcTexId], dx, dy)
			if p.Game.DebugMode {
				p.Font16.Draw(screen, fmt.Sprintf("%d", trackedNpc.Uid), float64(dx), float64(dy), color.White)
			}

			if trackedNpc.NpcData.NpcId == uint16(constants.GRPG_GUIDE) && p.Game.RenderExclamOnGuide {
				util.DrawImage(screen, p.ExclamTexture, dx, dy-TileSize+p.ExclamYOffset)
			}
		}
	}
}

func drawPlayer(p *Playground, screen *ebiten.Image) {
	player := p.Game.Player

	srcX := int(player.CurrFrame) * TileSize
	sourceRec := image.Rectangle{
		Min: image.Point{
			X: srcX,
			Y: 0,
		},
		Max: image.Point{
			X: srcX + TileSize,
			Y: TileSize,
		},
	}
	sub := util.SubImage(p.PlayerTextures[player.Facing], sourceRec)
	util.DrawImage(screen, sub, player.RealX, player.RealY)

	p.Font16.Draw(screen, player.Name, float64(player.RealX), float64(player.RealY), color.White)
}

func drawOtherPlayers(p *Playground, screen *ebiten.Image) {
	for _, player := range p.Game.OtherPlayers {
		srcX := int(player.CurrFrame) * TileSize
		sourceRec := image.Rectangle{
			Min: image.Point{
				X: srcX,
				Y: 0,
			},
			Max: image.Point{
				X: srcX + TileSize,
				Y: TileSize,
			},
		}
		sub := util.SubImage(p.PlayerTextures[player.Facing], sourceRec)
		util.DrawImage(screen, sub, player.RealX, player.RealY)

		p.Font16.Draw(screen, player.Name, float64(player.RealX), float64(player.RealY), util.Red)
	}
}

func drawGameFrame(p *Playground, screen *ebiten.Image) {
	player := p.Game.Player
	util.DrawImage(screen, p.GameframeRight, RightGameframeX, 0)
	if p.IsTypingCommand {
		p.Font16.Draw(screen, "Command: "+p.CommandString, 0, CommandY, color.White)
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

			p.Font16.Draw(screen, fmt.Sprintf("%d", item.Count), float64(currItemRealPosX+ItemCountXOffset), float64(currItemRealPosY + ItemCountYOffset), color.White)

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
