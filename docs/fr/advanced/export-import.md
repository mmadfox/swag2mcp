# Exportation et Importation

## Aperçu

swag2mcp prend en charge le transfert complet d'espace de travail via des archives ZIP. Vous pouvez exporter la totalité de votre espace de travail (configuration, fichiers de spécification, scripts d'authentification) vers un fichier ZIP et le restaurer sur une autre machine.

## Exportation

Crée une sauvegarde ZIP portable de votre espace de travail.

```bash
# Exportation vers le fichier par défaut (swag2mcp-backup-&lt;horodatage&gt;.zip)
swag2mcp export

# Exportation avec un chemin personnalisé
swag2mcp export --output ~/sauvegardes/swag2mcp-sauvegarde.zip

# Exportation de specs spécifiques uniquement
swag2mcp export --spec meteo
swag2mcp export --spec meteo,store
```

### Ce qui est inclus dans l'exportation

| Élément | Description |
|---------|-------------|
| `swag2mcp.yaml` | Fichier de configuration |
| `specs/` | Tous les fichiers de spécification (OpenAPI/Swagger/Postman) |
| `auth_scripts/` | Scripts d'authentification |
| `swag2mcp.meta` | Métadonnées (informations de version pour la compatibilité) |

Le cache et les réponses ne sont **pas** exportés — ils sont transitoires et seraient obsolètes lors de la restauration.

### Nom de fichier par défaut

Si vous ne spécifiez pas de chemin de sortie, le fichier est sauvegardé sous `swag2mcp-backup-<AAAA-MM-JJ-HHMMSS>.zip` dans le répertoire courant (horodatage UTC).

## Importation

Restaurez un espace de travail à partir d'une sauvegarde ZIP ou importez des fichiers de spécification.

### Restauration depuis un ZIP

```bash
# Restauration complète de l'espace de travail
swag2mcp import --from-zip /chemin/vers/sauvegarde.zip

# Restauration avec écrasement
swag2mcp import --from-zip /chemin/vers/sauvegarde.zip -f
```

Le ZIP doit être créé par `swag2mcp export` — les fichiers ZIP arbitraires ne fonctionneront pas.

### Importation d'un seul fichier de spécification

Téléchargez un fichier de spécification et ajoutez-le à l'espace de travail :

```bash
swag2mcp import https://example.com/spec.yaml maspec
swag2mcp import /chemin/vers/espace-travail https://example.com/spec.yaml maspec
```

### Importation en masse depuis une configuration existante

Téléchargez tous les fichiers de spécification de collection pour les specs spécifiées (domaines) :

```bash
swag2mcp import --spec meteo
swag2mcp import /chemin/vers/espace-travail --spec meteo,store
```

Cela télécharge le fichier de spécification de chaque collection, le sauvegarde dans `specs/` et met à jour la configuration pour pointer vers la copie locale.

## Cas d'utilisation

### Sauvegarde

```bash
swag2mcp export --output swag2mcp-$(date +%Y-%m-%d).zip
```

### Transfert vers une autre machine

```bash
# Sur l'ancienne machine
swag2mcp export --output swag2mcp.zip

# Copiez le ZIP sur la nouvelle machine, puis :
swag2mcp import --from-zip swag2mcp.zip
```

### Partage de configuration

```bash
swag2mcp init
swag2mcp export --output template.zip
# Partagez template.zip avec un collègue
```

## Vérification post-exportation

Vérifiez toujours que le fichier ZIP a été créé :

```bash
ls -la swag2mcp-backup-*.zip
```

## Notes importantes

- **La sortie doit être un chemin de fichier se terminant par `.zip`** — ne passez pas un répertoire
- **Le cache et les réponses sont exclus** — seuls la configuration, les specs et les scripts d'authentification sont conservés
- **Le ZIP est autonome** — il peut être restauré sur n'importe quelle machine avec swag2mcp installé
- **Filtre de spec** — utilisez `--spec` pour exporter ou importer uniquement des specs spécifiques
