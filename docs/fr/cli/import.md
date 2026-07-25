# import

## Objectif

Importer des fichiers de spécification dans l'espace de travail ou restaurer un espace de travail complet à partir d'une sauvegarde ZIP. Trois modes couvrent différents scénarios : ajout d'une seule spec, importation en masse depuis une configuration existante ou restauration d'un espace de travail complet.

## Quand l'utiliser

- Vous avez une URL ou un fichier de spécification et voulez l'ajouter à l'espace de travail
- Vous voulez télécharger tous les fichiers de spécification référencés dans la configuration
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
| `nom` | 3 | Variable | Nom de domaine pour la nouvelle spec |

## Drapeaux

| Drapeau | Raccourci | Type | Défaut | Description |
|---------|-----------|------|--------|-------------|
| `--spec` | `-s` | `stringSlice` | `nil` | Importer les collections des specs spécifiées (séparées par des virgules) |
| `--from-zip` | | `string` | `""` | Restaurer l'espace de travail à partir d'un ZIP de sauvegarde swag2mcp |

## Comment cela fonctionne

### Mode 1 — Importation unique depuis une URL ou un fichier

Téléchargez un fichier de spécification et ajoutez-le à l'espace de travail avec un nom de domaine :

```bash
swag2mcp import https://example.com/spec.yaml maspec
swag2mcp import /chemin/vers/espace-travail https://example.com/spec.yaml maspec
swag2mcp import ./spec-locale.yaml maspec
```

Le fichier de spécification est sauvegardé dans `specs/` et la configuration est mise à jour avec la nouvelle entrée de spec.

### Mode 2 — Importation en masse depuis une configuration existante

Téléchargez toutes les collections pour les domaines spécifiés à partir de leurs URL configurées :

```bash
swag2mcp import --spec meteo
swag2mcp import /chemin/vers/espace-travail --spec meteo,store
```

Le fichier de spécification de chaque collection est téléchargé et sauvegardé dans `specs/`. La configuration est mise à jour pour pointer vers les copies locales.

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
- **`--force` :** Disponible pour la restauration ZIP afin d'écraser un espace de travail existant.
- **Client HTTP :** Les paramètres globaux du client HTTP de la configuration sont appliqués pendant l'importation (délai d'attente, proxy, en-têtes, etc.).
