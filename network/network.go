package network

import "time"
import "paxos"

// Network defines the paxos reliable network.
type Network interface {
	send(m paxos.Message)
	recv(timeout time.Duration) (paxos.Message, bool)
}
