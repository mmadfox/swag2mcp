# Démarrage rapide

Mettez swag2mcp en fonctionnement en 2 minutes.

## 1. Initialisation

### Répertoire personnel (recommandé)

Configuration unique pour l'ensemble de votre système. La configuration est stockée dans votre dossier personnel.

::: code-group

```bash [macOS / Linux]
swag2mcp init
# Crée ~/.swag2mcp/swag2mcp.yaml
```

```powershell [Windows]
swag2mcp.exe init
# Crée %USERPROFILE%\.swag2mcp\swag2mcp.yaml
```

:::

### Répertoire du projet

Pour un espace de travail isolé au sein de votre projet.

::: code-group

```bash [macOS / Linux]
mkdir -p ./swag2mcp && swag2mcp init ./swag2mcp
```

```powershell [Windows]
mkdir ./swag2mcp; swag2mcp.exe init ./swag2mcp
```

:::

### Depuis un fichier ZIP

Si vous disposez d'un espace de travail prêt à l'emploi (par exemple, d'un collègue) :

```bash
swag2mcp import --from-zip workspace.zip
```

## 2. Installation des compétences de l'agent (recommandé)

Installez les compétences swag2mcp pour apprendre à votre agent IA toutes les commandes, les indicateurs, le format de configuration et les exemples concrets.

Demandez à votre agent :

```bash
"Ajoutez la compétence swag2mcp-cli depuis https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-cli/SKILL.md"
"Ajoutez la compétence swag2mcp-format depuis https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md"
```

> Certains IDE nécessitent un redémarrage après l'ajout de compétences.

## 3. Configuration du client LLM / IDE

Configurez votre IDE pour se connecter à swag2mcp. L'IDE démarrera le serveur MCP automatiquement si nécessaire.

::: code-group

```json [OpenCode]
{
  "mcp": {
    "swag2mcp": {
      "type": "local",
      "command": ["swag2mcp", "mcp"],
      "enabled": true
    }
  }
}
```

```json [Claude Desktop]
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

```json [Crush]
{
  "mcp": {
    "swag2mcp": {
      "type": "stdio",
      "command": "swag2mcp",
      "args": ["mcp"]
    }
  }
}
```

:::

Pour les autres IDE (Cursor, VS Code, JetBrains), consultez le [guide d'intégration](../integration/opencode.md).

> Si vous avez initialisé l'espace de travail dans un chemin personnalisé (par exemple `./swag2mcp`), utilisez le chemin complet dans la commande :
> `"command": ["swag2mcp", "mcp", "/chemin/absolu/vers/swag2mcp"]`

> **Après toute modification de configuration, redémarrez le serveur MCP** pour que les changements prennent effet.

## 4. Démarrage du serveur MCP

### stdio (par défaut) — pour IDE local

Rien à configurer. Votre IDE démarre swag2mcp automatiquement via la configuration ci-dessus.

```bash
swag2mcp mcp
```

### SSE / Streamable HTTP — pour accès distant

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

Ou configurez dans `swag2mcp.yaml` :

```yaml
mcp:
  transport: sse
  addr: ":8080"
  path: "/mcp"
```

Consultez la [référence du serveur MCP](../configuration/mcp-server.md) pour tous les indicateurs.

### Filtrer les spécifications par balises

```bash
swag2mcp mcp --tags=weather,public
```

Seules les spécifications avec les balises correspondantes seront disponibles pour le LLM.

### Vérifier le fonctionnement

Après la connexion, demandez à votre agent LLM :

```bash
"Quels outils MCP prenez-vous en charge ?"
```

Si l'agent liste les outils swag2mcp (`spec_list`, `search`, `invoke`, etc.) — tout fonctionne.

### Exemples de requêtes à essayer

| Demandez à votre agent | Ce qui se passe |
|-------|-------------|
| "Quel temps fait-il à Paris ?" | `invoke` — appelle l'API de prévisions Open-Meteo |
| "Quel est le prix actuel du BTC ?" | `invoke` — appelle l'API de cours Binance |
| "Racontez-moi une blague de papa" | `invoke` — appelle l'API icanhazdadjoke |
| "Montrez-moi Pikachu" | `invoke` — appelle l'API PokéAPI par nom |
| "Qui est Rick Sanchez ?" | `invoke` — appelle l'API des personnages Rick et Morty |
| "Quelle est la qualité de l'air à Pékin ?" | `invoke` — appelle l'API de qualité de l'air Open-Meteo |
| "Quelle est la hauteur des vagues près du Portugal ?" | `invoke` — appelle l'API marine Open-Meteo |
| "Recherchez des blagues sur les chiens" | `invoke` — appelle l'API de recherche dadjoke |
| "Listez tous les Pokémon" | `invoke` — appelle l'API de liste PokéAPI |
| "Quelle est l'altitude du Mont Blanc ?" | `invoke` — appelle l'API d'altitude Open-Meteo |

## 5. Prochaines étapes

- [Concepts](../concepts/overview.md) — comprendre l'architecture
- [Configuration](../configuration/config-file.md) — personnaliser les paramètres
- [Commandes CLI](../cli/overview.md) — référence complète des commandes
