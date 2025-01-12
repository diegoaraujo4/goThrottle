package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func TestResponder(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(responder)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestMainFunction(t *testing.T) {
	client, mock := redismock.NewClientMock()
	redisClient = client

	// Start the server in a goroutine
	go func() {
		main()
	}()
	time.Sleep(1 * time.Second) // Give the server a second to start

	mock.ExpectGet("ip:::1:block").RedisNil()
	mock.ExpectIncr("ip:::1").SetVal(int64(1))
	mock.ExpectExpire("ip:::1", time.Second).SetVal(true)

	// Test if the server is running and responding
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}
}
