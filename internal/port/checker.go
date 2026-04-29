package port

import (
	"fmt"
	"net"
	"time"
)

type Result struct {
	Address string
	Open    bool
	Latency time.Duration
	Error   string
}

func (r Result) String() string {
	if r.Open {
		return fmt.Sprintf("OPEN   %s (%dms)", r.Address, r.Latency.Milliseconds())
	}
	if r.Error != "" {
		return fmt.Sprintf("CLOSED %s — %s", r.Address, r.Error)
	}
	return fmt.Sprintf("CLOSED %s", r.Address)
}

func Check(address string, timeoutSeconds int) Result {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeoutSeconds)*time.Second)
	latency := time.Since(start)

	if err != nil {
		return Result{Address: address, Open: false, Latency: latency, Error: err.Error()}
	}
	conn.Close()
	return Result{Address: address, Open: true, Latency: latency}
}
