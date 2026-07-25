# Commandes CLI

## Aperçu

La CLI `swag2mcp` est le point d'entrée unique pour toutes les opérations — de l'initialisation d'un espace de travail et la gestion des spécifications API au démarrage d'un serveur MCP pour l'intégration LLM. Elle fournit **13 commandes** qui couvrent l'ensemble du cycle de vie du travail avec les specs OpenAPI/Swagger/Postman.

### Ce que la CLI résout

- **Cycle de vie de l'espace de travail** — créer (`init`), inspecter (`info`, `ls`), nettoyer (`clean`), mettre à jour (`update`) et supprimer (`delete`) les espaces de travail et leur contenu
- **Gestion des specs et collections** — ajouter (`add`), lister (`ls`) et supprimer (`delete`) les spécifications API et leurs collections
- **Modes d'exécution** — démarrer le serveur MCP pour l'accès aux outils LLM (`mcp`) ou lancer l'explorateur TUI interactif (`run`)
- **Diagnostics** — valider la configuration (`validate`), afficher la version (`version`), afficher les informations d'exécution (`info`)
- **Sauvegarde et restauration** — transfert complet d'espace de travail via ZIP (`export`, `import`)

### Nuances clés

- **Résolution de chemin** — les commandes qui acceptent `[chemin]` attendent un **répertoire d'espace de travail** (pas un chemin de fichier). Ordre de résolution : `[chemin]` explicite → répertoire courant (`./`) → `~/.swag2mcp/`. La CLI ajoute `swag2mcp.yaml` automatiquement. Passez toujours un chemin explicite lorsque vous l'exécutez comme service ou dans une configuration IDE pour éviter de charger le mauvais espace de travail.
- **Spec vs Collection** — une **spec** représente un service API logique (par exemple « API Open-Meteo »), tandis qu'une **collection** est un fichier OpenAPI/Swagger/Postman. Une spec peut avoir plusieurs collections.
- **`--version`** est pris en charge à la fois comme drapeau (`swag2mcp --version`) et comme sous-commande (`swag2mcp version`).
- **`add spec` / `add collection`** acceptent une entrée YAML via `--yaml` (chaîne en ligne ou `-` pour stdin). L'utilisation d'un pipe depuis un fichier ou d'un heredoc évite les problèmes de guillemets du shell avec les caractères spéciaux.
- **`delete`** nécessite un TTY (terminal interactif). Il n'y a pas de drapeau `--force` ou `--yes` — il demande toujours une sélection et une confirmation.
- **`mcp`** est la commande principale pour l'intégration LLM. Elle prend en charge trois transports : `stdio` (défaut), `sse` et `streamable-http`. Le drapeau `--disable-llm-auth` (défaut : `true`) supprime l'outil `auth` de la liste des outils MCP, empêchant le LLM de voir ou de demander des jetons. L'authentification fonctionne toujours — les jetons sont obtenus via le mécanisme de configuration standard, pas via le LLM. Ce mode est recommandé pour la **production** (le LLM n'a jamais accès aux identifiants). Pour le **débogage** ou lors de l'utilisation de jetons de courte durée, définissez `--disable-llm-auth=false` pour permettre au LLM de demander des jetons frais via l'outil `auth`.
- **`validate`** vérifie la syntaxe YAML, la structure de la configuration, l'existence des fichiers de spécification, l'accessibilité des URL, le format de la spec (OpenAPI/Swagger/Postman), les paramètres d'authentification et l'exactitude du client HTTP. Il ne **teste pas** les points d'accès d'authentification ni la disponibilité des points d'accès API.
- **`export` / `import`** fournissent un transfert complet d'espace de travail — fichier de configuration, fichiers de spécification, cache et scripts d'authentification sont tous inclus dans l'archive ZIP.
- **`clean`** supprime les répertoires `cache/` et `responses/` mais préserve `specs/` et `auth_scripts/`. Les anciennes réponses (>48h) sont également nettoyées automatiquement au démarrage de `mcp`.

## Commandes

| Commande | Description |
|----------|-------------|
| [`init`](/cli/init) | Initialiser un répertoire d'espace de travail avec la configuration par défaut |
| [`add`](/cli/add) | Ajouter une spec ou une collection à la configuration |
| [`delete`](/cli/delete) | Supprimer une spec ou une collection de manière interactive |
| [`ls`](/cli/ls) | Lister toutes les specs et leurs collections |
| [`run`](/cli/run) | Lancer l'explorateur TUI interactif d'API |
| [`validate`](/cli/validate) | Valider la configuration et les fichiers de spécification |
| [`clean`](/cli/clean) | Vider les specs en cache et les réponses d'invocation |
| [`update`](/cli/update) | Revalider, re-mettre en cache et réindexer toutes les specs |
| [`mcp`](/cli/mcp) | Démarrer le serveur MCP pour l'accès aux outils LLM |
| [`version`](/cli/version) | Afficher la version de swag2mcp |
| [`info`](/cli/info) | Afficher les informations détaillées de configuration et d'exécution |
| [`import`](/cli/import) | Importer des fichiers de spécification ou restaurer l'espace de travail depuis ZIP |
| [`export`](/cli/export) | Exporter l'espace de travail comme sauvegarde ZIP portable |

## Drapeaux globaux

| Drapeau | Description |
|---------|-------------|
| `--version` | Afficher la version (identique à la sous-commande `version`) |
| `--help` | Afficher l'aide pour n'importe quelle commande |
