package npcs

import (
	"grpg/data-go/grpgnpc"
	"os"

	"github.com/grian32/gcfg"
)

type ManifestConfig struct {
	Npcs []GRPGNpcManifestEntry `gcfg:"Npc"`
}

type GRPGNpcManifestEntry struct {
	Name    string `gcfg:"name"`
	NpcId   uint16 `gcfg:"id"`
	Texture string `gcfg:"texture"`
}

func BuildGRPGNpcFromManifest(entries []GRPGNpcManifestEntry, texMap map[string]uint16) []grpgnpc.Npc {
	npcArr := make([]grpgnpc.Npc, len(entries))

	for idx, entry := range entries {
		npcArr[idx] = grpgnpc.Npc{
			NpcId:     entry.NpcId,
			Name:      entry.Name,
			TextureId: texMap[entry.Texture],
		}
	}

	return npcArr
}

func ParseManifestFile(path string) ([]GRPGNpcManifestEntry, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ManifestConfig
	err = gcfg.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg.Npcs, nil
}
