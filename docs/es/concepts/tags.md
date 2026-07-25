# Etiquetas

Una etiqueta es una categoría que agrupa endpoints relacionados dentro de una colección. Las etiquetas pueden existir o no — no todas las colecciones las tienen, y una colección puede tener cualquier número de etiquetas.

Las etiquetas provienen del propio archivo OpenAPI/Swagger/Postman. **No hay configuraciones YAML** para las etiquetas — no puede crear, renombrar ni eliminar etiquetas en `swag2mcp.yaml`. La única forma de cambiar las etiquetas es editar el archivo de especificación original.

## Jerarquía

```
Especificación (dominio, por ejemplo "meteo")
  └── Colección (archivo de especificación, por ejemplo forecast.yml)
        └── Etiqueta "weather"
              └── GET /forecast
              └── GET /forecast/hourly
        └── Etiqueta "alerts"
              └── GET /alerts
```

## Cómo se Crean las Etiquetas

Las etiquetas se extraen del documento de especificación durante el análisis:

**OpenAPI 3.x / Swagger 2.0** — la lista `tags` de cada operación se convierte en etiquetas:

```yaml
paths:
  /pet:
    get:
      tags: ["pets"]
      summary: "Find pet by ID"
    post:
      tags: ["pets"]
      summary: "Add a new pet"
  /pet/{petId}/uploadImage:
    post:
      tags: ["pet_images"]
      summary: "Uploads an image"
```

**Postman** — cada carpeta de nivel superior se convierte en una etiqueta. Las carpetas anidadas usan el nombre de la última carpeta.

Si un endpoint no tiene etiquetas, se coloca bajo una etiqueta `"default"`.

## Propósito

Las etiquetas ayudan al LLM a encontrar grupos de endpoints relacionados. En lugar de buscar en todos los endpoints de una colección, el LLM puede primero encontrar la etiqueta correcta, luego listar solo los endpoints dentro de ella.

## Herramientas MCP para Etiquetas

| Herramienta | Descripción |
|-------------|-------------|
| `tag_by_spec` | Todas las etiquetas en una especificación completa |
| `tag_by_collection` | Etiquetas dentro de una colección específica |
| `tag_by_id` | Detalles de la etiqueta (título, recuento de métodos) |
| `endpoint_by_tag` | Endpoints agrupados bajo una etiqueta |

## Ejemplo

```
Consulta: "Muestra todas las etiquetas en la colección de mascotas"
→ tag_by_collection(collectionId: "...")
→ Resultado: pets (5 métodos), pet_images (1 método)
```

## Limitaciones

- Las etiquetas son de solo lectura desde la perspectiva de la configuración. Para agregar, renombrar o eliminar etiquetas, edite el archivo OpenAPI/Swagger/Postman original y ejecute `swag2mcp update`.
- Las etiquetas no se pueden filtrar ni deshabilitar por colección en la configuración YAML.
