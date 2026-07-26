# apps/console — Console OpenHands

Cette app héberge l'interface de l'agent (agent-canvas d'OpenHands).

## Deux façons d'intégrer OpenHands

### Option A — git submodule (contrôle total, recommandé)

```bash
# depuis la racine du monorepo
git submodule add https://github.com/All-Hands-AI/agent-canvas apps/console
git submodule add https://github.com/All-Hands-AI/software-agent-sdk packages/agent-sdk
```

### Option B — dépendance npm (si publiée)

```bash
cd apps/console
pnpm add @openhands/agent-sdk
```

## Déploiement Vercel

- Project: `axon-console`
- Root Directory: `apps/console`
