.PHONY: build clean run

BINARY_NAME=k8s-installer

build:
	go build -o $(BINARY_NAME) cmd/installer/main.go

clean:
	rm -f $(BINARY_NAME)

run:
	go run cmd/installer/main.go
