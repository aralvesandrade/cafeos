# Plano: Admin Panel CafeOS

## Stack

| Tecnologia       | Versão           | Nota                                   |
| ---------------- | ---------------- | -------------------------------------- |
| React            | 19               |                                        |
| Vite             | 8                |                                        |
| TypeScript       | 6                |                                        |
| Tailwind CSS     | 4                | Mesma paleta do frontend               |
| shadcn/ui        | latest           | Sidebar, Table, Dialog, Form, etc.     |
| Lucide React     | 1.x              |                                        |
| React Router DOM | 7                | Nested layouts, lazy routes            |
| Recharts         | 2.x              | Gráficos do dashboard                  |

## Estrutura (`apps/admin/`)

```
apps/admin/
├── public/
│   └── favicon.svg
├── src/
│   ├── components/
│   │   ├── ui/                # shadcn/ui (Button, Badge, Card, Table, Dialog, Input, Select, etc.)
│   │   ├── layout/
│   │   │   ├── Sidebar.tsx       # Navegação lateral (collapsible)
│   │   │   ├── Header.tsx        # Topo com avatar, notificações, logout
│   │   │   ├── AdminLayout.tsx   # Sidebar + Header + <Outlet />
│   │   │   └── AuthLayout.tsx    # Layout limpo para login
│   │   ├── farms/
│   │   │   ├── FarmList.tsx
│   │   │   ├── FarmForm.tsx
│   │   │   └── FarmDetail.tsx
│   │   ├── plots/
│   │   │   ├── PlotList.tsx
│   │   │   ├── PlotForm.tsx
│   │   │   └── PlotDetail.tsx
│   │   ├── operations/
│   │   │   ├── OperationList.tsx
│   │   │   └── OperationDetail.tsx
│   │   ├── harvests/
│   │   │   ├── HarvestList.tsx
│   │   │   ├── HarvestForm.tsx
│   │   │   └── HarvestDetail.tsx
│   │   ├── tenants/
│   │   │   ├── TenantList.tsx
│   │   │   └── TenantForm.tsx
│   │   ├── users/
│   │   │   ├── UserList.tsx
│   │   │   └── UserForm.tsx
│   │   └── dashboard/
│   │       ├── StatsCards.tsx
│   │       ├── ProductionChart.tsx
│   │       ├── CostChart.tsx
│   │       └── RecentOperations.tsx
│   ├── lib/
│   │   ├── utils.ts             # cn() utility
│   │   ├── api.ts               # Cliente HTTP (fetch wrapper)
│   │   └── auth.ts              # Context + hook de autenticação
│   ├── pages/
│   │   ├── Login.tsx
│   │   ├── Dashboard.tsx
│   │   ├── Farms.tsx
│   │   ├── FarmDetail.tsx
│   │   ├── Plots.tsx
│   │   ├── PlotDetail.tsx
│   │   ├── Operations.tsx
│   │   ├── Harvests.tsx
│   │   ├── HarvestDetail.tsx
│   │   ├── Tenants.tsx
│   │   ├── Users.tsx
│   │   └── NotFound.tsx
│   ├── router.tsx               # React Router config
│   ├── App.tsx                  # Provider wrapper
│   ├── main.tsx                 # Entry point
│   └── index.css                # Tailwind v4 @theme (mesmo do frontend)
├── index.html
├── tsconfig.json
├── tsconfig.app.json
├── tsconfig.node.json
├── vite.config.ts
└── package.json
```

## Páginas

| #  | Rota            | Página           | Descrição                                    |
| -- | --------------- | ---------------- | -------------------------------------------- |
| 1  | `/login`        | Login            | Formulário de autenticação JWT               |
| 2  | `/`             | Dashboard        | Cards + gráficos de produção, custos, safra  |
| 3  | `/farms`        | Farms            | Lista de fazendas com CRUD                   |
| 4  | `/farms/:id`    | FarmDetail       | Detalhe da fazenda + talhões vinculados      |
| 5  | `/plots`        | Plots            | Lista de talhões com CRUD                    |
| 6  | `/plots/:id`    | PlotDetail       | Detalhe do talhão + operações vinculadas     |
| 7  | `/operations`   | Operations       | Lista de operações com filtros               |
| 8  | `/harvests`     | Harvests         | Lista de safras com CRUD                     |
| 9  | `/harvests/:id` | HarvestDetail    | Detalhe da safra + produção por talhão       |
| 10 | `/tenants`      | Tenants          | [Platform Owner] CRUD de tenants             |
| 11 | `/users`        | Users            | [Platform Owner / Tenant Admin] CRUD de users |

## Rotas da API consumidas

### Farms
| Método | Rota                                     |
| ------ | ---------------------------------------- |
| GET    | `/api/v1/{tenant_id}/farms`              |
| POST   | `/api/v1/{tenant_id}/farms`              |
| GET    | `/api/v1/{tenant_id}/farms/{id}`         |
| PUT    | `/api/v1/{tenant_id}/farms/{id}`         |
| DELETE | `/api/v1/{tenant_id}/farms/{id}`         |

### Plots
| Método | Rota                                        |
| ------ | ------------------------------------------- |
| GET    | `/api/v1/{tenant_id}/plots`                 |
| POST   | `/api/v1/{tenant_id}/plots`                 |
| GET    | `/api/v1/{tenant_id}/plots/{id}`            |
| PUT    | `/api/v1/{tenant_id}/plots/{id}`            |
| DELETE | `/api/v1/{tenant_id}/plots/{id}`            |
| GET    | `/api/v1/{tenant_id}/farms/{farm_id}/plots` |

### Operations
| Método | Rota                                           |
| ------ | ---------------------------------------------- |
| GET    | `/api/v1/{tenant_id}/operations`               |
| POST   | `/api/v1/{tenant_id}/operations`               |
| GET    | `/api/v1/{tenant_id}/operations/recent`        |
| GET    | `/api/v1/{tenant_id}/operations/{id}`          |
| DELETE | `/api/v1/{tenant_id}/operations/{id}`          |
| GET    | `/api/v1/{tenant_id}/plots/{plot_id}/operations` |

### Harvests
| Método | Rota                                             |
| ------ | ------------------------------------------------ |
| GET    | `/api/v1/{tenant_id}/harvests`                   |
| POST   | `/api/v1/{tenant_id}/harvests`                   |
| GET    | `/api/v1/{tenant_id}/harvests/{id}`              |
| PUT    | `/api/v1/{tenant_id}/harvests/{id}/finalize`     |
| POST   | `/api/v1/{tenant_id}/harvests/{id}/production`   |
| GET    | `/api/v1/{tenant_id}/harvests/{id}/production`   |

### Dashboard
| Método | Rota                                 |
| ------ | ------------------------------------ |
| GET    | `/api/v1/{tenant_id}/dashboard`      |

## Sidebar — Itens de Navegação

```
☰ CafeOS
├── 📊 Dashboard
├── 🏠 Fazendas
├── 🗺️ Talhões
├── 🚜 Operações
├── 🌾 Safras
│
└── ⚙️ Administração (platform_owner / tenant_admin)
    ├── 🏢 Tenants
    └── 👥 Usuários
```

## Fluxo de Autenticação

1. Usuário acessa `/login` → form com email + senha
2. POST `/api/v1/auth/login` → retorna JWT + user info (endpoint a ser criado no backend)
3. Token armazenado em `localStorage` / contexto React
4. `AuthContext` expõe `{ user, token, login, logout, isAuthenticated }`
5. Rotas privadas protegidas por `<PrivateRoute>` que redireciona para `/login`
6. Header mostra avatar + nome do usuário + botão de logout

## Identidade Visual

- **Design system:** CafeOS Design System (shadcn-style) — tokens semânticos (`--background`, `--foreground`, `--card`, `--primary`, `--muted`, `--sidebar`, etc.) em `src/index.css`, com variantes light e dark
- **Cor de marca:** verde `primary` (`#15803d` light / `#22c55e` dark) sobre base neutra zinc
- **Sidebar:** fundo `--sidebar`, item ativo `--sidebar-active-bg`/`--sidebar-active-foreground`
- **Header:** fundo `--background` com borda inferior; inclui toggle de tema (claro/escuro)
- **Tema:** `ThemeProvider` (`src/lib/theme.tsx`) persiste escolha em `localStorage`, com fallback a `prefers-color-scheme`; toggle também presente na tela de login
- **Fonte:** Geist (interface) / Geist Mono (valores e código)
- **Ícones:** Lucide React

## Passos de Implementação

1. Criar Vite project com template react-ts em `apps/admin/`
2. Instalar dependências: tailwindcss, shadcn/ui, lucide-react, react-router-dom, recharts
3. Configurar Tailwind v4 (copiar `@theme` do frontend), `vite.config.ts` com alias `@/`, `tsconfig.app.json`
4. Criar `lib/utils.ts` (cn), `lib/api.ts` (fetch wrapper), `lib/auth.ts` (AuthContext + hook)
5. Criar `router.tsx` com React Router DOM v7 (nested routes)
6. Criar componentes de layout: `AdminLayout`, `AuthLayout`, `Sidebar`, `Header`
7. Criar página `Login.tsx` com formulário
8. Criar página `Dashboard.tsx` com cards + gráficos (Recharts)
9. Criar CRUD pages: Farms, Plots, Operations, Harvests, Tenants, Users
10. Criar componentes de formulário modais (`FarmForm`, `PlotForm`, etc.)
11. Responsivo (sidebar colapsa em mobile)
12. Build e verificação
