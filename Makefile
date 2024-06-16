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

css:
	@npx tailwindcss -i view/css/app.css -o public/styles.css --watch

templ:
	@templ generate --watch --proxy=http://localhost:3000

build:
	@templ generate view
	@go build -o bin/picturethisai main.go

up: ## DB migrate up
	@go run cmd/migrate/main.go up

down: ## DB migrate down
	@go run cmd/migrate/main.go down

reset:
	@go run cmd/reset/main.go up

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))