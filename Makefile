run:
	go run cmd/scarecrow/main.go

build:
	go build -o scarecrow-cli cmd/scarecrow/main.go

format:
	gofmt -w .

save-deps:
	godep save -r ./ ./listeners/console ./listeners/slack ./listeners/xmpp
