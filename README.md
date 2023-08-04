# Start app

1. docker-compose up
2. Enjoy ;)

# krakenD

## If you exec app services not in docker for mac:

1. In ./krakend/config/krakend.json change backend.host to `http://docker.for.mac.localhost:portServices` if you use mac.
2. Start services as `go run ./cmd/auth/main.go -l "Debug" -a "127.0.0.1:8081"`

# for build apps:

1. run `make build` in root directory

# swag

1. go to root folder project
2. exec `~/go/bin/swag init -g "cmd/auth/main.go" -o "docs/auth/"`
3. swagger on page `/swagger/index.html`

# create new migration

1. go to /db/{service_name}.
2. In `docker-compose.migrations.yml` change service `migration_db_add`.
3. Go to /Makefiles/{service_name} and call `make migration_new `.

# Отличия реализации от задания

Все контейнеры скрыты от пользователя в докеркомпоузе. Этим самым соблюдён один из принципов REST Многоуровневость,
в коммуникации участвуют двое: клиент и сервер. Каждый компонент должен видеть только свой уровень. Клиент общается 
только с auth (регистрация и аутентификация) и gophermart (логика приложения доступные для клиента). Для доступа к 
ручкам gophermart необходимо пройти аутентификацию. 

Сервиса accrual является внутренним для системы. Ручки сервиса accrual не проверяют аутентификацию клиента и поэтому 
не доступны для пользователя. При обращении к ним пользователя они должны быть недоступны и такой результат 
интеграционных тестов считается ожидаемым.

Для тестирования сервиса accrual его можно поднять отдельно. Например, так.

```bash
go run cmd/accrual/main.go -a "localhost:8081" -d "postgres://postgres:postgres@localhost:5432/accrual?sslmode=disable"
```

Для проверки результатов работы сервиса есть своеобразный backdoor. Наружу прокинут порт 5432 и в контейнере поднят 
pgAdmin можно посмотреть данные в базе.

В репозитории есть коллекция которую рекомендуем использовать для интеграционного тестирования. Она очень близка к 
техническому заданию, покрывая 100% реализованных ручек со всеми вариантами ответов предусмотренными нашим сервисом.
Кроме этого код на 70% покрыт unit тестами на хендлеры.

