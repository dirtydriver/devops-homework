package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	server := NewServer()
	handler := server.Handler("test")

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	expected := `{"status":"ok"}` + "\n"

	if recorder.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, recorder.Body.String())
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected application/json, got %q", contentType)
	}

}

func TestVersionHandler(t *testing.T) {
	server := NewServer()
	handler := server.Handler("test")

	request := httptest.NewRequest(http.MethodGet, "/version", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	expected := `{"version":"1.0.0"}` + "\n"

	if recorder.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, recorder.Body.String())
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected application/json, got %q", contentType)
	}

}

func TestEnvironmentnHandler(t *testing.T) {
	server := NewServer()
	handler := server.Handler("test")

	request := httptest.NewRequest(http.MethodGet, "/env", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	expected := `{"environment":"test"}` + "\n"

	if recorder.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, recorder.Body.String())
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected application/json, got %q", contentType)
	}

}

func TestCreateConfigHandler(t *testing.T) {
	server := NewServer()
	handler := server.Handler("test")
	body := `{
		"name": "database_url",
		"value": "postgres://example"
	}`

	request := httptest.NewRequest(http.MethodPost, "/config", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	expected := `{"name":"database_url","value":"postgres://example"}` + "\n"

	if recorder.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, recorder.Body.String())
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected application/json, got %q", contentType)
	}

}

func TestGetConfigHandler(t *testing.T) {
	server := NewServer()
	server.config["database_url"] = "postgres://example"
	handler := server.Handler("test")

	request := httptest.NewRequest(http.MethodGet, "/config/database_url", nil)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response configResponse

	json.NewDecoder(recorder.Body).Decode(&response)

	if response.Name != "database_url" {
		t.Errorf("expacted name %q, got %q", "database_url", response.Name)
	}

	if response.Value != "postgres://example" {
		t.Errorf("expacted name %q, got %q", "postgres://example", response.Value)
	}

}

func TestDeleteConfigHandler(t *testing.T) {
	server := NewServer()
	server.config["database_url"] = "postgres://example"
	handler := server.Handler("test")

	request := httptest.NewRequest(http.MethodDelete, "/config/database_url", nil)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response deleteResponse

	json.NewDecoder(recorder.Body).Decode(&response)

	if !response.Deleted {
		t.Errorf("expacted deleted to be true, got %t", response.Deleted)
	}

	getRequest := httptest.NewRequest(
		http.MethodGet, "/config/database_url", nil,
	)

	getRecorder := httptest.NewRecorder()

	handler.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusNotFound {
		t.Errorf(
			"expected deleted config to return status %d, got %d", http.StatusNotFound, getRecorder.Code,
		)
	}

}
