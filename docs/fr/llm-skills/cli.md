# swag2mcp-cli

Le skill **swag2mcp-cli** donne à votre LLM la référence complète de la CLI swag2mcp — chaque commande, drapeau, argument et option de configuration. Avec ce skill, le LLM peut répondre précisément aux questions "comment faire..." sans deviner.

## Ce qu'il couvre

Les 13 commandes CLI :

| Commande | Objectif |
|----------|----------|
| `init` | Initialiser un espace de travail et la configuration |
| `add` | Ajouter une spécification ou une collection |
| `delete` | Supprimer une spécification ou une collection |
| `ls` | Lister les spécifications configurées |
| `run` | Démarrer l'explorateur API TUI |
| `validate` | Valider le fichier de configuration |
| `clean` | Vider les données en cache |
| `update` | Mettre à jour le cache depuis la configuration |
| `mcp` | Démarrer le serveur MCP |
| `version` | Afficher les informations de version |
| `info` | Afficher les informations d'exécution |
| `import` | Importer un espace de travail depuis un fichier ZIP |
| `export` | Exporter l'espace de travail vers un fichier ZIP |

Plus tous les drapeaux, la structure du fichier de configuration, les méthodes d'authentification et les options avancées.

## Lien direct

<https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md>

## Installation via un agent LLM

Copiez cette demande dans votre IDE alimenté par l'IA :

```
Crée le répertoire .agents/skills/swag2mcp-cli/ et ajoute le skill depuis
https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Installation manuelle

```bash
mkdir -p .agents/skills/swag2mcp-cli
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Redémarrage requis

Après avoir ajouté le skill, redémarrez votre client LLM ou votre IDE (voir [Aperçu](overview.md#redémarrage-requis)).
