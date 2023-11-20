include .env

export POSTGRES_USER
export POSTGRES_PASSWORD
export POSTGRES_DB
export POSTGRES_PORT

DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

initdb:
	migrate -path migration -database "$(DB_URL)" -verbose up

dropdb:
	migrate -path migration -database "$(DB_URL)" -verbose down 1

.PHONY: initdb dropdb