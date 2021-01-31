build:
	mkdir build
	go build -o build/dashboard

run:
	touch .env
	go run main.go