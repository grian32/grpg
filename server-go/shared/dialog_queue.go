package shared

type DialogueType byte

const (
	PLAYER DialogueType = iota
	NPC
)

type Dialogue struct {
	Type    DialogueType
	Content string
}
type DialogueQueue struct {
	Index       uint16
	MaxIndex    uint16
	ActiveNpcId uint16
	Dialogues   []Dialogue
}

func (dq *DialogueQueue) Clear() {
	dq.Index = 0
	dq.MaxIndex = 0
	dq.ActiveNpcId = 0
	dq.Dialogues = []Dialogue{}
}
