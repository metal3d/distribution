
build:
	CGO_ENABLED=0 go build -tags netgo -o dist-example
	docker-compose build

run: build
	docker-compose up
