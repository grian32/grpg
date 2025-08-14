package game

import (
	"client/shared"
	"client/util"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgitem"
	"grpg/data-go/grpgmap"
	"grpg/data-go/grpgobj"
	"grpg/data-go/grpgtex"
	"grpg/data-go/grpgtile"
	"log"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func loadTextures(path string) map[uint16]rl.Texture2D {
	rlTextures := make(map[uint16]rl.Texture2D)

	grpgTexBytes, err := os.ReadFile(path)

	if err != nil {
		log.Fatal("Failed reading GRPGTEX file")
	}

	buf := gbuf.NewGBuf(grpgTexBytes)
	header, err := grpgtex.ReadHeader(buf)
	if err != nil {
		log.Fatalf("failed reading grpgtex header: %v", err)
	}

	if string(header.Magic[:]) != "GRPGTEX\x00" {
		log.Fatal("File is not GRPGTEX file.")
	}

	textures, err := grpgtex.ReadTextures(buf)
	if err != nil {
		log.Fatalf("failed reading grpgtex textures: %v", err)
	}

	for _, tex := range textures {
		rlImage := rl.LoadImageFromMemory(".png", tex.PNGBytes, int32(len(tex.PNGBytes)))
		rlTex := rl.LoadTextureFromImage(rlImage)

		rlTextures[tex.InternalIdInt] = rlTex
	}

	return rlTextures
}

// loadMaps returns a map of zone, while mutating the passed in game to set collision maps and max x/y
func loadMaps(dirPath string, game *shared.Game) map[util.Vector2I]grpgmap.Zone {
	zoneMap := make(map[util.Vector2I]grpgmap.Zone)
	entries, err := os.ReadDir(dirPath)

	if err != nil {
		log.Fatal("Error reading maps directory")
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fullPath := filepath.Join(dirPath, entry.Name())

			bytes, err := os.ReadFile(fullPath)

			if err != nil {
				log.Fatalf("Error reading map file, %v", err)
			}

			buf := gbuf.NewGBuf(bytes)
			header, err := grpgmap.ReadHeader(buf)
			if err != nil {
				log.Fatalf("reading grpgmap header errored: %v. file: %s", err, fullPath)
			}

			if string(header.Magic[:]) != "GRPGMAP\x00" {
				log.Fatalf("File %s isn't GRPGMAP", fullPath)
			}

			zone, err := grpgmap.ReadZone(buf)
			if err != nil {
				log.Fatalf("reading grpgmap tiles errored: %v. file: %s", err, fullPath)
			}

			chunkPos := util.Vector2I{X: int32(header.ChunkX), Y: int32(header.ChunkY)}

			zoneMap[chunkPos] = zone

			for idx, obj := range zone.Objs {
				data := game.Objs[uint16(obj)]
				if obj != 0 {
					x := (idx % 16) + (int(header.ChunkX) * 16)
					y := (idx / 16) + (int(header.ChunkY) * 16)

					vec := util.Vector2I{X: int32(x), Y: int32(y)}

					if grpgobj.IsFlagSet(data.Flags, grpgobj.STATE) {
						game.ObjIdByLoc[vec] = data.ObjId
					}

					game.CollisionMap[vec] = struct{}{}
				}
			}

			if ((header.ChunkX+1)*16)-1 > game.MaxX {
				game.MaxX = ((header.ChunkX + 1) * 16) - 1
			}

			if ((header.ChunkY+1)*16)-1 > game.MaxY {
				game.MaxY = ((header.ChunkY + 1) * 16) - 1
			}
		}
	}

	return zoneMap
}

func loadObjs(path string) map[uint16]*grpgobj.Obj {
	objMap := make(map[uint16]*grpgobj.Obj)

	grpgObjBytes, err := os.ReadFile(path)

	if err != nil {
		log.Fatal("Failed reading GRPGOBJ file")
	}

	buf := gbuf.NewGBuf(grpgObjBytes)

	header, err := grpgobj.ReadHeader(buf)
	if err != nil {
		log.Fatal(err)
	}

	if header.Magic != [8]byte{'G', 'R', 'P', 'G', 'O', 'B', 'J', 0x00} {
		log.Fatal("file does not have GRPGOBJ header")
	}

	objs, err := grpgobj.ReadObjs(buf)
	if err != nil {
		log.Fatal(err)
	}

	for _, obj := range objs {
		objMap[obj.ObjId] = &obj
	}

	return objMap
}

func loadTiles(path string) map[uint16]*grpgtile.Tile {
	tileMap := make(map[uint16]*grpgtile.Tile)

	grpgTileBytes, err := os.ReadFile(path)

	if err != nil {
		log.Fatal("Failed reading GRPGTILE file")
	}

	buf := gbuf.NewGBuf(grpgTileBytes)

	header, err := grpgtile.ReadHeader(buf)
	if err != nil {
		log.Fatal(err)
	}

	if header.Magic != [8]byte{'G', 'R', 'P', 'G', 'T', 'I', 'L', 'E'} {
		log.Fatal("file does not have GRPGTILE header")
	}

	tiles, err := grpgtile.ReadTiles(buf)
	if err != nil {
		log.Fatal(err)
	}

	for _, tile := range tiles {
		tileMap[tile.TileId] = &tile
	}

	return tileMap
}

func loadItems(path string) map[uint16]grpgitem.Item {
	itemMap := make(map[uint16]grpgitem.Item)

	grpgItemBytes, err := os.ReadFile(path)

	if err != nil {
		log.Fatal("Failed reading GRPGTILE file")
	}

	buf := gbuf.NewGBuf(grpgItemBytes)

	header, err := grpgitem.ReadHeader(buf)
	if err != nil {
		log.Fatal(err)
	}

	if header.Magic != [8]byte{'G', 'R', 'P', 'G', 'I', 'T', 'E', 'M'} {
		log.Fatal("file does not have GRPGITEM header")
	}

	items, err := grpgitem.ReadItems(buf)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		itemMap[item.ItemId] = item
	}

	return itemMap
}

func loadGameframeRightTexture(texturePath string) rl.Texture2D {
	bytes, err := os.ReadFile(texturePath)

	if err != nil {
		log.Fatalf("errored while trying to load gameframe right texture %v", err)
	}

	image := rl.LoadImageFromMemory(".png", bytes, int32(len(bytes)))
	return rl.LoadTextureFromImage(image)
}

func loadPlayerTextures(dirPath string) map[shared.Direction]rl.Texture2D {
	textureMap := map[shared.Direction]string{
		shared.UP:    "player_back.png",
		shared.DOWN:  "player_down.png",
		shared.LEFT:  "player_left.png",
		shared.RIGHT: "player_right.png",
	}
	textures := map[shared.Direction]rl.Texture2D{}

	for direction, texPath := range textureMap {
		fullPath := filepath.Join(dirPath, texPath)
		bytes, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatal("errored while trying to load player texture %s", fullPath)
		}

		image := rl.LoadImageFromMemory(".png", bytes, int32(len(bytes)))
		textures[direction] = rl.LoadTextureFromImage(image)
	}

	return textures
}
