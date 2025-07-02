package main

import (
	g "github.com/AllenDang/giu"
	"grpg/data-go/grpgtex"
	"sort"
)

type GTexKV struct {
	key string
	val GiuTextureTyped
}

var (
	tiles                    = make([]GTexKV, 0)
	objs                     = make([]GTexKV, 0)
	currentlySelected GTexKV = GTexKV{
		key: "_undefined",
		val: GiuTextureTyped{},
	}
	textureSelected = false
)

func BuildSelectorTypeMaps() {
	for name, tex := range textures {
		switch tex.TextureType {
		case grpgtex.TILE:
			tiles = append(tiles, GTexKV{
				key: name,
				val: tex,
			})
		case grpgtex.OBJ:
			objs = append(objs, GTexKV{
				key: name,
				val: tex,
			})
		}
	}

	// sort to provide some consistency since maps are unordered.
	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i].key < tiles[j].key
	})
	sort.Slice(objs, func(i, j int) bool {
		return objs[i].key < objs[j].key
	})
}

func BuildSelectorTab(data []GTexKV) g.Widget {
	col1 := make([]g.Widget, 0)
	col2 := make([]g.Widget, 0)

	for i := range len(data) {
		// check even in case data is of uneven length
		if i%2 == 0 {
			col1 = append(col1, buildTextureColumn(data[i]))
		} else {
			col2 = append(col2, buildTextureColumn(data[i]))
		}
	}

	return g.Column(g.Row(g.Column(col1...), g.Column(col2...)), buildCurrentlySelected())
}

func buildTextureColumn(kv GTexKV) g.Widget {
	return g.Column(
		g.ImageButton(kv.val.Texture).OnClick(func() {
			currentlySelected = kv
			textureSelected = true
			g.Update()
		}),
		g.Label(kv.val.FormattedStringId),
	)
}

func buildCurrentlySelected() g.Widget {
	if currentlySelected.key != "_undefined" { // kind of ugly but can't set structs to nil :(
		return g.Column(
			g.Label("Currently Selected: "),
			g.Image(currentlySelected.val.Texture),
			g.Label(currentlySelected.val.FormattedStringId),
		)
	} else if eraserEnabled {
		return g.Column(
			g.Label("Currently Selected: "),
			g.Label("Eraser!"),
		)
	} else {
		return g.Column(
			g.Label("Currently Selected: "),
			g.Label("None :("),
		)
	}
}
