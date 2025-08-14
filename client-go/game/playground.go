package game

import (
	"client/shared"
	"client/util"
	"fmt"
	"grpg/data-go/grpgmap"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Playground struct {
	Font             rl.Font
	Game             *shared.Game
	GameframeRight   rl.Texture2D
	PlayerTextures   map[shared.Direction]rl.Texture2D
	Textures         map[uint16]rl.Texture2D
	Zones            map[util.Vector2I]grpgmap.Zone
	CameraTarget     rl.Vector2
	PrevCameraTarget rl.Vector2
	CurrActionString string
}

var assetsDirectory = "../../grpg-assets/"

func (p *Playground) Setup() {
	// TODO: move font out
	p.Font = rl.LoadFont("./assets/font.ttf")
	p.CurrActionString = "Current Action: None :("

	p.Textures = loadTextures(assetsDirectory + "textures.grpgtex")
	p.Game.Objs = loadObjs(assetsDirectory + "objs.grpgobj")
	p.Game.Tiles = loadTiles(assetsDirectory + "tiles.grpgtile")
	p.Game.Items = loadItems(assetsDirectory + "items.grpgitem")
	p.Zones = loadMaps(assetsDirectory+"maps/", p.Game)

	p.GameframeRight = loadGameframeRightTexture(assetsDirectory + "used/gameframe_right_2.png")

	p.PlayerTextures = loadPlayerTextures(assetsDirectory + "used/")
}

func (p *Playground) Cleanup() {
	if p.Font.Texture.ID != 0 {
		rl.UnloadFont(p.Font)
	}

	for _, tex := range p.Textures {
		rl.UnloadTexture(tex)
	}

	rl.UnloadTexture(p.GameframeRight)

	for _, tex := range p.PlayerTextures {
		rl.UnloadTexture(tex)
	}
}

func (p *Playground) Loop() {
	player := p.Game.Player

	if rl.IsKeyPressed(rl.KeyW) {
		player.SendMovePacket(p.Game, player.X, player.Y-1, shared.UP)
	} else if rl.IsKeyPressed(rl.KeyS) {
		player.SendMovePacket(p.Game, player.X, player.Y+1, shared.DOWN)
	} else if rl.IsKeyPressed(rl.KeyA) {
		player.SendMovePacket(p.Game, player.X-1, player.Y, shared.LEFT)
	} else if rl.IsKeyPressed(rl.KeyD) {
		player.SendMovePacket(p.Game, player.X+1, player.Y, shared.RIGHT)
	} else if rl.IsKeyPressed(rl.KeyQ) {
		player.SendInteractPacket(p.Game)
	}

	crossedZone := player.PrevX/16 != player.ChunkX || player.PrevY/16 != player.ChunkY

	// pass crossed zone here as im already computing it for camera
	player.Update(p.Game, crossedZone)

	for _, rp := range p.Game.OtherPlayers {
		rp.Update(p.Game)
	}

	updateCurrActionString(p)

	// needs to be done last but crossed zone check must be doing before player is updated as that changes prev x/y
	updateCamera(p, crossedZone)
}

func (p *Playground) Render() {
	rl.ClearBackground(rl.Black)

	camera := rl.Camera2D{
		Offset:   rl.Vector2{X: 0, Y: 0},
		Target:   p.CameraTarget,
		Rotation: 0,
		Zoom:     1,
	}

	rl.BeginMode2D(camera)

	drawWorld(p)
	drawOtherPlayers(p)
	drawPlayer(p)

	rl.EndMode2D()

	drawGameFrame(p)
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
		p.CameraTarget.X = float32(cameraX)
		p.CameraTarget.Y = float32(cameraY)
	} else {
		if p.CameraTarget.X < float32(cameraX) {
			p.CameraTarget.X += float32(speed)
		} else if p.CameraTarget.X > float32(cameraX) {
			p.CameraTarget.X -= float32(speed)
		}

		if p.CameraTarget.Y < float32(cameraY) {
			p.CameraTarget.Y += float32(speed)
		} else if p.CameraTarget.Y > float32(cameraY) {
			p.CameraTarget.Y -= float32(speed)
		}
	}

	p.PrevCameraTarget = p.CameraTarget
}

func updateCurrActionString(p *Playground) {
	facingCoord := p.Game.Player.GetFacingCoord()
	trackedObj, exists := p.Game.TrackedObjs[facingCoord]
	if !exists {
		p.CurrActionString = "None :("
	} else {
		p.CurrActionString = trackedObj.DataObj.InteractText
	}
}

func drawWorld(p *Playground) {
	player := p.Game.Player

	mapTiles := p.Zones[util.Vector2I{X: p.Game.Player.ChunkX, Y: p.Game.Player.ChunkY}]

	for i := range 256 {
		localX := int32(i % 16)
		localY := int32(i / 16)

		dx := localX * p.Game.TileSize
		dy := localY * p.Game.TileSize

		texId := p.Game.Tiles[uint16(mapTiles.Tiles[i])].TexId

		tex := p.Textures[texId]
		rl.DrawTexture(tex, dx, dy, rl.White)

		obj := mapTiles.Objs[i]
		if obj != 0 {
			trackedObj, ok := p.Game.TrackedObjs[util.Vector2I{
				X: localX + (player.ChunkX * 16),
				Y: localY + (player.ChunkY * 16),
			}]

			// fallback pretty much, might not be necessary in the future
			var state uint16 = 0
			if ok {
				state = uint16(trackedObj.State)
			}

			objTexId := p.Game.Objs[uint16(mapTiles.Objs[i])].Textures[state]

			objTex := p.Textures[objTexId]
			rl.DrawTexture(objTex, dx, dy, rl.White)
		}
	}
}

// TODO: generalize this code
func drawPlayer(p *Playground) {
	player := p.Game.Player

	rl.DrawTexture(p.PlayerTextures[player.Facing], player.RealX, player.RealY, rl.White)
	rl.DrawTextEx(
		p.Font,
		player.Name,
		rl.Vector2{X: float32(player.RealX), Y: float32(player.RealY)},
		16,
		0,
		rl.White,
	)
}

func drawOtherPlayers(p *Playground) {
	for _, player := range p.Game.OtherPlayers {
		rl.DrawTexture(p.PlayerTextures[player.Facing], player.RealX, player.RealY, rl.White)
		rl.DrawTextEx(
			p.Font,
			player.Name,
			rl.Vector2{X: float32(player.RealX), Y: float32(player.RealY)},
			16,
			0,
			rl.Red,
		)
	}
}

func drawGameFrame(p *Playground) {
	player := p.Game.Player

	rl.DrawTexture(p.GameframeRight, 768, 0, rl.White)

	currItemRealPos := rl.Vector2{X: 768 + 64, Y: 64}

	for idx, item := range p.Game.Player.Inventory {
		if item.ItemId == 0 {
			continue
		}

		data := p.Game.Items[item.ItemId]
		tex := p.Textures[data.Texture]
		rl.DrawTexture(tex, int32(currItemRealPos.X), int32(currItemRealPos.Y), rl.White)

		textPos := rl.Vector2Add(currItemRealPos, rl.Vector2{X: 16, Y: 0})
		rl.DrawTextEx(p.Font, fmt.Sprintf("%d", item.Count), textPos, 18, 0, rl.White)

		currItemRealPos.X += 64
		if (idx+1)%4 == 0 {
			currItemRealPos.Y += 64
			currItemRealPos.X = 768 + 64
		}
	}

	rl.DrawRectangle(0, 768, 960-192, 192, rl.Blue)

	rl.DrawTextEx(p.Font, "Current Action: "+p.CurrActionString, rl.Vector2{X: 0, Y: 768}, 24, 0, rl.White)

	playerCoords := fmt.Sprintf("X: %d, Y: %d, Facing: %s", player.X, player.Y, player.Facing.String())
	rl.DrawTextEx(p.Font, playerCoords, rl.Vector2{X: 0, Y: 800}, 24, 0, rl.White)
}
