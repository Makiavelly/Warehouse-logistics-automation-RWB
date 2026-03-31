container_runtime := $(shell which podman || which docker)

$(info using ${container_runtime})

up:
	${container_runtime} compose up --build -d

down:
	${container_runtime} compose down

clean:
	${container_runtime} compose down -v

# Регенерация gRPC кода из .proto файла
proto:
	cd backend && make proto

proto-install:
	cd backend && make proto-install

# Инструменты разработки
tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "checking protoc, if it fails: https://protobuf.dev/installation/"
	@protoc --version