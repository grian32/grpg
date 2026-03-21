package constants

type PlayerVarId uint16
type NpcConstant uint16

const (
	_ PlayerVarId = iota
	SHOULD_SHOW_TUTORIAL_INDICATOR
	LAST_PV = SHOULD_SHOW_TUTORIAL_INDICATOR
)

const (
	_ NpcConstant = iota
	GRPG_GUIDE
)
