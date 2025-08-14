package main

import (
	"cmp"
	"errors"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgitem"
	"grpg/data-go/grpgmap"
	"grpg/data-go/grpgobj"
	"io"
	"log"
	"os"
	"path/filepath"
	"server/shared"
	"server/util"
)

// LoadMaps loads all collisions from maps & adds tracked objs to the map on game
func LoadMaps(dir string, game *shared.Game, objs []grpgobj.Obj) {
	game.CollisionMap = make(map[util.Vector2I]struct{})

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

			zone, err := grpgmap.ReadZone(buf)
			if err != nil {
				log.Fatalf("Error reading grpgmap tiles for file %s", fullPath)
			}

			for idx, obj := range zone.Objs {
				if obj != 0 {
					x := (idx % 16) + (int(header.ChunkX) * 16)
					y := (idx / 16) + (int(header.ChunkY) * 16)

					game.CollisionMap[util.Vector2I{X: uint32(x), Y: uint32(y)}] = struct{}{}

					data := objs[obj-1]

					if grpgobj.IsFlagSet(data.Flags, grpgobj.STATE) {
						game.TrackedObjs[util.Vector2I{X: uint32(x), Y: uint32(y)}] = &shared.GameObj{
							ObjData:  data,
							ChunkPos: util.Vector2I{X: uint32(header.ChunkX), Y: uint32(header.ChunkY)},
							State:    0,
						}
					}
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

func LoadObjs(path string) ([]grpgobj.Obj, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	buf := gbuf.NewGBuf(bytes)
	header, err := grpgobj.ReadHeader(buf)
	if err != nil {
		return nil, err
	}

	if header.Magic != [8]byte{'G', 'R', 'P', 'G', 'O', 'B', 'J', 0x00} {
		return nil, errors.New("file provided does not have GRPGOBJ magic")
	}

	objs, err := grpgobj.ReadObjs(buf)
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func LoadItems(path string) ([]grpgitem.Item, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	buf := gbuf.NewGBuf(bytes)
	header, err := grpgitem.ReadHeader(buf)
	if err != nil {
		return nil, err
	}

	if header.Magic != [8]byte{'G', 'R', 'P', 'G', 'I', 'T', 'E', 'M'} {
		return nil, errors.New("file provided does not have GRPGITEM magic")
	}

	items, err := grpgitem.ReadItems(buf)
	if err != nil {
		return nil, err
	}

	return items, nil
}
