# version

## Propósito

Imprimir la versión de swag2mcp. Útil para verificar la versión instalada, reportar errores o verificar la compatibilidad.

## Cuándo usarlo

- Desea verificar qué versión de swag2mcp está instalada
- Está reportando un error y necesita incluir la versión
- Desea verificar una instalación exitosa

## Sintaxis

```bash
swag2mcp version
swag2mcp --version
```

## Argumentos

Ninguno.

## Banderas

Ninguna.

## Cómo funciona

```bash
swag2mcp version
# swag2mcp v1.2.0

swag2mcp --version
# swag2mcp v1.2.0
```

## Formato de salida

```
swag2mcp &lt;version&gt;
```

La versión se establece en tiempo de compilación mediante `ldflags`. Si no se establece, el valor predeterminado es `"dev"`.

## Matices

- **Dos formas:** Tanto `swag2mcp version` (subcomando) como `swag2mcp --version` (bandera global) producen la misma salida.
- **No requiere configuración:** Este comando funciona sin un espacio de trabajo o archivo de configuración.
