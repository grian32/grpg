package scene

import rl "github.com/gen2brain/raylib-go/raylib"

type Playground struct {
	Game *Game
}

func (p Playground) Setup() {
}

func (p Playground) Cleanup() {
}

func (p Playground) Loop() {
}

func (p Playground) Render() {
	rl.ClearBackground(rl.Black)
}
