# Ajouter une nouvelle méthode d'authentification

## Étapes

1. **Créer le client d'authentification** dans `internal/auth/&lt;nom&gt;.go`
2. **Implémenter l'interface `Authenticator`**
3. **Ajouter une constante de type** dans `internal/auth/auth.go`
4. **Ajouter un décodeur YAML** dans `internal/config/auth.go`
5. **Enregistrer le décodeur** dans la map `authDecoders`
6. **Écrire les tests**

## 1. Client d'authentification

Créez `internal/auth/mon_auth.go` :

```go
package auth

import "net/http"

type MonAuthClient struct {
    Token string `yaml:"token" validate:"required"`
}

func (c *MonAuthClient) New() error {
    c.Token = resolveEnv(c.Token)
    return nil
}

func (c *MonAuthClient) Type() Type {
    return MonAuth
}

func (c *MonAuthClient) Apply(req *http.Request, out *Info) error {
    if c.Token == "" {
        return nil
    }
    setAuthHeader(req, out, "X-Mon-Auth", c.Token)
    return nil
}

func (c *MonAuthClient) Validate() error {
    return authValidator.Struct(c)
}
```

## 2. Interface Authenticator

Chaque client d'authentification doit implémenter :

```go
type Authenticator interface {
    New() error                    // Initialiser, résoudre les variables d'environnement
    Type() Type                    // Renvoyer l'identifiant du type d'authentification
    Apply(req *http.Request, out *Info) error  // Appliquer l'authentification à la requête
    Validate() error               // Valider les champs obligatoires
}
```

## 3. Constante de type

Ajoutez à `internal/auth/auth.go` :

```go
const MonAuth Type = "mon-auth"
```

## 4. Décodeur YAML

Ajoutez une fonction de décodage dans `internal/config/auth.go`. Le décodeur reçoit un `*yaml.Node` et doit le décoder dans votre structure de client d'authentification :

```go
func decodeMonAuth(node *yaml.Node) (auth.Authenticator, error) {
    var client auth.MonAuthClient
    if err := decodeConfig(node, &client); err != nil {
        return nil, err
    }
    return &client, nil
}
```

Le helper `decodeConfig` gère le modèle commun : il vérifie que le nœud n'est pas vide, décode le YAML dans la structure et renvoie une erreur descriptive en cas d'échec.

## 5. Enregistrer le décodeur

Ajoutez votre décodeur à la map `authDecoders` dans `internal/config/auth.go` :

```go
var authDecoders = map[string]authDecoder{
    // ... décodeurs existants
    auth.MonAuth.String(): decodeMonAuth,
}
```

La méthode `UnmarshalYAML` sur `Auth` lit le champ `type` du YAML, normalise les underscores en traits d'union, recherche le décodeur dans `authDecoders` et l'appelle avec le nœud `config`. C'est ainsi que swag2mcp sait quel client d'authentification instancier pour chaque spécification.

## 6. Tests

Créez `internal/auth/mon_auth_test.go` avec des tests pilotés par tableaux couvrant :

- `New()` résout correctement les variables d'environnement
- `Type()` renvoie le type correct
- `Apply()` définit les bons en-têtes/paramètres de requête
- `Apply()` gère les valeurs vides avec élégance
- `Validate()` réussit pour une configuration valide
- `Validate()` échoue pour les champs obligatoires manquants
