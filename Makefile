run:
	go run cmd/scarecrow/main.go

build:
	go build -o scarecrow cmd/scarecrow/main.go

format:
	gofmt -w .
