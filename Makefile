build:
	go run main.go

run: tailwindcss templ build

test:
	go test -v ./...

format:
	templ fmt .

tailwindcss:
	npx tailwindcss --config config/tailwind.config.js -i config/input.css -o static/css/styles.css

templ:
	templ generate