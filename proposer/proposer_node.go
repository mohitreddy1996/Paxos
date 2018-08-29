package proposer

import (
	"log"
	"paxos"
	"paxos/network"
	"time"
)

// Node acts as the proposers in the paxos network.
type Node struct {
	id        int
	lastSeq   int
	value     string
	key       int
	acceptors map[int]paxos.PromiseInterface
	network   network.Network
}

// NewProposer returns a proposer node with the given arguments.
func NewProposer(id int, value string, network network.Network, acceptors ...int) *Node {
	p := &Node{id: id, network: network, lastSeq: 0, value: value, acceptors: make(map[int]paxos.PromiseInterface)}
	for _, a := range acceptors {
		p.acceptors[a] = paxos.Message{}
	}
	return p
}

func (p *Node) majority() int {
	return len(p.acceptors)/2 + 1
}

func (p *Node) uniqueid() int {
	return p.lastSeq<<16 | p.id
}

func (p *Node) ismajority() bool {
	m := 0
	for _, promise := range p.acceptors {
		if promise.Key() == p.uniqueid() {
			m++
		}
	}
	if m >= p.majority() {
		return true
	}
	return false
}

func (p *Node) prepare() []paxos.Message {
	p.lastSeq++
	msgs := make([]paxos.Message, p.majority())
	indx := 0
	for to := range p.acceptors {
		msgs[indx] = paxos.Message{From: p.id, To: to, Type: paxos.Prepare, Key_: p.uniqueid()}
		indx++
		if indx == p.majority() {
			break
		}
	}
	return msgs
}

func (p *Node) propose() []paxos.Message {
	msgs := make([]paxos.Message, p.majority())
	indx := 0
	for to, promise := range p.acceptors {
		if promise.Key() == p.uniqueid() {
			msgs[indx] = paxos.Message{From: p.id, To: to, Type: paxos.Propose, Key_: p.uniqueid(), Value: p.value}
		}
		if indx == p.majority() {
			break
		}
	}
	return msgs
}

func (p *Node) receivePromise(promise paxos.Message) {
	prevPromise := p.acceptors[promise.From]
	if prevPromise.Key() < promise.Key() {
		log.Printf("proposer: received a new proposer with larger key. content: %+v", promise)
		p.acceptors[promise.From] = promise

		if promise.ProposalKey() > p.key {
			log.Printf("proposer: received a promise with larger key. Updating value to : %s", promise.ProposalValue())
			p.key = promise.ProposalKey()
			p.value = promise.ProposalValue()
		}
	}
}

// Run starts the proposer node.
func (p *Node) Run() {
	var ok bool
	var m paxos.Message
	// wait till majority, prepare stage.
	for !p.ismajority() {
		if !ok {
			msgs := p.prepare()
			for i := range msgs {
				p.network.Send(msgs[i])
			}
		}
		m, ok = p.network.Recv(time.Second)
		if !ok {
			continue
		}
		switch m.Type {
		case paxos.Promise:
			p.receivePromise(m)
		default:
			log.Panicf("proposer: Unexpected message type. %d %d", p.id, m.Type)
		}
	}
	// propose stage.
	log.Printf("proposer: Promise reached majority for proposer: %d", p.id)
	log.Printf("proposer: Starting to propose key: %d, value: %s", p.uniqueid(), p.value)
	msg := p.propose()
	for i := range msg {
		p.network.Send(msg[i])
	}
}
