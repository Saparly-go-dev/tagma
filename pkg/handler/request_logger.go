package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const loggedBodyLimit = 8 << 10

type loggingResponseWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *loggingResponseWriter) Write(data []byte) (int, error) {
	w.capture(data)
	return w.ResponseWriter.Write(data)
}

func (w *loggingResponseWriter) WriteString(data string) (int, error) {
	w.capture([]byte(data))
	return w.ResponseWriter.WriteString(data)
}

func (w *loggingResponseWriter) capture(data []byte) {
	remaining := loggedBodyLimit - w.body.Len()
	if remaining <= 0 {
		return
	}
	if len(data) > remaining {
		data = data[:remaining]
	}
	_, _ = w.body.Write(data)
}

func requestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		requestBody := readLoggableRequestBody(c.Request)
		writer := &loggingResponseWriter{ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		fields := logrus.Fields{
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"query":         sanitizeQuery(c.Request.URL.Query()),
			"status":        c.Writer.Status(),
			"response_size": c.Writer.Size(),
			"duration_ms":   time.Since(started).Milliseconds(),
			"client_ip":     c.ClientIP(),
		}
		if requestBody != "" {
			fields["request_body"] = requestBody
		}
		if responseBody := loggableResponseBody(c.Writer.Header().Get("Content-Type"), writer.body.Bytes()); responseBody != "" {
			fields["response_body"] = responseBody
		}
		if userID, ok := c.Get(userCtx); ok {
			fields["user_id"] = userID
		}
		if role, ok := c.Get(roleCtx); ok {
			fields["user_role"] = role
		}

		entry := logrus.WithFields(fields)
		if len(c.Errors) > 0 {
			entry = entry.WithField("errors", c.Errors.String())
		}
		entry.Info("http request")
	}
}

func readLoggableRequestBody(request *http.Request) string {
	if request.Body == nil || request.ContentLength <= 0 || request.ContentLength > loggedBodyLimit {
		return ""
	}
	if !strings.Contains(request.Header.Get("Content-Type"), "application/json") {
		return ""
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return ""
	}
	request.Body = io.NopCloser(bytes.NewReader(body))
	return sanitizeJSON(body)
}

func loggableResponseBody(contentType string, body []byte) string {
	if !strings.Contains(contentType, "application/json") || len(body) == 0 {
		return ""
	}
	return sanitizeJSON(body)
}

func sanitizeJSON(body []byte) string {
	var value interface{}
	if err := json.Unmarshal(body, &value); err != nil {
		return string(body)
	}
	redactSensitiveFields(value)
	sanitized, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(sanitized)
}

func redactSensitiveFields(value interface{}) {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, item := range typed {
			if isSensitiveKey(key) {
				typed[key] = "[REDACTED]"
				continue
			}
			redactSensitiveFields(item)
		}
	case []interface{}:
		for _, item := range typed {
			redactSensitiveFields(item)
		}
	}
}

func sanitizeQuery(query url.Values) string {
	clean := make(url.Values, len(query))
	for key, values := range query {
		clean[key] = append([]string(nil), values...)
	}
	for key := range clean {
		if isSensitiveKey(key) {
			clean.Set(key, "[REDACTED]")
		}
	}
	return clean.Encode()
}

func isSensitiveKey(key string) bool {
	key = strings.ToLower(key)
	return strings.Contains(key, "password") || strings.Contains(key, "token") ||
		strings.Contains(key, "authorization") || strings.Contains(key, "secret")
}
