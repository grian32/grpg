package game

import (
	"client/shared"
	"client/util"
	"cmp"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgmap"
	"grpg/data-go/grpgtex"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Playground struct {
	Font     rl.Font
	Game     *shared.Game
	Textures map[uint16]rl.Texture2D
	Maps     map[util.Vector2I][256]grpgmap.Tile
}

func (p *Playground) Setup() {
	p.Font = rl.LoadFont("./assets/font.ttf")

	p.Textures = make(map[uint16]rl.Texture2D)
	p.Maps = make(map[util.Vector2I][256]grpgmap.Tile)

	grpgTexFile, err1 := os.Open("../../grpg-assets/textures.pak")
	grpgTexBytes, err2 := io.ReadAll(grpgTexFile)

	buf := gbuf.NewGBuf(grpgTexBytes)
	header := grpgtex.ReadHeader(buf)
	if string(header.Magic[:]) != "GRPGTEX\x00" {
		log.Fatal("File is not GRPGTEX file.")
	}

	log.Printf("Succesfully loaded GRPGTEX file with version %d\n", header.Version)

	textures := grpgtex.ReadTextures(buf)
	if err := cmp.Or(err1, err2); err != nil {
		log.Fatal("Failed reading GRPGTEX file")
	}

	for _, tex := range textures {
		rlImage := rl.LoadImageFromMemory(".png", tex.PNGBytes, int32(len(tex.PNGBytes)))
		rlTex := rl.LoadTextureFromImage(rlImage)

		p.Textures[tex.InternalIdInt] = rlTex
	}

	dir := "../../grpg-assets/maps/"
	entries, err := os.ReadDir(dir)

	if err != nil {
		log.Fatal("Error reading maps directory")
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fullPath := filepath.Join(dir, entry.Name())

			file, err1 := os.Open(fullPath)
			bytes, err2 := io.ReadAll(file)

			if err := cmp.Or(err1, err2); err != nil {
				log.Fatalf("Error reading map file, %v", err)
			}

			buf := gbuf.NewGBuf(bytes)
			header := grpgmap.ReadHeader(buf)

			if string(header.Magic[:]) != "GRPGMAP\x00" {
				log.Fatalf("File %s isn't GRPGMAP", fullPath)
			}

			tiles := grpgmap.ReadTiles(buf)
			chunkPos := util.Vector2I{X: int32(header.ChunkX), Y: int32(header.ChunkY)}

			p.Maps[chunkPos] = tiles
		}
	}
}

func (p *Playground) Cleanup() {
	rl.UnloadFont(p.Font)
}

func (p *Playground) Loop() {
	player := p.Game.Player

	if rl.IsKeyPressed(rl.KeyW) {
		player.SendMovePacket(p.Game, player.X, player.Y-1)
	} else if rl.IsKeyPressed(rl.KeyS) {
		player.SendMovePacket(p.Game, player.X, player.Y+1)
	} else if rl.IsKeyPressed(rl.KeyA) {
		player.SendMovePacket(p.Game, player.X-1, player.Y)
	} else if rl.IsKeyPressed(rl.KeyD) {
		player.SendMovePacket(p.Game, player.X+1, player.Y)
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
	mapTiles := p.Maps[util.Vector2I{X: p.Game.Player.ChunkX, Y: p.Game.Player.ChunkY}]

	for i := range 256 {
		dx := int32(i%16) * p.Game.TileSize
		dy := int32(i/16) * p.Game.TileSize

		tex := p.Textures[mapTiles[i].InternalId]
		rl.DrawTexture(tex, dx, dy, rl.White)
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
			rl.Vector2{X: float32(player.RealX), Y: float32(player.RealY)},
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
