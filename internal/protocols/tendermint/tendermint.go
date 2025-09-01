// File: internal/protocols/tendermint/tendermint.go
package tendermint

import (
	"babel-bft/internal/types"
	"log"
)

// Tendermint is the implementation of the Tendermint consensus protocol.
type Tendermint struct {
	node      types.NodeInterface
	state     *State
	pacemaker *Pacemaker
	// More fields can be added here, like a logger, config, etc.
}

// NewTendermint creates a new instance of the Tendermint protocol engine.
func NewTendermint() *Tendermint {
	tm := &Tendermint{
		state: NewState(),
	}
	// The pacemaker will be initialized and started by the node
	// since it needs access to the node's messaging capabilities.
	tm.pacemaker = NewPacemaker(tm)
	return tm
}

// SetNode assigns the core node logic to the consensus protocol.
func (t *Tendermint) SetNode(node types.NodeInterface) {
	t.node = node
	t.pacemaker.node = node // Pacemaker also needs access to the node
}

// HandleMessage processes incoming consensus messages.
func (t *Tendermint) HandleMessage(senderID uint, msg *types.Message) bool {
	log.Printf("Node %d: Received message of type %d from %d", t.node.ID(), msg.Type, senderID)

	// Here we will expand the logic based on the message type and current state
	switch payload := msg.Payload.(type) {
	case *ProposeMessage:
		return t.handlePropose(senderID, payload)
	case *PrevoteMessage:
		return t.handlePrevote(senderID, payload)
	case *PrecommitMessage:
		return t.handlePrecommit(senderID, payload)
	default:
		log.Printf("Node %d: Received unknown message type", t.node.ID())
		return false
	}
}

// handlePropose contains the logic for processing a proposal message.
func (t *Tendermint) handlePropose(sender uint, proposal *ProposeMessage) bool {
	h, r, s := t.state.GetHeightRoundStep()
	log.Printf("Node %d handling Propose from %d for Height %d, Round %d (current state: H:%d, R:%d, S:%s)", t.node.ID(), sender, proposal.Height, proposal.Round, h, r, s)

	// Basic validation: is the proposal timely?
	if proposal.Height < h || (proposal.Height == h && proposal.Round < r) {
		log.Printf("Node %d: Discarding old proposal", t.node.ID())
		return false
	}

	// Further validation (is sender the correct proposer, is block valid?) should be added here.

	// If valid and in the propose or prevote step, we can act on it.
	if s == "propose" || s == "prevote" {
		t.state.SetStep("prevote")
		t.state.ProposalBlock = proposal.Block

		// Broadcast a prevote for this proposal
		prevote := &PrevoteMessage{
			Height: proposal.Height,
			Round:  proposal.Round,
			Hash:   proposal.Block.Hash(), // Assuming the block has a Hash() method
		}
		t.node.Broadcast(&types.Message{Type: PrevoteType, Payload: prevote})
		log.Printf("Node %d: Broadcasted Prevote for H:%d, R:%d", t.node.ID(), proposal.Height, proposal.Round)
		return true
	}

	return false
}

// handlePrevote contains the logic for processing a prevote message.
func (t *Tendermint) handlePrevote(sender uint, prevote *PrevoteMessage) bool {
	h, r, _ := t.state.GetHeightRoundStep()
	log.Printf("Node %d: Handling Prevote from %d for Height %d, Round %d", t.node.ID(), sender, prevote.Height, prevote.Round)

	// Check if the vote is for the current height and round
	if prevote.Height != h || prevote.Round != r {
		return false
	}

	t.state.AddVote(sender, prevote)

	// Check if we have +2/3 prevotes for this block
	// The quorum size should be configurable or passed in, hardcoding for now.
	quorum := (2*t.node.QuorumSize())/3 + 1
	if t.state.CountVotes(h, r, prevote.Hash) >= quorum {
		// We have a polka! Move to precommit step and broadcast precommit.
		t.state.SetStep("precommit")
		t.state.ValidHash = prevote.Hash
		t.state.ValidRound = r

		precommit := &PrecommitMessage{
			Height: h,
			Round:  r,
			Hash:   prevote.Hash,
		}
		t.node.Broadcast(&types.Message{Type: PrecommitType, Payload: precommit})
		log.Printf("Node %d: Reached Prevote quorum. Broadcasting Precommit for H:%d, R:%d", t.node.ID(), h, r)
	}

	return true
}

// handlePrecommit contains the logic for processing a precommit message.
func (t *Tendermint) handlePrecommit(sender uint, precommit *PrecommitMessage) bool {
	// Logic for handling a precommit:
	// 1. Collect precommits.
	// 2. If +2/3 precommits are received, commit the block, execute it, and start a new height.
	// This part is crucial for making progress and will be implemented next.
	log.Printf("Node %d: Handling Precommit from %d for Height %d, Round %d", t.node.ID(), sender, precommit.Height, precommit.Round)

	return true
}

// CurrentState returns the current internal state of the protocol.
func (t *Tendermint) CurrentState() interface{} {
	return t.state
}

// StartNewHeight is called to reset the state for a new consensus instance.
func (t *Tendermint) StartNewHeight() {
	t.state.mtx.Lock()
	t.state.Height++
	t.state.Round = 0
	t.state.Step = "propose"
	t.state.ProposalBlock = nil
	// Clear old votes and commits to prevent memory leaks
	// A more sophisticated garbage collection might be needed for a real implementation.
	t.state.Votes = make(map[int]map[int]map[uint]*PrevoteMessage)
	t.state.Commits = make(map[int]map[int]map[uint]*PrecommitMessage)
	t.state.mtx.Unlock()

	log.Printf("Node %d: Starting new height %d", t.node.ID(), t.state.Height)

	// If this node is the proposer for the new height/round, it should propose a block.
	// The logic to determine the proposer will be added.
}
