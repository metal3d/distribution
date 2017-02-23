SOURCES := $(shell find . -name '*.go')
BIN := dist-example

run: build
	docker-compose up

build: $(BIN)
	docker-compose build

$(BIN): $(SOURCES)
	CGO_ENABLED=0 go build -tags netgo -o $@
	strip $@

clean:
	docker-compose down -v

clean-all: clean
	rm $(BIN)
	
