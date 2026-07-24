# Explorateur TUI

## Aperçu

swag2mcp inclut une TUI (Interface Utilisateur Terminal) intégrée pour l'exploration interactive d'API. C'est une application terminal plein écran qui vous permet de rechercher, parcourir, inspecter et invoquer des points d'accès API sans quitter le terminal.

## Lancement

```bash
swag2mcp run
```

Si aucun fichier de configuration n'existe, la TUI démarrera automatiquement l'assistant d'initialisation d'abord.

## Modes

La TUI a trois modes, accessibles avec la touche `Tab` :

### Mode Recherche

Recherche en texte intégral dans tous les points d'accès de toutes les specs. Prend en charge la même syntaxe de requête que l'outil MCP `search`.

- Saisissez une requête pour rechercher les noms, chemins et descriptions des points d'accès
- Filtrez les résultats par méthode, étiquette ou chemin
- Affichez les détails d'un point d'accès avec une seule touche
- Naviguez dans les résultats avec la pagination (10 éléments par page)

### Mode Parcourir

Navigation arborescente dans la hiérarchie des specs :

```
Spec → Collection → Tag → Point d'accès
```

- Descendez dans l'arborescence pour trouver des points d'accès spécifiques
- Affichez les détails d'un point d'accès (paramètres, corps de requête, réponses)
- Invoquez l'API directement depuis la TUI
- Sauvegardez les détails d'un point d'accès sous forme de fichier JSON

### Mode Auth

Affichez les jetons d'authentification et les en-têtes pour n'importe quelle spec. Utile pour le débogage ou la génération de commandes curl.

## Contrôles

| Touche | Action |
|--------|--------|
| `↑` / `↓` | Naviguer vers le haut/bas |
| `Entrée` | Sélectionner ou ouvrir |
| `Échap` | Revenir d'un niveau |
| `Tab` | Basculer entre les modes Recherche, Parcourir et Auth |
| `/` | Activer la saisie de recherche |
| `N` / `P` | Page suivante / précédente |
| `B` | Retour à l'écran précédent |
| `M` | Retour au menu principal |
| `S` | Sauvegarder les détails du point d'accès en JSON |
| `q` / `Ctrl+C` | Quitter |

## États

La TUI passe par ces états au fur et à mesure de votre navigation :

1. **Chargement** — chargement des données depuis l'espace de travail
2. **Recherche** — mode recherche avec saisie de requête
3. **Parcourir** — mode parcourir avec liste des specs
4. **Liste des specs** — liste de toutes les specs
5. **Liste des collections** — collections dans une spec
6. **Liste des étiquettes** — étiquettes dans une collection
7. **Liste des points d'accès** — points d'accès dans une étiquette
8. **Détail du point d'accès** — informations complètes sur le point d'accès
9. **Résultat d'invocation** — résultat de l'appel API
10. **Erreur** — état d'erreur avec message

## Vue détaillée du point d'accès

Lorsque vous sélectionnez un point d'accès, la TUI affiche :

- Méthode HTTP et chemin
- URL de base et URL complète
- Résumé et description
- Tous les paramètres (nom, emplacement, type, requis)
- Schéma du corps de la requête (le cas échéant)
- Codes de réponse et schémas
- Statut de dépréciation

## Prérequis

- **Taille du terminal :** Au moins 80×24 caractères
- **Émulateur de terminal :** Fonctionne dans la plupart des terminaux modernes (iTerm2, Terminal.app, GNOME Terminal, Windows Terminal, etc.)
- **SSH :** Fonctionne via des connexions SSH

## Notes importantes

- **Auto-initialisation** — si aucun fichier de configuration n'existe, la TUI démarre automatiquement l'assistant d'initialisation
- **Pagination** — les listes sont paginées à 10 éléments par page. Utilisez `N` et `P` pour naviguer
- **Sauvegarde des détails du point d'accès** — appuyez sur `S` dans la vue détaillée du point d'accès pour sauvegarder les détails complets sous forme de fichier JSON dans le répertoire courant
- **Mode Auth** — affiche les jetons et en-têtes pour le débogage. En production, l'outil auth peut être désactivé avec `--disable-llm-auth`
