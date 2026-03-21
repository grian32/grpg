package game

import (
	"client/constants"
	"client/shared"
	"client/util"
	"fmt"
	"grpg/data-go/grpgmap"
	"image/color"

	gebiten_ui "github.com/grian32/gebiten-ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type PgWorld struct {
	Game *shared.Game

	Zones    map[util.Vector2I]grpgmap.Zone
	Textures map[uint16]*ebiten.Image
	Font     *gebiten_ui.GFont

	ExclamTexture *ebiten.Image
	WorldImage    *ebiten.Image

	ExclamYOffset int32
}

func NewPgWorld(
	game *shared.Game,
	worldImage *ebiten.Image,
	textures map[uint16]*ebiten.Image,
	exclamTex *ebiten.Image,
	font *gebiten_ui.GFont,
) *PgWorld {
	return &PgWorld{
		Game: game,

		Zones:    loadMaps(AssetsDir+"maps/", game),
		Textures: textures,
		Font:     font,

		ExclamTexture: exclamTex,
		WorldImage:    worldImage,
	}
}

func (p *PgWorld) Update(ticks uint32) {
	if (ticks/ExclamAnimTickInterval)%ExclamAnimFrameCount == 0 {
		p.ExclamYOffset = 0
	} else {
		p.ExclamYOffset = ExclamBobOffset
	}
}

func (p *PgWorld) Draw() {
	player := p.Game.Player

	mapTiles := p.Zones[util.Vector2I{X: p.Game.Player.ChunkX, Y: p.Game.Player.ChunkY}]

	for i := range TilesPerChunk {
		localX := int32(i % ChunkSize)
		localY := int32(i / ChunkSize)

		dx := localX * p.Game.TileSize
		dy := localY * p.Game.TileSize

		texId := p.Game.Tiles[uint16(mapTiles.Tiles[i])].TexId

		tex := p.Textures[texId]

		util.DrawImage(p.WorldImage, tex, dx, dy)

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
			util.DrawImage(p.WorldImage, objTex, dx, dy)
		}
		if trackedNpc, ok := p.Game.NpcsByPos[worldPos]; ok {
			// TODO: maybe don't render if player is standing over?
			npcTexId := trackedNpc.NpcData.TextureId
			util.DrawImage(p.WorldImage, p.Textures[npcTexId], dx, dy)
			if p.Game.DebugMode {
				p.Font.Draw(p.WorldImage, fmt.Sprintf("%d", trackedNpc.Uid), float64(dx), float64(dy), color.White)
			}

			if trackedNpc.NpcData.NpcId == uint16(constants.GRPG_GUIDE) && p.Game.RenderExclamOnGuide {
				util.DrawImage(p.WorldImage, p.ExclamTexture, dx, dy-TileSize+p.ExclamYOffset)
			}
		}
	}
}
