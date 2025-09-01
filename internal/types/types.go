package types

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

// Constants for message types
const (
	TxMsg = iota
	ConsensusMsg
)

// Message is the generic container for all communications between nodes.
type Message struct {
	Type    int
	From    uint
	Payload interface{}
}

// Transaction represents a client's request to be processed by the state machine.
type Transaction struct {
	ClientID  uint
	Timestamp int64
	Payload   []byte
}

// Block is a collection of transactions that will be atomically applied to the state machine.
type Block struct {
	ProposerID   uint
	Transactions []*Transaction
	HashCache    []byte
}

// Hash calculates and returns the SHA-256 hash of the block.
// The hash is cached for performance.
func (b *Block) Hash() []byte {
	if b.HashCache != nil {
		return b.HashCache
	}
	h := sha256.New()
	h.Write([]byte(strconv.Itoa(int(b.ProposerID))))
	for _, tx := range b.Transactions {
		h.Write([]byte(strconv.Itoa(int(tx.ClientID))))
		h.Write([]byte(strconv.FormatInt(tx.Timestamp, 10)))
		h.Write(tx.Payload)
	}
	b.HashCache = h.Sum(nil)
	return b.HashCache
}

// String provides a simple string representation of the block.
func (b *Block) String() string {
	return fmt.Sprintf("Block{Proposer: %d, Txs: %d, Hash: %x}", b.ProposerID, len(b.Transactions), b.Hash())
}

// NodeInterface defines the set of methods that the consensus engine can use
// to interact with the underlying node, abstracting away the network and core logic.
type NodeInterface interface {
	ID() uint
	QuorumSize() int
	Broadcast(msg *Message)
	Send(recipientID uint, msg *Message)
}
