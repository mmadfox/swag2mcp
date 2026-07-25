# clean

## Objectif

Supprimer les specs distantes en cache et les réponses d'invocation API sauvegardées. Cela libère de l'espace disque et force un téléchargement frais des fichiers de spécification lors du prochain démarrage de `update` ou `mcp`.

## Quand l'utiliser

- Les fichiers de spécification ont changé sur le serveur distant et vous voulez forcer un rafraîchissement
- Vous voulez libérer de l'espace disque
- Vous résolvez des problèmes de cache obsolète
- Avant d'exécuter `update` pour garantir un re-cache propre

## Syntaxe

```bash
swag2mcp clean [chemin]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |

## Drapeaux

Aucun.

## Comment cela fonctionne

```bash
swag2mcp clean
swag2mcp clean ./mon-espace-travail
```

## Ce qui est nettoyé

| Répertoire | Contenu | Pourquoi |
|------------|---------|----------|
| `cache/` | Fichiers de spécification distants téléchargés | Force le re-téléchargement au prochain accès |
| `responses/` | Réponses d'invocation API sauvegardées | Libère de l'espace disque |

## Ce qui est préservé

| Répertoire | Contenu | Pourquoi |
|------------|---------|----------|
| `specs/` | Fichiers de spécification locaux | Ce sont vos fichiers source, pas le cache |
| `auth_scripts/` | Scripts d'authentification | Ils sont créés par l'utilisateur, pas du cache |

## Nettoyage des scripts orphelins

Après le nettoyage, `clean` supprime également les scripts d'authentification pour les specs qui n'existent plus dans la configuration. Cela empêche l'accumulation de scripts obsolètes.

## Nettoyage automatique

Lorsque le serveur MCP démarre (`swag2mcp mcp`), les réponses de plus de 48 heures sont supprimées automatiquement. Vous n'avez généralement pas besoin d'exécuter `clean` manuellement pour la maintenance de routine.

## Vérification post-commande

```bash
ls ~/.swag2mcp/cache
# Devrait être vide (le répertoire existe mais n'a pas de fichiers)
```

## Nuances

- **Aucune configuration requise :** `clean` fonctionne même sans fichier de configuration valide. Il supprime simplement les répertoires de cache et de réponses.
- **Le nettoyage des orphelins est au mieux :** Si le fichier de configuration est corrompu ou illisible, le nettoyage des scripts d'authentification orphelins est ignoré (non fatal).
- **Les répertoires sont préservés :** Les répertoires `cache/` et `responses/` eux-mêmes sont conservés — seul leur contenu est supprimé.
