import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'swag2mcp',
  description: 'Bridging OpenAPI/Swagger/Postman specifications with LLM agents via MCP',
  lang: 'en-US',
  lastUpdated: true,
  cleanUrls: true,
  base: '/',

  head: [
    ['link', { rel: 'icon', href: '/logo.svg' }],
  ],

  themeConfig: {
    logo: '/logo.svg',
    siteTitle: 'swag2mcp',

    nav: [
      { text: 'Home', link: '/' },
      { text: 'Getting Started', link: '/getting-started/installation' },
      { text: 'GitHub', link: 'https://github.com/mmadfox/swag2mcp' },
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Installation', link: '/getting-started/installation' },
                { text: 'Quick Start', link: '/getting-started/quickstart' },
        ],
      },
      {
        text: 'Concepts',
        items: [
          { text: 'Overview', link: '/concepts/overview' },
          { text: 'Specs', link: '/concepts/specs' },
          { text: 'Collections', link: '/concepts/collections' },
          { text: 'Tags', link: '/concepts/tags' },
          { text: 'Endpoints', link: '/concepts/endpoints' },
          { text: 'Workspace', link: '/concepts/workspace' },
        ],
      },
      {
        text: 'Configuration',
        items: [
          { text: 'Config File', link: '/configuration/config-file' },
          { text: 'Global Settings', link: '/configuration/global-settings' },
          { text: 'Spec Settings', link: '/configuration/spec-settings' },
          { text: 'Collection Settings', link: '/configuration/collection-settings' },
          { text: 'HTTP Client', link: '/configuration/http-client' },
          { text: 'MCP Server', link: '/configuration/mcp-server' },
          { text: 'Cascade', link: '/configuration/cascade' },
        ],
      },
      {
        text: 'CLI',
        items: [
          { text: 'Overview', link: '/cli/overview' },
          { text: 'init', link: '/cli/init' },
          { text: 'add', link: '/cli/add' },
          { text: 'delete', link: '/cli/delete' },
          { text: 'ls', link: '/cli/ls' },
          { text: 'run', link: '/cli/run' },
          { text: 'validate', link: '/cli/validate' },
          { text: 'clean', link: '/cli/clean' },
          { text: 'update', link: '/cli/update' },
          { text: 'mcp', link: '/cli/mcp' },
          { text: 'version', link: '/cli/version' },
          { text: 'info', link: '/cli/info' },
          { text: 'import', link: '/cli/import' },
          { text: 'export', link: '/cli/export' },
        ],
      },
      {
        text: 'MCP Tools',
        items: [
          { text: 'Overview', link: '/mcp-tools/overview' },
          { text: 'Discovery', link: '/mcp-tools/discovery' },
          { text: 'Endpoints', link: '/mcp-tools/endpoints' },
          { text: 'Execution', link: '/mcp-tools/execution' },
          { text: 'Utilities', link: '/mcp-tools/utilities' },
          { text: 'Skills', link: '/mcp-tools/skills' },
        ],
      },
      {
        text: 'Authentication',
        items: [
          { text: 'Overview', link: '/auth/overview' },
          { text: 'None', link: '/auth/none' },
          { text: 'Basic', link: '/auth/basic' },
          { text: 'Bearer', link: '/auth/bearer' },
          { text: 'API Key', link: '/auth/api-key' },
          { text: 'Digest', link: '/auth/digest' },
          { text: 'HMAC', link: '/auth/hmac' },
          { text: 'OAuth2 CC', link: '/auth/oauth2-cc' },
          { text: 'OAuth2 Password', link: '/auth/oauth2-pwd' },
          { text: 'Script', link: '/auth/script' },
        ],
      },
      {
        text: 'Advanced',
        items: [
          { text: 'Search', link: '/advanced/search' },
          { text: 'Rate Limiting', link: '/advanced/ratelimit' },
          { text: 'Response Size', link: '/advanced/response-size' },
          { text: 'Caching', link: '/advanced/caching' },
          { text: 'Mock Server', link: '/advanced/mock-server' },
          { text: 'TUI', link: '/advanced/tui' },
          { text: 'Export/Import', link: '/advanced/export-import' },
          { text: 'Environment Variables', link: '/advanced/env-vars' },
        ],
      },
      {
        text: 'Integration',
        items: [
          { text: 'OpenCode', link: '/integration/opencode' },
          { text: 'Cursor', link: '/integration/cursor' },
          { text: 'Claude Desktop', link: '/integration/claude' },
          { text: 'VS Code', link: '/integration/vscode' },
          { text: 'Crush', link: '/integration/crush' },
        ],
      },
      {
        text: 'Configuration Examples',
        items: [
          { text: 'CLI Workflow', link: '/examples/cli-workflow' },
          { text: 'LLM Session', link: '/examples/llm-session' },
        ],
      },
      {
        text: 'Development',
        items: [
          { text: 'Overview', link: '/development/overview' },
          { text: 'Project Structure', link: '/development/project-structure' },
          { text: 'Building', link: '/development/building' },
          { text: 'Testing', link: '/development/testing' },
          { text: 'Conventions', link: '/development/conventions' },
          { text: 'New Auth Method', link: '/development/new-auth' },
          { text: 'New MCP Tool', link: '/development/new-tool' },
        ],
      },
      { text: 'FAQ', link: '/faq' },
      { text: 'Troubleshooting', link: '/troubleshooting' },
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/mmadfox/swag2mcp' },
    ],

    editLink: {
      pattern: 'https://github.com/mmadfox/swag2mcp/edit/main/docs/:path',
    },

    search: {
      provider: 'local',
    },

    i18n: [
      {
        locale: 'ru',
        text: 'Русский',
        lang: 'ru-RU',
        title: 'swag2mcp',
        description: 'Объединение OpenAPI/Swagger/Postman спецификаций с LLM-агентами через MCP',
        themeConfig: {
          nav: [
            { text: 'Главная', link: '/ru/' },
            { text: 'Начало работы', link: '/ru/getting-started/installation' },
            { text: 'GitHub', link: 'https://github.com/mmadfox/swag2mcp' },
          ],
          sidebar: [
            {
              text: 'Начало работы',
              items: [
                { text: 'Установка', link: '/ru/getting-started/installation' },
                { text: 'Быстрый старт', link: '/ru/getting-started/quickstart' },
              ],
            },
            {
              text: 'Концепции',
              items: [
                { text: 'Обзор', link: '/ru/concepts/overview' },
                { text: 'Спецификации', link: '/ru/concepts/specs' },
                { text: 'Коллекции', link: '/ru/concepts/collections' },
                { text: 'Теги', link: '/ru/concepts/tags' },
                { text: 'Эндпоинты', link: '/ru/concepts/endpoints' },
                { text: 'Рабочая область', link: '/ru/concepts/workspace' },
              ],
            },
            {
              text: 'Конфигурация',
              items: [
                { text: 'Файл конфигурации', link: '/ru/configuration/config-file' },
                { text: 'Глобальные настройки', link: '/ru/configuration/global-settings' },
                { text: 'Настройки спецификации', link: '/ru/configuration/spec-settings' },
                { text: 'Настройки коллекции', link: '/ru/configuration/collection-settings' },
                { text: 'HTTP клиент', link: '/ru/configuration/http-client' },
                { text: 'MCP сервер', link: '/ru/configuration/mcp-server' },
                { text: 'Каскад', link: '/ru/configuration/cascade' },
              ],
            },
            {
              text: 'CLI',
              items: [
                { text: 'Обзор', link: '/ru/cli/overview' },
                { text: 'init', link: '/ru/cli/init' },
                { text: 'add', link: '/ru/cli/add' },
                { text: 'delete', link: '/ru/cli/delete' },
                { text: 'ls', link: '/ru/cli/ls' },
                { text: 'run', link: '/ru/cli/run' },
                { text: 'validate', link: '/ru/cli/validate' },
                { text: 'clean', link: '/ru/cli/clean' },
                { text: 'update', link: '/ru/cli/update' },
                { text: 'mcp', link: '/ru/cli/mcp' },
                { text: 'version', link: '/ru/cli/version' },
                { text: 'info', link: '/ru/cli/info' },
                { text: 'import', link: '/ru/cli/import' },
                { text: 'export', link: '/ru/cli/export' },
              ],
            },
            {
              text: 'MCP Инструменты',
              items: [
                { text: 'Обзор', link: '/ru/mcp-tools/overview' },
                { text: 'Обнаружение', link: '/ru/mcp-tools/discovery' },
                { text: 'Эндпоинты', link: '/ru/mcp-tools/endpoints' },
                { text: 'Выполнение', link: '/ru/mcp-tools/execution' },
                { text: 'Утилиты', link: '/ru/mcp-tools/utilities' },
                { text: 'Скиллы', link: '/ru/mcp-tools/skills' },
              ],
            },
            {
              text: 'Аутентификация',
              items: [
                { text: 'Обзор', link: '/ru/auth/overview' },
                { text: 'None', link: '/ru/auth/none' },
                { text: 'Basic', link: '/ru/auth/basic' },
                { text: 'Bearer', link: '/ru/auth/bearer' },
                { text: 'API Key', link: '/ru/auth/api-key' },
                { text: 'Digest', link: '/ru/auth/digest' },
                { text: 'HMAC', link: '/ru/auth/hmac' },
                { text: 'OAuth2 CC', link: '/ru/auth/oauth2-cc' },
                { text: 'OAuth2 Password', link: '/ru/auth/oauth2-pwd' },
                { text: 'Script', link: '/ru/auth/script' },
              ],
            },
            {
              text: 'Продвинутое',
              items: [
                { text: 'Поиск', link: '/ru/advanced/search' },
                { text: 'Rate Limiting', link: '/ru/advanced/ratelimit' },
                { text: 'Размер ответов', link: '/ru/advanced/response-size' },
                { text: 'Кэширование', link: '/ru/advanced/caching' },
                { text: 'Mock сервер', link: '/ru/advanced/mock-server' },
                { text: 'TUI', link: '/ru/advanced/tui' },
                { text: 'Экспорт/Импорт', link: '/ru/advanced/export-import' },
                { text: 'Переменные окружения', link: '/ru/advanced/env-vars' },
              ],
            },
            {
              text: 'Интеграция',
              items: [
                { text: 'OpenCode', link: '/ru/integration/opencode' },
                { text: 'Cursor', link: '/ru/integration/cursor' },
                { text: 'Claude Desktop', link: '/ru/integration/claude' },
                { text: 'VS Code', link: '/ru/integration/vscode' },
                { text: 'Crush', link: '/ru/integration/crush' },
              ],
            },
            {
              text: 'Примеры конфигурации',
              items: [
                { text: 'CLI Workflow', link: '/ru/examples/cli-workflow' },
                { text: 'LLM Session', link: '/ru/examples/llm-session' },
              ],
            },
            {
              text: 'Разработка',
              items: [
                { text: 'Обзор', link: '/ru/development/overview' },
                { text: 'Структура проекта', link: '/ru/development/project-structure' },
                { text: 'Сборка', link: '/ru/development/building' },
                { text: 'Тестирование', link: '/ru/development/testing' },
                { text: 'Соглашения', link: '/ru/development/conventions' },
                { text: 'Новый auth метод', link: '/ru/development/new-auth' },
                { text: 'Новый MCP инструмент', link: '/ru/development/new-tool' },
              ],
            },
            { text: 'FAQ', link: '/ru/faq' },
            { text: 'Решение проблем', link: '/ru/troubleshooting' },
          ],
        },
      },
    ],
  },
})
