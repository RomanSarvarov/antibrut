go_build_flags=-tags=sqlite_unlock_notify

init:
	cp -n .env.example .env || true

run: init
	docker-compose up -d --build

stop:
	docker-compose stop

tool:
	docker-compose exec antibrut ./antibrut $(MAKECMDGOALS)

build:
	go build ${go_build_flags} -o ./bin/antibrut ./cmd/antibrut

generate:
	go generate ./...

lint: go-lint proto-lint

go-lint:
	golangci-lint run ./...

proto-lint:
	cd proto && go run github.com/bufbuild/buf/cmd/buf lint

test:
	go test -race ./... -count 10

PHONY: init run stop tool build generate lint go-lint proto-lint test