package shared

import "github.com/hajimehoshi/ebiten/v2"

type GScene interface {
	Setup()
	Cleanup()
	Update() error
	Draw(screen *ebiten.Image)
}

type GSceneManager struct {
	CurrentScene GScene
}

func (gsm *GSceneManager) SwitchTo(other GScene) {
	if gsm.CurrentScene != nil {
		gsm.CurrentScene.Cleanup()
	}

	other.Setup()
	gsm.CurrentScene = other
}
