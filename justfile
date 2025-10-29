default:
	just --list

migrate db:
	migrate -database sqlite3://{{db}}.db -path service/internal/db/migrations up

[working-directory: 'service']
build: build-css
	go build -o build/main cmd/*.go

[working-directory: 'service']
test:
	grc go test ./... -race

[working-directory: 'service']
test-verbose:
	grc go test -v ./... -race

export CONFIG_PATH := env_var_or_default("CONFIG_PATH", "service/config.yaml")

run application flags='' $LOG_LEVEL="debug": build
	./service/build/main {{flags}} {{application}}

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
build-css-watch:
	tailwindcss --input ./ui/static/css/tailwind.css --output ./ui/static/css/output.css --watch --optimize --minify

[working-directory: 'service']
build-css:
	tailwindcss --input ./ui/static/css/tailwind.css --output ./ui/static/css/output.css --optimize --minify

[working-directory: 'service']
mock:
  mockery
