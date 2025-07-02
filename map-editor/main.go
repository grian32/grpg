package main

import (
	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
	"image/color"
	"sync"
)

var (
	chunkX int32 = -1
	chunkY int32 = -1
)

func main() {
	wnd := g.NewMasterWindow("GRPG Map Editor", 1640, 1200, g.MasterWindowFlagsNotResizable)
	wnd.SetBgColor(color.RGBA{
		R: 17,
		G: 31,
		B: 86,
		A: 255,
	})
	wnd.Run(loop)
}

var syncOnce sync.Once

func loop() {
	syncOnce.Do(LoadDefaultGridTex)
	//imgui.ShowDemoWindow()
	var editorWindowPos imgui.Vec2
	var editorWindowSize imgui.Vec2
	g.Window("Editor").Flags(g.WindowFlagsNoCollapse|g.WindowFlagsNoMove|g.WindowFlagsNoResize).Size(1044, 1064).Layout(
		g.Custom(func() {
			editorWindowPos = imgui.WindowPos()
			editorWindowSize = imgui.WindowSize()
		}),
		BuildGrid(),
	)

	g.Window("Controls").Pos(editorWindowPos.X, editorWindowPos.Y+editorWindowSize.Y+10).Flags(g.WindowFlagsNoCollapse|g.WindowFlagsNoMove|g.WindowFlagsNoResize|g.WindowFlagsAlwaysAutoResize).Layout(
		g.Row(
			g.Button("Load Textures").OnClick(LoadTextures),
			g.Button("Save Map").OnClick(SaveMap),
			g.Button("Load Map").OnClick(LoadMap),
			g.Button("Set all empty tiles to currently selected").OnClick(SetAllEmptyTiles),
			g.Button("Clear Grid").OnClick(ClearGrid),
		),
		g.Row(
			g.Column(
				g.Label("Chunk X"),
				g.InputInt(&chunkX).Flags(g.InputTextFlagsCharsDecimal),
			),
			g.Column(
				g.Label("Chunk Y"),
				g.InputInt(&chunkY).Flags(g.InputTextFlagsCharsDecimal),
			),
		),
	)

	g.Window("Selector").Pos(editorWindowPos.X+editorWindowSize.X+10, editorWindowPos.Y).Flags(g.WindowFlagsNoCollapse | g.WindowFlagsNoMove | g.WindowFlagsNoResize | g.WindowFlagsAlwaysAutoResize).Layout(
		g.TabBar().TabItems(
			g.TabItem("Tiles").Layout(BuildSelectorTab(tiles)),
			g.TabItem("Objs").Layout(BuildSelectorTab(objs)),
		),
	)
}
