# Etapa de build com base Debian (golang:1.23 não possui versão Alpine ainda)
FROM golang:1.23 AS builder

# Instala certificados e git usando APT (correto para Debian)
RUN apt-get update && apt-get install -y ca-certificates git && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copia os arquivos de dependências primeiro (para cache)
COPY go.mod ./
RUN go mod download

# Copia o restante da aplicação
COPY . .

# Compila a aplicação de forma estática
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

# Etapa final usando Alpine
FROM alpine:latest

# Instala os certificados CA com apk (correto para Alpine)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copia o binário da etapa de build
COPY --from=builder /app/app .

EXPOSE 8080

# Comando de execução
CMD ["./app"]
