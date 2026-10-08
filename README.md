# Todo App

Менеджер задач: REST API на Go + PostgreSQL и веб-интерфейс.

- Пользователи и их задачи (CRUD)
- Статистика: процент выполнения и среднее время выполнения, с фильтром по автору и периоду
- Фронтенд — один файл `public/index.html`, без сборки

## Запуск

Нужны Go 1.27+, Docker и make.

```sh
cp .env.example .env     # заполнить POSTGRES_* и TIME_ZONE

make env-up              # PostgreSQL
make migrate-up          # миграции
make env-port-forward    # проброс порта 5432 для локального запуска
make todoapp-run         # приложение на :5050
```

Или всё в Docker: `make todoapp-deploy` — поднимет PostgreSQL, применит миграции и запустит приложение.

## Адреса

| | |
|---|---|
| Веб-интерфейс | http://localhost:5050/ |
| API | http://localhost:5050/api/v1 |
| Swagger | http://localhost:5050/swagger/index.html (если `HTTP_SWAGGER_ENABLED=true`) |
| Health check | http://localhost:5050/health |

Обновить Swagger-спецификацию: `make swagger-gen`.
