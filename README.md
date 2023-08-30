# Сервис динамической сегментации пользователей

[![Go Report Card](https://goreportcard.com/badge/github.com/vieux/docker-volume-sshfs)](https://goreportcard.com/report/github.com/OkDenAl/dynamic-user-segmentation)

Cервис для работы с сегментами, добавления и удаления пользователей в сегмент и последующем получении отчёта по операциям
добавления и удаления пользователя из сегмента.

Репозиторий содержит реализацию [тестового задания](https://github.com/avito-tech/backend-trainee-assignment-2023) на позицию стажера (бекенд).

#### Выполнено:
- Основное задание
- Дополнительное задание 1 (получения отчета по пользователю за определенный период)
- Дополнительное задание 2 (возможность задавать время автоматического удаления пользователя из сегмента)
- Дополнительное задание 3 (автоматическое добавление пользователей в сегмент при его создании)
- Наличие Swagger документации
- Покрытие кода тестами

#### Используемые технологии:
- [Gin](https://github.com/gin-gonic/gin) (веб фреймворк)
- pgx (драйвер для работы с PostgreSQL)
- golang/mock, testify (для тестирования)
- Google Drive API (для выгрузки отчётов на Google Drive)
- PostgreSQL (в качестве БД)
- Docker (для запуска сервиса)
- Swagger

Сервис был написан с Clean Architecture, что позволит в будущем легко наращивать функционал и расширять сервис.
Кроме того реализован Graceful Shutdown для грамотного завершения работы сервиса.

## Подготовка к запуску

Для подготовки к запуску требуется выполнить следующие шаги:
1. Скачать исходный код и перейти в директорию с проектом
    ```
    git clone https://github.com/OkDenAl/dynamic-user-segmentation.git
    cd dynamic-user-segmentation
    ```
2. Получить доступ к `Google Drive API`

   - Зарегистрировать приложение в `Google Cloud Platform`: [Документация](https://developers.google.com/workspace/guides/create-project)
   - Создать сервисный аккаунт и его секретный ключ: [Документация](https://developers.google.com/workspace/guides/create-credentials)
   - Добавить секретный ключ в директорию `config`

3. Добавить `.env` файл в директорию с проектом и заполнить его данными из `.env.example`

## Запуск
Запустить сервис можно с помощью команды `make start`

Swagger - Документацию после запуска сервиса можно посмотреть по адресу `http://localhost:8081/`

Для запуска тестов необходимо выполнить команду `make test`, для запуска тестов с покрытием `make cover` 
для получения информации о проценте покрытия

Для запуска линтера необходимо выполнить команду `make lint`

### Замечание
`При первом запуске` необходимо инициализировать SQL таблицы, код на создание таблиц находится [тут](postgresql/init.sql).
Необходимо сделать следующее:
- Подключиться к базе данных в запущенном контейнере
   ```
  docker exec -it pg psql -U raiden -d segmentation
  ```
- Вставить код из `postgresql/init.sql`

## Детали реализации
При удалении сегмента автоматически удаляются все записи об этом сегменте в таблице, 
связывающей пользователей и их сегменты

Раз в день (в 0:00) запускается скрипт с помощью `pg_cron`, который очищает записи в таблице
пользователей с их сегментами с истёкшим TTL

При добавлении или удалении пользователя из сегмента используется транзакция, чтобы
в случае ЧС информация об операции не была потеряна (т.к обновление таблицы с операциями
происходит при запросе на добавление пользователя в сегмент)

Т.к в задании не было сказано реализовывать API для работы с пользователями, то таблица с пользователями
считалась мною уже заполненной (сам добавил тестовые INSERT в таблицу users), я реализовал лишь метод репозитория
для получения количества пользователей

Все таблицы являются полностью безопасными к повторяющимся данным (на уровне таблиц в БД любое дублирование запрещено),
кроме таблицы с операциями

Было принято решение использовать Google Drive для хранения отчетов, мне показалось это интересной задачей,
также была идея реализовать какое-нибудь внутреннее хранилище

Для удобства при запросе на добавление пользователя в сегмент сегменты передаются строкой через запятую


## Пример работы API

Некоторые примеры запросов:
- [Создание нового сегмента](#create-segment)
- [Создание нового сегмента с автоматическим добавлением его пользователям](#create-auto-segment)
- [Удаление сегмента](#delete-segment)
- [Добавление пользователя в сегмент](#adding-user-to-segments)
- [Удаление пользователя из сегмента](#deleting-user-from-segments)
- [Получение активных сегментов пользователя](#get-user-segments)
- [Сводный отчёт по услугам с экспортом в Google Drive](#operations-report-link)


### Создание нового сегмента <a name="create-segment"></a>
Пример запроса:
```curl
curl --location 'http://localhost:8080/api/v1/segment/create' \
--header 'Content-Type: application/json' \
--data '
{
    "name":"test"
}'
```
Пример ответа:
```json
{
  "message": "success"
}
```

### Создание нового сегмента с автоматическим добавлением его пользователям <a name="create-auto-segment"></a>
Пример запроса:
```curl
curl --location 'http://localhost:8080/api/v1/segment/create' \
--header 'Content-Type: application/json' \
--data '
{
    "name":"test",
    "percent_of_users":50
}'
```
Пример ответа:
```json
{
  "message": "success"
}
```

### Удаление сегмента <a name="delete-segment"></a>
`Важно`: если какой-то пользователь состоит в этом сегменте, то при удалении сегмента
пользователь автоматически покидает сегмент

Пример запроса:
```curl
curl --location --request DELETE 'http://localhost:8080/api/v1/segment/delete' \
--header 'Content-Type: application/json' \
--data '
{
    "name":"test"
}'
```
Пример ответа:
```json
{
  "message": "success"
}
```

### Добавление пользователя в сегменты <a name="adding-user-to-segments"></a>
`Важно`: пользователь и сегменты уже должны быть в базе.

`Важно`: поле "expires_at" опционально и задаёт TTL для сегмента пользователя

Пример запроса:
```curl
curl --location 'http://localhost:8080/api/v1/user_segment/operation' \
--header 'Content-Type: application/json' \
--data '
{
    "user_id":1,
    "segments_to_add":"test,test1",
    "segments_to_delete":"",
    "expires_at": {
        "month":2
    }
}'
```
Пример ответа:
```json
{
  "message": "success"
}
```

### Удаление пользователя из сегментов <a name="deleting-user-from-segments"></a>
`Важно`: пользователь и сегменты уже должны быть в базе.

Пример запроса:
```curl
curl --location 'http://localhost:8080/api/v1/user_segment/operation' \
--header 'Content-Type: application/json' \
--data '
{
    "user_id":1,
    "segments_to_add":"test2",
    "segments_to_delete":"test,test1"
}'
```
Пример ответа:
```json
{
  "message": "success"
}
```

### Получение активных сегментов пользователя <a name="get-user-segments"></a>
Пример запроса:
```curl
curl --location 'http://localhost:8080/api/v1/user_segment/1'
```
Пример ответа:
```json
{
   "data": [
      {
         "name": "test1"
      },
      {
         "name": "test3"
      }
   ],
   "error": null
}
```

### Сводный отчёт по услугам с экспортом в Google Drive <a name="operations-report-link"></a>
Сервис формирует отчёт в разрезе каждой услуги, затем загружает его в Google Drive и возвращает ссылку на файл с открытым доступом на чтение

Пример запроса:
```curl
curl --location --request GET 'http://localhost:8080/api/v1/operations/report' \
--header 'Content-Type: application/json' \
--data '{
    "month":8,
    "year":2023
}'
```
Пример ответа:
```json
{
   "link": "https://drive.google.com/file/d/1mLVjeUMKAW0In8no5cSefoYe5iNZPCPL/view?usp=sharing"
}
```