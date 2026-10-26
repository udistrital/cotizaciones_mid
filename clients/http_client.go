package clients

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPClient centraliza las solicitudes HTTP a servicios externos.
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient crea un cliente HTTP con un tiempo máximo de espera.
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Do ejecuta una solicitud HTTP y devuelve el código, las cabeceras y el cuerpo.
// El llamador es responsable de interpretar los códigos HTTP no exitosos.
func (c *HTTPClient) Do(
	ctx context.Context,
	method string,
	url string,
	body []byte,
	headers map[string]string,
) (int, http.Header, []byte, error) {
	method = strings.TrimSpace(method)
	url = strings.TrimSpace(url)

	if method == "" {
		return 0, nil, nil, fmt.Errorf("el método HTTP es obligatorio")
	}
	if url == "" {
		return 0, nil, nil, fmt.Errorf("la URL es obligatoria")
	}

	var requestBody io.Reader
	if body != nil {
		requestBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("crear solicitud HTTP: %w", err)
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("ejecutar solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, resp.Header, nil,
			fmt.Errorf("leer respuesta HTTP: %w", err)
	}

	return resp.StatusCode, resp.Header, responseBody, nil
}
