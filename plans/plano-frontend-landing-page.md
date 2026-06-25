# Plano: Landing Page CafeOS

## Stack

| Tecnologia | Versão |
|---|---|
| React | 19 |
| Vite | 6 |
| TypeScript | 5 |
| Tailwind CSS | 4 |
| shadcn/ui | latest |
| Lucide React | 0.400+ |
| React Router DOM | 7 |

## Estrutura (`apps/frontend/`)

```
apps/frontend/
├── public/
│   └── logo.svg
├── src/
│   ├── components/
│   │   ├── ui/            # shadcn/ui components
│   │   ├── layout/
│   │   │   ├── Header.tsx
│   │   │   └── Footer.tsx
│   │   └── sections/
│   │       ├── Hero.tsx
│   │       ├── About.tsx
│   │       ├── Features.tsx
│   │       ├── CoffeeCycle.tsx
│   │       ├── Indicators.tsx
│   │       ├── Plans.tsx
│   │       ├── TechStack.tsx
│   │       ├── Roadmap.tsx
│   │       └── CtaSection.tsx
│   ├── lib/
│   │   └── utils.ts
│   ├── App.tsx
│   ├── main.tsx
│   └── index.css
├── index.html
├── tailwind.config.ts
├── tsconfig.json
├── vite.config.ts
└── package.json
```

## Seções da Landing Page

| # | Seção | Conteúdo |
|---|-------|----------|
| 1 | **Hero** | Título, subtítulo, CTA "Começar grátis", imagem ilustrativa |
| 2 | **Sobre** | Desafios da cafeicultura e como o CafeOS resolve |
| 3 | **Funcionalidades** | Cards: Fazendas, Talhões, Operações, Safras, Custos, Dashboard |
| 4 | **Ciclo do Café** | Timeline visual das fases fenológicas |
| 5 | **Indicadores** | Sacas/ha, Custo/saca, Rentabilidade, Bienalidade |
| 6 | **Planos** | Grid de 4 planos (Grátis, Pro, Cooperativa, Consultoria) |
| 7 | **Tecnologia** | Ícones da stack: Go, React, PostgreSQL, Docker, IoT, IA |
| 8 | **Roadmap** | Timeline das fases (MVP, Fase 2, 3, 4) |
| 9 | **CTA Final** | Chamada final com formulário de interesse |
| 10 | **Footer** | Redes, contato, links |

## Planos de Assinatura

| Feature | Grátis | Pro | Cooperativa | Consultoria |
|---------|--------|-----|-------------|-------------|
| Fazendas | 1 | 10 | Ilimitado | Multi-cliente |
| Talhões | 5 | 50 | Ilimitado | Ilimitado |
| Operações | 50/mês | Ilimitado | Ilimitado | Ilimitado |
| Safras | 3 | Ilimitado | Ilimitado | Ilimitado |
| Usuários | 2 | 10 | 50 | 30 |
| Dashboard | Básico | Avançado | Consolidado | Por cliente |
| Relatórios | — | CSV/PDF | Benchmarking | Técnicos |
| White Label | — | — | — | ✅ |
| Suporte | Comunidade | Email | Prioridade | Dedicado |
| **Preço** | **Grátis** | **R$ 97/mês** | **R$ 297/mês** | **R$ 497/mês** |

## Identidade Visual

- **Paleta:** `#2E7D32` (verde café), `#795548` (marrom), `#FAF5F0` (bege), `#1B5E20` (verde escuro), `#FFF` (branco)
- **Fonte:** Inter (sans-serif)
- **Ícones:** Lucide React
- **Ilustrações:** SVG placeholders

## Passos de Implementação

1. Criar Vite project com template react-ts
2. Adicionar Tailwind CSS v4 com Vite plugin
3. Inicializar shadcn/ui
4. Configurar globals.css com paleta de cores
5. Criar componentes de layout (Header, Footer)
6. Implementar seções na ordem listada
7. Responsivo mobile-first
8. Build e verificação
