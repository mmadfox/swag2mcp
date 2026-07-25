# Estructura del Proyecto

```
swag2mcp/
├── cmd/
│   ├── swag2mcp/          # Binario principal
│   │   └── main.go
│   └── swag2mcp-mock/     # Servidor simulado
│       └── main.go
├── internal/
│   ├── auth/              # 9 métodos de autenticación
│   ├── cache/             # Almacenamiento en caché de especificaciones
│   ├── commands/          # 13 comandos CLI (cobra)
│   ├── config/            # Configuración YAML
│   ├── env/               # Variables de entorno
│   ├── httpclient/        # Cliente HTTP
│   ├── id/                # Generación de IDs MD5
│   ├── index/             # Búsqueda de texto completo (bluge)
│   ├── model/             # Modelos de datos
│   ├── reader/            # Lectura de respuestas grandes
│   ├── server/
│   │   ├── mcp/           # Servidor MCP (19 herramientas)
│   │   └── mockserver/    # Servidor simulado
│   ├── service/           # Lógica de negocio
│   ├── spec/              # Analizadores de especificaciones
│   ├── tui/               # Interfaz TUI
│   └── workspace/         # Gestión del espacio de trabajo
├── specs/                 # Especificaciones de ejemplo
├── tests/                 # Pruebas de integración
├── docs/                  # Documentación
├── examples/              # Ejemplos de configuración
└── playground/            # Entorno de desarrollo
```

## Paquetes Clave

| Paquete | Descripción |
|---------|-------------|
| `auth` | 9 métodos de autenticación |
| `cache` | Almacenamiento en caché basado en disco con TTL |
| `commands` | Comandos CLI Cobra |
| `config` | Configuración YAML con cascada |
| `httpclient` | Cliente HTTP configurable |
| `index` | Búsqueda de texto completo (bluge) |
| `server/mcp` | Servidor MCP (3 transportes) |
| `service` | Lógica de negocio (núcleo) |
| `spec` | Analizadores OpenAPI/Swagger/Postman |
| `tui` | TUI Bubbletea |
| `workspace` | Gestión de archivos |
