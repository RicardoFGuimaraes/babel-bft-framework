package consensus

import (
	"github.com/ricardofguimaraes/babel-bft-framework/internal/types"
)

// Message é uma interface genérica para mensagens trocadas entre os nós.
type Message interface{}

// CoordinationModule define a interface para a seleção de líder.
type CoordinationModule interface {
	// GetLeader retorna o líder para uma determinada altura/rodada.
	GetLeader(height, round uint64) types.NodeID
}

// DisseminationModule define a interface para a disseminação de propostas.
type DisseminationModule interface {
	// DisseminateProposal envia a proposta para os nós relevantes.
	DisseminateProposal(proposal types.Block) error
}

// TopologyModule define a interface para a topologia da rede.
type TopologyModule interface {
	// GetPeers retorna a lista de pares para comunicação.
	GetPeers() []types.NodeID
}

// ProcessingModule define a interface para o processamento de mensagens.
type ProcessingModule interface {
	// ProcessProposal processa uma proposta recebida.
	ProcessProposal(proposal types.Block) error
	// ProcessVote processa um voto recebido.
	ProcessVote(vote Message) error
	// ProcessNewHeight processa uma nova altura.
	ProcessNewHeight(height uint64) error
}

// Pacemaker define a interface para o controle do avanço do protocolo.
type Pacemaker interface {
	// Start inicia o pacemaker.
	Start()
	// Stop para o pacemaker.
	Stop()
	// AdvanceRound avança para a próxima rodada.
	AdvanceRound()
	// OnReceiveProposal é chamado quando uma proposta é recebida.
	OnReceiveProposal(proposal types.Block)
	// OnReceiveVote é chamado quando um voto é recebido.
	OnReceiveVote(vote Message)
}

// Protocol é a interface que combina todos os módulos para formar um protocolo de consenso.
type Protocol struct {
	Coordination  CoordinationModule
	Dissemination DisseminationModule
	Topology      TopologyModule
	Processing    ProcessingModule
	Pacemaker     Pacemaker
}

// NewProtocol cria uma nova instância de um protocolo a partir dos módulos fornecidos.
func NewProtocol(
	coordination CoordinationModule,
	dissemination DisseminationModule,
	topology TopologyModule,
	processing ProcessingModule,
	pacemaker Pacemaker,
) *Protocol {
	return &Protocol{
		Coordination:  coordination,
		Dissemination: dissemination,
		Topology:      topology,
		Processing:    processing,
		Pacemaker:     pacemaker,
	}
}

// Start inicia o protocolo de consenso.
func (p *Protocol) Start() {
	p.Pacemaker.Start()
}
