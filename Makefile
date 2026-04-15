.PHONY: \
	run docker-dev-up docker-dev-down migrate-create migrate-up migrate-down db-info db-shell db-tables migrate-force
-include .env 
export

CONNECTION_STRING=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
run:
	go run ./cmd/api

docker-dev-up:
	docker compose -f docker-compose.dev.yaml up -d
	
docker-dev-down:
	docker compose -f docker-compose.dev.yaml down

migrate-create:
	migrate create -ext sql -dir internal/db/migrations -seq $(name)

migrate-up:
	migrate -path internal/db/migrations -database "$(CONNECTION_STRING)" up

migrate-down:
	migrate -path internal/db/migrations -database "$(CONNECTION_STRING)" down 1

migrate-force:
	migrate -path internal/db/migrations -database "$(CONNECTION_STRING)" force $(VERSION)

db-shell:
	docker exec -it gin-layered-architecture-db psql -U $(DB_USER) -d $(DB_NAME)

db-tables:
	docker exec -it gin-layered-architecture-db psql -U $(DB_USER) -d $(DB_NAME) -c "\dt"
	
db-info:
	docker exec -it gin-layered-architecture-db psql -U $(DB_USER) -d $(DB_NAME) -c "\d $(name)"
