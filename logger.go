package ginlogger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

var defaultLogger *logrus.Logger

func init() {
	defaultLogger = logrus.New()
	// Apply stackdriver format to logrus logger
	defaultLogger.SetFormatter(joonix.NewFormatter())
}

const (
	timeFormat = "2006/01/02 15:04:05.000"
)

var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	// More request info
	protocol := param.Request.Proto
	ua := param.Request.UserAgent()
	var bodyString string
	if body, err := ioutil.ReadAll(param.Request.Body); err == nil {
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
	// Logger is a logrus logger
	Logger *logrus.Logger

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
	logger := conf.Logger
	if logger == nil {
		logger = defaultLogger
	}

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

		// clone a body reader
		var bodyReader io.ReadCloser
		if body, err := ioutil.ReadAll(c.Request.Body); err == nil {
			bodyReader = ioutil.NopCloser(bytes.NewReader(body))
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		}

		// Process request
		c.Next()

		// Ignore when path is being skipped
		if ok := skip[path]; ok {
			return
		}

		// set body reader to the clone reader to read body again
		c.Request.Body = bodyReader
		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = time.Since(start)

		// Get request info
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.Path = path

		// Get response info
		param.StatusCode = c.Writer.Status()
		param.BodySize = c.Writer.Size()

		// Get error info
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

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
