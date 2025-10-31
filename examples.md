# Exemplos de Uso da API

## 1. Criar Usuário Estudante
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Santos",
    "email": "joao@email.com",
    "role": "student"
  }'
```

## 2. Criar Usuário Instrutor
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Maria Silva",
    "email": "maria@yoga.com",
    "role": "instructor"
  }'
```

**Roles disponíveis:**
- `student` (padrão se não especificado)
- `instructor`
- `admin`

## 3. Criar Aula
```bash
curl -X POST http://localhost:8080/api/v1/classes \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hatha Yoga - Manhã",
    "description": "Aula para iniciantes focada em posturas básicas",
    "instructor_id": "INSTRUCTOR_ID_AQUI",
    "instructor_name": "Maria Silva",
    "start_time": "2025-11-15T08:00:00Z",
    "end_time": "2025-11-15T09:30:00Z",
    "max_capacity": 10,
    "price_in_cents": 5000
  }'
```

## 4. Listar Aulas
```bash
curl http://localhost:8080/api/v1/classes
```

## 5. Inscrever Estudante em Aula
```bash
curl -X POST http://localhost:8080/api/v1/enrollments \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "USER_ID_AQUI",
    "class_id": "CLASS_ID_AQUI"
  }'
```

**Resposta inclui:**
- `enrollment`: Dados da inscrição
- `payment`: Dados do pagamento
- `payment_url`: URL do Mercado Pago para finalizar o pagamento

**Importante:** Apenas usuários com `role: "student"` podem se inscrever em aulas.

## 6. Cancelar Inscrição
```bash
curl -X DELETE http://localhost:8080/api/v1/enrollments/ENROLLMENT_ID_AQUI
```

## 7. Webhook Mercado Pago (Automático)
O Mercado Pago enviará notificações para:
```
POST http://seu-dominio.com/webhooks/mercadopago
```

Exemplo de payload:
```json
{
  "action": "payment.created",
  "type": "payment",
  "data.id": "1234567890",
  "live_mode": true
}
```

## Controle de Concorrência

A API garante que:
1. Múltiplos usuários podem tentar se inscrever simultaneamente
2. Apenas quem conseguir reservar a vaga dentro da transação terá sucesso
3. O versionamento otimista previne double-booking
4. Transações MongoDB garantem atomicidade

### Exemplo de Concorrência:
- Aula com 1 vaga disponível
- 3 usuários tentam se inscrever ao mesmo tempo
- Apenas 1 conseguirá, os outros receberão erro "aula sem vagas"

