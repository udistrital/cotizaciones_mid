package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

// HealthController expone la comprobación básica del servicio.
type HealthController struct {
	web.Controller
}

// Get responde cuando el servicio está disponible.
func (c *HealthController) Get() {
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]string{
		"status":  "ok",
		"service": "cotizaciones_mid",
	}
	c.ServeJSON()
}
