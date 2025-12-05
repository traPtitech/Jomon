up:
	@docker-compose up -d --build

.PHONY: down
down:
	@docker-compose down -v

.PHONY: server-test
server-test:
	docker-compose -f server-test.yml run --rm jomon-server

.PHONY: client
client:
	cd client; npm run lint; npm run dev

.PHONY: mock
mock:
	docker-compose -f mock.yml run --rm jomon-mock