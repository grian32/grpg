package game

import (
	"client/shared"
	"client/util"
	"cmp"
	"grpg/data-go/gbuf"
	"grpg/data-go/grpgmap"
	"grpg/data-go/grpgtex"
	"io"
	"log"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func loadTextures(path string) map[uint16]rl.Texture2D {
	rlTextures := make(map[uint16]rl.Texture2D)

	grpgTexFile, err1 := os.Open(path)
	grpgTexBytes, err2 := io.ReadAll(grpgTexFile)

	if err := cmp.Or(err1, err2); err != nil {
		log.Fatal("Failed reading GRPGTEX file")
	}

	defer grpgTexFile.Close()

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

			file, err1 := os.Open(fullPath)
			bytes, err2 := io.ReadAll(file)

			if err := cmp.Or(err1, err2); err != nil {
				log.Fatalf("Error reading map file, %v", err)
			}
			defer file.Close()

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
				if obj.InternalId != 0 && obj.Type == grpgmap.OBJ {
					x := (idx % 16) + (int(header.ChunkX) * 16)
					y := (idx / 16) + (int(header.ChunkY) * 16)

					game.CollisionMap[util.Vector2I{X: int32(x), Y: int32(y)}] = struct{}{}
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

func loadGameframeRightTexture(texturePath string) rl.Texture2D {
	file, err1 := os.Open(texturePath)
	bytes, err2 := io.ReadAll(file)

	if err := cmp.Or(err1, err2); err != nil {
		log.Fatalf("errored while trying to load gameframe right texture %s", err.Error())
	}

	image := rl.LoadImageFromMemory(".png", bytes, int32(len(bytes)))
	return rl.LoadTextureFromImage(image)
}

func loadPlayerTextures(dirPath string) map[shared.Direction]rl.Texture2D {
	// FIXME: kind of a shit way to map textures here but easier than loading each one manually so :shrug:
	textureFileNames := []string{"player_back.png", "player_down.png", "player_left.png", "player_right.png"}
	textures := []rl.Texture2D{}

	for _, texPath := range textureFileNames {
		file, err1 := os.Open(dirPath + texPath)
		bytes, err2 := io.ReadAll(file)

		if err := cmp.Or(err1, err2); err != nil {
			log.Fatalf("errored while trying to load player texture %s, %s", texPath, err.Error())
		}

		image := rl.LoadImageFromMemory(".png", bytes, int32(len(bytes)))
		textures = append(textures, rl.LoadTextureFromImage(image))
	}

	return map[shared.Direction]rl.Texture2D{
		shared.UP:    textures[0],
		shared.DOWN:  textures[1],
		shared.LEFT:  textures[2],
		shared.RIGHT: textures[3],
	}
}
