package main

import (
	"client/game"
	"client/network"
	"client/shared"
	"client/util"
	"flag"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	g = &shared.Game{
		ScreenWidth:     1152,
		ScreenHeight:    960,
		TileSize:        64,
		MaxX:            0,
		MaxY:            0,
		CollisionMap:    make(map[util.Vector2I]struct{}),
		TrackedObjs:     make(map[util.Vector2I]*shared.GameObj),
		SceneManager:    &shared.GSceneManager{},
		Player:          &shared.LocalPlayer{X: 0, Y: 0, RealX: 0, RealY: 0, Facing: shared.UP, Name: ""},
		OtherPlayers:    map[string]*shared.RemotePlayer{},
		Conn:            network.StartConn(),
		ShowFailedLogin: false,
	}
)

func main() {
	windowTitle := flag.String("title", "GRPG Client", "the window title")

	flag.Parse()

	rl.InitWindow(g.ScreenWidth, g.ScreenHeight, *windowTitle)

	g.SceneManager.SwitchTo(&game.LoginScreen{
		Game: g,
	})

	defer rl.CloseWindow()
	defer g.Conn.Close()

	serverPackets := make(chan network.ChanPacket, 100)

	go network.ReadServerPackets(g.Conn, serverPackets)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		processPackets(serverPackets, g)

		g.SceneManager.CurrentScene.Loop()
		rl.BeginDrawing()

		g.SceneManager.CurrentScene.Render()

		rl.EndDrawing()
	}

	// if i defer this it tries to double free for w/e reason, not sure why
	// it should be called first if i defer it last, but ¯\_(ツ)_/¯
	g.SceneManager.CurrentScene.Cleanup()
}

func processPackets(packetChan <-chan network.ChanPacket, g *shared.Game) {
	for {
		select {
		case packet := <-packetChan:
			packet.PacketData.Handler.Handle(packet.Buf, g)
		default:
			return
		}
	}
}
