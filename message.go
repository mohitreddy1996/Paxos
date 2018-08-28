package paxos

type messageType int

// Enums.
const (
	Prepare messageType = iota + 1
	Propose
	Accept
	Promise
)

type Message struct {
	From, To int
	Type     messageType
	Key_     int
	PrevKey  int
	Value    string
}

type PromiseInterface interface {
	Key() int
}

type AcceptInterface interface {
	ProposalKey() int
	ProposalValue() string
}

func (m Message) Key() int {
	return m.Key_
}

func (m Message) ProposalKey() int {
	switch m.Type {
	case Promise:
		return m.PrevKey
	case Accept:
		return m.Key
	default:
		panic("unexpected proposal key.")
	}
}

func (m Message) ProposalValue() string {
	switch m.Type {
	case Promise, Accept:
		return m.Value
	default:
		panic("unexpected proposal value.")
	}
}
