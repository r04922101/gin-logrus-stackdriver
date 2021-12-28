package ginlogger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006/01/02 15:04:05.000"
)

// NewLogger instances a middleware using logrus logger in stackdriver format
func NewLogger(notLogged ...string) gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(joonix.NewFormatter())

	var skip map[string]bool
	if length := len(notLogged); length > 0 {
		skip = make(map[string]bool, length)

		for _, path := range notLogged {
			skip[path] = true
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Ignore when path is being skipped
		if ok := skip[path]; ok {
			return
		}

		// Stop timer
		timeStamp := time.Now()
		latency := timeStamp.Sub(start)

		// Get request info
		protocol := c.Request.Proto
		clientIP := c.ClientIP()
		ua := c.Request.UserAgent()
		method := c.Request.Method
		var bodyString string
		body, err := ioutil.ReadAll(c.Request.Body)
		if err == nil {
			bodyString = string(body)
		}
		if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
			path = path + "?" + rawQuery
		}
		// Get response status and error
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		msg := fmt.Sprintf("%v - %s %s, req: \"%s %7s %s %s\", res: \"%3d\", latency: %.13v",
			timeStamp.Format(timeFormat),
			clientIP,
			ua,
			protocol,
			method,
			path,
			bodyString,
			statusCode,
			latency,
		)
		if statusCode >= http.StatusInternalServerError {
			logger.Error(msg)
		} else if statusCode >= http.StatusBadRequest {
			logger.Warn(msg)
		} else {
			logger.Info(msg)
		}
		if len(c.Errors) > 0 {
			logger.Errorf("err: %s", errorMessage)
		}
	}
}
