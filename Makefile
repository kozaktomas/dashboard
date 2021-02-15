.PHONY: build install run

build:
	go build -o build/ddboard

install:
	cp build/ddboard /usr/local/bin/ddboard

run:
	touch .env
	go run main.go