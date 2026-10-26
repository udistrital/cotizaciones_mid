package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Config contiene la configuración necesaria para ejecutar Cotizaciones MID.
type Config struct {
	AppName             string
	HTTPPort            int
	RunMode             string
	CotizacionesCRUDURL string
}

// Load lee y valida la configuración desde las variables de entorno.
func Load() (Config, error) {
	if err := LoadEnv(); err != nil {
		return Config{}, fmt.Errorf("error al cargar variables de entorno: %w", err)
	}
	cfg := Config{
		AppName:             strings.TrimSpace(os.Getenv("COTIZACIONES_MID_APP_NAME")),
		RunMode:             strings.TrimSpace(os.Getenv("COTIZACIONES_MID_RUN_MODE")),
		CotizacionesCRUDURL: strings.TrimRight(strings.TrimSpace(os.Getenv("COTIZACIONES_MID_COTIZACIONES_CRUD_URL")), "/"),
	}

	portValue := strings.TrimSpace(os.Getenv("COTIZACIONES_MID_HTTP_PORT"))
	if portValue == "" {
		return Config{}, fmt.Errorf("falta la variable COTIZACIONES_MID_HTTP_PORT")
	}

	port, err := strconv.Atoi(portValue)
	if err != nil || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("COTIZACIONES_MID_HTTP_PORT debe ser un puerto entre 1 y 65535")
	}
	cfg.HTTPPort = port

	if cfg.AppName == "" {
		return Config{}, fmt.Errorf("falta la variable COTIZACIONES_MID_APP_NAME")
	}

	if cfg.RunMode != "dev" && cfg.RunMode != "prod" {
		return Config{}, fmt.Errorf("COTIZACIONES_MID_RUN_MODE debe ser dev o prod")
	}

	parsedURL, err := url.ParseRequestURI(cfg.CotizacionesCRUDURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return Config{}, fmt.Errorf("COTIZACIONES_MID_COTIZACIONES_CRUD_URL debe ser una URL absoluta válida")
	}

	return cfg, nil
}
