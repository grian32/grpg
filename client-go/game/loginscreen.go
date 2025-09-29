package game

import (
	"client/network/c2s"
	"client/shared"
	"client/util"
	"cmp"
	"image/color"
	"log"

	gebitenui "github.com/grian32/gebiten-ui"

	"github.com/hajimehoshi/ebiten/v2"
)

type LoginScreen struct {
	LoginButton     *gebitenui.GButton
	UsernameTextbox *gebitenui.GTextbox
	Font48          *gebitenui.GFont
	Font24          *gebitenui.GFont
	HalfWidth       float64
	HalfHeight      float64
	GRPGTextX       float64
	EnterNameTextX  float64
	FailedTextX     float64
	LoginName       string
	Game            *shared.Game
}

func (l *LoginScreen) Cleanup() {
}

func (l *LoginScreen) Setup() {
	var assetsDirectory = "../../grpg-assets/"
	fontBig, err1 := gebitenui.NewGFont(assetsDirectory+"font.ttf", 48)
	fontSmall, err2 := gebitenui.NewGFont(assetsDirectory+"font.ttf", 24)
	// would be same error eitherway in this case lol, same font, diff size
	if err := cmp.Or(err1, err2); err != nil {
		log.Fatalf("failed creating font: %v\n", err)
	}

	l.Font48 = fontBig
	l.Font24 = fontSmall

	textures := loadTex(assetsDirectory + "assets/login.grpgtex")

	halfScreenWidth := float64(l.Game.ScreenWidth / 2)
	halfScreenHeight := float64(l.Game.ScreenHeight / 2)

	l.HalfWidth = halfScreenWidth
	l.HalfHeight = halfScreenHeight
	width, _ := fontBig.MeasureString("GRPG")
	l.GRPGTextX = halfScreenWidth - (width / 2.0)
	nameWidth, _ := fontSmall.MeasureString("Enter Name Below")
	l.EnterNameTextX = halfScreenWidth - (nameWidth / 2.0)
	loginFailedWidth, _ := fontSmall.MeasureString("Login failed. Name already taken.")
	l.FailedTextX = halfScreenWidth - (loginFailedWidth / 2.0)

	btnTex := textures["login_button"]

	btn, err := gebitenui.NewButton(
		"Login",
		halfScreenWidth-float64(btnTex.Bounds().Dx()/2.0),
		halfScreenHeight-125,
		btnTex,
		fontSmall,
		func() {
			if l.LoginName != "" {
				l.Game.Player.Name = l.LoginName
				shared.SendPacket(l.Game.Conn, &c2s.LoginPacket{
					PlayerName: l.LoginName,
				})
			}
		},
	)
	if err != nil {
		log.Fatalf("failed to intialize login button: %v\n", err)
	}
	l.LoginButton = btn

	textboxTex := textures["login_name_textbox"]

	l.UsernameTextbox = gebitenui.NewTextBox(
		halfScreenWidth-float64(textboxTex.Bounds().Dx()/2.0),
		halfScreenHeight-250,
		8,
		textboxTex,
		fontSmall,
		24,
		0,
		func(newText string) {
			l.LoginName = newText
		},
	)
}

func (l *LoginScreen) Update() error {
	l.LoginButton.Update()
	l.UsernameTextbox.Update()
	return nil
}

func (l *LoginScreen) Draw(screen *ebiten.Image) {
	bgColor := util.ValuesRGB(17, 33, 43)

	screen.Fill(bgColor)
	l.LoginButton.Draw(screen)
	l.UsernameTextbox.Draw(screen)

	l.Font48.Draw(screen, "GRPG", l.GRPGTextX, l.HalfHeight-375, color.White)
	l.Font24.Draw(screen, "Enter Name Below", l.EnterNameTextX, l.HalfHeight-275, color.White)
	if l.Game.ShowFailedLogin {
		l.Font24.Draw(screen, "Login failed. Name already taken.", l.FailedTextX, float64(l.Game.ScreenHeight)*0.95, util.ValuesRGB(255, 0, 0))
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
