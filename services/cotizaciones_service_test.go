package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/udistrital/cotizaciones_mid/config"
)

func TestCotizacionesServiceRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/solicitud_cotizacion" {
			t.Errorf("ruta recibida = %q; valor inesperado", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("método recibido = %q; se esperaba GET", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"resultado":"ok"}`))
	}))
	defer server.Close()

	service, err := NewCotizacionesService(config.Config{
		CotizacionesCRUDURL: server.URL,
	})
	if err != nil {
		t.Fatalf("no se pudo crear el servicio: %v", err)
	}

	status, _, body, err := service.Request(
		context.Background(),
		http.MethodGet,
		"/v1/solicitud_cotizacion",
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("Request devolvió un error inesperado: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("status = %d; se esperaba 200", status)
	}
	if string(body) != `{"resultado":"ok"}` {
		t.Errorf("cuerpo inesperado: %s", body)
	}
}

func TestCotizacionesServiceRejectsAbsolutePath(t *testing.T) {
	service, err := NewCotizacionesService(config.Config{
		CotizacionesCRUDURL: "http://localhost:8085",
	})
	if err != nil {
		t.Fatalf("no se pudo crear el servicio: %v", err)
	}

	_, _, _, err = service.Request(
		context.Background(),
		http.MethodGet,
		"http://otro-servidor/recurso",
		nil,
		nil,
	)
	if err == nil {
		t.Fatal("se esperaba un error para una ruta absoluta")
	}
}

func TestNewCotizacionesServiceRejectsInvalidBaseURL(t *testing.T) {
	_, err := NewCotizacionesService(config.Config{
		CotizacionesCRUDURL: "no-es-una-url",
	})
	if err == nil {
		t.Fatal("se esperaba un error para una URL base inválida")
	}
}

func TestCotizacionesServicePreservesDownstreamHTTPStatus(t *testing.T) {
	statuses := []int{
		http.StatusBadRequest,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	}

	for _, expectedStatus := range statuses {
		t.Run(http.StatusText(expectedStatus), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(expectedStatus)
				_, _ = w.Write([]byte(`{"error":"respuesta del CRUD"}`))
			}))
			defer server.Close()

			service, err := NewCotizacionesService(config.Config{
				CotizacionesCRUDURL: server.URL,
			})
			if err != nil {
				t.Fatalf("no se pudo crear el servicio: %v", err)
			}

			status, _, body, err := service.Request(
				context.Background(),
				http.MethodGet,
				"/v1/solicitud_cotizacion",
				nil,
				nil,
			)
			if err != nil {
				t.Fatalf("no se esperaba un error de transporte: %v", err)
			}
			if status != expectedStatus {
				t.Errorf("status = %d; se esperaba %d", status, expectedStatus)
			}
			if string(body) != `{"error":"respuesta del CRUD"}` {
				t.Errorf("cuerpo inesperado: %s", body)
			}
		})
	}
}
