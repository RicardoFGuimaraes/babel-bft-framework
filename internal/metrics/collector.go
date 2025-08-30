package metrics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Collector é responsável por coletar e agregar métricas de performance.
type Collector struct {
	mu        sync.Mutex
	events    map[string]int
	latencies map[string][]time.Duration
}

// NewCollector cria uma nova instância do coletor de métricas.
func NewCollector() *Collector {
	return &Collector{
		events:    make(map[string]int),
		latencies: make(map[string][]time.Duration),
	}
}

// RecordEvent incrementa um contador para um evento nomeado.
func (c *Collector) RecordEvent(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events[name]++
}

// RecordLatency registra uma medição de latência para uma operação nomeada.
func (c *Collector) RecordLatency(name string, d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.latencies[name] = append(c.latencies[name], d)
}

// Save serializa as métricas coletadas para um arquivo JSON.
func (c *Collector) Save(filePath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Cria o diretório se ele não existir
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Estrutura para a saída JSON
	output := struct {
		Events    map[string]int
		Latencies map[string][]time.Duration
	}{
		Events:    c.events,
		Latencies: c.latencies,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
