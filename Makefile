go_build_flags=-tags=sqlite_unlock_notify

run:
	go run ${go_build_flags} ./cmd/antibrut run

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
	go test -race ./... -count 1