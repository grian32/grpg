package main

import (
	"errors"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgitem"
	"grpg/data-go/grpgmap"
	"grpg/data-go/grpgnpc"
	"grpg/data-go/grpgobj"
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

			bytes, err := os.ReadFile(fullPath)

			if err != nil {
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
	bytes, err := os.ReadFile(path)
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

func LoadNpcs(path string) (map[uint16]*grpgnpc.Npc, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buf := gbuf.NewGBuf(bytes)
	header, err := grpgnpc.ReadHeader(buf)
	if err != nil {
		return nil, err
	}

	if header.Magic != [8]byte{'G', 'R', 'P', 'G', 'N', 'P', 'C', 0x00} {
		return nil, errors.New("file provided does not have GRPGNPC magic")
	}

	npcs, err := grpgnpc.ReadNpcs(buf)
	if err != nil {
		return nil, err
	}

	npcMap := make(map[uint16]*grpgnpc.Npc)
	for _, npc := range npcs {
		npcMap[npc.NpcId] = &npc
	}

	return npcMap, nil
}

func LoadItems(path string) ([]grpgitem.Item, error) {
	bytes, err := os.ReadFile(path)
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
