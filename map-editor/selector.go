package main

import (
	"grpg/data-go/grpgobj"
	"grpg/data-go/grpgtile"

	g "github.com/AllenDang/giu"
)

type PlaceTypeSelected byte

const (
	UNDEFINED PlaceTypeSelected = iota
	TILE
	OBJ
)

var (
	currentlySelected      string            = "_undefined"
	currentlySelectedTexId int32             = -1
	typeSelected           PlaceTypeSelected = UNDEFINED
	textureSelected        bool              = false
)

func BuildTileSelectorTab(data []grpgtile.Tile) g.Widget {
	col1 := make([]g.Widget, 0)
	col2 := make([]g.Widget, 0)

	for i := range len(data) {
		if i%2 == 0 {
			col1 = append(col1, buildTileColElem(data[i]))
		} else {
			col2 = append(col2, buildTileColElem(data[i]))
		}
	}

	return g.Column(g.Row(g.Column(col1...), g.Column(col2...)), buildCurrentlySelected())
}

func buildTileColElem(d grpgtile.Tile) g.Widget {
	return g.Column(
		g.ImageButton(textures[int32(d.TexId)].Texture).OnClick(func() {
			currentlySelected = d.Name
			currentlySelectedTexId = int32(d.TexId)
			typeSelected = TILE
			textureSelected = true
			eraserEnabled = false
			g.Update()
		}),
		g.Label(d.Name),
	)
}

func BuildObjSelectorTabs(data []grpgobj.Obj) g.Widget {
	col1 := make([]g.Widget, 0)
	col2 := make([]g.Widget, 0)

	for i := range len(data) {
		if i%2 == 0 {
			col1 = append(col1, buildObjColElem(data[i]))
		} else {
			col2 = append(col2, buildObjColElem(data[i]))
		}
	}

	return g.Column(g.Row(g.Column(col1...), g.Column(col2...)), buildCurrentlySelected())
}

func buildObjColElem(d grpgobj.Obj) g.Widget {
	name := d.Name + getFlagsName(d.Flags)
	return g.Column(
		g.ImageButton(textures[int32(d.Textures[0])].Texture).OnClick(func() {
			currentlySelected = name
			currentlySelectedTexId = int32(d.Textures[0])
			typeSelected = OBJ
			textureSelected = true
			eraserEnabled = false
			g.Update()
		}),
		g.Label(name),
	)
}

func buildCurrentlySelected() g.Widget {
	if eraserEnabled {
		return g.Column(
			g.Label("Currently Selected: "),
			g.Label("Eraser!"),
		)
	} else if currentlySelected != "_undefined" {
		return g.Column(
			g.Label("Currently Selected: "),
			g.Image(textures[currentlySelectedTexId].Texture),
			g.Label(currentlySelected),
		)
	} else {
		return g.Column(
			g.Label("Currently Selected: "),
			g.Label("None :("),
		)
	}
}

func getFlagsName(flags grpgobj.ObjFlags) string {
	str := "("

	if grpgobj.IsFlagSet(flags, grpgobj.STATE) {
		str += "s"
	}

	if grpgobj.IsFlagSet(flags, grpgobj.INTERACT) {
		str += "|i"
	}

	return str + ")"
}
