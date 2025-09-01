package run

import (
	"log"
	"time"

	"babel-bft/internal/core"
	"babel-bft/internal/network"
	"babel-bft/internal/protocols/tendermint"
)

// LocalSimulation sets up and runs a BFT consensus simulation in-process.
// It creates a specified number of nodes and clients, connects them via an
// in-memory transport layer, and runs the simulation for a fixed duration.
func LocalSimulation(numNodes, numClients uint, duration time.Duration) {
	log.Printf("Starting local simulation with %d nodes, %d clients for %s.", numNodes, numClients, duration)

	// 1. Initialize the local network transport
	transport := network.NewLocalTransport(numNodes + numClients)

	// 2. Create and start the consensus nodes (replicas)
	nodes := make([]*core.Node, numNodes)
	for i := uint(0); i < numNodes; i++ {
		// Each node gets its own instance of the consensus engine
		engine := tendermint.NewTendermint()
		nodes[i] = core.NewNode(i, transport, engine, int(numNodes))
		nodes[i].Start()
	}

	// 3. Create and start the clients
	clients := make([]*core.Client, numClients)
	for i := uint(0); i < numClients; i++ {
		// Client IDs start after the last node ID
		clientID := numNodes + i
		clients[i] = core.NewClient(clientID, transport)
		clients[i].Start()
	}

	// 4. Run the simulation for the specified duration
	log.Printf("Simulation running for %s...", duration)
	time.Sleep(duration)

	// 5. Stop all clients and nodes
	log.Println("Simulation duration ended. Stopping all components...")
	for _, client := range clients {
		client.Stop()
	}
	for _, node := range nodes {
		node.Stop()
	}

	log.Println("Simulation finished.")
	// In a real scenario, we would collect and report metrics here.
}
