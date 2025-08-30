# ---- Estágio de Build ----
FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependência e baixa as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código-fonte
COPY . .

# Compila o orquestrador de forma estática para não depender de bibliotecas do sistema
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orchestrator ./cmd/orchestrator

# ---- Estágio Final ----
FROM alpine:latest

WORKDIR /app

# Copia o binário compilado do estágio de build
COPY --from=builder /app/orchestrator .

# Cria um diretório para os resultados dentro do contêiner
RUN mkdir -p /app/results

# O entrypoint é o nosso binário. As flags serão passadas no `docker run`.
ENTRYPOINT ["./orchestrator"]
