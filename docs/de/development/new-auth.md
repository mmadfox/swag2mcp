# Hinzufügen einer neuen Auth-Methode

## Schritte

1. **Auth-Client erstellen** in `internal/auth/<name>.go`
2. **`Authenticator`-Interface implementieren**
3. **Typ-Konstante hinzufügen** in `internal/auth/auth.go`
4. **YAML-Decoder hinzufügen** in `internal/config/auth.go`
5. **Decoder registrieren** in der `authDecoders`-Map
6. **Tests schreiben**

## 1. Auth-Client

Erstellen Sie `internal/auth/my_auth.go`:

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

## 2. Authenticator-Interface

Jeder Auth-Client muss implementieren:

```go
type Authenticator interface {
    New() error                    // Initialisieren, Umgebungsvariablen auflösen
    Type() Type                    // Den Auth-Typ-Identifikator zurückgeben
    Apply(req *http.Request, out *Info) error  // Auth auf Anfrage anwenden
    Validate() error               // Erforderliche Felder validieren
}
```

## 3. Typ-Konstante

Fügen Sie in `internal/auth/auth.go` hinzu:

```go
const MyAuth Type = "my-auth"
```

## 4. YAML-Decoder

Fügen Sie eine Decoder-Funktion in `internal/config/auth.go` hinzu. Der Decoder empfängt einen `*yaml.Node` und muss ihn in Ihre Auth-Client-Struktur dekodieren:

```go
func decodeMyAuth(node *yaml.Node) (auth.Authenticator, error) {
    var client auth.MyAuthClient
    if err := decodeConfig(node, &client); err != nil {
        return nil, err
    }
    return &client, nil
}
```

Der Helfer `decodeConfig` behandelt das übliche Muster: er prüft, dass der Knoten nicht leer ist, dekodiert YAML in die Struktur und gibt bei Fehlschlag einen beschreibenden Fehler zurück.

## 5. Decoder registrieren

Fügen Sie Ihren Decoder zur `authDecoders`-Map in `internal/config/auth.go` hinzu:

```go
var authDecoders = map[string]authDecoder{
    // ... bestehende Decoder
    auth.MyAuth.String(): decodeMyAuth,
}
```

Die `UnmarshalYAML`-Methode auf `Auth` liest das `type`-Feld aus dem YAML, normalisiert Unterstriche zu Bindestrichen, sucht den Decoder in `authDecoders` und ruft ihn mit dem `config`-Knoten auf. So weiß swag2mcp, welchen Auth-Client es für jede Spec instanziieren muss.

## 6. Tests

Erstellen Sie `internal/auth/my_auth_test.go` mit tabellengesteuerten Tests, die Folgendes abdecken:

- `New()` löst Umgebungsvariablen korrekt auf
- `Type()` gibt den korrekten Typ zurück
- `Apply()` setzt die richtigen Header/Abfrageparameter
- `Apply()` behandelt leere Werte ordnungsgemäß
- `Validate()` besteht für gültige Konfiguration
- `Validate()` schlägt fehl für fehlende erforderliche Felder
