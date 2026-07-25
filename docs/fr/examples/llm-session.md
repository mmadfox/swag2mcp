# Exemples de sessions LLM

Cette page présente des dialogues réels entre un utilisateur et un agent LLM propulsé par swag2mcp. Chaque exemple inclut la demande de l'utilisateur, l'appel d'outil MCP et la réponse.

## Configuration

L'espace de travail a deux API configurées :

- **icanhazdadjoke** — API publique de blagues de papa
- **Open-Meteo** — API de prévisions météorologiques

---

## Découverte des API disponibles

**Utilisateur :** À quelles API avez-vous accès ?

**LLM → MCP :** `spec_list()`

**MCP → LLM :**
```
Spécifications :
  dadjoke (icanhazdadjoke)
  meteo (API météo Open-Meteo)
```

**LLM :** J'ai accès à deux API : icanhazdadjoke pour les blagues de papa et Open-Meteo pour les données météorologiques.

---

## Obtention des détails d'une spécification

**Utilisateur :** Parlez-moi plus de l'API dadjoke.

**LLM → MCP :** `spec_by_id(id: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM :**
```
Spécification : dadjoke (icanhazdadjoke)
Collections :
  - Jokes (3 points de terminaison, 1 balise)
```

**LLM :** L'API dadjoke a une collection appelée « Jokes » avec 3 points de terminaison sous une seule balise.

---

## Recherche de points de terminaison

**Utilisateur :** Trouvez-moi un point de terminaison pour obtenir une blague de papa aléatoire.

**LLM → MCP :** `search(query: "blague aléatoire", limit: 5)`

**MCP → LLM :**
```
Points de terminaison :
  GET /  — Renvoyer une blague de papa aléatoire
  GET /j/{id} — Renvoyer une blague de papa spécifique par ID
  GET /search — Rechercher des blagues de papa
```

**LLM :** J'ai trouvé trois points de terminaison. Le principal est `GET /` qui renvoie une blague de papa aléatoire.

---

## Inspection d'un point de terminaison

**Utilisateur :** De quels paramètres le point de terminaison de blague aléatoire a-t-il besoin ?

**LLM → MCP :** `inspect(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM :**
```
GET /
URL de base : https://icanhazdadjoke.com
URL complète : https://icanhazdadjoke.com/

Paramètres : aucun

Réponses :
  200 : Une blague de papa aléatoire
    En-têtes : Accept: text/plain
```

**LLM :** Le point de terminaison ne prend aucun paramètre. Il suffit de l'appeler et vous obtiendrez une blague aléatoire.

---

## Invocation d'une API

**Utilisateur :** Donnez-moi une blague de papa aléatoire !

**LLM → MCP :** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM :**
```json
{
  "statusCode": 200,
  "body": "Pourquoi les scientifiques ne font-ils pas confiance aux atomes ? Parce qu'ils inventent tout !"
}
```

**LLM :** Voici votre blague de papa : « Pourquoi les scientifiques ne font-ils pas confiance aux atomes ? Parce qu'ils inventent tout ! »

---

## Travail avec de grandes réponses

**Utilisateur :** Obtenez-moi les prévisions météorologiques pour Paris.

**LLM → MCP :** `search(query: "prévisions météorologiques", limit: 5)`

**MCP → LLM :**
```
Points de terminaison :
  GET /v1/forecast — Prévisions météorologiques
```

**LLM → MCP :** `inspect(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7")`

**MCP → LLM :**
```
GET /v1/forecast
Paramètres :
  latitude (query, nombre, obligatoire)
  longitude (query, nombre, obligatoire)
  hourly (query, chaîne, optionnel)

Réponses :
  200 : Données de prévisions
```

**LLM → MCP :** `invoke(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7", parameters: {"latitude": 48.85, "longitude": 2.35, "hourly": "temperature_2m"})`

**MCP → LLM :**
```json
{
  "statusCode": 200,
  "fileRef": {
    "path": "/home/utilisateur/.swag2mcp/responses/response_x7k9m2.json",
    "size": 1572864,
    "sizeHint": "1.5 Mo",
    "maxSizeHint": "2 Ko",
    "message": "La réponse dépasse la limite de 2 Ko et a été enregistrée sur le disque."
  }
}
```

**LLM :** La réponse est volumineuse (1.5 Mo). Laissez-moi explorer sa structure.

**LLM → MCP :** `response_outline(path: "/home/utilisateur/.swag2mcp/responses/response_x7k9m2.json")`

**MCP → LLM :**
```
Type : object
Taille : 1.5 Mo
Clés : latitude, longitude, hourly, daily
Tableaux :
  hourly.time (168 éléments)
  hourly.temperature_2m (168 éléments)
```

**LLM → MCP :** `response_compress(path: "/home/utilisateur/.swag2mcp/responses/response_x7k9m2.json", mode: "sample_array", jsonPath: "hourly", arrayHead: 3, arrayTail: 2)`

**MCP → LLM :**
```json
{
  "compressed": {
    "type": "array",
    "length": 5,
    "original": 168,
    "skipped": 163,
    "sample": [
      {"time": "2026-07-24T00:00", "temperature_2m": 18.5},
      {"time": "2026-07-24T01:00", "temperature_2m": 17.8},
      {"time": "2026-07-24T02:00", "temperature_2m": 16.9},
      {"time": "2026-07-24T22:00", "temperature_2m": 20.1},
      {"time": "2026-07-24T23:00", "temperature_2m": 19.4}
    ]
  }
}
```

**LLM :** Voici les prévisions météorologiques pour Paris. Les températures varient de 16,9°C la nuit à environ 20°C le soir.

---

## Récupération d'un jeton d'authentification

**Utilisateur :** Affichez-moi le jeton d'authentification pour mon API afin que je puisse la tester avec curl.

**LLM → MCP :** `auth(specId: "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6")`

**MCP → LLM :**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**LLM :** Voici le jeton. Vous pouvez l'utiliser dans curl :

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." https://api.example.com/v1/users
```

---

## Informations d'exécution

**Utilisateur :** Quel est l'état actuel du système ?

**LLM → MCP :** `info()`

**MCP → LLM :**
```json
{
  "version": "v1.2.0",
  "workspace": "/home/utilisateur/.swag2mcp",
  "uptime": "2h 15m",
  "specs": {
    "total": 2,
    "active": 2,
    "endpoints": 20
  },
  "http_client": {
    "timeout": "30s",
    "max_response_size": "2 Ko"
  },
  "mcp": {
    "transport": "stdio"
  }
}
```

**LLM :** Le système fonctionne depuis 2 heures et 15 minutes. Deux spécifications sont actives avec 20 points de terminaison au total. Le client HTTP a un délai d'attente de 30 secondes et une limite de taille de réponse de 2 Ko.

---

## Conversation complète : planification d'un voyage

**Utilisateur :** Je planifie un voyage à Tokyo la semaine prochaine. Pouvez-vous m'aider à me préparer ?

**LLM → MCP :** `search(query: "prévisions météorologiques", limit: 5)`

**MCP → LLM :**
```
Points de terminaison :
  GET /v1/forecast — Prévisions météorologiques
```

**LLM → MCP :** `inspect(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7")`

**MCP → LLM :**
```
GET /v1/forecast
Paramètres :
  latitude (query, nombre, obligatoire)
  longitude (query, nombre, obligatoire)
  daily (query, chaîne, optionnel)
```

**LLM → MCP :** `invoke(endpointId: "e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7", parameters: {"latitude": 35.68, "longitude": 139.69, "daily": "temperature_2m_max,temperature_2m_min,precipitation_sum"})`

**MCP → LLM :**
```json
{
  "statusCode": 200,
  "body": {
    "latitude": 35.68,
    "longitude": 139.69,
    "daily": {
      "time": ["2026-07-28", "2026-07-29", "2026-07-30", "2026-07-31", "2026-08-01"],
      "temperature_2m_max": [32, 33, 31, 30, 29],
      "temperature_2m_min": [25, 26, 24, 23, 22],
      "precipitation_sum": [0.0, 0.5, 2.1, 0.0, 0.0]
    }
  }
}
```

**LLM :** Voici les prévisions météorologiques pour Tokyo pour la semaine prochaine :

| Date | Max | Min | Pluie |
|------|-----|-----|------|
| 28 juil. | 32°C | 25°C | 0 mm |
| 29 juil. | 33°C | 26°C | 0,5 mm |
| 30 juil. | 31°C | 24°C | 2,1 mm |
| 31 juil. | 30°C | 23°C | 0 mm |
| 1 août | 29°C | 22°C | 0 mm |

Il fera chaud et le temps sera généralement sec. Prévoyez des vêtements légers et un parapluie au cas où pour le 30 juillet !

---

## Conversation complète : trouver une blague

**Utilisateur :** Dites-moi quelque chose de drôle !

**LLM → MCP :** `spec_list()`

**MCP → LLM :**
```
Spécifications :
  dadjoke (icanhazdadjoke)
  meteo (API météo Open-Meteo)
```

**LLM → MCP :** `search(query: "blague aléatoire", limit: 5)`

**MCP → LLM :**
```
Points de terminaison :
  GET /  — Renvoyer une blague de papa aléatoire
```

**LLM → MCP :** `invoke(endpointId: "f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6")`

**MCP → LLM :**
```json
{
  "statusCode": 200,
  "body": "J'ai dit à ma femme qu'elle dessinait ses sourcils trop haut. Elle avait l'air surprise."
}
```

**LLM :** La voici : « J'ai dit à ma femme qu'elle dessinait ses sourcils trop haut. Elle avait l'air surprise. »
