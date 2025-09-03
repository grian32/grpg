package main

import (
	"github.com/sqweek/dialog"
)

var (
	eraserEnabled = false
)

func SetAllEmptyTiles() {
	if currentlySelected == "_undefined" || len(textures) == 0 || !assetsLoaded || typeSelected != TILE {
		dialog.Message("No tile currently selected/assets not loaded.").Error()
	}

	for idx, _ := range gridTileTextures {
		if gridTileTextures[idx] == -1 {
			gridTileTextures[idx] = currentlySelectedTexId
			gridTiles[idx] = currentlySelectedInternalId
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
		gridTiles[idx] = -1
		gridObjTextures[idx] = -1
		gridObjs[idx] = -1
	}
}

func EnableEraser() {
	textureSelected = true
	eraserEnabled = true
	typeSelected = UNDEFINED
	currentlySelected = ""
	currentlySelectedTexId = -1
}
