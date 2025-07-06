package game

import (
	"client/shared"
	"client/util"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Playground struct {
	Font rl.Font
	Game *shared.Game
}

func (p *Playground) Setup() {
	p.Font = rl.LoadFont("./assets/font.ttf")
}

func (p *Playground) Cleanup() {
	rl.UnloadFont(p.Font)
}

func (p *Playground) Loop() {
	player := p.Game.Player

	// TODO: figure out some way to send remove duplicate packet code
	if rl.IsKeyPressed(rl.KeyW) {
		player.Move(player.X, player.Y-1, p.Game)
		player.SendMovePacket(p.Game)
	} else if rl.IsKeyPressed(rl.KeyS) {
		player.Move(player.X, player.Y+1, p.Game)
		player.SendMovePacket(p.Game)
	} else if rl.IsKeyPressed(rl.KeyA) {
		player.Move(player.X-1, player.Y, p.Game)
		player.SendMovePacket(p.Game)
	} else if rl.IsKeyPressed(rl.KeyD) {
		player.Move(player.X+1, player.Y, p.Game)
		player.SendMovePacket(p.Game)
	}
}

func (p *Playground) Render() {
	rl.ClearBackground(rl.Black)

	player := p.Game.Player

	var cameraX = 4 * p.Game.TileSize
	var cameraY = 4 * p.Game.TileSize

	// eh just hardcode these prob
	if player.RealX <= 12*p.Game.TileSize {
		cameraX = util.MinI(player.RealX-(9*p.Game.TileSize), 0)
	}

	if player.RealY <= 12*p.Game.TileSize {
		cameraY = util.MinI(player.RealY-(9*p.Game.TileSize), 0)
	}

	camera := rl.Camera2D{
		Offset:   rl.Vector2{X: 0, Y: 0},
		Target:   rl.Vector2{X: float32(cameraX), Y: float32(cameraY)},
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
	for x := range 16 {
		for y := range 16 {
			dx := int32(x) * p.Game.TileSize
			dy := int32(y) * p.Game.TileSize

			if y%2 == 0 || x%2 == 0 {
				rl.DrawRectangle(dx, dy, p.Game.TileSize, p.Game.TileSize, rl.White)
			} else {
				rl.DrawRectangle(dx, dy, p.Game.TileSize, p.Game.TileSize, rl.Gray)
			}
			rl.DrawRectangleLines(dx, dy, p.Game.TileSize, p.Game.TileSize, rl.Black)

			if p.Game.Player.ChunkY == 1 && x == 8 && y == 8 {
				rl.DrawRectangle(dx, dy, p.Game.TileSize, p.Game.TileSize, rl.Red)
			}
		}
	}
}

// TODO: generalize this code
func drawPlayer(p *Playground) {
	rl.DrawRectangle(p.Game.Player.RealX, p.Game.Player.RealY, 64, 64, rl.SkyBlue)
	rl.DrawTextEx(
		p.Font,
		p.Game.Player.Name,
		rl.Vector2{X: float32(p.Game.Player.RealX), Y: float32(p.Game.Player.RealY)},
		16,
		0,
		rl.Red,
	)
}

func drawOtherPlayers(p *Playground) {
	for _, player := range p.Game.OtherPlayers {
		rl.DrawRectangle(player.RealX, player.RealY, 64, 64, rl.SkyBlue)
		rl.DrawTextEx(
			p.Font,
			player.Name,
			rl.Vector2{X: float32(p.Game.Player.RealX), Y: float32(p.Game.Player.RealY)},
			16,
			0,
			rl.Red,
		)
	}
}

func drawGameFrame(p *Playground) {
	rl.DrawRectangle(768, 0, 192, 960, rl.Blue)
	rl.DrawTextEx(p.Font, "inventory or something", rl.Vector2{X: 768, Y: 0}, 24, 0, rl.White)
	rl.DrawRectangle(0, 768, 960-192, 192, rl.Blue)
	rl.DrawTextEx(p.Font, "something else eventually", rl.Vector2{X: 0, Y: 768}, 24, 0, rl.White)
	playerCoords := fmt.Sprintf("X: %d, Y: %d", p.Game.Player.X, p.Game.Player.Y)
	rl.DrawTextEx(p.Font, playerCoords, rl.Vector2{X: 0, Y: 800}, 24, 0, rl.White)
}
