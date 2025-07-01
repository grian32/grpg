package main

import (
	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
	"image/color"
)

func main() {
	wnd := g.NewMasterWindow("GRPG Map Editor", 1640, 1240, g.MasterWindowFlagsNotResizable)
	wnd.SetBgColor(color.RGBA{
		R: 17,
		G: 31,
		B: 86,
		A: 255,
	})
	wnd.Run(loop)
}

func loop() {
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

	g.Window("Buttons").Pos(editorWindowPos.X, editorWindowPos.Y+editorWindowSize.Y+10).Flags(g.WindowFlagsNoCollapse | g.WindowFlagsNoMove | g.WindowFlagsNoResize | g.WindowFlagsAlwaysAutoResize).Layout(
		g.Button("Load Textures").OnClick(func() {
			LoadTextures()
		}),
	)

	g.Window("Selector").Pos(editorWindowPos.X+editorWindowSize.X+10, editorWindowPos.Y).Flags(g.WindowFlagsNoCollapse | g.WindowFlagsNoMove | g.WindowFlagsNoResize | g.WindowFlagsAlwaysAutoResize).Layout(
		g.TabBar().TabItems(
			g.TabItem("Tiles").Layout(BuildSelectorTab(tiles)),
			g.TabItem("Objs").Layout(BuildSelectorTab(objs)),
		),
	)
}
