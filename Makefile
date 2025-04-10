
COMPOSE = docker-compose -f docker-compose.deploy.yml

start:
	$(COMPOSE) up -d --build

stop:
	$(COMPOSE) down

test-integration:
	$(COMPOSE) -f docker-compose.test.integration.yml up --build --abort-on-container-exit --exit-code-from backend
	$(COMPOSE) -f docker-compose.test.integration.yml down -v

test-unit:
	docker-compose -f docker-compose.test.unit.yml up
	docker-compose -f docker-compose.test.unit.yml down