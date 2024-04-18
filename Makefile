include .env

.PHONY:

compose-up:
	docker-compose up --build -d postgres rabbitmq

up:
	docker-compose up -d

run:
	go run cmd/app/main.go

get:
	go get -d -v ./...

test:
	go test -cover ./...   

swag:
	swag init -dir internal/controller/http/v1/ -generalInfo router.go --parseDependency internal/entity/ 

lint:
	golangci-lint run

migration-new-db:
	migrate -path migrations -database '$(PG_URL)' down
	migrate -path migrations -database '$(PG_URL)' goto 20240418133357

migration-add-testdata:
	migrate -path migrations -database '$(PG_URL)' goto 20240418133358