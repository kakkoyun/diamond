all: build

.PHONY: build
build: deps
	go build .

.PHONY: deps
deps:
	go mod tidy

.PHONY: dev/up
dev/up: manifest.yaml
	source ./local-dev.sh && up

.PHONY: dev/down
dev/down:
	source ./local-dev.sh && down

.PHONY: container
container:
	docker build -t kakkoyun/diamond:dev .
