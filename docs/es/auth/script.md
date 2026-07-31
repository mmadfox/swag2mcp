# Autenticación mediante Script

## Propósito

Autenticación mediante un script externo. swag2mcp admite dos tipos de scripts:
- **`.sh`** — scripts de shell en Linux y macOS
- **`.bat`** — scripts por lotes en Windows

El script debe colocarse en el directorio `auth_scripts` de su espacio de trabajo y generar un token JSON en stdout.

## Cuándo usarlo

- Esquemas de autenticación personalizados o no estándar
- Lógica compleja de obtención de tokens (multipaso, con verificaciones adicionales)
- Cuando ninguno de los métodos estándar se ajusta a sus necesidades

## Configuración

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: script
      config:
        domain: "jokes"
```

## Parámetros

| Parámetro | Requerido | Descripción |
|-----------|-----------|-------------|
| `domain` | Sí | Nombre del archivo de script (sin extensión). Debe coincidir con el `domain` de la spec — p. ej. si `domain: jokes`, el archivo debe ser `jokes.sh` (Unix) o `jokes.bat` (Windows). |

## Ubicación del script

El script debe colocarse en el directorio `auth_scripts` de su espacio de trabajo:

- **Linux / macOS:** `{workspace}/auth_scripts/{domain}.sh`
- **Windows:** `{workspace}/auth_scripts/{domain}.bat`

## Formato de salida del script

El script debe generar JSON en stdout con el token y su tiempo de expiración:

```bash
#!/bin/bash
# auth_scripts/jokes.sh

TOKEN=$(curl -s -X POST https://auth.example.com/token \
  -d "grant_type=client_credentials" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET" | jq -r '.access_token')

echo "{\"token\": \"$TOKEN\", \"expires_in\": 3600}"
```

### Campos JSON

| Campo | Requerido | Descripción |
|-------|-----------|-------------|
| `token` | Sí | Token de autenticación |
| `expires_in` | No | Vida útil del token en segundos (predeterminado: 3600) |

## Notas

- swag2mcp ejecuta el script en cada solicitud si el token en caché ha expirado
- El script debe completarse dentro de 30 segundos
- El token se almacena en caché hasta su tiempo de expiración
- Nombre del script = `{domain}.sh` (Unix) o `{domain}.bat` (Windows)
- `domain` no debe contener `/` ni `\`
