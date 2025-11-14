# Users Service

Простой HTTP-сервис на Go для управления пользователями, выполнения задач, начисления очков и работы с реферальной системой. Проект построен в стиле **гексагональной архитектуры (Hexagonal / Ports & Adapters)** и включает авторизацию через JWT, PostgreSQL, миграции и сборку через Docker Compose.

---

## Возможности сервиса

- Регистрация и аутентификация пользователей  
- Авторизация через JWT Access Token (middleware)  
- Получение статуса пользователя  
- Лидерборд по пользователям (по количеству очков)  
- Выполнение заданий (начисление награды)  
- Ввод реферального кода  
- Хранение данных в PostgreSQL  
- Миграции через **golang-migrate**  
- Поддержка Docker Compose  

---

## Структура проекта

```
cmd/users-service/
  main.go                    — точка входа
internal/app/
  config/                    — загрузка конфигурации
  logger/                    — настройка slog
  domain/                    — доменные сущности + порты
    user_domain/
  usecase/                   — бизнес-логика
    user_usecase/
    auth_usecase/
    jwt_usecase/
  adapters/
    http/                    — HTTP хэндлеры, роутеры, middleware JWT + logger
      handlers/
      middlewares/
      utils/
    storage/
      postgres/              — реализация репозиториев для PostgreSQL
migrations/                  — миграции golang-migrate (.sql)
Dockerfile
docker-compose.yml
```

## API

### POST `/register`

Регистрация нового пользователя

**Пример тела (JSON):**

```json
{
  "email": "user@example.org",
  "password": password // от 8 до 100 символов
}
```

**Успешный ответ:** `"status": "success"`

**Ошибки:**

* `400` — некорректный входной JSON / валидация

### POST `/login`

Позволяет залогиниться зарегистрированному пользователю

**Пример тела (JSON):**

```json
{
  "email": "user@example.org",
  "password": password
}
```

**Успешный ответ:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiODBmNzg5ZTAtMWZiNC00NjM2LWIzY2QtOWRjY2I3YzE1NTU2IiwiZXhwIjoxNzYzMTE4MjUwLCJpYXQiOjE3NjMxMTczNTB9.zsY8uo1cPJchQzFVKcWi2rACCaGgpoBTTHSG-egK4Ls",
  "refresh_token": "4883db53b30bc1ab454a609ef3c3490298b0e505052107c4c0fec3c6654e6a5b"
}
```

**Ошибки:**

* `400` — некорректный входной JSON / валидация

### GET `users/{id}/status`

Получить инфомрацию о пользователе по его id, который указан в URL запроса

**Успешный ответ:**
```json
{
  "status": "success",
  "user_id": "86313830-32b8-4023-bec0-3f314c983376",
  "score": 150
}
```

**Ошибки:**

* `400` — некорректный входной JSON / валидация

### GET `/users/leaderboard`

Получить список из топ-10 пользователях с наибольшим количеством очков

**Успешный ответ:**
```json
{
    "1": {
        "score": 200,
        "user_id": "86313830-32b8-4023-bec0-3f314c983376"
    },
    "2": {
        "score": 100,
        "user_id": "6677f49f-fb62-4cee-85f4-a0851977321e"
    }
}
```

**Ошибки:**

* `400` — некорректный входной JSON / валидация

### UPDATE `users/{id}/task/complete`

Изменяет значение `complete` в таблице `users_tasks` на true и увеличивает `score` пользователя на указанный `reward`

**Успешный ответ:**  `"status": "success"`

**Ошибки:**

* `400` — некорректный входной JSON / валидация

### UPDATE `users/{id}/referrer`

Позволяет пользователю `{id}` указать `referrer_id` другого пользователя, после чего оба пользователя получат раннее установленный `reward`

**Успешный ответ:**  `"status": "success"`

**Ошибки:**

* `400` — некорректный входной JSON / валидация

Для проверки эндпоинтов использовалась утилита httpie:
```
http POST http://localhost:8080/register email=user@example.com password=1234
http POST http://localhost:8080/login email=user@example.com password=1234
http GET http://localhost:8080/users/{id}/status Authorization:"Bearer $ACCESS_TOKEN"
http PATCH http://localhost:8080/users/{id}/task/complete Authorization:"Bearer $ACCESS_TOKEN" task="task_name"
http PATCH http://localhost:8080/users/{id}/referrer Authorization:"Bearer $ACCESS_TOKEN" referrer_id={id} task="task_name"
http GET http://localhost:8080:8080/users/leaderboard Authorization:"Bearer $ACCESS_TOKEN"
```

