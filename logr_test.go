package ginlogr

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	buffer := new(bytes.Buffer)
	zapLog := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(buffer),
			zapcore.DebugLevel,
		),
	)
	log := zapr.NewLogger(zapLog)
	r := gin.New()
	r.Use(Logger(log))
	r.GET("/example", func(c *gin.Context) {})
	r.GET("/example-error", func(c *gin.Context) { c.Error(errors.New("lorem")) })

	buffer.Reset()
	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/example", nil)
	r.ServeHTTP(res1, req1)
	logMsg := buffer.String()
	assert.Contains(t, logMsg, "200")

	buffer.Reset()
	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/example-error", nil)
	r.ServeHTTP(res2, req2)
	logMsg = buffer.String()
	assert.Contains(t, logMsg, "lorem")
}
