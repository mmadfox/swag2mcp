# version

## Objectif

Afficher la version de swag2mcp. Utile pour vérifier la version installée, signaler des bogues ou vérifier la compatibilité.

## Quand l'utiliser

- Vous voulez vérifier quelle version de swag2mcp est installée
- Vous signalez un bogue et devez inclure la version
- Vous voulez vérifier une installation réussie

## Syntaxe

```bash
swag2mcp version
swag2mcp --version
```

## Arguments

Aucun.

## Drapeaux

Aucun.

## Comment cela fonctionne

```bash
swag2mcp version
# swag2mcp v1.2.0

swag2mcp --version
# swag2mcp v1.2.0
```

## Format de sortie

```
swag2mcp <version>
```

La version est définie au moment de la construction via `ldflags`. Si elle n'est pas définie, elle prend par défaut `"dev"`.

## Nuances

- **Deux formes :** `swag2mcp version` (sous-commande) et `swag2mcp --version` (drapeau global) produisent la même sortie.
- **Aucune configuration requise :** Cette commande fonctionne sans espace de travail ni fichier de configuration.
