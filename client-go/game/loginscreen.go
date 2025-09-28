package game

import (
	"client/shared"
	"client/util"
	"cmp"
	"fmt"
	"image/color"
	"log"

	gebitenui "github.com/grian32/gebiten-ui"

	"github.com/hajimehoshi/ebiten/v2"
)

type LoginScreen struct {
	LoginButton     *gebitenui.GButton
	UsernameTextbox *gebitenui.GTextbox
	LoginName       string
	Game            *shared.Game
}

func (l *LoginScreen) Cleanup() {
	//if l.Font.Texture.ID != 0 {
	//	rl.UnloadFont(l.Font)
	//}
}

func (l *LoginScreen) Setup() {
	var assetsDirectory = "../../grpg-assets/"
	fontBig, err1 := gebitenui.NewGFont(assetsDirectory+"font.ttf", 48)
	fontSmall, err2 := gebitenui.NewGFont(assetsDirectory+"font.ttf", 24)
	// would be same error eitherway in this case lol
	if err := cmp.Or(err1, err2); err != nil {
		log.Fatalf("failed creating font: %v\n", err)
	}

	_ = fontBig
	textures := loadTex(assetsDirectory + "assets/login.grpgtex")

	btn, err := gebitenui.NewButton(
		"Login",
		float64(l.Game.ScreenWidth/2)-60,
		float64(l.Game.ScreenHeight/2)-200,
		textures["login_button"],
		fontSmall,
		func() {
			fmt.Println("logging in!")
		},
	)
	if err != nil {
		log.Fatalf("failed to intialize login button: %v\n\n", err)
	}
	l.LoginButton = btn

}

func (l *LoginScreen) Update() error {
	l.LoginButton.Update()
	return nil
}

func (l *LoginScreen) Draw(screen *ebiten.Image) {
	screen.Fill(util.ValuesRGB(17, 33, 43))
	l.LoginButton.Draw(screen)
	//rl.ClearBackground(rl.Black)

	halfWidth := l.Game.ScreenWidth / 2
	halfHeight := l.Game.ScreenHeight / 2

	//drawCenteredText(l, halfWidth, float32(halfHeight-375), 48.0, 0.0, "GRPG Client", rl.White)
	//rl.DrawRectangle(halfWidth-200, halfHeight-300, 400, 200, rl.NewColor(186, 109, 22, 255))
	//drawCenteredText(l, halfWidth, float32(halfHeight-275), 24.0, 0.4, "Enter Name Below:", rl.White)
	drawLayout(l, halfWidth, halfHeight)

	if l.Game.ShowFailedLogin {
		//drawCenteredText(l, halfWidth, 900.0, 24.0, 0.0, "Login failed, name most likely already taken", rl.Red)
	}
}

func drawCenteredText(l *LoginScreen, halfWidth int32, yPos, size, spacing float32, text string, color color.RGBA) {
	//calculatedSize := rl.MeasureTextEx(l.Font, text, size, spacing)

	//textPos := rl.Vector2{
	//	X: float32(halfWidth) - (calculatedSize.X / 2),
	//	Y: yPos,
	//}
	//
	//rl.DrawTextEx(l.Font, text, textPos, size, spacing, color)
}

func drawLayout(l *LoginScreen, halfWidth, halfHeight int32) {
	//loginTextPos := rl.Rectangle{
	//	X:      float32(halfWidth) - 50,
	//	Y:      float32(halfHeight) - 250,
	//	Width:  100,
	//	Height: 30,
	//}
	//
	//rg.TextBox(
	//	loginTextPos,
	//	&l.LoginName,
	//	// one of the all time stupidest variable names/variables, it takes the input - 1 as the max chars in the
	//	// textbox, 317 pi level shit :(
	//	9,
	//	true,
	//)
	//
	//loginButtonPos := rl.Rectangle{
	//	X:      float32(halfWidth) - 30,
	//	Y:      float32(halfHeight) - 200,
	//	Width:  60,
	//	Height: 30,
	//}
	//
	//if rg.Button(loginButtonPos, "Login") {
	//	l.Game.Player.Name = l.LoginName
	//	shared.SendPacket(l.Game.Conn, &c2s.LoginPacket{
	//		PlayerName: l.LoginName,
	//	})
	//}
}
