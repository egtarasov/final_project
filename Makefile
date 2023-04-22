ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/app
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/database/migrations
UNIT_TEST_SERVER=$(CURDIR)/internal/app/server
.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql


.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up
	

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down	
	

.PHONY: compose-up
compose-up:
	docker-compose build
	docker-compose up -d


.PHONY: compose-down
compose-down:
	docker-compose down

.PHONY: test-unit
test-unit:
	go test -v ./internal/app/server

.PHONY: test-group
test-group:
	go test -v ./internal/app/group

.PHONY: test-student
test-student:
	go test -v ./internal/app/student

.PHONY: test-server
test-server:
	go test -v ./internal/app/server/test_integration