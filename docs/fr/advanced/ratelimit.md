# Limitation de débit

## Aperçu

swag2mcp dispose d'un limiteur de débit intégré qui empêche le LLM d'appeler le même point d'accès API trop fréquemment. Cela protège contre les appels en double accidentels et respecte les limites de débit des API.

## Comment cela fonctionne

Chaque point d'accès a une période de refroidissement. Si le LLM tente d'appeler le même point d'accès pendant le refroidissement, l'appel est rejeté avec une erreur structurée.

```
t=0s  → invoke(point d'accès) → exécute
t=2s  → invoke(point d'accès) → rejeté avec erreur rate_limit
t=12s → invoke(point d'accès) → exécute (le refroidissement est passé)
```

### Comportement par défaut

- **Refroidissement :** 10 secondes par point d'accès
- **Portée :** Par point d'accès — appeler le point d'accès A n'affecte pas le point d'accès B
- **Réponse d'erreur :** Le LLM reçoit une `LLMError` avec le code `rate_limit` et un message indiquant combien de temps attendre
- **Réinitialisation :** Après 10 secondes d'inactivité sur ce point d'accès

### Format d'erreur

Lorsque la limite de débit est atteinte, le LLM reçoit :

```json
{
  "code": "rate_limit",
  "message": "limite de débit dépassée pour le point d'accès \"abc123\" : réessayez dans 8 secondes",
  "hint": "Attendez la fin de la période de refroidissement, puis essayez à nouveau d'invoquer le point d'accès. Utilisez l'outil de recherche pour trouver d'autres points d'accès que vous pouvez appeler entre-temps."
}
```

Le LLM peut utiliser ces informations pour attendre et réessayer, ou passer à un point d'accès différent.

### Pourquoi cela existe

- **Empêche les appels en double accidentels** — le LLM pourrait appeler le même point d'accès plusieurs fois en succession rapide
- **Protège contre les limites de débit des API** — de nombreuses API ont leurs propres limites de débit, et les atteindre provoquerait des erreurs
- **Économise les ressources** — réduit le trafic réseau inutile

## Configuration

Vous pouvez désactiver le limiteur de débit ou modifier l'intervalle de refroidissement :

```yaml
# Désactiver complètement le limiteur de débit
disable_ratelimiter: true

# Intervalle de refroidissement personnalisé
rate_limit_interval: 30s
```

### disable_ratelimiter

- **Type :** `bool`
- **Défaut :** `false`
- **Effet :** Lorsque `true`, le limiteur de débit par point d'accès est désactivé. Le LLM peut appeler le même point d'accès de manière répétée sans attendre.
- **Quand l'activer :** Tests, débogage, ou lorsque vous devez appeler le même point d'accès plusieurs fois en succession rapide.
- **Quand le laisser désactivé (recommandé) :** Production. Le limiteur de débit empêche les abus accidentels.

### rate_limit_interval

- **Type :** durée (format Go : `10s`, `30s`, `1m`)
- **Défaut :** `10s`
- **Effet :** Définit la période de refroidissement entre les appels au même point d'accès.
- **Quand augmenter :** API avec des limites de débit strictes (par exemple, 10 requêtes par minute).
- **Quand diminuer :** API internes où vous contrôlez la charge.
- **Exemples :** `5s`, `30s`, `1m`, `2m`

## Notes importantes

- **Suivi par point d'accès** — chaque point d'accès est suivi indépendamment. Appeler un point d'accès n'affecte pas les autres.
- **Erreur retournée au LLM** — le second appel pendant le refroidissement est rejeté avec une erreur `rate_limit`. Le LLM reçoit la durée du refroidissement et peut réessayer après avoir attendu.
- **Aucun nettoyage nécessaire** — le limiteur de débit suit les points d'accès automatiquement et ne nécessite pas de maintenance.
