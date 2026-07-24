# Compétences

## Personnalisation du format de sortie

Chaque outil MCP de swag2mcp renvoie des données JSON structurées. La manière dont ces données sont **présentées** à l'utilisateur dépend de la compétence de formatage du LLM — et vous pouvez la contrôler complètement.

### La compétence de formatage par défaut

swag2mcp est livré avec une compétence de formatage intégrée qui définit un markdown compact et lisible pour chaque réponse d'outil :

[swag2mcp-format SKILL.md](https://github.com/mmadfox/swag2mcp/blob/main/.agents/skills/swag2mcp-format/SKILL.md)

Cette compétence couvre les 19 outils MCP avec :
- Tableaux compacts pour les listes (spécifications, collections, balises, points de terminaison)
- En-têtes en ligne pour les vues détaillées
- Représentation compacte des schémas pour `inspect`
- Style cohérent pour toutes les réponses

### Pourquoi les compétences sont importantes

Les mêmes données peuvent être présentées de manière radicalement différente selon la compétence :

| Style | Exemple de sortie |
|-------|---------------|
| **Tableaux compacts** (par défaut) | `GET /pet/{petId}` — Trouver un animal par ID |
| **Verbeux** | `Méthode : GET, Chemin : /pet/{petId}, Résumé : Trouver un animal par ID, Obsolète : false` |
| **Minimal** | `GET /pet/{petId}` |
| **Technique** | `GET /pet/{petId} → 200 : Objet Pet, 404 : Non trouvé` |
| **Personnalisé** | Tout format que vous pouvez décrire |

### Créer votre propre compétence

Vous pouvez écrire votre propre compétence de formatage en décrivant le format de sortie exact que vous souhaitez. La compétence est un fichier markdown avec des règles de formatage pour chaque outil. Voici quelques idées :

- **Sortie JSON** — renvoyer le JSON brut pour une consommation machine
- **Style CSV** — données tabulaires pour l'importation dans un tableur
- **Adapté aux diagrammes** — diagrammes Mermaid ou ASCII de la structure API
- **Minimal** — juste la méthode et le chemin, rien d'autre
- **Style documentation** — descriptions complètes, exemples et notes

### La seule limite est le modèle

La qualité de la sortie formatée dépend entièrement de la capacité du LLM à suivre vos règles de formatage. Une compétence bien écrite avec des exemples clairs produit une sortie cohérente et fiable. Une compétence vague produit des résultats incohérents.

Vous pouvez :
- Utiliser la compétence par défaut telle quelle
- La forker et ajuster le formatage à votre goût
- Écrire la vôtre à partir de zéro
- Changer de compétence selon la tâche

### Comment utiliser une compétence

Les compétences sont chargées par le client LLM (OpenCode, Cursor, Claude Desktop, etc.) dans le cadre de son invite système ou de sa configuration d'agent. Consultez la documentation de votre client pour savoir comment attacher un fichier de compétence.

Pour OpenCode, les compétences sont configurées dans `opencode.json` :

```json
{
  "skills": [
    {
      "name": "swag2mcp-format",
      "sourceURL": "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/.agents/skills/swag2mcp-format/SKILL.md"
    }
  ]
}
```
