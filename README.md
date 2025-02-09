# Сервис для мониторинга докер контейнеров

## Содержание


- [Введение](#введение)
- [Начало работы](#начало-работы)
    - [Установка](#установка)
    - [Использование](#использование)

## Введение
В рамках задания требовалось разработать сервис для получения данных о состоянии контейнеров.

1. Backend - RESTful API сервис для работы с БД:
  - реализованы эндпоинты для healthcheck работы сервиса, добавления, обновления, удаления, просмотра и получения списка всех контейнеров;
  - добавлена валидация данных по контейнерам.
2. Pinger - HTTP.client который получает данные через API, собирает необходимую информацию о контейнерах(ping) и обновляет данные на backend.
3. Frontend - реализован на go и чистом HTML/CSS. Отображает данные по всем IP адресам контейнеров.
4. БД - PostgreSQL

Использован следующий стек:
- Golang
- PostgreSQL (в качестве хранилища данных)
- golang-migrate (для миграций БД)
- pq (драйвер для работы с PostgreSQL)
- viper (для настройки конфигурации)
- pro-bing (для пинга контейнеров)
- Docker и docker compose (для запуска сервиса)

## Начало работы

### Установка

Перед началом, убедитесь, что у вас установлен `Docker`

Для запуска сервиса выполните следующие шаги:

1. Склонируйте этот репозиторий.
2. Перейдите в каталог проекта.
3. Для удобства тестирования уже добавлены данные о контейнерах из этого докер компоуза.
4. Выполните следующие команды для запуска сервиса:

    ```bash
    docker compose up -d
    ```

После запуска сервиса, информацию о контейнерах можно посмотреть, перейдя по адресу:
[`http://localhost:8080`](http://localhost:8080).
