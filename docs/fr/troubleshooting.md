# Dépannage

## Problèmes d'installation

### swag2mcp : commande introuvable

Le binaire n'est pas dans votre PATH.

```bash
# Vérifiez si Go est installé
go version

# Trouvez où Go installe les binaires
go env GOPATH
# Généralement ~/go ou ~/go/bin

# Ajoutez au PATH (ajoutez ceci à ~/.zshrc ou ~/.bashrc)
export PATH=$PATH:$(go env GOPATH)/bin

# Ou utilisez le chemin complet
~/go/bin/swag2mcp --version
```

Si vous avez téléchargé un binaire depuis GitHub Releases, assurez-vous qu'il se trouve dans un répertoire qui est dans votre PATH :

```bash
# Déplacez vers /usr/local/bin (macOS/Linux)
sudo mv swag2mcp /usr/local/bin/
```

### permission refusée

Le binaire n'a pas les permissions d'exécution.

```bash
# Pour go install (corriger la propriété)
sudo chown -R $(whoami) $(go env GOPATH)

# Pour un binaire téléchargé
chmod +x /chemin/vers/swag2mcp
```

### Version de Go trop ancienne

swag2mcp nécessite Go 1.23+.

```bash
go version
# Si version < 1.23, mettez à jour Go :
# https://go.dev/dl/
```

### Serveur mock introuvable

Le serveur mock est un binaire séparé. Installez-le explicitement :

```bash
go install github.com/mmadfox/swag2mcp/cmd/swag2mcp-mock@latest
```

## Problèmes de configuration

### Fichier de configuration introuvable

swag2mcp ne trouve pas `swag2mcp.yaml`.

```bash
# Créez une nouvelle configuration
swag2mcp init

# Ou spécifiez le chemin explicitement
swag2mcp mcp /chemin/vers/espace-travail
swag2mcp ls /chemin/vers/espace-travail
```

**Cause fréquente :** Vous avez exécuté `swag2mcp mcp` depuis un répertoire quelconque et il a cherché `~/.swag2mcp/` au lieu de l'espace de travail de votre projet. Passez toujours le chemin explicitement.

### Mauvais espace de travail chargé

swag2mcp a chargé un espace de travail différent de celui attendu.

**Ordre de résolution :** `[chemin]` explicite → répertoire courant (`./`) → `~/.swag2mcp/`. Si vous exécutez `swag2mcp mcp` sans chemin depuis un répertoire qui n'a pas `swag2mcp.yaml`, il utilise `~/.swag2mcp/`.

**Correctif :** Passez toujours le chemin de l'espace de travail : `swag2mcp mcp /chemin/vers/votre/espace-travail`

### Erreur d'analyse YAML

Le fichier de configuration a une syntaxe YAML invalide.

```bash
# Validez la configuration
swag2mcp validate

# Erreurs courantes :
# - Tabulations au lieu d'espaces (YAML nécessite des espaces)
# - Indentation manquante pour les champs imbriqués
# - Chaînes non citées avec des caractères spéciaux (: # & {)
```

**Astuce :** Utilisez un linter YAML ou un éditeur avec support YAML pour détecter les erreurs de syntaxe.

### La validation échoue : « aucune spécification définie »

Le fichier de configuration existe mais n'a pas de specs.

```bash
# Ajoutez une spec
swag2mcp add spec

# Ou modifiez swag2mcp.yaml et ajoutez au moins une spec
```

### La validation échoue : « domaine en double »

Deux specs ont la même valeur `domain`. Les domaines doivent être uniques.

```bash
# Listez les specs actuelles
swag2mcp ls

# Vérifiez les domaines en double dans swag2mcp.yaml
```

### La validation échoue : « emplacement de spec invalide »

L'URL ou le chemin de fichier `location` n'est pas accessible ou n'est pas un fichier de spécification valide.

```bash
# Vérifiez si l'URL est accessible
curl -I https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml

# Vérifiez si le fichier local existe
ls -la ./specs/mon-api.yaml

# Vérifiez que le fichier est un OpenAPI/Swagger/Postman valide
# (pas n'importe quel JSON ou page HTML)
```

**Cause fréquente :** Le champ `location` pointe vers le point d'accès API lui-même (par exemple, `https://api.example.com/v1/users`) au lieu de l'URL du fichier de spécification. L'emplacement doit pointer vers un fichier OpenAPI/Swagger/Postman.

## Problèmes du serveur MCP

### Port déjà utilisé

Un autre processus utilise le port.

```bash
# Trouvez le processus
lsof -i :8080

# Tuez-le
kill <PID>

# Ou utilisez un port différent
swag2mcp mcp --transport sse --http-addr :9090
```

### Connexion refusée

Le serveur MCP ne fonctionne pas ou n'est pas accessible.

```bash
# Assurez-vous que le serveur fonctionne
swag2mcp mcp --transport sse --http-addr 127.0.0.1:8080

# Dans un autre terminal, vérifiez le point d'accès health
curl http://127.0.0.1:8080/health

# Si vous utilisez un chemin personnalisé
curl http://127.0.0.1:8080/chemin-personnalise/health
```

### Les outils MCP n'apparaissent pas dans le client LLM

Le client LLM ne voit aucun outil.

```bash
# Vérifiez que les specs sont chargées
swag2mcp ls

# Vérifiez que les specs ne sont pas désactivées
swag2mcp validate

# Vérifiez les journaux du serveur
swag2mcp mcp --logfile /tmp/swag2mcp.log
cat /tmp/swag2mcp.log

# Vérifiez que le chemin de l'espace de travail dans la configuration de votre IDE est correct
# (doit être un chemin absolu)
```

**Causes fréquentes :**
- Mauvais chemin d'espace de travail dans la configuration IDE
- Toutes les specs ont `disable: true`
- Les specs sont filtrées par `--tags`
- Le fichier de configuration n'existe pas au chemin spécifié

### La poignée de main MCP échoue (transport HTTP)

Pour les transports SSE et Streamable HTTP, le protocole MCP nécessite une initialisation avant que les appels d'outils fonctionnent.

```
Étape 1 : POST /mcp → {"method":"initialize", ...}
Étape 2 : POST /mcp → {"method":"notifications/initialized"}
Étape 3 : POST /mcp → {"method":"tools/list", ...}  ← fonctionne maintenant
```

Assurez-vous que votre client LLM effectue la poignée de main avant d'appeler des outils.

### Le health check retourne 404

Le chemin du point d'accès health peut différer du chemin MCP.

```bash
# Point d'accès health par défaut
curl http://127.0.0.1:8080/health

# Si vous avez changé le chemin MCP, health est toujours à /health
# (non affecté par --http-path)
```

### Outil auth non disponible

L'outil MCP `auth` n'apparaît pas.

L'outil `auth` est **désactivé par défaut** (`--disable-llm-auth=true`). C'est intentionnel pour la sécurité en production.

```bash
# Activez l'outil auth
swag2mcp mcp --disable-llm-auth=false
```

## Problèmes d'authentification

### 401 Non autorisé

L'API a rejeté la requête en raison d'identifiants manquants ou invalides.

```bash
# Vérifiez que l'authentification est configurée
swag2mcp info

# Validez la configuration
swag2mcp validate

# Vérifiez que les variables d'environnement sont définies
echo $MON_JETON

# Vérifiez que le jeton n'est pas expiré (les jetons bearer sont statiques)
```

**Causes fréquentes :**
- Jeton manquant ou vide
- Variable d'environnement non définie
- Jeton expiré (les jetons bearer ne s'actualisent pas automatiquement)
- Mauvais type d'authentification configuré

### 403 Interdit

L'API a rejeté la requête en raison de permissions insuffisantes.

- Le jeton peut ne pas avoir les portées requises
- La clé API peut ne pas avoir accès à cette ressource
- Consultez la documentation de l'API pour les permissions requises

### Point d'accès du jeton OAuth2 inaccessible

swag2mcp ne peut pas atteindre l'URL du jeton OAuth2.

```bash
# Vérifiez le token_url dans votre configuration
# Vérifiez que l'URL est correcte et accessible
curl -X POST https://auth.example.com/oauth/token \
  -d "grant_type=client_credentials" \
  -d "client_id=test" \
  -d "client_secret=test"

# Vérifiez la connectivité réseau
# Vérifiez les paramètres de proxy si derrière un proxy d'entreprise
```

### L'authentification Digest échoue

swag2mcp ne peut pas effectuer la poignée de main d'authentification Digest.

- Le serveur doit retourner un en-tête `WWW-Authenticate: Digest ...` avec une réponse 401
- Le défi est mis en cache pendant 5 minutes — si le serveur change son nonce, attendez l'expiration du cache
- Vérifiez que le nom d'utilisateur et le mot de passe sont corrects

### Non-concordance de signature HMAC

L'API a rejeté la requête signée HMAC.

- Vérifiez que `api_key` et `secret_key` sont corrects
- Vérifiez que l'API utilise la signature HMAC-SHA256 de style Binance
- Certains échanges utilisent des méthodes de signature différentes — l'authentification HMAC est spécifiquement pour les API compatibles Binance

### L'authentification par script échoue

Le script d'authentification externe a échoué.

```bash
# Vérifiez que le script existe
ls -la ~/.swag2mcp/auth_scripts/mon-domaine.sh

# Exécutez le script manuellement pour tester
sh ~/.swag2mcp/auth_scripts/mon-domaine.sh

# Vérifiez le format de sortie du script (doit être JSON : {"token": "...", "expires_in": 3600})
# Vérifiez que le script se termine dans les 30 secondes
# Vérifiez que le script a les permissions d'exécution
chmod +x ~/.swag2mcp/auth_scripts/mon-domaine.sh
```

## Problèmes de recherche

### Aucun résultat de recherche

La recherche n'a retourné aucun point d'accès.

```bash
# Vérifiez que les specs sont chargées
swag2mcp ls

# Vérifiez que les specs ne sont pas désactivées
swag2mcp validate

# Essayez une requête plus simple
# Essayez de rechercher par méthode : method:GET
# Essayez de rechercher par étiquette : tag:pets

# L'index est reconstruit à chaque démarrage du serveur MCP
# Si vous venez d'ajouter une spec, redémarrez le serveur
```

### La recherche retourne des résultats non pertinents

La requête est trop large ou ambiguë.

- Utilisez des filtres de champ pour affiner : `method:GET +tag:pets`
- Utilisez des phrases exactes : `"find pet by status"`
- Utilisez le paramètre `limit` pour obtenir des résultats plus ciblés

## Problèmes d'appel API

### invoke retourne une erreur

L'appel API a échoué.

```bash
# Vérifiez le message d'erreur — il inclut le code d'état HTTP
# Erreurs 4xx : vérifiez les paramètres, l'authentification ou les permissions
# Erreurs 5xx : le serveur API a un problème

# Inspectez toujours le point d'accès avant d'invoquer
inspect(endpointId: "...")

# Vérifiez que tous les paramètres requis sont fournis
# Vérifiez les types de paramètres (chaîne, nombre, booléen)
```

### Erreur de limite de débit

Le LLM a appelé le même point d'accès trop rapidement.

Chaque point d'accès a un délai de 10 secondes. Attendez avant d'appeler à nouveau, ou désactivez le limiteur de débit :

```yaml
disable_ratelimiter: true
```

### Réponse trop volumineuse (fileRef retourné)

La réponse a dépassé `max_response_size`.

C'est normal. Utilisez les outils de réponse pour explorer les données :

```
1. response_outline(chemin) → comprendre la structure
2. response_compress(chemin, mode: "first_of_array") → obtenir un échantillon
3. response_slice(chemin, jsonPath: "data.0") → obtenir des données spécifiques
```

Ou augmentez la limite :

```yaml
http_client:
  max_response_size: 4194304  # 4 Mo
```

### Réponses API lentes

L'API prend trop de temps à répondre.

```yaml
http_client:
  timeout: 120s  # Augmentez par rapport à la valeur par défaut de 30s
```

## Problèmes d'espace de travail

### swag2mcp init échoue : « le répertoire n'est pas vide »

Le répertoire cible contient déjà des fichiers.

```bash
# Utilisez --force pour écraser
swag2mcp init --force

# Ou utilisez un répertoire différent
swag2mcp init ./nouvel-espace-travail
```

### swag2mcp update échoue

Un ou plusieurs fichiers de spécification n'ont pas pu être téléchargés.

```bash
# Vérifiez le message d'erreur pour savoir quelle URL a échoué
# Vérifiez que l'URL est accessible
curl -I <url-echouee>

# Vérifiez la connectivité réseau
# Vérifiez les paramètres de proxy
```

### L'exportation ne crée pas de ZIP

L'argument `[output]` doit être un chemin de fichier se terminant par `.zip`, pas un répertoire.

```bash
# Correct
swag2mcp export /chemin/vers/espace-travail /chemin/vers/sauvegarde.zip

# Erroné (aucun ZIP ne sera créé)
swag2mcp export /chemin/vers/espace-travail /un/repertoire
```

### L'importation échoue : « n'est pas une sauvegarde swag2mcp valide »

Le fichier ZIP n'a pas été créé par `swag2mcp export`.

Seules les archives ZIP créées par `swag2mcp export` peuvent être importées. L'archive a une structure interne spécifique (`swag2mcp.yaml`, `specs/`, `auth_scripts/`).

## Problèmes de TUI

### La TUI ne s'affiche pas correctement

Le terminal est trop petit ou ne prend pas en charge les fonctionnalités requises.

- Taille minimale du terminal : 80×24 caractères
- La TUI utilise Bubbletea et fonctionne dans la plupart des terminaux modernes
- Essayez de redimensionner votre fenêtre de terminal
- Essayez un émulateur de terminal différent

### La TUI affiche « aucune spec trouvée »

L'espace de travail n'a pas de specs configurées.

```bash
# Vérifiez les specs
swag2mcp ls

# Ajoutez une spec
swag2mcp add spec
```

## Problèmes du serveur mock

### Le serveur mock ne démarre pas

```bash
# Vérifiez que mock_enabled: true dans la configuration
# Vérifiez que chaque collection a base_mock_url défini
# Vérifiez que les ports ne sont pas utilisés
lsof -i :9090

# Vérifiez les journaux du serveur mock
swag2mcp-mock mockserver
```

### Le serveur mock retourne des réponses vides

Le fichier de spécification peut ne pas avoir de schémas de réponse définis.

- Le serveur mock génère des données à partir des schémas de réponse
- Si aucun schéma n'est trouvé, il retourne `{}`
- Vérifiez que votre spécification OpenAPI a des `responses` avec `schema` défini

## Problèmes réseau

### Échec de la connexion proxy

swag2mcp ne peut pas se connecter via le proxy configuré.

```bash
# Vérifiez le format de l'URL du proxy (doit inclure le schéma : http://, https://, socks5://)
# Vérifiez les identifiants du proxy
# Vérifiez la liste de contournement — la cible peut être dans la liste de contournement
# Testez le proxy avec curl
curl -x http://proxy.entreprise.com:8080 https://api.example.com
```

### Erreurs TLS/SSL

La vérification du certificat a échoué.

- Si vous utilisez un certificat auto-signé pour le serveur MCP, le client doit lui faire confiance
- Pour le serveur mock avec `--tls`, un certificat auto-signé est généré automatiquement
- Pour les appels API, swag2mcp utilise le magasin de certificats du système

## Autres problèmes

### Utilisation élevée du disque

Les répertoires de cache et de réponses peuvent croître avec le temps.

```bash
# Nettoyez tout
swag2mcp clean

# Les anciennes réponses (>48h) sont nettoyées automatiquement au démarrage du serveur MCP
# Les fichiers de cache expirent aléatoirement entre 1 et 48 heures
```

### « commande introuvable » après go install

Le répertoire `go install` n'est pas dans votre PATH.

```bash
# Trouvez où Go installe les binaires
go env GOPATH
# Ajoutez au PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Le LLM n'utilise pas correctement les outils

Le LLM peut avoir besoin de meilleures instructions ou d'une compétence de formatage.

- Utilisez `llm_instruction` dans votre configuration de spec pour décrire ce que fait l'API
- Envisagez d'utiliser la [compétence swag2mcp-format](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md) pour un formatage cohérent des sorties
- La qualité des réponses du LLM dépend du modèle et des instructions qu'il reçoit

### Comment signaler un bogue ?

Ouvrez un problème sur [GitHub](https://github.com/mmadfox/swag2mcp/issues) avec :
- La version de swag2mcp (`swag2mcp --version`)
- Votre système d'exploitation et architecture
- La commande exacte que vous avez exécutée
- Le message d'erreur complet
- Votre fichier de configuration (avec les secrets supprimés)
