run: build
	@./bin/picturethisai

install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D tailwindcss
	@npm install -D daisyui@latest

build:
	## tailwindcss -i view/css/input.css -o public/styles.css
	@templ generate view
	@go build -o bin/picturethisai

up: ## DB migrate up
	@go run cmd/migrate/main.go up

down: ## DB migrate down
	@go run cmd/migrate/main.go down

drop:
	@go run cmd/drop/main.go up

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))