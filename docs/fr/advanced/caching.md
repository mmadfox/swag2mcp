# Mise en cache

## Aperçu

swag2mcp met en cache les fichiers de spécification téléchargés afin que le serveur MCP démarre plus rapidement lors des exécutions suivantes. Au lieu de télécharger le même fichier de spécification à chaque fois, il réutilise la copie en cache.

## Comment fonctionne la mise en cache

Lorsque vous ajoutez une spec avec une URL distante, swag2mcp la télécharge et la sauvegarde dans le répertoire `cache/`. Au prochain démarrage, il vérifie si la copie en cache est encore fraîche. Si c'est le cas, le téléchargement est ignoré.

### Ce qui est mis en cache

| Source | Comportement |
|--------|--------------|
| **URL distante** (http/https) | Toujours mise en cache. Téléchargée une fois, réutilisée jusqu'à l'expiration du cache. |
| **Fichier local dans `specs/`** | Utilisé directement depuis le répertoire `specs/`. Jamais mis en cache — les modifications sont immédiatement visibles. |
| **Fichier local hors de `specs/`** | Copié dans le cache. Si le fichier source change (heure de modification), le cache est invalidé. |

### Expiration du cache (TTL)

Chaque fichier en cache reçoit un temps d'expiration aléatoire entre **1 heure et 48 heures**. Le caractère aléatoire empêche tous les fichiers en cache d'expirer en même temps (ce qui provoquerait une ruée de téléchargements).

- Le TTL est réinitialisé à chaque démarrage du serveur MCP
- Si un fichier en cache est encore dans son TTL, il est réutilisé
- Si le TTL a expiré, le fichier est téléchargé à nouveau

### Structure du cache

```
~/.swag2mcp/cache/
├── a1b2c3d4e5f6a7b8.spec    # Fichier de spécification en cache
├── a1b2c3d4e5f6a7b8.meta    # Métadonnées (source, TTL, date de mise en cache)
├── b2c3d4e5f6a7b8c9.spec
├── b2c3d4e5f6a7b8c9.meta
└── ...
```

La clé de cache est dérivée de l'URL ou du chemin du fichier de spécification. Chaque fichier en cache a un fichier `.meta` compagnon qui stocke quand il a été mis en cache et quand il expire.

## Gestion du cache

### Forcer un rafraîchissement

Exécutez `swag2mcp update` pour vider tout le cache et retélécharger tous les fichiers de spécification :

```bash
swag2mcp update
```

Cela valide la configuration, vide le cache et télécharge tout à nouveau.

### Vider le cache manuellement

```bash
swag2mcp clean
```

Cela supprime tous les fichiers de spécification en cache et les réponses d'API sauvegardées. La prochaine fois que vous démarrerez le serveur MCP, toutes les specs seront téléchargées à nouveau.

### Nettoyage automatique

Lorsque le serveur MCP démarre (`swag2mcp mcp`), les réponses d'API sauvegardées de plus de 48 heures sont automatiquement supprimées. Cela empêche le répertoire `responses/` de croître indéfiniment.

## Notes importantes

- **Les fichiers locaux dans `specs/` ne sont jamais mis en cache** — si vous modifiez un fichier de spécification directement dans le répertoire `specs/`, les modifications sont immédiatement visibles sans vider le cache
- **Les URL distantes sont toujours mises en cache** — il n'y a aucun moyen de contourner le cache pour les URL distantes, sauf en exécutant `swag2mcp update` ou `swag2mcp clean`
- **Le cache est local** — il est stocké sur le disque et ne se synchronise pas entre les machines. Utilisez `swag2mcp export` et `swag2mcp import` pour transférer les specs entre machines
