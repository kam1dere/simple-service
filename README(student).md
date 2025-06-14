## **Описание проекта**

Simple Service - простой REST API сервис на GO для создания задач с использованием Fiber и PostgreSQL

Реализовано:
- Загрузка конфигурации 
- Логирование через zap
- Подключение к PostgreSQL с пулом соединений (pgxpool)
- Слой репозитория для работы с базой данных
- Слой бизнес-логики (сервис) с валидацией входящих данных
- REST API с аутентификацией и CORS через middleware Fiber
- Корректное завершение работы по системным сигналам

## **Стек технологий**

- Go 1.23.6
- Fiber (веб-фреймворк)
- PostgreSQL (БД)
- pgxpool (пул соединений PostgreSQL)
- zap (логирование)
- go-playground/validator (валидация данных)
- godotenv и envconfig (загрузка конфигурации из .env и переменных окружения)

## **Установка и запуск**

1. Клонировать репозиторий:
```
git clone https://github.com/yourusername/simple-service.git
cd simple-service
```
2. Создать и заполнить .env файл (пример в config/.env.example):
````
POSTGRES_USER=youruser
POSTGRES_PASSWORD=yourpassword
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=yourdb
REST_LISTEN_ADDRESS=:8080
REST_TOKEN=your_api_token
LOG_LEVEL=info
````
3. Собрать и запустить приложение:
``````
go build -o simple-service
./simple-service
``````
## **Тестирование API**

- Создание задачи

Запрос:
``````
POST http://localhost:8080/v1/create_task
Content-Type: application/json
Authorization: Bearer your_secret_token
``````
``````
{
"title": "New Feature",
"description": "Develop new API endpoint"
}
``````
Ответ:
``````
{
"status": "success",
"data": {
"task_id": 1
}
}
``````
- Получение списка задач

Запрос:
```
GET http://localhost:8080/v1/tasks
Authorization: Bearer your_secret_token
```
Пример запроса через curl:
```
curl -H "Authorization: Bearer your_secret_token" http://localhost:8080/v1/tasks
```
Ответ:
```
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "title": "New Feature",
      "description": "Develop new API endpoint",
      "status": "pending",
      "created_at": "2025-06-13T14:00:00Z"
    },
    {
      "id": 2,
      "title": "Fix Bug",
      "description": "Fix issue with login",
      "status": "completed",
      "created_at": "2025-06-12T10:30:00Z"
    }
  ]
}
```
- Получение одной задачи 

Запрос:
```
http://localhost:8080/v1/task/1
Authorization: Bearer your_secret_token
```
Пример запроса через curl:
```
curl -H "Authorization: Bearer your_secret_token" http://localhost:8080/v1/task/1
```
Ответ:
```
{
    "data": {
        "ID": 1,
        "Title": "Test task",
        "Description": "Test description",
        "Status": "new",
        "CreatedAt": "2025-06-13T12:51:15.868139Z",
        "UpdatedAt": "2025-06-13T14:18:49.60839Z"
    },
    "status": "success"
}
```


