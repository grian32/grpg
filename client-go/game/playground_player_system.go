package game

import (
	"client/shared"
	"client/util"
	"image"
	"image/color"

	gebiten_ui "github.com/grian32/gebiten-ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type PgPlayerSystem struct {
	WorldImage    *ebiten.Image
	PlayerIdleRun *ebiten.Image
	Game          *shared.Game
	InputHandler  *PgInputHandler
	Font16        *gebiten_ui.GFont
}

func NewPgPlayerSystem(
	worldImage *ebiten.Image,
	otherTex map[string]*ebiten.Image,
	game *shared.Game,
	font16 *gebiten_ui.GFont,
	inputHandler *PgInputHandler,
) *PgPlayerSystem {

	return &PgPlayerSystem{
		WorldImage:    worldImage,
		PlayerIdleRun: otherTex["player_idle_run"],
		Game:          game,
		InputHandler:  inputHandler,
		Font16:        font16,
	}
}

func (r *PgPlayerSystem) Update(crossedZone bool) {
	r.Game.Player.Update(r.Game, crossedZone, r.InputHandler.MovementHeld)

	for _, rp := range r.Game.OtherPlayers {
		rp.Update(r.Game)
	}
}

func (r *PgPlayerSystem) Draw() {
	lp := r.Game.Player
	r.drawPlayer(lp.CurrFrame, lp.Facing, lp.RealX, lp.RealY, lp.Name, color.White)

	for _, rp := range r.Game.OtherPlayers {
		r.drawPlayer(rp.CurrFrame, rp.Facing, rp.RealX, rp.RealY, rp.Name, util.Red)
	}
}

func (r *PgPlayerSystem) drawPlayer(currFrame uint8, facing shared.Direction, realX int32, realY int32, name string, textColor color.Color) {
	srcX := int(currFrame) * TileSize
	sourceRec := image.Rectangle{
		Min: image.Point{
			X: srcX,
			Y: int(facing) * TileSize,
		},
		Max: image.Point{
			X: srcX + TileSize,
			Y: int(facing)*TileSize + TileSize,
		},
	}
	sub := util.SubImage(r.PlayerIdleRun, sourceRec)
	util.DrawImage(r.WorldImage, sub, realX, realY)

	r.Font16.Draw(r.WorldImage, name, float64(realX), float64(realY), textColor)
}
