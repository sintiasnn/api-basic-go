APP=api-basic-go
BIN_DIR=bin
BIN=$(BIN_DIR)/$(APP)

.PHONY: build run clean docker-build docker-run compose-up compose-down

build:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -o $(BIN) .

run:
	go run .

clean:
	rm -rf $(BIN_DIR)

docker-build:
	docker build -t $(APP):latest .

docker-run:
	docker run --rm -p 8080:8080 -e PORT=8080 -e CORS_ALLOWED_ORIGINS=* $(APP):latest

compose-up:
	docker compose up --build -d

compose-down:
	docker compose down

