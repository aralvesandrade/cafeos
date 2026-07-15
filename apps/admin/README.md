# CafeOS — Admin Panel

Painel administrativo do CafeOS, plataforma SaaS multi-tenant para cafeicultura.

## Stack

- **Vite** + **React** + **TypeScript**
- **Tailwind CSS v4** com tokens semânticos do CafeOS Design System (shadcn-style, verde `primary` + neutros zinc), suporte a light/dark
- **React Router DOM v7** (nested layouts)
- **Recharts** (gráficos)
- **Lucide React** (ícones)
- **shadcn/ui** (Button, Badge, Card, Table, Dialog, Input, Select)

## Estrutura

```
src/
├── components/
│   ├── layout/       # Sidebar (theme toggle, colapsável), Header (NotificationBell), AdminLayout, AuthLayout (theme toggle)
│   ├── ui/           # Button, Badge, Card, Table, Dialog, Input, Select
│   ├── dashboard/    # StatsCards, ProductionChart, CostChart, RecentOperations
│   └── farms/        # FarmList, FarmForm
├── lib/
│   ├── utils.ts      # cn() utility (class-variance-authority + tailwind-merge)
│   ├── api.ts        # Cliente HTTP com JWT + suporte a rotas admin
│   ├── auth.tsx       # AuthContext + hook
│   └── theme.tsx      # ThemeProvider + useTheme (light/dark, persistido em localStorage)
├── pages/            # Login, Dashboard, Farms, FarmDetail, FarmEdit, Plots, PlotDetail, PlotEdit,
│                     # Operations, OperationDetail, OperationTypes, Harvests, HarvestDetail,
│                     # Financial, CostCenters, Budget, CostAllocations, Stock, Fleet, Labor,
│                     # Organizations, Users, NotFound
├── router.tsx        # React Router DOM v7 (nested layouts, role-based guards)
├── App.tsx
├── main.tsx          # aplica classe .dark antes do primeiro render (evita flash)
└── index.css         # Tokens CSS (:root / .dark) + Tailwind v4 @theme inline
```

## Tema (light/dark)

Tokens semânticos (`background`, `foreground`, `card`, `primary`, `sidebar`, etc.) definidos em `src/index.css` como variáveis CSS, com overrides em `.dark`. `ThemeProvider` (`src/lib/theme.tsx`) alterna a classe `.dark` no `<html>` e persiste a escolha em `localStorage` (`cafeos_theme`); sem preferência salva, segue `prefers-color-scheme`. Toggle disponível no rodapé do Sidebar (área logada) e na tela de login (`AuthLayout`).

## Sidebar colapsável

Botão no rodapé do Sidebar alterna entre modo expandido (ícone + label) e
colapsado (só ícones, com tooltip no hover) — desktop apenas (`lg:`), a
gaveta mobile sempre mostra os labels. Preferência persistida em
`localStorage` (`cafeos-sidebar-collapsed`).

## Notificação de alertas

`NotificationBell` (`components/layout/`) no Header faz polling de
`GET /alerts` a cada 60s e mostra um badge com a contagem de alertas em
aberto. Cada alerta pode ser resolvido ou descartado direto no dropdown
(`PUT /alerts/{id}`).

## Perfis RBAC

| Perfil | Visibilidade |
|--------|-------------|
| `platform_owner` | Dashboard + CRUD fazendas/talhões/operações/safras + Administração (Organizações, Usuários) |
| Demais perfis | Dashboard + dados da própria organização apenas |

## Comandos

```bash
npm run dev      # Servidor de desenvolvimento (Vite, http://localhost:5174)
npm run build    # Build de produção
npm run preview  # Preview do build
```

## Login Rápido (dev)

Botões de acesso rápido na tela de login preenchem email/senha automaticamente para cada perfil RBAC.
