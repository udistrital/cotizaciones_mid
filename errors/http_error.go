package errors

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

// DownstreamStatus traduce un fallo de comunicación con un servicio dependiente
// al código HTTP que debe devolver el MID.
func DownstreamStatus(err error) int {
	if err == nil {
		return http.StatusBadGateway
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return http.StatusGatewayTimeout
	}

	if errors.Is(err, context.Canceled) {
		// La cancelación del contexto no equivale necesariamente a un timeout.
		return http.StatusBadGateway
	}

	return http.StatusBadGateway
}

// IsTimeout permite identificar explícitamente errores de timeout.
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// TimeoutError crea un contexto con timeout para una operación dependiente.
func TimeoutError(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return context.WithTimeout(parent, timeout)
}
