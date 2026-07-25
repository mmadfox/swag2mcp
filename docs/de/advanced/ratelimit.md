# Ratenbegrenzung

## Übersicht

swag2mcp verfügt über einen integrierten Ratenbegrenzer, der verhindert, dass der LLM denselben API-Endpunkt zu häufig aufruft. Dies schützt vor versehentlichen doppelten Aufrufen und respektiert API-Ratenlimits.

## Wie es funktioniert

Jeder Endpunkt hat eine Abklingzeit. Wenn der LLM versucht, denselben Endpunkt innerhalb der Abklingzeit erneut aufzurufen, wird der Aufruf mit einem strukturierten Fehler abgelehnt.

```
t=0s  → invoke(endpoint) → wird ausgeführt
t=2s  → invoke(endpoint) → mit rate_limit-Fehler abgelehnt
t=12s → invoke(endpoint) → wird ausgeführt (Abklingzeit abgelaufen)
```

### Standardverhalten

- **Abklingzeit:** 10 Sekunden pro Endpunkt
- **Gültigkeitsbereich:** Pro Endpunkt — der Aufruf von Endpunkt A hat keine Auswirkung auf Endpunkt B
- **Fehlerantwort:** Der LLM erhält einen `LLMError` mit dem Code `rate_limit` und einer Nachricht, die angibt, wie lange gewartet werden muss
- **Zurücksetzen:** Nach 10 Sekunden Inaktivität auf diesem Endpunkt

### Fehlerformat

Bei Ratenbegrenzung erhält der LLM:

```json
{
  "code": "rate_limit",
  "message": "Ratenlimit für Endpunkt \"abc123\" überschritten: versuchen Sie es in 8 Sekunden erneut",
  "hint": "Warten Sie, bis die Abklingzeit abgelaufen ist, und versuchen Sie dann, den Endpunkt erneut aufzurufen. Verwenden Sie das Suchtool, um andere Endpunkte zu finden, die Sie in der Zwischenzeit aufrufen können."
}
```

Der LLM kann diese Informationen nutzen, um zu warten und es erneut zu versuchen oder zu einem anderen Endpunkt zu wechseln.

### Warum es existiert

- **Verhindert versehentliche doppelte Aufrufe** — der LLM könnte denselben Endpunkt mehrmals schnell hintereinander aufrufen
- **Schützt vor API-Ratenlimits** — viele APIs haben eigene Ratenlimits, deren Überschreitung zu Fehlern führen würde
- **Spart Ressourcen** — reduziert unnötigen Netzwerkverkehr

## Konfiguration

Sie können den Ratenbegrenzer deaktivieren oder das Abklingintervall ändern:

```yaml
# Ratenbegrenzer vollständig deaktivieren
disable_ratelimiter: true

# Benutzerdefiniertes Abklingintervall
rate_limit_interval: 30s
```

### disable_ratelimiter

- **Typ:** `bool`
- **Standard:** `false`
- **Wirkung:** Wenn `true`, wird der Pro-Endpunkt-Ratenbegrenzer deaktiviert. Der LLM kann denselben Endpunkt wiederholt ohne Wartezeit aufrufen.
- **Wann aktivieren:** Testen, Debuggen oder wenn Sie denselben Endpunkt mehrmals schnell hintereinander aufrufen müssen.
- **Wann deaktiviert lassen (empfohlen):** Produktion. Der Ratenbegrenzer verhindert versehentlichen Missbrauch.

### rate_limit_interval

- **Typ:** Dauer (Go-Format: `10s`, `30s`, `1m`)
- **Standard:** `10s`
- **Wirkung:** Legt die Abklingzeit zwischen Aufrufen desselben Endpunkts fest.
- **Wann erhöhen:** APIs mit strengen Ratenlimits (z. B. 10 Anfragen pro Minute).
- **Wann verringern:** Interne APIs, bei denen Sie die Last kontrollieren.
- **Beispiele:** `5s`, `30s`, `1m`, `2m`

## Wichtige Hinweise

- **Pro-Endpunkt-Verfolgung** — jeder Endpunkt wird unabhängig verfolgt. Der Aufruf eines Endpunkts hat keine Auswirkung auf andere.
- **Fehler wird an LLM zurückgegeben** — der zweite Aufruf innerhalb der Abklingzeit wird mit einem `rate_limit`-Fehler abgelehnt. Der LLM erhält die Abklingdauer und kann nach dem Warten erneut versuchen.
- **Keine Bereinigung erforderlich** — der Ratenbegrenzer verfolgt Endpunkte automatisch und erfordert keine Wartung.
