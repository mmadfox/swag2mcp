# Skills for LLM

Les skills sont des fichiers Markdown qui apprennent à votre agent LLM à travailler plus efficacement avec swag2mcp. Ils sont chargés dans le prompt système de l'agent et donnent au LLM des instructions précises pour formater les réponses et comprendre les commandes CLI.

## Skills disponibles

| Skill | Description | Télécharger |
|-------|-------------|-------------|
| **swag2mcp-format** | Formate les réponses des outils MCP en tableaux Markdown compacts et lisibles | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md) |
| **swag2mcp-cli** | Référence CLI complète — le LLM connaît chaque commande, drapeau et option de configuration | [SKILL.md](https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md) |

## Pourquoi les skills sont importants

Sans skill de formatage, le LLM décide lui-même comment afficher les résultats — souvent de manière verbeuse et incohérente. Le skill de formatage garantit un style propre et uniforme : tableaux compacts pour les listes, en-têtes en ligne pour les détails et schémas concis.

Le skill CLI permet au LLM de répondre précisément aux questions "comment faire..." sur les commandes swag2mcp, sans deviner.

## Installation via un agent LLM

Copiez cette demande dans votre IDE alimenté par l'IA (OpenCode, Cursor, Claude Desktop, VS Code, etc.) :

```
Ajoute les skills swag2mcp à mon projet :

1. Crée le répertoire .agents/skills/swag2mcp-format/ et ajoute le skill depuis https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
2. Crée le répertoire .agents/skills/swag2mcp-cli/ et ajoute le skill depuis https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

L'agent téléchargera les deux fichiers de skills et les placera aux bons endroits.

## Installation manuelle

Si votre client LLM ne prend pas en charge la configuration via agent, téléchargez les fichiers manuellement :

```bash
mkdir -p .agents/skills/swag2mcp-format
mkdir -p .agents/skills/swag2mcp-cli

curl -o .agents/skills/swag2mcp-format/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md
curl -o .agents/skills/swag2mcp-cli/SKILL.md \
  https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md
```

## Configurer votre client LLM

Chaque client LLM et IDE a sa propre méthode d'installation des skills. L'exemple ci-dessous est pour **OpenCode** — consultez la documentation de votre client pour la méthode correcte.

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    },
    {
      "name": "swag2mcp-cli",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-cli/SKILL.md"
    }
  ]
}
```

## Redémarrage requis

**Après avoir ajouté les skills, redémarrez votre client LLM ou votre IDE.** Certains outils ne chargent les skills qu'au démarrage. Si les skills ne semblent pas prendre effet, essayez :

- **OpenCode** : Redémarrez l'application ou exécutez à nouveau la commande opencode
- **Cursor** : Fermez et rouvrez la fenêtre (`Cmd+Shift+W` / `Ctrl+Shift+W`)
- **Claude Desktop** : Quittez et relancez l'application
- **VS Code** : Rechargez la fenêtre (`Ctrl+Shift+P` → "Developer: Reload Window")
