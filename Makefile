dev:
	reflex -c ./reflex.conf

build:
	go build -o gouter

up:
	go build -o bin/gouter
	./bin/gouter