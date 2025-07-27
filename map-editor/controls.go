package main

import "github.com/sqweek/dialog"

var (
	eraserEnabled = false
)

func SetAllEmptyTiles() {
	if currentlySelected.key == "_undefined" || len(textures) == 0 {
		dialog.Message("No tile currently selected/textures not loaded.").Error()
	}

	for idx, _ := range gridTileTextures {
		if gridTileTextures[idx] == -1 {
			gridTileTextures[idx] = int32(currentlySelected.val.InternalId)
		}
	}
}

func ClearGrid() {
	ask := dialog.Message("Are you sure you want to completely wipe the grid?").YesNo()

	if !ask {
		return
	}

	for idx, _ := range gridTileTextures {
		gridTileTextures[idx] = -1
		gridObjTextures[idx] = -1
	}
}

func EnableEraser() {
	textureSelected = true
	eraserEnabled = true
	currentlySelected = GTexKV{
		key: "",
		val: GiuTextureTyped{},
	}
}
