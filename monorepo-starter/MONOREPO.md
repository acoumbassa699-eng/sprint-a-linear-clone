# Axon Monorepo — Guide complet (de A à Z)

Ce dossier contient la **base prête à l'emploi** de ton monorepo. Copie tout le
contenu de `monorepo-starter/` à la racine de ton nouveau dépôt Git.

## Structure

```
axon-monorepo/
├── apps/
│   ├── web/              ← ton site Axon (copie le projet v0 ici)
│   └── console/          ← l'UI OpenHands (agent-canvas)
├── packages/
│   ├── ui/               ← composants partagés (Logo Axon)
│   └── agent-sdk/        ← software-agent-sdk d'OpenHands (submodule)
├── package.json          ← racine (workspaces + scripts)
├── pnpm-workspace.yaml
└── turbo.json
```

## Étapes

### 1. Initialiser le dépôt
```bash
mkdir axon-monorepo && cd axon-monorepo
# copie ici le contenu de monorepo-starter/
git init
```

### 2. Installer les dépendances
```bash
pnpm install
```

### 3. Placer le site Axon dans apps/web
Copie tout le projet v0 actuel (app/, components/, public/, package.json,
next.config.mjs, etc.) dans `apps/web/`. Voir `apps/web/README.md`.

### 4. Ajouter OpenHands dans apps/console + packages/agent-sdk
```bash
git submodule add https://github.com/All-Hands-AI/agent-canvas apps/console
git submodule add https://github.com/All-Hands-AI/software-agent-sdk packages/agent-sdk
```

### 5. Lancer en dev (les deux apps en parallèle)
```bash
pnpm dev
```

### 6. Build de production
```bash
pnpm build
```

## Déploiement Vercel (2 projets, 1 dépôt)

| Projet Vercel   | Root Directory   | Rôle                  |
|-----------------|------------------|-----------------------|
| `axon-web`      | `apps/web`       | site marketing        |
| `axon-console`  | `apps/console`   | console OpenHands      |

Vercel détecte automatiquement Turborepo et ne rebuild que ce qui a changé.

## Package UI partagé

Le logo Axon vit dans `packages/ui`. Chaque app l'utilise via :
```json
{ "dependencies": { "@axon/ui": "workspace:*" } }
```
```tsx
import { Logo } from "@axon/ui/logo"
```

## Références OpenHands (juillet 2026)

- `software-agent-sdk` → cœur de l'agent + serveur
- `agent-canvas` → interface utilisateur
- Repo principal : https://github.com/All-Hands-AI/OpenHands
