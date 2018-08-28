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
	Key      int
	PrevKey  int
	Value    string
}

type PromiseInterface interface {
	key() int
}

type AcceptInterface interface {
	proposalKey() int
	proposalValue() string
}

func (m Message) key() int {
	return m.Key
}

func (m Message) proposalKey() int {
	switch m.Type {
	case Promise:
		return m.PrevKey
	case Accept:
		return m.Key
	default:
		panic("unexpected proposal key.")
	}
}

func (m Message) proposalValue() string {
	switch m.Type {
	case Promise, Accept:
		return m.Value
	default:
		panic("unexpected proposal value.")
	}
}
