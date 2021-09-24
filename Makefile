.PHONY: up
up:
	@docker-compose up -d --build

.PHONY: down
down:
	@docker-compose down -v

.PHONY: test
test:
	@mkdir -p $(CURDIR)/uploads
	@MARIADB_HOSTNAME=localhost UPLOAD_DIR=$(CURDIR)/uploads go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@rmdir $(CURDIR)/uploads

.PHONY: docs
docs:
	@enter ./ent/schema ./docs/er.html
