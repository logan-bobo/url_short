fmt:
		go fmt ./...
.PHONY:fmt

lint: fmt
		golangci-lint run -v
.PHONY:lint

vet: lint
		go vet ./...
.PHONY:vet

build:
		docker build . -t "url-short:latest"
.PHONY:build

run:
		docker compose up -d
.PHONY:run

stop:
		docker compose down
.PHONY:stop

test:
		docker build . -t "url-short:test" --target tester
		docker run -t "url-short:test"
.PHONY:test

migrate: 
		goose -dir sql/schema postgres "postgres://url_short:password@localhost:5002/url_short" up 
.PHONY:migrate
