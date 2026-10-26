package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONError(t *testing.T) {
	recorder := httptest.NewRecorder()

	WriteJSONError(
		recorder,
		http.StatusBadRequest,
		"DATOS_INVALIDOS",
		"La solicitud contiene datos inválidos.",
	)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}

	if got := recorder.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q; valor inesperado", got)
	}

	var response APIError
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}

	if response.Error.Codigo != "DATOS_INVALIDOS" {
		t.Errorf("código inesperado: %q", response.Error.Codigo)
	}
	if response.Error.Mensaje != "La solicitud contiene datos inválidos." {
		t.Errorf("mensaje inesperado: %q", response.Error.Mensaje)
	}
}
