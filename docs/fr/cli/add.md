# add

## Objectif

Ajouter une nouvelle **spec** (service API) ou **collection** (fichier OpenAPI/Swagger/Postman) à une configuration existante. C'est le moyen principal d'enrichir votre espace de travail avec de nouvelles API.

## Quand l'utiliser

- Vous avez une nouvelle API à connecter à votre agent LLM
- Vous avez trouvé une URL de spécification OpenAPI et souhaitez l'ajouter
- Vous voulez ajouter un fichier de spécification supplémentaire (collection) à une spec existante
- Vous préférez écrire du YAML directement plutôt que d'utiliser l'assistant interactif

## Syntaxe

```bash
swag2mcp add spec [chemin] [drapeaux]
swag2mcp add collection [chemin] [drapeaux]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |

## Drapeaux

### `add spec`

| Drapeau | Raccourci | Type | Défaut | Description |
|---------|-----------|------|--------|-------------|
| `--yaml` | `-y` | `string` | `""` | Entrée YAML en ligne ou `-` pour stdin |
| `--example` | `-e` | `bool` | `false` | Afficher un modèle YAML et quitter |

### `add collection`

| Drapeau | Raccourci | Type | Défaut | Description |
|---------|-----------|------|--------|-------------|
| `--yaml` | `-y` | `string` | `""` | Entrée YAML en ligne ou `-` pour stdin |
| `--example` | `-e` | `bool` | `false` | Afficher un modèle YAML et quitter |

## Comment cela fonctionne

### Mode interactif (par défaut)

Lance un assistant TUI qui vous permet de remplir les champs de la spec ou de la collection étape par étape.

```bash
swag2mcp add spec
swag2mcp add collection
```

### Mode YAML en ligne

Passez le YAML directement sous forme de chaîne. **Faites attention aux guillemets du shell** — les caractères spéciaux comme `:`, `#`, `&`, `{` peuvent casser la commande.

```bash
swag2mcp add spec --yaml 'domain: meteo
llm_title: API Open-Meteo
base_url: https://meteo.swagger.io/v2
collections:
  - llm_title: Principal
    location: https://example.com/spec.json'
```

### YAML depuis stdin (recommandé pour les YAML complexes)

Utilisez un pipe depuis un fichier ou un heredoc pour éviter complètement les problèmes de guillemets du shell :

```bash
# Pipe depuis un fichier
cat spec.yaml | swag2mcp add spec --yaml -

# Heredoc
swag2mcp add spec --yaml - <<EOF
domain: mon-api
llm_title: Mon API
llm_instruction: "Utilisez cette API pour X & Y # important"
base_url: https://api.example.com/v1
collections:
  - llm_title: Principal
    location: https://raw.githubusercontent.com/org/repo/main/spec.yaml
EOF
```

### Modèle YAML

Affiche la structure YAML attendue et quitte :

```bash
swag2mcp add spec --example
swag2mcp add collection --example
```

## Format YAML

### Spec

```yaml
domain: meteo
llm_title: API Open-Meteo
llm_instruction: Utilisez cette API pour gérer les animaux.
base_url: https://meteo.swagger.io/v2
tags: [public, demo]
auth:
  type: bearer
  config:
    token: $(JETON)
collections:
  - llm_title: Open-Meteo Swagger
    location: https://example.com/spec.json
```

### Collection

```yaml
spec_domain: meteo
llm_title: Collection Commandes
location: https://example.com/orders.json
```

## Vérification post-commande

```bash
swag2mcp ls [chemin]
# La nouvelle spec ou collection devrait apparaître dans la liste
```

## Nuances

- **Auto-initialisation :** Si aucun fichier de configuration n'existe, `add` exécute automatiquement l'assistant d'initialisation d'abord. Vous n'avez pas besoin d'exécuter `init` séparément.
- **Guillemets du shell :** Le YAML en ligne (`--yaml '...'`) est fragile avec les caractères spéciaux. Préférez `--yaml -` avec un heredoc ou un pipe pour tout ce qui dépasse les valeurs simples.
- **`--example` quitte immédiatement** sans vérifier une configuration existante ni modifier quoi que ce soit.
- **`add spec` vs `add collection` :** Utilisez `add spec` pour un nouveau service API (nouveau domaine). Utilisez `add collection` pour ajouter un autre fichier de spécification à une spec existante.
