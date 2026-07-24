# Autenticación

## Descripción General

swag2mcp admite **9 métodos de autenticación** para trabajar con APIs que requieren autorización. Lo configura una vez en el archivo de configuración — después de eso, cada llamada a la API a través de `invoke` incluye automáticamente los tokens y encabezados correctos.

### Dónde configurar

La autenticación se establece a nivel de **especificación** en `swag2mcp.yaml`:

```yaml
specs:
  - domain: jokes
    llm_title: Dad Joke API
    base_url: https://icanhazdadjoke.com
    collections:
      - llm_title: Jokes
        location: https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/dadjoke.yaml
    auth:
      type: bearer
      config:
        token: "my-token"
```

### Cómo funciona

- Usted especifica el tipo de autenticación y los parámetros en la configuración
- swag2mcp los aplica automáticamente a cada solicitud cuando llama a `invoke`
- **No necesita** solicitar un token antes de llamar a una API — ocurre automáticamente
- Si un token expira (OAuth2, Script), swag2mcp lo renueva por su cuenta

### Variables de entorno

Los datos sensibles (tokens, contraseñas, claves) pueden almacenarse en variables de entorno usando la sintaxis `$(VAR_NAME)`:

```yaml
auth:
  type: bearer
  config:
    token: "$(MY_API_TOKEN)"
```

swag2mcp sustituye el valor de `MY_API_TOKEN` al iniciar.

### Herramienta MCP auth

El agente LLM puede recuperar un token o encabezados a través de la herramienta MCP `auth` — por ejemplo, para construir un comando curl o mostrarlo al usuario.

En **producción**, esta herramienta debe deshabilitarse con `--disable-llm-auth` (habilitado por defecto) para que el LLM nunca tenga acceso a los tokens.

### Métodos

| Método | Descripción | Mejor para |
|--------|-------------|------------|
| [`none`](/auth/none) | Sin autenticación | APIs públicas |
| [`basic`](/auth/basic) | HTTP Basic (usuario + contraseña) | APIs heredadas, autenticación simple |
| [`bearer`](/auth/bearer) | Token Bearer (JWT, token) | APIs REST modernas |
| [`api-key`](/auth/api-key) | Clave de API en encabezado o parámetro de consulta | Servicios con claves de API |
| [`digest`](/auth/digest) | HTTP Digest (usuario + contraseña) | APIs heredadas, más seguro que Basic |
| [`hmac`](/auth/hmac) | Firma HMAC-SHA256 (estilo Binance) | Exchanges de criptomonedas |
| [`oauth2-cc`](/auth/oauth2-cc) | Credenciales de Cliente OAuth2 | Servidor a servidor, microservicios |
| [`oauth2-pwd`](/auth/oauth2-pwd) | Concesión de Contraseña OAuth2 | Aplicaciones con inicio de sesión de usuario |
| [`script`](/auth/script) | Script externo para obtener un token | Cualquier esquema de autenticación personalizado |
