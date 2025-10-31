# IsaYoga API

API REST para gerenciamento de inscrições em aulas de yoga com pagamentos via Mercado Pago.

## Stack
- Go 1.25.3
- go-chi (router)
- MongoDB (com transações)
- Zap (logs estruturados)
- OpenTelemetry + Zipkin (tracing distribuído)
- godotenv
- Google Wire (DI)
- Mercado Pago SDK
- Docker & Docker Compose
- Air (hot-reload)

## Funcionalidades
- 🧘 Gerenciamento de aulas com vagas limitadas
- 🔒 Controle de concorrência otimista (versioning)
- 💳 Integração com Mercado Pago para pagamentos
- 🔄 Processamento de webhooks
- 📦 Clean Architecture
- 🔌 Dependency Injection com Wire

## Estrutura
```
api/
├── cmd/http/                  # Aplicação principal + Wire
├── internal/
│   ├── domain/               # Entidades e interfaces
│   ├── usecase/              # Regras de negócio
│   ├── infrastructure/       # Database, HTTP, Payment
│   └── interface/            # Handlers HTTP
└── pkg/                      # Config e Logger
```

## Instalação
```bash
cp .env.example .env
# Configure MERCADOPAGO_ACCESS_TOKEN no .env
docker-compose -f docker-compose.dev.yml up
```

## Interfaces Web
- **API**: http://localhost:8080
- **Zipkin UI**: http://localhost:9411 (Tracing distribuído)
- **Health Check**: http://localhost:8080/health

## Comandos
```bash
make run        # Rodar aplicação
make build      # Build
make wire       # Gerar DI
```

## Endpoints

### Health
```
GET  /health
```

### Usuários
```
GET    /api/v1/users               # Listar usuários
POST   /api/v1/users               # Criar usuário (role: student, instructor, admin)
GET    /api/v1/users/{id}          # Obter usuário
PUT    /api/v1/users/{id}          # Atualizar usuário
DELETE /api/v1/users/{id}          # Deletar usuário
```

**Roles disponíveis:**
- `student` - Pode se inscrever em aulas
- `instructor` - Pode criar e ministrar aulas
- `admin` - Acesso total ao sistema

### Aulas
```
GET  /api/v1/classes          # Listar aulas
POST /api/v1/classes          # Criar aula
```

### Inscrições
```
POST   /api/v1/enrollments     # Inscrever aluno (retorna URL de pagamento)
DELETE /api/v1/enrollments/{id} # Cancelar inscrição
```

### Webhooks
```
POST /webhooks/mercadopago     # Webhook Mercado Pago
```

## Controle de Concorrência
A API utiliza versionamento otimista para garantir que múltiplos usuários não reservem a mesma vaga simultaneamente. Transações MongoDB garantem atomicidade das operações.
