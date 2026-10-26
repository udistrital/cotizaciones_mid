package main

import (
	"log"

	"github.com/beego/beego/v2/server/web"
	"github.com/udistrital/cotizaciones_mid/config"
	"github.com/udistrital/cotizaciones_mid/controllers"
	"github.com/udistrital/cotizaciones_mid/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("No se pudo cargar la configuración: %v", err)
	}

	web.Router("/health", &controllers.HealthController{})

	cotizacionesService, err := services.NewCotizacionesService(cfg)
	if err != nil {
		log.Fatalf("No se pudo inicializar el servicio de Cotizaciones CRUD: %v", err)
	}
	web.Router(
		"/solicitud/:id/detalle",
		&controllers.SolicitudController{Service: cotizacionesService},
		"get:GetDetalle",
	)
	web.BConfig.AppName = cfg.AppName
	web.BConfig.RunMode = cfg.RunMode
	web.BConfig.Listen.HTTPPort = cfg.HTTPPort

	log.Printf(
		"Iniciando %s en modo %s, puerto %d",
		cfg.AppName,
		cfg.RunMode,
		cfg.HTTPPort,
	)

	web.Run()
}
