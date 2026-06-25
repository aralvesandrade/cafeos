# Fase 3 — Mobile Offline + Sincronização com Fila

## Objetivo

App mobile React Native com operação offline. Lançamentos em campo salvos localmente
e sincronizados via fila RabbitMQ quando houver conexão. Worker dedicado processa
os registros em lote no backend.

---

## 1. Infraestrutura — RabbitMQ

### Docker Compose

Add rabbitmq ao infra/dev/docker-compose.yml:

```yaml
rabbitmq:
  image: rabbitmq:4-management-alpine
  ports:
    - "5672:5672"   # AMQP
    - "15672:15672" # Management UI
  environment:
    RABBITMQ_DEFAULT_USER: cafeos
    RABBITMQ_DEFAULT_PASS: cafeos
  healthcheck:
    test: ["CMD", "rabbitmq-diagnostics", "check_port_connectivity"]
    interval: 10s
    timeout: 5s
    retries: 5
```

### Filas

| Fila | Propósito | Routing Key |
|------|-----------|-------------|
| sync.operations | Operações offline | sync.operation.* |
| sync.stock | Movimentações estoque | sync.stock.* |
| sync.harvest | Produção de safra | sync.harvest.* |
| sync.financial | Transações financeiras | sync.financial.* |
| sync.labor | Apontamento mão de obra | sync.labor.* |

Dead-letter queue após 3 retentativas: sync.dlq

---

## 2. Backend — RabbitMQ Integration

### Config (infra/config/config.go)

Add: `RabbitMQURL` (default amqp://cafeos:cafeos@localhost:5672/)

### Pacote infra/messaging/

```
internal/infra/messaging/
├── connection.go    # Conexão AMQP com reconnect automático
├── publisher.go     # Publica eventos na fila
└── consumer.go      # Consome + delega ao handler
```

### Sync Handler (api/handler/sync_handler.go)

`POST /api/v1/{tenant_id}/sync`

Recebe lote:
```json
{
  "batch": [{
    "client_id": "uuid",
    "event_type": "operation.created",
    "payload": {...},
    "client_timestamp": "2026-06-25T10:00:00Z"
  }]
}
```

Resposta: `202 Accepted { "accepted": N, "errors": [] }`

### Sync Worker (cmd/worker/main.go)

Processo separado:
1. Conecta RabbitMQ
2. Consome da fila sync.operations  
3. Reusa services existentes (OperationService, etc.)
4. Persiste no PostgreSQL
5. Publica evento no event bus
6. Ack a mensagem
7. Log + métricas

---

## 3. Mobile App — React Native

### Stack

| Lib | Uso |
|-----|-----|
| expo | Framework RN |
| expo-sqlite | Banco local |
| @react-native-community/netinfo | Conectividade |
| @react-navigation/native | Navegação |
| expo-secure-store | JWT storage |
| react-native-background-fetch | Sync background |

### Sync Engine

```
1. App inicia → NetInfo detecta conexão
2. Se online: busca sync_queue pendentes, envia lote p/ POST /sync
3. Se 202: marca synced. Se erro: retry, após 3 → failed
4. Se offline: salva local + enfileira p/ sync
5. NetInfo listener agenda sync quando reconectar
```

### Telas

| Tela | Offline |
|------|---------|
| Login | Precisa rede (JWT) |
| Operações | CRUD completo offline |
| Talhões | Cache local |
| Safras | CRUD offline |
| Estoque | CRUD offline |
| Financeiro | CRUD offline |
| Pendências | Leitura sync_queue local |

---

## 4. Entrega em Pacotes

### Pacote 1 — Infra
- RabbitMQ docker-compose + config

### Pacote 2 — Backend Messaging
- infra/messaging/ (connection, publisher, consumer)
- Sync handler + worker

### Pacote 3 — Mobile Base
- Scaffold Expo
- SQLite local + migrations
- Login + JWT
- NetInfo + sync engine
- Tela operações (CRUD offline)
- Tela pendências

### Pacote 4 — Mobile Demais Telas
- Estoque, Financeiro, Safras, Dashboard

---

## 5. Observações

- **Conflitos**: client_timestamp + updated_at — last-write-wins
- **Idempotência**: client_id único; worker ignora duplicatas
- **Volume**: RabbitMQ escala horizontalmente
- **Segurança**: JWT em toda requisição
- **Background sync**: react-native-background-fetch p/ sincronizar com app fechado
