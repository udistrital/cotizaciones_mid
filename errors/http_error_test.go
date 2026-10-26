package errors

import (
	"context"
	"net/http"
	"testing"
)

func TestDownstreamStatus(t *testing.T) {
	if got := DownstreamStatus(context.DeadlineExceeded); got != http.StatusGatewayTimeout {
		t.Errorf("DownstreamStatus(timeout) = %d; se esperaba %d", got, http.StatusGatewayTimeout)
	}

	if got := DownstreamStatus(context.Canceled); got != http.StatusBadGateway {
		t.Errorf("DownstreamStatus(canceled) = %d; se esperaba %d", got, http.StatusBadGateway)
	}

	if got := DownstreamStatus(nil); got != http.StatusBadGateway {
		t.Errorf("DownstreamStatus(nil) = %d; se esperaba %d", got, http.StatusBadGateway)
	}
}

func TestIsTimeout(t *testing.T) {
	if !IsTimeout(context.DeadlineExceeded) {
		t.Error("se esperaba que DeadlineExceeded fuera un timeout")
	}

	if IsTimeout(context.Canceled) {
		t.Error("Canceled no debe identificarse automáticamente como timeout")
	}

	if IsTimeout(nil) {
		t.Error("nil no debe identificarse como timeout")
	}
}
