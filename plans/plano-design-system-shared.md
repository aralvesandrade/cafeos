# Plano: CafeOS Design System v2 em packages/shared

## Objetivo

Implementar o CafeOS Design System v2 (fonte: `CafeOS Design System v2.dc.html` no projeto Claude Design) como biblioteca compartilhada em `packages/shared`, consumida por `apps/frontend` e `apps/admin`. Escopo desta rodada: **tokens + componentes primitivos**, sem tocar telas já existentes.

## Estado atual

- `apps/frontend` e `apps/admin` já têm o rebrand v2 aplicado (paleta terracota/marrom, IBM Plex Mono + Playfair Display) em `src/index.css` — commit `a59b2e9`.
- Tokens quase idênticos entre os dois apps, mas com pequenas divergências (frontend tem `--danger-*`/`--leaf` que admin não tem; admin tem `--sidebar-*` que frontend não tem).
- Faltam em ambos: `--shadow-sm/--shadow/--shadow-md/--shadow-lg`, `--bg-alt`, `--gold-foreground`, keyframes `ds-fade`/`ds-pop`/`ds-shimmer` (admin só tem `toast-in`).
- Componentes shadcn-style já existentes (duplicados entre os dois apps): `button`, `badge`, `input`, `select`, `dialog`, `tabs`, `field`, `required-legend`, `card` (só admin), `table` (só admin).
- Componentes do design system ainda **não implementados** em nenhum app: checkbox, radio, switch, range slider, combobox (select com busca), tooltip, dropdown-menu, alert, toast, avatar, progress bar, skeleton, accordion, pagination.
- Não existe workspace no monorepo hoje — sem `package.json` na raiz, sem `workspaces`. `packages/shared` está vazio.

## Decisão de arquitetura (confirmada com usuário)

Popular `packages/shared` com tokens CSS + componentes React reutilizáveis por `frontend` e `admin`, via npm workspaces.

## Passo a passo

### 1. Infra de workspace
- Criar `package.json` na raiz do monorepo com `"workspaces": ["apps/*", "packages/*"]`.
- Criar `packages/shared/package.json` (`name: "@cafeos/shared"`, `type: module`, exports para `./styles.css` e componentes TS).
- Criar `packages/shared/tsconfig.json` compatível com os `tsconfig.app.json` dos apps (mesmo `moduleResolution: bundler`, `jsx: react-jsx`).
- Ajustar `apps/frontend/tsconfig.app.json` e `apps/admin/tsconfig.app.json` para adicionar path `@cafeos/shared/*` (ou consumir via workspace symlink direto, sem path alias — a definir na implementação, o que for mais simples com Vite).
- Rodar `npm install` na raiz para linkar os workspaces (não vai quebrar builds existentes: apps continuam com seus `package.json` próprios).

### 2. Tokens (packages/shared/styles/tokens.css)
- Unificar os CSS custom properties dos dois apps + os que faltam do design doc:
  - Base surface/text, brand, destructive, semantic status (success/warning/info/danger), sidebar (theme-invariant), `--bg-alt`, `--gold`/`--gold-foreground`.
  - Shadow scale: `--shadow-sm`, `--shadow`, `--shadow-md`, `--shadow-lg` (mesmos valores em light/dark, conforme doc).
  - Keyframes: `ds-fade`, `ds-pop`, `ds-toast-in` (renomear/alinhar com `toast-in` já existente no admin), `ds-shimmer`.
  - Fontes: `--font-sans`/`--font-mono` = IBM Plex Mono, `--font-display` = Playfair Display (já corretos, só centralizar).
- `apps/frontend/src/index.css` e `apps/admin/src/index.css` passam a `@import "@cafeos/shared/styles/tokens.css"` antes do `@import "tailwindcss"` e mantêm só overrides específicos de cada app (ex: `--danger-*`/`--leaf` no frontend, `--sidebar-*` no admin — ou migrar tudo pro shared se fizer sentido unificar).

### 3. Componentes primitivos novos (packages/shared/src/components/ui/)
Implementar como componentes React + Tailwind (padrão shadcn, `cva` + `cn` util), cada um com variantes/estados descritos no design doc:
- `checkbox.tsx`, `radio-group.tsx`, `switch.tsx`, `slider.tsx` (range)
- `tooltip.tsx`, `dropdown-menu.tsx`
- `alert.tsx` (info/success/warning/error), `toast.tsx` + hook de toast stack
- `avatar.tsx` (single + stacked group)
- `progress.tsx`, `skeleton.tsx`, `accordion.tsx`, `pagination.tsx`
- `combobox.tsx` (select com busca, padrão "Município")

Reaproveitar utilitário `cn` (`lib/utils.ts`) — mover/centralizar em `packages/shared/src/lib/utils.ts`.

### 4. Migração dos componentes já duplicados (opcional/fora do escopo desta rodada)
Não mexer em `button`/`badge`/`input`/`select`/`dialog`/`tabs`/`field`/`card`/`table` existentes agora — ficam candidatos a uma segunda rodada de consolidação, pra não quebrar telas já em produção.

### 5. Validação
- `npm run build` em `apps/frontend` e `apps/admin` (tsc -b + vite build) pra confirmar que o workspace linka certo.
- Visual smoke test manual: abrir cada app, checar que tokens não regrediram (cores, fontes) e que os novos componentes renderizam corretamente em uma página de teste/story local (sem stories configuradas — validar via uma página temporária ou Storybook não incluso).

## Fora de escopo (a confirmar depois)
- Aplicar mudanças em telas existentes do admin (sidebar active state, badges de operação sempre terracota, tabelas, forms de cadastro).
- Migrar componentes shadcn já existentes pro shared.
- Storybook ou catálogo visual dos componentes.
