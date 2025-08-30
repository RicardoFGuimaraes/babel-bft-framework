package tendermint

import (
	"babel-bft/internal/core"
	"babel-bft/internal/types"
	"log"
)

// Pacemaker implementa a lógica de coordenação do Tendermint.
type Pacemaker struct {
	node *core.Node

	// Estado do Tendermint
	height int
	round  int
	step   types.ConsensusStep
	// Mais estados: lockedValue, lockedRound, etc.
}

// NewPacemaker cria um novo pacemaker para o Tendermint.
func NewPacemaker(node *core.Node) *Pacemaker {
	return &Pacemaker{
		node:   node,
		height: 1,
		round:  0,
		step:   types.StepPropose,
	}
}

func (p *Pacemaker) HandleTransaction(event *types.TransactionEvent) {
	// Lógica do líder para criar um bloco e iniciar uma nova rodada
	if p.isLeader() {
		log.Printf("[Nó %d] [H:%d, R:%d] Sou o líder, propondo novo bloco.", p.node.ID, p.height, p.round)
		proposal := types.ProposalMessage{
			Height: p.height,
			Round:  p.round,
			Block:  &types.Block{Transactions: []*types.Transaction{event.Tx}},
		}
		p.node.GetNetwork().Broadcast(p.node.ID, proposal)
		p.step = types.StepPrevote
		p.node.Metrics.RecordEvent("proposal_sent")
	}
}

func (p *Pacemaker) HandleConsensusMessage(event *types.MessageEvent) {
	// Lógica para processar mensagens de Propose, Prevote, Precommit
	switch msg := event.Message.(type) {
	case types.ProposalMessage:
		log.Printf("[Nó %d] [H:%d, R:%d] Proposta recebida do nó %d.", p.node.ID, msg.Height, msg.Round, event.SenderID)
		// 1. Validar proposta
		// 2. Enviar Prevote
		prevote := types.VoteMessage{
			Height: msg.Height,
			Round:  msg.Round,
			Type:   types.Prevote,
		}
		p.node.GetNetwork().Broadcast(p.node.ID, prevote)
		p.step = types.StepPrecommit
		p.node.Metrics.RecordEvent("prevote_sent")

	case types.VoteMessage:
		// Lógica para contar votos e avançar para o próximo passo ou decidir
		log.Printf("[Nó %d] [H:%d, R:%d] Voto '%s' recebido do nó %d.", p.node.ID, msg.Height, msg.Round, msg.Type, event.SenderID)
		// TODO: Implementar a lógica de contagem de quórum (2f+1)
		// Se quórum de Prevotes -> envia Precommit
		// Se quórum de Precommits -> decide o bloco, incrementa a altura e inicia nova rodada
		p.node.Metrics.RecordEvent("block_decided")
		p.height++
		p.round = 0
		p.step = types.StepPropose
		log.Printf("[Nó %d] Bloco decidido! Avançando para a altura %d.", p.node.ID, p.height)
	}
}

func (p *Pacemaker) HandleTimeout(event *types.TimeoutEvent) {
	// Lógica para lidar com timeouts, e.g., avançar para a próxima rodada
	log.Printf("[Nó %d] [H:%d, R:%d] Timeout! Avançando para a rodada %d.", p.node.ID, p.height, p.round, p.round+1)
	p.round++
	p.step = types.StepPropose
	// O novo líder (se for este nó) irá propor.
}

// isLeader verifica se o nó atual é o líder para a rodada atual.
func (p *Pacemaker) isLeader() bool {
	// Lógica de rotação de líder simples (round-robin)
	// TODO: Obter o número total de nós da configuração
	numNodes := 4
	return p.node.ID == p.round%numNodes
}
