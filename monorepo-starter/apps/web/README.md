# apps/web — Site Axon

Copie ici tout le contenu du projet v0 actuel (le site marketing Axon) :

- `app/`
- `components/`
- `public/`
- `package.json` (garde ses dépendances Next.js)
- `next.config.mjs`, `tsconfig.json`, `postcss.config.mjs`, `app/globals.css`

## Utiliser le package UI partagé

Dans le `package.json` de cette app, ajoute :

```json
{
  "dependencies": {
    "@axon/ui": "workspace:*"
  }
}
```

Puis importe le logo depuis le package partagé :

```tsx
import { Logo } from "@axon/ui/logo"
```

## Déploiement Vercel

- Project: `axon-web`
- Root Directory: `apps/web`
