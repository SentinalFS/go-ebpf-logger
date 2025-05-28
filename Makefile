APP_NAME=listener

.PHONY: build run clean docker

build:
	go build -o $(APP_NAME) ./cmd/listener

run:
	./$(APP_NAME) monitor.bpf.o

clean:
	rm -f $(APP_NAME)

docker:
	docker build -t go-ebp-logger -f build/docker/Dockerfile .
