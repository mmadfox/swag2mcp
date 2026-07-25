# swag2mcp

<div style="background: #dc2626; color: white; padding: 20px 24px; border-radius: 12px; text-align: center; font-size: 1.4em; font-weight: 700; margin: 24px 0;">
  🚧 TRAVAUX EN COURS — version bientôt disponible !
</div>

Fait le pont entre les spécifications d'API OpenAPI/Swagger/Postman et les agents LLM via le Model Context Protocol (MCP).

<a href="https://www.youtube.com/watch?v=1Da4UmE2f9U" target="_blank">
  <img src="https://raw.githubusercontent.com/mmadfox/swag2mcp/main/docs/cover.png" alt="Aperçu">
</a>

## Votre API parle LLM

Une ligne de configuration transforme n'importe quel fichier OpenAPI/Swagger/Postman en serveur MCP. Les agents LLM découvrent, inspectent et invoquent vos API — zéro code d'intégration.

<img src="/architecture.svg" width="700" alt="Architecture swag2mcp">

## Arrêtez d'écrire des adaptateurs

Chaque fois que vous connectez une nouvelle API à un LLM, vous écrivez le même code passe-partout : analyse de spécification, authentification, gestion des erreurs, limitation de débit. swag2mcp le fait pour vous — 19 outils MCP prêts à l'emploi.

## Qui a besoin de cela

| Rôle | Pourquoi |
|------|---------|
| **Développeur d'agent IA** | Connectez n'importe quelle API en 2 minutes, pas en 2 jours |
| **Ingénieur MCP** | Pas de code de gestion — pointez simplement vers une spécification |
| **Architecte** | Couche d'intégration API unique pour tous les LLM de votre entreprise |
| **Analyste de données** | Accédez aux API en langage naturel, sans codage |
| **DevOps / SRE** | Surveillance et automatisation via LLM sans services supplémentaires |
| **Intégrateur** | 9 méthodes d'authentification prêtes à l'emploi — Basic à OAuth2 en passant par HMAC |
| **Ingénieur QA** | Serveur mock pour des tests isolés sans API réelles |
| **Chef de produit** | Prototypes rapides de fonctionnalités IA sans travail backend |
| **et bien d'autres** | |

---

## Licence

Sous licence **GNU Affero General Public License v3.0** (AGPL v3).

Voir [LICENSE](https://github.com/mmadfox/swag2mcp/blob/main/LICENSE) pour le texte complet de la licence.

```
SPDX-License-Identifier: AGPL-3.0-only
```
