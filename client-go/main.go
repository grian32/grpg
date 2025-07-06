package main

import (
	"client/game"
	"client/network"
	"client/shared"
	"flag"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	g = &shared.Game{
		ScreenWidth:  960,
		ScreenHeight: 960,
		TileSize:     64,
		SceneManager: &shared.GSceneManager{},
		Player:       &shared.Player{X: 15, Y: 15, RealX: 960, RealY: 960, Name: ""},
		OtherPlayers: []shared.Player{
			{
				X:      12,
				Y:      12,
				RealX:  768,
				RealY:  768,
				ChunkX: 0,
				ChunkY: 0,
				Name:   "OtherDu",
			},
			{
				X:      12,
				Y:      11,
				RealX:  768,
				RealY:  704,
				ChunkX: 0,
				ChunkY: 0,
				Name:   "OtherD",
			},
		},
		Conn:            network.StartConn(),
		ShowFailedLogin: false,
	}
)

func main() {
	windowTitle := flag.String("title", "GRPG Client", "the window title")

	flag.Parse()

	rl.InitWindow(960, 960, *windowTitle)

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
