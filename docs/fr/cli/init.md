# init

## Objectif

La commande `init` crée un **espace de travail** — un répertoire avec un fichier de configuration `swag2mcp.yaml` et des sous-répertoires pour le cache, les specs, les réponses et les scripts d'authentification. C'est la première commande à exécuter lors de la configuration de swag2mcp.

## Quand l'utiliser

- Vous configurez swag2mcp pour la première fois
- Vous voulez créer un nouvel espace de travail dans un répertoire spécifique
- Vous devez réinitialiser un espace de travail corrompu ou manquant

## Syntaxe

```bash
swag2mcp init [chemin] [drapeaux]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, par défaut `~/.swag2mcp`. |

## Drapeaux

| Drapeau | Raccourci | Type | Défaut | Description |
|---------|-----------|------|--------|-------------|
| `--interactive` | `-i` | `bool` | `false` | Exécuter l'assistant TUI interactif |
| `--force` | `-f` | `bool` | `false` | Écraser la configuration existante dans un répertoire non vide |

## Comment cela fonctionne

### Mode non interactif (par défaut)

Crée un `swag2mcp.yaml` minimal sans specs. Vous modifiez le fichier manuellement ensuite.

```bash
swag2mcp init
# Crée ~/.swag2mcp/swag2mcp.yaml

swag2mcp init ./mon-projet
# Crée ./mon-projet/swag2mcp.yaml

swag2mcp init /chemin/absolu
# Crée /chemin/absolu/swag2mcp.yaml
```

### Mode interactif (`-i`)

Lance un assistant TUI en 18 étapes qui vous guide à travers :

1. Le choix du répertoire de l'espace de travail
2. L'ajout de specs avec domaine, titre, URL de base
3. La configuration des collections avec les URL d'emplacement
4. La configuration de l'authentification (les 9 méthodes)
5. La configuration des paramètres du client HTTP (délai d'attente, proxy, en-têtes, etc.)

```bash
swag2mcp init -i
```

### Mode force (`--force`)

Par défaut, `init` refuse de s'exécuter dans un répertoire non vide. Utilisez `--force` pour écraser :

```bash
swag2mcp init -f
swag2mcp init ./repertoire-existant -f
```

## Ce qui est créé

```
~/.swag2mcp/
├── swag2mcp.yaml       # Fichier de configuration
├── cache/               # Fichiers de spécification distants téléchargés
├── specs/               # Fichiers de spécification locaux
├── responses/           # Réponses d'invocation API sauvegardées
└── auth_scripts/        # Scripts d'authentification (pour le type ScriptAuth)
```

## Vérification post-commande

```bash
ls ~/.swag2mcp/swag2mcp.yaml
# Si le fichier existe, init a réussi
```

## Nuances

- **Résolution de chemin :** `[chemin]` est un **répertoire d'espace de travail**, pas un chemin de fichier. La CLI ajoute `swag2mcp.yaml` automatiquement. Ordre de résolution : `[chemin]` explicite → répertoire courant (`./`) → `~/.swag2mcp/`.
- **Vérification de répertoire non vide :** Sans `--force`, `init` retourne une erreur si le répertoire cible existe et n'est pas vide. Cela empêche les écrasements accidentels.
- **Stubs de scripts d'authentification :** Si une spec utilise `ScriptAuth`, `init` crée des fichiers de script stub (`.sh` sur Unix, `.bat` sur Windows) dans `auth_scripts/`.
- **Sortie :** En cas de succès, affiche le chemin de la configuration et un conseil : « Prochaine étape : modifiez swag2mcp.yaml ou exécutez 'swag2mcp ls' pour lister les specs configurées ».
