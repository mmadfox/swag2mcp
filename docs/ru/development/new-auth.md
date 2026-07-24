# Добавление нового метода аутентификации

## Шаги

1. **Создайте клиент аутентификации** в `internal/auth/<name>.go`
2. **Реализуйте интерфейс `Authenticator`**
3. **Добавьте константу типа** в `internal/auth/auth.go`
4. **Добавьте YAML-декодер** в `internal/config/auth.go`
5. **Зарегистрируйте декодер** в карте `authDecoders`
6. **Напишите тесты**

## 1. Клиент аутентификации

Создайте `internal/auth/my_auth.go`:

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

## 2. Интерфейс Authenticator

Каждый клиент аутентификации должен реализовать:

```go
type Authenticator interface {
    New() error                    // Инициализация, разрешение env-переменных
    Type() Type                    // Возврат идентификатора типа аутентификации
    Apply(req *http.Request, out *Info) error  // Применение аутентификации к запросу
    Validate() error               // Валидация обязательных полей
}
```

## 3. Константа типа

Добавьте в `internal/auth/auth.go`:

```go
const MyAuth Type = "my-auth"
```

## 4. YAML-декодер

Добавьте функцию декодера в `internal/config/auth.go`. Декодер получает `*yaml.Node` и должен декодировать его в структуру вашего клиента аутентификации:

```go
func decodeMyAuth(node *yaml.Node) (auth.Authenticator, error) {
    var client auth.MyAuthClient
    if err := decodeConfig(node, &client); err != nil {
        return nil, err
    }
    return &client, nil
}
```

Хелпер `decodeConfig` обрабатывает общий шаблон: проверяет, что узел не пуст, декодирует YAML в структуру и возвращает описательную ошибку при неудаче.

## 5. Регистрация декодера

Добавьте ваш декодер в карту `authDecoders` в `internal/config/auth.go`:

```go
var authDecoders = map[string]authDecoder{
    // ... существующие декодеры
    auth.MyAuth.String(): decodeMyAuth,
}
```

Метод `UnmarshalYAML` на `Auth` читает поле `type` из YAML, нормализует подчёркивания в дефисы, ищет декодер в `authDecoders` и вызывает его с узлом `config`. Так swag2mcp узнаёт, какой клиент аутентификации создать для каждой спецификации.

## 6. Тесты

Создайте `internal/auth/my_auth_test.go` с table-driven тестами, покрывающими:

- `New()` правильно разрешает env-переменные
- `Type()` возвращает правильный тип
- `Apply()` устанавливает правильные заголовки/query-параметры
- `Apply()` корректно обрабатывает пустые значения
- `Validate()` проходит для валидной конфигурации
- `Validate()` возвращает ошибку для отсутствующих обязательных полей
