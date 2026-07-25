# export

## Objectif

Créer une sauvegarde ZIP portable de l'espace de travail. L'archive contient le fichier de configuration, tous les fichiers de spécification et les scripts d'authentification — tout ce qui est nécessaire pour restaurer l'espace de travail sur une autre machine.

## Quand l'utiliser

- Vous voulez sauvegarder votre espace de travail avant d'apporter des modifications
- Vous migrez swag2mcp vers une autre machine
- Vous voulez partager votre configuration API avec un collègue
- Vous préparez un environnement reproductible

## Syntaxe

```bash
swag2mcp export [chemin] [sortie] [drapeaux]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |
| `sortie` | 2 | Non | Chemin complet pour le fichier ZIP de sortie. S'il est omis, par défaut `./swag2mcp-backup-&lt;horodatage&gt;.zip`. |

## Drapeaux

| Drapeau | Raccourci | Type | Défaut | Description |
|---------|-----------|------|--------|-------------|
| `--spec` | `-s` | `stringSlice` | `nil` | Exporter uniquement les specs spécifiées (séparées par des virgules) |

## Comment cela fonctionne

### Exportation par défaut

Crée un ZIP dans le répertoire courant avec un nom horodaté :

```bash
swag2mcp export
# Crée ./swag2mcp-backup-2026-07-22-143022.zip
```

### Chemin de sortie personnalisé

```bash
swag2mcp export /chemin/vers/espace-travail /chemin/vers/sauvegarde.zip
```

### Exporter des specs spécifiques

```bash
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

## Contenu du ZIP

| Entrée | Description |
|--------|-------------|
| `swag2mcp.meta` | Métadonnées sur l'exportation |
| `swag2mcp.yaml` | Fichier de configuration |
| `specs/` | Tous les fichiers de spécification (OpenAPI/Swagger/Postman) |
| `auth_scripts/` | Scripts d'authentification |
| `cache/` | Vide (le cache n'est pas exporté) |
| `responses/` | Vide (les réponses ne sont pas exportées) |

## Restauration

Utilisez `import` pour restaurer à partir d'une sauvegarde :

```bash
swag2mcp import --from-zip /chemin/vers/sauvegarde.zip
```

## Vérification post-commande

Vérifiez toujours que le fichier ZIP a été créé :

```bash
ls -la swag2mcp-backup-*.zip
# ou pour un chemin de sortie personnalisé :
ls -la /chemin/vers/sauvegarde.zip
```

## Nuances

- **La sortie doit être un chemin de fichier :** L'argument `[sortie]` doit être un chemin de fichier complet se terminant par `.zip`. Ne passez **pas** un répertoire — la commande ne créera pas de ZIP si un chemin de répertoire est donné.
- **Nom de fichier par défaut :** `swag2mcp-backup-<AAAA-MM-JJ-HHMMSS>.zip` utilisant l'horodatage UTC.
- **Filtre `--spec` :** Lorsqu'il est défini, seules les specs spécifiées sont incluses. Les autres specs sont exclues de l'archive.
- **Aucune configuration requise :** `export` fonctionne même sans fichier de configuration valide. Il exporte tout ce qui existe dans l'espace de travail.
- **Le cache et les réponses sont exclus :** Ce sont des données transitoires qui seraient obsolètes lors de la restauration. Seuls la configuration, les specs et les scripts d'authentification sont conservés.
