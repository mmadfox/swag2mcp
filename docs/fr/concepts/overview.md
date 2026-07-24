# Concepts

## Architecture

swag2mcp agit comme un pont entre les spécifications API et les agents LLM :

<img src="/architecture.svg" width="800" alt="Architecture swag2mcp">

## Concepts fondamentaux

**Spec** — un conteneur logique représentant un domaine ou service API (par exemple, YouTube, Binance, Open-Meteo). Chaque spec a un `domain` unique, une `base_url`, une `auth` optionnelle, et contient une ou plusieurs collections. Vous pouvez également définir `llm_instruction` — un indice court injecté dans l'invite système swag2mcp qui indique au LLM à quoi sert cette spec et quand l'utiliser. En savoir plus : [Specs](./specs).

**Collection** — un fichier OpenAPI/Swagger/Postman unique décrivant une API spécifique. Elle pointe vers un `location` (URL ou chemin de fichier local). Une spec peut avoir plusieurs collections — par exemple, la spec « meteo » pourrait avoir les collections « Prévisions », « Qualité de l'air » et « Maritime », chacune pointant vers un fichier de spécification différent. En savoir plus : [Collections](./collections).

**Étiquette (Tag)** — une catégorie de points d'accès dans une collection. Aide le LLM à trouver les bonnes opérations plus précisément. En savoir plus : [Étiquettes](./tags).

**Point d'accès (Endpoint)** — une méthode HTTP + chemin spécifique (par exemple, `GET /api/users`). Le LLM peut trouver un point d'accès par description, inspecter ses paramètres et schémas, puis l'invoquer. En savoir plus : [Points d'accès](./endpoints).

**Espace de travail (Workspace)** — le répertoire où swag2mcp stocke la configuration, le cache des specs, les réponses sauvegardées et les scripts d'authentification. En savoir plus : [Espace de travail](./workspace).

## Comment cela fonctionne

1. **Ajoutez une spec ou une collection** — définissez-la dans la configuration YAML (`~/.swag2mcp/swag2mcp.yaml`). Par exemple :

   ```yaml
   specs:
     - domain: jokes
       llm_title: API Dad Joke
       base_url: https://icanhazdadjoke.com
       collections:
         - llm_title: Blagues
           location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
   ```
2. **swag2mcp analyse chaque collection** — crée des étiquettes et des points d'accès, les indexe pour la recherche.
3. **Le LLM trouve le bon point d'accès** — via les outils MCP (`search`, `endpoint_by_tag`, `inspect`), le LLM recherche un point d'accès correspondant par description, examine ses paramètres et son schéma de requête.
4. **Le LLM invoque le point d'accès** — via l'outil MCP `invoke`, le LLM envoie la requête. swag2mcp valide chaque paramètre d'entrée par rapport au schéma OpenAPI du point d'accès (paramètres de chemin, paramètres de requête, en-têtes, corps de la requête) avant d'effectuer l'appel. Si quelque chose ne correspond pas au schéma, le LLM reçoit une erreur claire expliquant ce qui ne va pas. Une fois validé, swag2mcp exécute l'appel HTTP réel et retourne le résultat.
5. **Le résultat retourne au LLM** — la réponse API est transmise à l'agent. Les grandes réponses sont sauvegardées dans l'espace de travail et peuvent être explorées avec trois outils MCP dédiés : `response_outline` (voir la structure), `response_compress` (réduire à un échantillon représentatif) et `response_slice` (extraire des fragments spécifiques).

swag2mcp est un pont entre les LLM et le monde des API. Vous ajoutez des spécifications API, et le LLM — via le protocole MCP — trouve les bons points d'accès, inspecte leur documentation et les appelle. Tout ce que vous avez à faire est d'ajouter une spec et de démarrer le serveur MCP.

> **La configuration est modifiable à tout moment.** Le fichier de configuration YAML (`~/.swag2mcp/swag2mcp.yaml`) peut être modifié à la main — ajoutez des specs, changez l'authentification, ajustez les paramètres. Après chaque modification, redémarrez le serveur MCP (`swag2mcp mcp`) pour que les changements prennent effet.

## Hiérarchie

```
Spec (domaine, par ex. « meteo »)
  └── Collection 1 (fichier de spec, par ex. forecast.yml)
        └── Étiquette 1 (catégorie)
              └── Point d'accès (GET /api/forecast)
              └── Point d'accès (POST /api/forecast)
        └── Étiquette 2
              └── Point d'accès (GET /api/forecast/{id})
  └── Collection 2 (fichier de spec, par ex. air-quality.yml)
        └── Étiquette 3
              └── Point d'accès (GET /api/air-quality)
```
