.PHONY: air tailwind clean

run: tailwind
	go run ./cmd/web

air: tailwind
	go build -o ./tmp/main ./cmd/web

tailwind:
	./bin/tailwindcss -c ./ui/tailwind/tailwind.config.js \
	-i ./ui/tailwind/input.css -o ./ui/static/main.css

clean:
	rm -rf tmp
