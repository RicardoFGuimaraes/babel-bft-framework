package core

import (
	"log"
	"time"

	"babel-bft/internal/network"
	"babel-bft/internal/protocols"
	"babel-bft/internal/types"
)

// Node represents a single replica in the BFT system. It is the central component
// that connects the network transport, the consensus protocol, and the application logic.
type Node struct {
	id         uint
	Transport  network.Transport
	Engine     protocols.Consensus
	msgChan    chan *types.Message
	stopChan   chan struct{}
	quorumSize int
}

// NewNode creates and initializes a new consensus node.
func NewNode(id uint, transport network.Transport, engine protocols.Consensus, quorum int) *Node {
	return &Node{
		id:         id,
		Transport:  transport,
		Engine:     engine,
		msgChan:    make(chan *types.Message, 100), // Buffered channel
		stopChan:   make(chan struct{}),
		quorumSize: quorum,
	}
}

// Start initiates the node's main event loop in a separate goroutine.
func (n *Node) Start() {
	log.Printf("Node %d starting...", n.ID)
	n.Transport.RegisterNodeChan(n.id, n.msgChan)
	n.Engine.SetNode(n) // Provide the consensus engine with access to the node's interface
	go n.run()
}

// Stop terminates the node's event loop.
func (n *Node) Stop() {
	close(n.stopChan)
}

// The main event loop of the node. It listens for incoming messages
// and passes them to the consensus engine for processing.
func (n *Node) run() {
	log.Printf("Node %d is running.", n.ID)
	for {
		select {
		case msg := <-n.msgChan:
			// Forward the message to the consensus engine
			n.Engine.HandleMessage(msg.From, msg)
		case <-n.stopChan:
			log.Printf("Node %d stopping.", n.ID)
			return
		}
	}
}

// Broadcast sends a message to all other nodes in the network.
// This method implements the types.NodeInterface.
func (n *Node) Broadcast(msg *types.Message) {
	msg.From = n.ID
	n.Transport.Broadcast(msg)
}

// Send directs a message to a specific recipient node.
// This method implements the types.NodeInterface.
func (n *Node) Send(recipientID uint, msg *types.Message) {
	msg.From = n.ID
	n.Transport.Send(recipientID, msg)
}

// id returns the node's unique identifier.
// This method implements the types.NodeInterface.
func (n *Node) ID() uint {
	return n.id
}

// QuorumSize returns the number of replicas in the system.
// This method implements the types.NodeInterface.
func (n *Node) QuorumSize() int {
	return n.quorumSize
}
