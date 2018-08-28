package network

import "time"
import "paxos"

type network interface {
	send(m paxos.Message)
	recv(timeout time.Duration) (paxos.Message, bool)
}
