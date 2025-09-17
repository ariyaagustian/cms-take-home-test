package e2e

import (
	"net/http"
	"testing"
)

func TestHealthz(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/healthz")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
