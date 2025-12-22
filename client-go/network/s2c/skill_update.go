package s2c

import (
	"client/shared"
	"cmp"
	"fmt"
	"grpg/data-go/gbuf"
	"log"
)

type SkillUpdate struct {

}

func (s *SkillUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	len, err := buf.ReadByte();
	if err != nil {
		log.Printf("failed to read length skill update %v\n", err)
		return
	}

	for _ = range len {
		skillId, err1 := buf.ReadByte();
		level, err2 := buf.ReadByte();
		xp, err3 := buf.ReadUint32();

		if err := cmp.Or(err1, err2, err3); err != nil {
			log.Printf("failed to read skill update %v\n", err)
			return
		}

		game.Skills[shared.Skill(skillId)].Level = level;
		game.Skills[shared.Skill(skillId)].XP = xp;
		fmt.Printf("skills update: %v\n", game.Skills[shared.Foraging]);
	}
}
