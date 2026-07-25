# Authentification

## Aperçu

swag2mcp prend en charge **9 méthodes d'authentification** pour travailler avec des API qui nécessitent une autorisation. Vous la configurez une fois dans le fichier de configuration — après cela, chaque appel API via `invoke` inclut automatiquement les bons jetons et en-têtes.

### Où configurer

L'authentification est définie au niveau de la **spec** dans `swag2mcp.yaml` :

```yaml
specs:
  - domain: jokes
    llm_title: API Dad Joke
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Blagues
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: bearer
      config:
        token: "mon-jeton"
```

### Comment cela fonctionne

- Vous spécifiez le type d'authentification et les paramètres dans la configuration
- swag2mcp les applique automatiquement à chaque requête lorsque vous appelez `invoke`
- Vous **n'avez pas besoin** de demander un jeton avant d'appeler une API — cela se fait automatiquement
- Si un jeton expire (OAuth2, Script), swag2mcp le renouvelle automatiquement

### Variables d'environnement

Les données sensibles (jetons, mots de passe, clés) peuvent être stockées dans des variables d'environnement en utilisant la syntaxe `$(NOM_VAR)` :

```yaml
auth:
  type: bearer
  config:
    token: "$(MON_JETON_API)"
```

swag2mcp substitue la valeur de `MON_JETON_API` au démarrage.

### Outil MCP auth

L'agent LLM peut récupérer un jeton ou des en-têtes via l'outil MCP `auth` — par exemple, pour construire une commande curl ou les montrer à l'utilisateur.

En **production**, cet outil doit être désactivé avec `--disable-llm-auth` (activé par défaut) afin que le LLM n'ait jamais accès aux jetons.

### Méthodes

| Méthode | Description | Meilleur pour |
|---------|-------------|---------------|
| [`none`](/auth/none) | Aucune authentification | API publiques |
| [`basic`](/auth/basic) | HTTP Basic (nom d'utilisateur + mot de passe) | API legacy, auth simple |
| [`bearer`](/auth/bearer) | Jeton Bearer (JWT, jeton) | API REST modernes |
| [`api-key`](/auth/api-key) | Clé API dans l'en-tête ou le paramètre de requête | Services avec clés API |
| [`digest`](/auth/digest) | HTTP Digest (nom d'utilisateur + mot de passe) | API legacy, plus sécurisé que Basic |
| [`hmac`](/auth/hmac) | Signature HMAC-SHA256 (style Binance) | Échanges de cryptomonnaies |
| [`oauth2-cc`](/auth/oauth2-cc) | OAuth2 Client Credentials | Serveur à serveur, microservices |
| [`oauth2-pwd`](/auth/oauth2-pwd) | OAuth2 Password Grant | Applications avec connexion utilisateur |
| [`script`](/auth/script) | Script externe pour obtenir un jeton | Tout schéma d'authentification personnalisé |
