package orchestrator

import (
	"babel-bft/internal/run"
	"babel-bft/pkg/orchestration"
	"flag"
	"log"
	"os"
	"time"
)

// main é o ponto de entrada principal para o orquestrador do framework.
// Ele pode operar em três modos: remote (mestre), local, ou worker (escravo).
func main() {
	// Definição das flags da linha de comando
	mode := flag.String("mode", "local", "Modo de operação: remote, local, ou worker.")
	protocol := flag.String("protocol", "tendermint", "Protocolo a ser executado: tendermint ou hotstuff.")
	duration := flag.Duration("duration", 10*time.Second, "Duração do experimento (ex: 30s, 1m).")
	nodes := flag.Int("nodes", 4, "Número de nós para executar no modo local.")
	hostsFile := flag.String("hosts", "configs/hosts/local_hosts.txt", "Caminho para o arquivo de hosts para o modo remoto.")
	configFile := flag.String("config", "configs/protocols/tendermint.json", "Caminho para o arquivo de configuração do protocolo.")

	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Lshortfile)

	switch *mode {
	case "remote":
		// Modo Mestre: Orquestra a execução em máquinas remotas
		if *hostsFile == "" {
			log.Fatal("O modo 'remote' requer a flag --hosts com o caminho para o arquivo de hosts.")
		}
		master := orchestration.NewMaster(*protocol, *duration, *hostsFile, *configFile)
		if err := master.Run(); err != nil {
			log.Fatalf("Erro ao executar o mestre: %v", err)
		}
		log.Println("Experimento remoto concluído com sucesso.")

	case "local":
		// Modo Local: Simula N nós na máquina local para testes e depuração
		if err := run.RunLocal(*nodes, *protocol, *duration, *configFile); err != nil {
			log.Fatalf("Erro ao executar a simulação local: %v", err)
		}

	case "worker":
		// Modo Escravo: Executado em máquinas remotas, aguarda comandos do mestre
		log.Println("Iniciando em modo worker...")
		worker := orchestration.NewWorker()
		if err := worker.Run(); err != nil {
			log.Fatalf("Erro ao executar o worker: %v", err)
		}

	default:
		log.Fatalf("Modo desconhecido: %s. Use 'remote', 'local', ou 'worker'.", *mode)
	}
}
