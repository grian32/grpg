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

func NewGSceneManager(initialScene GScene) *GSceneManager {
	initialScene.Setup()
	return &GSceneManager{
		CurrentScene: initialScene,
	}
}

func (gsm *GSceneManager) SwitchTo(other GScene) {
	gsm.CurrentScene.Cleanup()
	other.Setup()
	gsm.CurrentScene = other
}
