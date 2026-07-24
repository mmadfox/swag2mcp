# FAQ

## Général

### Qu'est-ce que swag2mcp et quel problème résout-il ?

swag2mcp fait le pont entre les spécifications d'API OpenAPI/Swagger/Postman et les agents LLM via le Model Context Protocol (MCP). Au lieu d'écrire du code personnalisé pour connecter chaque API à un agent IA, vous la configurez une fois dans un fichier YAML et le LLM obtient 19 outils pour découvrir, inspecter et appeler vos API.

### En quoi est-ce différent des autres outils API-LLM ?

- **Aucun codage requis** — configurez les API en YAML, pas de code d'intégration nécessaire
- **19 outils MCP** — boîte à outils complète de la découverte à l'invocation en passant par la gestion des grandes réponses
- **9 méthodes d'authentification** — fonctionne avec tout schéma d'authentification d'API
- **Recherche en texte intégral** — moteur de recherche bluge sur tous les points d'accès
- **Explorateur TUI** — interface terminal interactive pour naviguer et tester
- **Serveur mock** — testez sans appels API réels

### Quels formats de spécification d'API sont pris en charge ?

OpenAPI 3.x, Swagger 2.0 et Postman Collections v2.1.

### Quelle est la différence entre une spec et une collection ?

Une **spec** représente un service API logique (par exemple, « Open-Meteo Weather APIs »). Une **collection** est un fichier OpenAPI/Swagger/Postman. Une spec peut avoir plusieurs collections — par exemple, lorsqu'une API a des fichiers de spécification distincts pour différents services (prévisions, qualité de l'air, maritime).

### Quels transports MCP sont pris en charge ?

Trois transports : `stdio` (par défaut, pour les clients LLM locaux), `sse` (Server-Sent Events pour les clients distants) et `streamable-http` (streaming HTTP moderne).

### Puis-je utiliser swag2mcp avec n'importe quel LLM ?

Oui, tout client LLM qui prend en charge le protocole MCP : Claude Desktop, VS Code, Cursor, Windsurf, IDE JetBrains, OpenCode et autres.

## Installation

### Comment installer swag2mcp ?

```bash
# Option 1 : Téléchargement depuis GitHub Releases
# Allez sur https://github.com/mmadfox/swag2mcp/releases/latest
# Téléchargez l'archive pour votre OS et architecture

# Option 2 : Installation avec Go
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp@latest
```

### Dois-je avoir Go installé ?

Non. Des binaires pré-construits sont disponibles pour Linux (amd64, arm64), macOS (amd64, arm64) et Windows (amd64) sur la [page GitHub Releases](https://github.com/mmadfox/swag2mcp/releases).

### Comment installer le serveur mock ?

Le serveur mock est un binaire séparé :

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

Ou téléchargez `swag2mcp-mock_<version>_<os>_<arch>.tar.gz` depuis GitHub Releases.

## Premiers pas

### Comment démarrer rapidement ?

```bash
# 1. Initialisez un espace de travail
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. Démarrez le serveur MCP (des exemples de spécifications publiques sont inclus après init)
swag2mcp mcp
```

Après `init`, l'espace de travail inclut déjà plusieurs exemples de spécifications publiques (icanhazdadjoke, Open-Meteo, Binance, PokéAPI). Vous pouvez démarrer le serveur MCP immédiatement — pas besoin d'ajouter des spécifications manuellement.

Si vous souhaitez ajouter votre propre API à la place :

```bash
swag2mcp add spec --yaml - <<EOF
domain: dadjoke
llm_title: API icanhazdadjoke
base_url: https://icanhazdadjoke.com
collections:
  - llm_title: Blagues
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
EOF
```

### Comment connecter swag2mcp à mon IDE ?

**VS Code** (`.vscode/settings.json`) :
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

**Cursor** (`~/.cursor/mcp.json`) :
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

**Claude Desktop** (`claude_desktop_config.json`) :
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

Utilisez toujours un chemin absolu vers le répertoire de l'espace de travail.

## Configuration

### Où se trouve le fichier de configuration ?

Par défaut : `~/.swag2mcp/swag2mcp.yaml`. Vous pouvez également le créer dans n'importe quel répertoire et passer le chemin aux commandes.

### Comment ajouter une API ?

```bash
# Mode interactif
swag2mcp add spec

# Avec YAML (recommandé pour les scripts)
swag2mcp add spec --yaml - <<EOF
domain: mon-api
llm_title: Mon API
base_url: https://api.example.com/v1
collections:
  - llm_title: Principal
    location: https://example.com/spec.yaml
EOF
```

### Comment ajouter une collection à une spec existante ?

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Qualité de l'air
location: https://example.com/air-quality.yaml
EOF
```

### Comment désactiver temporairement une spec ?

Définissez `disable: true` dans la configuration de la spec. La spec ne sera pas chargée ni indexée.

### Puis-je filtrer les specs chargées ?

Oui, utilisez le drapeau `--tags` : `swag2mcp mcp --tags=public`. Seules les specs avec des étiquettes correspondantes seront chargées.

### Comment utiliser les variables d'environnement pour les secrets ?

Utilisez la syntaxe `$(NOM_VAR)` dans les champs d'authentification :

```yaml
auth:
  type: bearer
  config:
    token: "$(MON_JETON_API)"
```

Définissez la variable avant de démarrer : `export MON_JETON_API="eyJhbGci..."`

## Authentification

### Quelles méthodes d'authentification sont prises en charge ?

Neuf méthodes : `none`, `basic`, `bearer`, `digest`, `hmac`, `oauth2-cc` (identifiants client), `oauth2-pwd` (mot de passe), `api-key` et `script`.

### Comment transmettre un jeton ?

Via le fichier de configuration ou les variables d'environnement :

```yaml
auth:
  type: bearer
  config:
    token: "$(MON_JETON)"
```

### Dois-je appeler auth avant invoke ?

Non. L'outil `invoke` applique automatiquement l'authentification à partir de la configuration de la spec. Vous n'avez besoin de l'outil MCP `auth` que si vous souhaitez afficher le jeton à l'utilisateur (par exemple, pour une commande curl).

### Pourquoi l'outil auth n'apparaît-il pas ?

L'outil `auth` est désactivé par défaut (`--disable-llm-auth=true`). C'est une mesure de sécurité pour la production. Pour l'activer : `swag2mcp mcp --disable-llm-auth=false`.

### Comment les jetons OAuth2 sont-ils renouvelés ?

Les jetons OAuth2 Client Credentials et Password Grant sont automatiquement renouvelés à leur expiration. Les jetons Bearer sont statiques et doivent être mis à jour manuellement.

## Serveur MCP

### Comment démarrer le serveur MCP ?

```bash
# Par défaut (transport stdio)
swag2mcp mcp

# Avec transport HTTP
swag2mcp mcp --transport sse --http-addr :8080
```

### Comment changer le port ?

```bash
swag2mcp mcp --transport sse --http-addr 0.0.0.0:9090
```

### Comment sécuriser le point d'accès HTTP MCP ?

Définissez un jeton Bearer :

```bash
swag2mcp mcp --transport sse --http-addr :8080 --auth-token "mon-secret"
```

Le client LLM doit inclure `Authorization: Bearer mon-secret` dans chaque requête.

### Qu'est-ce que la poignée de main MCP pour le transport HTTP ?

Pour les transports SSE et Streamable HTTP, le protocole MCP nécessite une poignée de main en trois étapes :

```
Étape 1 : POST /mcp → {"method":"initialize", ...}
Étape 2 : POST /mcp → {"method":"notifications/initialized"}
Étape 3 : POST /mcp → {"method":"tools/list", ...}  ← fonctionne maintenant
```

Les appels d'outils échoueront avant l'initialisation.

## Utilisation

### Comment rechercher des points d'accès ?

Utilisez l'outil MCP `search` ou la TUI (`swag2mcp run`). La recherche prend en charge les filtres de champ (`method:GET`, `tag:pets`), la recherche floue, les jokers et les opérateurs booléens.

### Comment appeler une API ?

Le LLM utilise l'outil MCP `invoke`. Inspectez toujours d'abord le point d'accès pour comprendre les paramètres requis :

```
inspect(endpointId: "...")  → comprendre le contrat
invoke(endpointId: "...", parameters: {...})  → effectuer l'appel
```

### Que se passe-t-il si une réponse est trop volumineuse ?

Les réponses dépassant `max_response_size` (par défaut 1 Mo) sont sauvegardées sur le disque. Le LLM reçoit une référence de fichier et peut l'explorer avec les outils `response_outline`, `response_compress` et `response_slice`.

### Comment fonctionne le limiteur de débit ?

Chaque point d'accès a un délai de 10 secondes. Si le LLM appelle le même point d'accès deux fois en 10 secondes, le second appel est silencieusement bloqué. Vous pouvez désactiver ou ajuster cela dans la configuration.

### Puis-je tester sans effectuer d'appels API réels ?

Oui, utilisez le serveur mock :

```bash
swag2mcp-mock mockserver
```

Il génère des réponses factices basées sur les schémas OpenAPI.

## Gestion de l'espace de travail

### Comment sauvegarder ma configuration ?

```bash
swag2mcp export --output ~/sauvegardes/swag2mcp-2026-07-24.zip
```

### Comment transférer vers une autre machine ?

```bash
# Sur l'ancienne machine
swag2mcp export --output swag2mcp.zip

# Copiez le ZIP, puis sur la nouvelle machine
swag2mcp import --from-zip swag2mcp.zip
```

### Comment mettre à jour les fichiers de spécification ?

```bash
swag2mcp update
```

Cela revalide la configuration, vide le cache et retélécharge tous les fichiers de spécification.

### Comment libérer de l'espace disque ?

```bash
swag2mcp clean
```

Supprime les fichiers de spécification en cache et les réponses d'API sauvegardées. Les anciennes réponses (>48h) sont également nettoyées automatiquement au démarrage du serveur MCP.

## TUI

### Qu'est-ce que la TUI et comment l'utiliser ?

La TUI (Interface Utilisateur Terminal) est un explorateur d'API interactif. Lancez-la avec `swag2mcp run`. Elle a trois modes : Recherche (recherche en texte intégral), Parcourir (navigation arborescente : Spec → Collection → Tag → Point d'accès) et Auth (afficher les jetons).

### Quels sont les raccourcis clavier ?

| Touche | Action |
|--------|--------|
| `↑/↓` | Naviguer |
| `Entrée` | Sélectionner |
| `Échap` | Retour |
| `Tab` | Changer de mode |
| `/` | Rechercher |
| `N/P` | Page suivante/précédente |
| `q` | Quitter |

## Avancé

### Puis-je utiliser un proxy ?

Oui, configurez-le dans `http_client.proxy` :

```yaml
http_client:
  proxy:
    url: "http://proxy.entreprise.com:8080"
    username: "$(UTILISATEUR_PROXY)"
    password: "$(MOT_DE_PASSE_PROXY)"
    bypass:
      - "localhost"
      - "*.interne.com"
```

### Puis-je ajouter une méthode d'authentification personnalisée ?

Oui, implémentez l'interface `Authenticator` dans `internal/auth/` et enregistrez-la dans l'analyseur de configuration. Voir la section Développement pour plus de détails.

### Puis-je ajouter un outil MCP personnalisé ?

Oui, ajoutez une méthode à l'interface `Svc`, implémentez-la dans la couche service, ajoutez un gestionnaire et enregistrez-le. Voir la section Développement pour plus de détails.

### Quelle est la différence entre `swag2mcp` et `swag2mcp-mock` ?

`swag2mcp` est le binaire principal avec les commandes CLI et le serveur MCP. `swag2mcp-mock` est un binaire séparé qui démarre des serveurs mock pour les tests sans appels API réels.
