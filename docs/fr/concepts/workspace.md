# Espace de travail

L'espace de travail est le répertoire où swag2mcp stocke toutes ses données — configuration, specs en cache, fichiers de spécification locaux, réponses sauvegardées et scripts d'authentification.

## Structure

```
~/.swag2mcp/                          # Racine de l'espace de travail (par défaut)
├── swag2mcp.yaml                     # Fichier de configuration
├── cache/                            # Fichiers de spécification distants en cache
│   ├── a1b2c3d4e5f6...spec          # Contenu de la spec en cache
│   └── a1b2c3d4e5f6...meta          # Métadonnées du cache (JSON)
├── specs/                            # Fichiers de spécification locaux
│   └── mon-api.yaml
├── responses/                        # Réponses API sauvegardées (grandes réponses)
│   ├── meteo-get-forecast-abc123.json
│   └── response-fragment-def456.json
└── auth_scripts/                     # Scripts d'authentification
    ├── meteo.sh                      # Script shell Unix
    └── meteo.bat                     # Script batch Windows
```

## Chemin par défaut

- **Linux/macOS** : `~/.swag2mcp/`
- **Windows** : `%USERPROFILE%\.swag2mcp\`

## Chemin personnalisé

```bash
swag2mcp mcp /chemin/vers/espace-travail
swag2mcp mcp ./mon-espace-travail
```

## Répertoires

### cache/

Stocke les fichiers de spécification distants téléchargés. Chaque fichier est mis en cache avec un hachage SHA-256 de son URL comme nom de fichier :

- `{hachage}.spec` — le contenu du fichier de spécification en cache
- `{hachage}.meta` — métadonnées JSON (URL source, heure de mise en cache, TTL)

Chaque fichier en cache a un TTL aléatoire entre 1 heure et 48 heures. Le cache est automatiquement vérifié à chaque démarrage — si une entrée valide (non expirée) existe, elle est réutilisée sans téléchargement.

**Commandes :**
- `swag2mcp update` — vide le cache et retélécharge toutes les specs
- `swag2mcp clean` — vide le cache et les réponses

### specs/

Stocke les fichiers de spécification locaux vers lesquels les collections pointent via `location: specs/{nom}`. Les fichiers ici sont utilisés directement sans mise en cache.

Ce répertoire est rempli par :
- `swag2mcp import &lt;source&gt; &lt;nom&gt;` — télécharge une spec distante et la sauvegarde ici
- `swag2mcp export` — copie les specs ici dans le ZIP d'exportation
- Placement manuel — vous pouvez copier vous-même des fichiers de spécification ici

### responses/

Stocke les réponses API qui dépassent la limite `max_response_size` (par défaut 1 Mo). Lorsque le LLM invoque un point d'accès et que la réponse est trop volumineuse, swag2mcp la sauvegarde ici et retourne une référence de fichier à la place.

Convention de nommage : `{domaine}-{méthode}-{chemin_avec_tirets_bas}-{6car_hex}.json`

Les anciennes réponses sont nettoyées automatiquement après 48 heures au démarrage du serveur MCP.

### auth_scripts/

Stocke les scripts d'authentification pour le type d'auth `script`. Chaque script est nommé d'après le domaine de la spec.

#### Convention de nommage

| Plateforme | Nom de fichier | Exemple |
|------------|----------------|---------|
| Unix (Linux, macOS) | `{domaine}.sh` | `meteo.sh` |
| Windows | `{domaine}.bat` | `meteo.bat` |

Le domaine ne doit pas contenir de caractères `/` ou `\`.

#### Comment les scripts fonctionnent

1. swag2mcp exécute le script avec un délai d'attente de 30 secondes
2. Le script doit produire du JSON valide sur stdout
3. swag2mcp analyse le JSON et utilise le jeton pour les requêtes API

#### Format de sortie attendu

```json
{
  "token": "votre-jeton-ici",
  "expires_in": 3600
}
```

| Champ | Type | Requis | Description |
|-------|------|--------|-------------|
| `token` | string | ✅ | Le jeton d'authentification |
| `access_token` | string | ❌ | Alternative à `token` (vérifié en premier) |
| `token_type` | string | ❌ | Type de jeton (par ex., « Bearer ») |
| `expires_in` | number | ❌ | Durée de vie du jeton en secondes (défaut : 3600) |

#### Exécution

| Plateforme | Commande |
|------------|---------|
| Unix | `sh {domaine}.sh` |
| Windows | `cmd /c {domaine}.bat` |

#### Mise en cache du jeton

Le jeton est mis en cache en mémoire jusqu'à son expiration. À chaque appel API, swag2mcp vérifie d'abord le cache — le script n'est exécuté que lorsque le jeton en cache a expiré.

#### Création de stub

Lorsque vous configurez `auth: { type: script, config: { domain: "monapi" } }`, swag2mcp crée un script stub automatiquement :

**Unix (`auth_scripts/monapi.sh`) :**
```bash
#!/bin/sh
echo '{"token": "votre-jeton-ici", "expires_in": 3600}'
```

**Windows (`auth_scripts/monapi.bat`) :**
```bat
@echo off
echo {"token": "votre-jeton-ici", "expires_in": 3600}
```

Remplacez le jeton factice par votre logique d'authentification réelle.

#### Nettoyage des orphelins

Lorsque vous supprimez une spec, son script d'authentification devient orphelin. swag2mcp supprime automatiquement les scripts orphelins lors de :
- `swag2mcp update`
- `swag2mcp clean`

## Commandes

### update

```bash
swag2mcp update [chemin]
```

Valide la configuration, vide le cache et les réponses, puis retélécharge tous les fichiers de spécification. Assure également que les scripts d'authentification existent et supprime les scripts orphelins.

Utilisez cette commande après :
- Ajout ou suppression de collections
- Changement d'emplacements de collections
- Modification de fichiers de spécification qui nécessitent un re-cache

### clean

```bash
swag2mcp clean [chemin]
```

Supprime tout le contenu de `cache/` et `responses/`, plus les scripts d'authentification orphelins. Ne remet PAS en cache les specs — utilisez `update` pour cela.

### validate

```bash
swag2mcp validate [chemin]
```

Valide la configuration, y compris tous les emplacements de collection. Voir [CLI : validate](../cli/validate.md).

## Exportation et importation

```bash
# Exporter l'espace de travail vers ZIP (nom par défaut : swag2mcp-backup-{date}.zip)
swag2mcp export

# Exporter vers un chemin spécifique
swag2mcp export /chemin/vers/espace-travail /chemin/vers/sauvegarde.zip

# Exporter uniquement des specs spécifiques
swag2mcp export --spec meteo

# Restaurer depuis une sauvegarde
swag2mcp import --from-zip /chemin/vers/sauvegarde.zip
swag2mcp import /chemin/vers/espace-travail /chemin/vers/sauvegarde.zip
```

L'exportation inclut : `swag2mcp.yaml`, `specs/`, `auth_scripts/`. Le cache et les réponses sont exclus (ce sont des données locales).

## .gitignore

Si votre espace de travail se trouve dans un dépôt Git, ajoutez ces entrées à `.gitignore` :

```gitignore
# swag2mcp — données locales uniquement
.swag2mcp/cache/
.swag2mcp/responses/
```

Les répertoires `cache/` et `responses/` contiennent des données locales spécifiques à la machine qui ne doivent pas être commitées. Tout le reste (`swag2mcp.yaml`, `specs/`, `auth_scripts/`) doit être dans le dépôt afin que la configuration soit partagée dans toute l'équipe.
