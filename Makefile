.PHONY: up
up:
	@docker-compose up --build -d

.PHONY: dev-up
dev-up:
	@docker-compose -f docker-compose.dev.yml up --build

.PHONY: down
down:
	@docker-compose down -v

.PHONY: test
test:
	docker-compose -f server-test.yml run --rm jomon-server

.PHONY: docs
docs:
	enter ./ent/schema ./docs/er.html
