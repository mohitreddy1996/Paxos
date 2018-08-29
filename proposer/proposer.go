package proposer

import "paxos"

// Proposer defines the functions to be implemented by the proposer node of paxos.
type Proposer interface {
	Propose() []paxos.Message
	Prepare() []paxos.Message
	Run()
}
