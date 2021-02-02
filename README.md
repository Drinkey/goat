# goat
GO At manages cron jobs of running machine and support view and run via RESTFul API

# Installation

## Linux running systemd

1. Run the following commands to build binary
   ```bash
    $ export GOARCH=amd64
	$ export GOOS=linux
	$ export CGO_ENABLED=1
	$ go build -v -o goat-linux-amd64
    ```
2. Move the binary to $PATH using the following command
    ```bash
    $ mv goat-linux-amd64 /usr/local/bin/goat
    $ chmod +x /usr/local/bin/goat
    ```
3. Update `goat.service` to satisfy your environment. **Especially if you don't want goat run as root**
4. Use `systemctl` to control the goat service
    ```bash
    $ sudo systemctl enable goat # start the service after system boot
    $ sudo systemctl start goat # start the service
    $ sudo systemctl status goat # query service status
    $ sudo systemctl stop goat # shutdown the service
    ```

# Usage

## Service
1. Use `sudo systemctl start goat` to start service
2. Use RESTFul API to access the service
3. Swagger document available at http://\<ip\>:8090/api/v1/swagger/index.html

## crontab configuration

To have better maintainability to your cron jobs, it is recommended to have comment before each task. `goat` will read this value as task title. It is not mandatory but recommended.

For example, if you have only one cron task:
```
02 14 * * * python3 -u /home/example/thg/scripts/uid.py
```

Add a comment prior to this line as the title
```
# Update whatever info 2:02PM every day
02 14 * * * python3 -u /home/example/thg/scripts/uid.py
```

**Please use exactly `# ` (# and one space)at the beginning of comment line**. `goat` will understand it as the cron task title.

> Even without using goat, you should also give a comment for each task in case you forgot what the task does.

## Task command definition

`os/exec` would use the first parameter to find executable in `$PATH`. If the first word of `command` is not executable in `$PATH`, this will trigger error.

And if your command contains `&&`, `|`, `;`, `>`, `>>`, or other symbols in shell command line, the execution will fail. The solution is wrap the commands like this in a shell script and call the shell script in crontab.

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
