## Тестовое задание stat4market.

### Задача 1:

Написанные sql запросы расположены в `task_1/queries.sql`

### Задача 2:

Генератор данных расположен в `generator.go`

### Задача 3:

Сервер с API можно запустить командой 
```shell
docker compose up
```

- POST /api/event
- POST /api/events

Пример тела запроса `POST /api/events`:

```json
{
    "event_type": "some",
    "events_from": "2024-04-21 07:05:49",
    "events_to": "2024-04-24 07:05:49"
}
```