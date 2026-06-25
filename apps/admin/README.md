# CafeOS — Admin Panel

Painel administrativo do CafeOS, plataforma SaaS multi-tenant para cafeicultura.

## Stack

- **Vite** + **React** + **TypeScript**
- **Tailwind CSS v4** com paleta personalizada (coffee-green, coffee-brown, coffee-beige)
- **React Router DOM v7** (nested layouts)
- **Recharts** (gráficos)
- **Lucide React** (ícones)
- **shadcn/ui** (Button, Badge, Card, Table, Dialog, Input)

## Estrutura

```
src/
├── components/
│   ├── layout/       # Sidebar, Header, AdminLayout, AuthLayout
│   ├── ui/           # Button, Badge, Card, Table, Dialog, Input, Select
│   ├── dashboard/    # StatsCards, ProductionChart, CostChart, RecentOperations
│   └── farms/        # FarmList, FarmForm
├── lib/
│   ├── utils.ts      # cn() utility (class-variance-authority + tailwind-merge)
│   ├── api.ts        # Cliente HTTP com JWT + suporte a rotas admin
│   └── auth.tsx      # AuthContext + hook
├── pages/            # Login, Dashboard, Farms, Plots, Operations, Harvests, Tenants, Users, NotFound
├── router.tsx        # React Router DOM v7 (nested layouts, role-based guards)
├── App.tsx
├── main.tsx
└── index.css         # Tailwind v4 @theme
```

## Perfis RBAC

| Perfil | Visibilidade |
|--------|-------------|
| `platform_owner` | Dashboard + CRUD fazendas/talhões/operações/safras + Administração (Tenants, Usuários) |
| Demais perfis | Dashboard + dados do próprio tenant apenas |

## Comandos

```bash
npm run dev      # Servidor de desenvolvimento (Vite, http://localhost:5174)
npm run build    # Build de produção
npm run preview  # Preview do build
```

## Login Rápido (dev)

Botões de acesso rápido na tela de login preenchem email/senha automaticamente para cada perfil RBAC.
