# update

## Objectif

Revalider la configuration, vider le cache et retélécharger tous les fichiers de spécification. C'est un **rafraîchissement complet** de l'espace de travail — il garantit que toutes les specs en cache sont à jour et que l'index est reconstruit.

## Quand l'utiliser

- Les fichiers de spécification distants ont changé et vous voulez la dernière version
- Après avoir modifié `swag2mcp.yaml` pour ajouter ou changer des emplacements de spécification
- Lors du dépannage d'un cache obsolète ou corrompu
- Avant d'exécuter `mcp` pour garantir que tout est frais

## Syntaxe

```bash
swag2mcp update [chemin]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |

## Drapeaux

Aucun.

## Comment cela fonctionne

La commande `update` exécute un pipeline d'opérations :

1. **Charger la configuration** — lit `swag2mcp.yaml` depuis l'espace de travail
2. **Valider** — exécute les mêmes vérifications que `validate` (syntaxe YAML, structure, accessibilité des fichiers de spécification, format, authentification, client HTTP)
3. **Nettoyer** — supprime tout le contenu de `cache/` et `responses/`
4. **Re-mettre en cache** — télécharge tous les fichiers de spécification distants et copie les fichiers locaux dans le cache
5. **Réindexer** — reconstruit l'index de recherche en texte intégral pour tous les points d'accès
6. **Scripts d'authentification** — crée des scripts stub pour les specs utilisant `ScriptAuth`
7. **Nettoyage des orphelins** — supprime les scripts d'authentification pour les specs qui n'existent plus

```bash
swag2mcp update
swag2mcp update ./mon-espace-travail
```

## Ce qui arrive aux collections désactivées

Les collections avec `disable: true` sont complètement ignorées — elles ne sont ni mises en cache ni indexées.

## Vérification post-commande

```bash
swag2mcp ls [chemin]
# Toutes les specs devraient toujours être listées et accessibles
```

## Nuances

- **Pas d'auto-initialisation :** Si le fichier de configuration n'existe pas, `update` retourne une erreur : « configuration introuvable à <chemin> ». Exécutez `init` d'abord.
- **Dépendance réseau :** Toutes les URL de spécification distantes doivent être accessibles. Si un téléchargement échoue, la mise à jour entière échoue avec un message d'erreur clair.
- **Création de script d'authentification :** Si une spec utilise `ScriptAuth` et que le script stub n'existe pas, `update` le crée. Si la création échoue, la mise à jour échoue.
- **`update` vs `clean` :** `clean` supprime uniquement le cache. `update` supprime le cache **et** retélécharge tout. Utilisez `clean` lorsque vous voulez juste libérer de l'espace ; utilisez `update` lorsque vous voulez rafraîchir.
