# Start app

1. make build_docker
2. docker-compose up
3. Enjoy ;)

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

# go-musthave-group-diploma-tpl

Шаблон репозитория для группового дипломного проекта курса "Go-разработчик"

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без
   префикса `https://`) для создания модуля

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-group-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.
