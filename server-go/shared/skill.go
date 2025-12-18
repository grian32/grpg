package shared

import (
	"cmp"
	"grpg/data-go/gbuf"
)

type Skill uint8

const (
	FORAGING Skill = iota
)

type SkillInfo struct {
	Level uint8
	XP    uint32
}

func EncodeSkillsToBlob(skills map[Skill]*SkillInfo) []byte {
	buf := gbuf.NewEmptyGBuf()

	// need to update rhs when i add more skills
	for i := FORAGING; i <= FORAGING; i++ {
		buf.WriteByte(skills[i].Level)
		buf.WriteUint32(skills[i].XP)
	}

	return buf.Bytes();
}

func DecodeSkillsFromBlob(blob []byte) (map[Skill]*SkillInfo, error) {
	// since pre-existing players won't have this field, can probably get rid of it at some point but needed during dev
	if len(blob) == 0 {
		skills := make(map[Skill]*SkillInfo)
		for i := FORAGING; i <= FORAGING; i++ {
			skills[i] = &SkillInfo{
				Level: 1,
				XP:    0,
			}
		}
		return skills, nil
	}

	buf := gbuf.NewGBuf(blob);
	skills := make(map[Skill]*SkillInfo);

	for i := FORAGING; i <= FORAGING; i++ {
		level, err1 := buf.ReadByte();
		xp, err2:= buf.ReadUint32();

		if err := cmp.Or(err1, err2); err != nil {
			return make(map[Skill]*SkillInfo), err;
		}

		skills[i] = &SkillInfo{
			Level: level,
			XP: xp,
		}
	}

	return skills, nil
}
