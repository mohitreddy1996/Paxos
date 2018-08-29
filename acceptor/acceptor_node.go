package acceptor

import (
	"log"
	"paxos"
	"paxos/network"
	"time"
)

// Node represents an acceptor in the paxos network.
type Node struct {
	id       int
	learners []int
	accept   paxos.Message
	promised paxos.PromiseInterface

	network network.Network
}

func NewAcceptor(id int, network network.Network, learners ...int) *Node {
	return &Node{id: id, learners: learners, network: network, promised: paxos.Message{}}
}

func (n *Node) ReceivePrepare(prepare paxos.Message) (paxos.Message, bool) {
	if n.promised.Key() >= prepare.Key() {
		return paxos.Message{}, false
	}
	log.Printf("acceptor: Accepting prepare. %+v", prepare)
	n.promised = prepare
	m := paxos.Message{
		Type:    paxos.Promise,
		From:    n.id,
		To:      prepare.From,
		Key_:    n.promised.Key(),
		PrevKey: n.accept.Key_,
		Value:   n.accept.Value,
	}
	return m, true
}

func (n *Node) ReceivePropose(propose paxos.Message) bool {
	if n.promised.Key() > propose.Key() {
		return false
	}
	if n.promised.Key() < propose.Key() {
		log.Panic("This should not happen.")
	}
	log.Printf("acceptor: Accepting proposal: %+v. Promised: %+v and accept: %+v", propose, n.promised, n.accept)
	n.accept = propose
	n.accept.Type = paxos.Accept
	return true
}

func (n *Node) run() {
	for {
		m, ok := n.network.Recv(time.Hour)
		if !ok {
			continue
		}
		switch m.Type {
		case paxos.Prepare:
			promise, ok := n.ReceivePrepare(m)
			if ok {
				n.network.Send(promise)
			}
		case paxos.Promise:
			accepted := n.ReceivePropose(m)
			if accepted {
				for _, l := range n.learners {
					m := n.accept
					m.From = n.id
					m.To = l
					n.network.Send(m)
				}
			}
		default:
			log.Panic("acceptor: Undefined message type received.")
		}
	}
}
