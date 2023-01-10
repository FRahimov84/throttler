build:
	go build -o ./bin/service .

run:
	./bin/service run

all:	build run
