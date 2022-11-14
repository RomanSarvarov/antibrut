package mock

//go:generate find . ! -name gen.go -type f -o -type d -maxdepth 1 -mindepth 1 -exec rm -rf {} +

// sqlite
//go:generate go run github.com/vektra/mockery/v2 --name database --keeptree --case underscore --dir "../sqlite" --output "./sqlite" --outpkg "mock" --exported

// leaky bucket
//go:generate go run github.com/vektra/mockery/v2 --name Repository --keeptree --case underscore --dir "../leakybucket" --output "./leakybucket" --outpkg "mock" --exported

// service
//go:generate go run github.com/vektra/mockery/v2 --name repository --keeptree --case underscore --dir "../" --output "./" --outpkg "mock" --exported
//go:generate go run github.com/vektra/mockery/v2 --name rateLimiter --keeptree --case underscore --dir "../" --output "./" --outpkg "mock" --exported
