package network

import (
	"babel-bft/internal/types"
	"sync"
	"time"
)

type LocalNetwork struct {
	nodeQueues map[int]chan types.Event
	numNodes   int
	mu         sync.RWMutex
}

func NewLocalNetwork(numNodes int) *LocalNetwork {
	return &LocalNetwork{
		nodeQueues: make(map[int]chan types.Event),
		numNodes:   numNodes,
	}
}

func (ln *LocalNetwork) Register(id int, queue chan types.Event) {
	ln.mu.Lock()
	defer ln.mu.Unlock()
	ln.nodeQueues[id] = queue
}

func (ln *LocalNetwork) Broadcast(senderID int, msg types.ConsensusMessage) {
	ln.mu.RLock()
	defer ln.mu.RUnlock()
	event := &types.MessageEvent{SenderID: senderID, Message: msg}
	for id, queue := range ln.nodeQueues {
		if id != senderID {
			// Simula uma pequena latência de rede
			go func(q chan types.Event) {
				time.Sleep(10 * time.Millisecond)
				q <- event
			}(queue)
		}
	}
}
