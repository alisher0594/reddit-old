
# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

compose-up: ### Run docker-compose
	docker-compose up --build -d postgres && docker-compose logs -f
.PHONY: compose-up

compose-up-integration-test: ### Run docker-compose with integration test
	docker-compose up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

test: ### run test
	go test -v -cover -race ./...
.PHONY: test

run: ### run app
	go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate ./cmd/api
.PHONY: run

migrate-create:  ### create new migration
	migrate create -seq -ext=.sql -dir=./migrations 'migrate_name'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path=./migrations -database=$(PG_URL) up
.PHONY: migrate-up

migrate-down: ### migration down
	migrate -path=./migrations -database=$(PG_URL) down 1
.PHONY: migrate-down
