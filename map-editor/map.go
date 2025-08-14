package main

import (
	"fmt"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgmap"
	"os"

	"github.com/sqweek/dialog"
)

func SaveMap() {
	tileArr := [256]grpgmap.Tile{}
	objArr := [256]grpgmap.Obj{}

	if chunkX == -1 || chunkY == -1 {
		dialog.Message("Both Chunk X & Chunk Y must be set to save a map.").Error()
		return
	}

	for idx, id := range gridTiles {
		if id == -1 {
			dialog.Message("All tiles must be filled in to save a map.").Error()
			return
		}
		tileArr[idx] = grpgmap.Tile(id)
	}

	for idx, id := range gridObjs {
		var structId uint16 = 0
		if id != -1 {
			structId = uint16(id)
		}
		objArr[idx] = grpgmap.Obj(structId)
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
		Magic:  [8]byte{'G', 'R', 'P', 'G', 'M', 'A', 'P', 0x00},
		ChunkX: uint16(chunkX),
		ChunkY: uint16(chunkY),
	})

	zone := grpgmap.Zone{
		Tiles: tileArr,
		Objs:  objArr,
	}

	grpgmap.WriteZone(buf, zone)

	_, err = file.Write(buf.Bytes())
	if err != nil {
		dialog.Message("Error writing to file").Error()
		return
	}
}

func LoadMap() {
	if len(textures) == 0 {
		dialog.Message("No textures loaded.").Error()
		return
	}

	fileToLoad, err := dialog.File().Title("Please select a .grpgmap file").Load()
	if err != nil {
		dialog.Message("Error finding file to load").Error()
		return
	}

	fileBytes, err := os.ReadFile(fileToLoad)
	if err != nil {
		dialog.Message("Error reading file").Error()
		return
	}

	buf := gbuf.NewGBuf(fileBytes)
	header, err := grpgmap.ReadHeader(buf)
	if err != nil {
		fmt.Println("reading grpgmap header errored: %w. file: %s", err, fileToLoad)
		return
	}

	if string(header.Magic[:]) != "GRPGMAP\x00" {
		dialog.Message("File isn't valid GRPGMAP format.").Error()
		return
	}

	chunkX = int32(header.ChunkX)
	chunkY = int32(header.ChunkY)

	zone, err := grpgmap.ReadZone(buf)
	if err != nil {
		fmt.Println("reading grpgmap tiles errored: %w. file: %s", err, fileToLoad)
		return
	}

	for idx, tile := range zone.Tiles {
		gridTiles[idx] = int32(tile)
		gridTileTextures[idx] = int32(tiles[int32(tile)].TexId)
	}

	for idx, obj := range zone.Objs {
		if obj == 0 {
			continue
		}
		gridObjs[idx] = int32(obj)
		gridObjTextures[idx] = int32(objs[int32(obj)].Textures[0])
	}
}
