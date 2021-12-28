# gin-logrus-stackdriver

## Installation

Download Go module:

```sh
go get -u github.com/r04922101/gin-logrus-stackdriver
```

## How to Use

### Import it in your code

```go
import "github.com/gin-gonic/gin"
```

### Example

Assume the following code in `main.go`

```go
package main

import (
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
  ginlogger "github.com/r04922101/gin-logrus-stackdriver"
)

func main() {
  r := gin.New()
  r.Use(ginlogger.NewLogger())

  r.POST("/ping", func(c *gin.Context) {
    c.String(http.StatusOK, "pong")
  })

  if err := r.Run(); err != nil {
    log.Fatalf("failed to run gin: %v", err)
  }
}

```

Run `main.go`

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
{"message":"2021/12/28 17:35:01.467 - ::1 curl/7.64.1, req: \"HTTP/1.1    POST /ping pang\", res: \"200\", latency: 26.403µs","severity":"INFO","timestamp":{"seconds":1640684101,"nanos":467857000}}
```
