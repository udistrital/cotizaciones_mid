package controllers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/cotizaciones_mid/services"
)

// SolicitudController expone las operaciones de consulta de solicitudes.
type SolicitudController struct {
	web.Controller
	Service *services.CotizacionesService
}

// GetDetalle consulta el detalle de una solicitud en Cotizaciones CRUD.
func (c *SolicitudController) GetDetalle() {
	id := strings.TrimSpace(c.Ctx.Input.Param(":id"))
	if id == "" {
		JSONError(&c.Controller, http.StatusBadRequest,
			"ID_REQUERIDO", "El identificador de la solicitud es obligatorio.")
		return
	}

	if c.Service == nil {
		JSONError(&c.Controller, http.StatusInternalServerError,
			"CONFIGURACION_INTERNA", "El servicio de consulta no está disponible.")
		return
	}

	path := "v1/solicitud_cotizacion/" + url.PathEscape(id) + "/detalle"
	status, headers, body, err := c.Service.Request(
		c.Ctx.Request.Context(),
		http.MethodGet,
		path,
		nil,
		map[string]string{"Accept": "application/json"},
	)
	if err != nil {
		DownstreamError(&c.Controller, err)
		return
	}

	if contentType := headers.Get("Content-Type"); contentType != "" {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", contentType)
	} else {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	c.Ctx.ResponseWriter.WriteHeader(status)
	_, _ = c.Ctx.ResponseWriter.Write(body)
}
