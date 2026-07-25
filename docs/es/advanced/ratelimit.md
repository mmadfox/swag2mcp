# Limitación de Velocidad

## Descripción General

swag2mcp tiene un limitador de velocidad incorporado que evita que el LLM llame al mismo endpoint de API con demasiada frecuencia. Esto protege contra llamadas duplicadas accidentales y respeta los límites de velocidad de la API.

## Cómo funciona

Cada endpoint tiene un período de enfriamiento. Si el LLM intenta llamar al mismo endpoint nuevamente dentro del período de enfriamiento, la llamada se rechaza con un error estructurado.

```
t=0s  → invoke(endpoint) → se ejecuta
t=2s  → invoke(endpoint) → rechazado con error rate_limit
t=12s → invoke(endpoint) → se ejecuta (el enfriamiento ha pasado)
```

### Comportamiento predeterminado

- **Enfriamiento:** 10 segundos por endpoint
- **Alcance:** Por endpoint — llamar al endpoint A no afecta al endpoint B
- **Respuesta de error:** El LLM recibe un `LLMError` con código `rate_limit` y un mensaje que indica cuánto tiempo esperar
- **Restablecimiento:** Después de 10 segundos de inactividad en ese endpoint

### Formato de error

Cuando se limita la velocidad, el LLM recibe:

```json
{
  "code": "rate_limit",
  "message": "límite de velocidad excedido para el endpoint \"abc123\": intente de nuevo en 8 segundos",
  "hint": "Espere a que expire el período de enfriamiento, luego intente invocar el endpoint nuevamente. Use la herramienta de búsqueda para encontrar otros endpoints que pueda llamar mientras tanto."
}
```

El LLM puede usar esta información para esperar y reintentar, o cambiar a un endpoint diferente.

### Por qué existe

- **Evita llamadas duplicadas accidentales** — el LLM podría llamar al mismo endpoint múltiples veces en rápida sucesión
- **Protege contra los límites de velocidad de la API** — muchas APIs tienen sus propios límites de velocidad, y alcanzarlos causaría errores
- **Ahorra recursos** — reduce el tráfico de red innecesario

## Configuración

Puede deshabilitar el limitador de velocidad o cambiar el intervalo de enfriamiento:

```yaml
# Deshabilitar el limitador de velocidad por completo
disable_ratelimiter: true

# Intervalo de enfriamiento personalizado
rate_limit_interval: 30s
```

### disable_ratelimiter

- **Tipo:** `bool`
- **Valor predeterminado:** `false`
- **Efecto:** Cuando es `true`, el limitador de velocidad por endpoint está deshabilitado. El LLM puede llamar al mismo endpoint repetidamente sin esperar.
- **Cuándo habilitar:** Pruebas, depuración, o cuando necesita llamar al mismo endpoint múltiples veces en rápida sucesión.
- **Cuándo mantenerlo deshabilitado (recomendado):** Producción. El limitador de velocidad evita el abuso accidental.

### rate_limit_interval

- **Tipo:** duración (formato Go: `10s`, `30s`, `1m`)
- **Valor predeterminado:** `10s`
- **Efecto:** Establece el período de enfriamiento entre llamadas al mismo endpoint.
- **Cuándo aumentar:** APIs con límites de velocidad estrictos (por ejemplo, 10 solicitudes por minuto).
- **Cuándo disminuir:** APIs internas donde usted controla la carga.
- **Ejemplos:** `5s`, `30s`, `1m`, `2m`

## Notas importantes

- **Seguimiento por endpoint** — cada endpoint se rastrea de forma independiente. Llamar a un endpoint no afecta a los demás.
- **Error devuelto al LLM** — la segunda llamada dentro del enfriamiento se rechaza con un error `rate_limit`. El LLM recibe la duración del enfriamiento y puede reintentar después de esperar.
- **No se necesita limpieza** — el limitador de velocidad rastrea los endpoints automáticamente y no requiere mantenimiento.
