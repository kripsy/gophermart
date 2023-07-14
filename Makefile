.PHONY: build_docker
build_docker: build_docker_auth build_docker_accrual build_docker_gophermart


.PHONY: build_docker_auth
build_docker_auth:
	docker build --tag docker-auth -f ./Dockerfiles/auth/Dockerfile .

.PHONY: build_docker_accrual
build_docker_accrual:
	docker build --tag docker-accrual -f ./Dockerfiles/accrual/Dockerfile .

.PHONY: build_docker_gophermart
build_docker_gophermart:
	docker build --tag docker-gophermart -f ./Dockerfiles/gophermart/Dockerfile .




.PHONY: build
build: build_auth build_accrual build_gophermart


.PHONY: build_auth
build_auth:
	go build -o ./bin/auth ./cmd/auth/main.go

.PHONY: build_accrual
build_accrual:
	go build -o ./bin/accrual ./cmd/accrual/main.go

.PHONY: build_gophermart
build_gophermart:
	go build -o ./bin/gophermart ./cmd/gophermart/main.go