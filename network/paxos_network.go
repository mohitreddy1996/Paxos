package network

import (
	"log"
	"paxos"
	"time"
)

// PaxosNetwork is the network in which all the nodes operate.
type PaxosNetwork struct {
	receiveQueues map[int]chan paxos.Message
}

// NodeNetwork is the network instance for each node with id.
type NodeNetwork struct {
	id int
	*PaxosNetwork
}

// NewPaxosNetwork creates a new paxos network with the given agents as nodes in it.
func NewPaxosNetwork(agents ...int) *PaxosNetwork {
	paxosNetwork := &PaxosNetwork{
		receiveQueues: make(map[int]chan paxos.Message, 0),
	}

	for _, a := range agents {
		paxosNetwork.receiveQueues[a] = make(chan paxos.Message, 512)
	}
	return paxosNetwork
}

func (paxosNetwork *PaxosNetwork) nodeNetwork(id int) *NodeNetwork {
	return &NodeNetwork{id: id, PaxosNetwork: paxosNetwork}
}

func (paxosNetwork *PaxosNetwork) send(msg paxos.Message) {
	log.Printf("networklog: send %+v", msg)
	paxosNetwork.receiveQueues[msg.To] <- msg
}

func (paxosNetwork *PaxosNetwork) empty() bool {
	var n int
	for _, q := range paxosNetwork.receiveQueues {
		n += len(q)
	}
	return n == 0
}

func (paxosNetwork *PaxosNetwork) recvFrom(from int, timeout time.Duration) (paxos.Message, bool) {
	select {
	case m := <-paxosNetwork.receiveQueues[from]:
		log.Printf("networklog: recv %+v", m)
		return m, true
	case <-time.After(timeout):
		return paxos.Message{}, false
	}
}

func (nodeNetwork *NodeNetwork) send(m paxos.Message) {
	nodeNetwork.PaxosNetwork.send(m)
}

func (nodeNetwork *NodeNetwork) recv(timeout time.Duration) (paxos.Message, bool) {
	return nodeNetwork.recvFrom(nodeNetwork.id, timeout)
}
