package scripts

// TODO: auto gen this from .gcfg file, include as part of data-packer?

type ObjConstant uint16
type NpcConstant uint16
type ItemConstant uint16

const (
	_ ObjConstant = iota
	BERRY_BUSH
)

const (
	_ NpcConstant = iota
	TEST
)

const (
	_ ItemConstant = iota
	BERRIES
)
