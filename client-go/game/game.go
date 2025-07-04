package game

type Game struct {
	ScreenWidth  int32
	ScreenHeight int32
	TileSize     int32
	SceneManager *GSceneManager
	Player       *Player
}
