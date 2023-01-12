include app.env
export

swag-v1: ### swag init
	swag init -g internal/controller/http/v1/router.go

build:
	go build -o ./bin/service .

run:
	./bin/service run

all: build run
