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

	p.Textures = loadTextures(assetsDirectory + "textures.pak")
	p.Zones = loadMaps(assetsDirectory+"maps/", p.Game)
	p.GameframeRight = loadGameframeRightTexture(assetsDirectory + "used/gameframe_right_2.png")

	p.PlayerTextures = loadPlayerTextures(assetsDirectory + "used/")
}

func (p *Playground) Cleanup() {
	rl.UnloadFont(p.Font)
	// TODO: unload all textures :S
}

func (p *Playground) Loop() {
	player := p.Game.Player

	if rl.IsKeyPressed(rl.KeyW) {
		player.Facing = shared.UP
		player.SendMovePacket(p.Game, player.X, player.Y-1)
	} else if rl.IsKeyPressed(rl.KeyS) {
		player.Facing = shared.DOWN
		player.SendMovePacket(p.Game, player.X, player.Y+1)
	} else if rl.IsKeyPressed(rl.KeyA) {
		player.Facing = shared.LEFT
		player.SendMovePacket(p.Game, player.X-1, player.Y)
	} else if rl.IsKeyPressed(rl.KeyD) {
		player.Facing = shared.RIGHT
		player.SendMovePacket(p.Game, player.X+1, player.Y)
	}

	player.Update(p.Game)

	for _, rp := range p.Game.OtherPlayers {
		rp.Update(p.Game)
	}

	var cameraX = 4 * p.Game.TileSize
	var cameraY = 4 * p.Game.TileSize

	if player.RealX <= 12*p.Game.TileSize {
		cameraX = util.MinI(player.RealX-(9*p.Game.TileSize), 0)
	}

	if player.RealY <= 12*p.Game.TileSize {
		cameraY = util.MinI(player.RealY-(9*p.Game.TileSize), 0)
	}

	const speed = 16.0

	crossedZone := player.PrevX/16 != player.ChunkX || player.PrevY/16 != player.ChunkY

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

func drawWorld(p *Playground) {
	mapTiles := p.Zones[util.Vector2I{X: p.Game.Player.ChunkX, Y: p.Game.Player.ChunkY}]

	for i := range 256 {
		dx := int32(i%16) * p.Game.TileSize
		dy := int32(i/16) * p.Game.TileSize

		tex := p.Textures[uint16(mapTiles.Tiles[i])]
		rl.DrawTexture(tex, dx, dy, rl.White)

		obj := mapTiles.Objs[i]
		if obj.InternalId != 0 && obj.Type == grpgmap.OBJ {
			objTex := p.Textures[obj.InternalId]
			rl.DrawTexture(objTex, dx, dy, rl.White)
		}
	}
}

// TODO: generalize this code
func drawPlayer(p *Playground) {
	// rl.DrawRectangle(p.Game.Player.RealX, p.Game.Player.RealY, 64, 64, rl.SkyBlue)
	rl.DrawTexture(p.PlayerTextures[p.Game.Player.Facing], p.Game.Player.RealX, p.Game.Player.RealY, rl.White)
	rl.DrawTextEx(
		p.Font,
		p.Game.Player.Name,
		rl.Vector2{X: float32(p.Game.Player.RealX), Y: float32(p.Game.Player.RealY)},
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
	// rl.DrawRectangle(768, 0, 320, 960, rl.Blue)
	// rl.DrawTextEx(p.Font, "inventory or something", rl.Vector2{X: 768, Y: 0}, 24, 0, rl.White)
	rl.DrawTexture(p.GameframeRight, 768, 0, rl.White)
	rl.DrawRectangle(0, 768, 960-192, 192, rl.Blue)
	rl.DrawTextEx(p.Font, p.CurrActionString, rl.Vector2{X: 0, Y: 768}, 24, 0, rl.White)
	playerCoords := fmt.Sprintf("X: %d, Y: %d", p.Game.Player.X, p.Game.Player.Y)
	rl.DrawTextEx(p.Font, playerCoords, rl.Vector2{X: 0, Y: 800}, 24, 0, rl.White)
}
