package health

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckServerHealth_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	err := CheckServerHealth(mockServer.URL, "test.log", 5)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

func TestCheckServerHealth_Unhealthy(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	err := CheckServerHealth(mockServer.URL, "test.log", 5)
	if err == nil {
		t.Fatalf("Expected error, but got none")
	}

	expectedErrorMsg := fmt.Sprintf("ERROR: Server %s is unreachable or returned status %d", mockServer.URL, 500)
	if err.Error() != expectedErrorMsg {
		t.Fatalf("Expected error message to be '%s', but got '%s'", expectedErrorMsg, err.Error())
	}
}

func TestCheckServerHealth_Unreachable(t *testing.T) {
	err := CheckServerHealth("https://nonexistent-server.com", "test.log", 2)
	if err == nil {
		t.Fatalf("Expected error for unreachable server, but got none")
	}
}

func TestCheckServerHealth_Timeout(t *testing.T) {
	slowServer := http.Server{
		Addr: ":8081",
	}

	go func() {
		time.Sleep(2 * time.Second)
		slowServer.Close()
	}()

	err := CheckServerHealth("http://localhost:8081", "test.log", 1) // 1 second timeout

	if err == nil || err.Error() != "ERROR: Server http://localhost:8081 is unreachable or returned status 502" {
		t.Errorf("Expected timeout error, got %v", err)
	}
}
