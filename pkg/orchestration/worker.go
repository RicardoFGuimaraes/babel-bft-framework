package orchestration

import (
	"babel-bft/internal/core"
	"babel-bft/internal/metrics"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Worker representa um nó escravo que executa o protocolo BFT.
type Worker struct {
	node *core.Node
}

// NewWorker cria uma nova instância de um worker.
func NewWorker() *Worker {
	// TODO: A inicialização do nó deve ser mais sofisticada,
	// recebendo configuração do mestre.
	metricsCollector := metrics.NewCollector()
	node := core.NewNode(0, metricsCollector) // id 0 é um placeholder
	return &Worker{node: node}
}

// Run inicia o worker, que escuta por comandos do mestre.
func (w *Worker) Run() error {
	// Inicia o loop de eventos do nó em uma goroutine
	w.node.Start()

	// Configura um servidor HTTP simples para receber comandos do mestre
	http.HandleFunc("/start", w.handleStart)
	http.HandleFunc("/stop", w.handleStop)

	log.Println("Worker escutando por comandos na porta 8080...")
	return http.ListenAndServe(":8080", nil)
}

// handleStart é o handler para o comando de início do experimento.
func (w *Worker) handleStart(rw http.ResponseWriter, r *http.Request) {
	log.Println("Comando 'start' recebido do mestre.")
	// TODO: Aqui, o worker começaria a lógica de consenso ativa,
	// possivelmente após receber a configuração completa.
	// Por enquanto, apenas registramos o evento.
	w.node.Metrics.RecordEvent("experiment_started")
	fmt.Fprintln(rw, "Experimento iniciado.")
}

// handleStop é o handler para o comando de término do experimento.
func (w *Worker) handleStop(rw http.ResponseWriter, r *http.Request) {
	log.Println("Comando 'stop' recebido do mestre.")
	w.node.Metrics.RecordEvent("experiment_stopped")

	// Salva as métricas em um arquivo
	if err := w.node.Metrics.Save("/app/results/metrics.json"); err != nil {
		log.Printf("ERRO: Falha ao salvar métricas: %v", err)
		http.Error(rw, "Falha ao salvar métricas", http.StatusInternalServerError)
		return
	}

	log.Println("Métricas salvas com sucesso.")
	fmt.Fprintln(rw, "Métricas salvas. Encerrando.")

	// Dá um tempo para a resposta HTTP ser enviada antes de sair
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}
