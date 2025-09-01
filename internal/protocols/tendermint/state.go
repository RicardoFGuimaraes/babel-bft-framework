// File: internal/protocols/tendermint/state.go
package tendermint

import (
	"babel-bft/internal/types"
	"sync"
)

// State holds the current consensus state for a replica.
// It tracks the current height, round, and step of the consensus process.
type State struct {
	mtx sync.RWMutex

	Height int
	Round  int
	Step   string // Propose, Prevote, Precommit

	// Locked block hash and round
	LockedHash  []byte
	LockedRound int

	// Valid block hash and round (the one with +2/3 prevotes)
	ValidHash  []byte
	ValidRound int

	ProposalBlock *types.Block
	Votes         map[int]map[int]map[uint]*PrevoteMessage   // height -> round -> validatorId -> vote
	Commits       map[int]map[int]map[uint]*PrecommitMessage // height -> round -> validatorId -> commit
}

// NewState creates a new state machine for the Tendermint protocol.
func NewState() *State {
	return &State{
		Height:      1,
		Round:       0,
		Step:        "propose",
		LockedRound: -1,
		ValidRound:  -1,
		Votes:       make(map[int]map[int]map[uint]*PrevoteMessage),
		Commits:     make(map[int]map[int]map[uint]*PrecommitMessage),
	}
}

// SetStep updates the current step of the consensus.
func (s *State) SetStep(step string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Step = step
}

// GetHeightRoundStep returns the current height, round, and step.
func (s *State) GetHeightRoundStep() (int, int, string) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.Height, s.Round, s.Step
}

// IncrementRound moves the state to the next round.
func (s *State) IncrementRound() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Round++
	s.Step = "propose"
}

// AddVote stores a prevote message for a given height and round.
func (s *State) AddVote(senderID uint, vote *PrevoteMessage) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.Votes[vote.Height]; !ok {
		s.Votes[vote.Height] = make(map[int]map[uint]*PrevoteMessage)
	}
	if _, ok := s.Votes[vote.Height][vote.Round]; !ok {
		s.Votes[vote.Height][vote.Round] = make(map[uint]*PrevoteMessage)
	}
	s.Votes[vote.Height][vote.Round][senderID] = vote
}

// CountVotes returns the number of prevotes for a specific block hash at a given height and round.
func (s *State) CountVotes(height, round int, hash []byte) int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	count := 0
	if roundVotes, ok := s.Votes[height][round]; ok {
		for _, vote := range roundVotes {
			// A simple byte comparison is sufficient here. For production, use a constant-time comparison.
			if string(vote.Hash) == string(hash) {
				count++
			}
		}
	}
	return count
}
