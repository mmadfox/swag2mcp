# Client HTTP

swag2mcp utilise un client HTTP configurable pour tous les appels API. Ces paramètres sont définis globalement et peuvent être remplacés aux niveaux de la spécification et de la collection.

## Configuration

```yaml
http_client:
  timeout: 30s
  max_response_size: 1048576
  user_agent: "swag2mcp-global/1.0"
  follow_redirects: true
  max_redirects: 10
  random: false
  proxy:
    url: ""
    username: ""
    password: ""
    bypass: []
  headers:
    "Accept": "application/json"
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
```

## Délai d'attente

Contrôle le temps d'attente de swag2mcp pour une réponse API avant d'abandonner.

- **Type :** durée (format Go : `30s`, `60s`, `2m`)
- **Valeur par défaut :** `30s`
- **Plage :** 1 seconde à 5 minutes
- **Effet :** Si l'API ne répond pas dans ce délai, la requête échoue avec une erreur de délai d'attente.
- **Quand l'augmenter :** API lentes, charges utiles volumineuses, réseaux peu fiables.
- **Quand le diminuer :** API internes, vérifications de santé, scénarios d'échec rapide.

```yaml
http_client:
  timeout: 60s
```

## Taille maximale de réponse

Limite la taille d'une réponse avant que swag2mcp ne l'enregistre sur le disque au lieu de la renvoyer en ligne au LLM.

- **Type :** `int` (octets)
- **Valeur par défaut :** `1048576` (1 Mo)
- **Plage :** 256 à 10 485 760 octets (10 Mo)
- **Effet :** Lorsqu'une réponse dépasse cette limite, elle est enregistrée dans `{workspace}/responses/` sous forme de fichier JSON. Le LLM reçoit une référence de fichier et peut l'explorer avec les outils `response_outline`, `response_compress` et `response_slice`.
- **Quand l'augmenter :** API qui renvoient de grands ensembles de données (rapports, journaux, analyses).
- **Quand la diminuer :** Fenêtre de contexte LLM limitée, ou lorsque vous préférez un accès par fichier pour toutes les réponses.

```yaml
http_client:
  max_response_size: 4194304  # 4 Mo
```

## User-Agent

L'en-tête `User-Agent` envoyé avec chaque requête. Certaines API nécessitent un user-agent spécifique ou bloquent les user-agents de robots connus.

- **Type :** `string`
- **Valeur par défaut :** `"swag2mcp-global/1.0"`
- **Effet :** Identifie votre application auprès du serveur API.
- **Quand le modifier :** L'API nécessite un user-agent spécifique, ou vous souhaitez identifier votre application pour les analyses.

```yaml
http_client:
  user_agent: "MonApp/1.0"
```

## Suivi des redirections

Contrôle si swag2mcp suit automatiquement les redirections HTTP (codes d'état 3xx).

- **Type :** `bool`
- **Valeur par défaut :** `true`
- **Effet :** Lorsqu'il est `true`, swag2mcp suit les redirections jusqu'à `max_redirects` fois. Lorsqu'il est `false`, la réponse de redirection est renvoyée telle quelle.
- **Quand le désactiver :** API qui redirigent en boucle, points de terminaison sensibles à la sécurité où vous souhaitez inspecter manuellement les cibles de redirection.

```yaml
http_client:
  follow_redirects: false
```

## Nombre maximal de redirections

Limite le nombre de redirections que swag2mcp suit avant de s'arrêter.

- **Type :** `int`
- **Valeur par défaut :** `10`
- **Plage :** 0 à 50
- **Effet :** Si l'API redirige plus de fois que cette limite, la requête échoue.
- **Quand le modifier :** API avec de longues chaînes de redirection, ou réduisez pour un échec plus rapide en cas de boucles de redirection.

```yaml
http_client:
  max_redirects: 5
```

## Randomiseur

Ajoute des en-têtes aléatoires de type navigateur à chaque requête pour éviter le pistage et le blocage.

- **Type :** `bool`
- **Valeur par défaut :** `false`
- **Effet :** Lorsqu'il est `true`, swag2mcp génère des en-têtes aléatoires pour chaque requête : `User-Agent` (à partir d'un pool de chaînes de navigateurs réels), `Accept`, `Accept-Language`, `Accept-Encoding`, `Cache-Control`. Cela remplace le paramètre `user_agent`.
- **Quand l'activer :** API qui bloquent les requêtes en fonction du User-Agent ou des modèles d'en-têtes, scénarios de collecte de données.

```yaml
http_client:
  random: true
```

## Proxy

Un serveur proxy agit comme intermédiaire entre swag2mcp et l'API cible. Tout le trafic HTTP est acheminé à travers lui.

**Quand vous pourriez avoir besoin d'un proxy :**
- **Réseau d'entreprise** — tout le trafic sortant doit passer par un proxy d'entreprise
- **Restrictions géographiques** — certaines API sont verrouillées par région, un proxy dans la bonne région contourne cela
- **IP statique** — API qui nécessitent une liste blanche d'adresses IP
- **Anonymat** — masquer l'adresse IP d'origine auprès de l'API cible

### URL du proxy

- **Type :** `string`
- **Valeur par défaut :** `""` (pas de proxy)
- **Schémas pris en charge :** `http`, `https`, `socks5`, `socks5h`
- **Prend en charge `$(VAR)` :** ✅ résolu à l'exécution

| Schéma | Description | Cas d'utilisation |
|--------|-------------|----------|
| `http` | Proxy HTTP pour le trafic HTTP | Proxies d'entreprise, proxy de base |
| `https` | Proxy HTTPS (tunnel CONNECT) | Proxies d'entreprise sécurisés |
| `socks5` | Proxy SOCKS5 (DNS résolu localement) | Usage général, tout protocole |
| `socks5h` | Proxy SOCKS5 (DNS résolu sur le proxy) | Lorsque le proxy a une meilleure résolution DNS |

### Authentification du proxy

Si le proxy nécessite une authentification, fournissez `username` et `password` :

- **Prend en charge `$(VAR)` :** ✅ résolu à l'exécution pour les trois champs (`url`, `username`, `password`)

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    username: "utilisateurproxy"
    password: "$(PROXY_PASSWORD)"
```

### Contournement du proxy

Une liste de domaines qui ne doivent **pas** passer par le proxy. Utile pour les services internes, localhost ou les API qui ne sont accessibles que directement.

```yaml
http_client:
  proxy:
    url: "http://proxy.example.com:8080"
    bypass:
      - "localhost"
      - "127.0.0.1"
      - "*.internal.company.com"
      - "api.local"
```

Le contournement prend en charge les motifs génériques (`*.example.com` correspond à tout sous-domaine).

## En-têtes

En-têtes HTTP personnalisés ajoutés à chaque requête. Les en-têtes sont fusionnés entre les niveaux de cascade :

```
En-têtes globaux → En-têtes de spécification (fusionnés) → En-têtes de collection (fusionnés)
```

Les en-têtes de collection remplacent les en-têtes de spécification, qui remplacent les en-têtes globaux pour la même clé.

```yaml
http_client:
  headers:
    "Accept": "application/json"
    "Accept-Language": "fr-FR"
```

Les valeurs d'en-tête prennent en charge la résolution `$(ENV_VAR)`.

## Cookies

Cookies envoyés avec chaque requête. Les cookies sont fusionnés entre les niveaux de cascade (le niveau inférieur remplace le global pour le même nom de cookie).

```yaml
http_client:
  cookies:
    - name: "session"
      value: "abc123"
      domain: ".example.com"
      path: "/"
      secure: false
      http_only: false
```

### Champs de cookie

| Champ | Obligatoire | Description |
|-------|----------|-------------|
| `name` | Oui | Nom du cookie |
| `value` | Oui | Valeur du cookie (prend en charge la résolution `$(ENV_VAR)`) |
| `domain` | Non | Domaine du cookie (par exemple, `.example.com`) |
| `path` | Non | Chemin du cookie (par exemple, `/`) |
| `secure` | Non | Envoyer uniquement via HTTPS |
| `http_only` | Non | Non accessible via JavaScript |

## En-têtes personnalisés au niveau de la spécification

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    http_client:
      headers:
        "Accept": "application/json"
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Cookies au niveau de la spécification

```yaml
specs:
  - domain: example
    llm_title: API exemple
    base_url: https://api.example.com
    http_client:
      cookies:
        - name: "session"
          value: "abc123"
        - name: "csrf"
          value: "$(CSRF_TOKEN)"
    collections:
      - llm_title: Default
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
```

## Cascade

Les paramètres du client HTTP se propagent en cascade du global à la spécification puis à la collection. Tous les paramètres peuvent être remplacés à chaque niveau :

```
Global (http_client)
    ↓ remplace (tous les paramètres)
Spécification (specs[].http_client)
    ↓ remplace (tous les paramètres)
Collection (specs[].collections[].http_client)
```

**Tous les paramètres du client HTTP** (délai d'attente, proxy, user-agent, redirections, taille de réponse, randomiseur, en-têtes, cookies) peuvent être remplacés aux niveaux spécification et collection.

Consultez [Cascade de configuration](./cascade) pour plus de détails.
