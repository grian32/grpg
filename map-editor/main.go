package main

import (
	"fmt"
	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
	"image/color"
)

func loop() {
	//imgui.ShowDemoWindow()
	var editorWindowPos imgui.Vec2
	var editorWindowSize imgui.Vec2
	g.Window("Editor").Flags(g.WindowFlagsNoCollapse|g.WindowFlagsNoMove|g.WindowFlagsNoResize|g.WindowFlagsAlwaysAutoResize).Layout(
		g.Custom(func() {
			editorWindowPos = imgui.WindowPos()
			editorWindowSize = imgui.WindowSize()
		}),
		g.Style().SetStyle(g.StyleVarItemSpacing, 0, 0).SetStyle(g.StyleVarItemInnerSpacing, 0, 0).SetStyleFloat(g.StyleVarFrameBorderSize, 0).SetStyle(g.StyleVarFramePadding, 1, 1).To(
			BuildGrid()...,
		),
	)

	g.Window("Buttons").Pos(editorWindowPos.X, editorWindowPos.Y+editorWindowSize.Y+10).Flags(g.WindowFlagsNoCollapse | g.WindowFlagsNoMove | g.WindowFlagsNoResize | g.WindowFlagsAlwaysAutoResize).Layout(
		g.Button("Load Textures").OnClick(func() {
			fmt.Println("loading tex")
		}),
	)
}

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
