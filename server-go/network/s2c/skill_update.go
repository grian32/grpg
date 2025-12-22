package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
)

type SkillUpdate struct {
	Player *shared.Player
	SkillIds []shared.Skill
}

func (s *SkillUpdate) Opcode() byte {
	return 0x08
}

func (s *SkillUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	buf.WriteByte(byte(len(s.SkillIds)))

	for _, skillId := range s.SkillIds {
		buf.WriteByte(byte(skillId))
		buf.WriteByte(s.Player.Skills[skillId].Level)
		buf.WriteUint32(s.Player.Skills[skillId].XP)
	}
}
