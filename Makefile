include .env

.PHONY: start-dev
start-dev:
	@clear
	@echo "starting the server"
	watchexec -r -e go go run ./cmd/app -env=development

.PHONY: start
start:
	@clear
	@echo "starting the server"
	go run ./cmd/app -env=production

.PHONY: migration/create
migration/create:
	@echo "creating new migration files for ${name}"
	migrate create -ext sql -dir ./internal/db/migrations -seq -digits 8 ${name}

.PHONY: migration/up
migration/up:
	@echo "running migration"
	migrate -path ./internal/db/migrations -database ${DSN} up

.PHONY: migration/down
migration/down:
	@echo "running migration"
	migrate -path ./internal/db/migrations -database ${DSN} down

.PHONY: tests
tests:
	@clear
	@echo "running all tests"
	go test -v ./...

.PHONY: tests-e2e
tests-e2e:
	@clear
	@echo "running e2e tests"
	go test -v ./tests

.PHONY: test_domain
test_domain:
	@clear
	@echo "running e2e tests for ${domain}"
	go test -v ./tests/${domain}

.PHONY: run-sqlc
run-sqlc:
	sqlc generate