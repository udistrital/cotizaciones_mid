package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/udistrital/cotizaciones_mid/clients"
	"github.com/udistrital/cotizaciones_mid/config"
)

// CotizacionesService centraliza la comunicación HTTP con Cotizaciones CRUD.
type CotizacionesService struct {
	baseURL    string
	httpClient *clients.HTTPClient
}

// NewCotizacionesService crea el servicio a partir de la configuración validada.
func NewCotizacionesService(cfg config.Config) (*CotizacionesService, error) {
	parsedURL, err := url.ParseRequestURI(cfg.CotizacionesCRUDURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("la URL base de Cotizaciones CRUD no es válida")
	}

	return &CotizacionesService{
		baseURL:    strings.TrimRight(cfg.CotizacionesCRUDURL, "/"),
		httpClient: clients.NewHTTPClient(30 * time.Second),
	}, nil
}

// Request ejecuta una solicitud al CRUD usando una ruta relativa a la URL base.
// Los códigos HTTP de error se devuelven al llamador para que este los interprete.
func (s *CotizacionesService) Request(
	ctx context.Context,
	method string,
	path string,
	body []byte,
	headers map[string]string,
) (int, http.Header, []byte, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return 0, nil, nil, fmt.Errorf("la ruta del CRUD es obligatoria")
	}

	if strings.Contains(path, "://") {
		return 0, nil, nil, fmt.Errorf("la ruta debe ser relativa a la URL base del CRUD")
	}

	targetURL := s.baseURL + "/" + strings.TrimLeft(path, "/")

	return s.httpClient.Do(ctx, method, targetURL, body, headers)
}
