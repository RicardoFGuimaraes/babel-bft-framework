package core

import (
	"babel-bft/internal/types"
	"time"
)

type Client struct {
	nodes    []*core.Node
	interval time.Duration
}

func NewClient(nodes []*core.Node, interval time.Duration) *Client {
	return &Client{nodes: nodes, interval: interval}
}

func (c *Client) Run() chan<- struct{} {
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				tx := &types.Transaction{
					ID:        time.Now().String(),
					Payload:   []byte("dados da transacao"),
					Timestamp: time.Now(),
				}
				// Envia a transação para um nó aleatório (simulando um ponto de entrada)
				node := c.nodes[0]
				node.SubmitEvent(&types.TransactionEvent{Tx: tx})
			case <-stop:
				return
			}
		}
	}()
	return stop
}
