package orchestration

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Master gerencia a orquestração de um experimento remoto.
type Master struct {
	protocol   string
	duration   time.Duration
	hostsFile  string
	configFile string
	hosts      []string
}

// NewMaster cria uma nova instância do orquestrador mestre.
func NewMaster(protocol string, duration time.Duration, hostsFile string, configFile string) *Master {
	return &Master{
		protocol:   protocol,
		duration:   duration,
		hostsFile:  hostsFile,
		configFile: configFile,
	}
}

// Run executa o fluxo de orquestração completo.
func (m *Master) Run() error {
	log.Println("Iniciando orquestração em modo Mestre...")

	// 1. Carregar a lista de hosts
	if err := m.loadHosts(); err != nil {
		return fmt.Errorf("falha ao carregar hosts: %w", err)
	}
	log.Printf("Hosts carregados: %v", m.hosts)

	// 2. Implantar e iniciar workers em contêineres Docker via SSH
	log.Println("Implantando workers...")
	if err := m.deployWorkers(); err != nil {
		// Tenta limpar em caso de erro
		m.cleanupWorkers()
		return fmt.Errorf("falha ao implantar workers: %w", err)
	}
	log.Println("Todos os workers foram implantados e estão em execução.")

	// TODO: Adicionar lógica para enviar configuração e sinal de início aos workers.
	// Por enquanto, apenas esperamos a duração do experimento.

	log.Printf("Experimento em andamento por %v...", m.duration)
	time.Sleep(m.duration)

	log.Println("Tempo de experimento esgotado. Coletando métricas...")

	// 3. Coletar métricas
	if err := m.collectMetrics(); err != nil {
		log.Printf("Aviso: falha ao coletar métricas: %v", err)
	}

	// 4. Limpar o ambiente
	log.Println("Limpando ambiente...")
	if err := m.cleanupWorkers(); err != nil {
		return fmt.Errorf("falha ao limpar workers: %w", err)
	}

	return nil
}

func (m *Master) loadHosts() error {
	file, err := os.Open(m.hostsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			m.hosts = append(m.hosts, line)
		}
	}
	return scanner.Err()
}

// deployWorkers usa SSH para iniciar contêineres Docker em cada host.
func (m *Master) deployWorkers() error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(m.hosts))
	imageName := "ricardogui/bft-analysis-framework:latest" // Exemplo de nome da imagem

	for i, host := range m.hosts {
		wg.Add(1)
		go func(h string, id int) {
			defer wg.Done()
			log.Printf("Implantando worker %d em %s", id, h)

			// Comando para iniciar o worker em um contêiner Docker
			// Mapeia uma porta para comunicação com o mestre e monta um volume para resultados
			cmdStr := fmt.Sprintf(
				"docker run -d --rm --name worker-%d -p 8080:8080 %s --mode=worker",
				id, imageName,
			)

			cmd := exec.Command("ssh", h, cmdStr)
			output, err := cmd.CombinedOutput()
			if err != nil {
				errChan <- fmt.Errorf("falha ao implantar em %s: %v\nOutput: %s", h, err, string(output))
				return
			}
			log.Printf("Worker %d implantado com sucesso em %s", id, h)
		}(host, i)
	}

	wg.Wait()
	close(errChan)

	// Verifica se houve algum erro durante a implantação
	for err := range errChan {
		return err // Retorna o primeiro erro encontrado
	}

	return nil
}

// collectMetrics copia os arquivos de métricas dos workers para a máquina mestre.
func (m *Master) collectMetrics() error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(m.hosts))
	resultsDir := "results"

	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório de resultados: %w", err)
	}

	for i, host := range m.hosts {
		wg.Add(1)
		go func(h string, id int) {
			defer wg.Done()
			containerName := fmt.Sprintf("worker-%d", id)
			remoteTempPath := fmt.Sprintf("/tmp/metrics-%d.json", id)
			localPath := filepath.Join(resultsDir, fmt.Sprintf("metrics-node-%d.json", id))

			// 1. Copia do contêiner para o host remoto
			copyCmdStr := fmt.Sprintf("docker cp %s:/app/results/metrics.json %s", containerName, remoteTempPath)
			cmd := exec.Command("ssh", h, copyCmdStr)
			if output, err := cmd.CombinedOutput(); err != nil {
				errChan <- fmt.Errorf("falha ao copiar métricas do contêiner em %s: %v\nOutput: %s", h, err, string(output))
				return
			}

			// 2. Copia do host remoto para a máquina mestre
			scpCmd := exec.Command("scp", fmt.Sprintf("%s:%s", h, remoteTempPath), localPath)
			if output, err := scpCmd.CombinedOutput(); err != nil {
				errChan <- fmt.Errorf("falha ao fazer scp de %s: %v\nOutput: %s", h, err, string(output))
				return
			}
			log.Printf("Métricas do worker %d coletadas com sucesso.", id)
		}(host, i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}
	return nil
}

// cleanupWorkers para e remove todos os contêineres dos workers.
func (m *Master) cleanupWorkers() error {
	var wg sync.WaitGroup
	for i, host := range m.hosts {
		wg.Add(1)
		go func(h string, id int) {
			defer wg.Done()
			cmdStr := fmt.Sprintf("docker stop worker-%d", id)
			cmd := exec.Command("ssh", h, cmdStr)
			cmd.Run() // Ignora erros, pois o contêiner pode já não existir
			log.Printf("Worker %d em %s parado.", id, h)
		}(host, i)
	}
	wg.Wait()
	return nil
}
