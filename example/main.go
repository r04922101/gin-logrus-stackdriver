package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginlogger "github.com/r04922101/gin-logrus-stackdriver"
)

// Custom log formatter
func customLog(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[My log] %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

// Req defines a data structure for client request
type Req struct {
	Value int `json:"value"`
}

func main() {
	r := gin.New()
	r.Use(ginlogger.NewLogger(), gin.Recovery())

	// customize formatter
	// conf := ginlogger.LoggerConfig{
	// 	Formatter: customLog,
	// }
	// r.Use(ginlogger.NewLoggerWithConfig(conf), gin.Recovery())

	r.POST("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/req", func(c *gin.Context) {
		req := new(Req)
		if err := c.BindJSON(req); err != nil {
			c.String(http.StatusBadRequest, "no value")
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("Got value %d", req.Value))
	})

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run gin: %v", err)
	}
}
