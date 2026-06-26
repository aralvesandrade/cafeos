# Plano de Migração: Expo (React Native) → Flutter

## Objetivo

Substituir o app mobile atual (`apps/mobile/`, Expo/React Native) por um novo app Flutter no diretório `apps/mobile_flutter/`. Primeiro release mantém o mesmo escopo do app atual (Login, Operações, Pendências). Releases futuros adicionam Estoque, Financeiro, Safras, Frotas, Mão de Obra.

## Justificativa

Equipe possui mais experiência e produtividade com Dart/Flutter. A migração aproveita melhor as habilidades do time sem alterar contratos de API ou backend.

## Stack Proposto

| Camada | Expo (atual) | Flutter (novo) |
|--------|-------------|----------------|
| Linguagem | TypeScript | Dart 3 |
| Framework | Expo SDK 56 | Flutter 3.x |
| SQLite local | expo-sqlite | drift |
| Storage seguro | expo-secure-store | flutter_secure_storage |
| HTTP | fetch nativo | dio |
| Conectividade | @react-native-community/netinfo | connectivity_plus |
| Navegação | @react-navigation/bottom-tabs | go_router |
| Estado | useState/useEffect | riverpod |
| UUID | crypto.randomUUID() | pacote uuid |
| Background sync | react-native-background-fetch | workmanager (v2) |
| Testes | Jest | flutter_test + mocktail |

> **Decisão**: riverpod escolhido para estado por ser mais moderno e leve que bloc.

## Estrutura de Diretórios

```
apps/mobile_flutter/
├── android/
├── ios/
├── lib/
│   ├── main.dart
│   ├── app.dart
│   ├── api/
│   │   ├── client.dart
│   │   └── storage.dart
│   ├── db/
│   │   ├── database.dart
│   │   ├── schema.dart
│   │   └── migrations.dart
│   ├── models/
│   │   ├── operation.dart
│   │   ├── plot.dart
│   │   ├── farm.dart
│   │   └── sync_item.dart
│   ├── repos/
│   │   ├── operation_repo.dart
│   │   ├── plot_repo.dart
│   │   └── sync_queue_repo.dart
│   ├── services/
│   │   ├── auth_service.dart
│   │   ├── sync_service.dart
│   │   └── connectivity_service.dart
│   ├── screens/
│   │   ├── login/
│   │   │   ├── login_screen.dart
│   │   │   └── login_controller.dart
│   │   ├── home/
│   │   │   └── home_screen.dart
│   │   ├── operations/
│   │   │   ├── operations_screen.dart
│   │   │   ├── operation_form.dart
│   │   │   └── operations_controller.dart
│   │   └── pending_sync/
│   │       ├── pending_sync_screen.dart
│   │       └── pending_sync_controller.dart
│   ├── router/
│   │   └── app_router.dart
│   └── shared/
│       ├── widgets/
│       │   ├── status_badge.dart
│       │   └── sync_button.dart
│       └── theme/
│           └── app_theme.dart
├── test/
│   ├── unit/
│   │   ├── services/sync_service_test.dart
│   │   └── repos/operation_repo_test.dart
│   └── widget/
│       ├── screens/login_screen_test.dart
│       └── screens/operations_screen_test.dart
├── pubspec.yaml
├── analysis_options.yaml
└── .env.example
```

## Arquitetura

### Offline-First (mesmo padrão do app atual)

```
User action -> Local DB (drift) + Enqueue sync (sync_queue)
     ^                    |
     +-- Online ------> POST /api/v1/{tenant_id}/sync
     |                    |
     |                 202 Accepted -> mark synced
     |                    |
     |                 RabbitMQ -> Worker -> PostgreSQL
     |
     +-- Offline ------> Salva local, pendente na sync_queue
                          |
                         connectivity_plus detecta reconexao
                          |
                         syncAll() processa pendentes
```

### Fluxo de Autenticação (idêntico ao atual)

```
App start -> flutter_secure_storage.read('cafeos_token')
  +-- Token existe -> home screen
  +-- Token vazio -> login screen
                      |
                   POST /auth/login { email, password }
                      |
                   { token, tenant_id, user }
                      |
                   storage.save('cafeos_token', token)
                   storage.save('cafeos_tenant_id', tenant_id)
                      |
                   home screen
```

### Contrato de Sincronização (INALTERADO)

```
POST /api/v1/{tenant_id}/sync
Authorization: Bearer {token}
Content-Type: application/json

{
  "batch": [{
    "client_id": "uuid-v4",
    "event_type": "operation.created",
    "payload": {
      "plot_id": "uuid",
      "type": "adubacao",
      "date": "2026-06-25T10:00:00Z",
      "responsible": "Joao",
      "product_used": "NPK 20-05-20",
      "quantity": 50.0,
      "cost": 2500.0,
      "notes": "Aplicacao manual"
    },
    "client_timestamp": "2026-06-25T10:00:00Z"
  }]
}

Response: 202 Accepted
{ "accepted": 1, "errors": [] }
```

## Telas (v1 — Mesmo escopo do app atual)

### 1. Login Screen
- Campos: email, senha
- Botao "Entrar" com loading
- Pre-preenchido com joao@cafeos.com.br / 123456 (dev)
- Validacao: campos obrigatorios
- Erro exibido em SnackBar
- Sucesso -> navega para home

### 2. Home Screen (Dashboard placeholder)
- Titulo "CafeOS"
- Subtitulo "App offline para operacoes de campo"
- Bottom nav com 3 abas: Inicio, Operacoes, Pendencias

### 3. Operations Screen
- Lista de operacoes ordenada por data DESC
- Cada item: tipo (color badge), data, responsavel, custo, badge "Pendente"
- FAB "+" para nova operacao
- Empty state + Loading state + Pull-to-refresh

### 4. Operation Form (Modal/BottomSheet)
- Tipo: chips (Adubacao, Pulverizacao, Irrigacao, Poda, Colheita)
- Data: date picker (YYYY-MM-DD)
- Responsavel: text field
- Quantidade: decimal input
- Custo: decimal input com prefixo R$
- Notas: text field multiline
- Botoes: Cancelar + Salvar

### 5. Pending Sync Screen
- Lista sync_queue com status badge (pending/amber, synced/green, failed/red)
- Botao "Sincronizar" -> syncAll() -> reload
- Empty state + Loading state

## SQLite Local (drift) — Esquema

```sql
CREATE TABLE sync_queue (
  id TEXT NOT NULL PRIMARY KEY,
  event_type TEXT NOT NULL,
  payload TEXT NOT NULL,
  client_timestamp TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  retry_count INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE operations (
  id TEXT NOT NULL PRIMARY KEY,
  plot_id TEXT NOT NULL,
  type TEXT NOT NULL,
  date TEXT NOT NULL,
  responsible TEXT NOT NULL DEFAULT '',
  product_used TEXT NOT NULL DEFAULT '',
  quantity REAL NOT NULL DEFAULT 0,
  cost REAL NOT NULL DEFAULT 0,
  notes TEXT NOT NULL DEFAULT '',
  synced INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE plots (
  id TEXT NOT NULL PRIMARY KEY,
  farm_id TEXT NOT NULL,
  name TEXT NOT NULL,
  area_ha REAL NOT NULL DEFAULT 0,
  cultivar TEXT NOT NULL DEFAULT '',
  synced INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE farms (
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  owner TEXT NOT NULL DEFAULT '',
  location TEXT NOT NULL DEFAULT '',
  synced INTEGER NOT NULL DEFAULT 0
);
```

## Pacotes Flutter (pubspec.yaml)

```yaml
dependencies:
  flutter:
    sdk: flutter
  dio: ^5.x
  flutter_secure_storage: ^9.x
  drift: ^2.x
  sqlite3_flutter_libs: ^0.x
  path: ^1.x
  path_provider: ^2.x
  connectivity_plus: ^6.x
  go_router: ^14.x
  flutter_riverpod: ^2.x
  riverpod_annotation: ^2.x
  uuid: ^4.x
  workmanager: ^0.x    # v2
  intl: ^0.x

dev_dependencies:
  flutter_test:
    sdk: flutter
  drift_dev: ^2.x
  build_runner: ^2.x
  riverpod_generator: ^2.x
  mocktail: ^1.x
  flutter_lints: ^5.x
```

## Fases de Implementacao

### Fase 1 — Scaffold + Base
- `flutter create apps/mobile_flutter`
- pubspec.yaml com dependencias
- Drift database setup com schema (4 tabelas)
- flutter_secure_storage wrapper
- Dio client com interceptors (auth header, tenant prefix)
- go_router com 3 tabs (placeholder screens)
- App theme (cores CafeOS: #2E7D32, #F5F0EB)
- Config API_BASE_URL

### Fase 2 — Autenticacao
- LoginScreen completo
- AuthService: login(), getToken(), setToken(), clearToken()
- Storage: cafeos_token, cafeos_tenant_id
- App startup: check token -> home ou login
- Error handling (SnackBar)
- Loading states

### Fase 3 — Sync Engine
- ConnectivityService (connectivity_plus listener)
- SyncQueueRepo (CRUD sync_queue)
- SyncService: enqueue(), syncAll(), retry (max 3)
- Batch POST /sync com batch <= 50
- Auto-sync on reconnect

### Fase 4 — Tela de Operacoes
- PlotRepo (cache local de talhoes)
- OperationRepo (CRUD local)
- OperationsScreen + OperationForm
- createOperation: enqueue + local insert + reload
- Color-coded type badges
- Sync status badges
- Pull-to-refresh

### Fase 5 — Tela de Pendencias
- PendingSyncScreen
- Visualizar sync_queue
- Botao "Sincronizar"
- Status badges: pending/amber, synced/green, failed/red

### Fase 6 — Finalizacao
- Testes unitarios (sync_service, operation_repo)
- Testes widget (login_screen, operations_screen)
- Testar fluxo completo offline -> online
- Remover apps/mobile/ (apos validacao)

## Observacoes

1. **Contratos de API inalterados**: backend nao precisa de mudancas. Flutter consome mesmas rotas, headers e payloads.
2. **Mesma semantica offline**: event sourcing identico. sync_queue como tabela de eventos, synced flag nas tabelas de dados.
3. **Idempotencia**: client_id UUID v4 gerado no cliente. Worker ignora duplicatas.
4. **Conflitos**: client_timestamp + updated_at -- last-write-wins (mesmo do plano original).
5. **Background sync**: postergado para v2. workmanager adicionado nas dependencias mas nao implementado na v1.
6. **Retencao do app antigo**: apps/mobile/ mantido ate validacao completa do Flutter em producao.

## Event Types Suportados (v1)

| Event Type | Tela | Release |
|-----------|------|---------|
| operation.created | Operacoes | v1 |
| operation.updated | Operacoes | v1 |
| stock.created, stock.moved | Estoque | v2 |
| financial.created | Financeiro | v2 |
| harvest.production | Safras | v2 |
| labor.shift | Mao de Obra | v3 |
