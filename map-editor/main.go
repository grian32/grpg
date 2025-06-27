package main

import (
	g "github.com/AllenDang/giu"
	"image/color"
)

func loop() {
	//imgui.ShowDemoWindow()
	g.Window("Editor").Flags(g.WindowFlagsNoCollapse | g.WindowFlagsNoMove | g.WindowFlagsNoResize | g.WindowFlagsAlwaysAutoResize).Layout(
		g.Style().SetStyle(g.StyleVarItemSpacing, 0, 0).SetStyle(g.StyleVarItemInnerSpacing, 0, 0).SetStyleFloat(g.StyleVarFrameBorderSize, 0).SetStyle(g.StyleVarFramePadding, 1, 1).To(
			BuildGrid()...,
		),
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
