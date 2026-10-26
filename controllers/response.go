package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/beego/beego/v2/server/web"
	miderrors "github.com/udistrital/cotizaciones_mid/errors"
)

// APIError representa el formato público de un error del MID.
type APIError struct {
	Error APIErrorDetail `json:"error"`
}

// APIErrorDetail contiene el código y el mensaje del error.
type APIErrorDetail struct {
	Codigo  string `json:"codigo"`
	Mensaje string `json:"mensaje"`
}

// WriteJSONError escribe una respuesta de error JSON sin depender de Beego.
func WriteJSONError(w http.ResponseWriter, status int, codigo string, mensaje string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(APIError{
		Error: APIErrorDetail{
			Codigo:  codigo,
			Mensaje: mensaje,
		},
	})
}

// JSONError permite utilizar el formato común desde controladores Beego.
func JSONError(c *web.Controller, status int, codigo string, mensaje string) {
	WriteJSONError(c.Ctx.ResponseWriter, status, codigo, mensaje)
}

// DownstreamError responde ante fallos de comunicación con un servicio dependiente.
func DownstreamError(c *web.Controller, err error) {
	status := miderrors.DownstreamStatus(err)

	if status == http.StatusGatewayTimeout {
		JSONError(c, status, "TIEMPO_ESPERA_SERVICIO", "El servicio dependiente no respondió a tiempo.")
		return
	}

	JSONError(c, http.StatusBadGateway, "SERVICIO_DEPENDIENTE_ERROR", "No fue posible completar la solicitud.")
}
