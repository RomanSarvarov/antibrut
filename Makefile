go_build_flags=-tags=sqlite_unlock_notify

init:
	cp -n .env.example .env || true

build:
	docker-compose build

run: init
	docker-compose up -d

stop:
	docker-compose stop

restart:
	docker-compose restart

tool:
	docker-compose exec antibrut ./antibrut $(MAKECMDGOALS)

compile:
	go build ${go_build_flags} -o ./bin/antibrut ./cmd/antibrut

generate:
	go generate ./...

lint: go-lint proto-lint

go-lint:
	golangci-lint run ./...

proto-lint:
	cd proto && go run github.com/bufbuild/buf/cmd/buf lint

test:
	go test -race -count 1 ./...

supertest:
	go test -race -count 10 ./...

integration-test:
	go test -race -tags integration -count 1 ./...

tests: supertest integration-test

.PHONY: init run stop restart tool build generate compile lint go-lint proto-lint test supertest integration-test tests testshort
