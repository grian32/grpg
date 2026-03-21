package game

import (
	"client/shared"
	"client/util"
)

type PgCamera struct {
	Player *shared.LocalPlayer
	CameraTarget util.Vector2
	PrevCameraTarget util.Vector2
}

func NewPgCamera(p *shared.LocalPlayer) *PgCamera {
	return &PgCamera{
		Player: p,
	}
}

func (c *PgCamera) Update(crossedZone bool) {
	var cameraX int32 = CameraOffsetTiles * TileSize
	var cameraY int32 = CameraOffsetTiles * TileSize

	if c.Player.RealX <= CameraBoundaryTiles*TileSize {
		cameraX = util.MinI(c.Player.RealX-(CameraMinOffsetTiles*TileSize), 0)
	}

	if c.Player.RealY <= CameraBoundaryTiles*TileSize {
		cameraY = util.MinI(c.Player.RealY-(CameraMinOffsetTiles*TileSize), 0)
	}

	if crossedZone {
		c.CameraTarget.X = float64(cameraX)
		c.CameraTarget.Y = float64(cameraY)
	} else {
		if c.CameraTarget.X < float64(cameraX) {
			c.CameraTarget.X += CameraPanSpeed
		} else if c.CameraTarget.X > float64(cameraX) {
			c.CameraTarget.X -= CameraPanSpeed
		}

		if c.CameraTarget.Y < float64(cameraY) {
			c.CameraTarget.Y += CameraPanSpeed
		} else if c.CameraTarget.Y > float64(cameraY) {
			c.CameraTarget.Y -= CameraPanSpeed
		}
	}

	c.PrevCameraTarget = c.CameraTarget
}
