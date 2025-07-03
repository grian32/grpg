package main

import (
	"client/scene"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	g = &scene.Game{
		ScreenWidth:  960,
		ScreenHeight: 960,
		SceneManager: scene.GSceneManager{},
		PlayerName:   "",
	}
)

func main() {
	rl.InitWindow(960, 960, "GRPG Client")

	alkhemikalFont := rl.LoadFont("./assets/font.ttf")
	defer rl.UnloadFont(alkhemikalFont)

	g.SceneManager.SwitchTo(&scene.LoginScreen{
		Font: alkhemikalFont,
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
