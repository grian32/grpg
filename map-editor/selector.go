package main

import (
	"grpg/data-go/grpgobj"
	"grpg/data-go/grpgtile"
	"sort"

	g "github.com/AllenDang/giu"
)

type PlaceTypeSelected byte

const (
	UNDEFINED PlaceTypeSelected = iota
	TILE
	OBJ
)

var (
	currentlySelected           string            = "_undefined"
	currentlySelectedTexId      int32             = -1
	currentlySelectedInternalId int32             = -1
	typeSelected                PlaceTypeSelected = UNDEFINED
	textureSelected             bool              = false
)

func BuildTileSelectorTab(data map[int32]grpgtile.Tile) g.Widget {
	col1 := make([]g.Widget, 0)
	col2 := make([]g.Widget, 0)

	count := 0

	keys := make([]int, 0, len(data))

	for k, _ := range data {
		keys = append(keys, int(k))
	}

	sort.Ints(keys)

	for _, k := range keys {
		if count%2 == 0 {
			col1 = append(col1, buildTileColElem(data[int32(k)]))
		} else {
			col2 = append(col2, buildTileColElem(data[int32(k)]))
		}
		count++
	}
	return g.Column(g.Row(g.Column(col1...), g.Column(col2...)), buildCurrentlySelected())
}

func buildTileColElem(d grpgtile.Tile) g.Widget {
	return g.Column(
		g.ImageButton(textures[int32(d.TexId)].Texture).OnClick(func() {
			currentlySelected = d.Name
			currentlySelectedInternalId = int32(d.TileId)
			currentlySelectedTexId = int32(d.TexId)
			typeSelected = TILE
			textureSelected = true
			eraserEnabled = false
			g.Update()
		}),
		g.Label(d.Name),
	)
}

func BuildObjSelectorTabs(data map[int32]grpgobj.Obj) g.Widget {
	col1 := make([]g.Widget, 0)
	col2 := make([]g.Widget, 0)

	count := 0

	keys := make([]int, 0, len(data))

	for k, _ := range data {
		keys = append(keys, int(k))
	}

	sort.Ints(keys)

	for _, k := range keys {
		if count%2 == 0 {
			col1 = append(col1, buildObjColElem(data[int32(k)]))
		} else {
			col2 = append(col2, buildObjColElem(data[int32(k)]))
		}
		count++
	}

	return g.Column(g.Row(g.Column(col1...), g.Column(col2...)), buildCurrentlySelected())
}

func buildObjColElem(d grpgobj.Obj) g.Widget {
	name := d.Name + getFlagsName(d.Flags)
	return g.Column(
		g.ImageButton(textures[int32(d.Textures[0])].Texture).OnClick(func() {
			currentlySelected = name
			currentlySelectedInternalId = int32(d.ObjId)
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
