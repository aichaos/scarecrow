init:
	go get ./...

run:
	go run cmd/scarecrow/main.go

build:
	go build -o scarecrow-cli cmd/scarecrow/main.go

format:
	gofmt -w .
