# Explorador TUI

## Descripción General

swag2mcp incluye una TUI (Interfaz de Usuario de Terminal) incorporada para la exploración interactiva de APIs. Es una aplicación de terminal de pantalla completa que le permite buscar, navegar, inspeccionar e invocar endpoints de API sin salir de la terminal.

## Inicio

```bash
swag2mcp run
```

Si no existe un archivo de configuración, la TUI iniciará automáticamente el asistente de inicialización primero.

## Modos

La TUI tiene tres modos, intercambiables con la tecla `Tab`:

### Modo de búsqueda

Búsqueda de texto completo en todos los endpoints de todas las especificaciones. Admite la misma sintaxis de consulta que la herramienta MCP `search`.

- Escriba una consulta para buscar nombres, rutas y descripciones de endpoints
- Filtre resultados por método, etiqueta o ruta
- Vea los detalles del endpoint con una sola pulsación de tecla
- Navegue por los resultados con paginación (10 elementos por página)

### Modo de navegación

Navegación en árbol a través de la jerarquía de especificaciones:

```
Especificación → Colección → Etiqueta → Endpoint
```

- Navegue hacia abajo en el árbol para encontrar endpoints específicos
- Vea los detalles del endpoint (parámetros, cuerpo de solicitud, respuestas)
- Invocar la API directamente desde la TUI
- Guarde los detalles del endpoint como un archivo JSON

### Modo de autenticación

Vea los tokens de autenticación y encabezados para cualquier especificación. Útil para depuración o generación de comandos curl.

## Controles

| Tecla | Acción |
|-------|--------|
| `↑` / `↓` | Navegar arriba/abajo |
| `Enter` | Seleccionar o abrir |
| `Esc` | Retroceder un nivel |
| `Tab` | Cambiar entre los modos Búsqueda, Navegación y Autenticación |
| `/` | Enfocar entrada de búsqueda |
| `N` / `P` | Siguiente / página anterior |
| `B` | Volver a la pantalla anterior |
| `M` | Volver al menú principal |
| `S` | Guardar detalle del endpoint como archivo JSON |
| `q` / `Ctrl+C` | Salir |

## Estados

La TUI pasa por estos estados a medida que navega:

1. **Cargando** — cargando datos del espacio de trabajo
2. **Búsqueda** — modo de búsqueda con entrada de consulta
3. **Navegación** — modo de navegación con lista de especificaciones
4. **Lista de Especificaciones** — lista de todas las especificaciones
5. **Lista de Colecciones** — colecciones dentro de una especificación
6. **Lista de Etiquetas** — etiquetas dentro de una colección
7. **Lista de Endpoints** — endpoints dentro de una etiqueta
8. **Detalle del Endpoint** — información completa del endpoint
9. **Resultado de Invocación** — resultado de la llamada a la API
10. **Error** — estado de error con mensaje

## Vista de detalle del endpoint

Cuando selecciona un endpoint, la TUI muestra:

- Método HTTP y ruta
- URL base y URL completa
- Resumen y descripción
- Todos los parámetros (nombre, ubicación, tipo, requerido)
- Esquema del cuerpo de solicitud (si aplica)
- Códigos de respuesta y esquemas
- Estado de obsolescencia

## Requisitos

- **Tamaño de terminal:** Al menos 80×24 caracteres
- **Emulador de terminal:** Funciona en la mayoría de terminales modernas (iTerm2, Terminal.app, GNOME Terminal, Windows Terminal, etc.)
- **SSH:** Funciona a través de conexiones SSH

## Notas importantes

- **Auto-inicio** — si no existe un archivo de configuración, la TUI inicia automáticamente el asistente de inicialización
- **Paginación** — las listas se paginan a 10 elementos por página. Use `N` y `P` para navegar
- **Guardar detalles del endpoint** — presione `S` en la vista de detalle del endpoint para guardar el detalle completo como un archivo JSON en el directorio actual
- **Modo de autenticación** — muestra tokens y encabezados para depuración. En producción, la herramienta de autenticación puede deshabilitarse con `--disable-llm-auth`
