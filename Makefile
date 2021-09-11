.PHONY: up
up:
	@docker-compose up -d --build

.PHONY: down
down:
	@docker-compose down -v

.PHONY: test
test:
	@mkdir -p $(CURDIR)/uplaods
	@MARIADB_HOSTNAME=localhost UPLOAD_DIR=$(CURDIR)/uplaods go test -cover -race ./...
	@rmdir $(CURDIR)/uplaods

.PHONY: docs
docs:
	@enter ./ent/schema ./docs/er.html
