package main

import (
	"github.com/sqweek/dialog"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgmap"
	"os"
)

func SaveMap() {
	tileArr := [256]grpgmap.Tile{}

	if chunkX == -1 || chunkY == -1 {
		dialog.Message("Both Chunk X & Chunk Y must be set to save a map.").Error()
		return
	}

	for idx, name := range gridTextures {
		if name == "" {
			dialog.Message("All tiles must be filled in to save a map.").Error()
			return
		}
		tileArr[idx] = grpgmap.Tile{
			InternalId: textures[name].InternalId,
			Type:       textures[name].TextureType,
		}
	}

	fileToSave, err := dialog.File().Title("Please make a file to save to. Warning, this may wipe the file if it already exists.").Save()
	if err != nil {
		dialog.Message("Error finding directory to save to.").Error()
		return
	}

	_, err = os.Stat(fileToSave)
	if err == nil {
		dialog.Message("File already exists.").Error()
		return
	}

	file, err := os.Create(fileToSave)
	if err != nil {
		dialog.Message("Error creating file to save to").Error()
		return
	}
	defer file.Close()

	buf := gbuf.NewEmptyGBuf()
	grpgmap.WriteHeader(buf, grpgmap.Header{
		Magic:   [8]byte{'G', 'R', 'P', 'G', 'M', 'A', 'P', 0x00},
		Version: 1,
		ChunkX:  uint16(chunkX),
		ChunkY:  uint16(chunkY),
	})

	grpgmap.WriteTiles(buf, tileArr)

	_, err = file.Write(buf.Bytes())
	if err != nil {
		dialog.Message("Error writing to file").Error()
		return
	}
}

func LoadMap() {

}
