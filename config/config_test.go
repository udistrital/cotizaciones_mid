package config

import (
	"os"
	"testing"
)

func setValidEnv(t *testing.T) {
	t.Helper()

	t.Setenv("COTIZACIONES_MID_APP_NAME", "cotizaciones_mid")
	t.Setenv("COTIZACIONES_MID_HTTP_PORT", "8080")
	t.Setenv("COTIZACIONES_MID_RUN_MODE", "prod")
	t.Setenv("COTIZACIONES_MID_COTIZACIONES_CRUD_URL", "http://localhost:8085")
}

func TestLoadValidConfig(t *testing.T) {
	setValidEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() devolvió un error inesperado: %v", err)
	}

	if cfg.AppName != "cotizaciones_mid" {
		t.Errorf("AppName = %q; se esperaba %q", cfg.AppName, "cotizaciones_mid")
	}
	if cfg.HTTPPort != 8080 {
		t.Errorf("HTTPPort = %d; se esperaba 8080", cfg.HTTPPort)
	}
	if cfg.RunMode != "prod" {
		t.Errorf("RunMode = %q; se esperaba %q", cfg.RunMode, "prod")
	}
	if cfg.CotizacionesCRUDURL != "http://localhost:8085" {
		t.Errorf("URL = %q; valor inesperado", cfg.CotizacionesCRUDURL)
	}
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	setValidEnv(t)
	t.Setenv("COTIZACIONES_MID_HTTP_PORT", "70000")

	_, err := Load()
	if err == nil {
		t.Fatal("se esperaba un error para un puerto inválido")
	}
}

func TestLoadRejectsInvalidRunMode(t *testing.T) {
	setValidEnv(t)
	t.Setenv("COTIZACIONES_MID_RUN_MODE", "testing")

	_, err := Load()
	if err == nil {
		t.Fatal("se esperaba un error para un modo inválido")
	}
}

func TestLoadRejectsInvalidURL(t *testing.T) {
	setValidEnv(t)
	t.Setenv("COTIZACIONES_MID_COTIZACIONES_CRUD_URL", "no-es-una-url")

	_, err := Load()
	if err == nil {
		t.Fatal("se esperaba un error para una URL inválida")
	}
}

func TestLoadRejectsMissingAppName(t *testing.T) {
	setValidEnv(t)
	t.Setenv("COTIZACIONES_MID_APP_NAME", "")

	_, err := Load()
	if err == nil {
		t.Fatal("se esperaba un error si falta el nombre de la aplicación")
	}
}

func TestLoadRejectsMissingPort(t *testing.T) {
	setValidEnv(t)
	if err := os.Unsetenv("COTIZACIONES_MID_HTTP_PORT"); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("se esperaba un error si falta el puerto")
	}
}
