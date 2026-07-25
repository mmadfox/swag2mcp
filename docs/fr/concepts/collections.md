# Collections

Une collection est un fichier OpenAPI/Swagger/Postman unique qui décrit une API spécifique. Elle pointe vers un `location` (URL ou chemin de fichier local) et appartient à une spec (domaine).

Une spec peut avoir plusieurs collections — par exemple, la spec « meteo » pourrait avoir les collections « Prévisions », « Qualité de l'air » et « Maritime », chacune pointant vers un fichier de spécification différent.

## Champs de la collection

| Champ | Clé YAML | Requis | Description |
|-------|----------|--------|-------------|
| [Titre LLM](#instruction-llm) | `llm_title` | ❌ | Nom d'affichage de la collection pour le LLM (max 120 caractères). Rempli automatiquement à partir du document de spec si non défini |
| [Instruction LLM](#instruction-llm) | `llm_instruction` | ❌ | Indice court pour le LLM (max 360 caractères). Rempli automatiquement à partir du document de spec si non défini |
| Titre | `title` | ❌ | Remplacement du titre original de la spec (rempli automatiquement à partir du document analysé) |
| [Emplacement](#emplacement--comment-les-fichiers-de-specification-sont-resolus) | `location` | ✅ | URL ou chemin vers le fichier de spécification (5–250 caractères) |
| [Désactiver](#desactiver) | `disable` | ❌ | Ignorer cette collection lors du chargement |
| [Client HTTP](#remplacement-du-client-http) | `http_client` | ❌ | Paramètres HTTP par collection (en-têtes, cookies) |
| [URL de base](#remplacement-de-lurl-de-base) | `base_url` | ❌ | Remplacer l'URL de base de la spec pour cette collection |
| [Serveur mock](#serveur-mock) | `base_mock_url` | ❌ | Adresse du serveur mock au format `hôte:port`. Requis lorsque `mock_enabled: true` |

## Emplacement — Comment les fichiers de spécification sont résolus

Le champ `location` indique à swag2mcp où trouver le fichier OpenAPI/Swagger/Postman. Il prend en charge plusieurs types de sources :

| Source | Exemple | Description |
|--------|---------|-------------|
| **URL distante** | `https://raw.githubusercontent.com/.../spec.yaml` | Téléchargé et mis en cache |
| **Fichier local (absolu)** | `/home/utilisateur/mon-api.yaml` | Lu depuis le système de fichiers, mis en cache |
| **Fichier local (relatif)** | `./mon-api.yaml` | Résolu en chemin absolu, mis en cache |
| **Fichier local de l'espace de travail** | `specs/mon-api.yaml` | Stocké dans `~/.swag2mcp/specs/`, utilisé directement (non mis en cache) |
| **URI file://** | `file:///home/utilisateur/spec.yaml` | Converti en chemin local, mis en cache |

swag2mcp détecte automatiquement le type de source :

- `https://` ou `http://` → URL distante (mise en cache)
- `file://` → fichier local (converti en chemin du système de fichiers)
- Tout le reste → fichier local (avec expansion `~` pour le répertoire personnel)

### URL distantes

Lorsque vous utilisez une URL distante, swag2mcp télécharge le fichier et le met en cache localement. Le cache est réutilisé lors des démarrages suivants pour éviter des téléchargements répétés.

### Fichiers locaux

Les fichiers locaux sont lus directement depuis le système de fichiers. Si le fichier est en dehors du répertoire `specs/` de l'espace de travail, il est copié dans le cache pour des raisons de cohérence.

### Fichiers locaux de l'espace de travail

Le répertoire `specs/` dans l'espace de travail (`~/.swag2mcp/specs/`) est l'emplacement recommandé pour les fichiers de spécification locaux. Les fichiers stockés ici sont utilisés directement sans mise en cache. Utilisez un chemin relatif commençant par `specs/` pour les référencer.

> **Remarque :** `specs/` est simplement un nom de répertoire (comme `cache/` ou `responses/`), pas le concept de « spec ». Il stocke les fichiers OpenAPI/Swagger/Postman réels vers lesquels les collections pointent.

```bash
# Importer un fichier de spécification dans l'espace de travail
swag2mcp import https://example.com/api.yaml maspec

# Après l'importation, l'emplacement devient :
# specs/maspec.yaml
```

## Système de cache

swag2mcp met en cache les fichiers de spécification distants pour éviter de les télécharger à chaque démarrage.

### Comment cela fonctionne

1. Lorsqu'une collection avec une URL distante est chargée, swag2mcp vérifie le cache
2. Si une entrée de cache valide (non expirée) existe, elle est utilisée directement
3. Sinon, le fichier est téléchargé, analysé et stocké dans le cache

### Structure du cache

```
~/.swag2mcp/
  cache/
    {sha256_hash}.spec    # Contenu du fichier de spécification en cache
    {sha256_hash}.meta    # Métadonnées du cache (JSON)
```

Chaque fichier en cache a un fichier de métadonnées contenant :

```json
{
  "source": "https://example.com/api.yaml",
  "source_type": "url",
  "cached_at": "2024-01-01T00:00:00Z",
  "mod_time": "2024-01-01T00:00:00Z",
  "ttl_sec": 3600
}
```

### TTL du cache

Chaque fichier en cache reçoit un **TTL aléatoire** entre 1 heure et 48 heures. Cela empêche tous les fichiers en cache d'expirer en même temps (problème de ruée).

### Clé de cache

La clé de cache est un hachage SHA-256 de la chaîne d'emplacement brute (16 premiers octets = 32 caractères hex).

### Gestion du cache

```bash
# Vider le cache et les réponses, retélécharger tous les fichiers de spécification
swag2mcp update

# Vider le cache et les réponses uniquement
swag2mcp clean
```

- `swag2mcp update` — valide la configuration, vide `cache/` et `responses/`, puis re-met en cache tous les emplacements de collection
- `swag2mcp clean` — supprime tout le contenu de `cache/` et `responses/`, plus les scripts d'authentification orphelins
- Les anciennes réponses sont nettoyées automatiquement après 48 heures au démarrage du serveur MCP

## Validation

Chaque collection est validée lorsque la configuration est chargée. La validation s'exécute à chaque démarrage de `swag2mcp mcp`. Si elle échoue, le serveur MCP ne démarrera pas — dans certains IDE, cela signifie que le serveur ne se connectera tout simplement pas, et le LLM reçoit un message d'erreur clair expliquant quoi corriger.

| Vérification | Règle |
|--------------|-------|
| **Emplacement** | Requis, 5–250 caractères |
| **Accessibilité de l'emplacement** | Doit être une URL accessible ou un fichier existant |
| **Validité de l'emplacement** | Doit être un fichier OpenAPI 3.x, Swagger 2.0 ou Postman valide |
| **Titre LLM** | Max 120 caractères, lettres/chiffres/ponctuation de base |
| **Instruction LLM** | Max 360 caractères, même jeu de caractères que le titre |
| **URL de base** | Doit être une URL valide si définie |
| **URL mock de base** | Doit être `hôte:port` ou `hôte:port/chemin` où l'hôte est `localhost`, `127.0.0.1` ou `0.0.0.0` |
| **Mock requis** | Si `mock_enabled: true`, chaque collection doit avoir `base_mock_url` |
| **Ports mock en double** | Deux collections ne peuvent pas partager le même port mock |

Pour diagnostiquer les problèmes avant de démarrer le serveur, utilisez la commande [`validate`](../cli/validate.md) :

```bash
# Valider l'espace de travail par défaut (~/.swag2mcp)
swag2mcp validate

# Valider un espace de travail de projet personnalisé
swag2mcp validate ./mon-projet
```

## Ajout de collections

### Via la configuration YAML

Modifiez `~/.swag2mcp/swag2mcp.yaml` directement :

```yaml
specs:
  - domain: meteo
    llm_title: API Météo Open-Meteo
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Prévisions
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
```

Après modification, redémarrez le serveur MCP (`swag2mcp mcp`) pour que les changements prennent effet.

### Via la CLI

```bash
# Mode interactif
swag2mcp add collection

# Non interactif avec YAML
swag2mcp add collection --yaml 'spec_domain: meteo
llm_title: Prévisions
location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml'

# Pipe depuis stdin
cat collection.yaml | swag2mcp add collection --yaml -

# Afficher l'exemple YAML
swag2mcp add collection --example
```

### Via l'importation

```bash
# Importer un fichier de spécification dans l'espace de travail
swag2mcp import https://example.com/api.yaml
```

## Instruction LLM

Les collections peuvent avoir leur propre `llm_instruction` (jusqu'à 360 caractères) pour des conseils plus spécifiques. Ceci est injecté dans l'invite système swag2mcp en plus de l'instruction au niveau de la spec.

```yaml
specs:
  - domain: meteo
    llm_title: API Météo Open-Meteo
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Prévisions
        llm_instruction: "Utilisez cette collection pour la météo actuelle et les prévisions quotidiennes."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Qualité de l'air
        llm_instruction: "Utilisez cette collection pour l'indice de qualité de l'air et les données de pollution."
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
```

Si `llm_title` n'est pas défini, il est automatiquement rempli à partir du champ `title` du document de spécification. Si `llm_instruction` n'est pas défini, il est rempli à partir du champ `description` du document de spécification.

## Désactiver

Définissez `disable: true` pour ignorer une collection. Elle ne sera pas chargée, indexée ni disponible pour le LLM.

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
        disable: true
```

## Remplacement de l'URL de base

Chaque collection peut remplacer la `base_url` de la spec. C'est utile lorsque différentes collections dans la même spec utilisent des points d'accès API différents.

```yaml
specs:
  - domain: meteo
    llm_title: API Météo Open-Meteo
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Prévisions
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Qualité de l'air
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Maritime
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

## Remplacement du client HTTP

Les collections peuvent remplacer les paramètres HTTP (en-têtes, cookies) des niveaux spec et global.

```yaml
specs:
  - domain: meteo
    llm_title: API Météo Open-Meteo
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Prévisions
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        http_client:
          headers:
            X-API-Version: "2"
          cookies:
            - name: session
              value: abc123
```

Les paramètres en cascade : global → spec → collection. Voir [Cascade de configuration](../configuration/cascade.md) pour plus de détails.

## Serveur mock

Lorsque `mock_enabled: true` est défini au niveau de la configuration, chaque collection doit avoir `base_mock_url` défini. Cela indique à swag2mcp où le serveur mock s'exécute pour cette collection.

```yaml
mock_enabled: true
specs:
  - domain: meteo
    llm_title: API Météo Open-Meteo
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Prévisions
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        base_mock_url: localhost:8080
```

Voir [Serveur mock](../advanced/mock-server.md) pour tous les détails.

## Exemples

### Collection minimale

```yaml
specs:
  - domain: dadjokes
    llm_title: API Dad Joke
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Blagues
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

### Collection complète avec tous les champs

```yaml
specs:
  - domain: meteo
    llm_title: API Météo Open-Meteo
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Prévisions
        llm_instruction: "Utiliser pour la météo actuelle et les prévisions quotidiennes."
        title: "Titre personnalisé"
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
        disable: false
        base_url: https://forecast-api.open-meteo.com
        base_mock_url: localhost:8080
        http_client:
          headers:
            X-Custom: valeur
```

### Plusieurs collections par spec

```yaml
specs:
  - domain: meteo
    llm_title: API Météo Open-Meteo
    base_url: https://api.open-meteo.com
    collections:
      - llm_title: Prévisions
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml
      - llm_title: Qualité de l'air
        base_url: https://air-quality-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/air-quality.yml
      - llm_title: Maritime
        base_url: https://marine-api.open-meteo.com
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/marine.yml
```

### Fichier local dans l'espace de travail (répertoire specs/)

```yaml
specs:
  - domain: monapi
    llm_title: Mon API Interne
    base_url: https://api.maentreprise.com
    collections:
      - llm_title: Utilisateurs
        location: specs/users.openapi.json
      - llm_title: Commandes
        location: specs/orders.openapi.json
```

### Collection désactivée

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
        disable: true
```

## Associé

- [Paramètres de collection (config)](../configuration/collection-settings.md) — référence YAML complète
- [Cascade de configuration](../configuration/cascade.md) — comment les paramètres se remplacent mutuellement
- [Specs](./specs) — conteneurs logiques pour les collections
- [Client HTTP](../configuration/http-client.md) — configuration du client HTTP
- [Serveur mock](../advanced/mock-server.md) — configuration du serveur mock
- [CLI : validate](../cli/validate.md) — référence de la commande validate
- [CLI : update](../cli/update.md) — référence de la commande update
- [CLI : clean](../cli/clean.md) — référence de la commande clean
