# Flux de travail CLI

Cette page présente des exemples concrets d'utilisation de swag2mcp depuis le terminal — de l'initialisation aux opérations quotidiennes.

## Démarrage rapide

```bash
# 1. Initialiser un espace de travail
mkdir -p .swag2mcp && swag2mcp init ./.swag2mcp

# 2. Lister vos spécifications
swag2mcp ls
```

## Ajout d'une spécification avec YAML

### Spécification simple (API publique)

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: Open-Meteo Weather API
base_url: https://api.open-meteo.com
collections:
  - llm_title: Weather Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
EOF
```

### Spécification avec authentification (jeton bearer depuis l'environnement)

```bash
swag2mcp add spec --yaml - <<EOF
domain: mon-api
llm_title: Mon API protégée
base_url: https://api.example.com/v1
auth:
  type: bearer
  config:
    token: \$(MON_JETON)
collections:
  - llm_title: Utilisateurs
    location: https://raw.githubusercontent.com/mon-org/mon-api/main/users.yaml
EOF
```

### Spécification avec plusieurs collections

```bash
swag2mcp add spec --yaml - <<EOF
domain: meteo
llm_title: API Open-Meteo
base_url: https://api.open-meteo.com
collections:
  - llm_title: Forecast
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
  - llm_title: Air Quality
    base_url: https://air-quality-api.open-meteo.com
    location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
EOF
```

## Ajout d'une collection à une spécification existante

```bash
swag2mcp add collection --yaml - <<EOF
spec_domain: meteo
llm_title: Météo marine
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
EOF
```

## Liste des spécifications

```bash
$ swag2mcp ls
Spécifications :
  dadjoke (https://icanhazdadjoke.com)
    jokes (3 points de terminaison)
  meteo (https://api.open-meteo.com)
    forecast (5 points de terminaison)
    air-quality (8 points de terminaison)
    marine (4 points de terminaison)
```

### Filtrer par balises

```bash
swag2mcp ls --tags=public
```

## Affichage des informations d'exécution

```bash
$ swag2mcp info
{
  "version": "v1.2.0",
  "workspace": "/home/utilisateur/.swag2mcp",
  "specs": {
    "total": 2,
    "active": 2,
    "disabled": 0,
    "collections": 4,
    "endpoints": 20
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 Ko",
    "follow_redirects": true
  },
  "mcp": {
    "transport": "stdio"
  },
  "auth": {
    "methods": ["bearer"]
  }
}
```

## Validation de la configuration

```bash
$ swag2mcp validate
✅ La configuration est valide.
✓ Spécification dadjoke : OK
✓ Spécification meteo : OK
```

## Démarrage du serveur MCP

### stdio (pour l'intégration IDE)

```bash
swag2mcp mcp
```

### HTTP (pour l'accès distant)

```bash
swag2mcp mcp --transport sse --http-addr :8080
```

### Avec filtre par balises

```bash
swag2mcp mcp --tags=public
```

## Mise à jour des spécifications

Actualisez tous les fichiers de spécification en cache :

```bash
swag2mcp update
```

## Nettoyage du cache

```bash
swag2mcp clean
```

## Exportation et importation

### Sauvegarder votre espace de travail

```bash
swag2mcp export --output ~/sauvegardes/swag2mcp-2026-07-24.zip
```

### Restaurer sur une autre machine

```bash
# Sur la nouvelle machine
swag2mcp import --from-zip swag2mcp-2026-07-24.zip
```

## Explorateur TUI interactif

```bash
swag2mcp run
```

Ouvre une interface utilisateur terminal plein écran pour rechercher, parcourir et invoquer des API.

## Serveur de simulation

```bash
# Installer le binaire de simulation
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest

# Démarrer les serveurs de simulation
swag2mcp-mock mockserver
```
