# run

## Propósito

Iniciar el explorador de API **TUI (Interfaz de Usuario de Terminal)** interactivo. Es una aplicación de pantalla completa para buscar, navegar, inspeccionar e invocar endpoints de API sin salir de la terminal.

## Cuándo usarlo

- Desea explorar sus APIs de forma interactiva
- Necesita buscar un endpoint específico en todas las especificaciones
- Desea navegar por la jerarquía especificación → colección → etiqueta → endpoint
- Desea probar una llamada a la API antes de configurar el servidor MCP

## Sintaxis

```bash
swag2mcp run [path]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |

## Banderas

Ninguna.

## Modos

### Modo de búsqueda

Búsqueda de texto completo en todos los endpoints de todas las especificaciones. Admite filtrado por método HTTP, etiqueta y ruta.

- Escriba una consulta para buscar nombres, rutas y descripciones de endpoints
- Filtre resultados por método (GET, POST, PUT, DELETE, etc.)
- Vea los detalles del endpoint con una sola pulsación de tecla

### Modo de navegación

Navegación en árbol a través de la jerarquía de especificaciones:

```
Especificación → Colección → Etiqueta → Endpoint
```

- Navegue hacia abajo en el árbol para encontrar endpoints específicos
- Vea los detalles del endpoint (parámetros, cuerpo de solicitud, respuestas)
- Invocar la API directamente desde la TUI

## Navegación

| Tecla | Acción |
|-------|--------|
| `↑` / `↓` | Navegar arriba/abajo |
| `Enter` | Seleccionar o abrir |
| `Esc` | Retroceder |
| `Tab` | Cambiar entre los modos Búsqueda y Navegación |
| `/` | Enfocar entrada de búsqueda |
| `q` | Salir |

## Verificación posterior al comando

La TUI carga todas las especificaciones del espacio de trabajo. Si una especificación falla al cargar, se muestra un mensaje de error en la interfaz.

## Matices

- **Auto-inicio:** Si no existe un archivo de configuración, `run` ejecuta automáticamente el asistente de inicio primero.
- **Sin banderas:** El comando `run` no tiene banderas — toda la configuración proviene del espacio de trabajo.
- **Tamaño de terminal:** La TUI requiere una terminal con al menos 80×24 caracteres. Puede no renderizarse correctamente en terminales muy pequeñas.
- **Dependencias:** La TUI usa Bubbletea. Funciona a través de SSH y en la mayoría de emuladores de terminal.
