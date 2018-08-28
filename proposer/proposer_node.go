package proposer

import "paxos/network"
import "paxos"

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

func (p *Node) prepare() []paxos.Message {
	p.lastSeq++
	msgs := make([]paxos.Message, p.majority())
	indx := 0
	for to := range p.acceptors {
		msgs[indx] = paxos.Message{From: p.id, To: to, Type: paxos.Prepare, Key: p.uniqueid()}
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
		if promise.key() == p.uniqueid() {
			ms[indx] = paxos.Message{From: p.id, To: to, Type: paxos.Propose, Key: p.uniqueid(), Value: p.value}
		}
		if indx == p.majority() {
			break
		}
	}
	return msgs
}
