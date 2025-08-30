package types

import "time"

// --- Estruturas de Dados ---

type Transaction struct {
	ID        string
	Payload   []byte
	Timestamp time.Time
}

type Block struct {
	Transactions []*Transaction
	// Outros campos como Header, ProposerID, etc.
}

// --- Tipos de Mensagens de Consenso ---

type ConsensusMessage interface {
	GetHeight() int
	GetRound() int
}

type ProposalMessage struct {
	Height int
	Round  int
	Block  *Block
}

func (m ProposalMessage) GetHeight() int { return m.Height }
func (m ProposalMessage) GetRound() int  { return m.Round }

type VoteType string

const (
	Prevote   VoteType = "PREVOTE"
	Precommit VoteType = "PRECOMMIT"
)

type VoteMessage struct {
	Height int
	Round  int
	Type   VoteType
	// BlockID // Hash do bloco pelo qual se está votando
}

func (m VoteMessage) GetHeight() int { return m.Height }
func (m VoteMessage) GetRound() int  { return m.Round }

// --- Tipos de Eventos ---

// TransactionEvent representa uma nova transação chegando ao sistema.
type TransactionEvent struct {
	Tx *Transaction
}

// MessageEvent representa uma mensagem de consenso recebida da rede.
type MessageEvent struct {
	SenderID int
	Message  ConsensusMessage
}

// TimeoutEvent representa a expiração de um temporizador de consenso.
type TimeoutEvent struct {
	Height int
	Round  int
	Step   ConsensusStep
}

// ConsensusStep representa as fases dentro de uma rodada do Tendermint.
type ConsensusStep string

const (
	StepPropose   ConsensusStep = "PROPOSE"
	StepPrevote   ConsensusStep = "PREVOTE"
	StepPrecommit ConsensusStep = "PRECOMMIT"
	StepCommit    ConsensusStep = "COMMIT"
)

type Event interface {
	// Type retorna o tipo do evento para despacho.
	Type() string
}
