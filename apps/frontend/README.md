# CafeOS — Frontend

Landing page do CafeOS, plataforma SaaS especialista em cafeicultura.

## Stack

- **Vite** + **React** + **TypeScript**
- **Tailwind CSS v4** com identidade visual própria "armazém/torrefação" (base
  escura `ink`, acentos `terreiro`/`gold`, tipografia `Fraunces` + `Inter` +
  `IBM Plex Mono`)
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

Identidade enraizada no mundo material do café (terreiro de secagem, cereja
madura, ficha de manifesto/pesagem) em vez do verde-safra genérico de SaaS.

| Token                  | Cor     | Uso                                   |
| ---------------------- | ------- | -------------------------------------- |
| `--color-ink`           | #171009 | Fundo de página (marrom-torrado)       |
| `--color-ink-raised`    | #221807 | Fundo de seção/card alternado          |
| `--color-parchment`     | #F1E6D2 | Texto principal, superfícies claras    |
| `--color-muted`         | #C9B99C | Texto secundário sobre fundo escuro    |
| `--color-terreiro`      | #C1552F | Acento primário (terracota/cereja)     |
| `--color-gold`          | #D2A44C | Acento secundário (grão torrado)       |
| `--color-leaf`          | #5C7A52 | Acento terciário, uso moderado         |
| `--color-rule`          | #3A2C18 | Divisores/hairlines                    |

Tipografia: `Fraunces` (display), `Inter` (corpo), `IBM Plex Mono` (dados/
indicadores — usado no elemento de assinatura "tira de manifesto" do Hero e
na seção Indicadores).

