# validate

## Propósito

Verificar el archivo de configuración y todos los archivos de especificación referenciados en busca de errores. Este es un comando de diagnóstico **de solo lectura** — nunca modifica nada.

## Cuándo usarlo

- Después de editar `swag2mcp.yaml` manualmente
- Antes de ejecutar `mcp` o `update` para detectar problemas temprano
- Al solucionar problemas de por qué una especificación no se carga
- En pipelines CI/CD para validar cambios de configuración

## Sintaxis

```bash
swag2mcp validate [path] [flags]
```

## Argumentos

| Argumento | Posición | Requerido | Descripción |
|-----------|----------|-----------|-------------|
| `path` | 1 | No | Directorio del espacio de trabajo. Si se omite, se resuelve mediante reglas de resolución de ruta. |

## Banderas

| Bandera | Abreviatura | Tipo | Valor predeterminado | Descripción |
|---------|-------------|------|---------------------|-------------|
| `--tags` | `-t` | `string` | `""` | Validar solo especificaciones con etiquetas coincidentes (separadas por comas) |

## Cómo funciona

```bash
swag2mcp validate
swag2mcp validate ./my-workspace
swag2mcp validate --tags=public
```

## Qué se verifica

| Verificación | Descripción |
|--------------|-------------|
| Sintaxis YAML | El archivo de configuración debe ser YAML válido |
| Estructura de configuración | Todos los campos requeridos presentes, tipos correctos |
| Unicidad de dominio | Sin dominios duplicados |
| Formato de dominio | Solo minúsculas, dígitos, guiones |
| Existencia del archivo de especificación | El archivo o URL de `location` debe ser accesible |
| Formato de especificación | El archivo debe ser OpenAPI 3.x, Swagger 2.0 o Postman collection válido |
| Configuración de autenticación | El tipo de autenticación y la configuración son válidos para el método seleccionado |
| Cliente HTTP | La configuración del cliente HTTP es válida |

## Qué NO se verifica

| No verificado | Razón |
|---------------|-------|
| Endpoints de autenticación | `validate` verifica la sintaxis de configuración de autenticación pero no prueba el inicio de sesión/intercambio de tokens |
| Disponibilidad de endpoints de API | Solo se verifica la URL del archivo de especificación, no la `base_url` |
| Corrección de `base_url` | Se valida el formato, pero no se realiza ninguna solicitud de prueba |
| Configuración del servidor simulado | `base_mock_url` no se verifica para conectividad |

## Ejemplo de salida

```
✅ Configuration is valid.
✓ Spec petstore: OK
✓ Spec meteo: OK
✗ Spec old-api: file not found
```

## Verificación posterior al comando

Si la validación pasa, la configuración está lista para `mcp`, `update` o `run`.

## Matices

- **Sin auto-inicio:** A diferencia de `add`, `ls` o `run`, `validate` **no** se auto-inicializa si falta la configuración. Devuelve un error: `"configuration not found at <path>"`.
- **Acceso a red:** Las URL de especificaciones remotas se obtienen durante la validación. El comando puede tardar más si las especificaciones están alojadas en servidores lentos.
- **Filtrado por etiquetas:** Cuando se establece `--tags`, solo se validan las especificaciones que coinciden con las etiquetas especificadas. Otras especificaciones se omiten.
