Сервис для управления событиями с напоминаниями и автоматической архивацией.

## Запуск

```bash
go run main.go --port 8080 --clean-interval 5m
```

## API

### Создать событие

```bash
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"date":"2099-01-01","title":"Meeting"}'
```

С напоминанием:
```bash
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"date":"2099-01-01","title":"Meeting","remind_time":"2099-01-01T09:00:00"}'
```

### Обновить событие

```bash
curl -X POST http://localhost:8080/update_event \
  -H "Content-Type: application/json" \
  -d '{"id":1,"user_id":1,"date":"2099-06-01","title":"Updated Meeting"}'
```

### Удалить событие

```bash
curl -X POST http://localhost:8080/delete_event \
  -H "Content-Type: application/json" \
  -d '{"id":1,"user_id":1}'
```

### Получить события

За день
```bash
curl "http://localhost:8080/events_for_day?user_id=1&date=2099-01-01"
```
За неделю
```
curl "http://localhost:8080/events_for_week?user_id=1&date=2099-01-01"
```
За месяц

```
curl "http://localhost:8080/events_for_month?user_id=1&date=2099-01-01"
```