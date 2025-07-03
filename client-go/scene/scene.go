package scene

type GScene interface {
	Setup()
	Cleanup()
	Loop()
	Render()
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
