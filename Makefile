test:
	go test -v ./...

swagger:
	swag init -g goat.go

build:
	export GOARCH=amd64 && export CGO_ENABLED=1 && go build -v -o goat-linux-amd64

install:
	cp goat-linux-amd64 /usr/local/bin/goat && chmod +x /usr/local/bin/goat
	cp goat.service /etc/systemd/system/goat.service