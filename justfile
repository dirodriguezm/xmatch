default:
	just --list

migrate db:
	migrate -database sqlite3://{{db}}.db -path service/internal/db/migrations up

[working-directory: 'service']
build:
	go build -o build/main cmd/*.go

[working-directory: 'service']
test:
	grc go test ./... -race

[working-directory: 'service']
test-verbose:
	grc go test -v ./... -race

export CONFIG_PATH := env_var_or_default("CONFIG_PATH", "service/config.yaml")

[working-directory: 'service']
run application flags='' $LOG_LEVEL="debug": build
	./build/main {{flags}} {{application}}

clean-build:
	rm -r service/build

[working-directory: 'service']
clean-all:
	go clean
	go clean -testcache
	rm -r build

clean-db db:
	rm {{db}}.db

[working-directory: 'service']
docs:
	swag init --dir ./ --generalInfo ./cmd/*.go --output ./docs

export USE_LOGGER := "true"
export ENVIRONMENT := "local"
export LOG_LEVEL := "debug"

[working-directory: 'service']
live-server:
	air 



[working-directory: 'service']
mock:
  go run github.com/vektra/mockery/v3@v3.7.0
