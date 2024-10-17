# sms-service

Реализация сервиса, который эмулирует работу API сервиса по выдаче СМС с помощью API методов.
Описание протокола: https://7grizzlysms.com/partner-documentation

Методы:

- `GET_SERVICES` - отдача списка сгенерированных номеров
- `GET_NUMBER` - запрос номера по условиям (отдача случайного из сгенерированных)
- `FINISH_ACTIVATION` - завершение работы с номером
- `PUSH_SMS` - отправка сгенерированной смс

## Таблицы в БД (sqlite)

- `страны`
- `сервисы`
- `номера`
- `активации`
- `смс`

Страны и сервисы заранее определены, тестовые номера генерируются.

## Запуск

Перед запуском нужно скопировать `.env_example` в `.env` и указать значения переменных, в случае необходимости.
На данный момент используется только одна переменная `DATABASE_DSN`, не требующая изменений.

```shell
cp .env_example .env
# edit .env if needed
```

Запуск сервиса в Docker контейнере.

```shell
docker-compose up -d --build
```

Запуск тестов:

```shell
go test ./...
```

## Для разработки

Создание миграции golang-migrate:

```shell
migrate create -ext sql -dir internal/infrastructure/database/migrations -seq create_table_name
```

