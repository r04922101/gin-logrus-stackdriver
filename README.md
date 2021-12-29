# gin-logrus-stackdriver

## Installation

Download Go module:

```sh
go get -u github.com/r04922101/gin-logrus-stackdriver
```

## How to Use

### Import it in your code

```go
import ginlogger "github.com/r04922101/gin-logrus-stackdriver"
```

### Example

Check the example code in [example/main.go](./example/main.go)

Run `example/main.go`

```sh
go run main.go
```

Send an HTTP POST request to `localhost:8080/ping`

```sh
curl -X POST 'localhost:8080/ping' \
-H 'Content-Type: text/plain' \
-d 'pang'
```

### Result

Check your console, and see the gin log is in stackdriver format

```sh
{"message":"2021/12/29 14:20:09.428 - ::1 PostmanRuntime/7.26.8, req: \"HTTP/1.1    POST /ping pang\", res: \"200\", latency: 107.049µs","severity":"INFO","timestamp":{"seconds":1640758809,"nanos":428314000}}
```

## Customize Formatter

Refer to [Gin Custom Log Format](https://github.com/gin-gonic/gin#custom-log-format), write your custom formatter function. \
Construct a `LoggerConfig` with the formatter and pass it to `NewLoggerWithConfig` function. \
You can also find a example snippet of custom formatter in [example/main.go](./example/main.go).

## Customize Logger

If you want to use your own `logrus.Logger` instance instead of default logger in stackdriver format, just pass your logger to `LoggerConfig`.
