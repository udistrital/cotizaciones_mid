package clients

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHTTPClientDoSendsRequestAndReadsResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("método = %q; se esperaba POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q; se esperaba application/json", got)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":123}`))
	}))
	defer server.Close()

	client := NewHTTPClient(2 * time.Second)
	status, headers, body, err := client.Do(
		context.Background(),
		http.MethodPost,
		server.URL,
		[]byte(`{"titulo":"Prueba"}`),
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		t.Fatalf("Do devolvió un error inesperado: %v", err)
	}
	if status != http.StatusCreated {
		t.Errorf("status = %d; se esperaba 201", status)
	}
	if headers.Get("Content-Type") == "" {
		t.Error("se esperaba Content-Type en la respuesta")
	}
	if string(body) != `{"id":123}` {
		t.Errorf("cuerpo = %q; valor inesperado", body)
	}
}

func TestHTTPClientDoRejectsEmptyMethod(t *testing.T) {
	client := NewHTTPClient(time.Second)

	_, _, _, err := client.Do(
		context.Background(), "", "http://localhost", nil, nil,
	)
	if err == nil || !strings.Contains(err.Error(), "método HTTP") {
		t.Fatalf("se esperaba error por método vacío; recibido: %v", err)
	}
}

func TestHTTPClientDoRejectsEmptyURL(t *testing.T) {
	client := NewHTTPClient(time.Second)

	_, _, _, err := client.Do(
		context.Background(), http.MethodGet, " ", nil, nil,
	)
	if err == nil || !strings.Contains(err.Error(), "URL") {
		t.Fatalf("se esperaba error por URL vacía; recibido: %v", err)
	}
}

func TestHTTPClientDoReportsTransportError(t *testing.T) {
	client := NewHTTPClient(time.Second)

	_, _, _, err := client.Do(
		context.Background(),
		http.MethodGet,
		"http://127.0.0.1:1",
		nil,
		nil,
	)
	if err == nil {
		t.Fatal("se esperaba un error de conexión")
	}
	if !strings.Contains(err.Error(), "ejecutar solicitud HTTP") {
		t.Errorf("mensaje de error inesperado: %v", err)
	}
}

func TestNewHTTPClientUsesDefaultTimeout(t *testing.T) {
	client := NewHTTPClient(0)
	if client.client.Timeout != 30*time.Second {
		t.Errorf("timeout = %s; se esperaba 30s", fmt.Sprint(client.client.Timeout))
	}
}
