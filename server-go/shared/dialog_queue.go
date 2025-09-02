package shared

import "errors"

type DialogueType byte

const (
	PLAYER DialogueType = iota
	NPC
)

var DialogueEnd = errors.New("dialogue has ended")

type Dialogue struct {
	Type    DialogueType
	Content string
}
type DialogueQueue struct {
	Index     uint16
	MaxIndex  uint16
	Dialogues []Dialogue
}

func (dq *DialogueQueue) Next() (Dialogue, error) {
	if dq.Index >= dq.MaxIndex {
		return Dialogue{}, DialogueEnd
	}
	dialogue := dq.Dialogues[dq.Index]
	dq.Index++
	return dialogue, nil
}

func (dq *DialogueQueue) Clear() {
	dq.Index = 0
	dq.Dialogues = []Dialogue{}
}
