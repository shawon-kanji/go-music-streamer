.PHONY: run build test clean

docker-run:
	docker-compose up -d

docker-down:
	docker-compose down

run:
	go run cmd/api/main.go

seed:
	go run cmd/seed/main.go

build:
	go build -o bin/server cmd/api/main.go

test:
	go test ./...

clean:
	rm -rf bin
