# swag2mcp

<div style="background: #dc2626; color: white; padding: 20px 24px; border-radius: 12px; text-align: center; font-size: 1.4em; font-weight: 700; margin: 24px 0;">
  🚧 WORK IN PROGRESS — релиз скоро!
</div>

Объединяет спецификации OpenAPI/Swagger/Postman с LLM-агентами через протокол MCP.

<a href="https://www.youtube.com/watch?v=1Da4UmE2f9U" target="_blank">
  <img src="https://raw.githubusercontent.com/mmadfox/swag2mcp/main/docs/cover.png" alt="Превью">
</a>

## Ваше API говорит на языке LLM

Одна строка конфига превращает любой файл OpenAPI/Swagger/Postman в MCP-сервер. LLM-агенты обнаруживают, изучают и вызывают ваши API — без единой строки интеграционного кода.

<img src="/architecture.svg" width="700" alt="Архитектура swag2mcp">

## Хватит писать обёртки

Каждый раз, подключая новый API к LLM, вы пишете один и тот же шаблонный код: разбор спецификаций, аутентификация, обработка ошибок, ограничение запросов. swag2mcp делает это за вас — 19 готовых MCP-инструментов.

## Кому это нужно

| Роль | Зачем |
|------|-------|
| **Разработчик AI-агентов** | Подключить любой API за 2 минуты, а не за 2 дня |
| **MCP-инженер** | Ни строчки кода обработчика — просто укажите спецификацию |
| **Архитектор** | Единый слой интеграции API для всех LLM в компании |
| **Аналитик данных** | Доступ к API через естественный язык, без программирования |
| **DevOps / SRE** | Мониторинг и автоматизация через LLM без дополнительных сервисов |
| **Интегратор** | 9 методов аутентификации из коробки — от Basic до OAuth2 и HMAC |
| **QA-инженер** | Mock-сервер для изолированного тестирования без реальных API |
| **Продукт-менеджер** | Быстрые прототипы AI-функций без бэкенд-работы |
| **и многие другие** | |

---

## Лицензия

Распространяется под лицензией **GNU Affero General Public License v3.0** (AGPL v3).

Полный текст лицензии: [LICENSE](https://github.com/mmadfox/swag2mcp/blob/main/LICENSE).

```
SPDX-License-Identifier: AGPL-3.0-only
```
