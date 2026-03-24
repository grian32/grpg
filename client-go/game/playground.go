package game

import (
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
	Gameframe      *PgGameframe

	WorldImage *ebiten.Image

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

	p.Camera = NewPgCamera(p.Game.Player)
	p.World = NewPgWorld(p.Game, p.WorldImage, p.Textures, otherTex["exclam"], p.Font16)
	p.InputHandler = NewPgInputHandler(p.Game)
	p.PlayerRenderer = NewPgPlayerRenderer(p.WorldImage, otherTex, p.Game, p.Font16)
	p.Gameframe = NewPgGameframe(
		p.Game,
		p.InputHandler,
		p.Font16,
		p.Font20,
		p.Font24,
		p.Textures,
		otherTex,
	)
}

func (p *Playground) Cleanup() {
	// would usually dispose here but gc will take care of it since this is the last scene they're on
	// NOTE: will require disposal if i start switching scenes back to login or something else
}

func (p *Playground) Update() error {
	player := p.Game.Player
	p.InputHandler.Update()

	crossedZone := player.PrevX/ChunkSize != player.ChunkX || player.PrevY/ChunkSize != player.ChunkY

	// TODO: move this to player renderer and rename it to pg_players or something
	//pass crossed zone here as im already computing it for camera
	player.Update(p.Game, crossedZone)

	for _, rp := range p.Game.OtherPlayers {
		rp.Update(p.Game)
	}

	p.Camera.Update(crossedZone)
	p.World.Update(p.Ticks)
	p.Gameframe.Update()

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

	p.Gameframe.Draw(screen)
}
