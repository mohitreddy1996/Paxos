package network

import "time"
import "paxos"

// Network defines the paxos reliable network.
type Network interface {
	Send(m paxos.Message)
	Recv(timeout time.Duration) (paxos.Message, bool)
}
