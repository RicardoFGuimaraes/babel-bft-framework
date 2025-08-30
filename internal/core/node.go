package core

import (
	"babel-bft/internal/metrics"
	"babel-bft/internal/types"
	"log"
	"time"
)

// Event representa um evento genérico a ser processado pelo nó.

// Node é a estrutura central que representa um participante do consenso.
type Node struct {
	ID         int
	eventQueue chan types.Event
	Metrics    *metrics.Collector
	// Outros campos como:
	// modules *ModuleRegistry
	// config *ProtocolConfig
}

// NewNode cria e inicializa um novo nó.
func NewNode(id int, metricsCollector *metrics.Collector) *Node {
	return &Node{
		ID:         id,
		eventQueue: make(chan types.Event, 1024), // Buffer para eventos
		Metrics:    metricsCollector,
	}
}

// Start inicia o loop de processamento de eventos do nó.
func (n *Node) Start() {
	log.Printf("[Nó %d] Iniciando loop de eventos...", n.ID)
	go n.eventLoop()
}

// eventLoop é o laço principal que consome e processa eventos da fila.
func (n *Node) eventLoop() {
	for event := range n.eventQueue {
		// Simplesmente loga o evento por enquanto.
		// A lógica real despacharia o evento para o módulo correto.
		log.Printf("[Nó %d] Processando evento: %T", n.ID, event)

		// Exemplo de como um módulo seria invocado e metrificado
		startTime := time.Now()
		// n.modules.Coordination.HandleEvent(event)
		n.Metrics.RecordLatency("event_processing", time.Since(startTime))
	}
}

// SubmitEvent adiciona um novo evento à fila de processamento do nó.
func (n *Node) SubmitEvent(e types.Event) {
	n.eventQueue <- e
}
