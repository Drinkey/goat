FROM golang:1.15.5

WORKDIR /go/src/goat
COPY . .

ENV GOPROXY="https://mirrors.aliyun.com/goproxy/"

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["goat"]
