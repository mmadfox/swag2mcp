# swag2mcp-format

Le skill **swag2mcp-format** apprend à votre LLM à afficher les réponses des outils MCP swag2mcp dans un format Markdown propre, compact et lisible. Sans ce skill, le LLM décide lui-même comment formater les réponses — souvent de manière verbeuse et incohérente.

## Ce qu'il couvre

Tous les outils MCP swag2mcp :

- `spec_list`, `spec_by_id` — aperçu et détails des spécifications
- `collection_by_spec`, `collection_by_id` — collections avec étiquettes
- `tag_by_spec`, `tag_by_collection`, `tag_by_id` — listes d'étiquettes
- `endpoint_by_spec`, `endpoint_by_collection`, `endpoint_by_tag`, `endpoint_by_id` — listes de points de terminaison
- `search` — résultats de recherche
- `inspect` — détails complets de l'opération avec schémas compacts
- `invoke` — résultats d'appels API
- `auth` — informations d'authentification
- `info` — informations d'exécution

## Installation via un agent LLM

Copiez cette demande dans votre IDE alimenté par l'IA :

```
Crée le répertoire .agents/skills/swag2mcp-format/ et ajoute le skill depuis
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## Installation manuelle

```bash
mkdir -p .agents/skills/swag2mcp-format
curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
```

## Redémarrage requis

Après avoir ajouté le skill, redémarrez votre client LLM ou votre IDE (voir [Aperçu](overview.md#redémarrage-requis)).
