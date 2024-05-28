start:
	go run github.com/cosmtrek/air	

dev: generate
	go run .

generate:
	go generate ./...

build: generate
	mkdir -p build && go build -o build/coma

install: generate
	go install

clean:
	@find . -name **fakes -delete
	@rm -rf coma/coma.cfg docs/docs.go docs/swagger.json docs/swagger.yaml docs


# Docker
docker-start: docker-build docker-run
	docker start coma

docker-stop:
	docker stop coma

docker-build:
	docker build -f docker/Dockerfile -t coma .

docker-run:
	docker run -d --name coma -p 5899:5899 coma

docker-clean:
	docker rm coma && docker rmi coma

# End Docker