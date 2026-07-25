# Agregar un Nuevo Método de Autenticación

## Pasos

1. **Crear el cliente de autenticación** en `internal/auth/<nombre>.go`
2. **Implementar la interfaz `Authenticator`**
3. **Agregar constante de tipo** a `internal/auth/auth.go`
4. **Agregar decodificador YAML** a `internal/config/auth.go`
5. **Registrar decodificador** en el mapa `authDecoders`
6. **Escribir pruebas**

## 1. Cliente de autenticación

Cree `internal/auth/my_auth.go`:

```go
package auth

import "net/http"

type MyAuthClient struct {
    Token string `yaml:"token" validate:"required"`
}

func (c *MyAuthClient) New() error {
    c.Token = resolveEnv(c.Token)
    return nil
}

func (c *MyAuthClient) Type() Type {
    return MyAuth
}

func (c *MyAuthClient) Apply(req *http.Request, out *Info) error {
    if c.Token == "" {
        return nil
    }
    setAuthHeader(req, out, "X-My-Auth", c.Token)
    return nil
}

func (c *MyAuthClient) Validate() error {
    return authValidator.Struct(c)
}
```

## 2. Interfaz Authenticator

Cada cliente de autenticación debe implementar:

```go
type Authenticator interface {
    New() error                    // Inicializar, resolver variables de entorno
    Type() Type                    // Devolver el identificador de tipo de autenticación
    Apply(req *http.Request, out *Info) error  // Aplicar autenticación a la solicitud
    Validate() error               // Validar campos requeridos
}
```

## 3. Constante de tipo

Agregue a `internal/auth/auth.go`:

```go
const MyAuth Type = "my-auth"
```

## 4. Decodificador YAML

Agregue una función decodificadora en `internal/config/auth.go`. El decodificador recibe un `*yaml.Node` y debe decodificarlo en su estructura de cliente de autenticación:

```go
func decodeMyAuth(node *yaml.Node) (auth.Authenticator, error) {
    var client auth.MyAuthClient
    if err := decodeConfig(node, &client); err != nil {
        return nil, err
    }
    return &client, nil
}
```

El ayudante `decodeConfig` maneja el patrón común: verifica que el nodo no esté vacío, decodifica YAML en la estructura y devuelve un error descriptivo en caso de fallo.

## 5. Registrar decodificador

Agregue su decodificador al mapa `authDecoders` en `internal/config/auth.go`:

```go
var authDecoders = map[string]authDecoder{
    // ... decodificadores existentes
    auth.MyAuth.String(): decodeMyAuth,
}
```

El método `UnmarshalYAML` en `Auth` lee el campo `type` del YAML, normaliza los guiones bajos a guiones, busca el decodificador en `authDecoders` y lo llama con el nodo `config`. Así es como swag2mcp sabe qué cliente de autenticación instanciar para cada especificación.

## 6. Pruebas

Cree `internal/auth/my_auth_test.go` con pruebas basadas en tablas que cubran:

- `New()` resuelve las variables de entorno correctamente
- `Type()` devuelve el tipo correcto
- `Apply()` establece los encabezados/parámetros de consulta correctos
- `Apply()` maneja valores vacíos correctamente
- `Validate()` pasa para configuración válida
- `Validate()` falla para campos requeridos faltantes
