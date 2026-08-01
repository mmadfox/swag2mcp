# swag2mcp

Fait le pont entre les spécifications d'API OpenAPI/Swagger/Postman et les agents LLM via le Model Context Protocol (MCP).

Toutes les API ne parlent pas MCP — les points de terminaison privés, les services internes, les systèmes legacy et les API tierces le font rarement. swag2mcp enveloppe toute API REST dans une interface MCP, donnant aux agents LLM un accès instantané à l'ensemble de votre surface API sans modifier une seule ligne de code serveur. Grâce à des appels API en direct, le LLM acquiert des connaissances concrètes pour prendre des décisions éclairées, automatiser les flux de travail et agir sur vos données — pas seulement deviner.

<a href="https://www.youtube.com/watch?v=1Da4UmE2f9U" target="_blank">
  <img src="https://raw.githubusercontent.com/mmadfox/swag2mcp/main/docs/cover.jpg" alt="Aperçu">
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

Sous licence **Apache License, Version 2.0**.

Voir [LICENSE](https://github.com/mmadfox/swag2mcp/blob/main/LICENSE) pour le texte complet de la licence.
