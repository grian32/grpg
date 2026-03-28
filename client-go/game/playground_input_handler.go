package game

import (
	"client/network/c2s"
	"client/shared"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PgInputHandler struct {
	Player *shared.LocalPlayer
	Game   *shared.Game

	IsTypingCommand bool
	CommandString   string

	MoveFrameCounter int
	MovementHeld     bool

	minInvX, maxInvX, minInvY, maxInvY int
}

func NewPgInputHandler(g *shared.Game) *PgInputHandler {
	h := &PgInputHandler{Game: g, Player: g.Player}
	h.minInvX = RightGameframeX + TileSize
	h.maxInvX = h.minInvX + TileSize*4
	h.minInvY = TileSize
	h.maxInvY = h.minInvY + TileSize*6
	return h
}

func (h *PgInputHandler) Update() {
	if h.IsTypingCommand {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			h.IsTypingCommand = false
			h.CommandString = ""
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			h.IsTypingCommand = false
			h.Player.SendCmdPacket(h.Game, h.CommandString)
			h.CommandString = ""
		} else if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(h.CommandString) != 0 {
			h.CommandString = h.CommandString[:len(h.CommandString)-1]
		} else {
			h.CommandString = string(ebiten.AppendInputChars([]rune(h.CommandString)))
		}
		h.MovementHeld = false
		h.MoveFrameCounter = 0
		return
	}

	h.MovementHeld = ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD)

	if h.MovementHeld {
		if h.MoveFrameCounter%10 == 0 {
			if ebiten.IsKeyPressed(ebiten.KeyW) {
				h.Player.SendMovePacket(h.Game, h.Player.X, h.Player.Y-1, shared.UP)
			} else if ebiten.IsKeyPressed(ebiten.KeyS) {
				h.Player.SendMovePacket(h.Game, h.Player.X, h.Player.Y+1, shared.DOWN)
			} else if ebiten.IsKeyPressed(ebiten.KeyA) {
				h.Player.SendMovePacket(h.Game, h.Player.X-1, h.Player.Y, shared.LEFT)
			} else if ebiten.IsKeyPressed(ebiten.KeyD) {
				h.Player.SendMovePacket(h.Game, h.Player.X+1, h.Player.Y, shared.RIGHT)
			}
		}
		h.MoveFrameCounter++
	} else {
		h.MoveFrameCounter = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		h.Player.SendInteractPacket(h.Game)
	} else if h.Game.Talkbox.Active && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		shared.SendPacket(h.Game.Conn, &c2s.Continue{})
	} else if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		h.IsTypingCommand = true
	}
}

func (h *PgInputHandler) UpdateItemMove(renderType RenderType, outlineInvSpot *int) {
	if renderType == Inventory && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		if mouseX >= h.minInvX && mouseX < h.maxInvX && mouseY >= h.minInvY && mouseY < h.maxInvY {
			col := (mouseX - h.minInvX) / TileSize
			row := (mouseY - h.minInvY) / TileSize
			idx := row*ItemsPerRow + col

			if *outlineInvSpot != -1 {
				if *outlineInvSpot == idx || h.Player.Inventory[idx].ItemId != 0 {
					// deselect behaviour, basically
					*outlineInvSpot = -1
					return
				}

				shared.SendPacket(h.Game.Conn, &c2s.InvSwap{
					From: byte(*outlineInvSpot),
					To:   byte(idx),
				})
				*outlineInvSpot = -1

				return
			}

			if h.Player.Inventory[idx].ItemId != 0 {
				*outlineInvSpot = idx
			}
		}
	}
}
