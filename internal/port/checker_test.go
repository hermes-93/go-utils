package port_test

import (
	"net"
	"testing"

	"github.com/hermes-93/go-utils/internal/port"
)

func TestCheckOpen(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	result := port.Check(ln.Addr().String(), 5)
	if !result.Open {
		t.Fatalf("expected open, got: %s", result.Error)
	}
	if result.Latency == 0 {
		t.Fatal("expected non-zero latency")
	}
}

func TestCheckClosed(t *testing.T) {
	result := port.Check("127.0.0.1:1", 1)
	if result.Open {
		t.Fatal("expected closed for port 1")
	}
	if result.Error == "" {
		t.Fatal("expected error message")
	}
}

func TestResultStringOpen(t *testing.T) {
	r := port.Result{Address: "host:80", Open: true}
	s := r.String()
	if s == "" {
		t.Fatal("empty string result")
	}
}

func TestResultStringClosed(t *testing.T) {
	r := port.Result{Address: "host:80", Open: false, Error: "refused"}
	s := r.String()
	if s == "" {
		t.Fatal("empty string result")
	}
}
