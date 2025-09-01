package protocols

// File: internal/protocols/consensus.go

import "babel-bft/internal/types"

// Consensus is the interface that every BFT protocol must implement.
// It defines the basic operations for message handling and state transitions.
type Consensus interface {
	// HandleMessage processes an incoming message according to the protocol's logic.
	// It returns true if the message was valid and processed, false otherwise.
	HandleMessage(senderID uint, msg *types.Message) bool

	// CurrentState returns the current state of the consensus engine.
	// This is useful for metrics, logging, and debugging.
	CurrentState() interface{}

	// SetNode assigns the core node logic to the consensus protocol.
	// This allows the protocol to send messages and interact with the node's state.
	SetNode(node types.NodeInterface)
}
