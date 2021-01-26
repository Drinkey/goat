# goat
GO At manages cron jobs of running machine and support view and run via RESTFul API

# Installation


1. Update `goat.service` to satisfy your environment. **Especially if you don't want goat run as root**
2. Run `make build && make install`

# Usage

## Service
1. Use `sudo systemctl start goat` to start service
2. Use RESTFul API to access the service

## crontab configuration

To have better maintainability to your cron jobs, it is recommended to have comment before each task. `goat` will read this value as task title. It is not mandatory but recommended.

For example, if you have only one cron task:
```
02 14 * * * python3 -u /home/example/thg/scripts/uid.py >> /tmp/uid.log 2>&1
```

Add a comment prior to this line as the title
```
# Update whatever info 2:02PM every day
02 14 * * * python3 -u /home/example/thg/scripts/uid.py >> /tmp/uid.log 2>&1
```

**Please use exactly `# ` (# and one space)at the beginning of comment line**. `goat` will understand it as the cron task title.

> Even without using goat, you should also give a comment for each task in case you forgot what the task does.

# Development

## Build with Docker

```
docker run --rm \
    -v "$PWD":/usr/src/goat \
    -w /usr/src/goat \
    -e GOOS=linux \
    -e GOARCH=amd64 \
    -e CGO_ENABLED=1 \
    goat:latest \
    go build -v -o goat-linux-amd64
```
