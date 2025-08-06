.PHONY: build run test test-cover test-bench docker-build docker-run docker-compose-up docker-compose-down clean

# Build da aplicação
build:
	go build -o bin/weatherbycep main.go

# Executar localmente
run:
	go run main.go

# Executar testes
test:
	go test -v

# Testes com coverage
test-cover:
	go test -cover

# Benchmark
test-bench:
	go test -bench=.

# Todos os testes
test-all: test test-cover test-bench

# Build Docker
docker-build:
	docker build -t weatherbycep .

# Executar com Docker
docker-run:
	docker run -p 8080:80 weatherbycep

# Docker Compose up
docker-compose-up:
	docker-compose up --build

# Docker Compose down
docker-compose-down:
	docker-compose down

# Limpar arquivos de build
clean:
	rm -rf bin/
	docker-compose down
	docker rmi weatherbycep 2>/dev/null || true

# Help
help:
	@echo "Comandos disponíveis:"
	@echo "  build             - Build da aplicação"
	@echo "  run               - Executar localmente"
	@echo "  test              - Executar testes"
	@echo "  test-cover        - Testes com coverage"
	@echo "  test-bench        - Benchmark de performance"
	@echo "  test-all          - Todos os testes"
	@echo "  docker-build      - Build da imagem Docker"
	@echo "  docker-run        - Executar com Docker"
	@echo "  docker-compose-up - Executar com Docker Compose"
	@echo "  docker-compose-down - Parar Docker Compose"
	@echo "  clean             - Limpar arquivos de build"
	@echo "  help              - Mostrar esta ajuda"
