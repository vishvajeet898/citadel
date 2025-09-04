package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter
		requestPayload := getRequestPayload(c)
		c.Next()
		responsePayload := getResponsePayload(c, bodyLogWriter)
		if responsePayload["response_body"] == "" || responsePayload["response_body"] == nil {
			return
		}
		if isUrlSkippable(c.Request.Method, c.Request.RequestURI) {
			return
		}
		logMap := map[string]interface{}{}
		if traceID, ok := c.Get("trace_id"); ok {
			logMap["trace_id"] = traceID
		}
		logMap["request"] = requestPayload
		logMap["response"] = responsePayload
		logMap["duration"] = time.Since(start).Milliseconds()
		logData, _ := json.Marshal(logMap)
		log.Info().Str("log", string(logData)).Msg("API Log")
	}
}

func getRequestPayload(c *gin.Context) map[string]interface{} {
	startTime := time.Now()
	// env config
	environment := os.Getenv("API_ENV")
	if environment == "" {
		environment = "development"
	}

	var bodyBytes []byte
	data := make(map[string]interface{})
	bodyBytes, _ = io.ReadAll(c.Request.Body)
	buffer := io.NopCloser(bytes.NewBuffer(bodyBytes))
	headers := c.Request.Header.Clone()
	contentType := headers.Get("Content-Type")
	if headers != nil {
		delete(headers, "api_key")
	}

	if bodyBytes != nil {
		data["request_body"] = string(bodyBytes)
	} else {
		data["request_body"] = ""
	}
	c.Request.Body = buffer
	data["start_time"] = startTime
	data["environment"] = environment
	data["remote_address"] = c.ClientIP()
	data["request_method"] = c.Request.Method
	data["host_name"] = c.Request.Host
	data["content_type"] = contentType
	data["request_path"] = c.Request.RequestURI

	// masking authorization token in logs
	if headers.Get("Authorization") != "" {
		headers["Authorization"] = []string{"*****"}
	}
	data["headers"] = headers
	return data
}

func getResponsePayload(c *gin.Context, bodyLogWriter *bodyLogWriter) map[string]interface{} {
	endTime := time.Now()
	data := make(map[string]interface{})
	contentType := c.Request.Header.Get("Content-Type")

	// env config
	environment := os.Getenv("API_ENV")
	if environment == "" {
		environment = "development"
	}

	// response
	if contentType == constants.ContentTypeJson || contentType == constants.ContentTypeJsonWithCharsetUtf8 {
		responseBody := bodyLogWriter.body.String()
		data["response_body"] = responseBody
	} else if c.Request.Method == "GET" {
		data["response_body"] = ""
	} else {
		data["response_body"] = "Failed to fetch response"
	}
	data["remote_address"] = c.ClientIP()
	data["request_method"] = c.Request.Method
	data["host_name"] = c.Request.Host
	data["request_path"] = c.Request.RequestURI
	data["end_time"] = endTime
	data["status_code"] = c.Writer.Status()
	data["env"] = environment
	return data
}

func isUrlSkippable(method, url string) bool {
	skipUrls := []string{
		"/ping",
		"/api/v1/search/tasks",
	}

	for _, skipUrl := range skipUrls {
		if method == http.MethodGet && strings.HasPrefix(url, skipUrl) {
			return true
		}
	}

	return false
}
