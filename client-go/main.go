package main

import (
	"client/scene"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	screenWidth  int32  = 960
	screenHeight int32  = 960
	playerName   string = ""
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "GRPG Client")

	alkhemikalFont := rl.LoadFont("./assets/font.ttf")
	defer rl.UnloadFont(alkhemikalFont)

	sceneManager := scene.NewGSceneManager(&scene.LoginScreen{
		Font:         alkhemikalFont,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	})

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		sceneManager.CurrentScene.Loop()
		rl.BeginDrawing()

		sceneManager.CurrentScene.Render()

		rl.EndDrawing()
	}
}
