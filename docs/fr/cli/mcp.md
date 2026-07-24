# mcp

## Objectif

Démarrer le **serveur MCP (Model Context Protocol)** — le mode principal pour l'intégration LLM. C'est ce que vous exécutez pour donner à un agent IA (Claude, Cursor, OpenCode, etc.) l'accès à vos API via 16 outils MCP.

## Quand l'utiliser

- Vous voulez connecter un agent LLM à vos API
- Vous configurez un IDE (VS Code, Cursor, JetBrains) ou une application de bureau (Claude Desktop)
- Vous devez exposer vos API via le protocole MCP
- Vous testez le serveur MCP avant l'intégration

## Syntaxe

```bash
swag2mcp mcp [chemin] [drapeaux]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |

## Drapeaux

| Drapeau | Raccourci | Type | Défaut | Description |
|---------|-----------|------|--------|-------------|
| `--transport` | | `string` | `"stdio"` | Transport MCP : `stdio`, `sse`, `streamable-http` |
| `--http-addr` | | `string` | `":8080"` | Adresse du serveur HTTP (pour `sse` et `streamable-http`) |
| `--http-path` | | `string` | `"/mcp"` | Chemin HTTP pour le gestionnaire MCP |
| `--auth-token` | | `string` | `""` | Jeton Bearer pour l'authentification du transport HTTP |
| `--logfile` | `-f` | `string` | `""` | Chemin du fichier journal. S'il n'est pas défini, les journaux vont sur stderr. |
| `--disable-llm-auth` | | `bool` | `true` | Supprimer l'outil `auth` de la liste des outils MCP |
| `--dump-dir` | | `string` | `""` | Répertoire pour vider les requêtes HTTP pour le débogage |
| `--tags` | `-t` | `string` | `""` | Filtrer les specs par étiquettes (séparées par des virgules) |

## Comment cela fonctionne

### Transport stdio (par défaut)

Utilisé lorsque le serveur MCP est lancé comme sous-processus par le client LLM (IDE, Claude Desktop, etc.). Le serveur communique via l'entrée/sortie standard.

```bash
swag2mcp mcp
```

### Transport SSE

Transport Server-Sent Events pour la communication HTTP. Nécessite la séquence de poignée de main MCP.

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### Transport Streamable HTTP

Transport HTTP moderne qui prend en charge les réponses en streaming.

```bash
swag2mcp mcp --transport streamable-http --http-addr 0.0.0.0:8080
```

### Avec authentification

Protégez le point d'accès HTTP avec un jeton Bearer :

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "mon-secret"
```

### Avec filtrage par étiquettes

Chargez uniquement les specs avec des étiquettes spécifiques :

```bash
swag2mcp mcp --tags=public
```

### Avec outil auth activé (mode débogage)

Permettez au LLM de demander des jetons frais via l'outil `auth` :

```bash
swag2mcp mcp --disable-llm-auth=false
```

### Avec répertoire de vidage des requêtes

Sauvegardez toutes les requêtes HTTP pour le débogage :

```bash
swag2mcp mcp --dump-dir ./vidages
```

## Transport HTTP MCP — Protocole de poignée de main

Lors de l'utilisation de `sse` ou `streamable-http`, le protocole MCP nécessite une poignée de main spécifique. Les appels d'outils échoueront avant l'initialisation :

```
Étape 1 : POST /mcp → {"method":"initialize", ...}
Étape 2 : POST /mcp → {"method":"notifications/initialized"}
Étape 3 : POST /mcp → {"method":"tools/list", ...}   ← fonctionne maintenant
```

### Health check

Fonctionne sans initialisation :

```bash
curl http://localhost:8080/health
# → {"status":"ok","version":"v1.2.0"}
```

## Exemples de configuration IDE

### VS Code (`.vscode/settings.json` ou paramètres globaux)

```json
{
  "mcp": {
    "servers": {
      "swag2mcp": {
        "command": "swag2mcp",
        "args": ["mcp", "/chemin/absolu/vers/.swag2mcp"]
      }
    }
  }
}
```

### Cursor / Windsurf (`~/.cursor/mcp.json` ou projet `.cursor/mcp.json`)

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/chemin/absolu/vers/.swag2mcp"]
    }
  }
}
```

### Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json` sur macOS)

```json
{
  "mcpServers": {
    "swag2mcp": {
      "command": "swag2mcp",
      "args": ["mcp", "/chemin/absolu/vers/.swag2mcp"]
    }
  }
}
```

### IDE JetBrains (Paramètres → Outils → MCP)

- Nom : `swag2mcp`
- Commande : `swag2mcp`
- Arguments : `mcp /chemin/absolu/vers/.swag2mcp`

> **Utilisez toujours un chemin absolu** vers le répertoire de l'espace de travail dans la configuration IDE. Les chemins relatifs peuvent échouer selon le répertoire de travail de l'IDE.

## Sortie

En cas de succès, le serveur affiche :

```
Serveur MCP à l'écoute sur http://127.0.0.1:8080/mcp
```

## Nuances

- **Pas d'auto-initialisation :** Si le fichier de configuration n'existe pas, `mcp` retourne une erreur : « configuration introuvable à &lt;chemin&gt; ». Exécutez `init` d'abord.
- **`--disable-llm-auth` (défaut : `true`) :** Lorsqu'il est activé, l'outil `auth` est complètement supprimé de la liste des outils MCP. Le LLM ne peut pas voir ni demander de jetons. L'authentification fonctionne toujours — les jetons sont obtenus via le mécanisme de configuration standard, pas via le LLM. Ce mode est recommandé pour la **production**. Pour le **débogage** ou lors de l'utilisation de jetons de courte durée, définissez `--disable-llm-auth=false` pour permettre au LLM de demander des jetons frais via l'outil `auth`.
- **Recours à la configuration YAML :** Si un drapeau CLI n'est pas explicitement défini, la valeur est prise de la section `mcp` dans `swag2mcp.yaml` (si présente). Cela vous permet de configurer le serveur dans le fichier de configuration au lieu de passer des drapeaux à chaque fois.
- **Nettoyage des réponses :** Au démarrage, les réponses de plus de 48 heures sont automatiquement supprimées du répertoire `responses/`.
- **Avertissement de résolution de chemin :** Lorsque `[chemin]` est omis, `mcp` cherche `swag2mcp.yaml` dans le répertoire courant d'abord, puis utilise `~/.swag2mcp/` comme recours. Si vous exécutez la commande depuis le mauvais répertoire, elle peut charger un espace de travail différent de celui prévu. **Spécifiez toujours `[chemin]` explicitement lorsque vous l'exécutez comme service ou dans une configuration IDE.**
