# import

## Objectif

Importer des fichiers de spécification dans le répertoire `specs/` de l'espace de travail pour une utilisation locale, ou restaurer un espace de travail complet à partir d'une sauvegarde ZIP. Trois modes couvrent différents scénarios : ajout d'une seule spec, importation en masse depuis une configuration existante ou restauration d'un espace de travail complet.

## Quand l'utiliser

- Vous avez une URL ou un fichier de spécification et voulez le sauvegarder localement dans l'espace de travail
- Vous voulez télécharger tous les fichiers de spécification de collection depuis la configuration et rendre l'espace de travail autonome
- Vous devez restaurer un espace de travail à partir d'une sauvegarde ZIP créée par `export`
- Vous migrez swag2mcp vers une autre machine

## Syntaxe

```bash
swag2mcp import [chemin] [source] [nom] [drapeaux]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |
| `source` | 2 | Variable | URL ou chemin local vers un fichier de spécification, ou chemin vers une archive ZIP |
| `nom` | 3 | Variable | Nom de fichier pour l'enregistrement (ex. `example-api.yaml`). Dérivé de l'URL si omis. |

## Drapeaux

| Drapeau | Raccourci | Type | Défaut | Description |
|---------|-----------|------|--------|-------------|
| `--spec` | `-s` | `string` | `""` | Télécharger les fichiers de spécification de collection depuis la configuration. Sans valeur pour toutes les specs, ou spécifier des domaines comme `--spec meteo,github` |
| `--force` | `-f` | `bool` | `false` | Écraser les fichiers de spécification existants sans erreur |
| `--from-zip` | | `string` | `""` | Restaurer l'espace de travail à partir d'un ZIP de sauvegarde swag2mcp |

## Comment cela fonctionne

### Mode 1 — Importation unique depuis une URL ou un fichier

Téléchargez un fichier de spécification et sauvegardez-le dans `specs/` :

```bash
swag2mcp import https://example.com/spec.yaml example-api.yaml
swag2mcp import /chemin/vers/espace-travail https://example.com/spec.yaml example-api.yaml
swag2mcp import ./spec-locale.yaml example-api.yaml
```

Si `nom` est omis, il est dérivé du nom du fichier dans l'URL :
```bash
swag2mcp import https://example.com/specs/petstore.yaml
# → sauvegardé comme petstore.yaml
```

Écraser un fichier existant avec `--force` :
```bash
swag2mcp import --force https://example.com/spec.yaml example-api.yaml
```

Après l'importation, la sortie montre le chemin de l'espace de travail, le fichier sauvegardé et un modèle YAML à ajouter à `swag2mcp.yaml` :

```
✅ Imported to /chemin/vers/espace-travail
   specs/example-api.yaml

   Add to swag2mcp.yaml:
     specs:
       - domain: <your-domain>
         collections:
           - location: specs/example-api.yaml
```

### Mode 2 — Importation en masse depuis une configuration existante (`--spec`)

Téléchargez tous les fichiers de spécification de collection pour les domaines spécifiés depuis leurs URL `location` configurées, sauvegardez-les dans `specs/` et mettez à jour la configuration pour pointer vers les copies locales :

```bash
swag2mcp import --spec                # toutes les specs
swag2mcp import --spec meteo           # spec spécifique
swag2mcp import --spec meteo,github    # plusieurs specs
swag2mcp import /chemin/vers/espace-travail --spec meteo
```

Si un domaine spécifié n'existe pas dans la configuration, la commande renvoie une erreur :
```
Error: import_no_match
  Spec "nonexistent" not found in config.
```

Cela rend l'espace de travail autonome — aucune URL de spécification distante n'est nécessaire après l'importation.

### Mode 3 — Restauration depuis une sauvegarde ZIP

Restaurez un espace de travail complet à partir d'une archive ZIP créée par `swag2mcp export` :

```bash
swag2mcp import --from-zip /chemin/vers/sauvegarde.zip
swag2mcp import /chemin/vers/espace-travail /chemin/vers/sauvegarde.zip
```

> **Le ZIP doit être créé par `swag2mcp export`.** Les fichiers ZIP arbitraires ne fonctionneront pas — l'archive a une structure interne spécifique (`swag2mcp.yaml`, `specs/`, `auth_scripts/`).

## Vérification post-commande

```bash
# Importation unique ou en masse
swag2mcp ls [chemin]
# La nouvelle spec devrait apparaître dans la liste

# Restauration ZIP
swag2mcp ls [chemin]
# Toutes les specs de la sauvegarde devraient apparaître
```

## Nuances

- **Le mode en masse nécessite une configuration :** Lors de l'utilisation de `--spec`, le fichier de configuration doit exister. Exécutez `init` d'abord si nécessaire.
- **L'importation unique crée l'espace de travail :** Si l'espace de travail n'existe pas, il est créé automatiquement.
- **Détection ZIP :** Un argument positionnel se terminant par `.zip` est traité comme une source ZIP. Le drapeau `--from-zip` a priorité sur la détection positionnelle.
- **Client HTTP :** Les paramètres globaux du client HTTP de la configuration sont appliqués pendant l'importation (délai d'attente, proxy, en-têtes, etc.).
