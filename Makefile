.PHONY: up
up:
	@docker-compose up --build

.PHONY: down
down:
	@docker-compose down -v

.PHONY: test
test:
	docker-compose -f server-test.yml run --rm jomon-server

.PHONY: docs
docs:
	enter ./ent/schema ./docs/er.html
