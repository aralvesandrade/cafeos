# CafeOS — Frontend

Landing page do CafeOS, plataforma SaaS especialista em cafeicultura.

## Stack

- **Vite** + **React** + **TypeScript**
- **Tailwind CSS v4** com paleta personalizada (coffee-green, coffee-brown, coffee-beige)
- **shadcn/ui** (Button, Badge)
- **Lucide React** (ícones)

## Estrutura

```
src/
├── components/
│   ├── layout/       # Header, Footer
│   ├── sections/     # Hero, About, Features, CoffeeCycle, Indicators, Plans, CtaSection
│   └── ui/           # Button, Badge (shadcn/ui)
├── lib/
│   └── utils.ts      # cn() utility (class-variance-authority + tailwind-merge)
├── App.tsx           # Página principal com todas as seções
├── main.tsx          # Entry point
└── index.css         # Tailwind v4 @theme + @import "tailwindcss"
```

## Comandos

```bash
npm run dev      # Servidor de desenvolvimento (Vite)
npm run build    # Build de produção
npm run preview  # Preview do build
```

## Paleta de Cores

| Token             | Cor     | Uso                          |
| ----------------- | ------- | ---------------------------- |
| `--color-coffee-green`  | #2E7D32 | Ações, destaque              |
| `--color-coffee-brown`  | #795548 | Títulos, navbar              |
| `--color-coffee-beige`  | #FAF5F0 | Background de seções         |
```

