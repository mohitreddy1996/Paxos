package acceptor

import "paxos"

type acceptor interface {
	ReceivePrepare(prepare paxos.Message) (paxos.Message, bool)
	ReceivePropose(propose paxos.Message) bool
}
