# Serveur Mock

## Aperçu

Le serveur mock génère des réponses d'API factices basées sur vos schémas OpenAPI. Il vous permet de tester votre intégration API sans effectuer d'appels HTTP réels. C'est utile pour le développement, les tests d'agents LLM et les démonstrations.

Le serveur mock est un **binaire séparé** — `swag2mcp-mock`. Il n'est pas inclus dans le binaire principal `swag2mcp` et doit être installé séparément.

## Installation

```bash
# Option 1 : Téléchargement depuis GitHub Releases
# Cherchez swag2mcp-mock_<version>_<os>_<arch>.tar.gz

# Option 2 : Installation avec Go
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## Configuration

Activez le serveur mock dans votre configuration :

```yaml
mock_enabled: true
mock_auth:
  oauth2_port: 9090
  digest_port: 9091
  hmac_port: 9092

specs:
  - domain: jokes
    llm_title: API Dad Joke
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Blagues
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
        base_mock_url: "127.0.0.1:9090"
```

## Paramètres

### mock_enabled

- **Type :** `bool`
- **Défaut :** `false`
- **Effet :** Lorsque `true`, chaque collection active doit avoir `base_mock_url` défini. Le serveur mock démarre des serveurs HTTP pour chaque collection.

### mock_auth

Ports pour les serveurs d'authentification mock. Ils simulent les points d'accès d'authentification OAuth2, Digest et HMAC afin que vous puissiez tester des API authentifiées sans identifiants réels.

| Champ | Défaut | Description |
|-------|--------|-------------|
| `oauth2_port` | `9090` | Port pour le serveur de jeton OAuth2 mock |
| `digest_port` | `9091` | Port pour le serveur d'authentification Digest mock |
| `hmac_port` | `9092` | Port pour le serveur d'authentification HMAC mock |

### base_mock_url (par collection)

- **Type :** `string`
- **Requis :** Oui (lorsque `mock_enabled: true`)
- **Format :** `hôte:port` (par exemple, `localhost:8080`, `127.0.0.1:9000`)
- **Effet :** Chaque collection obtient son propre serveur HTTP sur cette adresse. Le serveur répond à tous les points d'accès définis dans la spec avec des données générées aléatoirement.

## Démarrage du serveur mock

```bash
# Démarrage avec la configuration par défaut
swag2mcp-mock mockserver

# Démarrage avec TLS
swag2mcp-mock mockserver --tls

# Démarrage avec un certificat TLS personnalisé
swag2mcp-mock mockserver --tls --tls-cert cert.pem --tls-key key.pem
```

### Drapeaux TLS

| Drapeau | Description |
|---------|-------------|
| `--tls` | Active TLS avec un certificat auto-signé |
| `--tls-cert` | Chemin vers le fichier de certificat TLS |
| `--tls-key` | Chemin vers le fichier de clé TLS |

Si `--tls` est défini sans `--tls-cert` et `--tls-key`, un certificat auto-signé est généré automatiquement pour `localhost`.

## Ce que fait le serveur mock

Lorsque vous démarrez le serveur mock, il :

1. **Analyse tous les fichiers de spécification** — lit la spécification OpenAPI/Swagger de chaque collection
2. **Enregistre les gestionnaires** — crée un gestionnaire HTTP pour chaque chemin et méthode définis dans la spec
3. **Génère des données factices** — répond avec des données générées aléatoirement qui correspondent au schéma de réponse (types, formats et structure corrects)
4. **Démarre les serveurs d'authentification** — simule les points d'accès d'authentification OAuth2, Digest et HMAC pour les tests

### Test du mock

```bash
# Dans un terminal :
swag2mcp-mock mockserver

# Dans un autre terminal :
curl http://localhost:8080/pets
# → [{"id":1,"name":"Pet_name","status":"available"}]
```

## Comment les données factices sont générées

Le serveur mock génère des données factices réalistes basées sur le schéma OpenAPI :

- **Chaînes** — mots, phrases ou valeurs spécifiques au format (email, URL, UUID, date, téléphone, etc.)
- **Nombres** — entiers et flottants aléatoires dans la plage spécifiée
- **Booléens** — vrai/faux aléatoire
- **Tableaux** — 1 à 3 éléments aléatoires
- **Objets** — toutes les propriétés remplies avec des valeurs aléatoires
- **Énumérations** — valeur aléatoire de la liste d'énumération
- **Champs nullables** — retourne parfois `null` (~10% de chance)

## Cas d'utilisation

- **Développement** — testez votre intégration sans accès API réel
- **Test d'agents LLM** — vérifiez que le LLM peut découvrir, inspecter et invoquer des points d'accès
- **Démonstrations** — montrez swag2mcp en fonctionnement sans configurer d'API réelles
- **Tests de charge** — testez le serveur MCP sous charge sans solliciter des API réelles

## Notes importantes

- **Binaire séparé** — `swag2mcp-mock` n'est pas inclus dans le binaire principal `swag2mcp`. Installez-le séparément.
- **Chaque collection obtient son propre port** — configurez `base_mock_url` par collection
- **Les serveurs d'authentification mock sont globaux** — les serveurs OAuth2, Digest et HMAC fonctionnent sur les ports configurés, quel que soit le nombre de collections
- **Les échecs d'analyse de spec ne sont pas fatals** — si la spec d'une collection ne peut pas être analysée, elle est ignorée avec un avertissement
- **TLS auto-signé** — lors de l'utilisation de `--tls` sans certificats, un certificat auto-signé est généré pour localhost uniquement
