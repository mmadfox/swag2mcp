# swag2mcp

Puente entre especificaciones de API OpenAPI/Swagger/Postman y agentes LLM mediante el Protocolo de Contexto de Modelo (MCP).

No todas las API hablan MCP — los endpoints privados, servicios internos, sistemas heredados y API de terceros rara vez lo hacen. swag2mcp envuelve cualquier API REST en una interfaz MCP, dando a los agentes LLM acceso instantáneo a toda su superficie API sin modificar una sola línea de código del servidor. A través de llamadas API en vivo, el LLM obtiene conocimiento del mundo real para tomar decisiones informadas, automatizar flujos de trabajo y actuar sobre sus datos — no solo adivinar.

<a href="https://www.youtube.com/watch?v=1Da4UmE2f9U" target="_blank">
  <img src="https://raw.githubusercontent.com/mmadfox/swag2mcp/main/docs/cover.jpg" alt="Vista previa">
</a>

## Su API habla con LLM

Una línea de configuración convierte cualquier archivo OpenAPI/Swagger/Postman en un servidor MCP. Los agentes LLM descubren, inspeccionan e invocan sus APIs — sin código de integración.

<img src="/architecture.svg" width="700" alt="Arquitectura de swag2mcp">

## Deje de escribir envoltorios

Cada vez que conecta una nueva API a un LLM, escribe el mismo código repetitivo: análisis de especificaciones, autenticación, manejo de errores, limitación de velocidad. swag2mcp lo hace por usted — 19 herramientas MCP listas para usar.

## Quién necesita esto

| Rol | Por qué |
|------|---------|
| **Desarrollador de Agentes IA** | Conecte cualquier API en 2 minutos, no en 2 días |
| **Ingeniero MCP** | Sin código de controlador — solo apunte a una especificación |
| **Arquitecto** | Capa única de integración de APIs para todos los LLM de su empresa |
| **Analista de Datos** | Acceda a APIs mediante lenguaje natural, sin programación |
| **DevOps / SRE** | Monitoreo y automatización a través de LLM sin servicios adicionales |
| **Integrador** | 9 métodos de autenticación listos para usar — desde Basic hasta OAuth2 y HMAC |
| **Ingeniero de QA** | Servidor simulado para pruebas aisladas sin APIs reales |
| **Gerente de Producto** | Prototipos rápidos de funciones IA sin trabajo de backend |
| **y muchos más** | |

---

## Licencia

Licenciado bajo la **Apache License, Version 2.0**.

Consulte [LICENSE](https://github.com/mmadfox/swag2mcp/blob/main/LICENSE) para el texto completo de la licencia.
