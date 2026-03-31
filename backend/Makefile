container_runtime := $(shell which podman || which docker)

$(info using ${container_runtime})

up:
	${container_runtime} compose up --build -d

down:
	${container_runtime} compose down

clean:
	${container_runtime} compose down -v

lint:
	golangci-lint run ./...

test:
	go test ./... -v

test-cover:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

build:
	go build -ldflags="-s -w" -o bin/logistics-service ./logistics-service/main.go

# Генерация Go кода из .proto файлов.
# Требует: protoc, protoc-gen-go, protoc-gen-go-grpc
# Установить плагины: make proto-install
proto:
	protoc \
		--proto_path=../proto \
		--go_out=. \
		--go_opt=module=logistics-service \
		--go-grpc_out=. \
		--go-grpc_opt=module=logistics-service \
		../proto/logistics/logistics.proto

proto-install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

tool:
	go get github.com/ilyakaznacheev/cleanenv
	go get github.com/golang-jwt/jwt/v5
	go get github.com/golang-migrate/migrate/v4
	go get github.com/jmoiron/sqlx
	go get github.com/go-playground/validator/v10
	go get github.com/DATA-DOG/go-sqlmock
	go get github.com/jackc/pgconn
	go get github.com/stretchr/testify
	go get google.golang.org/grpc
	go get google.golang.org/protobuf