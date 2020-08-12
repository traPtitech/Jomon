.PHONY: server-test
server-test:
	docker-compose -f server-test.yml run --rm jomon-server

.PHONY: client
client:
	docker-compose -f mock-for-client.yml up -d
	cd client; npm run lint; npm run serve