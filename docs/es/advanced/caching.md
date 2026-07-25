# Almacenamiento en Caché

## Descripción General

swag2mcp almacena en caché los archivos de especificación descargados para que el servidor MCP se inicie más rápido en ejecuciones posteriores. En lugar de descargar el mismo archivo de especificación cada vez, reutiliza la copia en caché.

## Cómo funciona el almacenamiento en caché

Cuando agrega una especificación con una URL remota, swag2mcp la descarga y la guarda en el directorio `cache/`. En el siguiente inicio, verifica si la copia en caché sigue siendo válida. Si lo es, se omite la descarga.

### Qué se almacena en caché

| Origen | Comportamiento |
|--------|---------------|
| **URL remota** (http/https) | Siempre se almacena en caché. Se descarga una vez, se reutiliza hasta que la caché expire. |
| **Archivo local en `specs/`** | Se usa directamente desde el directorio `specs/`. Nunca se almacena en caché — los cambios son inmediatamente visibles. |
| **Archivo local fuera de `specs/`** | Se copia a la caché. Si el archivo fuente cambia (hora de modificación), la caché se invalida. |

### Expiración de la caché (TTL)

Cada archivo en caché recibe un tiempo de expiración aleatorio entre **1 hora y 48 horas**. La aleatoriedad evita que todos los archivos en caché expiren al mismo tiempo (lo que causaría una avalancha de descargas).

- El TTL se restablece cada vez que se inicia el servidor MCP
- Si un archivo en caché aún está dentro de su TTL, se reutiliza
- Si el TTL ha expirado, el archivo se descarga nuevamente

### Estructura de la caché

```
~/.swag2mcp/cache/
├── a1b2c3d4e5f6a7b8.spec    # Archivo de especificación en caché
├── a1b2c3d4e5f6a7b8.meta    # Metadatos (origen, TTL, momento de almacenamiento)
├── b2c3d4e5f6a7b8c9.spec
├── b2c3d4e5f6a7b8c9.meta
└── ...
```

La clave de caché se deriva de la URL o ruta del archivo de especificación. Cada archivo en caché tiene un archivo `.meta` acompañante que almacena cuándo se almacenó en caché y cuándo expira.

## Gestión de la caché

### Forzar una actualización

Ejecute `swag2mcp update` para limpiar toda la caché y volver a descargar todos los archivos de especificación:

```bash
swag2mcp update
```

Esto valida la configuración, limpia la caché y descarga todo nuevamente.

### Limpiar la caché manualmente

```bash
swag2mcp clean
```

Esto elimina todos los archivos de especificación en caché y las respuestas de API guardadas. La próxima vez que inicie el servidor MCP, todas las especificaciones se descargarán nuevamente.

### Limpieza automática

Cuando se inicia el servidor MCP (`swag2mcp mcp`), las respuestas de API guardadas con más de 48 horas se eliminan automáticamente. Esto evita que el directorio `responses/` crezca indefinidamente.

## Notas importantes

- **Los archivos locales en `specs/` nunca se almacenan en caché** — si edita un archivo de especificación directamente en el directorio `specs/`, los cambios son inmediatamente visibles sin limpiar la caché
- **Las URL remotas siempre se almacenan en caché** — no hay forma de omitir la caché para URL remotas excepto ejecutando `swag2mcp update` o `swag2mcp clean`
- **La caché es local** — se almacena en disco y no se sincroniza entre máquinas. Use `swag2mcp export` y `swag2mcp import` para transferir especificaciones entre máquinas
