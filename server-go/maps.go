package main

import (
	"cmp"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgmap"
	"grpg/data-go/grpgtex"
	"io"
	"log"
	"os"
	"path/filepath"
	"server/shared"
	"server/util"
)

func LoadCollisionMaps(game *shared.Game) {
	game.CollisionMap = make(map[util.Vector2I]struct{})

	dir := "../../grpg-assets/maps/"
	entries, err := os.ReadDir(dir)

	if err != nil {
		log.Fatal("Error reading maps directory")
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fullPath := filepath.Join(dir, entry.Name())

			file, err1 := os.Open(fullPath)
			bytes, err2 := io.ReadAll(file)

			if err := cmp.Or(err1, err2); err != nil {
				log.Fatalf("Error reading map file, %v", err)
			}

			buf := gbuf.NewGBuf(bytes)
			header, err := grpgmap.ReadHeader(buf)
			if err != nil {
				log.Fatalf("Error reading grpgmap header for file %s", fullPath)
			}

			if string(header.Magic[:]) != "GRPGMAP\x00" {
				log.Fatalf("File %s isn't GRPGMAP", fullPath)
			}

			tiles, err := grpgmap.ReadTiles(buf)
			if err != nil {
				log.Fatalf("Error reading grpgmap tiles for file %s", fullPath)
			}

			for idx, tile := range tiles {
				if tile.Type == grpgtex.OBJ {
					x := (idx % 16) + (int(header.ChunkX) * 16)
					y := (idx / 16) + (int(header.ChunkY) * 16)

					game.CollisionMap[util.Vector2I{X: uint32(x), Y: uint32(y)}] = struct{}{}
				}
			}

			// todo: ugly cALC
			if uint32(((header.ChunkX+1)*16)-1) > game.MaxX {
				game.MaxX = uint32(((header.ChunkX + 1) * 16) - 1)
			}

			if uint32(((header.ChunkY+1)*16)-1) > game.MaxY {
				game.MaxY = uint32(((header.ChunkY + 1) * 16) - 1)
			}
		}
	}
}
