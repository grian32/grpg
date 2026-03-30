package constants

// TODO: auto gen this from .gcfg file, include as part of data-packer?

type ObjConstant uint16
type NpcConstant uint16
type ItemConstant uint16
type PlayerVarId uint16

const (
	_ ObjConstant = iota
	BERRY_BUSH
)

const (
	_ NpcConstant = iota
	GRPG_GUIDE
)

const (
	_ ItemConstant = iota
	BERRIES
	BRONZE_HELM
	BRONZE_CHESTPLATE
	BRONZE_LEGS
	BRONZE_RING
	BRONZE_DAGGER
)

const (
	_ PlayerVarId = iota
	SHOULD_SHOW_TUTORIAL_INDICATOR
	LAST_PV = SHOULD_SHOW_TUTORIAL_INDICATOR
)
