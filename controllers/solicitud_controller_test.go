package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/udistrital/cotizaciones_mid/config"
	"github.com/udistrital/cotizaciones_mid/services"
)

func TestSolicitudControllerDetalleRequiereID(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/solicitud//detalle", nil)
	ctx := context.NewContext()
	ctx.Reset(recorder, request)

	controller := &SolicitudController{}
	controller.Init(ctx, "SolicitudController", "GetDetalle", nil)
	controller.GetDetalle()

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; se esperaba %d", recorder.Code, http.StatusBadRequest)
	}

	var response APIError
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if response.Error.Codigo != "ID_REQUERIDO" {
		t.Errorf("código = %q; se esperaba ID_REQUERIDO", response.Error.Codigo)
	}
}

func TestSolicitudControllerDetalleSinServicio(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/solicitud/123/detalle", nil)
	ctx := context.NewContext()
	ctx.Reset(recorder, request)
	ctx.Input.SetParam(":id", "123")

	controller := &SolicitudController{}
	controller.Init(ctx, "SolicitudController", "GetDetalle", nil)
	controller.GetDetalle()

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d; se esperaba %d", recorder.Code, http.StatusInternalServerError)
	}

	var response APIError
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if response.Error.Codigo != "CONFIGURACION_INTERNA" {
		t.Errorf("código = %q; se esperaba CONFIGURACION_INTERNA", response.Error.Codigo)
	}
}

func TestSolicitudControllerDetalleExitoso(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("método = %q; se esperaba GET", r.Method)
		}
		if r.URL.Path != "/v1/solicitud_cotizacion/123/detalle" {
			t.Errorf("ruta = %q; valor inesperado", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"resultado":"detalle","id":"123"}`))
	}))
	defer server.Close()

	service, err := services.NewCotizacionesService(config.Config{
		CotizacionesCRUDURL: server.URL,
	})
	if err != nil {
		t.Fatalf("no se pudo crear el servicio: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/solicitud/123/detalle", nil)
	ctx := context.NewContext()
	ctx.Reset(recorder, request)
	ctx.Input.SetParam(":id", "123")

	controller := &SolicitudController{Service: service}
	controller.Init(ctx, "SolicitudController", "GetDetalle", nil)
	controller.GetDetalle()

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d; se esperaba 200; cuerpo: %s", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q; se esperaba application/json", got)
	}
	if got := recorder.Body.String(); got != `{"resultado":"detalle","id":"123"}` {
		t.Errorf("cuerpo inesperado: %s", got)
	}
}

func TestSolicitudControllerDetalleErrorDelCRUD(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"no encontrado"}`, http.StatusNotFound)
	}))
	defer server.Close()

	service, err := services.NewCotizacionesService(config.Config{
		CotizacionesCRUDURL: server.URL,
	})
	if err != nil {
		t.Fatalf("no se pudo crear el servicio: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/solicitud/123/detalle", nil)
	ctx := context.NewContext()
	ctx.Reset(recorder, request)
	ctx.Input.SetParam(":id", "123")

	controller := &SolicitudController{Service: service}
	controller.Init(ctx, "SolicitudController", "GetDetalle", nil)
	controller.GetDetalle()

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d; se esperaba 404", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "no encontrado") {
		t.Errorf("no se conservó el cuerpo del CRUD: %s", recorder.Body.String())
	}
}
