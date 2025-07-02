package main

import "github.com/sqweek/dialog"

var (
	eraserEnabled = false
)

func SetAllEmptyTiles() {
	if currentlySelected.key == "_undefined" || len(textures) == 0 {
		dialog.Message("No tile currently selected/textures not loaded.").Error()
	}

	for idx, _ := range gridTextures {
		if gridTextures[idx] == "" {
			gridTextures[idx] = currentlySelected.val.InternalIdString
		}
	}
}

func ClearGrid() {
	ask := dialog.Message("Are you sure you want to completely wipe the grid?").YesNo()

	if !ask {
		return
	}

	for idx, _ := range gridTextures {
		gridTextures[idx] = ""
	}
}

func EnableEraser() {
	eraserEnabled = true
	currentlySelected = GTexKV{
		key: "",
		val: GiuTextureTyped{},
	}
}
