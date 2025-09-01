// File: internal/protocols/tendermint/message.go
package tendermint

import "babel-bft/internal/types"

const (
	ProposeType = iota
	PrevoteType
	PrecommitType
)

// ProposeMessage is sent by the proposer for a given height/round.
type ProposeMessage struct {
	Height int
	Round  int
	Block  *types.Block
}

// PrevoteMessage is cast by validators after receiving a valid proposal.
type PrevoteMessage struct {
	Height int
	Round  int
	Hash   []byte // Hash of the proposed block
}

// PrecommitMessage is cast by validators after receiving +2/3 prevotes for a proposal.
type PrecommitMessage struct {
	Height int
	Round  int
	Hash   []byte // Hash of the proposed block
}
