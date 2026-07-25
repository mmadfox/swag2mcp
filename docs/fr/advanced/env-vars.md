# Variables d'environnement

## Aperçu

swag2mcp prend en charge la substitution de variables d'environnement dans le fichier de configuration en utilisant la syntaxe `$(NOM_VAR)`. Cela vous permet de garder les données sensibles (jetons, mots de passe, clés) hors du fichier YAML.

## Comment cela fonctionne

Lorsque swag2mcp démarre, il analyse la configuration pour les motifs `$(NOM_VAR)` et les remplace par la valeur de la variable d'environnement correspondante.

```yaml
specs:
  - domain: meteo
    auth:
      type: bearer
      config:
        token: "$(JETON_API)"
```

Si la variable d'environnement `JETON_API` est définie, elle sera substituée. Si elle n'est pas définie, la valeur devient vide.

## Où `$(VAR)` est résolu

| Champ | Exemple |
|-------|---------|
| Jeton `token` (bearer) | `token: "$(JETON_API)"` |
| `username` / `password` (basic, digest) | `password: "$(MOT_DE_PASSE_API)"` |
| `client_id` / `client_secret` (oauth2-cc, oauth2-pwd) | `client_secret: "$(SECRET_OAUTH)"` |
| `api_key` / `secret_key` (hmac) | `api_key: "$(CLE_API_BINANCE)"` |
| `domain` (script) | `domain: "$(DOMAINE_AUTH)"` |
| Jeton du serveur MCP | `token: "$(JETON_MCP)"` |
| En-têtes du client HTTP | `"X-API-Key": "$(CLE_API)"` |
| Valeurs de cookies du client HTTP | `value: "$(JETON_SESSION)"` |

## Où `$(VAR)` n'est PAS résolu

- URL de base (`base_url`)
- Emplacements de collections (`location`)
- Noms de domaine de spec (`domain`)

## Exemple

```bash
export JETON_API="eyJhbGciOiJIUzI1NiIs..."
export JETON_MCP="mon-jeton-secret"

swag2mcp mcp
```

## Bonnes pratiques de sécurité

- **Ne stockez jamais** les secrets directement dans le fichier YAML
- Utilisez des variables d'environnement ou un gestionnaire de secrets externe
- Ajoutez le fichier YAML à `.gitignore` s'il contient des secrets codés en dur
- Définissez les variables d'environnement dans votre profil shell, la configuration de votre IDE ou votre pipeline de déploiement

## Détails de syntaxe

- `$(NOM_VAR)` — syntaxe standard
- `$( NOM_VAR )` — les espaces entre parenthèses sont autorisés et supprimés
- `$()` — un nom de variable vide retourne la chaîne d'origine inchangée
- Les motifs `$(...)` imbriqués ne sont pas résolus
