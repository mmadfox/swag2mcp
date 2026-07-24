# Specs

Une spec est un conteneur logique représentant un domaine ou service API (par exemple, YouTube, Binance, Open-Meteo). Chaque spec a un `domain` unique, une `base_url`, une `auth` optionnelle, et contient une ou plusieurs collections.

Les [collections](./collections) pointent vers des fichiers OpenAPI/Swagger/Postman — la spec elle-même n'est pas un fichier, c'est le regroupement qui les entoure.

## Domaine — Règles de nommage

Le `domain` est l'identifiant unique d'une spec. Il est utilisé comme clé primaire dans tout le système.

| Règle | Contrainte |
|-------|------------|
| Caractères | `a-z`, `0-9`, `_`, `-` uniquement |
| Longueur | 1–60 caractères |
| Unicité | **Aucun doublon autorisé** — deux specs actives ne peuvent pas partager le même domaine |

**Exemples valides :** `meteo`, `binance`, `github-api`, `mon_service`, `openai-v1`

**Exemples invalides :** `Meteo` (majuscule), `mon api` (espace), `mon.api` (point), `un-nom-de-domaine-tres-long-qui-depasse-soixante-caracteres` (trop long)

## Champs de la spec

| Champ | Clé YAML | Requis | Description |
|-------|----------|--------|-------------|
| [Domaine](#domaine--regles-de-nommage) | `domain` | ✅ | Identifiant API unique (1–60 car., `a-z0-9_-`) |
| Titre LLM | `llm_title` | ✅ | Nom lisible par l'humain que le LLM utilise pour référencer cette API (5–120 car.) |
| [Instruction LLM](#instruction-llm) | `llm_instruction` | ❌ | Indice court injecté dans l'invite système swag2mcp (max 500 car.) |
| URL de base | `base_url` | ✅ | URL de base pour toutes les requêtes API (URL valide) |
| [Désactiver](#desactiver) | `disable` | ❌ | Ignorer cette spec lors du chargement et de l'indexation |
| [Étiquettes](#etiquettes) | `tags` | ❌ | Étiquettes pour le filtrage (par ex., `["public", "demo"]`) |
| [Auth](#auth) | `auth` | ❌ | Configuration de l'authentification |
| [Client HTTP](#client-http) | `http_client` | ❌ | Paramètres HTTP par spec (en-têtes, cookies) |
| [Collections](./collections) | `collections` | ✅ | Liste de 1–30 collections |

## Validation

Lorsque swag2mcp valide la configuration, ces règles sont vérifiées pour chaque spec :

| Vérification | Règle |
|--------------|-------|
| **Domaines en double** | Deux specs actives ne peuvent pas partager le même `domain` |
| **Format du domaine** | Doit correspondre à `^[a-z0-9_-]{1,60}$` |
| **Titre LLM** | Requis, 5–120 caractères, lettres/chiffres/espaces/ponctuation de base |
| **Instruction LLM** | Max 500 caractères, même jeu de caractères que le titre |
| **URL de base** | Requis, doit être une URL valide |
| **Collections** | Requis, 1–30 éléments |
| **Auth** | Validé par type d'auth (par ex., bearer nécessite `token`, basic nécessite `username` + `password`) |
| **Emplacement** | Le `location` de chaque collection doit être une URL ou un chemin de fichier valide (5–250 car.) |

La validation s'exécute à chaque démarrage de `swag2mcp mcp`. Si elle échoue, le serveur MCP ne démarrera pas — dans certains IDE, cela signifie que le serveur ne se connectera tout simplement pas, et le LLM reçoit un message d'erreur clair expliquant quoi corriger.

Pour diagnostiquer les problèmes avant de démarrer le serveur, utilisez la commande [`validate`](../cli/validate.md) :

```bash
# Valider l'espace de travail par défaut (~/.swag2mcp)
swag2mcp validate

# Valider un espace de travail de projet personnalisé
swag2mcp validate ./mon-projet
```

## Instruction LLM

Il est recommandé de définir `llm_instruction` sur chaque spec — un indice court (jusqu'à 500 caractères) qui indique au LLM à quoi sert cette API et quand l'utiliser. Cette instruction est injectée dans l'invite système swag2mcp, aidant le LLM à comprendre l'objectif de la spec sans contexte supplémentaire.

```yaml
specs:
  - domain: jokes
    llm_title: API Dad Joke
    llm_instruction: "Utilisez cette API pour obtenir des blagues de papa aléatoires ou rechercher des blagues spécifiques par mot-clé."
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Blagues
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Les collections peuvent également avoir leur propre `llm_instruction` (jusqu'à 360 caractères) pour des conseils plus spécifiques.

## Auth

L'authentification est configurée au niveau de la spec et s'applique à toutes ses collections. swag2mcp prend en charge 9 méthodes d'authentification :

| Méthode | Type YAML | Champs clés |
|---------|-----------|-------------|
| [Aucune](../auth/none.md) | `none` | — |
| [Basic](../auth/basic.md) | `basic` | `username`, `password` |
| [Bearer](../auth/bearer.md) | `bearer` | `token` |
| [Digest](../auth/digest.md) | `digest` | `username`, `password` |
| [OAuth2 Client Credentials](../auth/oauth2-cc.md) | `oauth2-cc` | `client_id`, `client_secret`, `token_url` |
| [OAuth2 Password](../auth/oauth2-pwd.md) | `oauth2-pwd` | `username`, `password`, `client_id`, `token_url` |
| [Clé API](../auth/api-key.md) | `api-key` | `key`, `value`, `in` (`header` ou `query`) |
| [HMAC](../auth/hmac.md) | `hmac` | `api_key`, `secret_key` |
| [Script](../auth/script.md) | `script` | `domain` |

Voir [Aperçu de l'authentification](../auth/overview.md) pour tous les détails sur chaque méthode.

## Client HTTP

Vous pouvez remplacer les paramètres HTTP au niveau de la spec. Ils s'appliquent à toutes les requêtes effectuées par les collections de cette spec.

```yaml
specs:
  - domain: api-lente
    llm_title: API Lente
    base_url: https://api-lente.example.com
    http_client:
      headers:
        X-API-Version: "2"
      cookies:
        - name: session
          value: abc123
    collections:
      - llm_title: Par défaut
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

Les paramètres en cascade : global → spec → collection. Voir [Cascade de configuration](../configuration/cascade.md) pour plus de détails.

## Étiquettes

Les étiquettes vous permettent de filtrer les specs par catégorie. Utilisez-les avec le drapeau `--tags` sur `swag2mcp ls` ou lors de l'amorçage.

```yaml
specs:
  - domain: meteo
    llm_title: API Météo Open-Meteo
    base_url: https://api.open-meteo.com
    tags: ["meteo", "public"]
    collections:
      - llm_title: Prévisions
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

```bash
# Lister uniquement les specs étiquetées « meteo »
swag2mcp ls --tags meteo
```

## Désactiver

Définissez `disable: true` pour ignorer complètement une spec. Elle ne sera pas chargée, indexée ni disponible pour le LLM.

```yaml
specs:
  - domain: ancienne-api
    llm_title: Ancienne API (dépréciée)
    base_url: https://ancienne-api.example.com
    disable: true
    collections:
      - llm_title: Par défaut
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Exemples

### Spec minimale

```yaml
specs:
  - domain: dadjokes
    llm_title: API Dad Joke
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Blagues
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Spec avec authentification

```yaml
specs:
  - domain: binance
    llm_title: API Données de marché Binance
    base_url: https://api.binance.com
    auth:
      type: hmac
      config:
        api_key: $(CLE_API_BINANCE)
        secret_key: $(CLE_SECRETE_BINANCE)
    collections:
      - llm_title: Données de marché
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/binance.yaml
```

### Spec avec plusieurs collections

```yaml
specs:
  - domain: meteo
    llm_title: API Météo Open-Meteo
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Prévisions
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Qualité de l'air
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Maritime
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### Spec avec instruction LLM et étiquettes

```yaml
specs:
  - domain: rickandmorty
    llm_title: API Rick et Morty
    llm_instruction: "Utilisez cette API pour obtenir des informations sur les personnages, épisodes et lieux de la série Rick et Morty."
    base_url: https://rickandmortyapi.com/api
    tags: ["divertissement", "public"]
    collections:
      - llm_title: Personnages
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/rick-and-morty.json
```

## Associé

- [Paramètres de spec (config)](../configuration/spec-settings.md) — référence YAML complète
- [Cascade de configuration](../configuration/cascade.md) — comment les paramètres se remplacent mutuellement
- [Aperçu de l'authentification](../auth/overview.md) — les 9 méthodes d'authentification
- [Client HTTP](../configuration/http-client.md) — configuration du client HTTP
- [Collections](./collections) — fichiers de spécification dans une spec
