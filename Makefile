.PHONY: run build clean test install dev dev-up dev-down dev-logs docker-mongo wire

# Executar a aplicação
run:
	go run cmd/http/server.go cmd/http/wire_gen.go

# Compilar a aplicação
build:
	~/go/bin/wire ./cmd/http
	go build -o isayoga-api ./cmd/http

# Gerar wire
wire:
	~/go/bin/wire ./cmd/http

# Executar seed (carga inicial de dados) - Requer containers rodando
seed:
	docker-compose -f docker-compose.dev.yml exec -T api sh -c "cd /app && go run scripts/seed.go"

# Executar seed localmente (requer MongoDB local)
seed-local:
	go run scripts/seed.go

# Executar com hot reload local (requer air)
dev:
	@which air > /dev/null || (echo "Air não encontrado. Instale com: go install github.com/air-verse/air@latest" && exit 1)
	air

# Iniciar ambiente de desenvolvimento com Docker Compose
dev-up:
	docker-compose -f docker-compose.dev.yml up -d

# Parar ambiente de desenvolvimento
dev-down:
	docker-compose -f docker-compose.dev.yml down

# Ver logs do ambiente de desenvolvimento
dev-logs:
	docker-compose -f docker-compose.dev.yml logs -f api

# Rebuildar e reiniciar ambiente de desenvolvimento
dev-restart:
	docker-compose -f docker-compose.dev.yml down
	docker-compose -f docker-compose.dev.yml up -d --build

# Instalar dependências
install:
	go mod download
	go mod tidy

# Limpar arquivos compilados
clean:
	rm -f isayoga-api
	go clean

# Executar testes
test:
	go test -v ./...

# Executar testes com coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Iniciar MongoDB com Docker
docker-mongo:
	docker run -d -p 27017:27017 --name isayoga-mongo \
		-e MONGO_INITDB_DATABASE=isayoga \
		mongo:latest

# Parar MongoDB Docker
docker-mongo-stop:
	docker stop isayoga-mongo
	docker rm isayoga-mongo

# Verificar formatação
fmt:
	go fmt ./...

# Verificar código com linter
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint não encontrado. Instale com: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo ""
	@echo "Desenvolvimento:"
	@echo "  make dev             - Executar com hot reload local (requer air)"
	@echo "  make dev-up          - Iniciar ambiente de desenvolvimento (Docker)"
	@echo "  make dev-down        - Parar ambiente de desenvolvimento"
	@echo "  make dev-logs        - Ver logs do ambiente de desenvolvimento"
	@echo "  make dev-restart     - Rebuildar e reiniciar ambiente dev"
	@echo ""
	@echo "Aplicação:"
	@echo "  make run             - Executar a aplicação"
	@echo "  make build           - Compilar a aplicação"
	@echo "  make install         - Instalar dependências"
	@echo "  make clean           - Limpar arquivos compilados"
	@echo ""
	@echo "Testes:"
	@echo "  make test            - Executar testes"
	@echo "  make test-coverage   - Executar testes com coverage"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-mongo    - Iniciar MongoDB com Docker"
	@echo "  make docker-mongo-stop - Parar MongoDB Docker"
	@echo ""
	@echo "Qualidade de Código:"
	@echo "  make fmt             - Formatar código"
	@echo "  make lint            - Verificar código com linter"

