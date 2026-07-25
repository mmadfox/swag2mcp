# Paramètres de collection

Les paramètres de collection définissent un seul fichier de spécification OpenAPI/Swagger/Postman et remplacent les paramètres de spécification pour ce fichier spécifique. Chaque collection appartient à une spécification et représente un document de spécification API.

## Section de collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        llm_instruction: "Utiliser pour les données météorologiques actuelles et prévisionnelles"
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8081
        http_client:
          timeout: 5s
```

## Paramètres

### llm_title

- **Type :** `string`
- **Obligatoire :** Non
- **Description :** Nom lisible pour cette collection. Affiché dans les réponses des outils MCP.
- **Règles :** 120 caractères maximum. Uniquement lettres, chiffres, espaces et ponctuation de base.
- **Exemple :** `Forecast`, `Air Quality`, `Market Data`

### llm_instruction

- **Type :** `string`
- **Valeur par défaut :** `""`
- **Description :** Instructions pour le LLM concernant cette collection spécifique. Décrit les points de terminaison que cette collection fournit.
- **Règles :** 360 caractères maximum. Uniquement lettres, chiffres, espaces et ponctuation de base.
- **Exemple :** `"Utiliser pour les données météorologiques actuelles et prévisionnelles."`

### title

- **Type :** `string`
- **Valeur par défaut :** `""`
- **Description :** Titre brut provenant du fichier de spécification. Rempli automatiquement à l'exécution. Vous n'avez généralement pas besoin de le définir dans le YAML.

### location

- **Type :** `string`
- **Obligatoire :** Oui
- **Description :** URL ou chemin de fichier local vers le fichier de spécification OpenAPI 3.x, Swagger 2.0 ou Postman.
- **Règles :** 5 à 250 caractères.
- **Exemples :**
  - URL : `https://raw.githubusercontent.com/org/repo/main/spec.yaml`
  - Local : `./specs/mon-api.json`
  - Local (absolu) : `/home/utilisateur/.swag2mcp/specs/mon-api.yaml`

### disable

- **Type :** `bool`
- **Valeur par défaut :** `false`
- **Description :** Lorsqu'il est `true`, cette collection est exclue des outils MCP. Elle n'est pas chargée ni indexée.
- **Quand l'utiliser :** Désactiver temporairement une collection sans la supprimer de la configuration. Utile lorsqu'un fichier de spécification est en cours de mise à jour ou qu'une version d'API est obsolète.

### http_client

- **Type :** `object`
- **Valeur par défaut :** hérite de la spécification (ou du global)
- **Description :** Remplace les paramètres du client HTTP pour cette collection. Tous les paramètres du `http_client` global peuvent être remplacés : `timeout`, `max_response_size`, `user_agent`, `follow_redirects`, `max_redirects`, `random`, `proxy`, `headers`, `cookies`.
- **Exemple :**
  ```yaml
  http_client:
    timeout: 120s
    headers:
      "X-Custom": "valeur"
    cookies:
      - name: "session"
        value: "abc123"
  ```

### base_url

- **Type :** `string`
- **Valeur par défaut :** `""` (hérite de la spécification)
- **Description :** Remplace le `base_url` du niveau spécification pour cette collection. À utiliser lorsque différentes collections au sein d'une même spécification utilisent des URL de base différentes.
- **Exemple :** Si la spécification a `base_url: https://api.open-meteo.com` mais qu'une collection utilise `https://air-quality-api.open-meteo.com`, définissez `base_url` au niveau de la collection.

### base_mock_url

- **Type :** `string`
- **Valeur par défaut :** `""`
- **Description :** Adresse du serveur de simulation au format `host:port`. Requis lorsque `mock_enabled: true` dans la configuration globale.
- **Règles :** L'hôte doit être `localhost`, `127.0.0.1` ou `0.0.0.0`. Le port doit être un numéro de port valide.
- **Exemple :** `localhost:8081`, `127.0.0.1:9000`
- **Quand l'utiliser :** Vous avez `mock_enabled: true` et vous souhaitez tester cette collection avec des réponses factices.

## Plusieurs collections à partir d'une même spécification

Une spécification peut avoir plusieurs collections — par exemple, lorsqu'une API a des fichiers de spécification distincts pour différents services :

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Marine
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## Désactivation d'une collection

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Air Quality
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
        disable: true
```

## Remplacement du client HTTP

Tous les paramètres `http_client` peuvent être remplacés au niveau de la collection. Les valeurs de la collection prévalent sur les valeurs de la spécification et du global pour cette collection uniquement.

```yaml
specs:
  - domain: meteo
    llm_title: Open-Meteo Weather APIs
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Forecast
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          timeout: 120s
          headers:
            "X-Custom": "valeur"
          cookies:
            - name: "session"
              value: "abc123"
```
