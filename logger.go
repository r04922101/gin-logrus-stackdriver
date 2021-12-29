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

var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	// More request info
	protocol := param.Request.Proto
	ua := param.Request.UserAgent()
	var bodyString string
	body, err := ioutil.ReadAll(param.Request.Body)
	if err == nil {
		bodyString = string(body)
	}

	return fmt.Sprintf("%v - %s %s, req: \"%s %7s %s %s\", res: \"%3d\", latency: %.13v",
		param.TimeStamp.Format(timeFormat),
		param.ClientIP,
		ua,
		protocol,
		param.Method,
		param.Path,
		bodyString,
		param.StatusCode,
		param.Latency,
	)
}

// LoggerConfig defines the config for Logger middleware
type LoggerConfig struct {
	// Optional. Default value is defaultLogFormatter
	Formatter gin.LogFormatter

	// SkipPaths is a url path array which logs are not written
	// Optional
	SkipPaths []string
}

// NewLogger instances a Logger middleware.
func NewLogger() gin.HandlerFunc {
	return NewLoggerWithConfig(LoggerConfig{})
}

// NewLoggerWithConfig instance a Logger middleware with config.
func NewLoggerWithConfig(conf LoggerConfig) func(c *gin.Context) {
	logger := logrus.New()
	logger.SetFormatter(joonix.NewFormatter())

	var skip map[string]bool
	if length := len(conf.SkipPaths); length > 0 {
		skip = make(map[string]bool, length)

		for _, path := range conf.SkipPaths {
			skip[path] = true
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Ignore when path is being skipped
		if ok := skip[path]; ok {
			return
		}

		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		// Get request info
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}
		param.Path = path

		// Get response info
		param.StatusCode = c.Writer.Status()
		param.BodySize = c.Writer.Size()

		// If no formatter assigned, use default formatter
		formatter := conf.Formatter
		if formatter == nil {
			formatter = defaultLogFormatter
		}

		msg := formatter(param)
		// Log level: error for code >= 500, warn for 500 > code >= 400, info for others
		if param.StatusCode >= http.StatusInternalServerError {
			logger.Error(msg)
		} else if param.StatusCode >= http.StatusBadRequest {
			logger.Warn(msg)
		} else {
			logger.Info(msg)
		}
		// If any error occurs, print error message in a new line
		if len(c.Errors) > 0 {
			logger.Errorf("err: %s", param.ErrorMessage)
		}
	}
}
