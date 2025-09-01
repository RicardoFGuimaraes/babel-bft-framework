package metrics

import (
	"babel-bft/internal/types"
	"log"
	"sync"
	"time"
)

// Collector is responsible for gathering and storing performance metrics
// during an experiment run, such as transaction latency and throughput.
type Collector struct {
	mu           sync.Mutex
	startTime    time.Time
	txTimestamps map[string]int64 // Map transaction hash/ID to its start time
	latencies    []float64
	txCount      int
}

// NewCollector creates a new metrics collector.
func NewCollector() *Collector {
	return &Collector{
		txTimestamps: make(map[string]int64),
		latencies:    make([]float64, 0),
	}
}

// Start begins the collection period.
func (c *Collector) Start() {
	c.startTime = time.Now()
	log.Println("Metrics collection started.")
}

// AddTransaction marks the submission time of a transaction.
func (c *Collector) AddTransaction(tx *types.Transaction) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Using a simple string representation as a key. A hash would be better.
	key := string(tx.Payload)
	c.txTimestamps[key] = tx.Timestamp
}

// FinalizeTransaction marks the commit time of a transaction and calculates its latency.
func (c *Collector) FinalizeTransaction(tx *types.Transaction) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := string(tx.Payload)
	if startTime, ok := c.txTimestamps[key]; ok {
		endTime := time.Now().UnixNano()
		latency := float64(endTime-startTime) / 1_000_000.0 // Latency in milliseconds
		c.latencies = append(c.latencies, latency)
		c.txCount++
		delete(c.txTimestamps, key)
	}
}

// Report calculates and prints the final performance summary.
func (c *Collector) Report() {
	c.mu.Lock()
	defer c.mu.Unlock()

	duration := time.Since(c.startTime).Seconds()
	throughput := float64(c.txCount) / duration

	var totalLatency float64
	for _, l := range c.latencies {
		totalLatency += l
	}
	avgLatency := totalLatency / float64(len(c.latencies))

	log.Println("------ Metrics Report ------")
	log.Printf("Total execution time: %.2f seconds\n", duration)
	log.Printf("Total committed transactions: %d\n", c.txCount)
	log.Printf("Throughput: %.2f TPS\n", throughput)
	log.Printf("Average latency: %.2f ms\n", avgLatency)
	log.Println("--------------------------")
}
