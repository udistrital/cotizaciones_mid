package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv carga las variables del archivo .env durante el desarrollo.
// Si el archivo no existe, conserva las variables definidas en el entorno.
func LoadEnv() error {
	// En producción, las variables deben venir del entorno del servidor.
	if os.Getenv("COTIZACIONES_MID_RUN_MODE") == "prod" {
		return nil
	}

	// En desarrollo, cargar .env si existe. Las variables ya definidas
	// en el entorno tienen prioridad sobre las del archivo.
	if err := godotenv.Load(); err != nil {
		if _, statErr := os.Stat(".env"); statErr == nil {
			return err
		}
	}

	return nil
}
