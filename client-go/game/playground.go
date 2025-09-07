package game

import (
	"client/network/c2s"
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
	GameframeBottom  rl.Texture2D
	PlayerTextures   map[shared.Direction]rl.Texture2D
	Textures         map[uint16]rl.Texture2D
	Zones            map[util.Vector2I]grpgmap.Zone
	CameraTarget     rl.Vector2
	PrevCameraTarget rl.Vector2
	CurrActionString string
}

// this is reused in loginscreen.go can't rly avoid however
var assetsDirectory = "../../grpg-assets/"

func (p *Playground) Setup() {
	p.Game.Talkbox.CurrentName = "hi!"
	p.Font = rl.LoadFont(assetsDirectory + "font.ttf")
	p.CurrActionString = "Current Action: None :("

	p.Textures = loadTextures(assetsDirectory + "textures.grpgtex")
	p.Game.Objs = loadObjs(assetsDirectory + "objs.grpgobj")
	p.Game.Tiles = loadTiles(assetsDirectory + "tiles.grpgtile")
	p.Game.Items = loadItems(assetsDirectory + "items.grpgitem")
	p.Game.Npcs = loadNpcs(assetsDirectory + "npcs.grpgnpc")
	p.Zones = loadMaps(assetsDirectory+"maps/", p.Game)

	otherTex := loadOtherTex(assetsDirectory + "other.grpgtex")

	p.GameframeRight = otherTex["gameframe_right"]
	p.GameframeBottom = otherTex["gameframe_bottom"]

	p.PlayerTextures = make(map[shared.Direction]rl.Texture2D)
	p.PlayerTextures[shared.UP] = otherTex["player_up"]
	p.PlayerTextures[shared.DOWN] = otherTex["player_down"]
	p.PlayerTextures[shared.LEFT] = otherTex["player_left"]
	p.PlayerTextures[shared.RIGHT] = otherTex["player_right"]
}

func (p *Playground) Cleanup() {
	if p.Font.Texture.ID != 0 {
		rl.UnloadFont(p.Font)
	}

	for _, tex := range p.Textures {
		rl.UnloadTexture(tex)
	}

	rl.UnloadTexture(p.GameframeRight)
	rl.UnloadTexture(p.GameframeBottom)

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
	} else if p.Game.Talkbox.Active && rl.IsKeyPressed(rl.KeySpace) {
		shared.SendPacket(p.Game.Conn, &c2s.Continue{})
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
	trackedObj, objExists := p.Game.TrackedObjs[facingCoord]
	trackedNpc, npcExists := p.Game.TrackedNpcs[facingCoord]
	if objExists {
		p.CurrActionString = trackedObj.DataObj.InteractText
	} else if npcExists {
		p.CurrActionString = "Talk to " + trackedNpc.NpcData.Name
	} else {
		p.CurrActionString = "None :("
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
			rl.DrawTexture(objTex, dx, dy, rl.White)
		} else {
			// TODO: maybe don't render if player is standing over?
			trackedNpc, ok := p.Game.TrackedNpcs[worldPos]
			if ok {
				npcTexId := trackedNpc.NpcData.TextureId
				rl.DrawTexture(p.Textures[npcTexId], dx, dy, rl.White)
			}
		}
	}
}

func drawPlayer(p *Playground) {
	player := p.Game.Player

	if player.Facing == shared.LEFT {
		sourceRec := rl.Rectangle{
			X:      float32(player.CurrFrame * 64),
			Y:      0,
			Width:  64,
			Height: 64,
		}
		rl.DrawTextureRec(p.PlayerTextures[player.Facing], sourceRec, rl.Vector2{X: float32(player.RealX), Y: float32(player.RealY)}, rl.White)
	} else {
		rl.DrawTexture(p.PlayerTextures[player.Facing], player.RealX, player.RealY, rl.White)
	}

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
	// TODO: render based on naim sprite sheet.
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

	rl.DrawTexture(p.GameframeBottom, 0, 768, rl.White)

	talkbox := p.Game.Talkbox
	// x is offset from 0, y has offset added, to be placed in the right spot
	rl.DrawTextEx(p.Font, "Current Action: "+p.CurrActionString, rl.Vector2{X: 110, Y: 768 + 28 + 3}, 20, 0, rl.White)
	if talkbox.Active {
		rl.DrawTextEx(p.Font, talkbox.CurrentName, rl.Vector2{X: 110 + 332, Y: 768 + 28 + 3}, 24, 0, rl.White)
		rl.DrawTextEx(p.Font, talkbox.CurrentMessage, rl.Vector2{X: 90, Y: 840}, 24, 0, rl.White)
	}

	playerCoords := fmt.Sprintf("X: %d, Y: %d, Facing: %s", player.X, player.Y, player.Facing.String())
	rl.DrawTextEx(p.Font, playerCoords, rl.Vector2{X: 768, Y: 800}, 24, 0, rl.White)
}
