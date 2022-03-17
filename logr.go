package ginlogr

import (
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"strings"
	"time"
)

// Config defines the config for logger middleware
type LoggerConfig struct {
	// UTC a boolean stating whether to use UTC time zone or local.
	UTC bool

	TimeFormat string

	// LogV is the param for logger.V()
	LogV int

	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

func Logger(l logr.Logger) gin.HandlerFunc {
	return LoggerWithConfig(l, LoggerConfig{TimeFormat: time.RFC3339})
}

func LoggerWithConfig(l logr.Logger, conf LoggerConfig) gin.HandlerFunc {
	l = l.WithName("GIN")

	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			l.WithName("GIN")
			end := time.Now()
			latency := end.Sub(start)
			if conf.UTC {
				end = end.UTC()
			}

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				l.Error(c.Errors.Last(), strings.Join(c.Errors.Errors(), " "))
			} else {
				fields := []interface{}{
					"time", end.Format(conf.TimeFormat),
					"status", c.Writer.Status(),
					"method", c.Request.Method,
					"path", path,
					"query", query,
					"ip", c.ClientIP(),
					"user-agent", c.Request.UserAgent(),
					"latency", latency,
				}
				l.V(conf.LogV).Info(path, fields...)
			}
		}
	}
}
