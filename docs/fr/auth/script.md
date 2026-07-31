# Authentification par script

## Objectif

Authentification via un script externe. swag2mcp prend en charge deux types de scripts :
- **`.sh`** — scripts shell sur Linux et macOS
- **`.bat`** — scripts batch sur Windows

Le script doit être placé dans le répertoire `auth_scripts` de votre espace de travail et afficher un jeton JSON sur stdout.

## Quand l'utiliser

- Schémas d'authentification personnalisés ou non standard
- Logique d'acquisition de jeton complexe (multi-étapes, avec vérifications supplémentaires)
- Quand aucune des méthodes standard ne correspond à vos besoins

## Configuration

```yaml
specs:
  - domain: jokes
    llm_title: API Dad Joke
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Blagues
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: script
      config:
        domain: "jokes"
```

## Paramètres

| Paramètre | Requis | Description |
|-----------|--------|-------------|
| `domain` | Oui | Nom du fichier de script (sans extension). Doit correspondre au `domain` de la spec — par ex. si `domain: jokes`, le fichier doit être `jokes.sh` (Unix) ou `jokes.bat` (Windows). |

## Emplacement du script

Le script doit être placé dans le répertoire `auth_scripts` de votre espace de travail :

- **Linux / macOS :** `{espace-travail}/auth_scripts/{domaine}.sh`
- **Windows :** `{espace-travail}/auth_scripts/{domaine}.bat`

## Format de sortie du script

Le script doit produire du JSON sur stdout avec le jeton et son temps d'expiration :

```bash
#!/bin/bash
# auth_scripts/jokes.sh

JETON=$(curl -s -X POST https://auth.example.com/token \
  -d "grant_type=client_credentials" \
  -d "client_id=$ID_CLIENT" \
  -d "client_secret=$SECRET_CLIENT" | jq -r '.access_token')

echo "{\"token\": \"$JETON\", \"expires_in\": 3600}"
```

### Champs JSON

| Champ | Requis | Description |
|-------|--------|-------------|
| `token` | Oui | Jeton d'authentification |
| `expires_in` | Non | Durée de vie du jeton en secondes (défaut : 3600) |

## Notes

- swag2mcp exécute le script à chaque requête si le jeton en cache a expiré
- Le script doit se terminer dans les 30 secondes
- Le jeton est mis en cache jusqu'à sa date d'expiration
- Nom du fichier de script = `{domaine}.sh` (Unix) ou `{domaine}.bat` (Windows)
- `domain` ne doit pas contenir `/` ou `\`
