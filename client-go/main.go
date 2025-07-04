package main

import (
	"client/game"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	g = &game.Game{
		ScreenWidth:  960,
		ScreenHeight: 960,
		TileSize:     64,
		SceneManager: &game.GSceneManager{},
		Player:       &game.Player{X: 15, Y: 15, RealX: 960, RealY: 960, Name: ""},
	}
)

func main() {
	rl.InitWindow(960, 960, "GRPG Client")

	g.SceneManager.SwitchTo(&game.LoginScreen{
		Game: g,
	})

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		g.SceneManager.CurrentScene.Loop()
		rl.BeginDrawing()

		g.SceneManager.CurrentScene.Render()

		rl.EndDrawing()
	}
}
