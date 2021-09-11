.PHONY: up
up:
	@docker-compose up -d --build

.PHONY: down
down:
	@docker-compose down -v

.PHONY: test
test:
	@MARIADB_HOSTNAME=localhost go test -cover -race ./...

.PHONY: docs
docs:
	@enter ./ent/schema ./docs/er.html
