package game

import (
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

type Playground struct {
	Font16          *gebitenui.GFont
	Font18          *gebitenui.GFont
	Font20          *gebitenui.GFont
	Font24          *gebitenui.GFont
	Game            *shared.Game
	GameframeRight  *ebiten.Image
	GameframeBottom *ebiten.Image
	SkillIcons      map[shared.Skill]*gebitenui.GHoverTexture
	InventoryButton  *gebitenui.GTextureButton
	SkillsButton     *gebitenui.GTextureButton
	PlayerTextures   map[shared.Direction]*ebiten.Image
	Textures         map[uint16]*ebiten.Image
	Zones            map[util.Vector2I]grpgmap.Zone
	CameraTarget     util.Vector2
	PrevCameraTarget util.Vector2
	CurrActionString string
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

	p.Textures = loadTextures(assetsDirectory + "assets/textures.grpgtex")
	p.Game.Objs = loadObjs(assetsDirectory + "assets/objs.grpgobj")
	p.Game.Tiles = loadTiles(assetsDirectory + "assets/tiles.grpgtile")
	p.Game.Items = loadItems(assetsDirectory + "assets/items.grpgitem")
	p.Game.Npcs = loadNpcs(assetsDirectory + "assets/npcs.grpgnpc")
	p.Zones = loadMaps(assetsDirectory+"maps/", p.Game)

	otherTex := loadTex(assetsDirectory + "assets/other.grpgtex")

	p.GameframeRight = otherTex["gameframe_right"]
	p.GameframeBottom = otherTex["gameframe_bottom"]

	p.PlayerTextures = make(map[shared.Direction]*ebiten.Image)
	p.PlayerTextures[shared.UP] = otherTex["player_up"]
	p.PlayerTextures[shared.DOWN] = otherTex["player_down"]
	p.PlayerTextures[shared.LEFT] = otherTex["player_left"]
	p.PlayerTextures[shared.RIGHT] = otherTex["player_right"]

	p.SkillIcons = make(map[shared.Skill]*gebitenui.GHoverTexture)

	// TODO: refactor this into its own function at some point when i add more skills, it'll probably end up being rather manual unfortunately, don't think there's much i can do
	hoverTex := otherTex["hover_tex"]
	foragingIconTex := otherTex["foraging_icon"]

	p.SkillIcons[shared.Foraging] = gebitenui.NewHoverTexture(768+64, 64, 768+(64*5), foragingIconTex, p.Game.SkillHoverMsgs[shared.Foraging], hoverTex, font16, color.White)

	p.InventoryButton = gebitenui.NewTextureButton(768+64+16, 0, otherTex["inv_button"], func() {
		p.Game.GameframeContainerRenderType = shared.Inventory
	})

	p.SkillsButton = gebitenui.NewTextureButton(768+128+32, 0, otherTex["skills_button"], func() {
		p.Game.GameframeContainerRenderType = shared.Skills
	})
}

func (p *Playground) Cleanup() {
	// would usually dispose here but gc will take care of it since this is the last scene they're on
	// NOTE: will require disposal if i start switching scenes back to login or something else
}

func (p *Playground) Update() error {
	player := p.Game.Player

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
	}

	crossedZone := player.PrevX/16 != player.ChunkX || player.PrevY/16 != player.ChunkY

	//pass crossed zone here as im already computing it for camera
	player.Update(p.Game, crossedZone)

	for _, rp := range p.Game.OtherPlayers {
		rp.Update(p.Game)
	}

	updateCurrActionString(p)
	updateCamera(p, crossedZone)
	p.InventoryButton.Update()
	p.SkillsButton.Update()

	for _, si := range p.SkillIcons {
		si.Update()
	}
	return nil
}

func (p *Playground) Draw(screen *ebiten.Image) {
	worldImage := ebiten.NewImage(1024, 1024)

	drawWorld(p, worldImage)
	drawOtherPlayers(p, worldImage)
	drawPlayer(p, worldImage)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-p.CameraTarget.X, -p.CameraTarget.Y)
	op.Filter = ebiten.FilterNearest

	screen.DrawImage(worldImage, op)

	drawGameFrame(p, screen)
}

func updateCamera(p *Playground, crossedZone bool) {
	player := p.Game.Player

	var cameraX = 4 * p.Game.TileSize
	var cameraY = 4 * p.Game.TileSize

	if player.RealX <= 12*p.Game.TileSize {
		cameraX = util.MinI(player.RealX-(9*p.Game.TileSize), 0)
	}

	if player.RealY <= 12*p.Game.TileSize {
		cameraY = util.MinI(player.RealY-(9*p.Game.TileSize), 0)
	}

	const speed = 16.0

	if crossedZone {
		p.CameraTarget.X = float64(cameraX)
		p.CameraTarget.Y = float64(cameraY)
	} else {
		if p.CameraTarget.X < float64(cameraX) {
			p.CameraTarget.X += speed
		} else if p.CameraTarget.X > float64(cameraX) {
			p.CameraTarget.X -= speed
		}

		if p.CameraTarget.Y < float64(cameraY) {
			p.CameraTarget.Y += speed
		} else if p.CameraTarget.Y > float64(cameraY) {
			p.CameraTarget.Y -= speed
		}
	}

	p.PrevCameraTarget = p.CameraTarget
}

func updateCurrActionString(p *Playground) {
	facingCoord := p.Game.Player.GetFacingCoord()
	trackedObj, objExists := p.Game.TrackedObjs[facingCoord]
	trackedNpc, npcExists := p.Game.TrackedNpcs[facingCoord]
	if objExists {
		p.CurrActionString = "Current Action: " + trackedObj.DataObj.InteractText
	} else if npcExists {
		p.CurrActionString = "Talk to " + trackedNpc.NpcData.Name
	} else {
		p.CurrActionString = "Curret Action: None :("
	}
}

func drawWorld(p *Playground, screen *ebiten.Image) {
	player := p.Game.Player

	mapTiles := p.Zones[util.Vector2I{X: p.Game.Player.ChunkX, Y: p.Game.Player.ChunkY}]

	for i := range 256 {
		localX := int32(i % 16)
		localY := int32(i / 16)

		dx := localX * p.Game.TileSize
		dy := localY * p.Game.TileSize

		texId := p.Game.Tiles[uint16(mapTiles.Tiles[i])].TexId

		tex := p.Textures[texId]

		util.DrawImage(screen, tex, dx, dy)

		worldPos := util.Vector2I{
			X: localX + (player.ChunkX * 16),
			Y: localY + (player.ChunkY * 16),
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
		} else {
			// TODO: maybe don't render if player is standing over?
			trackedNpc, ok := p.Game.TrackedNpcs[worldPos]
			if ok {
				npcTexId := trackedNpc.NpcData.TextureId
				util.DrawImage(screen, p.Textures[npcTexId], dx, dy)
			}
		}
	}
}

func drawPlayer(p *Playground, screen *ebiten.Image) {
	player := p.Game.Player

	const frameSize = 64
	srcX := int(player.CurrFrame) * frameSize
	sourceRec := image.Rectangle{
		Min: image.Point{
			X: srcX,
			Y: 0,
		},
		Max: image.Point{
			X: srcX + frameSize,
			Y: frameSize,
		},
	}
	sub := util.SubImage(p.PlayerTextures[player.Facing], sourceRec)
	util.DrawImage(screen, sub, player.RealX, player.RealY)

	p.Font16.Draw(screen, player.Name, float64(player.RealX), float64(player.RealY), color.White)
}

func drawOtherPlayers(p *Playground, screen *ebiten.Image) {
	for _, player := range p.Game.OtherPlayers {
		const frameSize = 64
		srcX := int(player.CurrFrame) * frameSize
		sourceRec := image.Rectangle{
			Min: image.Point{
				X: srcX,
				Y: 0,
			},
			Max: image.Point{
				X: srcX + frameSize,
				Y: frameSize,
			},
		}
		sub := util.SubImage(p.PlayerTextures[player.Facing], sourceRec)
		util.DrawImage(screen, sub, player.RealX, player.RealY)

		p.Font16.Draw(screen, player.Name, float64(player.RealX), float64(player.RealY), util.Red)
	}
}

func drawGameFrame(p *Playground, screen *ebiten.Image) {
	player := p.Game.Player
	util.DrawImage(screen, p.GameframeRight, 768, 0)

	if p.Game.GameframeContainerRenderType == shared.Inventory {
		var currItemRealPosX int32 = 768 + 64
		var currItemRealPosY int32 = 64

		for idx, item := range p.Game.Player.Inventory {
			if item.ItemId == 0 {
				continue
			}

			data := p.Game.Items[item.ItemId]
			tex := p.Textures[data.Texture]
			util.DrawImage(screen, tex, currItemRealPosX, currItemRealPosY)

			p.Font16.Draw(screen, fmt.Sprintf("%d", item.Count), float64(currItemRealPosX+16), float64(currItemRealPosY), color.White)

			currItemRealPosX += 64
			if (idx+1)%4 == 0 {
				currItemRealPosY += 64
				currItemRealPosX = 768 + 64
			}
		}
	} else if p.Game.GameframeContainerRenderType == shared.Skills {
		for _, si := range p.SkillIcons {
			si.Draw(screen)
		}
		for i := shared.Foraging; i <= shared.Foraging; i++ {
			// TODO: maybe string can be pre computed by packet here?
			p.Font16.Draw(screen, fmt.Sprintf("%d", p.Game.Skills[i].Level), 768+64+32, 64+48, util.Yellow)
		}
	}

	util.DrawImage(screen, p.GameframeBottom, 0, 768)

	talkbox := p.Game.Talkbox
	// x is offset from 0, y has offset added, to be placed in the right spot
	p.Font20.Draw(screen, p.CurrActionString, 110, 768+28+3, color.White)
	if talkbox.Active {
		p.Font24.Draw(screen, talkbox.CurrentName, 110+332, 768+28+3, color.White)
		p.Font24.Draw(screen, talkbox.CurrentMessage, 90, 840, color.White)
	}
	p.InventoryButton.Draw(screen)
	p.SkillsButton.Draw(screen)

	playerCoords := fmt.Sprintf("X: %d, Y: %d, Facing: %s", player.X, player.Y, player.Facing.String())
	p.Font24.Draw(screen, playerCoords, 768, 800, color.White)
}
