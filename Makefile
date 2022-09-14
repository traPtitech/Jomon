.PHONY: up
up:
	-@docker-compose up -d --build

.PHONY: dev-up
dev-up:
	-@docker-compose up --build

.PHONY: down
down:
	@docker-compose down -v

.PHONY: test
test:
	@mkdir -p $(CURDIR)/uploads
	@MARIADB_HOSTNAME=localhost UPLOAD_DIR=$(CURDIR)/uploads go test -cover -race ./...

.PHONY: docs
docs:
	@enter ./ent/schema
	@cp er.html ./docs/
	@rm er.html
