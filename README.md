# sms-service

Реализация сервиса, который эмулирует работу API сервиса по выдаче СМС с помощью API методов.
Описание протокола: https://7grizzlysms.com/partner-documentation

Методы:
- GET_SERVICES - отдача списка сгененированных номеров
- GET_NUMBER - запрос номера по условиям (отдача случайного из сгененированных)
- FINISH_ACTIVATION - завершение работы с номером 
- PUSH_SMS - отправка сгенерированной смс

Таблицы в БД (sqlite):
- страны
- сервисы
- номера
- активации
- смс

Страны и сервисы заранее определены, тестовые номера генерируются. 
Запуск сервиса в докер контейнере.
```shell
docker-compose up -d --build
```

## Для разработки
Установка golang-migrate:
```shell
go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
Создание миграции:
```shell
migrate create -ext sql -dir internal/infrastructure/database/migrations -seq create_tablename_table
```