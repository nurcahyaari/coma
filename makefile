start:
	go run github.com/cosmtrek/air

dev: generate
	go run github.com/cosmtrek/air

generate:
	go generate ./...

build: generate
	mkdir -p build && go build -o build/app
