// File: internal/protocols/tendermint/pacemaker.go
package tendermint

import (
	"babel-bft/internal/types"
	"log"
	"time"
)

// Pacemaker is responsible for ensuring the liveness of the Tendermint protocol.
// It uses timeouts to trigger round changes when progress is not being made.
type Pacemaker struct {
	protocol   *Tendermint
	node       types.NodeInterface
	timer      *time.Timer
	timeout    time.Duration
	active     bool
	round      int
	lastHeard  time.Time
	isProposer bool
}

// NewPacemaker creates a new Pacemaker instance.
func NewPacemaker(protocol *Tendermint) *Pacemaker {
	return &Pacemaker{
		protocol: protocol,
		// Timeout duration should be configurable
		timeout: 5 * time.Second,
		active:  false,
	}
}

// Start initiates the pacemaker for a given round.
func (p *Pacemaker) Start(isProposer bool) {
	p.isProposer = isProposer
	p.active = true
	p.resetTimer()
	log.Printf("Node %d: Pacemaker started for round %d", p.node.ID(), p.protocol.state.Round)
	go p.run()
}

// Stop deactivates the pacemaker.
func (p *Pacemaker) Stop() {
	p.active = false
	if p.timer != nil {
		p.timer.Stop()
	}
	log.Printf("Node %d: Pacemaker stopped for round %d", p.node.ID(), p.protocol.state.Round)
}

// Internal run loop for the pacemaker timer.
func (p *Pacemaker) run() {
	for p.active {
		<-p.timer.C
		if p.active {
			p.handleTimeout()
		}
	}
}

// handleTimeout is called when the timer expires. It triggers a new round.
func (p *Pacemaker) handleTimeout() {
	currentHeight, currentRound, _ := p.protocol.state.GetHeightRoundStep()
	log.Printf("Node %d: Pacemaker timeout! H:%d R:%d. Advancing to next round.", p.node.ID(), currentHeight, currentRound)

	// Advance to the next round in the protocol state
	p.protocol.state.IncrementRound()

	// Propose nil and broadcast prevote for nil
	// This is part of the Tendermint recovery mechanism
	prevote := &PrevoteMessage{
		Height: currentHeight,
		Round:  p.protocol.state.Round,
		Hash:   nil, // Nil prevote
	}
	p.node.Broadcast(&types.Message{Type: PrevoteType, Payload: prevote})

	// Reset the timer for the new round
	p.resetTimer()
}

// resetTimer resets the timeout timer.
func (p *Pacemaker) resetTimer() {
	if p.timer != nil {
		p.timer.Stop()
	}
	p.timer = time.NewTimer(p.timeout)
	p.lastHeard = time.Now()
}
