package run

import (
	"babel-bft/internal/core"
	"babel-bft/internal/metrics"
	network2 "babel-bft/internal/network"
	"log"
	"sync"
	"time"
)

// RunLocal executa uma simulação completa na máquina local.
func RunLocal(nodeCount int, protocol string, duration time.Duration, configFile string) error {
	log.Printf("Iniciando simulação local com %d nós...", nodeCount)

	var wg sync.WaitGroup
	nodes := make([]*core.Node, nodeCount)
	collectors := make([]*metrics.Collector, nodeCount)
	network := network2.NewLocalNetwork(nodeCount)

	// Inicializa todos os nós
	for i := 0; i < nodeCount; i++ {
		collectors[i] = metrics.NewCollector()
		nodes[i] = core.NewNode(i, nodeCount, protocol, collectors[i])
		nodes[i].SetNetwork(network)
		network.Register(i, nodes[i].EventQueue)
	}

	// Inicia os nós
	for _, node := range nodes {
		wg.Add(1)
		go func(n *core.Node) {
			defer wg.Done()
			n.Start()
		}(node)
	}

	// Inicia o cliente simulado para gerar carga
	client := NewClient(nodes, 10*time.Millisecond) // Envia uma transação a cada 10ms
	stopClient := client.Run()

	// Aguarda a duração do experimento
	log.Printf("Experimento em andamento por %s...", duration)
	time.Sleep(duration)

	// Para o cliente e os nós
	log.Println("Parando simulação...")
	close(stopClient)
	for _, node := range nodes {
		node.Stop()
	}
	wg.Wait()

	// Agrega e exibe as métricas
	aggregateMetrics(collectors)

	log.Println("Simulação local concluída.")
	return nil
}

func aggregateMetrics(collectors []*metrics.Collector) {
	log.Println("\n--- MÉTRICAS AGREGADAS ---")
	totalBlocks := 0
	for i, c := range collectors {
		blocks, ok := c.Events["block_decided"]
		if ok {
			log.Printf("Nó %d decidiu %d blocos.", i, blocks)
			totalBlocks += blocks
		}
	}
	log.Printf("Total de blocos decididos em toda a rede: %d", totalBlocks)
	log.Println("--------------------------")
}
