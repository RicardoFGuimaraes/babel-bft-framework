// File: internal/network/local.go
package network

import (
	"babel-bft/internal/types"
	"log"
	"sync"
)

// LocalTransport provides an in-memory, channel-based transport implementation.
// It is used for running multiple nodes within a single process for testing and simulation.
type LocalTransport struct {
	mu       sync.RWMutex
	nodeChs  map[uint]chan<- *types.Message
	numNodes uint
}

// NewLocalTransport creates a new LocalTransport.
// numNodes is the total number of nodes that will participate in the network.
func NewLocalTransport(numNodes uint) *LocalTransport {
	return &LocalTransport{
		nodeChs:  make(map[uint]chan<- *types.Message),
		numNodes: numNodes,
	}
}

// RegisterNodeChan registers a channel for a given node id.
func (lt *LocalTransport) RegisterNodeChan(nodeID uint, ch chan<- *types.Message) {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	lt.nodeChs[nodeID] = ch
}

// Start begins the transport layer. For LocalTransport, this is a no-op
// as the channels are managed directly.
func (lt *LocalTransport) Start() {
	log.Println("Local transport started.")
}

// Broadcast sends the message to all registered nodes except the sender.
func (lt *LocalTransport) Broadcast(msg *types.Message) {
	lt.mu.RLock()
	defer lt.mu.RUnlock()

	for id, ch := range lt.nodeChs {
		// Avoid sending the message back to the sender
		if id == msg.From {
			continue
		}
		// Send message in a non-blocking way
		go func(c chan<- *types.Message) {
			c <- msg
		}(ch)
	}
}

// Send delivers a message to a specific recipient.
func (lt *LocalTransport) Send(recipientID uint, msg *types.Message) {
	lt.mu.RLock()
	defer lt.mu.RUnlock()

	if ch, ok := lt.nodeChs[recipientID]; ok {
		// Send message in a non-blocking way
		go func() {
			ch <- msg
		}()
	} else {
		log.Printf("Error: Attempted to send message to unregistered node %d", recipientID)
	}
}
