test:
	go test -v ./...

swagger:
	swag init -g goat.go

build:
	docker run --rm -v "$(PWD)":/usr/src/goat -w /usr/src/goat -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=1 goat:latest go build -v -o goat-linux-amd64
	tar zcvf goat-linux-amd64.tar.gz goat-linux-amd64

build-release:
	export GOARCH=amd64
	export GOOS=linux
	export CGO_ENABLED=1
	go build -v -o goat-linux-amd64
	tar zcvf goat-`git tag |sort|tail -n1`-linux-amd64.tar.gz goat-linux-amd64

build-mac:
	export GOARCH=amd64
	export GOOS=darwin
	go build -v -o goat-darwin-amd64

install:
	cp goat-linux-amd64 /usr/local/bin/goat && chmod +x /usr/local/bin/goat
	cp goat.service /etc/systemd/system/goat.service