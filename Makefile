# Targets for Makefile

# Command to run PostgreSQL using Docker Compose
run-db:
	docker-compose up -d db redis

# Command to build and run the Go app using Docker Compose
run:
	docker-compose up -d app --build

local-run:
	go run ./cmd/main.go

# Command to stop all services
stop:
	docker-compose down

# Command to view logs for all services
logs:
	docker-compose logs -f

SWAGGER := $(shell which swag)
SWAGGER_OUT_DIR := docs
SWAGGER_GEN_SCRIPT := $(SWAGGER) init -g ./internal/http/router.go -o $(SWAGGER_OUT_DIR) --parseDependency --parseInternal --parseDepth 1

swag-gen:
	$(SWAGGER_GEN_SCRIPT)

gen-proto:
	./scripts/genProto.sh .


