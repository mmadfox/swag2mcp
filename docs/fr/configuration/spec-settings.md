# Paramètres de spécification

Les paramètres de spécification définissent un service API et remplacent les paramètres globaux pour cette API spécifique. Chaque spécification représente une API logique (par exemple, « Open-Meteo Weather APIs ») et peut contenir plusieurs collections (fichiers de spécification).

## Section de spécification

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    llm_instruction: "Utilisez cette API pour les prévisions météorologiques et les données climatiques"
    base_url: https://api.open-meteo.com
    disable: false
    tags: ["weather", "climate"]
    http_client:
      timeout: 10s
      max_response_size: 1024
    auth:
      type: bearer
      config:
        token: "$(TOKEN)"
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

## Paramètres

### domain

- **Type :** `string`
- **Obligatoire :** Oui
- **Description :** Identifiant unique pour cette spécification API. Utilisé en interne pour référencer la spécification.
- **Règles :** 1 à 60 caractères. Uniquement lettres minuscules (`a-z`), chiffres (`0-9`), traits d'union (`-`) et underscores (`_`).
- **Exemple :** `meteo`, `binance`, `mon-api`

### llm_title

- **Type :** `string`
- **Obligatoire :** Oui
- **Description :** Nom lisible que le LLM utilise pour référencer cette API. Affiché dans les réponses des outils MCP.
- **Règles :** 5 à 120 caractères. Uniquement lettres, chiffres, espaces et ponctuation de base.
- **Exemple :** `Open-Meteo Weather APIs`, `Binance Market Data`

### llm_instruction

- **Type :** `string`
- **Valeur par défaut :** `""`
- **Description :** Instructions pour le LLM sur la façon d'utiliser cette API. Décrit ce que fait l'API et quand l'utiliser.
- **Règles :** 500 caractères maximum. Uniquement lettres, chiffres, espaces et ponctuation de base.
- **Exemple :** `"Utilisez cette API pour les prévisions météorologiques, les conditions actuelles et les données climatiques."`

### base_url

- **Type :** `string`
- **Obligatoire :** Oui
- **Description :** URL de base pour toutes les requêtes API de cette spécification. Les chemins des points de terminaison de la spécification OpenAPI sont ajoutés à cette URL.
- **Exemple :** `https://api.open-meteo.com`, `https://api.binance.com`
- **Remarque :** Peut être remplacé au niveau de la collection si différentes collections utilisent des URL de base différentes.

### disable

- **Type :** `bool`
- **Valeur par défaut :** `false`
- **Description :** Lorsqu'il est `true`, cette spécification est exclue des outils MCP. Elle n'est pas chargée, indexée ni disponible pour le LLM.
- **Quand l'utiliser :** Désactiver temporairement une API sans la supprimer de la configuration. Utile pour les API qui sont hors service, obsolètes ou en maintenance.

### tags

- **Type :** `[]string` (tableau de chaînes)
- **Valeur par défaut :** `[]`
- **Description :** Balises pour filtrer les spécifications. Utilisé avec l'indicateur `--tags` dans les commandes CLI (`ls`, `validate`, `mcp`, `update`).
- **Exemple :** `["public", "weather"]`, `["internal", "production"]`
- **Effet :** Lorsque vous exécutez `swag2mcp mcp --tags=public`, seules les spécifications avec la balise `public` sont chargées.

### http_client

- **Type :** `object`
- **Valeur par défaut :** hérite du global
- **Description :** Remplace les paramètres globaux du client HTTP pour cette spécification. Tous les paramètres du `http_client` global peuvent être remplacés : `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **Exemple :**
  ```yaml
  http_client:
    timeout: 60s
    max_response_size: 4194304
    headers:
      "X-DC": "us-east-1"
  ```

### auth

- **Type :** `object`
- **Valeur par défaut :** `none` (aucune authentification)
- **Description :** Configuration d'authentification pour cette spécification. Consultez la section [Authentification](/auth/overview) pour les 9 méthodes et leurs paramètres.
- **Exemple :**
  ```yaml
  auth:
    type: bearer
    config:
      token: "$(API_TOKEN)"
  ```

### collections

- **Type :** `[]object` (tableau de collections)
- **Obligatoire :** Oui (au moins 1)
- **Description :** Liste des fichiers de spécification OpenAPI/Swagger/Postman qui appartiennent à cette spécification. Chaque collection est un fichier de spécification.
- **Règles :** 1 à 30 collections par spécification.
- **Voir :** [Paramètres de collection](./collection-settings) pour tous les paramètres de collection.

## Désactivation d'une spécification

Les spécifications désactivées ne sont pas chargées ni indexées. Le LLM ne peut pas les voir ni les utiliser.

```yaml
specs:
  - domain: old-api
    llm_title: Ancienne API
    base_url: https://old-api.example.com
    disable: true
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Remplacement du client HTTP

Tous les paramètres `http_client` du niveau global peuvent être remplacés au niveau de la spécification. Les valeurs de la spécification prévalent sur les valeurs globales pour cette spécification uniquement.

```yaml
specs:
  - domain: slow-api
    llm_title: API lente
    base_url: https://slow-api.example.com
    http_client:
      timeout: 120s
      max_response_size: 8388608
      headers:
        "X-DC": "us-east-1"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Remplacement du proxy

Si cette spécification nécessite un proxy différent du proxy global, configurez-le au niveau de la spécification :

```yaml
specs:
  - domain: proxied-api
    llm_title: API avec proxy
    base_url: https://api.example.com
    http_client:
      proxy:
        url: http://proxy.company.com:8080
        username: $(PROXY_USER)
        password: $(PROXY_PASS)
        bypass:
          - "*.local"
          - "10.0.0.0/8"
    collections:
      - llm_title: Main
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo.json
```
