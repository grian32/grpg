package shared

type Skill byte

type SkillInfo struct {
	Level uint8
	XP uint32
}

const (
	Foraging Skill = iota
)
