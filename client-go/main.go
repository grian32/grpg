package main

import (
	"client/game"
	"client/network"
	"client/shared"
	"client/util"
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	g = &shared.Game{
		ScreenWidth:  1152,
		ScreenHeight: 960,
		TileSize:     64,
		MaxX:         0,
		MaxY:         0,
		CollisionMap: make(map[util.Vector2I]struct{}),
		ObjIdByLoc:   make(map[util.Vector2I]uint16),
		TrackedObjs:  make(map[util.Vector2I]*shared.GameObj),
		TrackedNpcs:  make(map[util.Vector2I]*shared.GameNpc),
		SceneManager: &shared.GSceneManager{},
		Player: &shared.LocalPlayer{
			X:         0,
			Y:         0,
			RealX:     0,
			RealY:     0,
			Facing:    shared.UP,
			Inventory: [24]shared.InventoryItem{},
			Name:      "",
		},
		Talkbox:         shared.Talkbox{},
		OtherPlayers:    map[string]*shared.RemotePlayer{},
		Conn:            network.StartConn(),
		ShowFailedLogin: false,
	}
)

type GameWrapper struct {
	gsm     *shared.GSceneManager
	packets chan network.ChanPacket
	game    *shared.Game
}

func (g *GameWrapper) Update() error {
	processPackets(g.packets, g.game)
	return g.gsm.CurrentScene.Update()
}

func (g *GameWrapper) Draw(screen *ebiten.Image) {
	g.gsm.CurrentScene.Draw(screen)
}

func (g *GameWrapper) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.game.ScreenWidth), int(g.game.ScreenHeight)
}

func main() {
	windowTitle := flag.String("title", "GRPG Client", "the window title")

	flag.Parse()

	//rl.InitWindow(g.ScreenWidth, g.ScreenHeight, *windowTitle)
	ebiten.SetWindowSize(int(g.ScreenWidth), int(g.ScreenHeight))
	ebiten.SetWindowTitle(*windowTitle)

	g.SceneManager.SwitchTo(&game.LoginScreen{
		Game: g,
	})

	//defer rl.CloseWindow()
	defer g.Conn.Close()

	serverPackets := make(chan network.ChanPacket, 100)

	go network.ReadServerPackets(g.Conn, serverPackets)

	//rl.SetTargetFPS(60)

	ebGame := &GameWrapper{
		gsm:     g.SceneManager,
		packets: serverPackets,
		game:    g,
	}

	if err := ebiten.RunGame(ebGame); err != nil {
		log.Fatal(err)
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
