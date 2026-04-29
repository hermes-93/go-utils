package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hermes-93/go-utils/internal/health"
)

func TestCheckHealthy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	result := health.Check(srv.URL, health.Options{Timeout: 5, Retries: 0, ExpectedCode: 200})
	if !result.Healthy {
		t.Fatalf("expected healthy, got error: %s", result.Error)
	}
	if result.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", result.StatusCode)
	}
	if result.Attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", result.Attempts)
	}
}

func TestCheckUnhealthy503(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))
	defer srv.Close()

	result := health.Check(srv.URL, health.Options{Timeout: 5, Retries: 0, ExpectedCode: 200})
	if result.Healthy {
		t.Fatal("expected unhealthy for 503")
	}
	if result.StatusCode != 503 {
		t.Fatalf("expected 503, got %d", result.StatusCode)
	}
}

func TestCheckConnectionRefused(t *testing.T) {
	result := health.Check("http://127.0.0.1:1", health.Options{Timeout: 1, Retries: 0, ExpectedCode: 200})
	if result.Healthy {
		t.Fatal("expected unhealthy for refused connection")
	}
	if result.Error == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestCheckRetriesSucceed(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	result := health.Check(srv.URL, health.Options{
		Timeout:       5,
		Retries:       3,
		RetryInterval: 0,
		ExpectedCode:  200,
	})
	if !result.Healthy {
		t.Fatalf("expected healthy after retries, error: %s", result.Error)
	}
	if result.Attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", result.Attempts)
	}
}

func TestResultString(t *testing.T) {
	r := health.Result{URL: "http://example.com", Healthy: true, StatusCode: 200, Attempts: 1}
	s := r.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}

	r2 := health.Result{URL: "http://example.com", Healthy: false, Error: "connection refused", Attempts: 3}
	s2 := r2.String()
	if s2 == "" {
		t.Fatal("expected non-empty string for error result")
	}
}
