run:
	go run ./cmd/antibrut run

generate:
	go generate ./...

lint:
	golangci-lint run ./...
	make proto-lint

proto-lint:
	cd proto && go run github.com/bufbuild/buf/cmd/buf lint