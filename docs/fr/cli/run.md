# run

## Objectif

Lancer l'explorateur d'API **TUI (Interface Utilisateur Terminal)** interactif. C'est une application plein écran pour rechercher, parcourir, inspecter et invoquer des points d'accès API sans quitter le terminal.

## Quand l'utiliser

- Vous voulez explorer vos API de manière interactive
- Vous devez rechercher un point d'accès spécifique dans toutes les specs
- Vous voulez parcourir la hiérarchie spec → collection → tag → point d'accès
- Vous voulez tester un appel API avant de configurer le serveur MCP

## Syntaxe

```bash
swag2mcp run [chemin]
```

## Arguments

| Argument | Position | Requis | Description |
|----------|----------|--------|-------------|
| `chemin` | 1 | Non | Répertoire de l'espace de travail. S'il est omis, résolution via les règles de résolution de chemin. |

## Drapeaux

Aucun.

## Modes

### Mode Recherche

Recherche en texte intégral dans tous les points d'accès de toutes les specs. Prend en charge le filtrage par méthode HTTP, étiquette et chemin.

- Saisissez une requête pour rechercher les noms, chemins et descriptions des points d'accès
- Filtrez les résultats par méthode (GET, POST, PUT, DELETE, etc.)
- Affichez les détails d'un point d'accès avec une seule touche

### Mode Parcourir

Navigation arborescente dans la hiérarchie des specs :

```
Spec → Collection → Tag → Point d'accès
```

- Descendez dans l'arborescence pour trouver des points d'accès spécifiques
- Affichez les détails d'un point d'accès (paramètres, corps de requête, réponses)
- Invoquez l'API directement depuis la TUI

## Navigation

| Touche | Action |
|--------|--------|
| `↑` / `↓` | Naviguer vers le haut/bas |
| `Entrée` | Sélectionner ou ouvrir |
| `Échap` | Revenir en arrière |
| `Tab` | Basculer entre les modes Recherche et Parcourir |
| `/` | Activer la saisie de recherche |
| `q` | Quitter |

## Vérification post-commande

La TUI charge toutes les specs de l'espace de travail. Si une spec ne parvient pas à se charger, un message d'erreur est affiché dans l'interface.

## Nuances

- **Auto-initialisation :** Si aucun fichier de configuration n'existe, `run` exécute automatiquement l'assistant d'initialisation d'abord.
- **Pas de drapeaux :** La commande `run` n'a pas de drapeaux — toute la configuration provient de l'espace de travail.
- **Taille du terminal :** La TUI nécessite un terminal d'au moins 80×24 caractères. Elle peut ne pas s'afficher correctement dans les très petits terminaux.
- **Dépendances :** La TUI utilise Bubbletea. Elle fonctionne via SSH et dans la plupart des émulateurs de terminal.
