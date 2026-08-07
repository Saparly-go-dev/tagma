package handler

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadLoggableRequestBodyRestoresBodyAndRedactsSecrets(t *testing.T) {
	body := `{"username":"agent","password":"secret","currency":10}`
	request := httptest.NewRequest("POST", "/api/payments/save", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	logged := readLoggableRequestBody(request)
	if strings.Contains(logged, "secret") || !strings.Contains(logged, "[REDACTED]") {
		t.Fatalf("sensitive request body was not redacted: %s", logged)
	}

	restored, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatalf("read restored body: %v", err)
	}
	if string(restored) != body {
		t.Fatalf("restored body = %q, want %q", restored, body)
	}
}

func TestSanitizeQueryRedactsToken(t *testing.T) {
	request := httptest.NewRequest("GET", "/?token=secret&price1=10", nil)
	logged := sanitizeQuery(request.URL.Query())
	if strings.Contains(logged, "secret") || !strings.Contains(logged, "price1=10") {
		t.Fatalf("unexpected sanitized query: %s", logged)
	}
}
