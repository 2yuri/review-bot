#!make
include .env
export $(shell sed 's/=.*//' .env)

MIGRATIONS_PATH="$(shell pwd)/db/migrations"
DOCKER_NETWORK="$(DB_MIGRATIONS_NETWORK)"
DB_DSN="postgres://$(DB_USER):$(DB_PASS)@db:5432/$(DB_NAME)?sslmode=disable"

db:
	docker-compose up --build -d db

migrations-up:
	docker run --rm -v $(MIGRATIONS_PATH):/migrations --network "$(DOCKER_NETWORK)" migrate/migrate -path=/migrations -database $(DB_DSN) up $(VERSION)

migrations-down:
	docker run --rm -it -v $(MIGRATIONS_PATH):/migrations  --network "$(DOCKER_NETWORK)" migrate/migrate -path=/migrations -database $(DB_DSN) down

migrations-fix:
	docker run --rm -it -v $(MIGRATIONS_PATH):/migrations --network "$(DOCKER_NETWORK)" migrate/migrate -path=/migrations -database $(DB_DSN) force $(VERSION)

migrations-create:
	docker run --rm -v $(MIGRATIONS_PATH):/migrations --user $(shell id -u):$(shell id -g) \
					--network $(DOCKER_NETWORK) migrate/migrate create -ext sql -dir ./migrations -seq $(NAME)