run:
	go run ./cmd/antibrut run

generate:
	go generate ./...

lint:
	golangci-lint run ./...
	make proto-lint

proto-lint:
	cd proto && buf lint