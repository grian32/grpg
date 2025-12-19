package s2c

import (
	"grpg/data-go/gbuf"
	"server/shared"
)

type SkillUpdate struct {
	Player *shared.Player
	SkillId shared.Skill
}

func (s *SkillUpdate) Opcode() byte {
	return 0x08
}

func (s *SkillUpdate) Handle(buf *gbuf.GBuf, game *shared.Game) {
	buf.WriteByte(byte(s.SkillId))
	buf.WriteByte(s.Player.Skills[s.SkillId].Level)
	buf.WriteUint32(s.Player.Skills[s.SkillId].XP)
}
