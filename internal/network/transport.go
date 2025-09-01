// File: internal/network/transport.go
package network

import "babel-bft/internal/types"

// Transport is an interface for network communication between nodes.
// It allows broadcasting messages to all nodes or sending a message to a specific node.
type Transport interface {
	// Broadcast sends a message to all other nodes in the network.
	Broadcast(msg *types.Message)

	// Send sends a message to a specific recipient node.
	Send(recipientID uint, msg *types.Message)

	// RegisterNodeChan registers a channel for a specific node to receive messages.
	// This is essential for the transport layer to deliver incoming messages to the correct node.
	RegisterNodeChan(nodeID uint, ch chan<- *types.Message)

	// Start initializes the transport layer.
	Start()
}
