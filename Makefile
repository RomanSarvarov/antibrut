go_build_flags=-tags=sqlite_unlock_notify

run:
	go run ${go_build_flags} ./cmd/antibrut run

generate:
	go generate ./...

lint:
	golangci-lint run ./...
	make proto-lint

proto-lint:
	cd proto && go run github.com/bufbuild/buf/cmd/buf lint

test:
	go test -race ./...